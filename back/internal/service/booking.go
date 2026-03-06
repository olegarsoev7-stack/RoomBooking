package service

import (
	"net/http"
	"search-job/internal/models"
	"strconv"
	"time"

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

func (s *Service) RestoreBooking(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		// Убрали http.StatusBadRequest
		return c.JSON(s.NewError("InvalidParams"))
	}

	if err := s.bookingRepo.RRestoreBooking(c.Request().Context(), id); err != nil {
		s.logger.Error(err)
		// Убрали http.StatusInternalServerError
		return c.JSON(s.NewError("InternalServerError"))
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "booking restored"})
}

func (s *Service) GetResourceSchedule(c echo.Context) error {
	resourceID, _ := strconv.Atoi(c.Param("id"))
	from, _ := time.Parse(time.RFC3339, c.QueryParam("from"))
	to, _ := time.Parse(time.RFC3339, c.QueryParam("to"))

	schedule, err := s.bookingRepo.RGetResourceSchedule(c.Request().Context(), resourceID, from, to)
	if err != nil {
		s.logger.Error(err)
		// Убрали http.StatusInternalServerError
		return c.JSON(s.NewError("InternalServerError"))
	}

	return c.JSON(http.StatusOK, schedule)
}

func (s *Service) GetBookingsSummary(c echo.Context) error {
	resourceID, _ := strconv.Atoi(c.QueryParam("resourceId"))
	from, _ := time.Parse(time.RFC3339, c.QueryParam("from"))
	to, _ := time.Parse(time.RFC3339, c.QueryParam("to"))

	count, err := s.bookingRepo.RGetBookingsCount(c.Request().Context(), resourceID, from, to)
	if err != nil {
		s.logger.Error(err)
		// Убрали http.StatusInternalServerError
		return c.JSON(s.NewError("InternalServerError"))
	}

	return c.JSON(http.StatusOK, map[string]int{"count": count})
}
