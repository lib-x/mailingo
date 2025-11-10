package mailingo

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"

	"github.com/lib-x/mailingo/options"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed templates/*.html
var templatesFS embed.FS

// Mailer is a multi-language email generator that supports i18n
type Mailer struct {
	bundle    *i18n.Bundle
	product   Product
	theme     Theme
	template  *template.Template
	customCSS string
}

// Product represents the product/company information displayed in emails
type Product struct {
	Name      string // Product or company name
	Link      string // Product or company website URL
	Logo      string // URL to the logo image
	Copyright string // Copyright text (supports i18n key, e.g., "product.copyright")
}

// Theme defines the color scheme and styling for the email
type Theme struct {
	PrimaryColor    string // Primary brand color
	BackgroundColor string // Email background color
	TextColor       string // Main text color
	ButtonColor     string // Button background color
	ButtonTextColor string // Button text color
}

// Email represents the complete email structure
type Email struct {
	Body            Body             // Email body content
	SMTPAttachments []SMTPAttachment // Files to be attached when sending via SMTP (not rendered in template)
}

// Body contains the main content of the email
type Body struct {
	Name        string       // Recipient's name
	Intros      []string     // Introduction paragraphs (supports i18n keys)
	Dictionary  []Entry      // Key-value pairs for structured information
	Table       Table        // Table data
	Actions     []Action     // Action buttons
	Outros      []string     // Closing paragraphs (supports i18n keys)
	Attachments []Attachment // List of attachments with download links
	Greeting    string       // Greeting text (supports i18n key, defaults to "greeting")
	Signature   string       // Signature text (supports i18n key, defaults to "signature")
	Title       string       // Email title (supports i18n key)
}

// Entry represents a key-value pair entry
type Entry struct {
	Key   string // Key text (supports i18n key)
	Value string // Value text
}

// Table represents tabular data in the email
type Table struct {
	Data    [][]Entry // Table rows, first row is treated as headers
	Columns Columns   // Column definitions
}

// Columns defines custom column properties
type Columns struct {
	CustomWidth     map[string]string // Custom column widths (e.g., "50%", "100px")
	CustomAlignment map[string]string // Custom column alignments (e.g., "left", "center", "right")
}

// Action represents a call-to-action button with instructions
type Action struct {
	Instructions   string // Instruction text above the button (supports i18n key)
	Button         Button // The action button
	InvertedButton bool   // Whether to use inverted button style (outlined)
}

// Button represents a clickable button in the email
type Button struct {
	Text  string // Button text (supports i18n key)
	Link  string // Button URL
	Color string // Custom button color (optional, overrides theme color)
}

// Attachment represents a file attachment that can be downloaded
type Attachment struct {
	Name string // File name (e.g., "invoice.pdf")
	URL  string // Download URL for the attachment
	Size string // Human-readable file size (e.g., "2.5 MB")
	Type string // File type or description (e.g., "PDF Document", "Excel Spreadsheet")
}

// SMTPAttachment represents an actual file to be attached to the email when sending via SMTP.
// This is separate from Attachment which only displays download links in the email body.
// Use this with your SMTP sending library (e.g., gomail, go-mail, etc.)
type SMTPAttachment struct {
	Filename    string // Name of the file as it will appear in the email
	Content     []byte // File content bytes
	ContentType string // MIME type (e.g., "application/pdf", "image/png")
}

// DefaultTheme is the default color theme (similar to Hermes default theme)
var DefaultTheme = Theme{
	PrimaryColor:    "#3869D4",
	BackgroundColor: "#F2F4F6",
	TextColor:       "#51545E",
	ButtonColor:     "#3869D4",
	ButtonTextColor: "#FFFFFF",
}

// FlatTheme is a flat minimalist theme (similar to Hermes flat theme)
var FlatTheme = Theme{
	PrimaryColor:    "#2F3133",
	BackgroundColor: "#FFFFFF",
	TextColor:       "#2F3133",
	ButtonColor:     "#2F3133",
	ButtonTextColor: "#FFFFFF",
}

// New creates a new Mailer instance with the specified product info and theme.
// The default language is set to English.
// You can customize the mailer using functional options.
//
// Example with default template:
//
//	mailer := mailingo.New(product, theme)
//
// Example with custom template:
//
//	mailer := mailingo.New(product, theme, options.WithCustomCSS("..."))
func New(product Product, theme Theme, opts ...options.Option) *Mailer {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	// Apply options
	config := &options.Config{}
	for _, opt := range opts {
		opt(config)
	}

	// Determine which template to use
	var tmpl *template.Template
	var err error

	if config.CustomTemplate != nil {
		// User provided a parsed template
		tmpl = config.CustomTemplate
	} else if config.CustomTemplateText != "" {
		// User provided a template string
		tmpl, err = template.New("email").Parse(config.CustomTemplateText)
		if err != nil {
			panic(fmt.Sprintf("failed to parse custom template: %v", err))
		}
	} else if config.CustomTemplateFS != nil && config.CustomTemplatePath != "" {
		// User provided a template file from embedded FS
		content, err := fs.ReadFile(config.CustomTemplateFS, config.CustomTemplatePath)
		if err != nil {
			panic(fmt.Sprintf("failed to read custom template file: %v", err))
		}
		tmpl, err = template.New("email").Parse(string(content))
		if err != nil {
			panic(fmt.Sprintf("failed to parse custom template file: %v", err))
		}
	} else {
		// Use default embedded template
		content, err := fs.ReadFile(templatesFS, "templates/default.html")
		if err != nil {
			panic(fmt.Sprintf("failed to read default template: %v", err))
		}
		tmpl, err = template.New("email").Parse(string(content))
		if err != nil {
			panic(fmt.Sprintf("failed to parse default template: %v", err))
		}
	}

	return &Mailer{
		bundle:    bundle,
		product:   product,
		theme:     theme,
		template:  tmpl,
		customCSS: config.CustomCSS,
	}
}

// LoadMessageFile loads translation messages from a file.
// The file format is determined by its extension (e.g., .json, .toml, .yaml).
func (m *Mailer) LoadMessageFile(path string) error {
	_, err := m.bundle.LoadMessageFile(path)
	return err
}

// LoadMessageFileFS loads translation messages from an embedded filesystem.
// This is useful when you embed translation files using go:embed directive.
func (m *Mailer) LoadMessageFileFS(fs fs.FS, path string) error {
	_, err := m.bundle.LoadMessageFileFS(fs, path)
	return err
}

// GenerateHTML generates an HTML email from the given email structure and language.
// The lang parameter should be a BCP 47 language tag (e.g., "en", "zh-CN").
func (m *Mailer) GenerateHTML(email Email, lang string) (string, error) {
	localizer := i18n.NewLocalizer(m.bundle, lang)

	// Process all translations
	data := m.processTranslations(email, localizer)

	// Render the HTML template
	var buf bytes.Buffer
	err := m.template.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute email template: %w", err)
	}

	return buf.String(), nil
}

// GeneratePlainText generates a plain text email from the given email structure and language.
// The lang parameter should be a BCP 47 language tag (e.g., "en", "zh-CN").
func (m *Mailer) GeneratePlainText(email Email, lang string) (string, error) {
	localizer := i18n.NewLocalizer(m.bundle, lang)

	var buf bytes.Buffer

	// Greeting
	greeting := m.translate(localizer, email.Body.Greeting, "greeting")
	buf.WriteString(fmt.Sprintf("%s %s,\n\n", greeting, email.Body.Name))

	// Title
	if email.Body.Title != "" {
		title := m.translate(localizer, email.Body.Title, "")
		buf.WriteString(fmt.Sprintf("%s\n\n", title))
	}

	// Introduction paragraphs
	for _, intro := range email.Body.Intros {
		text := m.translate(localizer, intro, "")
		buf.WriteString(fmt.Sprintf("%s\n\n", text))
	}

	// Dictionary (key-value pairs)
	for _, entry := range email.Body.Dictionary {
		key := m.translate(localizer, entry.Key, "")
		buf.WriteString(fmt.Sprintf("%s: %s\n", key, entry.Value))
	}
	if len(email.Body.Dictionary) > 0 {
		buf.WriteString("\n")
	}

	// Actions
	for _, action := range email.Body.Actions {
		instructions := m.translate(localizer, action.Instructions, "")
		buttonText := m.translate(localizer, action.Button.Text, "")
		buf.WriteString(fmt.Sprintf("%s\n%s: %s\n\n", instructions, buttonText, action.Button.Link))
	}

	// Closing paragraphs
	for _, outro := range email.Body.Outros {
		text := m.translate(localizer, outro, "")
		buf.WriteString(fmt.Sprintf("%s\n\n", text))
	}

	// Attachments
	if len(email.Body.Attachments) > 0 {
		buf.WriteString("Attachments:\n")
		for _, attachment := range email.Body.Attachments {
			buf.WriteString(fmt.Sprintf("  - %s", attachment.Name))
			if attachment.Type != "" || attachment.Size != "" {
				buf.WriteString(" (")
				if attachment.Type != "" {
					buf.WriteString(attachment.Type)
				}
				if attachment.Size != "" {
					if attachment.Type != "" {
						buf.WriteString(", ")
					}
					buf.WriteString(attachment.Size)
				}
				buf.WriteString(")")
			}
			buf.WriteString(fmt.Sprintf("\n    %s\n", attachment.URL))
		}
		buf.WriteString("\n")
	}

	// Signature
	signature := m.translate(localizer, email.Body.Signature, "signature")
	buf.WriteString(fmt.Sprintf("%s,\n%s\n\n", signature, m.product.Name))

	// Copyright
	copyright := m.translate(localizer, m.product.Copyright, "product.copyright")
	buf.WriteString(copyright)

	return buf.String(), nil
}

// translate is a helper function that translates a message ID using the localizer.
// If the key is empty and a defaultKey is provided, it uses the defaultKey.
// If translation fails, it returns the original key as fallback.
func (m *Mailer) translate(localizer *i18n.Localizer, key string, defaultKey string) string {
	if key == "" && defaultKey != "" {
		key = defaultKey
	}
	if key == "" {
		return ""
	}

	// Try to localize the message
	result, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: key,
	})
	if err != nil {
		// If translation fails, return the original key as fallback
		return key
	}
	return result
}

// processTranslations processes all translations in the email structure
func (m *Mailer) processTranslations(email Email, localizer *i18n.Localizer) map[string]interface{} {
	body := email.Body

	// Translate introduction paragraphs
	intros := make([]string, len(body.Intros))
	for i, intro := range body.Intros {
		intros[i] = m.translate(localizer, intro, "")
	}

	// Translate closing paragraphs
	outros := make([]string, len(body.Outros))
	for i, outro := range body.Outros {
		outros[i] = m.translate(localizer, outro, "")
	}

	// Translate dictionary entries
	dictionary := make([]Entry, len(body.Dictionary))
	for i, entry := range body.Dictionary {
		dictionary[i] = Entry{
			Key:   m.translate(localizer, entry.Key, ""),
			Value: entry.Value,
		}
	}

	// Translate actions
	actions := make([]Action, len(body.Actions))
	for i, action := range body.Actions {
		actions[i] = Action{
			Instructions: m.translate(localizer, action.Instructions, ""),
			Button: Button{
				Text:  m.translate(localizer, action.Button.Text, ""),
				Link:  action.Button.Link,
				Color: action.Button.Color,
			},
			InvertedButton: action.InvertedButton,
		}
	}

	// Translate table data
	tableData := make([][]Entry, len(body.Table.Data))
	for i, row := range body.Table.Data {
		tableData[i] = make([]Entry, len(row))
		for j, cell := range row {
			tableData[i][j] = Entry{
				Key:   m.translate(localizer, cell.Key, ""),
				Value: cell.Value,
			}
		}
	}

	// Process attachments (no translation needed, but include for consistency)
	attachments := make([]Attachment, len(body.Attachments))
	for i, attachment := range body.Attachments {
		attachments[i] = attachment
	}

	return map[string]interface{}{
		"Product": map[string]interface{}{
			"Name":      m.product.Name,
			"Link":      m.product.Link,
			"Logo":      m.product.Logo,
			"Copyright": m.translate(localizer, m.product.Copyright, "product.copyright"),
		},
		"Theme":     m.theme,
		"CustomCSS": template.CSS(m.customCSS), // Use template.CSS for CSS context
		"Body": map[string]interface{}{
			"Name":       body.Name,
			"Greeting":   m.translate(localizer, body.Greeting, "greeting"),
			"Signature":  m.translate(localizer, body.Signature, "signature"),
			"Title":      m.translate(localizer, body.Title, ""),
			"Intros":     intros,
			"Dictionary": dictionary,
			"Table": map[string]interface{}{
				"Data":    tableData,
				"Columns": body.Table.Columns,
			},
			"Actions":     actions,
			"Outros":      outros,
			"Attachments": attachments,
		},
	}
}
