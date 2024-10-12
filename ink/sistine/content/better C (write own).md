[spec](spec)

zero runtime - in compiled binary there is no additional info about compiler, debug symbols, runtime, etc, if not explicitly added

generic types are functions of types:
```go
[](T) // slice of elements of type T
chan(T) // channel of elements of type T
map(K, V) // map with keys of type K and values of type V
struct{
    fieldA: A,
    fieldB: B,
    ...
} // struct with fieldA of type A and fieldB of type B
vector(T) = (T == bool ? [](int) with bit-mapping : [](T))
```

macroses
https://danielkeep.github.io/tlborm/book/README.html

`assert.Equal(expected T, actual T)` args might be differentiated based on if `expected` is known(computable) at compile time

compiled printf, scanf, etc.

metadata: e.g. position in AST, all positions of macros application(?)

pure functions: all instances are inlined and checked to call only other pure functions

VERY consistent syntax
```php
x :int // variable declaration without initialization
x :int = 3 // variable declaration with initialization
x := 3 // variable declaration with initialization and type inference
x = 3 // assignment
y := {
  field = 3, // struct field assignment
}
```

routing

```php
macro route(Macro applications) {
  fn _route(Request r) {
    `switch r.path {
      ${for p in applications.places {
                p.function(r);
            }}
    }`
  }
}

@route("/")
fn index(Request r) Response {...}

@route("/login")
fn login(Request r) Response {...}
```

