# AGENTS.md — Gitea Mail Templates

## Project Overview

A curated collection of email template themes (10 visual styles) for self-hosted Gitea instances. Each theme contains 11 Go `html/template` files covering all Gitea notification email types.

## Repository Layout

```
themes/             # Template themes (10 styles, 11 .tmpl each = 110 source files)
  aurora/           #   Ethereal / Dreamlike
  bloom/            #   Creative / Startup (glassmorphism)
  ember/            #   Community / Open Source
  heritage/         #   Education / Research
  horizon/          #   Enterprise / Corporate
  ink/              #   Editorial / Publishing
  mono/             #   Minimal / Swiss design
  neon/             #   Cyberpunk / Gaming
  terminal/         #   Developers / Tech
  terra/            #   Nature / Sustainability
tools/              # Go CLI tooling (modular, zero dependencies)
  tools.go          #   Main entry point
  cli/              #   CLI subcommands: list, create, delete, preview
  config/           #   Config types and templates_config.json loading
  data/             #   templates_config.json — single source of truth for template metadata
  preview/          #   Template rendering engine (funcs, locale, engine)
  go.mod            #   Go module (stdlib only, zero dependencies)
preview/            # Browser-based live preview
  index.html        # SPA with style/template/client/viewport switching
  rendered.js       # Pre-rendered HTML (generated, committed for clone-and-preview)
docs/               # Bilingual documentation (English + Simplified Chinese)
```

## Working With Templates

### Template Files
- All `.tmpl` files use Go `html/template` syntax
- Must use only Gitea's built-in template functions: `AppUrl`, `DotEscape`, `QueryEscape`, `ShortSha`, `HTMLFormat`, `PathEscapeSegments`, `FileSize`
- Must use only official Gitea translation keys (`mail.*` namespace)
- Never reference `.DisplayName` in templates where the data context lacks it (collaborator, transfer, release, workflow_run, assigned, default)
- Each style must have all 11 template types

### Adding a New Theme
1. Scaffold the new theme: `cd tools && go run . create <name>` — creates the full directory structure with placeholder `.tmpl` files for all 11 email types
2. Write all 11 `.tmpl` files with unique visual design
3. Run `cd tools && go run . preview all` to regenerate preview data
4. Update README.md style gallery table

### Preview System
- `preview/index.html` loads `preview/rendered.js` (pre-rendered by Go) and displays in iframes
- Supports theme/template switching, view mode (Modern/Source), and viewport toggle (Desktop 1386x780 / Mobile 390x780)
- Keyboard navigation: `←→` cycles focus between Theme/Template/View selects, `↑↓` selects within the focused dropdown, `d`/`m` toggles viewport
- `REGISTRY` and `PARAMS` are auto-generated from `templates_config.json` — no manual syncing needed
- Static preview (open `index.html` directly) — works via `file://` protocol with Modern and Source views
- Dev server (`go run . dev`) — pure Go HTTP server with SSE live reload; watches `themes/` for `.tmpl` changes, re-renders in-process, and pushes reload events to the browser

### Build Tool
- `tools/tools.go` is the main entry point for the modular CLI
- Subcommands: `list`, `create`, `delete`, `preview`, `dev`
- Template metadata lives in `tools/data/templates_config.json` — the single source of truth
- `tools/config/` handles config loading and data flattening
- `tools/preview/` implements the rendering engine (template funcs, locale, engine, markSafeHTML)
- `tools/cli/` implements CLI subcommands using `github.com/urfave/cli/v2`
- `tools/preview/server.go` pure Go dev server with SSE live reload and in-process template re-rendering
- Uses Go's native `html/template` package for template rendering

## Versioning

- Release numbers stay on the project's own scheme; the supported Gitea version is appended in parentheses — `v1.0.1(v1.27.2)` means release `v1.0.1` tracks Gitea `1.27.2`
- The current release is **v1.0.1** (tracks Gitea 1.27.2); the quick-reference table in `COMPATIBILITY.md` lists it first — the tracker workflow updates that row by position
- When a new Gitea version is verified compatible: update the parenthesized Gitea version in the active `COMPATIBILITY.md` row, plus the README release line and badge
- Tag a new release (`vX.Y.Z`) only when the template content itself changes — the release workflow packages automatically on tag push
- Keep the `TRACKER:` markers in `COMPATIBILITY.md` / `README.md` adjacent to their rows so the automated workflow keeps parsing them

## Commit Conventions
- `style(name):` — template changes for a specific theme
- `preview:` — preview tooling changes
- `tools:` — Go CLI/build tooling changes
- `docs:` — documentation and translations
- `fix:` — bug fixes
- `project:` — README, LICENSE, AGENTS.md, meta
- `refactor:` — code restructuring (e.g. modularization)
- `chore:` — maintenance (config updates, build scripts)

## Constraints
- No JavaScript framework dependencies — preview is vanilla JS
- No external Go dependencies — build uses stdlib only
- Templates must remain compatible with Gitea's `html/template` execution environment
- Preview works with `file://` protocol (no server needed)
