package service

import (
	"net/http"
	"search-job/internal/models"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (s *Service) CreateResource(c echo.Context) error {
	var res models.Resource
	if err := c.Bind(&res); err != nil {
		return c.JSON(s.NewError(InvalidParams))
	}

	if err := s.bookingRepo.CreateResource(c.Request().Context(), &res); err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InternalServerError))
	}
	return c.JSON(http.StatusCreated, res)
}

func (s *Service) ListResources(c echo.Context) error {
	list, err := s.bookingRepo.GetResources(c.Request().Context())
	if err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InternalServerError))
	}
	return c.JSON(http.StatusOK, Response{Object: list})
}

func (s *Service) UpdateResource(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var res models.Resource
	if err := c.Bind(&res); err != nil {
		return c.JSON(s.NewError(InvalidParams))
	}
	res.ID = id

	if err := s.bookingRepo.UpdateResource(c.Request().Context(), &res); err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InternalServerError))
	}
	return c.String(http.StatusOK, "Ok")
}

func (s *Service) DeleteResource(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := s.bookingRepo.DeleteResource(c.Request().Context(), id); err != nil {
		s.logger.Error(err)
		return c.JSON(s.NewError(InternalServerError))
	}
	return c.String(http.StatusOK, "Ok")
}
