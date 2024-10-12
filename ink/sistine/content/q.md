https://github.com/rprtr258/q

[machine learning](machine%20learning%20%26%20data%20science)
https://jqlang.github.io/jq/manual/
https://zed.brimdata.io https://zui.brimdata.io/

main idea is: many models are represented as graphs: vertices are entities, edges are relations between them, so as a graph can be modelled:
- all SQL databases
- package managers with package dependencies relations and other metadatas (package repos, authors, etc)
- resources in k8s
- command outputs like `ls`, `ps`, etc
there are plenty of scenarios to query these things in a similar way:
- find all groups which user is member of
- find dependencies of particular package
- find resources used by another k8s resource
- find processes cmd using file or port

example schema:
```graphql
type User {
    name String
    region Region
}

type Region {
    id Uint64
    name String
}
```
example queries:
```php
user[name,region.name]
user[name,region.name,name="Petya"?]
user[name,role.permission[name="Delete DB"?]]
user[name,region.name,name=/Petya/?]
user[name,role.permission[*, name=/Delete DB/?]]
```

in `jq`
```jq
$users | map([.name,.region.name])
$users | map(select(.name=="Petya")) | map([.name,.region.name])
$users | map(select($roles[.id].permission | contains("Delete DB"))) | map(.name)
$users | map(select(.name|match(".*Petya.*"))) | map([.name,.region.name])
$users | map(select($roles[.id].permission | any(match(".*Delete DB.*")))) | map(.name)
```

unxpectedly hard query
```sql
entities
    film
        id       ID
        title    string
        actors   [ID]
        director ID
    person
        id   ID
        name string

find man with film he acted in and other film he directed
logic is following:
SELECT
    m1.title
    m2.title
    p.name
FROM
    m1@movie
    m2@movie
    p@profile
WHERE
    m1.id!=m2.id
    m1.actors.*=p.id
    m2.director=p.id

jq:
(.film | map(.id as $film_id | .actors | map({($film_id): .})) | add | add) as $actoring |
(.film | map({(.director): .id}) as $directors |
$actoring | keys | map(. as $actor_id | select($directors | contains($actor_id)))
```

so it makes sense to make succint query language taking graph structure in account
- sql is way verbose and has no explicit ways to use relations(haha) between entities
- graphql is verbose, no way to dinamically filter data
- other graph query languages focuses on many entities involved in graph structure, so they are solving completely different tasks: find paths, find vertex on path between such that, etc

see nushell
https://miro.medium.com/v2/resize:fit:679/1*yjJs6XwNTu-qeliPdYbYYg.gif
https://home.adelphi.edu/~siegfried/cs443/443l9.pdf
https://habr.com/ru/articles/759342/
https://kroki.io/examples.html#excalidraw https://kroki.io/examples.html#erd graphics
![[data/static/old/projects/q/schema-sketch.png]]
![[data/static/old/projects/q/pb.png]]
https://github.com/alecthomas/participle
https://sq.io/
file-tree like tui

[Вулканический поросенок, или SQL своими руками](https://habr.com/ru/companies/badoo/articles/461699/)

additional literals
@now - the current datetime as string
@second - @now second number (0-59)
@minute - @now minute number (0-59)
@hour - @now hour number (0-23)
@weekday - @now weekday number (0-6)
@day - @now day number
@month - @now month number
@year - @now year number
@todayStart - beginning of the current day as datetime string
@todayEnd - end of the current day as datetime string
@monthStart - beginning of the current month as datetime string
@monthEnd - end of the current month as datetime string
@yearStart - beginning of the current year as datetime string
@yearEnd - end of the current year as datetime string
operators
length() - number of connected entities
- `&&`, `||` - and, or
- `=` Equal
- `!=` NOT equal
- `>` Greater than
- `>=` Greater than or equal
- `<` Less than
- `<=` Less than or equal
- `~` Like/Contains (if not specified auto wraps the right string OPERAND in a "%" for wildcard match)
- `!~` NOT Like/Contains (if not specified auto wraps the right string OPERAND in a "%" for wildcard match)
- `?=` _Any/At least one of_ Equal
- `?!=` _Any/At least one of_ NOT equal
- `?>` _Any/At least one of_ Greater than
- `?>=` _Any/At least one of_ Greater than or equal
- `?<` _Any/At least one of_ Less than
- `?<=` _Any/At least one of_ Less than or equal
- `?~` _Any/At least one of_ Like/Contains (if not specified auto wraps the right string OPERAND in a "%" for wildcard match)
- `?!~` _Any/At least one of_ NOT Like/Contains (if not specified auto wraps the right string OPERAND in a "%" for wildcard match)

To group and combine several expressions you could use brackets `(...)`, `&&` (AND) and `||` (OR) tokens.

sorting
`<field>` - ascending
`-<field>` - descending

https://www.aleksandra.codes/comments-db-model
https://github.com/getnokori/api
https://github.com/iximiuz/kexp
strings querying https://lucene.apache.org/core/2_9_4/queryparsersyntax.html
https://mermaid.js.org/syntax/entityRelationshipDiagram.html
sample queries https://www.orientdb.com/docs/last/gettingstarted/demodb/queries/index.html
https://dgraph.io/docs/query-language/functions/
https://docs.jargon.sh/#/pages/language
https://drawsql.app/ visualization of relational database schema
https://steampipe.io/docs/develop/writing-plugins https://hub.steampipe.io/plugins?categories=software+development
http://number-none.com/product/Lerp,%20Part%203/index.html
[F1 Query: Declarative Querying at Scale](https://storage.googleapis.com/gweb-research2023-media/pubtools/pdf/fa380016eccb33ac5e92c84f7b5eec136e73d3f1.pdf)
[harelba/q: q - Run SQL directly on delimited files and multi-file sqlite databases](https://github.com/harelba/q)
http://docs.jsonata.org/programming
https://jmespath.org/tutorial.html
https://prql-lang.org/book/overview.html