// Package options provides functional options for configuring Mailer instances.
package options

import (
	"html/template"
	"io/fs"
)

// Option is a function that configures a Mailer.
type Option func(*Config)

// Config holds the configuration for a Mailer.
type Config struct {
	CustomTemplate     *template.Template
	CustomTemplateText string
	CustomTemplateFS   fs.FS
	CustomTemplatePath string
	CustomCSS          string
}

// WithCustomTemplate allows you to provide your own HTML template.
// The template should use the same data structure as the default template.
//
// Example:
//
//	tmpl := template.Must(template.New("email").Parse(`<html>...</html>`))
//	mailer := mailingo.New(product, theme, options.WithCustomTemplate(tmpl))
func WithCustomTemplate(tmpl *template.Template) Option {
	return func(c *Config) {
		c.CustomTemplate = tmpl
	}
}

// WithCustomTemplateString allows you to provide your own HTML template as a string.
// The template string will be parsed when creating the Mailer.
//
// Example:
//
//	templateStr := `<!DOCTYPE html><html>...</html>`
//	mailer := mailingo.New(product, theme, options.WithCustomTemplateString(templateStr))
func WithCustomTemplateString(tmplStr string) Option {
	return func(c *Config) {
		c.CustomTemplateText = tmplStr
	}
}

// WithCustomTemplateFile allows you to load a custom template from an embedded filesystem.
// This is useful when you want to embed your custom templates using go:embed.
//
// Example:
//
//	//go:embed templates/*.html
//	var templatesFS embed.FS
//
//	mailer := mailingo.New(product, theme,
//	    options.WithCustomTemplateFile(templatesFS, "templates/mytemplate.html"))
func WithCustomTemplateFile(filesystem fs.FS, path string) Option {
	return func(c *Config) {
		c.CustomTemplateFS = filesystem
		c.CustomTemplatePath = path
	}
}

// WithCustomCSS allows you to add custom CSS that will be appended to the default styles.
// This is useful for small style tweaks without replacing the entire template.
//
// Example:
//
//	customCSS := `
//	    .email-title { font-size: 28px; }
//	    .email-button { border-radius: 8px; }
//	`
//	mailer := mailingo.New(product, theme, options.WithCustomCSS(customCSS))
func WithCustomCSS(css string) Option {
	return func(c *Config) {
		c.CustomCSS = css
	}
}
