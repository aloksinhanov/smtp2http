package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

func CreateTransactionEmail(msg EmailMessage, apiKey string) error {

	var (
		err  error
		resp *resty.Response
	)

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panicked %v", r)
		}
	}()

	resp, err = resty.New().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("X-MW-PUBLIC-KEY", apiKey).
		SetHeader("X-MW-TIMESTAMP", fmt.Sprintf("%v", (time.Now().Unix()))).
		SetBody(createRequestBody(msg)).Post(*flagWebhook)
	log.Println(resp)
	if resp.StatusCode() != 201 {
		return errors.New("E2: Cannot accept your message due to internal error, please report that to our engineers")
	}
	return err
}

func getName(address string) string {
	addr := strings.TrimSpace(address)
	parts := strings.Split(addr, "@")
	return parts[0]
}

func createRequestBody(msg EmailMessage) string {
	toName := msg.Addresses.To.Name
	if toName == "" {
		toName = getName(msg.Addresses.To.Address)
	}

	fromName := msg.Addresses.From.Name
	if fromName == "" {
		fromName = getName(msg.Addresses.From.Address)
	}

	return fmt.Sprintf(`email[to_name]=%v&email[to_email]=%v&email[from_name]=%v
	&email[from_email]=%v&email[subject]=%v&email[body]=%v&send_at=%v`,
		toName, msg.Addresses.To.Address, fromName, msg.Addresses.From.Address,
		msg.Subject, base64.StdEncoding.EncodeToString([]byte(msg.Body.Text)), time.Now().UTC().Format("2006-01-02 15:04:05"))
	//need to check what to do with text or html
}
