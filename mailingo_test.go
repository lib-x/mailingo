package mailingo

import (
	"embed"
	"os"
	"strings"
	"testing"

	"github.com/lib-x/mailingo/options"
)

//go:embed testdata/*.json
var testFS embed.FS

func TestNew(t *testing.T) {
	product := Product{
		Name:      "Test Product",
		Link:      "https://example.com",
		Logo:      "https://example.com/logo.png",
		Copyright: "© 2025 Test Product. All rights reserved.",
	}

	mailer := New(product, DefaultTheme)

	if mailer == nil {
		t.Fatal("Expected mailer to be non-nil")
	}

	if mailer.product.Name != product.Name {
		t.Errorf("Expected product name %s, got %s", product.Name, mailer.product.Name)
	}

	if mailer.theme.PrimaryColor != DefaultTheme.PrimaryColor {
		t.Errorf("Expected primary color %s, got %s", DefaultTheme.PrimaryColor, mailer.theme.PrimaryColor)
	}
}

func TestLoadMessageFile(t *testing.T) {
	// Create temporary test translation file
	content := `{
		"greeting": "Hello",
		"signature": "Best regards",
		"product.copyright": "© 2025 Test Product. All rights reserved."
	}`

	tmpFile, err := os.CreateTemp("", "test-*.en.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	tmpFile.Close()

	product := Product{
		Name:      "Test Product",
		Link:      "https://example.com",
		Copyright: "product.copyright",
	}

	mailer := New(product, DefaultTheme)
	err = mailer.LoadMessageFile(tmpFile.Name())
	if err != nil {
		t.Errorf("LoadMessageFile failed: %v", err)
	}
}

func TestLoadMessageFileFS(t *testing.T) {
	product := Product{
		Name:      "Test Product",
		Link:      "https://example.com",
		Copyright: "product.copyright",
	}

	mailer := New(product, DefaultTheme)
	err := mailer.LoadMessageFileFS(testFS, "testdata/en.json")
	if err != nil {
		t.Errorf("LoadMessageFileFS failed: %v", err)
	}
}

func TestGenerateHTML(t *testing.T) {
	product := Product{
		Name:      "Acme Corporation",
		Link:      "https://acme.com",
		Logo:      "https://acme.com/logo.png",
		Copyright: "product.copyright",
	}

	mailer := New(product, DefaultTheme)

	// Load test translations
	err := mailer.LoadMessageFileFS(testFS, "testdata/en.json")
	if err != nil {
		t.Fatalf("Failed to load translations: %v", err)
	}

	email := Email{
		Body: Body{
			Name:     "John Doe",
			Greeting: "greeting",
			Title:    "email.welcome.title",
			Intros: []string{
				"email.welcome.intro",
			},
			Dictionary: []Entry{
				{Key: "email.username", Value: "johndoe"},
				{Key: "email.email", Value: "john@example.com"},
			},
			Actions: []Action{
				{
					Instructions: "email.action.instructions",
					Button: Button{
						Text: "email.action.button",
						Link: "https://acme.com/confirm?token=abc123",
					},
				},
			},
			Outros: []string{
				"email.welcome.outro",
			},
			Signature: "signature",
		},
	}

	html, err := mailer.GenerateHTML(email, "en")
	if err != nil {
		t.Fatalf("GenerateHTML failed: %v", err)
	}

	// Verify HTML contains expected elements
	if !strings.Contains(html, "<!DOCTYPE html>") {
		t.Error("HTML should contain DOCTYPE")
	}

	if !strings.Contains(html, "John Doe") {
		t.Error("HTML should contain recipient name")
	}

	if !strings.Contains(html, "Acme Corporation") {
		t.Error("HTML should contain product name")
	}

	if !strings.Contains(html, DefaultTheme.PrimaryColor) {
		t.Error("HTML should contain theme primary color")
	}
}

func TestGeneratePlainText(t *testing.T) {
	product := Product{
		Name:      "Acme Corporation",
		Link:      "https://acme.com",
		Copyright: "product.copyright",
	}

	mailer := New(product, DefaultTheme)

	// Load test translations
	err := mailer.LoadMessageFileFS(testFS, "testdata/en.json")
	if err != nil {
		t.Fatalf("Failed to load translations: %v", err)
	}

	email := Email{
		Body: Body{
			Name:     "Jane Smith",
			Greeting: "greeting",
			Title:    "email.welcome.title",
			Intros: []string{
				"email.welcome.intro",
			},
			Dictionary: []Entry{
				{Key: "email.username", Value: "janesmith"},
				{Key: "email.email", Value: "jane@example.com"},
			},
			Actions: []Action{
				{
					Instructions: "email.action.instructions",
					Button: Button{
						Text: "email.action.button",
						Link: "https://acme.com/verify",
					},
				},
			},
			Outros: []string{
				"email.welcome.outro",
			},
			Signature: "signature",
		},
	}

	text, err := mailer.GeneratePlainText(email, "en")
	if err != nil {
		t.Fatalf("GeneratePlainText failed: %v", err)
	}

	// Verify plain text contains expected elements
	if !strings.Contains(text, "Jane Smith") {
		t.Error("Plain text should contain recipient name")
	}

	if !strings.Contains(text, "Acme Corporation") {
		t.Error("Plain text should contain product name")
	}

	if !strings.Contains(text, "janesmith") {
		t.Error("Plain text should contain username from dictionary")
	}
}

func TestMultipleLanguages(t *testing.T) {
	product := Product{
		Name:      "Global Corp",
		Link:      "https://globalcorp.com",
		Copyright: "product.copyright",
	}

	mailer := New(product, FlatTheme)

	// Load English translations
	err := mailer.LoadMessageFileFS(testFS, "testdata/en.json")
	if err != nil {
		t.Fatalf("Failed to load English translations: %v", err)
	}

	// Load Chinese translations
	err = mailer.LoadMessageFileFS(testFS, "testdata/zh.json")
	if err != nil {
		t.Fatalf("Failed to load Chinese translations: %v", err)
	}

	email := Email{
		Body: Body{
			Name:     "Zhang San",
			Greeting: "greeting",
			Intros: []string{
				"email.welcome.intro",
			},
			Signature: "signature",
		},
	}

	// Generate English version
	htmlEN, err := mailer.GenerateHTML(email, "en")
	if err != nil {
		t.Fatalf("Failed to generate English HTML: %v", err)
	}

	if !strings.Contains(htmlEN, "Hello") {
		t.Error("English HTML should contain 'Hello'")
	}

	// Generate Chinese version
	htmlZH, err := mailer.GenerateHTML(email, "zh")
	if err != nil {
		t.Fatalf("Failed to generate Chinese HTML: %v", err)
	}

	if !strings.Contains(htmlZH, "您好") {
		t.Error("Chinese HTML should contain '您好'")
	}
}

func TestThemes(t *testing.T) {
	product := Product{
		Name: "Test Product",
		Link: "https://example.com",
	}

	tests := []struct {
		name  string
		theme Theme
		color string
	}{
		{"DefaultTheme", DefaultTheme, "#3869D4"},
		{"FlatTheme", FlatTheme, "#2F3133"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mailer := New(product, tt.theme)
			email := Email{
				Body: Body{
					Name:   "Test User",
					Intros: []string{"Test message"},
				},
			}

			html, err := mailer.GenerateHTML(email, "en")
			if err != nil {
				t.Fatalf("GenerateHTML failed: %v", err)
			}

			if !strings.Contains(html, tt.color) {
				t.Errorf("HTML should contain theme color %s", tt.color)
			}
		})
	}
}

