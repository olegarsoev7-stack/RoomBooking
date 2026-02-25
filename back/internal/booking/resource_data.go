package booking

import (
	"context"
	"search-job/internal/models"
)

// CreateResource — сохраняет новую переговорку
func (r *Repo) CreateResource(ctx context.Context, res *models.Resource) error {
	query := `
		INSERT INTO resources (name, capacity, created_at)
		VALUES ($1, $2, NOW())
		RETURNING id, created_at`

	return r.db.QueryRow(ctx, query, res.Name, res.Capacity).Scan(&res.ID, &res.CreatedAt)
}

// GetResources — возвращает список всех ресурсов
func (r *Repo) GetResources(ctx context.Context) ([]models.Resource, error) {
	query := `SELECT id, name, capacity, created_at FROM resources ORDER BY id ASC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resources []models.Resource
	for rows.Next() {
		var res models.Resource
		if err := rows.Scan(&res.ID, &res.Name, &res.Capacity, &res.CreatedAt); err != nil {
			return nil, err
		}
		resources = append(resources, res)
	}
	return resources, nil
}

// UpdateResource — обновляет данные (название или вместимость)
func (r *Repo) UpdateResource(ctx context.Context, res *models.Resource) error {
	query := `UPDATE resources SET name = $1, capacity = $2 WHERE id = $3`
	_, err := r.db.Exec(ctx, query, res.Name, res.Capacity, res.ID)
	return err
}

// DeleteResource — удаляет ресурс из базы
func (r *Repo) DeleteResource(ctx context.Context, id int) error {
	query := `DELETE FROM resources WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
