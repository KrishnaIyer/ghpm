// Copyright © 2023 Krishna Iyer Easwaran
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

const httpTimeout = 10 * time.Second

// Config is the client configuration.
type Config struct {
	APIKey string `name:"api-key" description:"The GitHub API key"`
}

// Client is a GitHub client.
type Client struct {
	config *Config
	client *http.Client
}

// NewClient creates a new client.
func NewClient(config *Config) *Client {
	return &Client{
		config: config,
		client: &http.Client{
			Timeout: httpTimeout,
		},
	}
}

// Do executes a request.
func (c *Client) Do(ctx context.Context, method string, url string, body []byte) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "ghpm")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch {
	case resp.StatusCode >= 200 && resp.StatusCode < 300:
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return body, nil
	case resp.StatusCode == http.StatusNotFound:
		return nil, fmt.Errorf("not found")
	default:
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
}
