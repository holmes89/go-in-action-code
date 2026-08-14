# Recommended Terminals for Go Development

A good terminal emulator improves the day-to-day experience of running `go build`, `go test`, and related tools. Below are popular options by platform.

## Cross-Platform

### [Warp](https://www.warp.dev/)
Modern, GPU-accelerated terminal with built-in AI assistance, command history search, and blocks-based output. Available on macOS and Linux.

### [WezTerm](https://wezfurlong.org/wezterm/)
Highly configurable, GPU-accelerated terminal written in Rust. Supports multiplexing, ligatures, and is configured via Lua. Works on macOS, Linux, and Windows.

### [Alacritty](https://alacritty.org/)
Minimalist, GPU-accelerated terminal focused on performance. Cross-platform (macOS, Linux, Windows). Configured via YAML/TOML.

---

## macOS

### [iTerm2](https://iterm2.com/)
The most popular macOS terminal replacement. Offers split panes, search, shell integration, and extensive customization. Pairs well with [oh-my-zsh](https://ohmyz.sh/) or [Starship](https://starship.rs/).

### Terminal.app
The built-in macOS terminal. Perfectly capable for Go development with no setup required.

---

## Linux

### [Kitty](https://sw.kovidgoyal.net/kitty/)
Fast, feature-rich terminal with GPU rendering, tiling layouts, and a powerful extension system ("kittens"). Popular among developers on Linux.

### [GNOME Terminal](https://help.gnome.org/users/gnome-terminal/stable/)
Default terminal for GNOME-based distros (Ubuntu, Fedora). Reliable and easy to use.

### [Konsole](https://konsole.kde.org/)
Default terminal for KDE Plasma. Supports tabs, split views, and extensive profile management.

---

## Windows

### [Windows Terminal](https://aka.ms/terminal)
Microsoft's modern terminal supporting PowerShell, CMD, WSL, and more. Highly recommended for Go development on Windows, especially when combined with [WSL2](https://learn.microsoft.com/en-us/windows/wsl/).

### [WSL2 (Windows Subsystem for Linux)](https://learn.microsoft.com/en-us/windows/wsl/install)
Runs a full Linux environment on Windows. Lets you use Linux-native Go tooling and any of the Linux terminals listed above inside Windows Terminal.

---

## Shell Enhancements

These work inside any terminal and improve the Go development workflow:

| Tool | Purpose | Link |
|------|---------|------|
| [Starship](https://starship.rs/) | Cross-shell prompt with Go version display | https://starship.rs/ |
| [oh-my-zsh](https://ohmyz.sh/) | Zsh framework with plugins/themes | https://ohmyz.sh/ |
| [fzf](https://github.com/junegunn/fzf) | Fuzzy finder for history, files, and more | https://github.com/junegunn/fzf |
| [zoxide](https://github.com/ajeetdsouza/zoxide) | Smarter `cd` for navigating projects quickly | https://github.com/ajeetdsouza/zoxide |
