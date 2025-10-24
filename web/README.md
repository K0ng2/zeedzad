# Nuxt frontend (web)

This is the Nuxt 4 frontend for Zeedzad. It is a standard Nuxt app and can be run with bun, npm, or pnpm.

See the Nuxt docs for more details: https://nuxt.com/docs/getting-started/introduction

## Install dependencies

```bash
# bun
bun install

# npm
npm install

# pnpm
pnpm install
```

## Development

Start the dev server (default: http://localhost:3000):

```bash
# bun
bun dev

# npm
npm run dev

# pnpm
pnpm dev
```

The frontend development server expects the backend API at `http://localhost:8088` by default. Change `NUXT_PUBLIC_API_BASE` if you run the backend on a different host/port.

## Production build

```bash
# build
bun run build

# preview locally
bun run preview
```

After build, copy the generated static assets into the backend embed directory if you plan to serve them from the Go binary:

```bash
cp -r .output/public/* ../pkg/web/public/
```

## Notes

- The frontend uses Tailwind CSS and DaisyUI.
- Use `NUXT_PUBLIC_API_BASE` to point to the backend API when developing.
