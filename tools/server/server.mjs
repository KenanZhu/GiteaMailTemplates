import { readFileSync, writeFileSync, existsSync, watch } from 'fs';
import { resolve, dirname, basename, join, relative } from 'path';
import { fileURLToPath } from 'url';
import { spawn } from 'child_process';

import express from 'express';
import expressWs from 'express-ws';

import { inlineCSS, stripGmail, stripOutlook } from './inliner.mjs';

// ANSI color helpers
const G = '\x1b[32m';   // green (Info)
const Y = '\x1b[33m';   // yellow (Warn)
const R = '\x1b[31m';   // red (Error)
const Z = '\x1b[0m';    // reset
const I = `${G}[I]${Z}`;
const W = `${Y}[W]${Z}`;
const E = `${R}[E]${Z}`;

const __dirname = dirname(fileURLToPath(import.meta.url));
const ROOT = resolve(__dirname, '..', '..');       // project root
const TOOLS = join(ROOT, 'tools');                   // tools/ (for go run)
const PREVIEW = join(ROOT, 'preview');               // preview/
const THEMES = join(ROOT, 'themes');                 // themes/
const RENDERED_JS = join(PREVIEW, 'rendered.js');
const PORT = parseInt(process.env.PORT || '3456', 10);

// Express + WebSocket
const app = express();
expressWs(app);

// Track connected WebSocket clients
const clients = new Set();

app.ws('/ws', (ws) => {
  clients.add(ws);
  ws.on('close', () => clients.delete(ws));
});

// Broadcast to all connected browsers
function broadcast(type, payload = {}) {
  const msg = JSON.stringify({ type, ...payload });
  for (const c of clients) c.send(msg);
}

// Serves preview/ as static (index.html, rendered.js)
app.use(express.static(PREVIEW, { etag: false }));

// Rebuild pipeline
let building = false;
const CHANGED_THEMES = new Set();

function getThemeFromPath(filePath) {
  const rel = relative(THEMES, filePath).replace(/\\/g, '/');
  const parts = rel.split('/').filter(Boolean);
  return parts[0] || null;
}

async function rebuild() {
  if (building) { building = false; return rebuild(); } // debounce — restart
  building = true;

  const themes = [...CHANGED_THEMES];
  CHANGED_THEMES.clear();

  const start = Date.now();
  console.log(`${I} [Builder] Rebuilding: ${themes.length > 0 ? themes.join(', ') : 'all'} changed`);

  // Step 1: Go render
  try {
    await goPreview(themes.length > 0 ? themes : ['all']);
  } catch (err) {
    console.error(`${E} [Builder] Go render failed: ${err.message}`);
    building = false;
    broadcast('error', { message: err.message });
    return;
  }

  // Step 2: Read rendered.js, apply juice inlining
  try {
    const renderedJS = readFileSync(RENDERED_JS, 'utf8');
    const juiced = await juiceRenderedOutput(renderedJS);
    writeFileSync(RENDERED_JS, juiced, 'utf8');
  } catch (err) {
    console.error(`${E} [Builder] Juice step failed: ${err.message}`);
  }

  const elapsed = Date.now() - start;
  console.log(`${I} [Builder] Rebuild done in ${elapsed}ms`);

  building = false;

  // Step 3: Notify browsers
  broadcast('reload', { elapsed });
}

// Go preview (child process)
function goPreview(themes) {
  return new Promise((resolve, reject) => {
    const args = ['run', '.', 'preview', '--folder', '../themes', '--config', './data/templates_config.json', ...themes];
    const proc = spawn('go', args, { cwd: TOOLS, stdio: ['ignore', 'pipe', 'pipe'] });

    let stdout = '', stderr = '';
    proc.stdout.on('data', (d) => { stdout += d.toString(); });
    proc.stderr.on('data', (d) => { stderr += d.toString(); });

    proc.on('close', (code) => {
      if (code !== 0) reject(new Error(stderr.trim() || `exit ${code}`));
      else {
        // Print the summary lines
        const lines = stdout.split('\n').filter(l => l.includes('all required') || l.includes('only'));
        for (const l of lines) console.log('  ' + l.trim());
        resolve(stdout);
      }
    });

    proc.on('error', reject);
  });
}