func TestCustomButtonColor(t *testing.T) {
	product := Product{
		Name: "Test Product",
		Link: "https://example.com",
	}

	mailer := New(product, DefaultTheme)

	customColor := "#FF5733"
	email := Email{
		Body: Body{
			Name: "Test User",
			Actions: []Action{
				{
					Instructions: "Click the button",
					Button: Button{
						Text:  "Custom Button",
						Link:  "https://example.com/action",
						Color: customColor,
					},
				},
			},
		},
	}

	html, err := mailer.GenerateHTML(email, "en")
	if err != nil {
		t.Fatalf("GenerateHTML failed: %v", err)
	}

	if !strings.Contains(html, customColor) {
		t.Errorf("HTML should contain custom button color %s", customColor)
	}
}

func TestInvertedButton(t *testing.T) {
	product := Product{
		Name: "Test Product",
		Link: "https://example.com",
	}

	mailer := New(product, DefaultTheme)

	email := Email{
		Body: Body{
			Name: "Test User",
			Actions: []Action{
				{
					Instructions:   "Click the inverted button",
					InvertedButton: true,
					Button: Button{
						Text: "Inverted Button",
						Link: "https://example.com/action",
					},
				},
			},
		},
	}

	html, err := mailer.GenerateHTML(email, "en")
	if err != nil {
		t.Fatalf("GenerateHTML failed: %v", err)
	}

	if !strings.Contains(html, "email-button-inverted") {
		t.Error("HTML should contain inverted button class")
	}
}

