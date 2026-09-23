package preview

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gitea-mail-templates/tools/config"
)

type testCommitID string

func (id testCommitID) String() string { return string(id) }

// Exercise the push-only branch with Gitea's nested commit data shape.
// The normal preview fixture has no commits, so it cannot catch stale .ID paths.
func TestPushCommitTemplates(t *testing.T) {
	cfg, err := config.Load(filepath.Join("..", "data", "templates_config.json"))
	if err != nil {
		t.Fatal(err)
	}
	const sha = "0123456789abcdef0123456789abcdef01234567"
	data := cfg.Templates["default"].Params
	data["ActionName"] = "push"
	data["locale"] = Locale{}
	comment := data["Comment"].(map[string]any)
	comment["Commits"] = []any{map[string]any{
		"UserCommit": map[string]any{
			"GitCommit": map[string]any{
				"ID":           testCommitID(sha),
				"MessageTitle": "Test push commit",
			},
		},
	}}
	markSafeHTML(data)

	themesDir := filepath.Join("..", "..", "themes")
	themes, err := os.ReadDir(themesDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, theme := range themes {
		if !theme.IsDir() {
			continue
		}
		t.Run(theme.Name(), func(t *testing.T) {
			path := filepath.Join(themesDir, theme.Name(), "mail", "repo", "issue", "default.tmpl")
			content, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			tmpl, err := template.New("push").Funcs(TemplateFuncs).Parse(string(content))
			if err != nil {
				t.Fatal(err)
			}
			var html strings.Builder
			if err := tmpl.Execute(&html, data); err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(html.String(), "/commit/"+sha) || !strings.Contains(html.String(), sha[:10]) {
				t.Fatal("push commit link or abbreviated hash missing from rendered email")
			}
		})
	}
}