[Examples](https://docs.openvalidation.io/examples)

lambdas?

defer https://github.com/lichray/deferxx

with block

allocate/deallocate memory

open/close file/socket/etc

type synonym

after `typedef uint32 index` every use of `uint32` as `index` must be explicitly casted

dot notation for functions

```cpp
obj.do(1, 'c')
// is just
do(&obj, 1, 'c')
```

modules

import only certain object from another module

postfix type notation

```c
typedef sighandler = fn(int)->(); // ->void
// instead of
typedef void (*sighandler)(int);

sighandler signal(int signum, sighandler_t handler);
```

struct modification
```rust
x : Struct {x, y : int}
x.{
    .x = 2;
    z := 5;
    .y = .x + z;
}
```

loops

do while?

iter-while(name?)

```cpp
iter {
    c = getchar();
} while (c != EOF) {
    putchar(c ^ ' ');
}
// instead of
while ((c = getchar()) != EOF) {
    putchar(c ^ ' ');
}
// and
c = getchar();
while (c != EOF) {
    putchar(c ^ ' ');
    c = getchar();
}
```

infinite loop?

for-each loop?

`{}` are mandatory in `if`, `while`, `for`, `do-while`

switch case without breaks and cases in blocks

blocks return values

no braces in `if`, `while`, `for`, `do-while`

error handling

errors are either errno with info in global `void* info` or explicit in function signature, it is required to check error before using value

function signature containing possible errno values, when calling function use `try` like:

```cpp
FILE* fopen(const char*, Filemode) throws (NOT_EXISTS, NO_RIGHTS);
// somewhere in fopen
// if (not_exists(filename))
//     throw NOT_EXISTS; // sets errno to NOT_EXISTS, info to NULL and returns

FILE* fd = try fopen("file.txt", READ) {
  NOT_EXISTS {
    panic("file %s does not exist", "file.txt");
  }
  NO_RIGHTS {
    panic("no permissions to open %s for read", "file.txt");
  }
};
```

safety flag:
```c
@safety
get(arr []int, i int) int {
    if @safety && (i < 0 || i >= len(arr)) {
        panic("{i} is out of bounds 0..{len(arr)-1}");
    }
    return arr[i];
}

// then call like so
get(a, i); // no safety check
safe get(a, i); // safety check
```

on jai:
https://github.com/Ivo-Balbaert/The_Way_to_Jai
https://github.com/BSVino/JaiPrimer/blob/master/JaiPrimer.md
https://inductive.no/jai/
https://github.com/Jai-Community/Jai-Community-Library/wiki

[Андрей Викторович Столяров: сайт автора](http://www.stolyarov.info/guestbook#comment-6168)

[Cyclone: Wiki](http://cyclone.thelanguage.org/wiki/)

[What non-theoretical, practical programming language has no reserved keywords?](https://softwareengineering.stackexchange.com/questions/189699/what-non-theoretical-practical-programming-language-has-no-reserved-keywords)

[Why does a programming language need keywords?](https://stackoverflow.com/questions/2452365/why-does-a-programming-language-need-keywords)

[Overview](https://odin-lang.org/docs/overview/#advanced-idioms)

[How I program C](https://www.youtube.com/watch?v=443UNeGrFoM)

[Forge](https://gamepipeline.org/forge.html)

compile sql statements into structure of request result, instead of ORM

[x86 calling conventions - Wikipedia](https://en.wikipedia.org/wiki/X86_calling_conventions#cdecl)

[ANSI C Yacc grammar](https://www.lysator.liu.se/c/ANSI-C-grammar-y.html)

[C2x: будущий стандарт C](https://habr.com/en/company/badoo/blog/503140/)

[GitHub - PyvesB/asm-game-of-life: Assembly implementation of Conway's Game of Life, using NASM assembler for Linux x86-64.](https://github.com/PyvesB/asm-game-of-life)

[NASM Assembly Language Tutorials - asmtutor.com](https://asmtutor.com/)

[NASM - The Netwide Assembler](https://www.nasm.us/xdoc/2.13.01/html/nasmdoc0.html)

train writing parsers with

[If++](If++%2040ed2.md)

[THE DESIGN OF AN OPTIMIZING COMPILER](better%20C%20(%2028f84/file.pdf)

THE DESIGN OF AN OPTIMIZING COMPILER

[List of C-family programming languages - Wikipedia](https://en.wikipedia.org/wiki/List_of_C-family_programming_languages)

[ANTLR](https://www.antlr.org/)

[Си должен умереть](https://cmustdie.com/)

[comp.lang.c Frequently Asked Questions](http://c-faq.com/index.html)

[http://rustmustdie.com/](http://rustmustdie.com/)

[Meet Google Drive - One place for all your files](https://drive.google.com/drive/folders/1pYXIyW9Pr4C-8nGasqSr56aKz2NcGmn9?usp=sharing)

[Assembly language - Wikipedia](https://en.wikipedia.org/wiki/Assembly_language)

```c
#define СНЕСК(LINE, EXPECTED) \
    int гс = LINE; \
        if (гс != EXPECTED) \
            ut_abort(FILE__, LINE #LINE, гс, EXPECTED); }
void ut_abort(char *file, int In, char *line, int rc, int exp) {
    fprintf(stderr, "%s(%d)", file, In);
    fprintf(stderr, "    '%s': expected %d, got %d\n", line, exp, rc);
    exit(1);
}

// Тогда вы можете инкапсулировать вызовы, которые никогда подведут, с помощью строки:
CHECK(stat('/tmp", &stat_buff), 0);
// Если бы это не удалось, то вы получаете сообщение, записанное в stderr:
source.с(19)
    'stat("/tmp", &stat_buff)': expected 0, got -1
```

[https://git.musl-libc.org/cgit/musl/tree/](https://git.musl-libc.org/cgit/musl/tree/)

[Untitled](better%20C%20(%2028f84/Untitled%20D%2051f6a.csv)

[parsing materials](better%20C%20(%2028f84/parsing%20ma%20e20f9.md)

[](https://grosskurth.ca/bib/2006/dolstra-thesis.pdf)
https://talks.golang.org/2012/splash.article#TOC_5.
https://en.cppreference.com/w/c/language/generic

https://habr.com/en/company/ruvds/blog/583576/
https://habr.com/en/post/88101/
https://habr.com/en/post/137706/
https://www.opennet.ru/opennews/art.shtml?num=57081

https://engineering.purdue.edu/ece264/16au/hw/HW14
https://engineering.purdue.edu/ece264/16au/hw/HW15
https://commandcenter.blogspot.com/2012/06/less-is-exponentially-more.html
https://doc.rust-lang.org/cargo/reference/features.html#the-features-section
https://harelang.org/
https://ericniebler.github.io/range-v3/index.html
https://eax.me/cpp-smart-pointers/
https://eax.me/cpp-will-never-die/
https://eax.me/avoid-metaprogramming/
https://llvm.org/docs/tutorial/
https://www.muppetlabs.com/~breadbox/software/tiny/teensy.html
https://eax.me/c-vs-cpp/
https://www.sigbus.info/how-i-wrote-a-self-hosting-c-compiler-in-40-days
https://www.google.com/search?q=ld+and+gold+difference
https://docs.microsoft.com/ru-ru/cpp/code-quality/understanding-sal?view=msvc-170
https://ru.wikipedia.org/wiki/Cyclone_(%D1%8F%D0%B7%D1%8B%D0%BA_%D0%BF%D1%80%D0%BE%D0%B3%D1%80%D0%B0%D0%BC%D0%BC%D0%B8%D1%80%D0%BE%D0%B2%D0%B0%D0%BD%D0%B8%D1%8F)
https://en.wikipedia.org/wiki/Limbo_(programming_language)
https://vector-of-bool.github.io/2019/01/27/modules-doa.html
https://github.com/c3lang/c3c
https://sparrow-lang.readthedocs.io/en/latest/index.html
https://github.com/grassator/mass
https://fpl.handmade.network/
https://jiyu.handmade.network/
https://handmade.network/p/177/cuikc/
https://cakelisp.handmade.network/
https://www.youtube.com/watch?v=335ylTUlyng
https://go.googlesource.com/proposal/+/master/
https://habr.com/en/post/667164/
https://kotlinlang.org/docs/basic-syntax.html#program-entry-point
https://habr.com/en/post/270379/
https://habr.com/en/post/47878/

Представьте себе язык программирования на котором бы вам хотелось писать, такой что бы некоторый круг задач решался бы на нем наиболее эффективно по параметрам:
- скорости решения/написания, т.е. времени которое будет потрачено программистом,
- легкости, т.е. умственных усилий которые будут потрачены программистом,
- дизайн ориентированный на избавление от избыточной сложности,
- размера получаемых программ,
- размера используемой памяти при выполнении,
- скорости выполнения,
- пригодности и для системного программирования (например путем описания упрощенного подмножества)
- пригодности и для самого высокоуровнего программирования (например путем описания надмножества языка, даже допускаю введение раздельных режимов компиляции для любого модуля, т.е. модуль на низко или высоко уровневом подмножестве языка может быть написан и соответственно скомпилирован),
- созданы все условия для надежного и безопасного программирования.

И во вторую очередь, что бы этот круг задач которые решаются эффективно - был бы как можно больше.

- Oberon
- Nim
- Seed7
- D
- Ruby
- Basic (Euphoria, FreeBasic, LangMF)

библиотеки не должны зависеть друг от друга.

строки:
    строки, владеющие данными
        поиск, разбиение на токены
    строки, не владеющие данными

setjmp, longjmp based error handling

отдельные библиотечные функции с thread-safety и без thread-safety

https://stackoverflow.com/questions/479207/how-to-achieve-function-overloading-in-c
https://github.com/nothings/stb/blob/master/docs/stb_howto.txt

on compile time computation use cases
    https://nothings.org/gamedev/optimal_static_grid_pathfinding.txt

strlen
    https://cvsweb.openbsd.org/cgi-bin/cvsweb/src/lib/libc/string/strlen.c?rev=1.9&content-type=text/x-cvsweb-markup
    https://git.musl-libc.org/cgit/musl/tree/src/string/strlen.c
    https://svnweb.freebsd.org/base/head/lib/libc/string/strlen.c?revision=333449&view=markup
    https://sourceware.org/git/?p=glibc.git;a=blob;f=string/strlen.c;h=5d9366c2c0e09abc1c948ef00e2d843e3a28d909;hb=HEAD















see ocaml compilers/transpilers/iterpreters/wtfers
semicolon putting macro - https://go.dev/ref/spec#Semicolons
https://golangci-lint.run/usage/linters/
https://github.com/golang/go/tree/master/src/cmd/compile/internal/ssa/gen
https://github.com/golang/go/issues/32437
https://github.com/coreutils/coreutils/tree/master/src
https://habr.com/ru/post/672282/
https://www.youtube.com/watch?v=TH9VCN6UkyQ
https://gsd.web.elte.hu/lectures/bolyai/2018/grin/grin-pres.pdf
https://github.com/gingerBill/gb/
https://spin.atomicobject.com/2014/09/03/visualizing-garbage-collection-algorithms/

https://github.com/rswier/c4
https://compilers.iecc.com/crenshaw/
https://github.com/lotabout/write-a-C-interpreter
https://pvs-studio.com/ru/blog/posts/0908/
http://craftinginterpreters.com/
https://github.com/carbon-language/carbon-lang
https://github.com/carp-lang/Carp
http://maibriz.de/unix/ultrix/_root/porttour.pdf
https://en.wikipedia.org/wiki/Sethi%E2%80%93Ullman_algorithm
https://github.com/makuto/cakelisp/blob/master/doc/Cakelisp.org
https://github.com/rsashka/newlang
https://habr.com/ru/company/skillfactory/blog/688078/
https://developers.redhat.com/articles/2022/09/01/3-essentials-writing-linux-system-library-rust#
https://blog.replit.com/langjam
https://replit.com/talk/announcements/PL-Jam-Results/57498
https://www.youtube.com/watch?v=8Ab3ArE8W3s
https://tmewett.com/c-tips/
https://dannorth.net/2022/02/10/cupid-for-joyful-coding/
https://citw.dev/tutorial/create-your-own-compiler
https://fennel-lang.org/
http://wiki.call-cc.org/man/5/The%20User%27s%20Manual

```embed
title: 'Writing a toy compiler with Go and LLVM'
image: 'https://i.imgur.com/eQjSgoG.png'
description: 'adventures with LLVM IR'
url: 'https://ketansingh.me/posts/toy-compiler-with-llvm-and-go/'
```

[[data/static/old/someday_maybe/programming_projects/better C/The design of an optimising compiler.pdf]]
[[data/static/old/someday_maybe/programming_projects/better C/Metaprogramming in Modern Programming Languages.pdf]]
[[data/static/old/someday_maybe/programming_projects/better C/Modern Compiler Implemenation in C.djvu]]
[[data/static/old/someday_maybe/programming_projects/better C/r4rs.pdf]]
[[data/static/old/someday_maybe/programming_projects/better C/Парадигмы программирования и алгебраический подход к построению универсальных языков программирования.pdf]]
[[data/static/old/someday_maybe/programming_projects/better C/Advanced Compiler Design and Implementation.pdf]]
[[data/static/old/someday_maybe/programming_projects/better C/Engineering a Compiler.pdf]]
[[data/static/old/someday_maybe/programming_projects/better C/programming languages comparison.csv]]
[[data/static/old/someday_maybe/programming_projects/better C/Programming Language Pragmatics.pdf]]
[[data/static/old/someday_maybe/programming_projects/better C/The C Programming Language.pdf]]

[[data/static/old/someday_maybe/programming_projects/better C/Compiler_Construction.pdf]]

```embed
title: 'Немного о семантиках перемещения, копирования и заимствования'
image: 'https://habrastorage.org/getpro/habr/upload_files/745/8d9/a5c/7458d9a5c64e1c1a716385a5bf9ec4a5.png'
description: 'Существует три основных способа передачи данных в функции: перемещение (move), копирование (copy) и заимствование (borrow, иными словами, передача по ссылке). Поскольку изменяемость (мутабельность)…'
url: 'https://habr.com/ru/company/otus/blog/713910/'
```

```embed
title: 'Terra'
image: 'https://terralang.org/logo.png'
description: 'Terra is a low-level system programming language that is embedded in and meta-programmed by the Lua programming language:'
url: 'https://terralang.org/'
```

```embed
title: 'Models of Generics and Metaprogramming: Go, Rust, Swift, D and More - Tristan Hume'
image: 'https://thume.ca/assets/themes/thume/images/bubble-110.png'
description: 'In some domains of programming it’s common to want to write a data structure or algorithm that can work with elements of many different types, such as a generic list or a sorting algorithm that only needs a comparison function. Different programming languages have come up with all sorts of solutions…'
url: 'https://thume.ca/2019/07/14/a-tour-of-metaprogramming-models-for-generics/'
```

https://research.swtch.com/telemetry
[[data/programming ideas]]
https://www.youtube.com/watch?v=QwuHX-VjhXI
https://ru.wikipedia.org/wiki/Objective-C
https://github.com/TinyCC/tinycc
http://progopedia.ru/
https://github.com/carbon-language/carbon-lang/tree/trunk/proposals
https://go.dev/talks/2015/gophercon-goevolution.slide#1
https://habr.com/ru/post/721776/
https://werat.dev/blog/what-a-good-debugger-can-do/
https://habr.com/ru/post/723400/
https://github.com/SerenityOS/jakt
https://habr.com/ru/post/724010/
https://www.youtube.com/playlist?list=PLmV5I2fxaiCIZVTLzofsocka2LvWBFvBa
https://nothings.org/stb_ds/
https://github.com/sam-astro/Z-Sharp
https://vcpkg.io/en/
https://ballerina.io/learn/slides/language-walkthrough/Ballerina_Language_Presentation-2021-03-08.pdf
http://lucacardelli.name/Papers/DataOrientedLanguages.A4.pdf
https://www.infoq.com/articles/ballerina-data-oriented-language/
https://habr.com/ru/company/ruvds/blog/721686/
https://habr.com/ru/post/726410/
    https://github.com/stsaz/ffos
    https://github.com/stsaz/ffbase

https://news.ycombinator.com/item?id=35436600

[Google C++ Style Guide](https://google.github.io/styleguide/cppguide.html)

http://www.stolyarov.info/guestbook/archive/6#comment-6168
```
Вот уж ни фига подобного. Даже если применять традиционный подход с фиксацией в компиляторе всего, чего только можно, и даже если не думать о высокоуровневых абстракциях и их построении библиотечными средствами (естественно, именно библиотечными, а не средствами компилятора) — то даже в чистом Си есть что исправить: во-первых, ввести массивы как самостоятельную сущность, во-вторых, описания сделать линейными, т.е. читаемыми слева направо, а не как сейчас, да и много чего ещё. В первую очередь — ссылки, это из всего Си++ настоящий семантический бриллиант.

Останавливаться на этом никто не заставляет, следует ввести (хотя бы для структур) конструкторы с деструкторами и переопределение операций (да хотя бы даже присваивания), и вот тут нужно вовремя затормозить, поскольку виртуальность — это уже механизм достаточно сложный, чтобы в языке ему было не место. Дальше очевидным образом возникает некая операция вызова виртуального метода (в отличие от обычного), а её определение отдаётся на откуп библиотекам, то есть построение _vmt и вот это вот всё — должно быть там, а не здесь, и библиотек таких должно быть больше одной, хотя бы даже для того, чтобы множественное наследование можно было поддерживать или не поддерживать. Такой язык УЖЕ будет лучше, чем Си и Си++ (даже вместе взятые, пожалуй).

Следующий шаг — осознание того, что с перегрузкой операций нужно разбираться целиком во время компиляции, следовательно — перегруженные операции не должны быть функциями, они должны быть макросами, и конструкторы, кстати, тоже; отсюда следует, что макропроцессор должен быть не такой, разумеется, как в Си/Си++, а _нормальный_, работающий не до компиляции, а во время таковой, т.е. имеющий доступ к идентификаторам, их категориям, к системе типов и т.п. — вообще ко всей информации, какая есть у компилятора. Заодно избавляемся от манглинга. А ещё (ну это уже мелочь) this должен быть ссылкой, а не указателем.

Пока не ушли далеко от нижнего уровня, можно (и нужно) реифицировать (сделать доступным на уровне конструкций языка) стек — вот буквально снабдить язык возможностью напрямую работать со стековыми фреймами. Чтобы с каждой функцией был связан как-то там называющийся тип вроде структурного, соответствующий её стековому фрейму со всеми локальными переменными. Это чтобы фреймы можно было помечать — если кому-то потребуются исключения, то пусть они будут библиотечной возможностью, а не частью языка.

Следующий шаг — понять, что вообще-то не должно быть такого дебильного понятия, как "ключевое слово". Например, все слова, введённые компилятором, могут начинаться с какого-нибудь знака доллара. Или вообще с обратного слэша, как в TeX'е. Дальше — что вообще-то набор операций и их приоритеты компилятором фиксировать не надо, достаточно на уровне компилятора предусмотреть некое «применение функции/псевдофункции (читай — макроса) к кортежу параметров» и средства для построения этого кортежа, а с остальным справятся макросы, введённые библиотекой. Такие слова, как if, while и прочее, можно вводить как имена макросов. В итоге тот "язык", который видит программист, будет сформирован целиком библиотеками макросов, и таких вариантов будет много, но компилятор будет для всех один и объектные модули будут одного и того же вида, без всяких, естественно, сложностей с их линковкой и взаимными обращениями. А "уровень" тут будет доступен абсолютно любой, от почти ассемблера до почти питона (с точностью до стратегии исполнения, всё-таки питон интерпретируемый, а тут будет, естественно, чистая компиляция).

Здесь я бы, пожалуй, остановился, хотя можно и дальше пройти.

Я это всё к тому, что сугубо низкоуровневый язык может, во-первых, совершенно не быть "новым Си", а во-вторых, может содержать в себе не сами высокоуровневые абстракции, а средства для их создания. Собственно, Си++ подавал на это надежды, причём всерьёз. Примерно до того момента, как за него Степанов принялся.
```
https://www.hillelwayne.com/post/learning-a-language/
https://habr.com/ru/articles/730686/
https://amrdeveloper.github.io/Amun/
https://matklad.github.io/2023/05/02/implicits-for-mvs.html
https://vale.dev/ https://www.val-lang.dev https://whiley.org/learn/
https://ziglang.org/
https://github.com/true-grue/Compiler-Development
https://www.modular.com/mojo
https://habr.com/ru/articles/735152/
https://interpreterbook.com https://compilerbook.com/
https://www.youtube.com/watch?v=MnctEW1oL-E
https://ketansingh.me/posts/toy-compiler-with-llvm-and-go/
https://github.com/aalhour/awesome-compilers?tab=readme-ov-file#go
https://github.com/contextfreeinfo/rio/wiki/Motivation
https://www.reddit.com/r/ProgrammingLanguages/comments/7q8m4m/systems_language_that_compiles_fast_and_has_no_gc/

https://avestura.dev/blog/ideal-programming-language
https://quasilyte.dev/blog/post/c-broken-defaults/
https://quasilyte.dev/blog/post/naive-ssa-alternative/#where-exactly-to-insert-a-varkill
https://github.com/graphitemaster/codin odin -> c compiler
https://habr.com/ru/articles/741124/
http://thalassa.croco.net/doc/cpp_subset.html
https://gitlab.com/cznic/go0
https://github.com/rui314/chibicc
https://antelang.org/
http://stolyarov.info/guestbook/archive/8/#cmt569