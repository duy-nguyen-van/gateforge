package email

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRender_PasswordReset(t *testing.T) {
	link := "https://iam.example/reset-password?token=abc"
	msg, err := Render(TemplatePasswordReset, Vars{
		"CTAURL": link,
	})
	require.NoError(t, err)

	assert.Equal(t, "Reset your GateForge password", msg.Subject)
	assert.Contains(t, msg.HTMLBody, "GateForge")
	assert.Contains(t, msg.HTMLBody, "#3451B2")
	assert.Contains(t, msg.HTMLBody, link)
	assert.Contains(t, msg.HTMLBody, "Reset password")
	assert.Contains(t, msg.HTMLBody, `lang="en"`)
	assert.Contains(t, msg.TextBody, link)
	assert.Contains(t, msg.TextBody, "GateForge identity")
}

func TestRender_MemberAdded(t *testing.T) {
	link := "http://localhost:5173/login"
	msg, err := Render(TemplateMemberAdded, Vars{
		"OrgName": "Acme",
		"Role":    "member",
		"CTAURL":  link,
	})
	require.NoError(t, err)
	assert.Equal(t, "You were added to Acme", msg.Subject)
	assert.Contains(t, msg.HTMLBody, link)
	assert.Contains(t, msg.HTMLBody, "member")
	assert.Contains(t, msg.TextBody, "Acme")
	assert.Contains(t, msg.TextBody, "GateForge identity")
}

func TestRender_WelcomeOmitsButton(t *testing.T) {
	msg, err := Render(TemplateWelcome, Vars{"UserName": "Ada"})
	require.NoError(t, err)

	assert.Equal(t, "Welcome to GateForge", msg.Subject)
	assert.Contains(t, msg.TextBody, "Hello Ada.")
	assert.NotContains(t, msg.HTMLBody, "href=")
	assert.Contains(t, msg.TextBody, "GateForge identity")
}

func TestRender_EscapesHTML(t *testing.T) {
	msg, err := Render(TemplateNotification, Vars{
		"Subject": "Notice",
		"Message": `<script>alert("x")</script>`,
	})
	require.NoError(t, err)

	assert.Equal(t, "Notice", msg.Subject)
	assert.NotContains(t, msg.HTMLBody, "<script>alert")
	assert.Contains(t, msg.HTMLBody, "&lt;script&gt;")
	assert.Contains(t, msg.TextBody, `<script>alert("x")</script>`)
}

func TestRender_UnknownTemplate(t *testing.T) {
	_, err := Render("missing", nil)
	require.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "not found"))
}
