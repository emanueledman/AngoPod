package db

import (
	"angopod-back/internal/models"
)

// GetProvinces retorna todas as províncias registradas com o total atualizado
func GetProvinces() ([]models.Province, error) {
	query := `
		SELECT p.name,
		       COALESCE((SELECT COUNT(*) FROM podcasts WHERE LOWER(province) = LOWER(p.name)), p.count) as count,
		       p.x, p.y, p.description
		FROM provinces p
		ORDER BY count DESC
	`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Province
	for rows.Next() {
		var prov models.Province
		if err := rows.Scan(&prov.Name, &prov.Count, &prov.X, &prov.Y, &prov.Description); err != nil {
			return nil, err
		}
		result = append(result, prov)
	}
	return result, nil
}

