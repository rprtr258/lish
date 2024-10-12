`GET /shop`

возвращает список предметов в магазине, e.g.

```json
[
    {
        "id": 1,
        "name": "Forest bee",
        "obj": {
            "itemtype": "bee",
            "lifetime": 100,
            "type": "forest"
        },
        "price": 10
    },
    ...
]
```

`GET /money`

возвращает баланс, e.g.

```json
{
    "balance": 100
}
```

`POST /buy`

купить предмет, в запросе указывается $id$ предмета для покупки, после покупки вычитаются деньги из баланса и предмет добавляется в инвентарь

`GET /inv`

возвращает список предметов в инвентаре e.g.

```json
[
    {"itemtype": "bee", "lifetime": 100, "type": "forest"},
    {"itemtype": "comb", "type": "honey_comb", "count": 2},
    ...
]
```