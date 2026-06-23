# Spec: Remove Juice & Simplify Preview Module

**Date:** 2026-06-23
**Status:** approved

## Goal

Remove multi-email-client simulation (Gmail/Outlook) from the preview module, eliminate the `juice` CSS-inlining dependency, and rewrite the dev server in pure Go with SSE-based live reload.

## Motivation

- Multi-client CSS simulation has diminishing value — modern email clients render consistently
- `juice` + Node.js adds complexity (npm install, node_modules, separate runtime)
- The project already uses Go for all build tooling; a Node.js dev server is an outlier
- Simplifying to "Modern" and "Source" views covers the real use cases

## Changes

### 1. Delete `tools/server/` (entire directory)

Remove all Node.js artifacts:
- `inliner.mjs` — juice CSS inlining + Gmail/Outlook CSS stripping
- `server.mjs` — Express + WebSocket dev server
- `package.json` / `package-lock.json`
- `node_modules/` — all npm dependencies

### 2. Rewrite `tools/cli/dev.go`

- Remove Node.js dependency check
- Remove `exec.Command("node", "server.mjs")` subprocess
- Instead: create and start a pure Go HTTP server (calling into `tools/preview/server.go`)
- Keep the same CLI interface: `go run . dev [--port <port>]`

### 3. New file: `tools/preview/server.go`

Pure Go dev server with:
- **Static file serving** — serve `preview/` directory (index.html, rendered.js)
- **SSE endpoint** (`GET /events`) — Server-Sent Events for browser reload notifications
- **File watcher** — poll `themes/` every 500ms, compare mod times of `.tmpl` files
- **On change detected** — call `preview.RenderAll()` directly (in-process, no subprocess), write `rendered.js`, broadcast SSE `reload` event
- No external Go dependencies — uses `net/http` stdlib only

### 4. Modify `preview/index.html`

Frontend changes:
- **Client selector** — reduce from 4 options (Modern/Gmail/Outlook/Raw Source) to 2 (Modern/Source)
- **Remove** `renderedGmail` / `renderedOutlook` variables and all Gmail/Outlook JS logic
- **Simplify** `getRendered()` — always returns `rendered`
- **Simplify** `transform()` — only handles `source` mode (HTML escaping)
- **Remove** Client Simulation indicators panel (HTML section + JS in `updatePanel()`)
- **Remove** static-mode warning message ("Gmail / Outlook simulation is approximate...")
- **Replace** WebSocket with SSE (`new EventSource('/events')`)
- **SSE reload handler** — dynamically reload `rendered.js` on `reload` event, re-render iframe without losing current theme/template/viewport selection
- **Simplify** `setDevMode()` — remove multi-client disclaimer

### 5. Update `AGENTS.md`

- Remove references to Juice, Node.js, Gmail/Outlook variants
- Update dev server description to reflect pure Go implementation

## Non-Changes

- `tools/preview/engine.go` — template rendering engine unchanged
- `tools/preview/funcs.go` — template functions unchanged
- `tools/preview/locale.go` — locale data unchanged
- `tools/config/` — config loading unchanged
- `tools/data/templates_config.json` — unchanged
- `tools/cli/preview.go` — preview command unchanged
- `tools/cli/commands.go` — command registration unchanged
- All theme `.tmpl` files — unchanged

## Behavior

### Static preview (open `preview/index.html` directly)
- Works via `file://` protocol as before
- Two view modes: Modern (rendered HTML in iframe) and Source (escaped HTML source)
- No dev-mode warnings needed

### Dev mode (`go run . dev`)
- `go run . dev` starts the Go HTTP server on port 3456 (configurable)
- On startup: runs preview engine once, writes `rendered.js`
- Serves `preview/` as static files
- Watches `themes/` for `.tmpl` changes (500ms polling)
- On change: re-renders affected themes, updates `rendered.js`, pushes SSE event
- Browser auto-reloads preview content without page refresh (preserves UI state)

## Constraints

- No new Go dependencies — SSE and file polling use stdlib only
- No Node.js requirement
- `preview/rendered.js` format stays compatible: `window.__RENDERED__`, `window.__REGISTRY__`, `window.__PARAMS__`
- Preview still works with `file://` protocol (no server needed for basic use)
