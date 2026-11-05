# AngoPod — Plataforma de Podcasts Angolanos

Plataforma moderna para descoberta, ranking e reprodução de podcasts angolanos de todas as 18 províncias.

---

## 🏗️ Arquitetura do Projeto

```
AngoPod/
├── Front/          # Single Page Application (React 19 + TypeScript + Vite + Tailwind CSS v4)
└── Back/           # API REST em Go 1.24 (net/http nativo com roteamento Go 1.22+)
```

---

## 🚀 Como Executar

### 1. Back-end (Go)
```bash
cd Back
go run ./cmd/api
```
- A API iniciará em: **http://localhost:8080**
- Health Check: `http://localhost:8080/api/health`
- Endpoints REST disponíveis:
  - `GET /api/podcasts` (suporta `?category=...&province=...&search=...`)
  - `GET /api/podcasts/{slug}`
  - `GET /api/podcasts/{slug}/episodes`
  - `GET /api/episodes`
  - `GET /api/episodes/{id}`
  - `GET /api/presenters`
  - `GET /api/presenters/{slug}`
  - `GET /api/provinces`
  - `GET /api/ranking`
  - `GET /api/highlights`

### 2. Front-end (React + Vite)
```bash
cd Front
npm install
npm run dev
```
- A aplicação web iniciará em: **http://localhost:5174** (ou porta indicada no terminal)
- Comunicação com a API Go via `VITE_API_URL=http://localhost:8080/api` (com fallback automático).

