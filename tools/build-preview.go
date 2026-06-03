package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// TemplateFuncs maps Gitea's built-in template functions
var TemplateFuncs = template.FuncMap{
	"AppUrl":              func() string { return "https://about.gitea.com/" },
	"AppName":             func() string { return "Gitea" },
	"DotEscape":           func(s string) template.HTML { return template.HTML(template.HTMLEscapeString(s)) },
	"QueryEscape":         url.QueryEscape,
	"ShortSha":            ShortSha,
	"PathEscapeSegments":  PathEscapeSegments,
	"FileSize":            FileSize,
	"HTMLFormat":          HTMLFormat,
	"Dict":                Dict,
}

func ShortSha(s string) string {
	if len(s) > 10 { return s[:10] }
	return s
}

func PathEscapeSegments(s string) string {
	parts := strings.Split(s, "/")
	for i, p := range parts { parts[i] = url.PathEscape(p) }
	return strings.Join(parts, "/")
}

func FileSize(n int64) string {
	const unit = 1024
	if n < unit { return fmt.Sprintf("%d B", n) }
	div, exp := int64(unit), 0
	for nn := n / unit; nn >= unit && exp < 3; nn /= unit { div *= unit; exp++ }
	return fmt.Sprintf("%.1f %s", float64(n)/float64(div), []string{"KB","MB","GB","TB"}[exp])
}

func HTMLFormat(fmtStr string, args ...any) template.HTML {
	result := fmtStr
	for i, arg := range args {
		escaped := template.HTMLEscapeString(fmt.Sprint(arg))
		result = strings.Replace(result, fmt.Sprintf("%%[%d]s", i+1), escaped, 1)
		if !strings.Contains(result, fmt.Sprintf("%%[%d]s", i+1)) {
			result = strings.Replace(result, "%s", escaped, 1)
			result = strings.Replace(result, "%d", fmt.Sprint(arg), 1)
		}
	}
	return template.HTML(result)
}

func Dict(values ...any) (map[string]any, error) {
	if len(values)%2 != 0 { return nil, fmt.Errorf("odd number of Dict args") }
	dict := make(map[string]any, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok { return nil, fmt.Errorf("Dict keys must be strings") }
		dict[key] = values[i+1]
	}
	return dict, nil
}

// Locale provides English fallback translations
type Locale struct{}

func (l Locale) Tr(key string, args ...any) template.HTML {
	s := LocaleEN[key]
	if s == "" { s = key }
	for _, arg := range args {
		var val string
		if h, ok := arg.(template.HTML); ok {
			val = string(h) // already safe HTML — don't double-escape
		} else {
			val = template.HTMLEscapeString(fmt.Sprint(arg))
		}
		s = strings.Replace(s, "%s", val, 1)
		s = strings.Replace(s, "%d", fmt.Sprint(arg), 1)
	}
	return template.HTML(s)
}

func (l Locale) TrN(count any, singular, plural string, args ...any) template.HTML {
	n := 0
	switch v := count.(type) {
	case int: n = v
	case int64: n = int(v)
	case float64: n = int(v)
	default: n = 1
	}
	key := plural
	if n == 1 { key = singular }
	return l.Tr(key, args...)
}