// Juice post-processing
function juiceRenderedOutput(jsContent) {
  // rendered.js format: window.__RENDERED__ = { "theme": { "tpl": "<html>..." } };
  // We inline CSS and produce three client-specific variants:
  //   __RENDERED__         — modern  (Juice-inlined only)
  //   __RENDERED_GMAIL__   — Gmail   (Juice + strip unsupported CSS)
  //   __RENDERED_OUTLOOK__ — Outlook (Juice + aggressive CSS strip)
  const match = jsContent.match(/window\.__RENDERED__\s*=\s*(\{[\s\S]*?\});/);
  if (!match) return Promise.resolve(jsContent);

  let rendered;
  try {
    rendered = JSON.parse(match[1]);
  } catch {
    return Promise.resolve(jsContent);
  }

  // Deep-clone for Gmail and Outlook variants
  const renderedGmail = JSON.parse(JSON.stringify(rendered));
  const renderedOutlook = JSON.parse(JSON.stringify(rendered));

  let changed = false;
  for (const theme of Object.keys(rendered)) {
    for (const tpl of Object.keys(rendered[theme])) {
      const html = rendered[theme][tpl];
      if (typeof html !== 'string') continue;
      if (!html.includes('<style') && !html.includes('<html')) continue;
      try {
        // Modern — Juice inlining only
        const modern = inlineCSS(html);
        rendered[theme][tpl] = modern;
        changed = true;

        // Gmail — Juice + strip unsupported CSS
        try { renderedGmail[theme][tpl] = stripGmail(modern); } catch {}
        // Outlook — Juice + aggressive CSS strip
        try { renderedOutlook[theme][tpl] = stripOutlook(modern); } catch {}
      } catch {
        // keep original on failure
      }
    }
  }

  if (changed) {
    const rest = jsContent.replace(/window\.__RENDERED__\s*=\s*\{[\s\S]*?\};/, '');
    return Promise.resolve(
      'window.__RENDERED__ = ' + JSON.stringify(rendered, null, 2) + ';\n' +
      'window.__RENDERED_GMAIL__ = ' + JSON.stringify(renderedGmail, null, 2) + ';\n' +
      'window.__RENDERED_OUTLOOK__ = ' + JSON.stringify(renderedOutlook, null, 2) + ';\n' +
      rest
    );
  }
  return Promise.resolve(jsContent);
}

// File watcher (native fs.watch, recursive)
function relPath(absPath) {
  return relative(ROOT, absPath).replace(/\\/g, '/');
}

// Debounced rebuild trigger
let rebuildTimer = null;
function scheduleRebuild(filePath) {
  const theme = getThemeFromPath(filePath);
  if (!theme) return;

  const ext = filePath.slice(filePath.lastIndexOf('.'));
  if (ext !== '.tmpl') return;

  CHANGED_THEMES.add(theme);
  clearTimeout(rebuildTimer);
  rebuildTimer = setTimeout(rebuild, 300);
}

// Watch themes directory recursively (Windows supports recursive)
try {
  watch(THEMES, { recursive: true }, (event, fname) => {
    if (!fname || !fname.endsWith('.tmpl')) return;
    const fullPath = join(THEMES, fname);
    // event is 'rename' for both create and delete; 'change' for modifications
    if (event === 'change') {
      console.log(`${I} [Watcher] ${relPath(fullPath)} changed`);
      scheduleRebuild(fullPath);
    } else if (event === 'rename') {
      if (existsSync(fullPath)) {
        console.log(`${I} [Watcher] ${relPath(fullPath)} created`);
      } else {
        console.log(`${I} [Watcher] ${relPath(fullPath)} deleted`);
      }
    }
  });
  console.log(`${I} [Watcher] Watching themes/ (recursive)`);
} catch (err) {
  console.error(`${E} [Watcher] File watch failed: ${err.message}`);
}

// Initial Juice pass
// The Go `preview all` step (dev.go) already generated raw rendered.js.
// Only run juice inlining — no need to re-render with Go.
(async function initialJuice() {
  try {
    const js = readFileSync(RENDERED_JS, 'utf8');
    const juiced = await juiceRenderedOutput(js);
    writeFileSync(RENDERED_JS, juiced, 'utf8');
    console.log(`${I} [Builder] CSS inlining applied`);
  } catch (err) {
    console.error(`${E} [Builder] CSS inlining failed: ${err.message}`);
  }
})();

// Start server
app.get('/health', (_req, res) => res.json({ status: 'ok', port: PORT }));

app.listen(PORT, () => {
  console.log(`${I} [Server] Gitea Mail Templates — Dev Server`);
  console.log(`${I} [Server] http://localhost:${PORT}`);
  console.log(`${I} [Watcher] Watching themes/ for changes`);
  console.log(`${I} [Server] WebSocket ws://localhost:${PORT}/ws\n`);
});
