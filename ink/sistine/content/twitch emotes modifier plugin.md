#someday_maybe #project

[link](https://github.com/rprtr258/twitch-emotes-modifier-plugin/pull/1)

## backend
- [x] get 7tv id by emote name on channel ✅ 2023-02-11

## frontend
https://github.com/FrankerFaceZ/Add-Ons
- [ ] userscript/ffz plugin to replace emote strings to images

## filters
- [ ] hats (mostly covered by `hover` with zero-width emotes modifier?)
    - [ ] Top
    - [ ] Santa
    - [ ] comfy
    - [ ] sunglasses(EZ)
    - [ ] snow
    - [ ] rain
    - [ ] glitch https://7tv.app/emotes/60af306712d770149135895e
    - [ ] glitch2 https://7tv.app/emotes/624f8a706bb22d4119fc86f6
    - [ ] [7TV.APP](https://7tv.app/emotes?sortBy=popularity&page=0&filter=%5Bobject%20Object%5D)
- [ ] raw text, e.g. `СРАТЬ` or `ВЕЧНО`
- [ ] unicode emotes
- [ ] modifiers
    - [ ] wide(4HEader) - see `iscalex`
    - [ ] color modifiers
        - [x] grayscale ✅ 2023-02-11
        - [ ] change tone/brightness/contrast
- [ ] crowd https://7tv.app/emotes/63affc7cad551e191f8a0c18
- [ ] animation
    - [x] reverse time ✅ 2023-01-15
    - [x] mirror horizontally/vertically ✅ 2023-01-15
    - [x] `scalex`, `scaley` - resize, `scalet` - speedup/slowdown; particularly `{i,d,}scale{x,y,t}`: ✅ 2023-01-22
        - [x] `i` - `i`ncrease 2 times ✅ 2023-01-22
        - [x] `d` - `d`ecrease 2 times ✅ 2023-01-22
        - [x] empty - scale size is taken from stack ✅ 2023-01-22
    - [x] concat ✅ 2023-01-15
    - [x] rename `concat` to `stackt` ✅ 2023-01-16
    - [x] `stackx`, `stacky` ✅ 2023-01-16
    - [x] shake (x,y=random()) ✅ 2023-02-11
    - [x] rave (hue=time) ✅ 2023-02-11
    - [ ] roll (x=time)
    - [ ] spin
    - [x] slide in = `slide out > revt` ✅ 2023-02-11
    - [x] slide out = `slide in > revt` ✅ 2023-02-11
    - [ ] lurk https://7tv.app/emotes/61fe873323f0a55b0ba891d4
    - [ ] petpet https://github.com/tsoding/emoteJAM/issues/46
    - [ ] bonk https://github.com/tsoding/emoteJAM/issues/41
    - [ ] poof https://7tv.app/emotes/6040aa41cf6746000db1034e
    - [ ] exploding(waytoo) https://github.com/tsoding/emoteJAM/issues/9
    - [ ] linear interpolation from one to another
- [ ] rotations
- [ ] positoning
    - [ ] position using coordinates(e.g. Pog under Champ), maybe not needed, covered by `stacky` and precise coordinates seem not to be useful
- [ ] stack manipulation
    - [x] dup ✅ 2023-01-22
    - [x] swap/xchg ✅ 2023-01-22
    - [ ] see stack languages for examples (e.g. Forth)
- [ ] take emotes from current channel e.g. `peepoHappy` (OR by id e.g. `#b23d1a86cc2abf69ca3f0657b` - only in frontend site, for testing/demonstration purposes)
- [ ] diagrams - trees where childs are source emotes and result is emote with resulting emote from modifier

```embed
title: 'emoteJAM — Generate animated emotes from static images'
image: 'https://t0.gstatic.com/faviconV2?client=SOCIAL&type=FAVICON&fallback_opts=TYPE,SIZE,URL&url=https://github.io/emoteJAM/&size=128'
description: 'Welcome to emoteJAM — a simple website that generates animated BTTV emotes from static images. Let’s get started!'
url: 'https://tsoding.github.io/emoteJAM/'
```

[7tv zero width emotes](https://7tv.app/emotes?page=1&filter=zero_width)
https://api.frankerfacez.com/docs/#emote-effects
[Complete list of tags [Cheat Sheet]](https://www.reddit.com/r/betterponymotes/comments/1y7vel/complete_list_of_tags_cheat_sheet/)
[FOR ALL THOSE LOOKING FOR HALLOWEEN EMOTES HERES A LIST <3](https://www.reddit.com/r/Twitch/comments/j5gqg4/for_all_those_looking_for_halloween_emotes_heres/)
store emotes under identifiers to database as aliases to e.g.
`Kappa-long-shaking-rave-50%-30%-50px/50px`

sunglasses+shake https://github.com/tsoding/emoteJAM/issues/56
`OOOO` = `PogFish` > `-shake`
`sratVechno` = `Basedge,callmehand>stackx,срать,вечно>stacky>stackx`

![[data/static/old/someday_maybe/programming_projects/twitch emotes modifier plugin/PogFish.png]]
![[data/static/old/someday_maybe/programming_projects/twitch emotes modifier plugin/peepoHappy.png]]
![[data/static/old/someday_maybe/programming_projects/twitch emotes modifier plugin/coMMMMfy.png]]
https://7tv.app/emotes/64136af65997ea325d92b9db
https://7tv.app/emotes/645cc58178ffa284f316dd8f
```
ffmpeg -i input.mp4 -filter_complex "[0]reverse[r];[0][r]concat=n=2:v=1:a=0" output.gif
```