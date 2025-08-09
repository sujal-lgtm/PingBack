package services

import (
	"bytes"
	"net/http"
	"fmt"
)

func ForwardEvent(payload string, targetURL string) error {
	resp , err := http.Post(targetURL, "application/json",bytes.NewBuffer([]byte(payload)))

	if err != nil{
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		return fmt.Errorf("failed to forward event, status code: %d", resp.StatusCode)
	}

	return nil
}