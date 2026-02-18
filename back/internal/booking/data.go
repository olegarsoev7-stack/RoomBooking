package booking

import (
	"context"

	"search-job/internal/models"
	"search-job/pkg/postgres"
)

type Repo struct {
	db *postgres.DB
}

func NewRepo(db *postgres.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) RGetBookingById(ctx context.Context, id int) (*models.Booking, error) {
	var b models.Booking
	err := r.db.QueryRow(ctx, `
		SELECT id, user_id, resource_id, title, start_at, end_at, is_holiday, created_at, updated_at, deleted_at
		FROM bookings
		WHERE id = $1
	`, id).Scan(
		&b.ID,
		&b.UserID,
		&b.ResourceID,
		&b.Title,
		&b.StartAt,
		&b.EndAt,
		&b.IsHoliday,
		&b.CreatedAt,
		&b.UpdatedAt,
		&b.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *Repo) RCreateBooking(ctx context.Context, b *models.Booking) error {
	err := r.db.QueryRow(ctx, `
		INSERT INTO bookings (user_id, resource_id, title, start_at, end_at, is_holiday)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at, deleted_at
	`,
		b.UserID,
		b.ResourceID,
		b.Title,
		b.StartAt,
		b.EndAt,
		b.IsHoliday,
	).Scan(
		&b.ID,
		&b.CreatedAt,
		&b.UpdatedAt,
		&b.DeletedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) RListBookings(ctx context.Context, page, limit int) ([]models.Booking, error) {
	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	offset := (page - 1) * limit

	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, resource_id, title, start_at, end_at, is_holiday, created_at, updated_at, deleted_at
		FROM bookings
		ORDER BY id DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []models.Booking
	for rows.Next() {
		var b models.Booking
		if err := rows.Scan(
			&b.ID,
			&b.UserID,
			&b.ResourceID,
			&b.Title,
			&b.StartAt,
			&b.EndAt,
			&b.IsHoliday,
			&b.CreatedAt,
			&b.UpdatedAt,
			&b.DeletedAt,
		); err != nil {
			return nil, err
		}
		res = append(res, b)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (r *Repo) RUpdateBooking(ctx context.Context, b *models.Booking) error {
	err := r.db.QueryRow(ctx, `
		UPDATE bookings
		SET user_id = $1,
		    resource_id = $2,
		    title = $3,
		    start_at = $4,
		    end_at = $5,
		    is_holiday = $6,
		    updated_at = NOW()
		WHERE id = $7
		RETURNING created_at, updated_at, deleted_at
	`,
		b.UserID,
		b.ResourceID,
		b.Title,
		b.StartAt,
		b.EndAt,
		b.IsHoliday,
		b.ID,
	).Scan(
		&b.CreatedAt,
		&b.UpdatedAt,
		&b.DeletedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) RDeleteBooking(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `
		DELETE FROM bookings
		WHERE id = $1
	`, id)
	if err != nil {
		return err
	}
	return nil
}
