package page

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"golang.org/x/time/rate"
)

type NamuClient struct {
	clnt *http.Client
	lmtr *rate.Limiter
	mu   sync.Mutex
}

func NewNamuClient() *NamuClient {
	return &NamuClient{clnt: &http.Client{},
		lmtr: rate.NewLimiter(rate.Limit(10), 1)}
}

func (namu *NamuClient) GetHtml(url string) (*string, error) {
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("🛑error creating request: %w", err)
	}
	return namu.readHtml(req)
}

func (namu *NamuClient) readHtml(req *http.Request) (*string, error) {
	for {
		if !namu.lmtr.Allow() {
			continue
		}
		resp, err := namu.clnt.Do(req)
		if err != nil {
			return nil, fmt.Errorf("🛑error sending request: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("🛑error reading response body: %w", err)
		}
		bodyString := string(body)

		if namu.detectBlock(bodyString) {
			return nil, fmt.Errorf("🛑blocking is detected: %w", err)
		}
		return &bodyString, nil
	}

}

func (namu *NamuClient) detectBlock(bodyString string) bool {
	return strings.Contains(bodyString, "<title>비정상")
}
