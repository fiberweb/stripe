package stripe

import (
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/gofiber/fiber"
)

var (
	invalidSignature = "t=1583797522,v1=c41417a31dd074a324454ca17c27b698d65c9e1251a1ad06b282610dbfb95358,v0=115541604a038bdfed1ce950591d31af21d644828ae576009585d0501dbd3c6b"
	webhookBody      = `{
		"created": 1326853478,
		"livemode": false,
		"id": "evt_00000000000000",
		"type": "reporting.report_type.updated",
		"object": "event",
		"request": null,
		"pending_webhooks": 1,
		"api_version": "2019-05-16",
		"data": {
		  "object": {
			"id": "balance.summary.1_00000000000000",
			"object": "reporting.report_type",
			"data_available_end": 1583712000,
			"data_available_start": 1563753600,
			"default_columns": [
			  "category",
			  "description",
			  "net_amount",
			  "currency"
			],
			"name": "Balance summary",
			"updated": 1583770258,
			"version": 1
		  }
		}
	  }`
)

func makeRequest(body io.Reader) (req *http.Request) {
	req, _ = http.NewRequest("POST", "/", strings.NewReader(webhookBody))
	req.Header.Set("Content-Length", strconv.FormatInt(req.ContentLength, 10))
	return
}

func Test_MissingSignatureInHeader(t *testing.T) {
	req := makeRequest(strings.NewReader(webhookBody))

	app := fiber.New()
	app.Use(New(&Config{SigningSecret: "whsec_t4QaeaxpeR"}))
	app.Get("/", func(c *fiber.Ctx) {
		c.Send("ok")
	})
	resp, _ := app.Test(req)

	if http.StatusBadRequest != resp.StatusCode {
		t.Error("missing signature should return Bad Request")
	}
}

func Test_BadSignatureInHeader(t *testing.T) {
	req := makeRequest(strings.NewReader(webhookBody))
	req.Header.Set("Stripe-Signature", invalidSignature)

	app := fiber.New()
	app.Use(New(&Config{SigningSecret: "whsec_t4QaeaxpeR"}))
	app.Get("/", func(c *fiber.Ctx) {
		c.Send("ok")
	})
	resp, _ := app.Test(req)

	if http.StatusBadRequest != resp.StatusCode {
		t.Error("invalid signature should return Bad Request")
	}
}
