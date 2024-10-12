#someday_maybe #project

https://github.com/rprtr258/gena
https://p5js.org/reference/
https://p5js.org/reference/#/p5.Vector

[lang](lang)

# bot buffers commands

https://github.com/rprtr258/buffers

со стороны бота:
    добавление буффера
        e.g. `somebody: !buf add f x+rprtr258.y` вычислит `somebody.f = somebody.x + rprtr258.y`
    web syntax checker
    добавление комментария к буфферу
        e.g. `somebody: !buf desc f Very useful function` поменяет описание буффера `somebody.f` на `"Very useful function"`
    просмотр буффера
        e.g. `somebody: !buf view f` даст ссылку на картинку с отрендеренным буффером `somebody.f`
    добавление буффера с картинкой
        e.g. `somebody: !buf image img https://google.com/pic.png` добавляет буффер `img` с картинкой `pic.png`, пишет информацию о картинке
со стороны сайта (фронтенда / cli):
    просматривать список буфферов
    смотреть визуальное изображение буфферов
        https://docs.rs/image/latest/image/
    смотреть дерево вычисления буффера, начиная от самых базовых буфферов (координат пикселей)
    дебаг вычисления буффера, просмотр промежуточных шагов (в виде картинки или массива чисел)
    вывод полной формулы буффера (загруженные картинки суть константы), просто для красоты
со стороны сервера (бэкенд / cli / library):
    [x] парсинг формулы для буффера в AST
    [ ] вычисление AST
    [ ] сохранение в файл/бд/файловый сервер (на самом деле можно без этого, т.к. каждый буффер (кроме картинок) это в конечном счете формула)
        формат данных для файлов буфферов:
```
data ::= buffer_type shape_len shape buffer_data
buffer_type ::= 0x00 (for float) | 0x01 (for u8) | 0x02 (for bool) | 0x03 (for idx)
shape_len ::= 0x01 - 0xFF
shape ::= (0x01 - 0xFF) x `shape_len` times
buffer_data ::= (0x00000000-0xFFFFFFFF (for floats) | 0x00-0xFF (for u8) | 0x00-0xFE (for idx) | 0x00-0x01 (for bool)) x `shape[0] * shape[1] * ... shape[shape_len-1]` times
```
    [ ] (если нужно) вычисленное значение буффера
    [ ] генерация изображения по буфферу
    [ ] удаление (pruning) промежуточных/временных буфферов (?)

