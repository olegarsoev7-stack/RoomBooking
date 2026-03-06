package booking

import (
	"context"
	"time"

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
		UPDATE bookings 
		SET deleted_at = NOW(), updated_at = NOW() 
		WHERE id = $1 AND deleted_at IS NULL
	`, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) HasOverlap(ctx context.Context, resourceID int, startAt, endAt time.Time) (bool, error) {
	var count int
	// Ищем активные (не удаленные) бронирования, которые пересекаются с новым
	query := `
		SELECT COUNT(*) FROM bookings 
		WHERE resource_id = $1 
		  AND deleted_at IS NULL 
		  AND start_at < $3 
		  AND end_at > $2`

	err := r.db.QueryRow(ctx, query, resourceID, startAt, endAt).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repo) RRestoreBooking(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `
		UPDATE bookings 
		SET deleted_at = NULL, updated_at = NOW() 
		WHERE id = $1 AND deleted_at IS NOT NULL
	`, id)
	if err != nil {
		return err
	}
	return nil
}

// Получение расписания конкретного ресурса за период
func (r *Repo) RGetResourceSchedule(ctx context.Context, resourceID int, from, to time.Time) ([]models.Booking, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, resource_id, title, start_at, end_at, is_holiday, created_at, updated_at, deleted_at
		FROM bookings 
		WHERE resource_id = $1 
		  AND deleted_at IS NULL
		  AND start_at >= $2 
		  AND start_at <= $3 
		ORDER BY start_at ASC
	`, resourceID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []models.Booking
	for rows.Next() {
		var b models.Booking
		if err := rows.Scan(
			&b.ID, &b.UserID, &b.ResourceID, &b.Title,
			&b.StartAt, &b.EndAt, &b.IsHoliday,
			&b.CreatedAt, &b.UpdatedAt, &b.DeletedAt,
		); err != nil {
			return nil, err
		}
		res = append(res, b)
	}
	return res, rows.Err()
}

// Подсчет количества бронирований для сводки
func (r *Repo) RGetBookingsCount(ctx context.Context, resourceID int, from, to time.Time) (int, error) {
	var count int
	query := `
		SELECT COUNT(*) 
		FROM bookings 
		WHERE resource_id = $1 
		  AND deleted_at IS NULL
		  AND start_at >= $2 
		  AND start_at <= $3`

	err := r.db.QueryRow(ctx, query, resourceID, from, to).Scan(&count)
	return count, err
}
