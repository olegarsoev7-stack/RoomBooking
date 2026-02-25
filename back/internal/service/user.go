package service

import (
	"net/http"

	"search-job/internal/auth"
	"search-job/internal/models"
	"search-job/pkg/jwt"

	"github.com/labstack/echo/v4"
)

func (s *Service) Register(c echo.Context) error {
	var u models.User
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	// Шифруем пароль
	hashed, _ := auth.HashPassword(u.Password)
	u.Password = hashed

	// Сохраняем в базу
	if err := s.bookingRepo.CreateUser(c.Request().Context(), &u); err != nil {
		s.logger.Error(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not create user"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "User registered successfully"})
}

func (s *Service) Login(c echo.Context) error {
	var input models.User
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	// 1. Ищем пользователя по Email
	user, err := s.bookingRepo.GetUserByEmail(c.Request().Context(), input.Email)
	if err != nil {
		s.logger.Error(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Server error"})
	}
	if user == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not found"})
	}

	// 2. Проверяем пароль
	if !auth.CheckPassword(input.Password, user.Password) {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid password"})
	}

	// 3. Генерируем токен, подставляя реальный ID пользователя!
	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		s.logger.Error(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not generate token"})
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})
}
