#in

write https://elv.sh/ script to bruteforce receipts

```bash
curl 'https://alchemy.na4u.ru/alchemy/check-recipe' -X POST -H 'Content-Type: application/json' --data-raw '{"recipe":["water","fire"]}'
```
```json
{"uuid":"d0646071-e470-44e5-b299-2053ef67d5ce","id":"alcohol","name":"Спирт"}
```