{
  description = "recolor dynamic theming Go service";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    mcuc.url = "github:hoppxi/mcu-cli";
    mcuc.inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
      mcuc,
      ...
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        # build the Go binary
        packages.recolor = pkgs.buildGoModule {
          pname = "recolor";
          version = "1.0.0";

          src = ./.;
          vendorHash = "sha256-RzVtyevt/bFkuGkxQmgsDFHRV8eQcmLhZAzPyON3P4I=";

          subPackages = [ "." ];

          meta = with pkgs.lib; {
            description = "Dynamic theming tool for NixOS with swww wallpaper manager";
            longDescription = ''
              recolor is a tool that integrates with Home Manager and swww to
              dynamically update your system's color theme. It can generate Nix,
              CSS, SCSS, or JSON outputs, commit changes to your Nix config, and
              optionally trigger a Home Manager switch.
            '';
            homepage = "https://github.com/hoppxi/recolor";
            changelog = "https://github.com/hoppxi/recolor/releases";
            license = licenses.mit;
            maintainers = with maintainers; [ hoppxi ];
            platforms = platforms.linux;
            mainProgram = "recolor";
          };
        };
        packages.default = self.packages.${system}.recolor;
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            swww
            mcuc.packages.${system}.default
          ];
        };
      }
    )
    // {
      # Home Manager module
      homeModules.recolor =
        {
          config,
          lib,
          pkgs,
          inputs,
          ...
        }:
        let
          cfg = config.services.recolor;
          bin = lib.getExe self.packages.${pkgs.system}.recolor;
          hmColorCmd = lib.concatStringsSep " " (
            [
              "${bin}"
              "--swww-cache ${cfg.swww-cache}"
              "-t ${cfg.theme}"
            ]
            ++ lib.optional (cfg.nix-theme-file != null) "--nix-out ${cfg.nix-theme-file}"
            ++ lib.optional (cfg.scss-theme-file != null) "--scss-out ${cfg.scss-theme-file}"
            ++ lib.optional (cfg.css-theme-file != null) "--css-out ${cfg.css-theme-file}"
            ++ lib.optional (cfg.json-theme-file != null) "--json-out ${cfg.json-theme-file}"
          );
        in
        {
          options.services.recolor = {
            enable = lib.mkEnableOption "recolor dynamic theming";

            swww-cache = lib.mkOption {
              type = lib.types.path;
              default = "${config.xdg.cacheHome}/swww";
              description = "Path to swww cache directory.";
            };

            nix-theme-file = lib.mkOption {
              type = lib.types.path;
              default = null;
              description = "File where recolor writes the generated nix theme.";
            };

            scss-theme-file = lib.mkOption {
              type = lib.types.path;
              default = null;
              description = "File where recolor writes the generated scss theme.";
            };

            css-theme-file = lib.mkOption {
              type = lib.types.path;
              default = null;
              description = "File where recolor writes the generated css theme.";
            };

            json-theme-file = lib.mkOption {
              type = lib.types.path;
              default = null;
              description = "File where recolor writes the generated json theme.";
            };

            theme = lib.mkOption {
              type = lib.types.enum [
                "dark"
                "light"
                "system"
              ];
              default = "dark";
              description = "Theme to use dark or light or system";
            };

            start-with-systemd = lib.mkOption {
              type = lib.types.bool;
              default = true;
              description = "Run recolor as a systemd user service.";
            };

            start-with-hyprland = lib.mkOption {
              type = lib.types.bool;
              default = false;
              description = "Run recolor via Hyprland's exec-once.";
            };
          };

          config = lib.mkIf cfg.enable {
            home.packages = [
              self.packages.${pkgs.system}.recolor
              mcuc.packages.${pkgs.system}.default
              pkgs.swww
            ];
            # Hyprland exec-once
            wayland.windowManager.hyprland.settings.exec-once = lib.mkIf cfg.start-with-hyprland [ hmColorCmd ];

            # systemd user service
            systemd.user.services."recolor" = lib.mkIf cfg.start-with-systemd {
              Unit = {
                Description = "recolor dynamic theming";
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
