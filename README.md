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
  imports = [ inputs.hm-color.homeModules.default ];

  hm-color = {
    execOnceHyprland = true;
    swww-cache = "${config.xdg.cacheHome}/swww";
    nix-theme-file = "${config.xdg.configHome}/hm-theme.nix";
    flake-path = "/home/you/.config/home-manager"; # optional
    gitCommit.enable = true; # optional
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

```nix
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
				{
					home.packages = [ hm-color.packages.${system}.default ];
					# add to exec-once
  				wayland.windowManager.hyprland.settings.exec-once = [
						"hm-color -n ~/nix-config/theme.nix -f ~/nix-config -gc -gp"
					];
				}
			];
		};
	};
}
```
