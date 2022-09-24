package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/url"
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

	req := resty.New().R().SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("X-MW-PUBLIC-KEY", apiKey).
		SetHeader("X-MW-TIMESTAMP", fmt.Sprintf("%v", (time.Now().Unix()))).
		SetBody(createRequestBody(msg))

	// useful for capturing CURls to debug
	// log.Println(http2curl.GetCurlCommand(req.RawRequest))

	resp, err = req.Post(*flagWebhook)

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

	data := url.Values{}
	data.Set("email[to_name]", toName)
	data.Set("email[to_email]", msg.Addresses.To.Address)
	data.Set("email[from_name]", fromName)
	data.Set("email[from_email]", msg.Addresses.From.Address)
	data.Set("email[subject]", msg.Subject)
	//always reading HTML. Is there a use-case to read Text?
	data.Set("email[body]", base64.StdEncoding.EncodeToString([]byte(msg.Body.HTML)))
	data.Set("send_at", time.Now().UTC().Format("2006-01-02 15:04:05"))

	return data.Encode()
}
