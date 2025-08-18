{
  description = "hm-color - Dynamic theming tool for NixOS with swww";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs";
    flake-utils.url = "github:numtide/flake-utils";
    home-manager.url = "github:nix-community/home-manager";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
      home-manager,
      ...
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages.default = pkgs.buildGoModule {
          pname = "hm-color";
          version = "0.1.0";

          src = ./.;

          vendorHash = null; # fill with nix build
          subPackages = [ "." ];

          meta = with pkgs.lib; {
            description = "Dynamic theming for NixOS wallpapers using swww";
            homepage = "https://github.com/hoppxi/hm-color";
            license = licenses.mit;
            maintainers = [ maintainers.yourself ];
            platforms = platforms.linux;
          };
        };

        devShells.default = pkgs.mkShell {
          buildInputs = [
            pkgs.go
            pkgs.git
          ];
        };
      }
    )
    // {
      homeModules.default =
        {
          config,
          pkgs,
          lib,
          ...
        }:
        let
          cfg = config.hm-color;
          bin = "${self.packages.${config.system}.default}/bin/hm-color";
          hmColorCmd = lib.concatStringsSep " " (
            [
              "${bin}"
              "--swww-cache ${cfg.swww-cache}"
              "-n"
              "--nix-out ${cfg.nix-theme-file}"
            ]
            ++ lib.optional (cfg.flake-path != "") "-f ${cfg.flake-path}"
            ++ lib.optional cfg.gitCommit.enable "-g"
          );
        in
        {
          options.hm-color = {
            execOnceHyprland = lib.mkOption {
              type = lib.types.bool;
              default = false;
              description = "Run hm-color via Hyprland's exec-once.";
            };

            execOnceSystemd = lib.mkOption {
              type = lib.types.bool;
              default = false;
              description = "Run hm-color as a systemd user service.";
            };

            swww-cache = lib.mkOption {
              type = lib.types.path;
              default = "${config.xdg.cacheHome}/swww";
              description = "Path to swww cache directory.";
            };

            nix-theme-file = lib.mkOption {
              type = lib.types.path;
              default = "";
              description = "File where hm-color writes the generated nix theme.";
            };

            flake-path = lib.mkOption {
              type = lib.types.path;
              default = "";
              description = "Flake path for home-manager switch.";
            };

            gitCommit.enable = lib.mkOption {
              type = lib.types.bool;
              default = false;
              description = "Enable committing Nix config changes after updating colors.";
            };
          };

          config = {
            home.packages = [ self.packages.${config.system}.default ];

            wayland.windowManager.hyprland.settings.exec-once = lib.mkIf cfg.execOnceHyprland [ hmColorCmd ];

            systemd.user.services."hm-color" = lib.mkIf cfg.execOnceSystemd {
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
