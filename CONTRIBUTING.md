# Contributing to Gitea Mail Templates

Thanks for your interest in contributing! This project aims to provide a diverse, well-maintained collection of email templates for the Gitea ecosystem.

---

## Ways to Contribute

### Adding a New Style

1. Scaffold the new style: `cd tools && go run . create <your-style-name>` — this creates the directory structure with placeholder `.tmpl` files for all 11 email types
2. Edit each `.tmpl` file in `themes/<your-style-name>/` with your unique visual design
3. Regenerate the preview: `cd tools && go run . preview all` — the build script auto-discovers all theme directories under `themes/` and generates the theme selector dynamically
4. Submit a PR with screenshots of rendered emails

### Style Guidelines

- Each style must include **all 11 template types** listed in the README
- Use only Gitea's built-in template functions — check the [Gitea source](https://github.com/go-gitea/gitea) for reference
- Translation keys must come from Gitea's official locale files (`mail.*` namespace)
- **Never reference `.DisplayName`** in templates where the data context lacks it (collaborator, transfer, release, workflow_run, assigned, default)
- Design for 600px max-width email clients
- Test against major email clients (Gmail, Outlook, Apple Mail) when possible

### Bug Reports

If a template doesn't render correctly:

1. Check that all referenced Go template variables exist — compare against the Gitea source mail templates
2. Verify translation keys match Gitea's locale files
3. Confirm `.DisplayName` isn't used in templates that lack it
4. Regenerate the preview: `cd tools && go run . preview all`
5. Open an issue with: the style name, which email type, and the error or unexpected output

### Documentation Improvements

Documentation updates, preview screenshots, installation guides, and translations are always welcome.

---

## Development Setup

No build tools or dependencies are needed — these are raw Go HTML templates.

### Previewing Locally

1. Open `preview/index.html` directly in a browser — no server needed
2. Use the theme switcher, template selector, and client mode toggles to review designs
3. Toggle between Modern, Gmail, Outlook, and Raw source modes to verify degradation

### Regenerating Previews

```bash
cd tools && go run . preview all
```

This renders all templates (themes auto-discovered from the themes/ directory) using Go's native `html/template` package and writes the output to `preview/rendered.js`. The `--folder` and `--config` flags default to `../themes` and `./data/templates_config.json` respectively — override them only when using a custom layout.

### Integration Testing

Deploy the templates to a Gitea instance and use the admin test email feature:
**Site Administration > Configuration > Mailer > Send Test Email**

---

## Commit Conventions

- `style(horizon|terminal|ember|bloom|heritage|neon|mono|terra|ink|aurora):` — template changes
- `preview:` — preview tooling changes
- `tools:` — Go build script changes
- `docs:` — documentation and translations
- `fix:` — bug fixes
- `project:` — README, LICENSE, AGENTS.md, meta

---

## Translations

- [English](CONTRIBUTING.md)
- [Simplified Chinese](docs/CONTRIBUTING.zh-CN.md)
- [Traditional Chinese](docs/CONTRIBUTING.zh-TW.md)
- [Russian](docs/CONTRIBUTING.ru.md)
- [Japanese](docs/CONTRIBUTING.ja.md)
- [Korean](docs/CONTRIBUTING.ko.md)

---

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
