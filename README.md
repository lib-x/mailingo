# Mailingo

A powerful and flexible Go library for generating beautiful, multi-language email templates with built-in internationalization (i18n) support. Inspired by [Hermes](https://github.com/matcornic/hermes), Mailingo makes it easy to create professional HTML and plain text emails.

## Features

- **Multi-language Support**: Built-in i18n using [go-i18n](https://github.com/nicksnyder/go-i18n)
- **Beautiful Themes**: Pre-built themes (Default, Flat) with custom theme support
- **HTML & Plain Text**: Generate both HTML and plain text versions of emails
- **Rich Content**: Support for tables, dictionaries, action buttons, and more
- **Embedded Resources**: Works seamlessly with Go's `embed` package
- **Easy to Use**: Simple, intuitive API for creating professional emails
- **Customizable**: Full control over colors, styling, and content
- **Type-Safe**: Strong typing with comprehensive documentation

## Installation

```bash
go get github.com/lib-x/mailingo
```

## Quick Start

```go
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

    // Create a new mailer with a theme
    mailer := mailingo.New(product, mailingo.DefaultTheme)

    // Create an email
    email := mailingo.Email{
        Body: mailingo.Body{
            Name:     "John Doe",
            Greeting: "Hello",
            Title:    "Welcome to Acme!",
            Intros: []string{
                "Thank you for signing up!",
                "We're excited to have you on board.",
            },
            Actions: []mailingo.Action{
                {
                    Instructions: "Click the button below to get started:",
                    Button: mailingo.Button{
                        Text: "Get Started",
                        Link: "https://acme.com/get-started",
                    },
                },
            },
            Outros: []string{
                "Need help? Just reply to this email.",
            },
            Signature: "Best regards",
        },
    }

    // Generate HTML email
    html, err := mailer.GenerateHTML(email, "en")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(html)
}
```

## Internationalization (i18n)

Mailingo has powerful built-in support for multiple languages:

### 1. Create Translation Files

Create JSON files for each language (e.g., `locales/en.json`, `locales/zh.json`):

**en.json:**
```json
{
  "greeting": "Hello",
  "signature": "Best regards",
  "email.welcome.title": "Welcome!",
  "email.welcome.intro": "Thank you for signing up!"
}
```

**zh.json:**
```json
{
  "greeting": "您好",
  "signature": "此致敬礼",
  "email.welcome.title": "欢迎！",
  "email.welcome.intro": "感谢您的注册！"
}
```

### 2. Load Translations

```go
// Load from files
mailer.LoadMessageFile("locales/en.json")
mailer.LoadMessageFile("locales/zh.json")

// Or use embedded files
//go:embed locales/*.json
var localesFS embed.FS

mailer.LoadMessageFileFS(localesFS, "locales/en.json")
mailer.LoadMessageFileFS(localesFS, "locales/zh.json")
```

### 3. Use Translation Keys in Emails

```go
email := mailingo.Email{
    Body: mailingo.Body{
        Name:     "Alice",
        Greeting: "greeting",           // Will be translated
        Title:    "email.welcome.title", // Will be translated
        Intros: []string{
            "email.welcome.intro",       // Will be translated
        },
        Signature: "signature",          // Will be translated
    },
}

// Generate in different languages
htmlEN, _ := mailer.GenerateHTML(email, "en")
htmlZH, _ := mailer.GenerateHTML(email, "zh")
```

## Themes

Mailingo comes with two pre-built themes:

### Default Theme
```go
mailer := mailingo.New(product, mailingo.DefaultTheme)
```

### Flat Theme
```go
mailer := mailingo.New(product, mailingo.FlatTheme)
```

### Custom Theme
```go
customTheme := mailingo.Theme{
    PrimaryColor:    "#FF6B35",
    BackgroundColor: "#F7F7F7",
    TextColor:       "#333333",
    ButtonColor:     "#FF6B35",
    ButtonTextColor: "#FFFFFF",
}

mailer := mailingo.New(product, customTheme)
```

## Template Customization

Mailingo provides three ways to customize the email template to match your brand and requirements.

### 1. Custom CSS

Add your own CSS to override or extend the default template styles:

```go
import "github.com/lib-x/mailingo/options"

customCSS := `
    /* Override default styles */
    .email-title {
        font-size: 32px;
        font-weight: 900;
        text-transform: uppercase;
    }

    /* Add rounded buttons */
    .email-button {
        border-radius: 25px;
        box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    }

    /* Custom gradient for dictionary section */
    .email-dictionary {
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        color: white;
    }
`

mailer := mailingo.New(
    product,
    mailingo.DefaultTheme,
    options.WithCustomCSS(customCSS),
)
```

### 2. Custom Template String

Provide your own complete HTML template:

```go
minimalistTemplate := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body {
            font-family: 'Helvetica Neue', Arial, sans-serif;
            background-color: #f5f5f5;
            padding: 40px;
        }
        .container {
            max-width: 500px;
            margin: 0 auto;
            background: white;
            padding: 40px;
            border-radius: 8px;
        }
        h1 {
            color: #333;
            border-bottom: 2px solid {{.Theme.PrimaryColor}};
        }
        .button {
            background: {{.Theme.ButtonColor}};
            color: white;
            padding: 12px 30px;
            text-decoration: none;
            border-radius: 4px;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>{{.Body.Title}}</h1>
        <p>Hi {{.Body.Name}},</p>
        {{range .Body.Intros}}
        <p>{{.}}</p>
        {{end}}
        {{range .Body.Actions}}
        <a href="{{.Button.Link}}" class="button">{{.Button.Text}}</a>
        {{end}}
        <div class="footer">
            {{.Product.Copyright}}
        </div>
    </div>
</body>
</html>
`

mailer := mailingo.New(
    product,
    mailingo.DefaultTheme,
    options.WithCustomTemplateString(minimalistTemplate),
)
```

### 3. Custom Template from Embedded FS

Use a template file from your embedded filesystem:

```go
//go:embed templates/custom.html
var templatesFS embed.FS

mailer := mailingo.New(
    product,
    mailingo.DefaultTheme,
    options.WithCustomTemplateFS(templatesFS, "templates/custom.html"),
)
```

### Template Variables

When creating custom templates, you have access to these template variables:

```go
{{.Product.Name}}      // Product name
{{.Product.Link}}      // Product URL
{{.Product.Logo}}      // Logo URL
{{.Product.Copyright}} // Copyright text

{{.Theme.PrimaryColor}}    // Primary color
{{.Theme.BackgroundColor}} // Background color
{{.Theme.TextColor}}       // Text color
{{.Theme.ButtonColor}}     // Button color
{{.Theme.ButtonTextColor}} // Button text color

{{.CustomCSS}}             // Custom CSS (if provided)

{{.Body.Name}}             // Recipient name
{{.Body.Greeting}}         // Greeting text
{{.Body.Title}}            // Email title
{{.Body.Signature}}        // Signature text
{{.Body.Intros}}           // Array of intro paragraphs
{{.Body.Outros}}           // Array of outro paragraphs
{{.Body.Dictionary}}       // Array of Entry (Key/Value pairs)
{{.Body.Table.Data}}       // 2D array of table data
{{.Body.Actions}}          // Array of actions (Instructions, Button, InvertedButton)
{{.Body.Attachments}}      // Array of attachments (Name, URL, Size, Type)
```

### Combining Customizations

You can combine custom CSS with custom themes for maximum flexibility:

```go
customTheme := mailingo.Theme{
    PrimaryColor:    "#FF6B6B",
    BackgroundColor: "#FFF5F5",
    TextColor:       "#2C3E50",
    ButtonColor:     "#FF6B6B",
    ButtonTextColor: "#FFFFFF",
}

customCSS := `
    .email-body {
        background: linear-gradient(to bottom, #ffffff 0%, #f8f9fa 100%);
    }
    .email-greeting {
        font-size: 20px;
        font-style: italic;
    }
`

mailer := mailingo.New(
    product,
    customTheme,
    options.WithCustomCSS(customCSS),
)
```

See **[examples/custom-template](./examples/custom-template)** for complete working examples of all customization approaches.

## Email Components

### Basic Structure

```go
email := mailingo.Email{
    Body: mailingo.Body{
        Name:       "Recipient Name",    // Required
        Greeting:   "Hello",             // Optional (default: "greeting" i18n key)
        Title:      "Email Title",       // Optional
        Intros:     []string{},          // Introduction paragraphs
        Dictionary: []mailingo.Entry{},  // Key-value pairs
        Table:      mailingo.Table{},    // Tabular data
        Actions:    []mailingo.Action{}, // Call-to-action buttons
        Outros:     []string{},          // Closing paragraphs
        Signature:  "Best regards",      // Optional (default: "signature" i18n key)
    },
}
```

### Dictionary (Key-Value Pairs)

Display structured information:

```go
Dictionary: []mailingo.Entry{
    {Key: "Username", Value: "johndoe"},
    {Key: "Email", Value: "john@example.com"},
    {Key: "Account Type", Value: "Premium"},
}
```

### Tables

Perfect for order details, invoices, etc.:

```go
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
}
```

### Action Buttons

Add call-to-action buttons:

```go
Actions: []mailingo.Action{
    {
        Instructions: "Click the button below to confirm:",
        Button: mailingo.Button{
            Text: "Confirm Email",
            Link: "https://example.com/confirm",
        },
    },
}
```

#### Custom Button Colors

```go
Button: mailingo.Button{
    Text:  "Custom Button",
    Link:  "https://example.com",
    Color: "#FF5733", // Custom color
}
```

#### Inverted Buttons (Outlined)

```go
Action{
    Instructions:   "Or view your account:",
    InvertedButton: true, // Outlined style
    Button: mailingo.Button{
        Text: "View Account",
        Link: "https://example.com/account",
    },
}
```

### Attachments

Mailingo supports **two types** of attachments:

#### 1. Download Link Attachments (Displayed in Email Body)

Perfect for cloud file sharing and large files. Shows attachment cards in the email with download links.

```go
Body: mailingo.Body{
    Name: "John Doe",
    Intros: []string{
        "Your documents are ready for download:",
    },
    Attachments: []mailingo.Attachment{
        {
            Name: "invoice.pdf",
            URL:  "https://yourserver.com/download?id=123",
            Size: "2.5 MB",
            Type: "PDF Document",
        },
        {
            Name: "report.xlsx",
            URL:  "https://yourserver.com/download?id=456",
            Size: "1.8 MB",
            Type: "Excel Spreadsheet",
        },
    },
}
```

**Best for:**
- Large files (doesn't impact email size)
- Cloud-hosted files
- Better email deliverability
- Tracking downloads
- Expiring links

#### 2. SMTP Attachments (Real Email Attachments)

Real file attachments that are embedded in the email when sent via SMTP.

```go
email := mailingo.Email{
    Body: mailingo.Body{
        Name: "Jane Smith",
        Intros: []string{
            "Your report is attached.",
        },
    },
    // These files will be attached to the email
    SMTPAttachments: []mailingo.SMTPAttachment{
        {
            Filename:    "monthly_report.pdf",
            Content:     fileBytes, // File content as []byte
            ContentType: "application/pdf",
        },
    },
}

// Send with your SMTP library (see examples/smtp-integration)
```

**Best for:**
- Small files
- Important documents that should be immediately accessible
- Better user experience (no extra clicks)

#### 3. Combined Approach (Recommended)

Use both types together for optimal experience:

```go
email := mailingo.Email{
    Body: mailingo.Body{
        Name: "User Name",
        Intros: []string{
            "Summary attached, supporting docs available below:",
        },
        // Large files as download links
        Attachments: []mailingo.Attachment{
            {
                Name: "full_report.zip",
                URL:  "https://app.com/download/report-123",
                Size: "25 MB",
                Type: "ZIP Archive",
            },
        },
    },
    // Small important file as real attachment
    SMTPAttachments: []mailingo.SMTPAttachment{
        {
            Filename:    "summary.pdf",
            Content:     summaryBytes,
            ContentType: "application/pdf",
        },
    },
}
```

See **[examples/attachments](./examples/attachments)** and **[examples/smtp-integration](./examples/smtp-integration)** for complete examples.

## Common Use Cases

Mailingo supports all common email scenarios out of the box:

- ✅ **Email Verification Codes** - Use `Dictionary` to display verification codes
- ✅ **Magic Links / One-Time Login** - Use `Actions` for secure login buttons
- ✅ **Team Invitations** - Use `Dictionary` + `Actions` (accept/decline buttons)
- ✅ **Billing Statements** - Use `Table` for itemized billing + `Dictionary` for summary
- ✅ **File Sharing** - Use `Attachments` (download links) or `SMTPAttachments` (real files)

See **[SCENARIOS.md](./SCENARIOS.md)** for detailed examples of each scenario.

## Examples

Check out the [examples](./examples) directory for complete, runnable examples:

- **[basic](./examples/basic)**: Simple email with basic features
- **[common-scenarios](./examples/common-scenarios)**: Verification codes, magic links, team invites, billing
- **[multilingual](./examples/multilingual)**: Multi-language email with i18n (English, Chinese, Spanish)
- **[invoice](./examples/invoice)**: Order confirmation with table and custom theme
- **[attachments](./examples/attachments)**: File sharing with download links
- **[smtp-integration](./examples/smtp-integration)**: How to use SMTP attachments with go-mail/gomail
- **[custom-template](./examples/custom-template)**: Custom CSS, custom templates, and combined customizations

## API Reference

### Types

#### Mailer
```go
type Mailer struct {
    // Internal fields
}
```

#### Product
```go
type Product struct {
    Name      string // Product or company name
    Link      string // Product or company website URL
    Logo      string // URL to the logo image
    Copyright string // Copyright text (supports i18n key)
}
```

#### Theme
```go
type Theme struct {
    PrimaryColor    string // Primary brand color
    BackgroundColor string // Email background color
    TextColor       string // Main text color
    ButtonColor     string // Button background color
    ButtonTextColor string // Button text color
}
```

#### Email
```go
type Email struct {
    Body            Body             // Email body content
    SMTPAttachments []SMTPAttachment // Files to be attached when sending via SMTP
}
```

#### Body
```go
type Body struct {
    Name        string       // Recipient's name
    Intros      []string     // Introduction paragraphs
    Dictionary  []Entry      // Key-value pairs
    Table       Table        // Table data
    Actions     []Action     // Action buttons
    Outros      []string     // Closing paragraphs
    Attachments []Attachment // Download link attachments (displayed in email body)
    Greeting    string       // Greeting text
    Signature   string       // Signature text
    Title       string       // Email title
}
```

#### Attachment
```go
type Attachment struct {
    Name string // File name (e.g., "invoice.pdf")
    URL  string // Download URL
    Size string // Human-readable size (e.g., "2.5 MB")
    Type string // File type (e.g., "PDF Document")
}
```

#### SMTPAttachment
```go
type SMTPAttachment struct {
    Filename    string // Name as it appears in email
    Content     []byte // File content bytes
    ContentType string // MIME type (e.g., "application/pdf")
}
```

### Methods

#### New
```go
func New(product Product, theme Theme, opts ...options.Option) *Mailer
```
Creates a new Mailer instance with the specified product info and theme. Optionally accepts customization options:

- `options.WithCustomCSS(css string)`: Add custom CSS to the default template
- `options.WithCustomTemplateString(template string)`: Use a custom template string
- `options.WithCustomTemplateFS(fs fs.FS, path string)`: Use a custom template from embedded filesystem

Example:
```go
mailer := mailingo.New(
    product,
    theme,
    options.WithCustomCSS("/* custom styles */"),
)
```

#### LoadMessageFile
```go
func (m *Mailer) LoadMessageFile(path string) error
```
Loads translation messages from a file.

#### LoadMessageFileFS
```go
func (m *Mailer) LoadMessageFileFS(fs fs.FS, path string) error
```
Loads translation messages from an embedded filesystem.

#### GenerateHTML
```go
func (m *Mailer) GenerateHTML(email Email, lang string) (string, error)
```
Generates an HTML email. The `lang` parameter should be a BCP 47 language tag (e.g., "en", "zh-CN").

#### GeneratePlainText
```go
func (m *Mailer) GeneratePlainText(email Email, lang string) (string, error)
```
Generates a plain text email. The `lang` parameter should be a BCP 47 language tag (e.g., "en", "zh-CN").

## Testing

Run the test suite:

```bash
go test -v
```

Run tests with coverage:

```bash
go test -v -cover
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Inspired by [Hermes](https://github.com/matcornic/hermes)
- Uses [go-i18n](https://github.com/nicksnyder/go-i18n) for internationalization
- Built with Go's standard library and minimal dependencies

## Support

If you have questions or need help, please:
- Open an issue on GitHub
- Check the [examples](./examples) directory
- Review the test files for usage patterns
