# hm-color

Dynamic theming for **NixOS** wallpapers using [swww](https://github.com/LGFae/swww) and **Material You** color extraction.
Whenever your wallpaper changes, `hm-color` extracts a full Material Design 3 color scheme and updates your system theme automatically.

## Installation

Using flakes:

```bash
nix run github:hoppxi/hm-color
```

Or add to your `home-manager` config:

```nix
# example flake.nix
{
	inputs = {
		nixpkgs.url = "github:NixOS/nixpkgs";
		hm-color.url = "github:hoppxi/hm-color";
	};
	outputs = { self, nixpkgs, hm-color, ... }:
		let
			system = "x86_64-linux";
			pkgs = nixpkgs.legacyPackages.${system};
		in {
			homeConfigurations.yourUser = nixpkgs.lib.homeManagerConfiguration {
			inherit pkgs;
			modules = [
				./home.nix
			];
		};
	};
}
```

```nix
# example home.nix
{
  config,
  inputs,
  ...
}:

{
  imports = [ inputs.hm-color.homeModules.hm-color ];
  services.hm-color = {
    enable = true;
    run-in-hyprland = true;
		# All needs either run-in-hyprland or run-as-systemd to be true.
    swww-cache = "${config.xdg.cacheHome}/swww";
    nix-theme-file = "${config.home.homeDirectory}/nix-config/home/theme/default.nix";
    flake-path = "${config.home.homeDirectory}/nix-config#hoppxi@ea"; # flake path must have #fragment
    git-commit = false;
  };
}
```

## Usage

```bash
hm-color \
  --swww-cache ~/.cache/swww \
  --nix-out ~/.config/hm-theme.nix \
  -f /home/you/.config/home-manager \
  -g
```

- `--swww-cache` → path to swww cache (default: `$XDG_CACHE_HOME/swww`)
- `--nix-out` → write theme as nix file
- `-f` → optional flake path for `home-manager switch`
- `-g` → optionally commit config changes
