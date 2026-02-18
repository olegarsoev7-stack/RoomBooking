package service

import (
	"net/http"
	"search-job/internal/models"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (s *Service) CreateBooking(c echo.Context) error {
	var booking models.Booking
	if err := c.Bind(&booking); err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InvalidParams))
	}

	repo := s.bookingRepo
	if err := repo.RCreateBooking(c.Request().Context(), &booking); err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InternalServerError))
	}

	return c.String(http.StatusOK, "Ok")
}

func (s *Service) GetBooking(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InvalidParams))
	}

	repo := s.bookingRepo
	booking, err := repo.RGetBookingById(c.Request().Context(), id)
	if err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InternalServerError))
	}

	return c.JSON(http.StatusOK, Response{Object: booking})
}

func (s *Service) ListBookings(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page == 0 {
		page = 1
	}
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit == 0 {
		limit = 20
	}

	repo := s.bookingRepo
	list, err := repo.RListBookings(c.Request().Context(), page, limit)
	if err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InternalServerError))
	}

	return c.JSON(http.StatusOK, Response{Object: list})
}

func (s *Service) UpdateBooking(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InvalidParams))
	}

	var booking models.Booking
	if err := c.Bind(&booking); err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InvalidParams))
	}
	booking.ID = id

	repo := s.bookingRepo
	if err := repo.RUpdateBooking(c.Request().Context(), &booking); err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InternalServerError))
	}

	return c.String(http.StatusOK, "Ok")
}

func (s *Service) DeleteBooking(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InvalidParams))
	}

	repo := s.bookingRepo
	if err := repo.RDeleteBooking(c.Request().Context(), id); err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InternalServerError))
	}

	return c.String(http.StatusOK, "Ok")
}
