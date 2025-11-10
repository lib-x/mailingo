package main

import (
	"fmt"
	"log"

	"github.com/lib-x/mailingo"
)

func main() {
	// Define your product information
	product := mailingo.Product{
		Name:      "ShopMart",
		Link:      "https://shopmart.com",
		Logo:      "https://shopmart.com/logo.png",
		Copyright: "© 2025 ShopMart Inc. All rights reserved.",
	}

	// Create a custom theme
	customTheme := mailingo.Theme{
		PrimaryColor:    "#FF6B35",
		BackgroundColor: "#F7F7F7",
		TextColor:       "#333333",
		ButtonColor:     "#FF6B35",
		ButtonTextColor: "#FFFFFF",
	}

	// Create a new mailer with custom theme
	mailer := mailingo.New(product, customTheme)

	// Create an order confirmation email with a table
	email := mailingo.Email{
		Body: mailingo.Body{
			Name:     "Sarah Williams",
			Greeting: "Hello",
			Title:    "Order Confirmation",
			Intros: []string{
				"Thank you for your order! Your order has been confirmed and will be shipped soon.",
			},
			Dictionary: []mailingo.Entry{
				{Key: "Order Number", Value: "#123456"},
				{Key: "Order Date", Value: "January 15, 2025"},
				{Key: "Total Amount", Value: "$149.97"},
			},
			Table: mailingo.Table{
				Data: [][]mailingo.Entry{
					// Header row
					{
						{Key: "Product"},
						{Key: "Quantity"},
						{Key: "Price"},
					},
					// Data rows
					{
						{Value: "Wireless Headphones"},
						{Value: "1"},
						{Value: "$79.99"},
					},
					{
						{Value: "Phone Case"},
						{Value: "2"},
						{Value: "$19.99"},
					},
					{
						{Value: "USB Cable"},
						{Value: "3"},
						{Value: "$9.99"},
					},
				},
			},
			Actions: []mailingo.Action{
				{
					Instructions: "Track your order status:",
					Button: mailingo.Button{
						Text: "Track Order",
						Link: "https://shopmart.com/orders/123456",
					},
				},
				{
					Instructions:   "View your invoice:",
					InvertedButton: true,
					Button: mailingo.Button{
						Text: "View Invoice",
						Link: "https://shopmart.com/invoices/123456",
					},
				},
			},
			Outros: []string{
				"We'll send you a notification when your order ships.",
				"If you have any questions about your order, please contact our support team.",
			},
			Signature: "Best regards",
		},
	}

	// Generate HTML email
	html, err := mailer.GenerateHTML(email, "en")
	if err != nil {
		log.Fatalf("Failed to generate HTML email: %v", err)
	}

	fmt.Println("Order Confirmation Email (HTML):")
	fmt.Println(html)

	// Generate plain text email
	text, err := mailer.GeneratePlainText(email, "en")
	if err != nil {
		log.Fatalf("Failed to generate plain text email: %v", err)
	}

	fmt.Println("\n\nOrder Confirmation Email (Plain Text):")
	fmt.Println(text)
}
