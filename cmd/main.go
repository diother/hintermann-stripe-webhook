package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/diother/hintermann-stripe-webhook/internal/handler"
	"github.com/stripe/stripe-go/v79"
)

func main() {
	stripeKey := os.Getenv("STRIPE_SECRET")
	webhookSecret := os.Getenv("WEBHOOK_SECRET")

	if stripeKey == "" || webhookSecret == "" {
		log.Fatal("env variables are missing")
	}
	stripe.Key = stripeKey

	http.HandleFunc("/webhook", handler.HandleWebhook(webhookSecret))

	fmt.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
