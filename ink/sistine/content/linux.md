#project

- [x] пофиксить необходимость запускать `.xprofile` при перезагрузке линукса (это невозможно блять) ✅ 2023-01-01
```sh
#!/usr/bin/bash

setxkbmap -layout us,ru -option 'grp:alt_shift_toggle' -option 'ctrl:nocaps' -option 'numpad:microsoft'
parcellite &
just -f ~/cron/justfile wallpaper
```
- [ ] store dotfiles, configs, binary deps, scripts such as update software, etc as `ansible` playbooks, and later rewrite in `mk`
    https://github.com/rprtr258/dotfiles
    https://www.chezmoi.io/quick-start/#start-using-chezmoi-on-your-current-machine
    https://phelipetls.github.io/posts/introduction-to-ansible/ - Ansible for dotfiles: the introduction I wish I’ve had
    https://github.com/rprtr258/rprtr258/blob/master/install.sh + install.md + cron/*

## terminal
https://github.com/alacritty/alacritty/ - no scrollback, Ctrl-L clears scrollback instead moving up, though usable w/ https://zellij.dev/
https://github.com/LukeSmithxyz/st - `st` requires setting up through code, works well, maybe try some time later
https://www.waveterm.dev/

## shells
someday i wanted [my own shell](/note/lisp_shell) if one of below have all features of it, it won't be needed
[nushell](https://www.nushell.sh/)
[Next Generation Shell. A modern programming language for DevOps](https://ngs-lang.org/)
https://elv.sh/ https://elv.sh/learn/tour.html
[python-like shell](https://www.marceltheshell.org/)
https://github.com/curlywurlycraig/crsh
https://github.com/liljencrantz/crush/blob/master/docs/overview.md

https://github.com/hyperupcall/autoenv use arbitrary scripts: `.env` on dir enter, `.env.leave` on dir leave
https://direnv.net/docs/installation.html use env from local `.envrc` files, envs only, no arbitrary bash e.g. alias

https://risor.io/ scripting language
https://janet-lang.org/docs/index.html

https://murex.rocks/ - интересно, не хватает удобных биндов для редактирования типо ^w, ^u не работает, нужно кучу нажатий esc для перехода между режимами автокомплишна, хелпа, превью, предупреждений о вставке в промпт, два раза нужно нажимать ентер чтобы вставить из истории
https://xon.sh/ - сложная семантика, сложный синтаксис, рабочий питон, интегрируется с башем, иногда вылетает, в ранней разработке
https://www.oilshell.org - кал, непонятные переходы между выражениями, подстановками строк, подстановками переменных, нельзя элементарно взять вывод `ls` и разбить на строки, т.к. он сразу же пытается вычислить получающуюся команду типо `"a.py\n" "porn\n" ".bashrc\n" ...`, иногда вылетает, непонятно, зачем нужны C-style строки, когда есть ""-style строки, другие синтаксические приколы унаследованные то ли из bash то ли хз откуда, нет сохранения истории, плохой автокомплит

### interactive history
https://atuin.sh/docs/
https://github.com/ddworken/hishtory
https://github.com/cantino/mcfly
https://www.outcoldman.com/en/archive/2017/07/19/dbhist/

## ricing
[guide to ricing - Google Search](https://www.google.com/search?q=guide+to+ricing)
[Any good guides for a beginner?](https://www.reddit.com/r/unixporn/comments/2gwy5v/any_good_guides_for_a_beginner/)
[Ricing](https://thatnixguy.github.io/posts/ricing/)
https://github.com/fosslife/awesome-ricing

## file managers
https://github.com/sxyazi/yazi
https://dystroy.org/broot/
https://github.com/antonmedv/walk

## vim
https://astronvim.com/Basic%20Usage/walkthrough
https://www.lunarvim.org/docs/installation
https://nvchad.com/docs/quickstart/install
vim setup guides https://www.youtube.com/@chrisatmachine
https://www.lazyvim.org/

## remote file sending
https://github.com/SpatiumPortae/portal
https://github.com/schollz/croc
https://upspin.io/
## uncategorized
[exa](https://github.com/ogham/exa) modern ls
[awesome-cli](https://github.com/agarrharr/awesome-cli-apps#directory-navigation)
[autojump](https://github.com/wting/autojump)
[fzf](https://github.com/junegunn/fzf)
[nb](https://github.com/xwmx/nb) note taking and bookmarking

[googler](https://github.com/jarun/googler)
[s](https://github.com/zquestz/s) - web search in terminal

[up](https://github.com/akavel/up) - Ultimate Plumber is a tool for writing Linux pipes with instant live preview
https://github.com/ericfreese/rat - not working at all, though idea is interesting

[Guix](http://guix.gnu.org/)
[NixOS](https://nixos.org/) https://nixos-and-flakes.thiscute.world/preface

https://github.com/ibraheemdev/modern-unix
[A curated list of command-line utilities written in Rust](https://gist.github.com/sts10/daadbc2f403bdffad1b6d33aff016c0a)
https://dev.to/lissy93/cli-tools-you-cant-live-without-57f6#utils

https://github.com/wader/fq - jq for binary formats
https://github.com/antonmedv/fx - Terminal JSON viewer
https://github.com/noahgorstein/jqp - TUI playground to experiment with jq
https://github.com/tomnomnom/gron

[Как работают snap, flatpak, appimage](https://habr.com/ru/post/673488/)go-global-update
https://github.com/gopasspw/gopass
https://news.ycombinator.com/item?id=34069106
rewrite with pipes
[HH: Programs that saved you 100 hours?](    https://github.com/rprtr258/Most-frequent)-shell-commands
shttps://jonas.git io/tig/doc/manual.html - ncurses frontend for git

https://phelipetls.github.io/posts/introduction-to-ansible/ - Ansible 
https://github.com/TheR1D/shell_gpt

https://github.com/chubin/awesome-console-services
https://github.com/binpash/try - Inspect a command's effects before modifying your live system
https://github.com/89luca89/distrobox
https://evanhahn.github.io/ffmpeg-buddy/
https://github.com/mholt/archiver
https://github.com/BloopAI/bloop/releases `sqlite3 bleep.db 'select exchanges from conversations ORDER BY id DESC LIMIT 1;' | jq '.[0].outcome.Article' | xargs printf | clip`
https://github.com/noborus/ov pager
https://github.com/boyter/cs code fuzzy search
https://github.com/martinvonz/jj better(?) git
https://ffmpeg.app/ https://ffmpeg.lav.io/
[pkgx/run anything](https://pkgx.sh/)
https://dev.to/mebble/ltag-a-little-cli-tool-for-tagged-text-searching-31o3
https://9fans.github.io/plan9port/
https://github.com/muesli/beehive event/trigger system
https://github.com/heyman/heynote/ gui scratchpad