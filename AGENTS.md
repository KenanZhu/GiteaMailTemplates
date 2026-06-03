# AGENTS.md — Gitea Mail Templates

## Project Overview

A curated collection of email template themes (10 visual styles) for self-hosted Gitea instances. Each theme contains 11 Go `html/template` files covering all Gitea notification email types.

## Repository Layout

```
themes/             # Template themes (5 styles, 11 .tmpl each = 55 source files)
  horizon/          #   Enterprise / Corporate
  terminal/         #   Developers / Tech
  ember/            #   Community / Open Source
  bloom/            #   Creative / Startup (glassmorphism)
  heritage/         #   Education / Research
tools/              # Go build tooling
  build-preview.go  # Pre-renders all templates into preview/rendered.js
  go.mod            # Go module (stdlib only, zero dependencies)
preview/            # Browser-based live preview
  index.html        # SPA with style/template/client/viewport switching
  rendered.js       # Pre-rendered HTML (generated, committed for clone-and-preview)
docs/               # Multi-language documentation
```

## Working With Templates

### Template Files
- All `.tmpl` files use Go `html/template` syntax
- Must use only Gitea's built-in template functions: `AppUrl`, `DotEscape`, `QueryEscape`, `ShortSha`, `HTMLFormat`, `PathEscapeSegments`, `FileSize`
- Must use only official Gitea translation keys (`mail.*` namespace)
- Never reference `.DisplayName` in templates where the data context lacks it (collaborator, transfer, release, workflow_run, assigned, default)
- Each style must have all 11 template types

### Adding a New Theme
1. Create `themes/<name>/` with the full `mail/` directory structure
2. Write all 11 `.tmpl` files with unique visual design
3. Run `go run ./tools/build-preview.go` to regenerate preview data
4. Add the theme to the `<select id="sel-theme">` in `preview/index.html`
5. Update README.md style gallery table

### Preview System
- `preview/index.html` loads `preview/rendered.js` (pre-rendered by Go) and displays in iframes
- Supports theme switching, template type switching, client simulation (Modern/Gmail/Outlook/Raw), and viewport toggle (Desktop 1386x780 / Mobile 390x780)
- Entries in `REGISTRY` and `PARAMS` objects must match the Go build script's template definitions

### Build Script
- `tools/build-preview.go` uses Go's native `html/template` package
- Defines mock data for each template type matching Gitea's actual data contexts
- Post-processes favicon URLs for preview display
- Zero external dependencies (stdlib only)

## Commit Conventions
- `style(name):` — template changes for a specific theme
- `preview:` — preview tooling changes
- `tools:` — Go build script changes
- `docs:` — documentation and translations
- `fix:` — bug fixes
- `project:` — README, LICENSE, AGENTS.md, meta

## Constraints
- No JavaScript framework dependencies — preview is vanilla JS
- No external Go dependencies — build uses stdlib only
- Templates must remain compatible with Gitea's `html/template` execution environment
- Preview works with `file://` protocol (no server needed)
