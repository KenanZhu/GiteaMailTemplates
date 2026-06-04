import juice from 'juice';

/** Juice CSS inlining options optimized for HTML email. */
const JUICE_OPTIONS = {
  removeStyleTags: false,    // keep <style> for modern clients that support it
  preserveMediaQueries: true,
  preserveFontFaces: true,
  applyWidthAttributes: true,
  applyHeightAttributes: true,
  applyAttributesTableElements: true,
  xmlMode: false,
};

/**
 * Inline CSS from <style> blocks and <link> tags into inline style attributes.
 * Produces HTML that renders consistently across Gmail, Outlook, Apple Mail.
 */
export function inlineCSS(html) {
  return juice(html, JUICE_OPTIONS);
}

// Gmail / Outlook CSS stripping

const GMAIL_UNSUPPORTED = [
  'border-radius','box-shadow','text-shadow',
  'backdrop-filter','filter',
  'animation','transition','transform',
  'background-image\\s*:\\s*(linear|radial|conic)-gradient',
];
const OUTLOOK_UNSUPPORTED = GMAIL_UNSUPPORTED.concat([
  'display\\s*:\\s*(flex|grid|inline-flex|inline-grid)',
  'flex[^;]*','grid[^;]*','gap\\s*:','order\\s*:',
]);

function stripDeclarations(styleValue, patterns) {
  let v = styleValue;
  for (const p of patterns) {
    v = v.replace(new RegExp('(^|;\\s*)' + p + '[^;]*(;|$)', 'gi'), (m) => m.endsWith(';') ? ';' : '');
    v = v.replace(new RegExp(p + '[^;]*;', 'gi'), '');
  }
  return v.replace(/;;+/g, ';').replace(/^\s*;\s*/, '').replace(/\s*;\s*$/, '').trim();
}

/**
 * Strip CSS properties unsupported by Gmail from inline styles and <style> blocks.
 */
export function stripGmail(html) {
  return stripClientCSS(html, GMAIL_UNSUPPORTED, false);
}

/**
 * Strip CSS properties unsupported by Outlook (more aggressive than Gmail).
 */
export function stripOutlook(html) {
  return stripClientCSS(html, OUTLOOK_UNSUPPORTED, true);
}

function stripClientCSS(html, patterns, isOutlook) {
  let result = html;

  // 1. Strip from inline style="..." attributes
  result = result.replace(/style="([^"]*)"/gi, (_, styles) => {
    const cleaned = stripDeclarations(styles, patterns);
    return cleaned ? 'style="' + cleaned + '"' : '';
  });

  // 2. Strip from <style> blocks
  result = result.replace(/<style[^>]*>([\s\S]*?)<\/style>/gi, (_, cssContent) => {
    let cleaned = cssContent;
    for (const p of patterns) {
      const re = new RegExp(p + '[^;{}]*[;}]?', 'gi');
      cleaned = cleaned.replace(re, (m) => m.endsWith('}') ? '}' : '');
    }
    return '<style>' + cleaned + '</style>';
  });

  // 3. Outlook: replace CSS custom properties with initial
  if (isOutlook) {
    result = result.replace(/style="([^"]*)"/gi,
      (_, s) => 'style="' + s.replace(/var\(--[^)]*\)/g, 'initial') + '"');
  }

  return result;
}

/**
 * Apply juice inlining to a rendered.js-style HTML snippet and return the
 * result. Handles escaped HTML (JSON-stringified) by first unescaping,
 * inlining, then re-escaping for storage in rendered.js.
 */
export function inlineRenderedHTML(rawhtml) {
  try {
    // If this is already a standalone HTML document fragment
    if (rawhtml.includes('<style') || rawhtml.includes('</style>')) {
      return inlineCSS(rawhtml);
    }
    // Otherwise pass through unchanged (inline-only content)
    return rawhtml;
  } catch (err) {
    console.error('Juice inline error:', err.message);
    return rawhtml;
  }
}

if (process.argv[1] && import.meta.url.endsWith(process.argv[1])) {
  // CLI mode: juice <input> [output]
  const { readFileSync, writeFileSync } = await import('fs');
  const input = process.argv[2];
  const output = process.argv[3];
  if (!input) { console.error('Usage: node inliner.mjs <input.html> [output.html]'); process.exit(1); }
  const html = readFileSync(input, 'utf8');
  const result = inlineCSS(html);
  if (output) { writeFileSync(output, result, 'utf8'); console.log(`Inlined → ${output}`); }
  else { console.log(result); }
}
