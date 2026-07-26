package gateforge

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GetJWKS fetches the JSON Web Key Set from /.well-known/jwks.json.
func (c *Client) GetJWKS(ctx context.Context) (map[string]any, *http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/.well-known/jwks.json", nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("User-Agent", c.userAgent)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, resp, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, resp, NewAPIError(resp.StatusCode, Meta{Message: string(body)}, resp)
	}
	var out map[string]any
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, resp, fmt.Errorf("jwks: decode: %w", err)
	}
	return out, resp, nil
}
