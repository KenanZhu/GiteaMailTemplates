package preview

import (
	"fmt"
	"html/template"
	"strings"
)

// Locale provides English fallback translations, matching Gitea's locale.Tr interface.
type Locale struct{}

// Tr translates a key and substitutes positional %s/%d arguments.
// HTML arguments are preserved; plain strings are HTML-escaped.
func (l Locale) Tr(key string, args ...any) template.HTML {
	s := LocaleEN[key]
	if s == "" {
		s = key
	}
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

// TrN picks singular/plural by count, then delegates to Tr.
func (l Locale) TrN(count any, singular, plural string, args ...any) template.HTML {
	n := 0
	switch v := count.(type) {
	case int:
		n = v
	case int64:
		n = int(v)
	case float64:
		n = int(v)
	default:
		n = 1
	}
	key := plural
	if n == 1 {
		key = singular
	}
	return l.Tr(key, args...)
}

// LocaleEN holds all English translation strings used by Gitea mail templates.
var LocaleEN = map[string]string{
	"mail.activate_account.title":         "Activate your %s account",
	"mail.hi_user_x":                      "Hi %s,",
	"mail.activate_account.text_1":        "Welcome to %s! Please activate your account by clicking the button below.",
	"mail.activate_account.text_2":        "This activation link will expire in %s.",
	"mail.activate_email.title":           "Verify your email address for %s",
	"mail.activate_email.text":            "Please verify your email address. The link will expire in %s.",
	"mail.register_notify.title":          "Welcome to %s, %s",
	"mail.register_notify.text_1":         "Your account on %s has been created successfully.",
	"mail.register_notify.text_2":         "Your username is: %s",
	"mail.register_notify.text_3":         "If you need to set a password, visit %s.",
	"mail.reset_password.title":           "Reset your %s password",
	"mail.reset_password.text":            "Click the link below to reset your password. It will expire in %s.",
	"mail.link_not_working_do_paste":      "If the button doesn't work, copy and paste this link into your browser:",
	"mail.team_invite.text_1":             "%s has invited you to join the team %s in the organization %s.",
	"mail.team_invite.text_2":             "Click the button below to accept the invitation.",
	"mail.team_invite.text_3":             "This invitation was sent to %s.",
	"mail.repo.collaborator.added.text":   "You have been added as a collaborator on the repository",
	"mail.view_it_on":                     "View it on %s",
	"mail.repo.transfer.body":             "The repository %s has been transferred to you.",
	"mail.release.new.text":               "%s published a new release %s in the repository %s.",
	"mail.release.title":                  "Release: %s",
	"mail.release.note":                   "Release Notes:",
	"mail.release.downloads":              "Downloads",
	"mail.release.download.zip":           "Source Code (ZIP)",
	"mail.release.download.targz":         "Source Code (TAR.GZ)",
	"mail.issue_assigned.pull":            "%s assigned pull request %s to you in the repository %s.",
	"mail.issue_assigned.issue":           "%s assigned issue %s to you in the repository %s.",
	"mail.issue.x_mentioned_you":          "%s mentioned you.",
	"mail.issue.action.close":             "%s closed issue #%d.",
	"mail.issue.action.reopen":            "%s reopened issue #%d.",
	"mail.issue.action.merge":             "%s merged pull request #%d into %s.",
	"mail.issue.action.approve":           "%s approved this change.",
	"mail.issue.action.reject":            "%s requested changes.",
	"mail.issue.action.review":            "%s reviewed this change.",
	"mail.issue.action.review_dismissed":  "%s dismissed the review from %s.",
	"mail.issue.action.ready_for_review":  "%s marked this as ready for review.",
	"mail.issue.action.new":               "%s created issue #%d.",
	"mail.issue.action.force_push":        "%s force-pushed the %s branch from %s to %s.",
	"mail.issue.action.push_1":            "%s pushed %d commit to %s.",
	"mail.issue.action.push_n":            "%s pushed %d commits to %s.",
	"mail.issue.in_tree_path":             "In %s:",
	"mail.reply":                          "Reply to this email",
}
