#!/usr/bin/env python3
"""Generate Pysar vyshyvanka pixel mark + raster icons (transparent bg)."""

from __future__ import annotations

import subprocess
import sys
from pathlib import Path

try:
    from PIL import Image, ImageDraw
except ImportError:
    subprocess.check_call([sys.executable, "-m", "pip", "install", "pillow", "-q"])
    from PIL import Image, ImageDraw

RED = (252, 67, 43, 255)
BLK = (27, 27, 27, 255)
LIT = (241, 241, 241, 255)
INK = (23, 23, 23, 255)
EMP = (0, 0, 0, 0)

# 19×19 with 1-cell quiet margin so corners never feel cropped.
N = 19
g = [["." for _ in range(N)] for _ in range(N)]


def P(x: int, y: int, c: str) -> None:
    if 0 <= x < N and 0 <= y < N:
        g[y][x] = c


def build() -> None:
    # Tip crosses inset from edges (not flush)
    for cx, cy in ((2, 2), (16, 2), (2, 16), (16, 16)):
        for dx, dy in ((0, 0), (-1, 0), (1, 0), (0, -1), (0, 1)):
            P(cx + dx, cy + dy, "R")

    # Short L arms from each corner (leave quiet margin at 0/18)
    for x in range(1, 4):
        P(x, 1, "R")
        P(x, 17, "R")
        P(18 - x, 1, "R")
        P(18 - x, 17, "R")
    for y in range(1, 4):
        P(1, y, "R")
        P(17, y, "R")
        P(1, 18 - y, "R")
        P(17, 18 - y, "R")

    # Re-assert tip crosses after L paint
    for cx, cy in ((2, 2), (16, 2), (2, 16), (16, 16)):
        for dx, dy in ((0, 0), (-1, 0), (1, 0), (0, -1), (0, 1)):
            P(cx + dx, cy + dy, "R")

    def diamond(cx: int, cy: int, flank: str) -> None:
        P(cx, cy - 1, "K")
        P(cx, cy + 1, "K")
        P(cx - 1, cy, "K")
        P(cx + 1, cy, "K")
        if flank == "h":
            P(cx - 2, cy, "R")
            P(cx + 2, cy, "R")
        else:
            P(cx, cy - 2, "R")
            P(cx, cy + 2, "R")

    diamond(9, 3, "h")
    diamond(9, 15, "h")
    diamond(3, 9, "v")
    diamond(15, 9, "v")

    # N/S black chevrons, E/W red chevrons around center (9,9)
    nb = [
        (9, 4),
        (8, 5),
        (9, 5),
        (10, 5),
        (7, 6),
        (8, 6),
        (10, 6),
        (11, 6),
        (7, 7),
        (11, 7),
    ]
    for x, y in nb:
        P(x, y, "K")
        P(x, 18 - y, "K")

    er = [
        (14, 9),
        (13, 8),
        (13, 9),
        (13, 10),
        (12, 7),
        (12, 8),
        (12, 10),
        (12, 11),
        (11, 7),
        (11, 11),
    ]
    for x, y in er:
        P(x, y, "R")
        P(18 - x, y, "R")

    P(9, 9, "R")
    P(8, 9, "K")
    P(10, 9, "K")
    P(9, 8, "K")
    P(9, 10, "K")
    for x, y in ((7, 8), (7, 10), (11, 8), (11, 10), (8, 7), (10, 7), (8, 11), (10, 11)):
        if g[y][x] == ".":
            P(x, y, "R")


def render_mark(scale: int, black=BLK) -> Image.Image:
    im = Image.new("RGBA", (N * scale, N * scale), EMP)
    for y in range(N):
        for x in range(N):
            c = g[y][x]
            if c == ".":
                continue
            color = RED if c == "R" else black
            for dy in range(scale):
                for dx in range(scale):
                    im.putpixel((x * scale + dx, y * scale + dy), color)
    return im


