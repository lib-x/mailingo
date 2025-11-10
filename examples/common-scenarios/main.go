package main

import (
	"fmt"
	"log"

	"github.com/lib-x/mailingo"
)

func main() {
	product := mailingo.Product{
		Name:      "MyApp",
		Link:      "https://myapp.com",
		Logo:      "https://myapp.com/logo.png",
		Copyright: "© 2025 MyApp Inc. All rights reserved.",
	}

	mailer := mailingo.New(product, mailingo.DefaultTheme)

	// Scenario 1: Email Verification with Code
	verificationEmail := mailingo.Email{
		Body: mailingo.Body{
			Name:     "Alice Johnson",
			Greeting: "Hello",
			Title:    "Verify Your Email Address",
			Intros: []string{
				"Thank you for signing up! To complete your registration, please use the verification code below:",
			},
			Dictionary: []mailingo.Entry{
				{Key: "Verification Code", Value: "8 7 4 3 2 1"}, // Spaced for readability
				{Key: "Valid for", Value: "10 minutes"},
			},
			Outros: []string{
				"If you didn't create an account, you can safely ignore this email.",
				"For security reasons, never share this code with anyone.",
			},
			Signature: "Best regards",
		},
	}

	fmt.Println("========== Scenario 1: Email Verification Code ==========")
	html, err := mailer.GenerateHTML(verificationEmail, "en")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(html[:500] + "...")

	// Scenario 2: One-Time Login Link (Magic Link)
	magicLinkEmail := mailingo.Email{
		Body: mailingo.Body{
			Name:     "Bob Smith",
			Greeting: "Hello",
			Title:    "Sign In to Your Account",
			Intros: []string{
				"We received a request to sign in to your account.",
				"Click the button below to sign in securely without a password:",
			},
			Actions: []mailingo.Action{
				{
					Instructions: "This link will expire in 15 minutes:",
					Button: mailingo.Button{
						Text: "Sign In Now",
						Link: "https://myapp.com/auth/magic?token=abc123xyz789",
					},
				},
			},
			Outros: []string{
				"If you didn't request this, please ignore this email or contact support if you have concerns.",
			},
			Signature: "Best regards",
		},
	}

	fmt.Println("\n========== Scenario 2: One-Time Login (Magic Link) ==========")
	html, err = mailer.GenerateHTML(magicLinkEmail, "en")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(html[:500] + "...")

	// Scenario 3: Team Invitation
	teamInviteEmail := mailingo.Email{
		Body: mailingo.Body{
			Name:     "Carol Davis",
			Greeting: "Hello",
			Title:    "You've Been Invited to Join a Team!",
			Intros: []string{
				"John Doe has invited you to join the Engineering Team at Acme Corp.",
			},
			Dictionary: []mailingo.Entry{
				{Key: "Team Name", Value: "Engineering Team"},
				{Key: "Organization", Value: "Acme Corp"},
				{Key: "Invited by", Value: "John Doe (john@acme.com)"},
				{Key: "Role", Value: "Developer"},
			},
			Actions: []mailingo.Action{
				{
					Instructions: "Click the button below to accept the invitation:",
					Button: mailingo.Button{
						Text: "Accept Invitation",
						Link: "https://myapp.com/invites/accept?token=invite123",
					},
				},
				{
					Instructions:   "Or decline if you don't want to join:",
					InvertedButton: true,
					Button: mailingo.Button{
						Text: "Decline Invitation",
						Link: "https://myapp.com/invites/decline?token=invite123",
					},
				},
			},
			Outros: []string{
				"This invitation will expire in 7 days.",
			},
			Signature: "Best regards",
		},
	}

	fmt.Println("\n========== Scenario 3: Team Invitation ==========")
	html, err = mailer.GenerateHTML(teamInviteEmail, "en")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(html[:500] + "...")

	// Scenario 4: Monthly Billing Statement
	billingEmail := mailingo.Email{
		Body: mailingo.Body{
			Name:     "David Wilson",
			Greeting: "Hello",
			Title:    "Your January 2025 Billing Statement",
			Intros: []string{
				"Here's your billing summary for January 2025.",
			},
			Dictionary: []mailingo.Entry{
				{Key: "Billing Period", Value: "January 1 - January 31, 2025"},
				{Key: "Account ID", Value: "ACC-12345"},
				{Key: "Payment Method", Value: "•••• •••• •••• 4242"},
				{Key: "Next Billing Date", Value: "February 1, 2025"},
			},
			Table: mailingo.Table{
				Data: [][]mailingo.Entry{
					// Header
					{
						{Key: "Description"},
						{Key: "Quantity"},
						{Key: "Unit Price"},
						{Key: "Amount"},
					},
					// Items
					{
						{Value: "Pro Plan Subscription"},
						{Value: "1"},
						{Value: "$29.00"},
						{Value: "$29.00"},
					},
					{
						{Value: "Additional Users"},
						{Value: "5"},
						{Value: "$5.00"},
						{Value: "$25.00"},
					},
					{
						{Value: "API Calls (per 1000)"},
						{Value: "150"},
						{Value: "$0.10"},
						{Value: "$15.00"},
					},
					{
						{Value: "Storage Overage (GB)"},
						{Value: "20"},
						{Value: "$0.50"},
						{Value: "$10.00"},
					},
					// Subtotal
					{
						{Value: "Subtotal"},
						{Value: ""},
						{Value: ""},
						{Value: "$79.00"},
					},
					{
						{Value: "Tax (10%)"},
						{Value: ""},
						{Value: ""},
						{Value: "$7.90"},
					},
					{
						{Value: "Total Amount"},
						{Value: ""},
						{Value: ""},
						{Value: "$86.90"},
					},
				},
			},
			Actions: []mailingo.Action{
				{
					Instructions: "View detailed invoice or download PDF:",
					Button: mailingo.Button{
						Text: "View Invoice",
						Link: "https://myapp.com/billing/invoices/2025-01",
					},
				},
			},
			Outros: []string{
				"Payment was successfully processed on February 1, 2025.",
				"If you have any questions about your bill, please contact our billing support.",
			},
			Signature: "Best regards",
		},
	}

	fmt.Println("\n========== Scenario 4: Monthly Billing Statement ==========")
	html, err = mailer.GenerateHTML(billingEmail, "en")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(html[:500] + "...")

	// Plain text version for billing
	text, err := mailer.GeneratePlainText(billingEmail, "en")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\n\n--- Plain Text Version ---")
	fmt.Println(text)
}
