package metricsservice

import (
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"io"
	"net/http"
)

type Client struct {
	client  *http.Client
	address string
}

func NewClient(address string) *Client {
	return &Client{
		client:  retryablehttp.NewClient().StandardClient(),
		address: address,
	}
}

func doRequest(client *http.Client, method string, url string, body io.Reader, contentType string) error {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return fmt.Errorf("failed to create http request: %w", err)
	}

	req.Header.Set("Content-Type", contentType)
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to do http request: %w", err)
	}
	defer resp.Body.Close()

	_, _ = io.Copy(io.Discard, resp.Body)

	return nil
}
