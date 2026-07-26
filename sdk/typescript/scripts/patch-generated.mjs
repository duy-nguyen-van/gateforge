#!/usr/bin/env node
/**
 * Post-process OpenAPI Generator typescript-fetch output so it typechecks
 * under NodeNext (explicit .js extensions) and drops unused mapValues imports.
 */
import fs from 'node:fs'
import path from 'node:path'
import { fileURLToPath } from 'node:url'

const root = path.join(
  path.dirname(fileURLToPath(import.meta.url)),
  '..',
  'src',
  'generated',
)

const relImport = /(?<=(?:from|import)\s+)(['"])(\.[^'"]+)\1/g

function addJsExt(spec) {
  if (/\.(js|json|ts|mjs|cjs)$/.test(spec)) {
    return spec
  }
  return `${spec}.js`
}

function walk(dir, out = []) {
  for (const ent of fs.readdirSync(dir, { withFileTypes: true })) {
    const full = path.join(dir, ent.name)
    if (ent.isDirectory()) {
      walk(full, out)
    } else if (ent.name.endsWith('.ts')) {
      out.push(full)
    }
  }
  return out
}

if (!fs.existsSync(root)) {
  console.error('patch-generated: missing', root)
  process.exit(1)
}

for (const file of walk(root)) {
  let text = fs.readFileSync(file, 'utf8')
  const before = text
  if (file.includes(`${path.sep}models${path.sep}`) && !file.endsWith(`${path.sep}index.ts`)) {
    text = text.replace(/^import \{ mapValues \} from '\.\.\/runtime(?:\.js)?';\n/m, '')
  }
  text = text.replace(relImport, (_, quote, spec) => `${quote}${addJsExt(spec)}${quote}`)
  if (text !== before) {
    fs.writeFileSync(file, text)
  }
}

fs.writeFileSync(
  path.join(root, 'index.ts'),
  "export * from './src/index.js'\n",
)

console.log('patch-generated: ok')
