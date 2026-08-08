import { execFileSync } from "node:child_process";
import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";
import pngToIco from "png-to-ico";

const root = path.join(path.dirname(fileURLToPath(import.meta.url)), "..");
const publicDir = path.join(root, "public");
const appDir = path.join(root, "src/app");

execFileSync("python3", [path.join(root, "scripts/generate-mark.py")], {
  cwd: root,
  stdio: "inherit",
});

const toIco = pngToIco.default ?? pngToIco;
const buf = await toIco([
  path.join(publicDir, "favicon-16.png"),
  path.join(publicDir, "favicon-32.png"),
]);
fs.writeFileSync(path.join(appDir, "favicon.ico"), buf);
fs.writeFileSync(path.join(publicDir, "favicon.ico"), buf);
console.log("favicon.ico written");
