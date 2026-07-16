// Package client
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/henrywhitaker3/maas/pkg/locker"
	"github.com/henrywhitaker3/windowframe/duration"
	wuuid "github.com/henrywhitaker3/windowframe/uuid"
)

var (
	url = "https://maas.henrywhitaker.com"
)

type Client struct{}

func New() *Client {
	return &Client{}
}

func (c *Client) Lock(
	ctx context.Context,
	subject string,
	owner uuid.UUID,
	dur time.Duration,
) error {
	by, err := json.Marshal(locker.LockRequest{
		Subject:  subject,
		Owner:    wuuid.UUID(owner),
		Duration: duration.StringDuration(dur),
	})
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}
	reader := bytes.NewReader(by)
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v1/lock", url), reader)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("lock subject: %w", err)
	}
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("lock request: %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) Unlock(
	ctx context.Context,
	subject string,
	owner uuid.UUID,
) error {
	by, err := json.Marshal(locker.UnlockRequest{
		Subject: subject,
		Owner:   wuuid.UUID(owner),
	})
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}
	reader := bytes.NewReader(by)
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v1/unlock", url), reader)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("unlock subject: %w", err)
	}
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unlock request: %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) Refresh(
	ctx context.Context,
	subject string,
	owner uuid.UUID,
	dur time.Duration,
) error {
	by, err := json.Marshal(locker.RefreshRequest{
		Subject:  subject,
		Owner:    wuuid.UUID(owner),
		Duration: duration.StringDuration(dur),
	})
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}
	reader := bytes.NewReader(by)
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v1/refresh", url), reader)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("refresh subject: %w", err)
	}
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("refresh request: %d", resp.StatusCode)
	}
	return nil
}
