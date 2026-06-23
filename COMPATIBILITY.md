# Gitea Compatibility

This document tracks the compatibility between **Gitea Mail Templates** releases and **Gitea** versions.

## Quick Reference

| Template Release | Min Gitea | Max Tested Gitea | Status |
|-----------------|-----------| **1.26.4** | ⏳ Pending Verification |
<!-- TRACKER:QUICK-REF-MAX -->
| **v1.0.x**      | **1.25.0** | **1.26.4** | ✅ Active |

> **Latest verified:** All 11 templates pass validation against Gitea 1.26.4 data contexts. <!-- TRACKER:LATEST-VERIFIED -->

## Check Your Gitea Version

```bash
# On your Gitea server:
gitea --version
# Or check the web UI footer / Site Administration → Monitoring
```

## Gitea Version History — Mail Template Impact

| Gitea | Release Date | Mail Template Changes | Breaking? |
|-------|-------------|----------------------|-----------|
<!-- TRACKER:VERSION-INSERT -->
| **1.26.4** | 2026-06-23 | ⏳ Pending Verification | TBD |
| **1.26.4** | 2026-06-21 | None — hotfix release | No |
| **1.26.3** | 2026-06-20 | None — security release | No |
| **1.26.2** | 2026-05-20 | None — security + bug fixes | No |
| **1.26.1** | 2026-04-22 | None — bug fixes | No |
| **1.26.0** | 2026-04-19 | AppURL cleanup; SanitizeHTML deprecated → use HTMLFormat | No |
| **1.25.5** | 2026-03-10 | None — security + maintenance | No |
| **1.25.0** | 2025 | **Directory restructure** — templates moved to `mail/<category>/<type>.tmpl` (PR #35150); subject/body split with `---` separator; template preview support added | **Yes** (structural) |
| **≤ 1.24.x** | — | Flat directory structure under `custom/templates/mail/` | ❌ Unsupported |

## Template Variable Reference

All 11 template types use only Gitea built-in variables and functions. Verified against Gitea source (`services/mailer/` + `modules/templates/mail.go`).

### Template Functions (available in all templates)

| Function | Since Gitea | Notes |
|----------|------------|-------|
| `AppName` | ≤ 1.21 | Application name |
| `AppUrl` | ≤ 1.21 | Application base URL |
| `AppDomain` | ≤ 1.21 | Server domain |
| `DotEscape` | ≤ 1.21 | Prevents auto-linking of dotted text |
| `QueryEscape` | ≤ 1.21 | URL query encoding |
| `PathEscapeSegments` | ≤ 1.21 | Per-segment path encoding |
| `ShortSha` | ≤ 1.21 | Truncated commit hash |
| `FileSize` | ≤ 1.21 | Human-readable file size |
| `HTMLFormat` | ≤ 1.21 | Render string as safe HTML |
| `Iif` | ≤ 1.21 | Inline conditional |
| `dict` | ≤ 1.21 | Build maps from key-value pairs |
| `Eval` | ≤ 1.21 | Evaluate template tokens |
| `StringUtils` | ≤ 1.21 | String manipulation helpers |
| `SliceUtils` | ≤ 1.21 | Slice manipulation helpers |
| `JsonUtils` | ≤ 1.21 | JSON helpers |
| `DumpVar` | ≤ 1.21 | Debug variable dump |
| `SanitizeHTML` | ≤ 1.21 | **Deprecated in 1.26** — use `HTMLFormat` |

### Data Contexts by Template

| Template | Key Variables |
|----------|--------------|
| `user/auth/activate` | `DisplayName`, `Code`, `ActiveCodeLives` |
| `user/auth/activate_email` | `DisplayName`, `Code`, `Email`, `ActiveCodeLives` |
| `user/auth/register_notify` | `DisplayName`, `Username` |
| `user/auth/reset_passwd` | `DisplayName`, `Code`, `ResetPwdCodeLives` |
| `org/team_invite` | `Inviter`, `Team`, `Organization`, `InviteURL`, `Invite` |
| `repo/collaborator` | `Subject`, `RepoName`, `Link` |
| `repo/transfer` | `Doer`, `User`, `Repo`, `Link`, `Destination` |
| `repo/release` | `Release` (with `Publisher`, `TagName`, `Title`, `RenderedNote`, `Attachments`), `Link` |
| `repo/actions/workflow_run` | `Run` (with `WorkflowID`, `HTMLURL`), `Jobs`, `RunStatusText` |
| `repo/issue/assigned` | `Doer`, `Issue`, `Link`, `IsPull` |
| `repo/issue/default` | `Doer`, `Issue`, `Link`, `Body`, `ActionName`, `Comment`, `IsPull`, `IsMention`, `ReviewComments`, `CanReply` |

> ⚠️ **`.DisplayName`** is not available in collaborator, transfer, release, workflow_run, assigned, and default templates — do not reference it.

### Translation Keys

All templates use Gitea's official `mail.*` translation namespace. Keys are stable across Gitea 1.25+.

## How Compatibility Is Verified

1. **Automated lint** — CI renders all templates via `go run . preview all` on every push
2. **Source audit** — Template data contexts are cross-referenced against Gitea's `services/mailer/` package
3. **Release checklist** — Each release confirms the max-tested Gitea version in this file

## Automated Tracking

A [workflow](.github/workflows/gitea-tracker.yml) runs daily (UTC 08:00) to detect new Gitea releases. When a new version is found, it automatically creates a PR updating the matrix and badges with a **⏳ Pending Verification** status. Manual trigger is also available via `workflow_dispatch`.

## Reporting Issues

If you find a compatibility problem with a specific Gitea version:
1. Check the [Gitea changelog](https://github.com/go-gitea/gitea/blob/main/CHANGELOG.md) for recent mail template changes
2. Open an issue with: your Gitea version, which template, and the error
