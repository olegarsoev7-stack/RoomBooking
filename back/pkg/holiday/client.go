package holiday

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client - структура нашего HTTP-клиента для проверки праздников
type Client struct {
	httpClient *http.Client
}

// NewClient создает новый клиент с таймаутом, чтобы внешний API не повесил наш сервис
func NewClient(timeout time.Duration) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: timeout},
	}
}

// CheckDay ходит во внешнее API и проверяет, является ли дата праздником/выходным
func (c *Client) CheckDay(ctx context.Context, date time.Time) (bool, error) {
	// isdayoff.ru принимает дату в формате YYYYMMDD или отдельными параметрами
	url := fmt.Sprintf("https://isdayoff.ru/api/getdata?year=%d&month=%02d&day=%02d",
		date.Year(), date.Month(), date.Day())

	// Создаем запрос с контекстом (позволяет отменить запрос, если клиент отвалился)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, err
	}

	// Выполняем запрос
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close() // Не забываем закрывать тело ответа!

	// Читаем ответ
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	// isdayoff.ru возвращает "1", если это нерабочий день (выходной или праздник)
	if string(body) == "1" {
		return true, nil
	}

	return false, nil
}
