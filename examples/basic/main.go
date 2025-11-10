package main

import (
	"fmt"
	"log"

	"github.com/lib-x/mailingo"
)

func main() {
	// Define your product information
	product := mailingo.Product{
		Name:      "Acme Corporation",
		Link:      "https://acme.com",
		Logo:      "https://acme.com/logo.png",
		Copyright: "© 2025 Acme Corporation. All rights reserved.",
	}

	// Create a new mailer with the default theme
	mailer := mailingo.New(product, mailingo.DefaultTheme)

	// Create an email
	email := mailingo.Email{
		Body: mailingo.Body{
			Name:     "John Doe",
			Greeting: "Hello",
			Title:    "Welcome to Acme!",
			Intros: []string{
				"Thank you for signing up for Acme Corporation!",
				"We're excited to have you on board.",
			},
			Dictionary: []mailingo.Entry{
				{Key: "Username", Value: "johndoe"},
				{Key: "Email", Value: "john@example.com"},
				{Key: "Account Type", Value: "Premium"},
			},
			Actions: []mailingo.Action{
				{
					Instructions: "To get started with Acme, please click here:",
					Button: mailingo.Button{
						Text: "Get Started",
						Link: "https://acme.com/get-started",
					},
				},
			},
			Outros: []string{
				"Need help, or have questions? Just reply to this email, we'd love to help.",
			},
			Signature: "Best regards",
		},
	}

	// Generate HTML email
	html, err := mailer.GenerateHTML(email, "en")
	if err != nil {
		log.Fatalf("Failed to generate HTML email: %v", err)
	}

	fmt.Println("HTML Email Generated:")
	fmt.Println(html)

	// Generate plain text email
	text, err := mailer.GeneratePlainText(email, "en")
	if err != nil {
		log.Fatalf("Failed to generate plain text email: %v", err)
	}

	fmt.Println("\n\nPlain Text Email Generated:")
	fmt.Println(text)
}
