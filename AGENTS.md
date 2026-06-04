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
1. Scaffold the new theme: `cd tools && go run . create <name>` — creates the full directory structure with placeholder `.tmpl` files for all 11 email types
2. Write all 11 `.tmpl` files with unique visual design
3. Run `cd tools && go run . preview all` to regenerate preview data
4. Update README.md style gallery table

### Preview System
- `preview/index.html` loads `preview/rendered.js` (pre-rendered by Go) and displays in iframes
- Supports theme switching, template type switching, client simulation (Modern/Gmail/Outlook/Raw), and viewport toggle (Desktop 1386x780 / Mobile 390x780)
- `REGISTRY` and `PARAMS` are auto-generated from `templates_config.json` by the build tool — no manual syncing needed

### Build Tool
- `tools/tools.go` is the main entry point for the modular CLI
- Subcommands: `list`, `create`, `delete`, `preview`
- Template metadata lives in `tools/data/templates_config.json` — the single source of truth
- `tools/config/` handles config loading and data flattening
- `tools/preview/` implements the rendering engine (template funcs, locale, engine)
- `tools/cli/` implements CLI subcommands using `github.com/urfave/cli/v2`
- Uses Go's native `html/template` package for template rendering

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