func TestTableRendering(t *testing.T) {
	product := Product{
		Name: "Test Product",
		Link: "https://example.com",
	}

	mailer := New(product, DefaultTheme)

	email := Email{
		Body: Body{
			Name: "Test User",
			Table: Table{
				Data: [][]Entry{
					// Header row
					{
						{Key: "Product"},
						{Key: "Quantity"},
						{Key: "Price"},
					},
					// Data rows
					{
						{Value: "Widget A"},
						{Value: "2"},
						{Value: "$19.99"},
					},
					{
						{Value: "Widget B"},
						{Value: "1"},
						{Value: "$29.99"},
					},
				},
			},
		},
	}

	html, err := mailer.GenerateHTML(email, "en")
	if err != nil {
		t.Fatalf("GenerateHTML failed: %v", err)
	}

	// Verify table elements
	if !strings.Contains(html, "email-table") {
		t.Error("HTML should contain table class")
	}

	if !strings.Contains(html, "Widget A") {
		t.Error("HTML should contain table data")
	}

	if !strings.Contains(html, "Product") {
		t.Error("HTML should contain table header")
	}
}

func TestEmptyEmail(t *testing.T) {
	product := Product{
		Name: "Test Product",
		Link: "https://example.com",
	}

	mailer := New(product, DefaultTheme)

	email := Email{
		Body: Body{
			Name: "Test User",
		},
	}

	html, err := mailer.GenerateHTML(email, "en")
	if err != nil {
		t.Fatalf("GenerateHTML should not fail with empty email: %v", err)
	}

	if !strings.Contains(html, "Test User") {
		t.Error("HTML should still contain recipient name")
	}

	text, err := mailer.GeneratePlainText(email, "en")
	if err != nil {
		t.Fatalf("GeneratePlainText should not fail with empty email: %v", err)
	}

	if !strings.Contains(text, "Test User") {
		t.Error("Plain text should still contain recipient name")
	}
}

func TestTranslateFallback(t *testing.T) {
	product := Product{
		Name: "Test Product",
		Link: "https://example.com",
	}

	mailer := New(product, DefaultTheme)

	// Don't load any translations
	email := Email{
		Body: Body{
			Name:     "Test User",
			Greeting: "nonexistent.key",
			Intros: []string{
				"This is a literal string, not a translation key",
			},
		},
	}

	html, err := mailer.GenerateHTML(email, "en")
	if err != nil {
		t.Fatalf("GenerateHTML failed: %v", err)
	}

	// Should fall back to the original key/text
	if !strings.Contains(html, "nonexistent.key") {
		t.Error("HTML should contain fallback key when translation is missing")
	}

	if !strings.Contains(html, "This is a literal string") {
		t.Error("HTML should contain literal string")
	}
}

