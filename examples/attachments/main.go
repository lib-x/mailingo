package main

import (
	"fmt"
	"log"

	"github.com/lib-x/mailingo"
)

func main() {
	product := mailingo.Product{
		Name:      "DocShare",
		Link:      "https://docshare.com",
		Logo:      "https://docshare.com/logo.png",
		Copyright: "© 2025 DocShare Inc. All rights reserved.",
	}

	mailer := mailingo.New(product, mailingo.DefaultTheme)

	// Scenario: Document Sharing with Attachments
	email := mailingo.Email{
		Body: mailingo.Body{
			Name:     "Emily Chen",
			Greeting: "Hello",
			Title:    "Documents Shared With You",
			Intros: []string{
				"John Smith has shared the following documents with you:",
			},
			Dictionary: []mailingo.Entry{
				{Key: "Shared by", Value: "John Smith (john@company.com)"},
				{Key: "Shared on", Value: "January 15, 2025"},
				{Key: "Access Level", Value: "View and Download"},
				{Key: "Expires", Value: "February 15, 2025"},
			},
			Attachments: []mailingo.Attachment{
				{
					Name: "Q4_Financial_Report_2024.pdf",
					URL:  "https://docshare.com/files/download?id=abc123",
					Size: "2.5 MB",
					Type: "PDF Document",
				},
				{
					Name: "2025_Budget_Proposal.xlsx",
					URL:  "https://docshare.com/files/download?id=def456",
					Size: "1.8 MB",
					Type: "Excel Spreadsheet",
				},
				{
					Name: "Project_Timeline.png",
					URL:  "https://docshare.com/files/download?id=ghi789",
					Size: "456 KB",
					Type: "Image",
				},
				{
					Name: "Meeting_Notes.docx",
					URL:  "https://docshare.com/files/download?id=jkl012",
					Size: "123 KB",
					Type: "Word Document",
				},
			},
			Actions: []mailingo.Action{
				{
					Instructions: "Click below to view all documents in your shared folder:",
					Button: mailingo.Button{
						Text: "Open Shared Folder",
						Link: "https://docshare.com/shared/folder-xyz",
					},
				},
			},
			Outros: []string{
				"These documents will be available for 30 days.",
				"If you have any questions about these documents, you can reply directly to the sender.",
			},
			Signature: "Best regards",
		},
	}

	// Generate HTML version
	fmt.Println("========== Document Sharing Email (HTML) ==========")
	html, err := mailer.GenerateHTML(email, "en")
	if err != nil {
		log.Fatalf("Failed to generate HTML: %v", err)
	}
	fmt.Println(html)

	// Generate Plain Text version
	fmt.Println("\n\n========== Document Sharing Email (Plain Text) ==========")
	text, err := mailer.GeneratePlainText(email, "en")
	if err != nil {
		log.Fatalf("Failed to generate plain text: %v", err)
	}
	fmt.Println(text)
}
