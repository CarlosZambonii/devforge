package mlclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	baseURL string
	http    *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		http:    &http.Client{Timeout: 3 * time.Second},
	}
}

type predictRequest struct {
	URL string `json:"url"`
}

type PredictResult struct {
	URL   string  `json:"url"`
	Risk  string  `json:"risk"`
	Score float64 `json:"score"`
}

// Predict consulta o ml-service. Se falhar, retorna risco "unknown"
// (degrada com elegancia - nao quebra o encurtador se a IA estiver fora).
func (c *Client) Predict(ctx context.Context, url string) PredictResult {
	body, _ := json.Marshal(predictRequest{URL: url})

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.baseURL+"/predict", bytes.NewReader(body))
	if err != nil {
		return PredictResult{URL: url, Risk: "unknown", Score: 0}
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return PredictResult{URL: url, Risk: "unknown", Score: 0}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return PredictResult{URL: url, Risk: "unknown", Score: 0}
	}

	var result PredictResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return PredictResult{URL: url, Risk: "unknown", Score: 0}
	}
	return result
}

var _ = fmt.Sprintf