var LocaleEN = map[string]string{
	"mail.activate_account.title":    "Activate your %s account",
	"mail.hi_user_x":                 "Hi %s,",
	"mail.activate_account.text_1":   "Welcome to %s! Please activate your account by clicking the button below.",
	"mail.activate_account.text_2":   "This activation link will expire in %s.",
	"mail.activate_email.title":      "Verify your email address for %s",
	"mail.activate_email.text":       "Please verify your email address. The link will expire in %s.",
	"mail.register_notify.title":     "Welcome to %s, %s",
	"mail.register_notify.text_1":    "Your account on %s has been created successfully.",
	"mail.register_notify.text_2":    "Your username is: %s",
	"mail.register_notify.text_3":    "If you need to set a password, visit %s.",
	"mail.reset_password.title":      "Reset your %s password",
	"mail.reset_password.text":       "Click the link below to reset your password. It will expire in %s.",
	"mail.link_not_working_do_paste": "If the button doesn't work, copy and paste this link into your browser:",
	"mail.team_invite.text_1":        "%s has invited you to join the team %s in the organization %s.",
	"mail.team_invite.text_2":        "Click the button below to accept the invitation.",
	"mail.team_invite.text_3":        "This invitation was sent to %s.",
	"mail.repo.collaborator.added.text": "You have been added as a collaborator on the repository",
	"mail.view_it_on":                "View it on %s",
	"mail.repo.transfer.body":        "The repository %s has been transferred to you.",
	"mail.release.new.text":          "%s published a new release %s in the repository %s.",
	"mail.release.title":             "Release: %s",
	"mail.release.note":              "Release Notes:",
	"mail.release.downloads":         "Downloads",
	"mail.release.download.zip":      "Source Code (ZIP)",
	"mail.release.download.targz":    "Source Code (TAR.GZ)",
	"mail.issue_assigned.pull":       "%s assigned pull request %s to you in the repository %s.",
	"mail.issue_assigned.issue":      "%s assigned issue %s to you in the repository %s.",
	"mail.issue.x_mentioned_you":     "%s mentioned you.",
	"mail.issue.action.close":        "%s closed issue #%d.",
	"mail.issue.action.reopen":       "%s reopened issue #%d.",
	"mail.issue.action.merge":        "%s merged pull request #%d into %s.",
	"mail.issue.action.approve":      "%s approved this change.",
	"mail.issue.action.reject":       "%s requested changes.",
	"mail.issue.action.review":       "%s reviewed this change.",
	"mail.issue.action.review_dismissed": "%s dismissed the review from %s.",
	"mail.issue.action.ready_for_review": "%s marked this as ready for review.",
	"mail.issue.action.new":          "%s created issue #%d.",
	"mail.issue.action.force_push":   "%s force-pushed the %s branch from %s to %s.",
	"mail.issue.action.push_1":       "%s pushed %d commit to %s.",
	"mail.issue.action.push_n":       "%s pushed %d commits to %s.",
	"mail.issue.in_tree_path":        "In %s:",
	"mail.reply":                    "Reply to this email",
}

// MockData for each template type
type TemplateDef struct {
	Path   string
	Data   map[string]any
}

