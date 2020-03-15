# stripe
This is a [Fiber](https://github.com/gofiber/fiber) middleware to validates [Stripe webhooks](https://stripe.com/docs/webhooks) signed request.

## Install

```
go get -u github.com/fiberweb/stripe
```

## Usage

```
package main

import (
  "github.com/gofiber/fiber"
  "github.com/stripe/stripe-go"
  
  webhook "github.com/fiberweb/stripe"
)

func main() {
  app := fiber.New()
  
  // use the middleware
  app.Use(webhook.New(&webhook.Config{SigningSecret: "whsec_t4QaeaxpeR"}))
  
  // webhook handler
  app.Post("/webhook", func(c *fiber.Ctx) {
    event := c.Locals("StripeEvent").(stripe.Event)
    log.Println("Stripe webhook received with type:", event.Type)
    c.Send("Ok")
  })
  app.Listen("8080")
}
```

Please note that when the webhook request successfully pass the middleware checking, Stripe event will be available for the next handlers inside the Fiber context `Locals` called `StripeEvent`.

You could also specify the Signing Secret on Environment Variable `STRIPE_WEBHOOK_SIGNING_SECRET`, so you don't need to specify Signin Secret in the code, example:

```
app.Use(webhook.New())
```

When you run your app, you need to make sure the `STRIPE_WEBHOOK_SIGNING_SECRET` variable is set:

```
$ export STRIPE_WEBHOOK_SIGNING_SECRET=whsec_t4QaeaxpeR
$ ./MyGoApp

# for one liners:

$ STRIPE_WEBHOOK_SIGNING_SECRET=whsec_t4QaeaxpeR ./MyGoApp
```

## Configuration

This middleware has only two configuration options:

```
type Config struct {
	Skip          func(*fiber.Ctx) bool
	SigningSecret string
}
```

#### `Skip func(*fiber.Ctx) bool`
This is to skip Fiber from using this middleware based on certain condition, example:

```
app.Use(webhook.New(&webhook.Config{
    SigningSecret: "whsec_t4QaeaxpeR",
    Skip: func(c *fiber.Ctx) bool {
      // add your logic here
      return true // returning true will skip this middleware.
    }
}))
```

#### `SigningSecret string`
This is a webhook Signing Secret that you will get from Stripe dashboard.
