#someday_maybe #project

when developing locally it might be useful to skip (select) some tests to run only tests that might be affected by last changes. OFC that should be done only locally, cuz CI will have to run all tests anyway.

Skippable tests should be somehow specifiable at test-running time, not in code like it is done here

[Intermediate Testing in Golang](https://routley.io/posts/intermediate-go-testing/#skipping-tests)

(skipping groups of / specific tests?, using build tags?)

[The Rust Programming Language](https://doc.rust-lang.org/book/ch11-02-running-tests.html#controlling-how-tests-are-run)