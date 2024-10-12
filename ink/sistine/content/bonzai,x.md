#project

- random.org like funcs
- on bonzai
https://github.com/rwxrob/bonzai
https://github.com/rwxrob/z/blob/main/main.go
- blog posts creation
    - one small note, big note
    - repost to telegram
    - month, year report
- dns lookup with all record types by default
- save cheatsheets, notes for commands/configs/etc. https://github.com/cheat/cheat
    - в `x` и на сайт
    - ascii таблица
    - краткие справки по командам
    - git commands cheatsheets, ascii illustrations
- networking info dump https://github.com/rprtr258/x/pull/3
    
```embed
title: 'Dangit, Git!?!'
image: 'https://dangitgit.com/favicon-32x32.png'
description: 'Git is hard: messing up is easy, and figuring out how to fix your mistakes is impossible. Git documentation has this chicken and egg problem where you can’t search for how to get yourself out of a mess, unless you already know the name of the thing you need to know about in order to fix your problem…'
url: 'https://dangitgit.com/'
```
https://www.kdnuggets.com/publications/sheets/Git_Cheatsheet_KDnuggets.pdf
```embed
title: '7 tips for improving your productivity with Git'
image: 'https://res.cloudinary.com/practicaldev/image/fetch/s--5fM3nWOL--/c_imagga_scale,f_auto,fl_progressive,h_500,q_auto,w_1000/https://dev-to-uploads.s3.amazonaws.com/uploads/articles/m5yq9kgbcn6zr034k7r4.jpg'
description: 'Introduction Git is the most popular source control system with an incredible 93.87% of...'
url: 'https://dev.to/dgenezini/7-tips-for-improving-your-productivity-with-git-ajg'
```

```embed
title: 'GitHub - cheat/cheat: cheat allows you to create and view interactive cheatsheets on the command-line. It was designed to help remind *nix system administrators of options for commands that they use frequently, but not frequently enough to remember.'
image: 'https://opengraph.githubassets.com/995c69c68ab404bcc8cd97853be83c7fde97c34a9602d03b79788a41d0d22987/cheat/cheat'
description: 'cheat allows you to create and view interactive cheatsheets on the command-line. It was designed to help remind *nix system administrators of options for commands that they use frequently, but not ...'
url: 'https://github.com/cheat/cheat'
```

```sql
GRANT SELECT, INSERT ON TABLE users TO appuser -- Grant permissions on a table
REVOKE DELETE ON TABLE users FROM appuser      -- Revoke permission on a table
```

https://github.com/chubin/cheat.sheets/blob/master/sheets/apt-get
https://github.com/chubin/cheat.sheets/blob/master/sheets/convert
https://github.com/chubin/cheat.sheets/blob/master/sheets/df
https://github.com/chubin/cheat.sheets/blob/master/sheets/du
https://github.com/chubin/cheat.sheets/blob/master/sheets/htop
https://github.com/chubin/cheat.sheets/blob/master/sheets/proc
https://github.com/chubin/cheat.sheets/blob/master/sheets/ss
https://github.com/chubin/cheat.sheets/blob/master/sheets/ssh [[reference/devops, sre, admin/devops, sre, admin#ssh]]

- twitch chat integration
    `x twitch chat satont | rg coopert1no | x twitch send satont` repeats `coopert1no` messages
- gtd integration
    - reference dir
    - in dir
- key-value configs, file saving using [https://github.com/charmbracelet/charm](https://github.com/charmbracelet/charm)
- command configuration using [https://github.com/BoRuDar/configuration](https://github.com/BoRuDar/configuration) or similar
- store password, give back by master password
- random music from saved musics
```sh
find /mnt/hdd/music -name '*.mp3' | rg -iv '(BIG RUSSIAN BOSS|oxxxymiron)' | shuf -n1 | python3 -c 'import pathlib;print(pathlib.Path(input()).stem)' | python3 ~/abobus.py
```
- https://devhints.io/