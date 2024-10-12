## data oriented design

[CppCon 2014: Mike Acton "Data-Oriented Design and C++"](https://www.youtube.com/watch?v=rX0ItVEVjHc)
[Data-Oriented Design](https://www.dataorienteddesign.com/dodbook/)
[http://curtclifton.net/papers/MoseleyMarks06a.pdf](http://curtclifton.net/papers/MoseleyMarks06a.pdf)
https://odin-lang.org/docs/overview/
https://www.youtube.com/watch?v=yy8jQgmhbAU
https://github.com/CppCon/CppCon2018/blob/master/Presentations/oop_is_dead_long_live_dataoriented_design/oop_is_dead_long_live_dataoriented_design__stoyan_nikolov__cppcon_2018.pdf - last slide
https://github.com/dbartolini/data-oriented-design

```embed
title: 'Полиморфные структуры данных и производительность'
image: 'https://habrastorage.org/getpro/habr/upload_files/7c5/040/fcb/7c5040fcbca6df583d949a5b29b1c049.png'
description: 'В этой статье мы рассмотрим как обычно происходит работа с динамическим полиморфизмом, где теряется производительность и как её можно улучшить, используя интересные структуры данных. В С++ не так...'
url: 'https://habr.com/ru/post/703666/'
```

https://media.handmade-seattle.com/practical-data-oriented-design/

## pprof - golang profiler

```embed
title: 'Профилирование и оптимизация программ на Go'
image: 'https://habr.com/share/publication/301990/038c265a4ffa0169f8905c49fc19bb57/'
description: 'Введение В этой статье я расскажу, как профилировать и оптимизировать приложения на языке Go с использованием встроенных и общих инструментов, доступных в ОС Linux. Что такое профайлинг и...'
url: 'https://habr.com/en/company/badoo/blog/301990/'
```

```embed
title: 'Profiling Go programs with pprof'
image: 'https://jvns.ca/images/pprof.png'
description: 'Profiling Go programs with pprof'
url: 'https://jvns.ca/blog/2017/09/24/profiling-go-with-pprof/'
```
https://habr.com/ru/companies/badoo/articles/301990/
https://www.ardanlabs.com/blog/2023/07/getting-friendly-with-cpu-caches.html
https://www.benburwell.com/posts/flame-graphs-for-go-with-pprof/
https://tproger.ru/translations/memory-leaks-investigation-in-go-using-pprof/
https://github.com/rakyll/autopprof

## sql

```embed
title: 'MySQL :: MySQL 8.0 Reference Manual :: 8 Optimization'
image: 'https://labs.mysql.com/common/themes/sakila/favicon.ico'
description: 'The world’s most popular open source database'
url: 'https://dev.mysql.com/doc/refman/8.0/en/optimization.html'
```

```embed
title: 'Techniques for Optimising SQL Queries'
image: 'https://miro.medium.com/max/1200/1*IOupJk-yndymZusZiIjHuw.png'
description: 'Get the most of your queries'
url: 'https://itnext.io/techniques-for-optimising-sql-queries-c362dbe626b4'
```

```embed
title: 'Магия оптимизации SQL запросов'
image: 'https://habrastorage.org/getpro/habr/upload_files/bb6/d14/597/bb6d1459751aa11fa7b01d9016dc22a6.png'
description: 'Привет, Хабр! Думаю, каждый хоть раз использовал команду explain или хотя бы слышал про нее. Эта команда демонстрирует план выполнения запроса, но как именно СУБД приходит к нему остается загадкой....'
url: 'https://habr.com/ru/post/709898/'
```

```embed
title: 'Использование EXPLAIN. Улучшение запросов'
image: 'https://habr.com/share/publication/211022/ace6e7381907c9f979687914a74260b2/'
description: 'Когда вы выполняете какой-нибудь запрос, оптимизатор запросов MySQL пытается придумать оптимальный план выполнения этого запроса. Вы можете посмотреть этот самый план используя запрос с ключевым...'
url: 'https://habr.com/ru/post/211022/'
```

## unsorted

```embed
title: 'Высокопроизводительные вычисления: проблемы и решения'
image: 'https://habr.com/share/publication/117021/99acb74982aec2da3c0bf43fe4555422/'
description: 'Компьютеры, даже персональные, становятся все сложнее. Не так уж давно в гудящем на столе ящике все было просто — чем больше частота, тем больше производительность. Теперь же системы стали...'
url: 'https://habr.com/en/post/117021/'
```

https://nnethercote.github.io/perf-book/io.html

```embed
title: 'Почему не все тестовые задания одинаково полезны: С++ edition'
image: 'https://habrastorage.org/getpro/habr/upload_files/a15/49a/f49/a1549af497e2d970f88fea9a57b0fb8e.png'
description: 'Вначале было слово, и было два байта, и ничего больше не было. Но отделил Бог ноль от единицы, и понял, что это хорошо. Потом, опуская некоторые незначительные события мироздания, была вот эта статья...'
url: 'https://habr.com/en/post/572726/'
```

[Ультра быстрый Cron с шагом в миллисекунду, или когда тестовые задания такими прикидываются](https://habr.com/en/post/589667/)
[Intrinsics Guide](https://www.laruence.com/sse/)
[[data/static/old/someday_maybe/perfomance/x86 intrinsics cheat sheet v1.0.pdf]]
[GitHub - dendibakh/perf-ninja: This is an online course where you can learn and master the skill of low-level performance analysis and tuning.](https://github.com/dendibakh/perf-ninja)
[Performance Analysis and Tuning on Modern CPUs: Squeeze the last bit of performance from your application.](https://www.amazon.com/dp/B08R6MTM7K/ref=as_li_ss_tl?ref_=pe_3052080_397514860&linkCode=sl1&tag=dendibakh-20&linkId=5ed2c1237c2b8d3c35b482aa6c47e6ce&language=en_US)
[Векторизация](https://algorithmica.org/ru/sse)
[FizzBuzz по-сениорски](https://habr.com/en/post/540136/)

```embed
title: 'Optimising a small real-world C++ application - Hubert Matthews [ACCU 2019]'
image: 'https://img.youtube.com/vi/fDlE93hs_-U/maxresdefault.jpg'
description: '#Cpp #ACCU #ACCUConfThis is a hands-on demonstration of optimising a small real-world application written in C++.It shows how measurement tools such as strac...'
url: 'https://www.youtube.com/watch?v=fDlE93hs_-U'
```

## effective bttv (vptree?)
[compare languages performance (BTTV VPTree search?)](%D0%A1%D0%B0%D0%BC%D0%BE%D0%BE%D0%B1%D1%83%D1%87%D0%B5%D0%BD%208936d/compare%20la%20642cf.md)
https://blog.burntsushi.net/transducers/

## compare cache efficiency of heap and my heap
my heap is array [L + x + R], where L, R are subtrees of x
use advent of code solution

## fast rot13
(try to use memory map)
https://eax.me/linux-file-mapping/
[GitHub - rprtr258/rot13: Trying to make fast rot13. Because if rot13 can't be made fast easily, what could be?](https://github.com/rprtr258/rot13)
make fast rot13, bench on 1MB/GB file, compare with cat(no ciphering) and tr(with rot13), check strace 
[fastest rot13 implementation - Google Search](https://www.google.com/search?q=fastest+rot13+implementation)
[](https://hea-www.harvard.edu/~fine/Tech/rot13.html) 
[Looking at The Source Code for Function isalpha, isdigit, isalnum, isspace, islower, isupper, isxdigit, iscntrl, isprint, ispunct, isgraph, tolower, and, toupper in C Programming](https://grandidierite.github.io/looking-at-the-source-code-for-function-isalpha-isdigit-isalnum-isspace-islower-isupper-isxdigit-iscntrl-isprint-ispunct-isgraph-tolower-and-toupper-in-C-programming/)
[setvbuf](http://www.cplusplus.com/reference/cstdio/setvbuf/) 
[coreutils/tr.c at master · coreutils/coreutils](https://github.com/coreutils/coreutils/blob/master/src/tr.c)
fast io 
[https://www.codechef.com/viewsolution/32916606](https://www.codechef.com/viewsolution/32916606) 
[Solution: 34266958 | CodeChef](https://www.codechef.com/viewsolution/34266958) 
[](https://marc.info/?l=linux-kernel&m=95496636207616&w=2)
[Advent of Code](https://adventofcode.com/2017)

## fast advent of code
[Solving Advent of Code 2020 in under a second](https://timvisee.com/blog/solving-aoc-2020-in-under-a-second/)
[https://www.forrestthewoods.com/blog/solving-advent-of-code-in-under-a-second/](https://www.forrestthewoods.com/blog/solving-advent-of-code-in-under-a-second/) 
[[All years, all days] Solve them within the time limit](https://www.reddit.com/r/adventofcode/comments/7m9mg8/all_years_all_days_solve_them_within_the_time/)
[GitHub - Voltara/advent2018-fast: Advent of Code 2018 optimized solutions in C++](https://github.com/Voltara/advent2018-fast) 
[GitHub - Voltara/advent2017-fast: Advent of Code 2017 optimized solutions in C](https://github.com/Voltara/advent2017-fast)
[CRDTs go brrr](https://josephg.com/blog/crdts-go-brrr/) 
[crdt](https://crdt.tech)
[https://herbsutter.com/welcome-to-the-jungle/](https://herbsutter.com/welcome-to-the-jungle/)
[](https://people.freebsd.org/~lstewart/articles/cpumemory.pdf) 
[stolyarov: asm → C](%D0%A1%D0%B0%D0%BC%D0%BE%D0%BE%D0%B1%D1%83%D1%87%D0%B5%D0%BD%208936d/stolyarov%20%20d68b3.md)
[https://stackoverflow.com/questions/49947915/assembly-syscalls-in-64-bit-windows](https://stackoverflow.com/questions/49947915/assembly-syscalls-in-64-bit-windows) 
[Продуманная оптимизация](https://optimization.guide/)

## io-uring
https://unzip.dev/0x013-io_uring/
https://developers.mattermost.com/blog/hands-on-iouring-go/
https://unixism.net/loti/index.html
https://github.com/ii64/gouring/ https://github.com/godzie44/go-uring https://github.com/pawelgaczynski/gain
https://habr.com/ru/articles/589389/ https://habr.com/ru/articles/597745/ https://habr.com/ru/articles/649161/

https://go101.org/optimizations/101.html
https://go.dev/blog/pprof
https://github.com/alphadose/ZenQ

https://habr.com/ru/company/first/blog/442738/
https://github.com/sharkdp/hyperfine
https://github.com/miloyip/dtoa-benchmark
https://github.com/miloyip/itoa-benchmark
http://www.0x80.pl/articles/sse-itoa.html
http://0x80.pl/articles/simd-parsing-int-sequences.html
https://stackoverflow.com/questions/4351371/c-performance-challenge-integer-to-stdstring-conversion
https://graphics.stanford.edu/~seander/bithacks.html#IntegerLog10
https://git.musl-libc.org/cgit/musl/tree/
https://www.opennet.ru/man.shtml?topic=getc_unlocked&category=3&russian=0
http://www.open-std.org/jtc1/sc22/wg21/docs/papers/2016/p0479r0.html
https://eigen.tuxfamily.org/index.php?title=Main_Page
https://eax.me/c-cpp-profiling/
https://habr.com/ru/post/664044/
https://travisdowns.github.io/
https://benchmarksgame-team.pages.debian.net/benchmarksgame/index.html
https://habr.com/ru/post/673508/
https://www.google.com/search?q=golang%20lock%20free
https://habr.com/ru/post/679008/
https://habr.com/ru/post/682080/
https://habr.com/ru/post/682332/
https://habr.com/ru/post/686222/
https://habr.com/ru/post/688874/
https://github.com/quasilyte/qbenchstat
https://doc.rust-lang.org/std/simd/struct.Simd.html

[[data/static/old/someday_maybe/perfomance/fast_utf8_validation.pdf]]
https://github.com/felixge/fgtrace

```embed
title: 'N By N Skyscrapers'
image: 'https://www.codewars.com/packs/assets/og-image.7f5134fb.png'
description: 'Codewars is where developers achieve code mastery through challenge. Train on kata in the dojo and reach your highest potential.'
url: 'https://www.codewars.com/kata/5f7f023834659f0015581df1'
```

```embed
title: 'GitHub - janestreet/magic-trace: magic-trace collects and displays high-resolution traces of what a process is doing'
image: 'https://repository-images.githubusercontent.com/452433468/15ae30ee-7773-4fe0-a4ed-0d0d13751f70'
description: 'magic-trace collects and displays high-resolution traces of what a process is doing - GitHub - janestreet/magic-trace: magic-trace collects and displays high-resolution traces of what a process is ...'
url: 'https://github.com/janestreet/magic-trace'
```

```embed
title: 'How to Use AVX512 in Golang via C Compiler'
image: 'https://gorse.io/logo.png'
description: 'How to Use AVX512 in Golang via C Compiler AVX512 is the latest generation of SIMD instructions released by Intel, which can process 512 bits of data in one instruction cycle, equivalent to 16 single-precision floating point numbers or 8 double-precision floating point numbers. The training and infe…'
url: 'https://gorse.io/posts/avx512-in-golang.html'
```

https://highload.fun/
https://habr.com/ru/post/716292/
https://habr.com/ru/post/717716/
https://github.com/Highload-fun/platform/wiki
https://github.com/pkg/profile
https://dave.cheney.net/high-performance-json.html
https://www.youtube.com/watch?v=un-bZdyumog
https://philpearl.github.io/post/reader/
https://www.youtube.com/watch?v=un-bZdyumog
https://www.brendangregg.com/linuxperf.html
https://github.com/google/slowjam
https://matklad.github.io/2022/10/06/hard-mode-rust.html
https://github.com/mmcloughlin/avo
https://matklad.github.io/2022/10/06/hard-mode-rust.html
https://github.com/RazrFalcon/cargo-bloat
https://internals.rust-lang.org/t/is-custom-allocators-the-right-abstraction/13460
https://benchkram.de/blog/dev/profiling-go-programs
https://vk.com/vkteam?w=wall-147415323_15536 - Разработка высоконагруженного key-value хранилища https://github.com/recoilme/sniper
https://habr.com/ru/articles/733948/
https://avivcarmi.com/the-search-for-the-perfect-request-response-protocol/

https://github.com/dgryski/go-perfbook
https://github.com/aitsvet/debugcharts
https://golangconf.ru/2020/abstracts/7061
https://www.google.com/search?q=golang+rtdsc
https://www.agner.org/optimize/
https://rigtorp.se/ringbuffer/
https://lemire.me/blog/2023/02/07/bit-hacking-with-go-code/
https://totallygamerjet.hashnode.dev/the-smallest-go-binary-5kb
https://github.com/DataDog/go-profiler-notes
https://eblog.fly.dev/startfast.html
https://pkg.go.dev/modernc.org/memory
https://shane.ai/posts/threads-and-goroutines/
https://github.com/stackimpact/stackimpact-go
```go
reportMem := func() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	log.Println("total alloc:", m.TotalAlloc, "mallocs:", m.Mallocs)
}
reportMem()
defer reportMem()
f, err := os.Create("cpu.prof")
if err != nil {
	log.Fatal("could not create CPU profile: ", err)
}
defer f.Close()
if err := pprof.Lookup("allocs").WriteTo(f, 0); err != nil {
	// if err := pprof.StartCPUProfile(f); err != nil {
	log.Fatal("could not start CPU profile: ", err)
}
defer pprof.StopCPUProfile()
```
https://gitlab.com/gitlab-com/support/toolbox/strace-parser pretty strace output
https://github.com/samonzeweb/profilinggo
https://github.com/sebbbi/OffsetAllocator
https://medium.com/@lordmoma/6-tips-on-high-performance-go-advanced-go-topics-37b601fa329d
https://github.com/felixge/fgprof
https://github.com/zjc17/pprof-web
https://mazzo.li/posts/fast-pipes.html
https://nanxiao.gitbooks.io/perf-little-book/content/
https://github.com/dgraph-io/ristretto/tree/master/z
[flame graph online](https://flamegraph.com/)
https://jvns.ca/perf-zine.pdf