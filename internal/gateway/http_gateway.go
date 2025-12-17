package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	domain "github.com/josuesantos1/emv/internal/domain"
)

type HTTPGateway struct {
	BaseURL string
	Client  *http.Client
}

type authorizationRequest struct {
	Pan          string    `json:"pan"`
	DataValidade time.Time `json:"data_validade"`
	CVM          string    `json:"cvm"`
}

type authorizationResponse struct {
	Approved  bool      `json:"approved"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

func NewHTTPGateway(baseURL string) *HTTPGateway {
	return &HTTPGateway{
		BaseURL: baseURL,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (g *HTTPGateway) Authorize(transaction *domain.Tlv) (bool, error) {
	req := authorizationRequest{
		Pan:          transaction.Pan,
		DataValidade: transaction.DataValidade,
		CVM:          transaction.CVM,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return false, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := g.Client.Post(
		g.BaseURL+"/authorize",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return false, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var authResp authorizationResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return false, fmt.Errorf("failed to decode response: %w", err)
	}

	return authResp.Approved, nil
}