func TestAttachments(t *testing.T) {
	product := Product{
		Name: "Test Product",
		Link: "https://example.com",
	}

	mailer := New(product, DefaultTheme)

	email := Email{
		Body: Body{
			Name: "Test User",
			Intros: []string{
				"Your documents are ready.",
			},
			Attachments: []Attachment{
				{
					Name: "document.pdf",
					URL:  "https://example.com/download/doc.pdf",
					Size: "2.5 MB",
					Type: "PDF Document",
				},
				{
					Name: "spreadsheet.xlsx",
					URL:  "https://example.com/download/sheet.xlsx",
					Size: "1.2 MB",
					Type: "Excel Spreadsheet",
				},
			},
		},
	}

	// Test HTML generation with attachments
	html, err := mailer.GenerateHTML(email, "en")
	if err != nil {
		t.Fatalf("GenerateHTML failed: %v", err)
	}

	if !strings.Contains(html, "email-attachments") {
		t.Error("HTML should contain attachments section")
	}

	if !strings.Contains(html, "document.pdf") {
		t.Error("HTML should contain first attachment name")
	}

	if !strings.Contains(html, "spreadsheet.xlsx") {
		t.Error("HTML should contain second attachment name")
	}

	if !strings.Contains(html, "2.5 MB") {
		t.Error("HTML should contain attachment size")
	}

	// Test plain text generation with attachments
	text, err := mailer.GeneratePlainText(email, "en")
	if err != nil {
		t.Fatalf("GeneratePlainText failed: %v", err)
	}

	if !strings.Contains(text, "Attachments:") {
		t.Error("Plain text should contain attachments header")
	}

	if !strings.Contains(text, "document.pdf") {
		t.Error("Plain text should contain attachment name")
	}

	if !strings.Contains(text, "https://example.com/download/doc.pdf") {
		t.Error("Plain text should contain attachment URL")
	}
}

func TestSMTPAttachments(t *testing.T) {
	product := Product{
		Name: "Test Product",
		Link: "https://example.com",
	}

	mailer := New(product, DefaultTheme)

	email := Email{
		Body: Body{
			Name: "Test User",
			Intros: []string{
				"Please see attached file.",
			},
		},
		SMTPAttachments: []SMTPAttachment{
			{
				Filename:    "report.pdf",
				Content:     []byte("PDF content"),
				ContentType: "application/pdf",
			},
		},
	}

	// SMTP attachments should not appear in HTML body
	html, err := mailer.GenerateHTML(email, "en")
	if err != nil {
		t.Fatalf("GenerateHTML failed: %v", err)
	}

	if strings.Contains(html, "report.pdf") {
		t.Error("HTML should not contain SMTP attachment filename")
	}

	// But the Email struct should have them
	if len(email.SMTPAttachments) != 1 {
		t.Error("Email should have SMTP attachments")
	}

	if email.SMTPAttachments[0].Filename != "report.pdf" {
		t.Error("SMTP attachment filename should be preserved")
	}
}

func TestCombinedAttachments(t *testing.T) {
	product := Product{
		Name: "Test Product",
		Link: "https://example.com",
	}

	mailer := New(product, DefaultTheme)

	email := Email{
		Body: Body{
			Name: "Test User",
			Intros: []string{
				"Files included both ways.",
			},
			Attachments: []Attachment{
				{
					Name: "large_file.zip",
					URL:  "https://example.com/download/large.zip",
					Size: "50 MB",
					Type: "ZIP Archive",
				},
			},
		},
		SMTPAttachments: []SMTPAttachment{
			{
				Filename:    "small_summary.pdf",
				Content:     []byte("PDF content"),
				ContentType: "application/pdf",
			},
		},
	}

	html, err := mailer.GenerateHTML(email, "en")
	if err != nil {
		t.Fatalf("GenerateHTML failed: %v", err)
	}

	// Body attachments should appear in HTML
	if !strings.Contains(html, "large_file.zip") {
		t.Error("HTML should contain body attachment")
	}

	// SMTP attachments should not appear in HTML
	if strings.Contains(html, "small_summary.pdf") {
		t.Error("HTML should not contain SMTP attachment")
	}

	// Both should be present in the Email struct
	if len(email.Body.Attachments) != 1 {
		t.Error("Should have body attachment")
	}

	if len(email.SMTPAttachments) != 1 {
		t.Error("Should have SMTP attachment")
	}
}