var templates = map[string]TemplateDef{
	"activate": {
		Path: "mail/user/auth/activate.tmpl",
		Data: map[string]any{
			"DisplayName":     "Alex Johnson",
			"Code":            "activation-code-abc123",
			"ActiveCodeLives": "24 hours",
		},
	},
	"activate_email": {
		Path: "mail/user/auth/activate_email.tmpl",
		Data: map[string]any{
			"DisplayName":     "Alex Johnson",
			"Code":            "verify-code-xyz789",
			"Email":           "alex@example.com",
			"ActiveCodeLives": "24 hours",
		},
	},
	"register_notify": {
		Path: "mail/user/auth/register_notify.tmpl",
		Data: map[string]any{
			"DisplayName": "Alex Johnson",
			"Username":    "alexj",
		},
	},
	"reset_passwd": {
		Path: "mail/user/auth/reset_passwd.tmpl",
		Data: map[string]any{
			"DisplayName":       "Alex Johnson",
			"Code":              "reset-token-abc123",
			"ResetPwdCodeLives": "3 hours",
		},
	},
	"team_invite": {
		Path: "mail/org/team_invite.tmpl",
		Data: map[string]any{
			"Inviter":      map[string]any{"DisplayName": "Sarah Chen"},
			"Team":         map[string]any{"Name": "Core Developers"},
			"Organization": map[string]any{"DisplayName": "Acme Corp"},
			"InviteURL":    "https://gitea.example.com/org/team/invite?token=invite-abc",
			"Invite":       map[string]any{"Email": "alex@example.com"},
		},
	},
	"collaborator": {
		Path: "mail/repo/collaborator.tmpl",
		Data: map[string]any{
			"Subject":  "Collaborator Added to repo",
			"Link":     "https://gitea.example.com/repo/project",
			"RepoName": "acme/project",
		},
	},
	"transfer": {
		Path: "mail/repo/transfer.tmpl",
		Data: map[string]any{
			"Subject": "Repository Transfer",
			"Link":    "https://gitea.example.com/repo/transferred-project",
			"Repo":    "acme/transferred-project",
		},
	},
	"release": {
		Path: "mail/repo/release.tmpl",
		Data: map[string]any{
			"Subject": "New Release v2.0.0",
			"Link":    "https://gitea.example.com/repo/project/releases/tag/v2.0.0",
			"Release": map[string]any{
				"Publisher":    map[string]any{"Name": "Sarah Chen"},
				"HTMLURL":      "https://gitea.example.com/repo/project/releases/tag/v2.0.0",
				"TagName":      "v2.0.0",
				"Title":        "Version 2.0 -- Major Update",
				"RenderedNote": "<p>Major new features and performance improvements in this release.</p>",
				"Repo": map[string]any{
					"HTMLURL":  "https://gitea.example.com/repo/project",
					"FullName": "acme/project",
				},
				"Attachments": []any{},
			},
			"DisableDownloadSourceArchives": false,
		},
	},
	"workflow_run": {
		Path: "mail/repo/actions/workflow_run.tmpl",
		Data: map[string]any{
			"Subject":      "Workflow Run: CI",
			"RunStatusText": "All jobs passed",
			"Repo":          map[string]any{"FullName": "acme/project"},
			"Run": map[string]any{
				"WorkflowID": "ci.yml",
				"HTMLURL":    "https://gitea.example.com/repo/project/actions/runs/42",
			},
			"Jobs": []map[string]any{
				{"Name": "build", "Status": "success", "Attempt": 1, "HTMLURL": "https://gitea.example.com/job/1", "Duration": "2m34s"},
				{"Name": "test", "Status": "success", "Attempt": 1, "HTMLURL": "https://gitea.example.com/job/2", "Duration": "1m12s"},
				{"Name": "deploy", "Status": "success", "Attempt": 2, "HTMLURL": "https://gitea.example.com/job/3", "Duration": "45s"},
			},
		},
	},
	"assigned": {
		Path: "mail/repo/issue/assigned.tmpl",
		Data: map[string]any{
			"Subject": "Issue #42 Assigned",
			"Link":    "https://gitea.example.com/repo/project/issues/42",
			"Doer":    map[string]any{"Name": "Sarah Chen"},
			"IsPull":  false,
			"Issue": map[string]any{
				"Index": 42,
				"Repo": map[string]any{
					"HTMLURL":  "https://gitea.example.com/repo/project",
					"FullName": "acme/project",
				},
			},
		},
	},
	"default": {
		Path: "mail/repo/issue/default.tmpl",
		Data: map[string]any{
			"Subject":    "Re: Issue #42 -- Bug fix",
			"Link":       "https://gitea.example.com/repo/project/issues/42",
			"ActionName": "comment",
			"Doer":       map[string]any{"Name": "Sarah Chen"},
			"Body":       "This is a comment on the issue with details about the fix.",
			"CanReply":   true,
			"IsMention":  false,
			"IsPull":     false,
			"Issue": map[string]any{
				"Index": 42,
				"PullRequest": map[string]any{
					"BaseBranch": "main",
					"BaseRepo": map[string]any{
						"HTMLURL": "https://gitea.example.com/repo/project",
					},
					"HeadBranch": "feature/fix",
				},
			},
			"Comment": map[string]any{
				"Commits":    []any{},
				"IsForcePush": false,
				"Issue": map[string]any{
					"PullRequest": map[string]any{
						"BaseRepo": map[string]any{
							"HTMLURL": "https://gitea.example.com/repo/project",
						},
						"HeadBranch": "feature/fix",
					},
				},
			},
			"ReviewComments": []any{},
		},
	},
}

