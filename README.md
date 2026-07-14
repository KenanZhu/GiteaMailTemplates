# Gitea Mail Templates

A curated collection of professionally designed, audience-driven email templates for self-hosted [Gitea](https://about.gitea.com) instances.

[![Gitea](https://img.shields.io/badge/Gitea-1.25+%20%7C%201.27.0%20pending-yellow)](COMPATIBILITY.md) <!-- TRACKER:BADGE -->

> **110 template files — 10 visual styles, 11 email types each**

---

## Philosophy

Most self-hosted Gitea instances use the default plain email templates. This project provides **ready-to-deploy, visually polished alternatives** — each designed for a specific community or audience, so you can pick the one that feels right for your users.

Every template is a drop-in replacement. All Go template variables, translation keys, and Gitea data contexts are fully compatible. **No patches, no plugins, no forks required.**

---

## Style Gallery

| Preview | Style | Audience | Character |
|---|---|---|---|
| ![Horizon](docs/images/horizon.png) | **Horizon** | Enterprise / Corporate | Blue accent, slate typography, centered cards |
| ![Terminal](docs/images/terminal.png) | **Terminal** | Developers / Tech | Dark mode, monospace, green CLI accents |
| ![Ember](docs/images/ember.png) | **Ember** | Community / Open Source | Warm amber, rounded, humanist, inclusive |
| ![Bloom](docs/images/bloom.png) | **Bloom** | Creative / Startup | Frosted glass, soft blue light, iridescent accents |
| ![Heritage](docs/images/heritage.png) | **Heritage** | Education / Research | Navy and gold, serif, classic, authoritative |
| ![Neon](docs/images/neon.png) | **Neon** | Gaming / Web3 / Creative Tech | Cyberpunk neon glow, hot pink & cyan, synthwave energy |
| ![Mono](docs/images/mono.png) | **Mono** | Design Studios / Editorial | Swiss brutalist, black & white, red accent, zero radius |
| ![Terra](docs/images/terra.png) | **Terra** | Sustainability / Wellness | Warm earth tones, organic textures, humanist serif |
| ![Ink](docs/images/ink.png) | **Ink** | Publishing / News / Literature | Editorial print, navy & gold, newspaper layout, drop caps |
| ![Aurora](docs/images/aurora.png) | **Aurora** | Premium SaaS / Mindfulness | Ethereal light gradients, deep purple & teal, atmospheric glow |

> Images are screenshots from the [live preview](preview/index.html). See [docs/images/README.md](docs/images/README.md) for capture instructions.

[**Live preview gallery**](preview/index.html) — open in a browser for an interactive style switcher with desktop/mobile viewports and view mode (Modern, Source).

---

## Installation

### Quick Start

Choose a style, then copy the `mail/` directory into your Gitea custom templates path:

```bash
# Locate your Gitea custom directory
#   (set by GITEA_CUSTOM; defaults shown below)

# Copy templates (example: Horizon style)
cp -r themes/horizon/mail/* /var/lib/gitea/custom/templates/mail/

# Restart Gitea
systemctl restart gitea
```

### Custom Directory Location

| Platform | Default Path |
|---|---|
| Linux (binary) | `/var/lib/gitea/custom` |
| Linux (Docker) | `/data/gitea` |
| Windows | `C:\gitea\custom` |

### Switching Styles

Overwrite the files with a different style. All templates share the exact same variable structure — no configuration changes needed.

```bash
cp -r themes/terminal/mail/* /var/lib/gitea/custom/templates/mail/
systemctl restart gitea
```

### Confirming It Works

The admin test email does not use custom mail templates. To verify your templates
are active, trigger a real email notification. The quickest way is the password
reset flow: log out, click **"Forgot password"** on the login page, and check the
reset email — it will render with your custom styles.

---

## Preview

Two modes are available — a zero-dependency static preview for quick checks, and a live-reload dev server for design work.

### Static Preview

Generate the preview data once (Go only), then open in a browser:

```bash
cd tools && go run . preview all
open preview/index.html        # no server needed
```

### Dev Server (Live Reload)

Start a pure Go development server that watches for `.tmpl` changes, auto-rebuilds, and pushes live updates to the browser via SSE:

```bash
cd tools && go run . dev
open http://localhost:3456
```

| Capability | Static | Dev |
|-----------|--------|-----|
| Go template rendering | ✅ | ✅ |
| Theme/template switching | ✅ | ✅ |
| Live reload on save | — | ✅ |

### Features

- Theme switcher — browse all 10 visual styles
- Template switcher — all 11 email types
- View mode — Modern (rendered preview), Source (raw HTML)
- Viewport toggle — Desktop 1386×780 / Mobile 390×780
- Parameter panel — mock data per email type
- Keyboard shortcuts — `←→` tab between Theme/Template/View, `↑↓` select within, `d`/`m` viewport

---

## Directory Structure

```
gitea-mail-templates/
├── themes/               # 10 visual styles, 11 .tmpl files each
│   ├── ...               # Custom styles are added here as separate directories
├── preview/              # Live preview SPA
│   ├── index.html        # Style/template/client/viewport switcher
│   └── rendered.js       # Pre-rendered templates (generated by tools/)
├── tools/                # Modular CLI tooling
│   ├── tools.go          #   Main entry point
│   ├── cli/              #   CLI subcommands (list, create, delete, preview)
│   ├── config/           #   Config types and templates_config.json loading
│   ├── data/             #   templates_config.json — single source of truth
│   ├── preview/          #   Template rendering engine
│   └── go.mod
├── docs/                 # Multi-language documentation
├── AGENTS.md             # AI agent guidance
├── CONTRIBUTING.md
├── LICENSE
├── README.md
└── .gitignore
```

### Template Types

| File | Email Trigger |
|---|---|
| `mail/user/auth/activate.tmpl` | Account activation |
| `mail/user/auth/activate_email.tmpl` | Email address verification |
| `mail/user/auth/register_notify.tmpl` | New registration notification |
| `mail/user/auth/reset_passwd.tmpl` | Password reset |
| `mail/org/team_invite.tmpl` | Team invitation |
| `mail/repo/collaborator.tmpl` | Repository collaborator added |
| `mail/repo/transfer.tmpl` | Repository ownership transfer |
| `mail/repo/release.tmpl` | New release published |
| `mail/repo/actions/workflow_run.tmpl` | Actions workflow run |
| `mail/repo/issue/assigned.tmpl` | Issue / Pull Request assigned |
| `mail/repo/issue/default.tmpl` | Issue / Pull Request updates |

---

## Compatibility

- **Gitea 1.25+** — matches the refactored mail template directory structure (v1.25)
- **Latest tested:** Gitea 1.27.0 <!-- TRACKER:LATEST-TESTED -->
- 100% variable-compatible with official Gitea templates — see [COMPATIBILITY.md](COMPATIBILITY.md) for the full matrix
- Uses only built-in Gitea template functions and official translation keys
- No custom template functions or locale patches required

---

## Design Principles

1. **Responsive** — Max-width 600px cards; works in all email clients
2. **Accessible** — 4.5:1 contrast ratios; semantic HTML
3. **Graceful degradation** — Fallback link visible when buttons fail to render
4. **Logo support** — References `{{AppUrl}}assets/img/favicon.png` by default
5. **i18n-ready** — All user-facing strings use Gitea's `{{.locale.Tr}}` system

---

## Documentation

- [English](README.md)
- [Simplified Chinese](docs/README.zh-CN.md)
- [Traditional Chinese](docs/README.zh-TW.md)
- [Russian](docs/README.ru.md)
- [Japanese](docs/README.ja.md)
- [Korean](docs/README.ko.md)

---

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines. Multi-language versions available in [docs/](docs/).

## License

MIT — see [LICENSE](LICENSE). Free to use, modify, and distribute in any Gitea deployment.

---

<p align="center">
  <sub>Not affiliated with the Gitea project. Gitea is a community-managed lightweight code hosting solution written in Go.</sub>
</p>
