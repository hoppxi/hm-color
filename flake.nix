{
  description = "hm-color dynamic theming Go service";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
      ...
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        # build the Go binary
        packages.hm-color = pkgs.buildGoModule {
          pname = "hm-color";
          version = "0.1.0";

          src = ./.;
          vendorHash = "sha256-RzVtyevt/bFkuGkxQmgsDFHRV8eQcmLhZAzPyON3P4I=";

          subPackages = [ "." ];

          meta = with pkgs.lib; {
            description = "Dynamic theming tool for NixOS with swww wallpaper manager";
            longDescription = ''
              hm-color is a tool that integrates with Home Manager and swww to
              dynamically update your system's color theme. It can generate Nix,
              CSS, SCSS, or JSON outputs, commit changes to your Nix config, and
              optionally trigger a Home Manager switch.
            '';
            homepage = "https://github.com/hoppxi/hm-color";
            changelog = "https://github.com/hoppxi/hm-color/releases";
            license = licenses.mit;
            maintainers = with maintainers; [ hoppxi ];
            platforms = platforms.linux;
            mainProgram = "hm-color";
          };
        };
        packages.default = self.packages.${system}.hm-color;
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
          ];
        };
      }
    )
    // {
      # Home Manager module
      homeModules.hm-color =
        {
          config,
          lib,
          pkgs,
          inputs,
          ...
        }:
        let
          cfg = config.services.hm-color;
          bin = lib.getExe inputs.hm-color.packages.${pkgs.system}.hm-color;
          hmColorCmd = lib.concatStringsSep " " (
            [
              "${bin}"
              "--swww-cache ${cfg.swww-cache}"
              "--nix-out ${cfg.nix-theme-file}"
            ]
            ++ lib.optional cfg.theme "-t ${cfg.theme}"
            ++ lib.optional cfg.activate "-a"
          );
        in
        {
          options.services.hm-color = {
            enable = lib.mkEnableOption "hm-color dynamic theming";

            swww-cache = lib.mkOption {
              type = lib.types.path;
              default = "${config.xdg.cacheHome}/swww";
              description = "Path to swww cache directory.";
            };

            nix-theme-file = lib.mkOption {
              type = lib.types.path;
              default = null;
              description = "File where hm-color writes the generated nix theme.";
            };

            activate-hm = lib.mkOption {
              type = lib.types.bool;
              default = false;
              description = "Activate home-manager to apply the colors";
            };

            theme = lib.mkOption {
              type = lib.types.string;
              default = "dark";
              description = "Theme to use dark or light or system";
            };

            run-as-systemd = lib.mkOption {
              type = lib.types.bool;
              default = true;
              description = "Run hm-color as a systemd user service.";
            };

            run-in-hyprland = lib.mkOption {
              type = lib.types.bool;
              default = false;
              description = "Run hm-color via Hyprland's exec-once.";
            };
          };

          config = lib.mkIf cfg.enable {
            home.packages = [ inputs.hm-color.packages.${pkgs.system}.hm-color ];
            # Hyprland exec-once
            wayland.windowManager.hyprland.settings.exec-once = lib.mkIf cfg.run-in-hyprland [ hmColorCmd ];

            # systemd user service
            systemd.user.services."hm-color" = lib.mkIf cfg.run-as-systemd {
              Unit = {
                Description = "hm-color dynamic theming";
                After = [ "graphical-session.target" ];
              };
              Service = {
                ExecStart = hmColorCmd;
                Restart = "on-failure";
              };
              Install = {
                WantedBy = [ "default.target" ];
              };
            };
          };
        };
    };
}
