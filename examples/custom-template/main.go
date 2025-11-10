package main

import (
	"fmt"
	"log"

	"github.com/lib-x/mailingo"
	"github.com/lib-x/mailingo/options"
)

func main() {
	product := mailingo.Product{
		Name:      "CustomApp",
		Link:      "https://customapp.com",
		Logo:      "https://customapp.com/logo.png",
		Copyright: "© 2025 CustomApp Inc.",
	}

	// Example 1: Using default template with custom CSS
	fmt.Println("========== Example 1: Custom CSS ==========")

	customCSS := `
        /* Custom styling */
        .email-title {
            font-size: 32px;
            font-weight: 900;
            text-transform: uppercase;
            letter-spacing: 2px;
        }
        .email-button {
            border-radius: 25px;
            padding: 15px 40px;
            font-size: 18px;
            box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
        }
        .email-dictionary {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            border-radius: 10px;
        }
        .email-dictionary-key {
            color: #fff;
            font-weight: bold;
        }
    `

	mailer1 := mailingo.New(
		product,
		mailingo.DefaultTheme,
		options.WithCustomCSS(customCSS),
	)

	email := mailingo.Email{
		Body: mailingo.Body{
			Name:  "Alice",
			Title: "Welcome to CustomApp",
			Intros: []string{
				"We're excited to have you on board!",
			},
			Dictionary: []mailingo.Entry{
				{Key: "Username", Value: "alice123"},
				{Key: "Account Type", Value: "Premium"},
			},
			Actions: []mailingo.Action{
				{
					Instructions: "Get started now:",
					Button: mailingo.Button{
						Text: "Launch App",
						Link: "https://customapp.com/dashboard",
					},
				},
			},
			Signature: "Best regards",
		},
	}

	html, err := mailer1.GenerateHTML(email, "en")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("HTML with custom CSS generated successfully!")
	fmt.Println(html[:300] + "...\n")

	// Example 2: Using custom template string
	fmt.Println("========== Example 2: Custom Template String ==========")

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
            margin: 0;
        }
        .container {
            max-width: 500px;
            margin: 0 auto;
            background: white;
            padding: 40px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        h1 {
            color: #333;
            font-size: 24px;
            margin-bottom: 20px;
            border-bottom: 2px solid {{.Theme.PrimaryColor}};
            padding-bottom: 10px;
        }
        p {
            color: #666;
            line-height: 1.8;
        }
        .button {
            display: inline-block;
            background: {{.Theme.ButtonColor}};
            color: white;
            padding: 12px 30px;
            text-decoration: none;
            border-radius: 4px;
            margin-top: 20px;
        }
        .footer {
            margin-top: 30px;
            padding-top: 20px;
            border-top: 1px solid #eee;
            color: #999;
            font-size: 12px;
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
            {{.Product.Copyright}}<br>
            <a href="{{.Product.Link}}">{{.Product.Name}}</a>
        </div>
    </div>
</body>
</html>
`

	mailer2 := mailingo.New(
		product,
		mailingo.FlatTheme,
		options.WithCustomTemplateString(minimalistTemplate),
	)

	email2 := mailingo.Email{
		Body: mailingo.Body{
			Name:  "Bob",
			Title: "Simple Notification",
			Intros: []string{
				"This is a minimalist email template.",
				"It focuses on simplicity and readability.",
			},
			Actions: []mailingo.Action{
				{
					Button: mailingo.Button{
						Text: "View Details",
						Link: "https://customapp.com/details",
					},
				},
			},
		},
	}

	html2, err := mailer2.GenerateHTML(email2, "en")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("HTML with custom template generated successfully!")
	fmt.Println(html2[:300] + "...\n")

	// Example 3: Multiple customizations together
	fmt.Println("========== Example 3: Combined Customizations ==========")

	extraCSS := `
        .email-body {
            background: linear-gradient(to bottom, #ffffff 0%, #f8f9fa 100%);
        }
        .email-greeting {
            font-size: 20px;
            font-style: italic;
        }
    `

	mailer3 := mailingo.New(
		product,
		mailingo.Theme{
			PrimaryColor:    "#FF6B6B",
			BackgroundColor: "#FFF5F5",
			TextColor:       "#2C3E50",
			ButtonColor:     "#FF6B6B",
			ButtonTextColor: "#FFFFFF",
		},
		options.WithCustomCSS(extraCSS),
	)

	email3 := mailingo.Email{
		Body: mailingo.Body{
			Name:     "Carol",
			Greeting: "Greetings",
			Title:    "Customized Experience",
			Intros: []string{
				"This email combines custom theme colors with custom CSS.",
			},
			Actions: []mailingo.Action{
				{
					Instructions: "Experience it yourself:",
					Button: mailingo.Button{
						Text: "Try Now",
						Link: "https://customapp.com/try",
					},
				},
			},
			Signature: "Warm regards",
		},
	}

	html3, err := mailer3.GenerateHTML(email3, "en")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("HTML with combined customizations generated successfully!")
	fmt.Println(html3[:300] + "...")
}
