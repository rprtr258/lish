#project [files](http://192.168.1.139:5000/old/projects/backend/)

https://roadmap.sh/backend
https://github.com/Alikhll/golang-developer-roadmap
https://github.com/bzick/oh-my-backend
https://github.com/ByteByteGoHq/system-design-101

[Правильные ответы по криптографии](https://habr.com/ru/company/globalsign/blog/353576/)
[Делаем современное веб-приложение с нуля](https://habr.com/en/post/444446/)
[codgen](https://habr.com/ru/company/sbermarket/blog/676486/)
https://go-zero.dev/

[golang/dependency analysis](https://github.com/loov/goda)
[Visualize call graph of a Go program using Graphviz](https://github.com/ondrajz/go-callvis)

https://www.infoq.com/presentations/event-tracing-monitoring/
[[projects/backend/go-with-domain.pdf]]
[how cdn works article series](https://uploadcare.com/blog/author/igor-adamenko/)
[minimal base docker image](https://github.com/GoogleContainerTools/distroless)
[[reference/backend/_architecture/Large-scale cluster management at Google with Borg.pdf]]
[[reference/backend/_architecture/Создание микросервисов.pdf]]
[[reference/backend/_architecture/proektirovanie_protsessa_proektirovaniya.pdf]]
[[projects/backend/Облачный Go. Создание надежных служб в ненадежных окружениях [2022] Мэтью А. Титмус.pdf]]
[golang/tools](https://github.com/nikolaydubina/go-recipes)
https://deepflow.yunshan.net/docs/ automatic instrumentation with eBPF
[localhost with SSL cert](https://github.com/Upinel/localhost.direct)
[Service Mesh на стероидах, часть 2: Zero Deployment Downtime в корпоративных приложениях](https://habr.com/ru/companies/oleg-bunin/articles/708110/)
https://github.com/bouk/monkey
https://deep.foundation/
https://dev.to/hixdev/software-project-checklist-4chb
https://trunkbaseddevelopment.com/
https://github.com/contribsys/faktory Language-agnostic persistent background job server
https://github.com/cheatsnake/backend-cheats/blob/master/README_RUS.md
https://monorepo.tools/
https://github.com/aitsvet/debugcharts
https://jbrandhorst.com/post/grpc-errors/
https://go-kratos.dev/en/
https://www.youtube.com/watch?v=rCJvW2xgnk0 rest api in go
https://www.youtube.com/watch?v=tqQr2tNpJrA tg bot in go
https://habr.com/ru/articles/460535/ minimal docker image for golang
https://github.com/ulid/spec https://habr.com/ru/articles/658855/ https://github.com/jetpack-io/typeid uuid vs ulid vs typeid
https://zalopay-oss.github.io/go-advanced/ch5-distributed-system/ch5-03-delay-job.html
https://evilmartians.com/chronicles/speeding-up-go-modules-for-docker-and-ci
https://github.com/jetpack-io/devbox
https://www.inkandswitch.com/local-first/
https://squeaky.ai/blog/development/why-we-dont-use-a-staging-environment https://news.ycombinator.com/item?id=30899363 https://refactoring.fm/p/do-you-need-staging
[Best distributed task scheduling framework — Openjob 1.0.7 released](https://habr.com/ru/articles/760108/)
https://hub.docker.com/r/clickhouse/clickhouse-keeper - better zookeeper
https://quii.dev/How_to_go_fast
[golang/charts](https://github.com/go-echarts/go-echarts)
https://github.com/arl/statsviz simple builtin service dashboard for golang
[github graphite/client](https://graphite.dev/)
https://www.amazon.com/Designing-Distributed-Control-Systems-Language/dp/1118694155
[manage tools using go.mod](https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module)
https://github.com/nikolaydubina/go-cover-treemap
https://go-architect.github.io/docs/install

## tools
https://github.com/boyter/cs fuzzy code searcher
https://github.com/cue-lang/cue configuration and validation language
https://github.com/boyter/dcd/ duplicate code finder
https://go.dev/blog/deadcode dead code finder
https://github.com/go-delve/delve/tree/master/Documentation/installation debugger
https://github.com/ktr0731/evans grpc cli client
https://github.com/fullstorydev/grpcurl grpc cli client
https://github.com/antonmedv/fx terminal json viewer
https://github.com/ondrajz/go-callvis function calling diagram 
https://github.com/Gelio/go-global-update update binaries installed via `go install`
https://github.com/go-critic/go-critic opinionated linter
https://github.com/mvdan/gofumpt stricter gofmt
https://golangci-lint.run/usage/install/#binaries linters
https://github.com/gotestyourself/gotestsum appealing test runner
https://github.com/charmbracelet/gum bash prompts
https://github.com/google/go-jsonnet imperative templating language
https://github.com/birdayz/kaf kafka cli
https://github.com/jesseduffield/lazygit git tui
https://github.com/alajmo/mani checkout multiple git repos
https://confluence.o3.ru/pages/viewpage.action?pageId=382366442
https://github.com/noborus/ov pager
https://github.com/bep/punused find public unused symbols
https://github.com/cespare/reflex run command on file change
https://github.com/rprtr258/rwenv run w/ env
https://github.com/rprtr258/xxd colourful xxd
https://github.com/go-delve/delve/blob/master/Documentation/cli/README.md debugger

### kubernetes tools
`curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"`
https://github.com/ahmetb/kubectx
```
kubectl krew install ctx
kubectl krew install ns
```

## caching
https://github.com/dragonflydb/dragonfly
http://www.pelikan.io/

## [docs](/note/how%20to%20write%20documentation)

## k8s/remote/docker for development
https://mirrord.dev/
https://www.telepresence.io https://habr.com/ru/companies/flant/articles/446788/
https://www.google.com/search?q=k8s+for+developer
https://skaffold.dev/
https://www.devspace.sh/
https://tilt.dev/
https://kubernetes.io/blog/2018/05/01/developing-on-kubernetes/
https://threedots.tech/post/go-docker-dev-environment-with-go-modules-and-live-code-reloading/
https://devpod.sh/
https://rnemet.dev/posts/docker/new_stuff_init_watch/
[VSСode. Как настроить окружение для разработки в Docker на удаленном сервере через SSH](https://habr.com/ru/articles/734062/)
[serverside vscode](https://hub.docker.com/r/linuxserver/code-server)
https://unzip.dev/p/0x015-dev-containers

## file storage
https://blog.min.io/go-based-amazon-s3-cli/
[distributed file system](https://github.com/seaweedfs/seaweedfs)
[Зачем и как хранить объекты на примере MinIO](https://habr.com/ru/company/ozontech/blog/586024/)
https://cloud.yandex.ru/services/storage#calculator
https://bazil.org
https://gist.github.com/GabLeRoux/3b2168a88ceca43ac524e318839dd798
[lakeFS - Git-like capabilities for your object storage](https://github.com/treeverse/lakeFS)

## multiple git repos
https://github.com/isacikgoz/gitbatch
https://github.com/alajmo/mani

## services examples
https://github.com/thangchung/go-coffeeshop
https://github.com/ThreeDotsLabs/wild-workouts-go-ddd-example


## message brokers
[fast message broker](https://github.com/nats-io/nats.go)
https://github.com/emqx/nanomq
https://zeromq.org/get-started/
https://bloomberg.github.io/blazingmq/
https://github.com/birdayz/kaf
https://github.com/nsqio/nsq
[Как правильно выбирать очередь](https://www.youtube.com/watch?v=hEC8CX8Drac)

## databases
https://github.com/mgramin/awesome-db-tools
https://vitess.io/
https://dbsensei.com/
[Методы выявления ошибок в SQL приложении](http://sql-error.microbecal.com/)
https://github.com/d-tsuji/awesome-go-orms
![[reference/backend/amazon_dbs.webp]]
![[projects/backend/Pasted image 20230425183557.png]]
https://db-engines.com/en/ranking/document+store
http://go-database-sql.org/accessing.html
https://rethinkdb.com/faq
https://surrealdb.com/
https://docs.keydb.dev/
https://cloud.docs.scylladb.com/stable/scylladb-basics/
[git like db](https://www.dolthub.com/)
[git like db in ocaml](https://irmin.org/)
[distributed sqlite](https://github.com/rqlite/rqlite)
[in-memory(?) db](https://www.tarantool.io/ru/)
https://planetscale.com/blog/what-are-the-disadvantages-of-database-indexes
https://entgo.io/blog/2021/10/11/generating-ent-schemas-from-existing-sql-databases
https://planetscale.com/blog/how-read-mysql-explains
https://www.tigrisdata.com/ document db
https://github.com/paypal/junodb
https://github.com/rosedblabs/rosedb fast kv database
https://github.com/syndtr/goleveldb kv database, leveldb implemented in golang
https://github.com/danielealbano/cachegrand fast kv storage
[gitlab for db](https://www.bytebase.com/)
[sql diagrams](https://dbdiagram.io/home)
[post sql database](https://www.edgedb.com/)
https://antonz.org/sql-window-functions-book/
https://architecturenotes.co/database-sharding-explained/
https://github.com/kelindar/column in-memory columnar database
https://github.com/doug-martin/goqu query builder
https://justinjaffray.com/joins-13-ways/
http://dbcli.com/ cli clients
https://github.com/xo/usql
https://rethinkdb.com/
https://sqltools.beekeeperstudio.io/build
[Vector database built for scalable similarity search](https://milvus.io/)
postgres query plan visualizer with `EXPLAIN (ANALYZE, COSTS, VERBOSE, BUFFERS, FORMAT JSON)` https://explain.dalibo.com

### full text searc engines
https://docs.paradedb.com/quickstart

### clients
beekeeper studio
dbeaver
https://www.jetbrains.com/ru-ru/datagrip/features/
https://slashbase.com/
[redis gui](https://github.com/diego3g/rocketredis)
https://github.com/RedisInsight/RedisInsight

### orm-like frameworks
https://play.sqlc.dev/
https://entgo.io/docs/getting-started
https://github.com/go-jet/jet
https://bun.uptrace.dev/guide/#ent
https://github.com/volatiletech/sqlboiler

### migrations
https://atlasgo.io/
[migrations tool](https://github.com/pressly/goose)
https://planetscale.com/blog/backward-compatible-databases-changes


## authorization
[fresh token client implementation in go (repeatable Future)](https://appliedgo.net/refresh/)
https://casbin.org/
https://www.openpolicyagent.org/
https://checkmarx.gitbooks.io/go-scp/content/session-management/
https://dev.to/egeaytin/why-google-build-zanzibar--3kp5
[access control system](https://github.com/Permify/permify)
https://authzed.com/docs
https://github.com/openfga/openfga
https://github.com/warrant-dev/warrant
keycloak
[aserto-dev/topaz: Cloud-native authorization for modern applications and APIs](https://github.com/aserto-dev/topaz)
[cerbos/cerbos: Cerbos is the open core, language-agnostic, scalable authorization solution that makes user permissions and authorization simple to implement and manage by writing context-aware access control policies for your application resources](https://github.com/cerbos/cerbos)
[RBAC vs ReBAC: When to use them](https://dev.to/egeaytin/rbac-vs-rebac-when-to-use-them-47c4)

## authentication
https://casdoor.org/
https://github.com/go-pkgz/auth
https://testdriven.io/blog/web-authentication-methods/
https://dexidp.io/docs/getting-started/
https://openfga.dev/
https://habr.com/ru/articles/779534/
https://www.authelia.com/
https://www.youtube.com/watch?v=BQ0lNhIFQBk
[markbates/goth: Package goth provides a simple, clean, and idiomatic way to write authentication packages for Go web applications](https://github.com/markbates/goth)
[zitadel/zitadel: ZITADEL - The best of Auth0 and Keycloak combined](https://github.com/zitadel/zitadel)
[Creating an OAuth2 Client in Golang (With Full Examples)](https://www.sohamkamani.com/golang/oauth/)
[eko/authz: Authorization backend that comes with a UI for RBAC and ABAC permissions](https://github.com/eko/authz)

### jwt
https://infosecwriteups.com/json-web-tokens-409297c260a0
https://gist.github.com/zmts/802dc9c3510d79fd40f9dc38a12bccfc
https://jwt.io
https://stormpath.com/blog/where-to-store-your-jwts-cookies-vs-html5-web-storage
https://dev.to/siddheshk02/jwt-authentication-in-go-5dp7
https://roadmap.sh/guides/jwt-authentication
[How to test for JWT attacks](https://systemweakness.com/how-to-test-for-jwt-attacks-513da89abe94)

### oauth2
![](https://roadmap.sh/guides/oauth.png)
https://blog.oauth.io/choose-oauth2-flow-grant-types-for-app/
![](https://miro.medium.com/v2/resize:fit:4800/format:webp/1*_hL37NU36uUPePlkUjS3KQ.png)
https://dev.twitch.tv/docs/authentication/getting-tokens-oauth/


## микросервисы
https://frpc.io/introduction
[dumb pipes and smart endpoints - Google](https://www.google.ru/search?q=dumb+pipes+and+smart+endpoints)
[Microservices Integration Using Service Choreography](https://medium.com/@dreweaster/the-art-of-microservices-integration-using-service-choreography-69a4bbbf81c5)
https://micro.dev/
https://eax.me/micro-service-architecture/
https://microservices.io
https://www.google.com/search?q=go+kit
https://eightify.app/ru/summary/programming/mastering-clean-architecture-with-vertical-slices


## system design
https://www.youtube.com/playlist?list=PLrw6a1wE39_tb2fErI4-WkMbsvGQk9_UB#distributed MIT distributed systems lectures
https://pdos.csail.mit.edu/6.824/
https://roadmap.sh/design-system
https://github.com/karanpratapsingh/system-design
https://kps.hashnode.dev/system-design-the-complete-course
https://dev.to/karanpratapsingh/series/19332
[http://nealford.com/katas/list.html](http://nealford.com/katas/list.html)
[System Design Cheatsheet](https://gist.github.com/vasanthk/485d1c25737e8e72759f) 
[System Design Interviews: A step by step guide](https://www.educative.io/courses/grokking-the-system-design-interview/B8nMkqBWONo)
[Scalable System Design Patterns](http://horicky.blogspot.com/2010/10/scalable-system-design-patterns.html) 
[System Design - LeetCode Discuss](https://leetcode.com/discuss/interview-question/system-design/?currentPage=1&orderBy=recent_activity&query=)
[Systems Analysis and Design (SAD) Tutorial](https://www.w3computing.com/systemsanalysis/) 
[Systems Analysis](https://web.archive.org/web/20110722022042/http://web.simmons.edu/~benoit/LIS486/SystemsAnalysis.html)
[Бизнес и системный анализ для архитекторов](https://www.youtube.com/playlist?list=PLrCZzMib1e9ryE2N4xMHfbM0fWM8veVQ7) 
[system-design-primer](https://github.com/donnemartin/system-design-primer/blob/master/README.md)
[Как проходят архитектурные секции собеседования в Яндексе: практика дизайна распределённых систем](https://habr.com/en/company/yandex/blog/564132/) 
[Темные века разработки программного обеспечения](https://habr.com/en/company/cian/blog/569940/)
[Системный архитектор. Кто этот человек?](https://habr.com/en/post/593559/) 
[Программный архитектор. Кто этот человек?](https://habr.com/en/post/590883/)
https://12factor.net/
[The Architecture of Open Source Applications](https://aosabook.org/en/index.html)
https://book.mixu.net/distsys/single-page.html
[Архитектура приложений с открытым исходным кодом](http://rus-linux.net/MyLDP/BOOKS/Architecture-Open-Source-Applications/index.html)
https://faun.pub/top-30-system-design-interview-questions-and-problems-for-programmers-417e89eadd67
https://dev.to/gbengelebs/netflix-system-design-backend-architecture-10i3
https://www.educative.io/blog/complete-guide-to-system-design
https://github.com/InterviewReady/system-design-resources
https://blog.devgenius.io/system-design-blueprint-the-ultimate-guide-e27b914bf8f1
https://github.com/theanalyst/awesome-distributed-systems
https://www.distributedsystemscourse.com
https://www.freecodecamp.org/news/design-patterns-for-distributed-systems/
https://martinfowler.com/articles/patterns-of-distributed-systems/
[Паттерн Outbox: как не растерять сообщения в микросервисной архитектуре](https://habr.com/ru/companies/lamoda/articles/678932/)
https://dev.to/siy/nanoservices-or-alternative-to-monoliths-and-microservices-12bb
https://www.codereliant.io/the-art-of-building-fault-tolerant-software-systems/
https://github.com/binhnguyennus/awesome-scalability
https://www.youtube.com/watch?v=ppvuFdaYv3k
[distributed transactions strategies](https://developers.redhat.com/articles/2021/09/21/distributed-transaction-patterns-microservices-compared#how_to_choose_a_distributed_transactions_strategy)
https://github.com/bregman-arie/system-design-notebook
[Architecture Notes — System Design & Software Architectures Explained](https://architecturenotes.co/)
[[reference/backend/_architecture/Распределенные системы.pdf]]
[[reference/backend/_architecture/Распределенные системы2.pdf]]
[[reference/backend/_architecture/System Design Interviews.pdf]]
[[reference/backend/_architecture/distributed systems practitioners.pdf]]
[[reference/backend/_architecture/Идеальная архитектура.pdf]]
[[reference/backend/_architecture/Чистая архитектура.pdf]]
[[reference/backend/_architecture/Event_Streams_in_Action_Real_time.pdf]]
[[reference/backend/_architecture/97-things-cloud-engineers-know.pdf]]
[[reference/backend/_architecture/Распределенные системы. Принципы и парадигмы.pdf]]
[[reference/backend/_architecture/Stream_Processing_with_Apache_Spark.pdf]]
[5 Common Server Setups For Your Web Application](https://www.digitalocean.com/community/tutorials/5-common-server-setups-for-your-web-application)
https://temporal.io/
https://serviceweaver.dev/docs.html
https://github.com/inngest/inngesthttps://cadenceworkflow.io/docs/get-started/golang-hello-world/
[microservices patterns](https://microservices.io/patterns/microservice-chassis.html)
[System Design 101](https://habr.com/ru/articles/770564/)
https://github.com/mehdihadeli/awesome-software-architect

#### saga
https://habr.com/ru/companies/ozontech/articles/590709/
https://habr.com/ru/companies/oleg-bunin/articles/588488/
https://www.cs.cornell.edu/andru/cs711/2002fa/reading/sagas.pdf

## ci
https://stepci.com/
https://github.com/schemathesis/schemathesis
ghactions testing locally https://github.com/nektos/act
https://github.com/ory/dockertest https://eax.me/golang-dockertest/
https://golang.testcontainers.org/quickstart/
https://github.com/mxschmitt/action-tmate
https://agola.io/

## api design
https://github.com/ogen-go/ogen OpenAPI v3 code generator for go
https://abhinavg.net/2022/12/06/designing-go-libraries/
https://habr.com/ru/company/redmadrobot/blog/719222/
https://google.aip.dev/general
https://cloud.google.com/apis/design
![[projects/backend/Pasted image 20230404201636.png]]
library api design https://www.youtube.com/watch?v=ZQ5_u8Lgvyk = https://www.chrishecker.com/API_Design
[Стажёр Вася и его истории об идемпотентности API](https://habr.com/en/company/yandex/blog/442762/)
https://github.com/deepmap/oapi-codegen
https://habr.com/ru/companies/piter/articles/729874/
https://twirl.github.io/The-API-Book/API.ru.html
https://jargon.sh/
https://www.youtube.com/watch?v=KSBed4yyoDM
https://habr.com/ru/companies/piter/articles/472522/
https://github.com/brandur/heroku-http-api-design
https://appliedgo.net/sqlasapi/

## event drive architecture
https://www.asyncapi.com/ event driven application specification and tools
https://docs.dapr.io/concepts/overview/ [Микросервисы на основе событий с Dapr](https://habr.com/ru/company/otus/blog/706186/)
https://encore.dev/blog/event-driven-architecture https://encore.dev/blog/eda-business-case https://encore.dev/blog/building-for-failure https://encore.dev/blog/long-term-ownership

## optimization
[high-performance networking](https://hpbn.co/)
https://istlsfastyet.com/
https://www.marcobehler.com/guides/load-testing
[fast? go http framework](https://github.com/cloudwego/hertz)
[k6](https://k6.io/) [bombardier](https://github.com/codesenberg/bombardier) [vegeta](https://github.com/tsenart/vegeta) [wrk](https://github.com/wg/wrk)
[ab replacement](https://github.com/rakyll/hey)
[Обзор тестирования производительности](https://habr.com/ru/articles/735376/)
redis proto might be faster than grpc https://avivcarmi.com/the-search-for-th
https://github.com/alipay/fury fast proto format

## protocol formats
redis proto might be faster than grpc https://avivcarmi.com/the-search-for-the-perfect-request-response-protocol/
[zero rpc](https://capnproto.org/)
https://github.com/alipay/fury fast proto format
https://storj.github.io/drpc/


### grpc
https://github.com/bufbuild/buf
https://bob.build/blog/supercharge-grpc-workflows
https://connect.build/docs/introduction
https://github.com/bufbuild/buf


## testing
https://copyconstruct.medium.com/testing-in-production-the-safe-way-18ca102d0ef1
[[data/better go mock lib|better go mock lib]]
https://github.com/gtramontina/ooze mutation testing https://github.com/avito-tech/go-mutesting
https://github.com/buger/goreplay
https://github.com/rekby/fixenv fixtures
https://habr.com/ru/articles/751220/ database testing
https://github.com/mfridman/tparse
https://habr.com/ru/articles/758888/
https://habr.com/ru/companies/karuna/articles/764326/ coverage in prod/integration tests
[приемочные тесты](https://habr.com/ru/articles/765892/)
https://github.com/shoenig/test
https://github.com/mockoon/mockoon/ mock api by openapi spec

### http/grpc clients
https://github.com/Kong/insomnia/releases
https://hexmos.com/lama2/index.html http requests file format
https://github.com/usebruno/bruno http client and requests file format
[debugger](https://github.com/go-delve/delve)
[goroutine-inspect/interactive tool to analyze Golang goroutine dump](https://github.com/linuxerwang/goroutine-inspect)
https://github.com/carlmjohnson/requests simple http client

### integration testing
https://tavern.readthedocs.io/en/latest/basics.html
https://github.com/k1LoW/runn
https://github.com/zoncoen/scenarigo
https://habr.com/ru/companies/sbermarket/articles/739968/


## application structure
https://github.com/evrone/go-clean-template
https://github.com/ThreeDotsLabs/wild-workouts-go-ddd-example
https://git.codemonsters.team/guides/ddd-code-toolkit/src/branch/dev Руководство по внедрению Domain Driven Design
https://github.com/mikestefanello/pagoda


## libraries
[data structures](https://github.com/zyedidia/generic)
https://www.grank.io/ golang libraries trends
https://github.com/avelino/awesome-go
https://threedots.tech/post/list-of-recommended-libraries/
https://github.com/rwxrob/awesome-go
https://github.com/moznion/go-optional
[golang/fast json](https://github.com/bytedance/sonic)
https://trendshift.io/
[library for calling C functions from Go without Cgo](https://github.com/ebitengine/purego)
https://github.com/maragudk/gomponents html generator
https://github.com/go-faster
[golang/actor framework](https://github.com/ergo-services/ergo) [ergo/saga](https://github.com/ergo-services/ergo/blob/master/tests/saga_test.go )
[golang/scraping](https://go-colly.org/)### validation

https://github.com/cohesivestack/valgo
https://github.com/go-ozzo/ozzo-validation

### utils
[mega utils lib](https://github.com/duke-git/lancet)
https://github.com/gookit/goutil utils
https://github.com/dropbox/godropbox

### concurrency
[[data/asynchronous,parallel_programming]]
[[reference/backend/Concurrency in Go.pdf]]
[sourcegraph/conc: Better structured concurrency for go](https://github.com/sourcegraph/conc)

### di
https://habr.com/ru/company/kaspersky/blog/699994/
https://github.com/ompluscator/genjector

### logging
httpmetrics, traces, logshub.com/charmbracelet/log
[golang::exp/slog](https://pkg.go.dev/golang.org/x/exp/slog)

https://www.komu.engineer/blogs/11/opentelemetry-and-go
https://opentelemetry.io/docs/instrumentation/go/getting-started/
### http servers framework
https://github.com/valyala/fasthttp
https://github.com/aceld/zinx
https://betterprogramming.pub/an-introduction-to-gain-part-1-writing-high-performance-tcp-application-df5f7253e54a
https://haykot.dev/blog/reduce-boilerplate-in-go-http-handlers-with-go-generics/
https://github.com/go-chi/render
https://github.com/gofiber/fiber
[golang/web pages framework](https://pushup.adhoc.dev/)

### telegram clients
https://go-pkgz.umputun.dev/notify/#telegram has nil dereference panic
https://github.com/go-telebot/telebot
https://github.com/gotd/td
https://github.com/go-telegram-bot-api/telegram-bot-api

## formatters, pretty printers
https://github.com/kr/pretty pretty print values as go code
[golang/colored pretty print](https://github.com/k0kubun/pp)


## code formatters, linters
[gofumpt/stricter go fmt](https://github.com/mvdan/gofumpt)
[spellchecker](https://github.com/codespell-project/codespell)
```bash
docker run -v $(pwd):/data:ro python:3.9 sh -c "pip3 install codespell 2>&1 >/dev/null; codespell /data"
```
https://github.com/go-critic/go-critic
[project structure linter](https://habr.com/ru/articles/751174/)
[duplicate code finder](https://github.com/boyter/dcd)
https://github.com/bep/punused


## langs
[distributed lang](https://www.unison-lang.org/)

### scripting
https://expr-lang.org/
https://github.com/glycerine/zygomys golang/embedded lisp








