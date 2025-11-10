package main

import (
	"embed"
	"fmt"
	"log"

	"github.com/lib-x/mailingo"
)

//go:embed locales/*.json
var localesFS embed.FS

func main() {
	// Define your product information
	product := mailingo.Product{
		Name:      "Global Services",
		Link:      "https://globalservices.com",
		Logo:      "https://globalservices.com/logo.png",
		Copyright: "product.copyright", // This will be translated
	}

	// Create a new mailer with flat theme
	mailer := mailingo.New(product, mailingo.FlatTheme)

	// Load translations from embedded files
	if err := mailer.LoadMessageFileFS(localesFS, "locales/en.json"); err != nil {
		log.Fatalf("Failed to load English translations: %v", err)
	}

	if err := mailer.LoadMessageFileFS(localesFS, "locales/zh.json"); err != nil {
		log.Fatalf("Failed to load Chinese translations: %v", err)
	}

	if err := mailer.LoadMessageFileFS(localesFS, "locales/es.json"); err != nil {
		log.Fatalf("Failed to load Spanish translations: %v", err)
	}

	// Create an email with i18n keys
	email := mailingo.Email{
		Body: mailingo.Body{
			Name:     "Alice Johnson",
			Greeting: "greeting",
			Title:    "email.password_reset.title",
			Intros: []string{
				"email.password_reset.intro",
			},
			Actions: []mailingo.Action{
				{
					Instructions: "email.password_reset.instructions",
					Button: mailingo.Button{
						Text: "email.password_reset.button",
						Link: "https://globalservices.com/reset?token=xyz789",
					},
				},
			},
			Outros: []string{
				"email.password_reset.outro",
				"email.password_reset.warning",
			},
			Signature: "signature",
		},
	}

	// Generate emails in different languages
	languages := []string{"en", "zh", "es"}

	for _, lang := range languages {
		fmt.Printf("\n========== %s ==========\n", lang)

		html, err := mailer.GenerateHTML(email, lang)
		if err != nil {
			log.Fatalf("Failed to generate HTML for %s: %v", lang, err)
		}

		fmt.Printf("HTML Email Generated for %s\n", lang)
		fmt.Println(html[:200] + "...")

		text, err := mailer.GeneratePlainText(email, lang)
		if err != nil {
			log.Fatalf("Failed to generate plain text for %s: %v", lang, err)
		}

		fmt.Printf("\nPlain Text Email for %s:\n", lang)
		fmt.Println(text)
	}
}
