#someday_maybe #project

https://abhinavg.net/2022/05/15/hijack-testmain/

C - extract `main` to some `main_fun.c` as `main_fun` and `main.c` should be just

```c
#include "main_fun.c"
int main(int argc, char **argv) {
  return main_fun(argc, argv);
}
```

point is to have ability to define `Makefile` like so

```makefile
SOURCES = main_fun.c utils.c

build:
    gcc main.c $(SOURCES) -o main
test:
    gcc test.c $(SOURCES) -o test
    ./test
```

to build app executable and testable executable separately

`test` executable can test app by running test cases and can even use `main` defined as `main_fun`!

sample `test.c`

```c
#include "main_fun.c"
#include "utils.c"
#macro INIT_TESTS() = int err;
#macro TEST(test) = err = test(); if (err != SUCCES_CODE) return err;
#macro ASSERT_EQUALS(x, y, message) = if (x != y) {printf(message); return -1;}

const int SUCCES_CODE = 0;

int main_fun_tests() {
  return main_fun(2, ["1", "2"]);
}

int utils_tests() {
  ASSERT_EQUALS(add(1, 2), 3, "add(1, 2) != 3");
  ASSERT_EQUALS(parse_int("1"), 1, "parse_int(\"1\") != 1");
  return SUCCES_CODE;
}

int main() {
  INIT_TESTS();
  TEST(main_fun_tests);
  TEST(utils_tests);
  return SUCCES_CODE;
}
```
