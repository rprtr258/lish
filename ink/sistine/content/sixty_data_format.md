#someday_maybe #project

https://github.com/rprtr258/sixty

[[data/tree-diff]]

https://github.com/nin-jin/slides/tree/master/tree

- [ ]  Посмотреть форматы сериализации данных
    [SXML - Wikipedia](https://en.wikipedia.org/wiki/SXML) 
    [YANG - Wikipedia](https://en.wikipedia.org/wiki/YANG)
    [JSON](https://www.json.org/json-en.html) 
    [BSON - Википедия](https://ru.wikipedia.org/wiki/BSON)
    MessagePack, BSON, BJSON, UBJSON, BISON, Smile, WBXML, Fast Infoset
    [Tree - убийца JSON, XML, YAML и иже с ними](https://habr.com/en/post/248147/) 
    https://habr.com/en/company/vk/blog/504952/
- [ ]  Сравнить существующие форматы сериализации
    [Непричёсанные мысли по поводу формата сохранения: теория](https://habr.com/en/post/343078/) 
    [Protobuffers - это неправильно](https://habr.com/en/post/427265/)
- [ ]  (Если нужно) Написать программу и библиотеку на CommonLisp для конвертации s-expr в разные форматы сериализации и обратно
    [ ] json
    [ ] yaml
    [ ] ini
    [ ] toml
    [ ] css
    [ ] etc
    [https://serde.rs/](https://serde.rs/)
- [ ]  Протестировать библиотеку
    сгенерировать случайный `<json>`, преобразовать в `s-expr` и обратно, сравнить 
    сгенерировать случайный `<json>`, преобразовать `json`→`s-expr`→`yaml`→`s-expr`→`xml`→`s-expr`→`json` и сравнить
- [ ]  сделать селекторы (x-path/css like) для поиска/преобразования
https://github.com/janestreet/sexp
- recursive datatypes?, check that from any type at least one primitive type is accessible
- [ ] query/transform facility
    https://www.w3schools.com/cssref/css_selectors.asp
    https://webcache.googleusercontent.com/search?q=cache:AYUf1LG3U5cJ:https://www.mongodb.com/docs/manual/reference/operator/query/+&cd=1&hl=ru&ct=clnk&gl=ru
    https://nbviewer.org/github/RumbleDB/rumble/blob/master/RumbleSandbox.ipynb
    https://stedolan.github.io/jq/tutorial/
    https://www.jsonquerytool.com/
      https://jmespath.org/tutorial.html
      https://kubernetes.io/ru/docs/reference/kubectl/jsonpath/
      https://jsonpath.com/
      https://github.com/dfilatov/jspath
      https://datatracker.ietf.org/doc/html/rfc6901
      https://jsonata.org/
      https://github.com/dchester/jsonpath
      https://jsonpath-plus.github.io/JSONPath/docs/ts/

https://github.com/vlm/asn1c

```
html {
  head {
    title "Welcome page"
  }
  body {
    para "Hello, world"
  }
}
```

# map reduce
https://github.com/alecgorge/json-map-reduce-tool

https://ru.wikipedia.org/wiki/Re2c

see useful logical types
    https://developers.google.com/protocol-buffers/docs/reference/google.protobuf
any = schema * value/schema

https://en.wikipedia.org/wiki/JSON_streaming
https://ru.wikipedia.org/wiki/Obj
https://wiki.c2.com/?XmlIsaPoorCopyOfEssExpressions
https://github.com/hamba/avro/

- benches (just for fun and perf/memory metrics)
- streaming encoder/decoder, with size of element for lists (if it is not known)

clis
- encode/decode json
- encode/decode binary
- encode/decode text format
- generate go code with types and encoder/decoder
- formatting cli

### metadata (TODO)
can be used whenever type `A` is used
```
A `key:value,flag,key:value`
```
e.g.
```
rating = uint `<=:10`
position = (int, int, int) `name:pos` // pos is the name of ints product
positions = []position `soa` // `soa` is about list of products
```
auto define `name` value for each aliased type?
## alias/typedef
identifier is
```
[A-Za-z_:][A-Za-z0-9_-]*
```
like so (`std:` prefix is for predefined (avro calles them logical) types, maybe just `:type`)
```
std:u8 = [8]bit // really a byte
std:bool = bit // or {TRUE, FALSE}
std:f64 = (
  bit     `name:sign`,
  [11]bit `name:exponent`,
  [52]bit `name:mantissa`, // TODO: just [64]bit
)
month = [0..11]
week_day = {Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday}
// week_day might be described as [0..6] but then it is up to application to deduce what is first day of the week
date = (
  uint `name:day`,
  month,
  year,
)
second = [0..59]
minute = [0..59]
hour = [0..23]
time = (second, minute, hour)
// really it is just int - number of seconds since Unix epoch
datetime = (date, time)
timezone = (minute, hour)
positions = [](int, int, int) `soa` // `soa` is about list of products
```

https://cuelang.org/
https://en.wikipedia.org/wiki/Interface_description_language