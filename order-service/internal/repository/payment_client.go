package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type paymentClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewPaymentClient(baseURL string) *paymentClient {
	return &paymentClient{
		httpClient: &http.Client{
			Timeout: 2 * time.Second,
		},
		baseURL: baseURL,
	}
}

func (c *paymentClient) AuthorizePayment(ctx context.Context, orderID string, amount int64) (string, string, error) {
	url := fmt.Sprintf("%s/payments", c.baseURL)

	body, _ := json.Marshal(map[string]interface{}{
		"order_id": orderID,
		"amount":   amount,
	})

	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("payment service unavailable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("payment service returned error: %d", resp.StatusCode)
	}

	var result struct {
		TransactionID string `json:"transaction_id"`
		Status        string `json:"status"`
	}
	json.NewDecoder(resp.Body).Decode(&result)

	return result.TransactionID, result.Status, nil
}
