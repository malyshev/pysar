import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const siteRoot = path.join(path.dirname(fileURLToPath(import.meta.url)), "..");
const src = path.join(siteRoot, "..", "install.sh");
const dest = path.join(siteRoot, "public", "install.sh");

if (!fs.existsSync(src)) {
  console.error(`stage-install: missing ${src}`);
  process.exit(1);
}

fs.mkdirSync(path.dirname(dest), { recursive: true });
fs.copyFileSync(src, dest);
console.log(`stage-install: ${path.relative(siteRoot, dest)}`);
