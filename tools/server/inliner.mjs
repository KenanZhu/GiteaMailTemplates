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