// discoverThemes scans the themes directory and returns all subdirectory names.
func discoverThemes(themesDir string) []string {
	entries, err := os.ReadDir(themesDir)
	if err != nil {
		log.Fatalf("Cannot read themes directory %s: %v", themesDir, err)
	}
	var themes []string
	for _, e := range entries {
		if e.IsDir() {
			themes = append(themes, e.Name())
		}
	}
	if len(themes) == 0 {
		log.Fatalf("No theme directories found in %s", themesDir)
	}
	return themes
}

func main() {
	exe, _ := os.Executable()
	rootDir := filepath.Dir(exe)
	// When run with `go run`, use CWD's parent
	if _, err := os.Stat(filepath.Join(rootDir, "themes")); os.IsNotExist(err) {
		cwd, _ := os.Getwd()
		// Try cwd, then cwd's parent
		for _, d := range []string{cwd, filepath.Dir(cwd)} {
			if _, e := os.Stat(filepath.Join(d, "themes")); e == nil { rootDir = d; break }
		}
	}

	themesDir := filepath.Join(rootDir, "themes")
	previewDir := filepath.Join(rootDir, "preview")

	result := make(map[string]map[string]string)

	themes := discoverThemes(themesDir)
	for _, theme := range themes {
		themeDir := filepath.Join(themesDir, theme)
		result[theme] = make(map[string]string)

		for tplID, tplDef := range templates {
			tmplPath := filepath.Join(themeDir, tplDef.Path)
			tmplContent, err := os.ReadFile(tmplPath)
			if err != nil {
				log.Printf("WARN: cannot read %s: %v", tmplPath, err)
				result[theme][tplID] = fmt.Sprintf("<p style='color:red'>Failed to read template: %s</p>", err)
				continue
			}

			// Build data context
			data := make(map[string]any)
			for k, v := range tplDef.Data {
				data[k] = v
			}
			data["locale"] = Locale{}

			// Parse and execute
			tmpl, err := template.New(tplID).Funcs(TemplateFuncs).Parse(string(tmplContent))
			if err != nil {
				log.Printf("WARN: parse error %s/%s: %v", theme, tplID, err)
				result[theme][tplID] = fmt.Sprintf("<p style='color:red'>Parse error: %s</p>", err)
				continue
			}

			var buf strings.Builder
			if err := tmpl.Execute(&buf, data); err != nil {
				log.Printf("WARN: execute error %s/%s: %v", theme, tplID, err)
				result[theme][tplID] = fmt.Sprintf("<p style='color:red'>Render error: %s</p>", err)
				continue
			}

			// Fix favicon path for preview (official Gitea icon)
			rendered := strings.Replace(buf.String(),
				`src="https://about.gitea.com/assets/img/favicon.png"`,
				`src="https://about.gitea.com/gitea.svg"`, -1)
			result[theme][tplID] = rendered
			log.Printf("OK: %s/%s (%d bytes)", theme, tplID, buf.Len())
		}
	}

	// Write output as JS (script-loadable, avoids file:// CORS issues)
	os.MkdirAll(previewDir, 0755)
	outPath := filepath.Join(previewDir, "rendered.js")
	outFile, err := os.Create(outPath)
	if err != nil { log.Fatalf("Cannot create output: %v", err) }
	defer outFile.Close()

	jsonBytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil { log.Fatalf("JSON marshal error: %v", err) }

	outFile.WriteString("// Auto-generated by tools/build-preview.go — do not edit\n")
	outFile.WriteString("window.__RENDERED__ = ")
	outFile.Write(jsonBytes)
	outFile.WriteString(";\n")

	log.Printf("Wrote %d themes x %d templates to %s", len(themes), len(templates), outPath)
}
