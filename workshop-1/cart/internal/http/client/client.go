package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// client HTTP клиент для взаимодействия с API
type client struct {
	baseURL string
	client  *http.Client
}

// newClient создает новый экземпляр клиента
func newClient(baseURL string, requestTimeout time.Duration) client {
	return client{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: requestTimeout,
		},
	}
}

// doRequest выполняет POST запрос и обрабатывает ответ
func (c client) doRequest(ctx context.Context, endpoint string, req, resp any) (status int, _ error) {

	// Подготовка тела запроса
	reqBody, err := json.Marshal(req)
	if err != nil {
		return 0, fmt.Errorf("marshal request: %w", err)
	}

	// Создание запроса
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+endpoint, bytes.NewReader(reqBody))
	if err != nil {
		return 0, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	// Выполнение запроса
	httpResp, err := c.client.Do(httpReq)
	if err != nil {
		return 0, fmt.Errorf("execute request: %w", err)
	}
	defer httpResp.Body.Close()

	status = httpResp.StatusCode

	// Чтение ответа
	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return status, fmt.Errorf("read body failed: %w", err)
	}

	// Декодирование ответа
	if status == 200 || status == 201 {
		if err := json.Unmarshal(respBody, resp); err != nil {
			return status, fmt.Errorf("unmarshal response: %w", err)
		}
	}

	return status, nil
}
