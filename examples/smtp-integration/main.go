package main

import (
	"fmt"

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

	// Example 1: Body.Attachments (Download Links in Email Content)
	// This displays attachments as clickable links in the email body
	emailWithLinks := mailingo.Email{
		Body: mailingo.Body{
			Name:     "Alice",
			Greeting: "Hello",
			Title:    "Your Documents Are Ready",
			Intros: []string{
				"Your requested documents are now available for download.",
			},
			// These attachments will be shown as download links in the email
			Attachments: []mailingo.Attachment{
				{
					Name: "invoice_2025_01.pdf",
					URL:  "https://myapp.com/download/inv-123",
					Size: "245 KB",
					Type: "PDF Invoice",
				},
				{
					Name: "receipt.pdf",
					URL:  "https://myapp.com/download/rec-456",
					Size: "180 KB",
					Type: "PDF Receipt",
				},
			},
			Outros: []string{
				"These links will expire in 7 days.",
			},
			Signature: "Best regards",
		},
	}

	html1, _ := mailer.GenerateHTML(emailWithLinks, "en")
	fmt.Println("========== Example 1: Download Links in Email ==========")
	fmt.Println("Use this when files are hosted on your server and users click to download.")
	fmt.Println(html1[:500] + "...\n")

	// Example 2: SMTPAttachments (Real SMTP Attachments)
	// These are actual file attachments that will be sent with the email
	emailWithSMTPAttachments := mailingo.Email{
		Body: mailingo.Body{
			Name:     "Bob",
			Greeting: "Hello",
			Title:    "Monthly Report Attached",
			Intros: []string{
				"Please find your monthly report attached to this email.",
			},
			Outros: []string{
				"If you have any questions, please let us know.",
			},
			Signature: "Best regards",
		},
		// These are real file attachments (for SMTP sending)
		SMTPAttachments: []mailingo.SMTPAttachment{
			{
				Filename:    "monthly_report.pdf",
				Content:     []byte("PDF content here..."), // In real usage, read from file
				ContentType: "application/pdf",
			},
		},
	}

	html2, _ := mailer.GenerateHTML(emailWithSMTPAttachments, "en")
	fmt.Println("========== Example 2: Real SMTP Attachments ==========")
	fmt.Println("Use this when you want files attached directly to the email.")
	fmt.Println("Note: SMTPAttachments are not shown in the email body - they're for SMTP sending only.")
	fmt.Println(html2[:500] + "...\n")

	// Example 3: Both Types Together
	// Combine both: show some files as links + attach other files to email
	emailWithBoth := mailingo.Email{
		Body: mailingo.Body{
			Name:     "Carol",
			Greeting: "Hello",
			Title:    "Your Tax Documents",
			Intros: []string{
				"Your tax documents are ready. The summary is attached to this email.",
				"Additional supporting documents are available for download:",
			},
			// Large files as download links (saves email size)
			Attachments: []mailingo.Attachment{
				{
					Name: "supporting_docs.zip",
					URL:  "https://myapp.com/download/support-789",
					Size: "15.2 MB",
					Type: "ZIP Archive",
				},
			},
			Outros: []string{
				"The download link will remain active for 30 days.",
			},
			Signature: "Best regards",
		},
		// Small important file as real attachment
		SMTPAttachments: []mailingo.SMTPAttachment{
			{
				Filename:    "tax_summary.pdf",
				Content:     []byte("PDF content here..."),
				ContentType: "application/pdf",
			},
		},
	}

	html3, _ := mailer.GenerateHTML(emailWithBoth, "en")
	fmt.Println("========== Example 3: Combined Approach ==========")
	fmt.Println("Small files: SMTP attachments (convenient, immediately available)")
	fmt.Println("Large files: Download links (saves email size, better deliverability)")
	fmt.Println(html3[:500] + "...\n")

	// How to send with SMTP (example using go-mail library)
	fmt.Println("\n========== How to Send with SMTP ==========")
	fmt.Println(`// Example using go-mail library (github.com/wneessen/go-mail)
import "github.com/wneessen/go-mail"

func SendEmail(email mailingo.Email, htmlContent string) error {
    m := mail.NewMsg()
    m.From("sender@example.com")
    m.To("recipient@example.com")
    m.Subject("Subject")
    m.SetBodyString(mail.TypeTextHTML, htmlContent)

    // Add SMTP attachments if present
    for _, att := range email.SMTPAttachments {
        m.AttachReader(att.Filename, bytes.NewReader(att.Content))
    }

    c, _ := mail.NewClient("smtp.example.com",
        mail.WithPort(587),
        mail.WithSMTPAuth(mail.SMTPAuthPlain),
        mail.WithUsername("user"),
        mail.WithPassword("pass"))

    return c.DialAndSend(m)
}

// Example using gomail (gopkg.in/gomail.v2)
import "gopkg.in/gomail.v2"

func SendEmailGomail(email mailingo.Email, htmlContent string) error {
    m := gomail.NewMessage()
    m.SetHeader("From", "sender@example.com")
    m.SetHeader("To", "recipient@example.com")
    m.SetHeader("Subject", "Subject")
    m.SetBody("text/html", htmlContent)

    // Add SMTP attachments if present
    for _, att := range email.SMTPAttachments {
        m.Attach(att.Filename,
            gomail.SetCopyFunc(func(w io.Writer) error {
                _, err := w.Write(att.Content)
                return err
            }))
    }

    d := gomail.NewDialer("smtp.example.com", 587, "user", "pass")
    return d.DialAndSend(m)
}`)
}
