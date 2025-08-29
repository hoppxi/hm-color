# recolor

> [!IMPORTANT]  
> **Project renamed:** This project was previously called **hm-color**, built around NixOS and Home Manager.  
> However, the Home Manager workflow proved too slow and rigid — every wallpaper change required a full rebuild, which is not ideal for large configurations.
>
> To solve this, the project has been rebranded to **recolor**:
>
> - It no longer depends on Home Manager.
> - Instead, it listens to wallpaper changes (currently via `swww`) and generates theme files directly.
> - This makes it lightweight, cross-distro, and easy to integrate with any Linux app.

Dynamic theming for Linux wallpapers using [swww](https://github.com/LGFae/swww) and **Material You** color extraction.  
Whenever your wallpaper changes, `recolor` extracts a full Material Design 3 color scheme and regenerates theme files for your system.

`recolor` can output multiple formats (CSS, SCSS, JSON, Nix).

---

## Installation

### With Go

```bash
go install github.com/hoppxi/recolor
```

### With Nix (optional)

You can still use `recolor` as a flake if you’re on NixOS:

```bash
nix run github:hoppxi/recolor
```

---

## Usage

### Export once

```bash
recolor export \
  --swww-cache ~/.cache/swww \
  --css-out ~/.config/theme.css
```

### Watch for changes

```bash
recolor watch \
  --swww-cache ~/.cache/swww \
  --json-out ~/.config/theme.json
```

### Run a command after update

```bash
recolor watch \
  --swww-cache ~/.cache/swww \
  --scss-out ~/.config/theme.scss \
```

---

## Migration from hm-color

If you were using `hm-color`:

- The binary is now called **`recolor`**.
- The `--nix-out` flag still works, so you can continue to generate Nix theme files if needed.
- Home-Manager integration is no longer bundled; instead, you can run `recolor` in `watch` mode and export to any format you like.
