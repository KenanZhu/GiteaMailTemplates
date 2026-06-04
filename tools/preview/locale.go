package preview

import (
	"fmt"
	"html/template"
	"strings"
)

// Locale provides English fallback translations, matching Gitea's locale.Tr interface.
type Locale struct{}

// Tr translates a key and substitutes %s / %d / %[N]s / %[N]d arguments.
// HTML arguments are preserved; plain strings are HTML-escaped.
func (l Locale) Tr(key string, args ...any) template.HTML {
	s := LocaleEN[key]
	if s == "" {
		s = key
	}
	for i, arg := range args {
		var val string
		if h, ok := arg.(template.HTML); ok {
			val = string(h) // already safe HTML
		} else {
			val = template.HTMLEscapeString(fmt.Sprint(arg))
		}
		// Handle positional format specifiers: %[1]s, %[2]d, etc.
		pos := fmt.Sprintf("%%[%d]", i+1)
		s = strings.Replace(s, pos+"s", val, 1)
		s = strings.Replace(s, pos+"d", val, 1)
	}
	// Fallback for non-positional %s / %d (sequential replacement)
	for _, arg := range args {
		var val string
		if h, ok := arg.(template.HTML); ok {
			val = string(h)
		} else {
			val = template.HTMLEscapeString(fmt.Sprint(arg))
		}
		if !strings.Contains(s, "%s") && !strings.Contains(s, "%d") {
			break
		}
		s = strings.Replace(s, "%s", val, 1)
		s = strings.Replace(s, "%d", val, 1)
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
	"mail.hi_user_x":                         "Hi <b>%s</b>,",
	"mail.activate_account":                  "Please activate your account",
	"mail.activate_account.title":            "%s, please activate your account",
	"mail.activate_account.text_1":           "Hi <b>%[1]s</b>, thanks for registering at %[2]s!",
	"mail.activate_account.text_2":           "Please click the following link to activate your account within <b>%s</b>:",
	"mail.activate_email":                    "Verify your email address",
	"mail.activate_email.title":              "%s, please verify your email address",
	"mail.activate_email.text":               "Please click the following link to verify your email address within <b>%s</b>:",
	"mail.register_success":                  "Registration successful",
	"mail.register_notify":                   "Welcome to %s",
	"mail.register_notify.title":             "%[1]s, welcome to %[2]s",
	"mail.register_notify.text_1":            "This is your registration confirmation email for %s!",
	"mail.register_notify.text_2":            "You can now log in via username: %s.",
	"mail.register_notify.text_3":            "If this account has been created for you, please <a href=\"%s\">set your password</a> first.",
	"mail.reset_password":                    "Recover your account",
	"mail.reset_password.title":              "%s, you have requested to recover your account",
	"mail.reset_password.text":               "Please click the following link to recover your account within <b>%s</b>:",
	"mail.team_invite.subject":               "%[1]s has invited you to join the %[2]s organization",
	"mail.team_invite.text_1":                "%[1]s has invited you to join team %[2]s in organization %[3]s.",
	"mail.team_invite.text_2":                "Please click the following link to join the team:",
	"mail.team_invite.text_3":                "Note: This invitation was intended for %[1]s. If you were not expecting this invitation, you can ignore this email.",
	"mail.repo.collaborator.added.subject":   "%s added you to %s",
	"mail.repo.collaborator.added.text":      "You have been added as a collaborator of repository:",
	"mail.repo.transfer.subject_to":          "%s would like to transfer \"%s\" to %s",
	"mail.repo.transfer.subject_to_you":      "%s would like to transfer \"%s\" to you",
	"mail.repo.transfer.to_you":              "you",
	"mail.repo.transfer.body":                "To accept or reject it, visit %s or just ignore it.",
	"mail.release.new.subject":               "%s in %s released",
	"mail.release.new.text":                  "<b>@%[1]s</b> released %[2]s in %[3]s",
	"mail.release.title":                     "Title: %s",
	"mail.release.note":                      "Note:",
	"mail.release.downloads":                 "Downloads:",
	"mail.release.download.zip":              "Source Code (ZIP)",
	"mail.release.download.targz":            "Source Code (TAR.GZ)",
	"mail.repo.actions.run.failed":           "Run failed",
	"mail.repo.actions.run.succeeded":        "Run succeeded",
	"mail.repo.actions.run.cancelled":        "Run cancelled",
	"mail.repo.actions.jobs.all_succeeded":   "All jobs have succeeded",
	"mail.repo.actions.jobs.all_failed":      "All jobs have failed",
	"mail.repo.actions.jobs.some_not_successful": "Some jobs were not successful",
	"mail.repo.actions.jobs.all_cancelled":       "All jobs have been cancelled",
	"mail.issue_assigned.pull":               "@%[1]s assigned you to pull request %[2]s in repository %[3]s.",
	"mail.issue_assigned.issue":              "@%[1]s assigned you to issue %[2]s in repository %[3]s.",
	"mail.issue.x_mentioned_you":             "<b>@%s</b> mentioned you:",
	"mail.issue.action.close":                "<b>@%[1]s</b> closed #%[2]d.",
	"mail.issue.action.reopen":               "<b>@%[1]s</b> reopened #%[2]d.",
	"mail.issue.action.merge":                "<b>@%[1]s</b> merged #%[2]d into %[3]s.",
	"mail.issue.action.approve":              "<b>@%[1]s</b> approved this pull request.",
	"mail.issue.action.reject":               "<b>@%[1]s</b> requested changes on this pull request.",
	"mail.issue.action.review":               "<b>@%[1]s</b> commented on this pull request.",
	"mail.issue.action.review_dismissed":     "<b>@%[1]s</b> dismissed last review from %[2]s for this pull request.",
	"mail.issue.action.ready_for_review":     "<b>@%[1]s</b> marked this pull request ready for review.",
	"mail.issue.action.new":                  "<b>@%[1]s</b> created #%[2]d.",
	"mail.issue.action.force_push":           "<b>%[1]s</b> force-pushed the <b>%[2]s</b> from %[3]s to %[4]s.",
	"mail.issue.action.push_1":               "<b>@%[1]s</b> pushed %[3]d commit to %[2]s",
	"mail.issue.action.push_n":               "<b>@%[1]s</b> pushed %[3]d commits to %[2]s",
	"mail.issue.in_tree_path":                "In %s:",
	"mail.view_it_on":                        "View it on %s",
	"mail.reply":                             "or reply to this email directly",
	"mail.link_not_working_do_paste":         "Not working? Try copying and pasting it to your browser.",
}
