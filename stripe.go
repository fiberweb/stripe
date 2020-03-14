package stripe

import (
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber"
	"github.com/stripe/stripe-go/webhook"
)

// Config is used to configure the middleware
type Config struct {
	Skip          func(*fiber.Ctx) bool
	SigningSecret string
}

// New returns middleware for Stripe webhook
func New(config ...*Config) func(*fiber.Ctx) {
	var cfg *Config

	if len(config) == 0 {
		cfg = &Config{
			SigningSecret: os.Getenv("STRIPE_WEBHOOK_SIGNING_SECRET"),
		}
	} else {
		cfg = config[0]
	}

	if cfg.SigningSecret == "" {
		log.Fatalln("Stripe webhook: missing Signing Secret")
	}
	return func(c *fiber.Ctx) {
		if cfg.Skip != nil && cfg.Skip(c) {
			c.Next()
			return
		}
		event, err := webhook.ConstructEvent([]byte(c.Body()), c.Get("Stripe-Signature"), cfg.SigningSecret)
		if err != nil {
			c.Status(http.StatusBadRequest).SendString(err.Error())
			return
		}
		// put Stripe event to "stripeEvent" locals for next handler
		c.Locals("stripeEvent", event)
		c.Next()
	}
}
