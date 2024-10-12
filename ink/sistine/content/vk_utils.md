#someday_maybe #project

# VK UTILS

## План

- [ ] rewrite to `execute`
    [](https://dev.vk.com/method/execute)
    [Урок 41 - Бонус - Программирование на VKScript и метод Execute](https://vk.com/@autopilot_school-execute)
    [Анализ языка VKScript: JavaScript, ты ли это?](https://habr.com/en/post/464099/)
- [ ] todos, fixes
    [Коды ошибок](https://dev.vk.com/reference/errors)
    [API Requests | Developers | VK](https://vk.com/dev/api_requests)
    [Streaming API и ограничения API для поиска | Developers | VK](https://vk.com/dev/data_limits)
- posts search
    https://vk.com/app3876642
    https://vk.com/wall-2158488_651604
- [ ] dump groups/profiles posts into database (json format)

tests

```bash
curl 'http://localhost:8000/reposts' --form 'postUrl=https://vk.com/wall-149859311_975'
```
```
[](user_id=uint, repost_id=int)
(478324650, 333)
(269848835, 138)
(564879147, 73)
(45144307, "50848)
```

```bash
curl 'http://localhost:8000/reposts' --form 'postUrl=https://vk.com/wall-196725903_10304'
```
```
[](user_id=uint, repost_id=int)
(439667637, 709)
(226164478, 20837)
(28415286, 1591)
(260326982, 10100)
(262319680, 796)
(109595520, 1433)
(287419935, 6682)
(496862771, 2638)
(293566673, 582)
(144954480, 2970)
(144217682, 8710)
(170777868, 4855)
(555494732, 5505)
(364505590, 1572)
(365604588, 1348)
(626285320, 220)
(424251029, 1944)
```
