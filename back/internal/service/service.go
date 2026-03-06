package service

import (
	"search-job/internal/booking"
	"search-job/pkg/holiday"
	"search-job/pkg/postgres"

	"github.com/labstack/echo/v4"
)

const (
	InvalidParams       = "invalid params"
	InternalServerError = "internal error"
)

type Service struct {
	db            *postgres.DB
	logger        echo.Logger
	bookingRepo   *booking.Repo
	holidayClient *holiday.Client
}

func NewService(db *postgres.DB, logger echo.Logger, hc *holiday.Client) *Service {
	return &Service{
		db:            db,
		logger:        logger,
		bookingRepo:   booking.NewRepo(db),
		holidayClient: hc, // <-- Сохраняем в структуру
	}
}

func (s *Service) initRepositories(db *postgres.DB) {
	s.bookingRepo = booking.NewRepo(db)
}

// Пока можно не вдаваться в то что ниже

type Response struct {
	Object       any    `json:"object,omitempty"`
	ErrorMessage string `json:"error,omitempty"`
}

func (r *Response) Error() string {
	return r.ErrorMessage
}

func (s *Service) NewError(err string) (int, *Response) {
	return 400, &Response{ErrorMessage: err}
}
