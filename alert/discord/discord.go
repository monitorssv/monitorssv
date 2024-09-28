package discord

import (
	"bytes"
	"encoding/json"
	"golang.org/x/xerrors"
	"net/http"
)

type Client struct {
	webhook string
}

func NewDiscordClient(webhook string) *Client {
	return &Client{
		webhook: webhook,
	}
}

func (*Client) Platform() string {
	return "discord"
}

func (c *Client) Send(msg string) error {
	payload := map[string]string{
		"content": msg,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return xerrors.Errorf("failed to marshal JSON: %w", err)
	}

	req, err := http.NewRequest("POST", c.webhook, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		return nil
	}

	return xerrors.Errorf("%s", resp.Status)
}
