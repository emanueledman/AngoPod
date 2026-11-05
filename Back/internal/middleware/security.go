package middleware

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// SecurityHeaders adiciona cabeçalhos de proteção avançada (OWASP)
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Proteção contra MIME-sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")
		// Proteção contra Clickjacking
		w.Header().Set("X-Frame-Options", "DENY")
		// Filtro XSS para navegadores legados
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		// Política de referência segura
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		// Content Security Policy permissivo para imagens e mídia externas (YouTube / CDN)
		w.Header().Set("Content-Security-Policy", "default-src 'self' http: https: data:; img-src * data: blob:; media-src * data: blob:; frame-ancestors 'none'")
		// Restrição de APIs sensíveis do navegador
		w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=(), payment=()")
		// HSTS (HTTP Strict Transport Security)
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// Remove ou oculta informação do servidor
		w.Header().Del("Server")
		w.Header().Del("X-Powered-By")

		next.ServeHTTP(w, r)
	})
}

// CORSStrict valida a origem da requisição contra uma whitelist rigorosa
func CORSStrict(allowedOrigins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			isAllowed := false

			// Em ambiente local, qualquer porta de localhost ou 127.0.0.1 é autorizada
			if origin != "" && isLocalhost(origin) {
				isAllowed = true
			} else if origin != "" {
				for _, allowed := range allowedOrigins {
					if allowed == "*" || origin == allowed {
						isAllowed = true
						break
					}
				}
			}

			if isAllowed && origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Accept")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Max-Age", "86400") // 24h cache preflight
				w.Header().Set("Vary", "Origin")
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func isLocalhost(origin string) bool {
	return strings.HasPrefix(origin, "http://localhost:") ||
		strings.HasPrefix(origin, "http://127.0.0.1:") ||
		strings.HasPrefix(origin, "https://localhost:")
}

// ipRateLimiter mantém controlo do número de requisições por IP
type ipRateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

func newIPRateLimiter(limit int, window time.Duration) *ipRateLimiter {
	limiter := &ipRateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}

	// Rotina em segundo plano para limpar IPs inativos e evitar fuga de memória
	go func() {
		for {
			time.Sleep(limiter.window)
			limiter.mu.Lock()
			now := time.Now()
			for ip, times := range limiter.requests {
				var valid []time.Time
				for _, t := range times {
					if now.Sub(t) < limiter.window {
						valid = append(valid, t)
					}
				}
				if len(valid) == 0 {
					delete(limiter.requests, ip)
				} else {
					limiter.requests[ip] = valid
				}
			}
			limiter.mu.Unlock()
		}
	}()

	return limiter
}

// RateLimit limita a taxa de requisições por IP (DDoS/Brute-force protection)
func RateLimit(limit int, window time.Duration) func(http.Handler) http.Handler {
	limiter := newIPRateLimiter(limit, window)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := getClientIP(r)

			limiter.mu.Lock()
			now := time.Now()
			times := limiter.requests[ip]

			var validTimes []time.Time
			for _, t := range times {
				if now.Sub(t) < limiter.window {
					validTimes = append(validTimes, t)
				}
			}

			if len(validTimes) >= limiter.limit {
				limiter.mu.Unlock()
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Retry-After", "60")
				w.WriteHeader(http.StatusTooManyRequests)
				_ = json.NewEncoder(w).Encode(map[string]any{
					"error":   "Demasiadas requisições. Por favor, tente novamente mais tarde.",
					"status":  http.StatusTooManyRequests,
					"limite":  limiter.limit,
					"janela":  limiter.window.String(),
				})
				return
			}

			limiter.requests[ip] = append(validTimes, now)
			limiter.mu.Unlock()

			next.ServeHTTP(w, r)
		})
	}
}

// getClientIP extrai o IP real do cliente mesmo através de proxies seguros
func getClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// RecoverPanic captura falhas de execução e previne crash do servidor
func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("🚨 [PANIC RECUPERADO] %s %s: %v", r.Method, r.URL.Path, rec)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(map[string]string{
					"error": "Ocorreu um erro interno no servidor.",
				})
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// SanitizeRequest valida comprimentos e impede injeção/path traversal
func SanitizeRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Proteção contra URLs excessivamente longas (ataques de buffer)
		if len(r.URL.Path) > 1024 || len(r.URL.RawQuery) > 2048 {
			http.Error(w, `{"error":"URL demasiado longa"}`, http.StatusBadRequest)
			return
		}

		// Bloqueia tentativas de Path Traversal ou injeção de bytes nulos
		if strings.Contains(r.URL.Path, "..") || strings.Contains(r.URL.Path, "\x00") {
			http.Error(w, `{"error":"Caminho de URL inválido"}`, http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// LimitBodySize limita o tamanho do corpo da requisição (ex: 1MB)
func LimitBodySize(maxBytes int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Body != nil {
				r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			}
			next.ServeHTTP(w, r)
		})
	}
}

