# Style Preview Images

Place theme style screenshots of **Desktop** in PNG format here, captured from the [live preview](../preview/index.html).

## Naming Convention

```
horizon.png    terminal.png    ember.png    bloom.png    heritage.png
neon.png       mono.png        terra.png    ink.png      aurora.png
```

## Image Size Requirements

To ensure screenshots can be displayed in the README, please follow these size requirements:

- **Maximum:** 50 KiB per image
- **Recommended:** 30–40 KiB
- **Format:** PNG, optimised — run through `pngquant` or `optipng` before committing

## How to Capture

1. Start the dev server: `cd tools && go run . dev` (opens http://localhost:3456 in your browser)
2. For each style, select the **"Activate Account"** template and **"Modern"** client mode in **"Desktop"**
3. Take a screenshot of the rendered email (600px width recommended)
4. Save as `<style-name>.png` in this directory

> For static preview, run `cd tools && go run . preview all` then open `preview/index.html` directly.