[](https://github.com/rprtr258/various-scripts/blob/master/notebooks/twitch_shaders.ipynb)
examples:
    `(0.4 > max#2 abs (stack#2 x y) - _tmp) * fract 7 * x + y`
[The Book of Shaders](https://thebookofshaders.com/07/)
https://www.shadertoy.com/view/st2yW1
[necessary-disorder tutorials](https://necessarydisorder.wordpress.com/)
https://www.reddit.com/user/xponentialdesign/
https://twitter.com/etiennejcb
https://gist.github.com/Bleuje
https://bleuje.github.io
https://habr.com/en/post/240417/
https://www.reddit.com/r/perfectloops/comments/c6a26z/iris_variations/
https://en.wikipedia.org/wiki/List_of_fractals_by_Hausdorff_dimension
https://www.youtube.com/c/InigoQuilez/videos
[Live Stream #169: Perlin Noise Loops + JS Inheritance](https://www.youtube.com/watch?v=7k-iJyHq7-k)
[Composing Reactive Animations](http://conal.net/fran/tutorial.htm)
vector operations
compression(?) - remove buffer and inline its values into all dependent shaders
filters: if buffer `a` depends on `b`, there could be an ability to copy `a` with replacing `b` to some other buffer `c`, e.g. `a=convolve(b, gaussian_kernel)` is gaussian blur, so after replacing it will be gaussian blur of `c`
fractal functions / transformation degree
operations / computations cost bound - to prevent too heavy computations

https://numpy.org/doc/stable/reference/generated/numpy.einsum.html
# convolution in numpy
https://pastebin.com/V6AunsGU
https://habr.com/ru/post/240417/

https://github.com/rprtr258/various-scripts/blob/master/notebooks/twitch_shaders.ipynb
https://github.com/rprtr258/various-scripts/blob/master/notebooks/Arctangent_plane.ipynb
https://github.com/rprtr258/various-scripts/blob/master/notebooks/bin_pics.ipynb
https://twitter.com/KilledByAPixel/status/1517294627996545024
https://iquilezles.org/articles/
sell NFTs with images produced by this project
https://www.enlight.ru/demo/faq/
np.array([­[1], [2]]) + np.array([3, 4]) = ?

https://github.com/robpike/ivy
https://codegolf.stackexchange.com/a/53129
https://code.jsoftware.com/wiki/Studio/Gallery
http://demoscene.ru/forum/viewtopic.php?t=842&postdays=0&postorder=asc&start=0
![](/static/buffers_apl_example.png)
```python
from colorsys import hsv_to_rgb
from math import sin, cos, pi, fmod
from PIL import Image, ImageDraw

color = "red"
SIZE = 500
R = SIZE // 2

def get_color(k):
    r, g, b = hsv_to_rgb(fmod(k * 10, 360.) / 360, 1., .5)
    return (int(r * 255), int(g * 255), int(b * 255))

def drawer(f):
    def _f(n, k):
        im = Image.new("RGB", (SIZE, SIZE), (255, 255, 255))
        draw = ImageDraw.Draw(im)
        color = get_color(k)
        f(n, k, draw=draw, color=color)
        return im
    return _f

@drawer
def draw_star(n, k, draw, color):
    angle = 2 * pi * k / n
    def point(i):
        ang = angle * i
        return (R + R * cos(ang), R + R * sin(ang))
    draw.line([
        point(i)
        for i in range(n + 1)
    ], width=2, fill=color)

@drawer
def draw_multiplication_table(n, k, draw, color):
    angle = 2 * pi / n
    def point(i):
        ang = angle * i
        return (R + R * cos(ang), R + R * sin(ang))
    t = 1
    for i in range(1, n):
        draw.line(
            (point(i), point(t)),
            width=1,
            fill=color
        )
        t += k
        while t >= n:
            t -= n

n = 10
FRAMES = 100
ks = [n * kk / FRAMES for kk in range(FRAMES)]
images = [
    #draw_multiplication_table(n, k)
    draw_star(n, k)
    for k in ks
]
images[0].save(
    "star.gif",
    "gif",
    save_all=True,
    append_images=images[1:],
    loop=0,
    duration=70.,
    optimize=True
)
```

(ocaml) frontend for compiling GLSL from APL, then render GLSL

color++ mod 255
![](https://hsto.org/getpro/habr/upload_files/ce9/aaf/4d9/ce9aaf4d9f68a11ffcb65d483a621113.gif)
https://www.youtube.com/watch?v=KPoeNZZ6H4s
https://github.com/shader-slang/slang
https://www.youtube.com/watch?v=f4s1h2YETNY
https://demobasics.pixienop.net/tweetcarts/
https://pkg.go.dev/robpike.io/ivy
https://www.scratchapixel.com/index.html
https://habr.com/ru/articles/334580/ simple 3d engine
https://habr.com/ru/articles/763142/
https://www.uiua.org/

[perlin_noise](https://web.archive.org/web/20160530124230/http://freespace.virgin.net/hugo.elias/models/m_perlin.htm)
https://wiki.nikiv.dev/art/generative-art
https://google-research.github.io/dex-lang/
https://www.khoury.northeastern.edu/home/jrslepak/typed-j.pdf

https://inconvergent.net/
https://fronkonstin.com/
https://github.com/aschinchon/cyclic-cellular-automata
https://github.com/armdz/ProcessingSketchs
https://github.com/Mr-Slesser/Generative-Art-And-Fractals
https://github.com/cdr6934/Generative-Processing-Experiments
https://github.com/pkd2512/inktober2017
http://blog.dragonlab.de/2015/03/generative-art-week-1
https://editor.p5js.org/kenekk1/sketches/Ly-5XYvKX
http://paulbourke.net/fractals/peterdejong/
https://editor.p5js.org/kenekk1/sketches/O44Dln5oo
https://openprocessing.org/sketch/1071233
https://twitter.com/okazz_
https://openprocessing.org/sketch/738638
https://openprocessing.org/sketch/1102157
https://openprocessing.org/sketch/1071233
https://openprocessing.org/user/139364
https://openprocessing.org/sketch/792407
https://www.iquilezles.org/www/articles/warp/warp.htm