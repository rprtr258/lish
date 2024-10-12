[GitHub - rprtr258/async-parallel-exercises: Some programs to train asynchronous programming](https://github.com/rprtr258/async-parallel-exercises)

сделать скрипт для бенчмарка sync vs async/parallel

сделать количество процессов/потоков разным, делать в бенчмарке график

- [ ]  java

![[data/static/old/someday_maybe/asynchronous,parallel_programming/java_async_factorial.png]]
![[data/static/old/someday_maybe/asynchronous,parallel_programming/java_factorial_thread.png]]
![[data/static/old/someday_maybe/asynchronous,parallel_programming/java_fibonacci_threads.png]]
![[data/static/old/someday_maybe/asynchronous,parallel_programming/java_sumrange_thread.png]]

- [ ]  clojure

[Mastering Concurrent Processes with core.async](https://www.braveclojure.com/core-async/)

- [ ]  kotlin

[Coroutines guide | Kotlin](https://kotlinlang.org/docs/coroutines-guide.html)

[KEEP/coroutines.md at master · Kotlin/KEEP](https://github.com/Kotlin/KEEP/blob/master/proposals/coroutines.md)

- [ ]  go
- [ ]  javascript
- [ ]  rust
- [ ]  C++
- [ ]  python

[Trio: a friendly Python library for async concurrency and I/O - Trio 0.19.0 documentation](https://trio.readthedocs.io/en/stable/)

[Client Quickstart - aiohttp 3.7.4.post0 documentation](https://docs.aiohttp.org/en/stable/client_quickstart.html#make-a-request)

- [ ]  ruby

[Preface | ØMQ - The Guide](https://zguide.zeromq.org/docs/preface/)

# Задачи

```jsx
Implement a simple version of Promise.all. This function should accept an array of promises and return an array of resolved values. If any of the promises are rejected, the function should catch them.

function Promise_all(promises)
    {return new Promise((resolve, reject) =>
      {/* Your code here.*/;});}

// Test code.
Promise_all([]).then(array =>
    {console.log("This should be []:", array);});
function soon(val)
    {return new Promise(resolve =>
        {setTimeout(() => resolve(val), Math.random() * 500);});}
Promise_all([soon(1), soon(2), soon(3)]).then(array =>
  {console.log("This should be [1, 2, 3]:", array);});
Promise_all([soon(1), Promise.reject("X"), soon(3)])
  .then(array => {console.log("We should not get here");})
  .catch(error => {if (error != "X") {console.log("Unexpected failure:", error);}});
```

[go-concurrency-exercises/0-limit-crawler at master · loong/go-concurrency-exercises](https://github.com/loong/go-concurrency-exercises/tree/master/0-limit-crawler)

[go-concurrency-exercises/1-producer-consumer at master · loong/go-concurrency-exercises](https://github.com/loong/go-concurrency-exercises/tree/master/1-producer-consumer)

sequential producer - consumer&producer - consumer&producer - consumer

[go-concurrency-exercises/2-race-in-cache at master · loong/go-concurrency-exercises](https://github.com/loong/go-concurrency-exercises/tree/master/2-race-in-cache#race-condition-in-caching-szenario)

[go-concurrency-exercises/3-limit-service-time at master · loong/go-concurrency-exercises](https://github.com/loong/go-concurrency-exercises/tree/master/3-limit-service-time)

[go-concurrency-exercises/4-graceful-sigint at master · loong/go-concurrency-exercises](https://github.com/loong/go-concurrency-exercises/tree/master/4-graceful-sigint)

[go-concurrency-exercises/5-session-cleaner at master · loong/go-concurrency-exercises](https://github.com/loong/go-concurrency-exercises/tree/master/5-session-cleaner)

[golang-vs-clojure-async.md](https://gist.github.com/xyproto/6584125)

[Clojure core.async and Go: A Code Comparison](https://blog.drewolson.org/clojure-go-comparison)

```python
import asyncio

queue = []

async def handle_echo(reader, writer):
    global queue
    queue.append((reader, writer))
    print(f"Accepted {writer.get_extra_info('peername')} connection")
    if len(queue) == 2:
        r1, w1 = queue.pop()
        r2, w2 = queue.pop()
        w1.write(await r2.read(100))
        w2.write(await r1.read(100))
        await w1.drain()
        await w2.drain()

        w1.close()
        w2.close()

async def main():
    server = await asyncio.start_server(handle_echo, '127.0.0.1', 8888)

    addr = server.sockets[0].getsockname()
    print(f'Serving on {addr}')

    async with server:
        await server.serve_forever()

asyncio.run(main())
```

```python
import asyncio

async def tcp_echo_client(message):
    reader, writer = await asyncio.open_connection('127.0.0.1', 8888)
    _, port = writer.get_extra_info("sockname")
    message = f"{message} from {port}"

    print(f'Send: {message!r}')
    writer.write(message.encode())
    await writer.drain()

    data = await reader.read(100)
    print(f'Received: {data.decode()!r}')

    print('Close the connection')
    writer.close()

asyncio.run(tcp_echo_client('Hello World!'))
```

[Андрей Викторович Столяров: сайт автора](http://www.stolyarov.info/books/gameserv)

twitch live followed channels site

file downloader with progress bar

solve some [adventofcode tasks](https://adventofcode.com/) with async / parallelism

tasks parallelism

parallel prime sieve

matrix multiplication

parallel map

parallel reduce

бинарным сдваиванием

разделением на батчи

sequence manipulation library

async twitch bot

recursive fibonacci

[Building H2O - LeetCode](https://leetcode.com/problems/building-h2o/)

[Print FooBar Alternately - LeetCode](https://leetcode.com/problems/print-foobar-alternately/)

[Print in Order - LeetCode](https://leetcode.com/problems/print-in-order/)

[](https://leetcode.com/problems/print-zero-even-odd/)

[Fizz Buzz Multithreaded - LeetCode](https://leetcode.com/problems/fizz-buzz-multithreaded/)

[The Dining Philosophers - LeetCode](https://leetcode.com/problems/the-dining-philosophers/)

fetching twitch messages, processing and making dataset for quote generator; quote generation model training

directed acyclic graph of task dependencies execution

# Тестовые данные

[Django REST framework](https://swapi.dev/api/)

```jsx
['https://www.gutenberg.org/files/764/764-0.txt',
'https://www.gutenberg.org/files/15/15-0.txt', 
'https://www.gutenberg.org/files/1661/1661-0.txt',
'https://www.gutenberg.org/files/84/84-0.txt',
'https://www.gutenberg.org/files/345/345-0.txt',
'https://www.gutenberg.org/files/768/768-0.txt',
'https://www.gutenberg.org/files/1342/1342-0.txt',
'https://www.gutenberg.org/files/11/11-0.txt',
'https://www.gutenberg.org/files/61262/61262-0.txt']
```

# Учебные материалы

[[data/static/old/someday_maybe/asynchronous,parallel_programming/The Art of Multiprocessor Programming.pdf]]

p. 120

[Параллельное программирование в JAVA напрактике.pdf](asynchrono%201e6c9/___JAVA_.pdf)

[[data/static/old/someday_maybe/asynchronous,parallel_programming/How to write parallel programs. A first course.pdf]]

[The Deadlock Empire](https://deadlockempire.github.io/#menu)

[Java Memory Model Pragmatics (transcript)](https://shipilev.net/blog/2014/jmm-pragmatics/)

[Close Encounters of The Java Memory Model Kind](https://shipilev.net/blog/2016/close-encounters-of-jmm-kind/)

[Introduction to Parallel Computing Tutorial](https://hpc.llnl.gov/training/tutorials/introduction-parallel-computing-tutorial)

[njs blog](https://vorpus.org/blog/)

[Lock-free структуры данных. Concurrent map: разминка](https://habr.com/en/post/250383/)

[Nathaniel J Smith - Python Concurrency for Mere Mortals - Pyninsula #10](https://www.youtube.com/watch?v=i-R704I8ySE)

[Communicating Sequential Processes.pdf](asynchrono%201e6c9/Communicating_Sequential_Processes.pdf)

[Communicating Sequential Processes](https://swannodette.github.io/2013/07/12/communicating-sequential-processes/)

[distributed-development-stack.pdf](https://drive.google.com/file/d/18PLe4KicpLTU5GzHnkXhiLP3d3fIwS0K/view?usp=drivesdk)

[Designing Data-Intensive Applications. The Big Ideas Behind Reliable, Scalable, and Maintainable Systems ( PDFDrive.com ).pdf](https://drive.google.com/file/d/1UELHvq6Vc--PqUQL7BuCtLHcFZ4KkIVJ/view?usp=drivesdk)

[Распределенные системы. Паттерны проектирования.pdf](asynchrono%201e6c9/_.__.pdf)

[](https://github.com/rprtr258/fs/blob/main/Распределенные%20системы.%20Принципы%20и%20парадигмы.pdf)
[[data/static/old/someday_maybe/asynchronous,parallel_programming/Распределенные системы.pdf]]
[[data/static/old/someday_maybe/asynchronous,parallel_programming/progintro_dmkv2.pdf]]
[[data/static/old/someday_maybe/asynchronous,parallel_programming/Communication and Mobile Systems.pdf]]
[[data/static/old/someday_maybe/asynchronous,parallel_programming/Communication and Concurrency.pdf]]
https://github.com/luk4z7/go-concurrency-guide
https://www.youtube.com/playlist?list=PL4_hYwCyhAva37lNnoMuBcKRELso5nvBm