def sized(target: int, black=BLK, bg=EMP) -> Image.Image:
    scale = max(1, target // N)
    while scale * N > target:
        scale -= 1
    scale = max(1, scale)
    mark = render_mark(scale, black=black)
    canvas = Image.new("RGBA", (target, target), bg)
    ox = (target - mark.size[0]) // 2
    oy = (target - mark.size[1]) // 2
    canvas.paste(mark, (ox, oy), mark)
    return canvas


def on_ink(target: int) -> Image.Image:
    im = Image.new("RGBA", (target, target), EMP)
    draw = ImageDraw.Draw(im)
    r = int(target * 0.1875)
    draw.rounded_rectangle((0, 0, target - 1, target - 1), radius=r, fill=INK)
    # Keep mark inside rounded corners — ~58% of canvas
    mark = sized(int(target * 0.58), black=LIT, bg=EMP)
    ox = (target - mark.size[0]) // 2
    oy = (target - mark.size[1]) // 2
    im.paste(mark, (ox, oy), mark)
    return im


# Each logical stitch is CELL×CELL SVG units — integer geometry, no fractional scale().
CELL = 16


def mark_rects(ink: str) -> str:
    rects: list[str] = []
    for y in range(N):
        for x in range(N):
            c = g[y][x]
            if c == "R":
                fill = "#FC432B"
            elif c == "K":
                fill = ink
            else:
                continue
            rects.append(
                f'<rect x="{x * CELL}" y="{y * CELL}" width="{CELL}" height="{CELL}" fill="{fill}"/>'
            )
    return "\n".join(rects)


def write_mark_svg(path: Path, dark: bool = False) -> None:
    """Mark only — wordmark text is composed in React (sans, full 'Pysar')."""
    ink = "#F1F1F1" if dark else "#1B1B1B"
    body = mark_rects(ink)
    mark_px = N * CELL  # 304
    path.write_text(
        f'''<svg xmlns="http://www.w3.org/2000/svg" width="{mark_px}" height="{mark_px}" viewBox="0 0 {mark_px} {mark_px}" fill="none" shape-rendering="crispEdges" role="img" aria-label="Pysar">
{body}
</svg>
'''
    )


def main() -> None:
    root = Path(__file__).resolve().parent.parent
    public = root / "public"
    app = root / "src" / "app"
    build()

    for row in g:
        print("".join(row))

    write_mark_svg(public / "logo-mark.svg", dark=False)
    write_mark_svg(public / "logo-mark-dark.svg", dark=True)
    # Aliases used by older paths / OG tooling
    write_mark_svg(public / "logo.svg", dark=False)
    write_mark_svg(public / "logo-dark.svg", dark=True)

    # Transparent mark PNG — native stitch grid + high-DPI nearest-neighbor sizes
    render_mark(1).save(public / "logo-mark.png")
    render_mark(CELL).save(public / "logo-mark@16.png")  # 304² crisp source
    sized(16).save(public / "favicon-16.png")
    sized(32).save(public / "favicon-32.png")
    sized(48).save(public / "favicon-48.png")
    sized(192).save(public / "icon-192.png")
    sized(512).save(public / "icon-512.png")
    # Prefer a denser app icon for retina tabs (Next serves /icon.png)
    sized(64).save(app / "icon.png")

    on_ink(180).save(public / "apple-touch-icon.png")
    on_ink(180).save(app / "apple-icon.png")
    on_ink(512).save(public / "logo.png")
    on_ink(1024).save(public / "logo@2x.png")

    og = Image.new("RGBA", (1200, 630), INK)
    mark = sized(N * 20, black=LIT, bg=EMP)  # 380
    og.paste(mark, ((1200 - mark.size[0]) // 2, (630 - mark.size[1]) // 2), mark)
    og.save(public / "og-default.png")

    # Drop debug preview if present
    preview = public / "logo-mark-preview.png"
    if preview.exists():
        preview.unlink()

    print("mark + rasters written")


if __name__ == "__main__":
    main()