func TestCustomCSS(t *testing.T) {
	product := Product{
		Name: "Test Product",
		Link: "https://example.com",
	}

	customCSS := `
		.email-title { font-size: 32px; }
		.custom-class { color: red; }
	`

	mailer := New(product, DefaultTheme, options.WithCustomCSS(customCSS))

	email := Email{
		Body: Body{
			Name:  "Test User",
			Title: "Custom CSS Test",
			Intros: []string{
				"Testing custom CSS",
			},
		},
	}

	html, err := mailer.GenerateHTML(email, "en")
	if err != nil {
		t.Fatalf("GenerateHTML failed: %v", err)
	}

	// Debug: find and print the style section
	startIdx := strings.Index(html, "<style>")
	endIdx := strings.Index(html, "</style>")

	if startIdx != -1 && endIdx != -1 {
		styleSection := html[startIdx : endIdx+8]
		t.Logf("Style section (last 300 chars): %s", styleSection[len(styleSection)-300:])
	}

	// The customCSS should be in the style section
	if !strings.Contains(html, ".custom-class") {
		t.Error("HTML should contain custom class")
	}

	if !strings.Contains(html, "font-size: 32px") {
		t.Error("HTML should contain custom CSS font-size")
	}
}

func TestCustomTemplateString(t *testing.T) {
	product := Product{
		Name: "Test Product",
		Link: "https://example.com",
	}

	customTemplate := `
<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<h1>{{.Body.Title}}</h1>
	<p>Hello {{.Body.Name}}</p>
	<p>{{.Product.Name}}</p>
</body>
</html>
`

	mailer := New(product, DefaultTheme, options.WithCustomTemplateString(customTemplate))

	email := Email{
		Body: Body{
			Name:  "Alice",
			Title: "Test Title",
		},
	}

	html, err := mailer.GenerateHTML(email, "en")
	if err != nil {
		t.Fatalf("GenerateHTML failed: %v", err)
	}

	if !strings.Contains(html, "<h1>Test Title</h1>") {
		t.Error("HTML should use custom template")
	}

	if !strings.Contains(html, "<p>Hello Alice</p>") {
		t.Error("HTML should contain recipient name from custom template")
	}

	if !strings.Contains(html, "<p>Test Product</p>") {
		t.Error("HTML should contain product name from custom template")
	}
}

func TestCustomTemplateWithEmbedFS(t *testing.T) {
	product := Product{
		Name: "Test Product",
		Link: "https://example.com",
	}

	// Use the default template from embedded FS as a test
	mailer := New(product, DefaultTheme, options.WithCustomTemplateFile(templatesFS, "templates/default.html"))

	email := Email{
		Body: Body{
			Name: "Bob",
			Intros: []string{
				"Test message",
			},
		},
	}

	html, err := mailer.GenerateHTML(email, "en")
	if err != nil {
		t.Fatalf("GenerateHTML failed: %v", err)
	}

	if !strings.Contains(html, "Bob") {
		t.Error("HTML should contain recipient name")
	}

	if !strings.Contains(html, "Test message") {
		t.Error("HTML should contain intro text")
	}
}

func TestDefaultTemplateWithoutOptions(t *testing.T) {
	product := Product{
		Name: "Test Product",
		Link: "https://example.com",
	}

	// Create mailer without any options - should use default template
	mailer := New(product, DefaultTheme)

	email := Email{
		Body: Body{
			Name:  "Test User",
			Title: "Default Template Test",
		},
	}

	html, err := mailer.GenerateHTML(email, "en")
	if err != nil {
		t.Fatalf("GenerateHTML failed: %v", err)
	}

	// Should use default template
	if !strings.Contains(html, "email-container") {
		t.Error("Should use default template with standard classes")
	}

	if !strings.Contains(html, "Test User") {
		t.Error("HTML should contain recipient name")
	}
}
