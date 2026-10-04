package email

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"path"
	"strings"
	texttemplate "text/template"
)

//go:embed templates/*
var templateFS embed.FS

const (
	TemplateWelcome       = "welcome"
	TemplatePasswordReset = "password_reset"
	TemplateNotification  = "notification"
	TemplateMemberAdded   = "member_added"
	TemplateMemberInvite  = "member_invite"
)

// Vars are the runtime values a template may interpolate.
type Vars map[string]string

// Map returns a copy for EmailRequest.TemplateData.
func (v Vars) Map() map[string]interface{} {
	out := make(map[string]interface{}, len(v))
	for k, val := range v {
		out[k] = val
	}
	return out
}

// Message is a rendered transactional email ready to send.
type Message struct {
	Subject  string
	HTMLBody string
	TextBody string
}

type parsed struct {
	html *template.Template
	text *texttemplate.Template
}

// Catalog is the set of named email templates loaded from templates/.
type Catalog struct {
	templates map[string]*parsed
}

func parseCatalog() (*Catalog, error) {
	entries, err := fs.ReadDir(templateFS, "templates")
	if err != nil {
		return nil, fmt.Errorf("read email templates: %w", err)
	}
	c := &Catalog{templates: make(map[string]*parsed)}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".html") {
			continue
		}
		name := strings.TrimSuffix(entry.Name(), ".html")
		if name == "layout" || strings.HasPrefix(name, "_") {
			continue
		}
		item, err := parseTemplate(name)
		if err != nil {
			return nil, err
		}
		c.templates[name] = item
	}
	if len(c.templates) == 0 {
		return nil, fmt.Errorf("no email templates found")
	}
	return c, nil
}

func parseTemplate(name string) (*parsed, error) {
	htmlT, err := template.New(name).ParseFS(templateFS, "templates/layout.html", path.Join("templates", name+".html"))
	if err != nil {
		return nil, fmt.Errorf("parse html template %s: %w", name, err)
	}
	if htmlT.Lookup("content") == nil {
		return nil, fmt.Errorf("email template %q missing content define", name)
	}
	textT, err := texttemplate.New(name).ParseFS(templateFS, "templates/_footer.txt", path.Join("templates", name+".txt"))
	if err != nil {
		return nil, fmt.Errorf("parse text template %s: %w", name, err)
	}
	if textT.Lookup("subject") == nil || textT.Lookup("text") == nil {
		return nil, fmt.Errorf("email template %q must define subject and text", name)
	}
	return &parsed{html: htmlT, text: textT}, nil
}

var defaultCatalog = mustCatalog()

func mustCatalog() *Catalog {
	c, err := parseCatalog()
	if err != nil {
		panic(err)
	}
	return c
}

// Render loads a named template from templates/ and interpolates vars.
func Render(name string, vars Vars) (Message, error) {
	return defaultCatalog.Render(name, vars)
}

// Render interpolates vars into the named template.
func (c *Catalog) Render(name string, vars Vars) (Message, error) {
	item, ok := c.templates[name]
	if !ok {
		return Message{}, fmt.Errorf("email template %q not found", name)
	}
	data := cloneVars(vars)
	subject, err := execText(item.text, "subject", data)
	if err != nil {
		return Message{}, fmt.Errorf("render %s subject: %w", name, err)
	}
	data["EmailTitle"] = subject
	htmlBody, err := execHTML(item.html, "layout", data)
	if err != nil {
		return Message{}, fmt.Errorf("render %s html: %w", name, err)
	}
	textBody, err := execText(item.text, "text", data)
	if err != nil {
		return Message{}, fmt.Errorf("render %s text: %w", name, err)
	}
	return Message{Subject: subject, HTMLBody: htmlBody, TextBody: textBody}, nil
}

func cloneVars(vars Vars) Vars {
	out := make(Vars, len(vars)+1)
	for k, v := range vars {
		out[k] = v
	}
	return out
}

func execHTML(t *template.Template, name string, data any) (string, error) {
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, name, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func execText(t *texttemplate.Template, name string, data any) (string, error) {
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, name, data); err != nil {
		return "", err
	}
	return strings.TrimSpace(buf.String()), nil
}
