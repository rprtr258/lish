`#23005c`

https://github.com/catppuccin
```json
background,foreground,color1,color2,color3,color4,color5,color6,color7
#3f393b,#f5cbdc,#e98799,#db84ac,#ed90b2,#f8a3bc,#de87c1,#f89ac6,#ab8e9a
#0a0e10,#b9c1cf,#52606a,#586975,#5d7183,#677a8b,#728499,#7a8ca6,#818790
#030309,#beb5d8,#4c4d75,#655676,#82627e,#3b5783,#67658f,#8876a4,#857e97
#120f1f,#e4d0dc,#5b4c66,#a64e68,#655a94,#b26d9d,#8d979c,#c59bb3,#9f919a
#0c0f16,#a5ced9,#384650,#36576c,#505d64,#8e8370,#587685,#6c8c8f,#739097
#0e040c,#aac6e2,#4c516f,#555066,#8e5069,#3e57b4,#566a95,#6a8cb0,#768a9e
#00000c,#d4c5dc,#8a5479,#260593,#632697,#73429d,#9467a0,#a428c9,#94899a
#05080e,#a2b2ba,#243a51,#283f60,#2f4b61,#4c6270,#9c491a,#938e72,#717c82
```

```go
package main

import (
	"fmt"
	"strings"

	"github.com/rprtr258/scuf"
)

func hex(s string) (r, g, b uint8) {
	fmt.Sscanf(s, "#%02x%02x%02x", &r, &g, &b)
	return
}

func printTheme(
	background,
	foreground,
	color1, color2, color3, color4, color5, color6, color7 string,
) {
	bg := scuf.BgRGB(hex(background))
	fmt.Println(scuf.String(foreground, bg, scuf.FgRGB(hex(foreground))))
	fmt.Println(scuf.String(color1, bg, scuf.FgRGB(hex(color1))))
	fmt.Println(scuf.String(color2, bg, scuf.FgRGB(hex(color2))))
	fmt.Println(scuf.String(color3, bg, scuf.FgRGB(hex(color3))))
	fmt.Println(scuf.String(color4, bg, scuf.FgRGB(hex(color4))))
	fmt.Println(scuf.String(color5, bg, scuf.FgRGB(hex(color5))))
	fmt.Println(scuf.String(color6, bg, scuf.FgRGB(hex(color6))))
	fmt.Println(scuf.String(color7, bg, scuf.FgRGB(hex(color7))))
}

func main() {
	for name, theme := range map[string]string{
		"strawberry": "#3f393b,#f5cbdc,#e98799,#db84ac,#ed90b2,#f8a3bc,#de87c1,#f89ac6,#ab8e9a",
		"vixen blue": "#0a0e10,#b9c1cf,#52606a,#586975,#5d7183,#677a8b,#728499,#7a8ca6,#818790",
		"c":          "#030309,#beb5d8,#4c4d75,#655676,#82627e,#3b5783,#67658f,#8876a4,#857e97",
		"pink whore": "#120f1f,#e4d0dc,#5b4c66,#a64e68,#655a94,#b26d9d,#8d979c,#c59bb3,#9f919a",
		"e":          "#0c0f16,#a5ced9,#384650,#36576c,#505d64,#8e8370,#587685,#6c8c8f,#739097",
		"f":          "#0e040c,#aac6e2,#4c516f,#555066,#8e5069,#3e57b4,#566a95,#6a8cb0,#768a9e",
		"g":          "#00000c,#d4c5dc,#8a5479,#260593,#632697,#73429d,#9467a0,#a428c9,#94899a",
		"h":          "#05080e,#a2b2ba,#243a51,#283f60,#2f4b61,#4c6270,#9c491a,#938e72,#717c82",
	} {
		parts := strings.Split(theme, ",")
		fmt.Println(name)
		printTheme(parts[0], parts[1], parts[2], parts[3], parts[4], parts[5], parts[6], parts[7], parts[8])
		fmt.Println()
	}
}
```

[Vim Colors](https://vimcolors.com/712/strawberry-dark/dark)

![](/static/pink_hoe.png)
![](/static/vixen.png)
![](/static/material_palenight.png)
![](/static/sredan.jpg)
![](/static/hollow_knight.jpg)
![](/static/red.png)
![](/static/green.png)
