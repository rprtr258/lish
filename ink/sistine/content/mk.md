#project

https://github.com/rprtr258/mk/pull/1

----
https://github.com/rprtr258/mk/pull/2
- [x] host agent with listing docker containers with statuses
- [x] [[next_actions/mk -- in client, check remote agent version (if present), if not actual, replace with newer agent binary on remote|mk -- in client, check remote agent version (if present), if not actual, replace with newer agent binary on remote]]
- [x] implement docker container state reconciliation in agent
- [x] execute docker container state reconciliation using agent
----
- [x] rewrite https://github.com/rprtr258/fimgs/blob/master/justfile and template ✅ 2023-05-16
- [ ] update makefile guide https://github.com/rprtr258/rprtr258/blob/master/how-to-write-justfiles.md
- [ ] generate readme from docstrings like here https://github.com/skratchdot/open-golang/blob/master/Makefile
- [ ] add to styleguide
```
use [`rprtr258/mk`](https://github.com/rprtr258/mk) over `Makefile`
```

## cmake
https://eax.me/cmake/
https://gist.github.com/mbinna/c61dbb39bca0e4fb7d1f73b0d66a4fd1
https://github.com/onqtam/awesome-cmake
https://aosabook.org/en/v1/cmake.html

## existing build systems
https://github.com/gopinath-langote/1build
https://go-gilbert.github.io/
https://github.com/goyek/goyek
https://github.com/magefile/mage
https://github.com/tj/mmake
https://github.com/oxequa/realize
https://github.com/go-task/task
https://github.com/taskctl/taskctl
https://github.com/joerdav/xc
https://gruntjs.com/configuring-tasks
https://gulpjs.com/
https://pydoit.org/
https://www.pyinvoke.org/
https://buck2.build/docs/why/
https://premake.github.io/docs/your-first-script
http://doc.cat-v.org/plan_9/4th_edition/papers/mk
https://bob.build/
https://shakebuild.com/manual
https://www.pantsbuild.org/

## articles
[Build Systems à la Carte](https://dl.acm.org/doi/pdf/10.1145/3236774)
https://blogs.ncl.ac.uk/andreymokhov/the-task-abstraction/
[Build Systems à la Carte. Theory and practice](https://www.cambridge.org/core/services/aop-cambridge-core/content/view/097CE52C750E69BD16B78C318754C7A4/S0956796820000088a.pdf/div-class-title-build-systems-a-la-carte-theory-and-practice-div.pdf)
[Build Scripts with Perfect Dependencies](https://arxiv.org/pdf/2007.12737.pdf)
[Non-recursive Make Considered Harmful](https://eprints.ncl.ac.uk/file_store/production/226639/BA1AC951-8C2E-46C8-92F1-74FA95507892.pdf)
https://web.archive.org/web/20100411063824/http://www.conifersystems.com/whitepapers/gnu-make/

## example makefiles
https://github.com/Thiht/smocker/blob/master/Makefile
https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
https://danishpraka.sh/posts/using-makefiles-for-go/
https://www.olioapps.com/blog/the-lost-art-of-the-makefile
https://rosszurowski.com/log/2022/makefiles
https://github.com/stripe/stripe-cli/blob/master/Makefile

## docs generation
```makefile
gen-docs: # generate readme.md
```
POC компиляции `readme.md` из разных элементов (картинки, вывод `--help`, статический текст) можно посмотреть [тут](https://github.com/rprtr258/fimgs/blob/master/cmd/mk/mk.go#L79).

### env table generation
можно (в теории) генерировать из конфига с доп.комментариями таблицу используемых переменных окружения, например из:
```go
type Config struct {
	// service name
	ServiceName string `env:"SERVICE_NAME,required"`
	// redis
	Redis struct {
		# host
		RedisHost string `env:"HOST" envDefault:"redis"`
		# port
		RedisPort int `env:"PORT" envDefault:"1234"`
	} `envPrefix:"REDIS_"`
}
```
генерировать:

|var name|default value|required?|description|
|-|-|-|-|
|`SERVICE_NAME`||:heavy_check_mark:|service name|
|`REDIS_HOST`|`redis`||redis host|
|`REDIS_PORT`|`1234`||redis port|

также можно сделать cli, которая заставляет заполнять все обязательные поля, и мб необязательные

## load dotenv
можно добавить [include .env](https://lithic.tech/blog/2020-05/makefile-dot-env)
https://stackoverflow.com/questions/4225497/include-files-depended-on-target
https://stackoverflow.com/questions/30300830/load-env-file-in-makefile

## other
написать статью на хабр?
https://earthly.dev/blog/make-tutorial/
[[data/machine learning & data science]]
- https://github.com/caarlos0/svu
    - `release-major`
    - `release-minor`
    - `release-patch`

https://github.com/umputun/spot https://simplotask.com/ - cool docs
https://github.com/geerlingguy/ansible-for-devops
https://github.com/gotestyourself/gotestsum
https://github.com/princjef/gomarkdoc
https://github.com/braintree/runbook
https://github.com/ansible-semaphore/semaphore web ui
https://www.zackproser.com/blog/bubbletea-state-machine
https://habr.com/ru/companies/skyeng/articles/743458/
https://github.com/oxequa/realize#config-sample
https://arslan.io/2019/07/03/how-to-write-idempotent-bash-scripts/
https://stackoverflow.com/questions/19390600/ansible-lineinfile-duplicates-line
https://encore.dev/
https://gist.github.com/rprtr258/54c4f6f943e15880ad03b9a3c23fec16 makefile guide
https://brunch.io/

```makefile
lint-dockerfile:   # lint Dockerfile
	docker run --rm -i hadolint/hadolint < Dockerfile
```

https://gist.github.com/rprtr258/54c4f6f943e15880ad03b9a3c23fec16
https://sakecli.com/
https://habr.com/ru/companies/southbridge/articles/748788/
https://docs.docker.com/build/bake/reference/
https://github.com/run-x/opta
https://github.com/google/zx
https://rnemet.dev/posts/tools/taskfile/
https://www.gnu.org/software/make/manual/html_node/Functions.html
https://www.pluralith.com/
https://medium.com/cloud-native-daily/build-ci-cd-pipeline-using-your-favorite-language-with-dagger-42ccfc43e073
[[data/static/old/projects/devops, sre, admin/devops, sre, admin|devops, sre, admin#ansible]]
https://gitlab.com/ita1024/waf
https://peace.mk/
https://pkg.go.dev/github.com/bitfield/script
https://www.pulumi.com/registry/packages/docker/
https://www.winglang.io/
https://cli.github.com/manual/
https://hexmos.com/ansika
https://concourse-ci.org/
https://github.com/xorpaul/g10k
[lama2](https://hexmos.com/lama2/index.html) but based on jsonnet (actualy just use mk)

[GitHub - cl-adams/adams: UNIX system administration in Common Lisp](https://github.com/cl-adams/adams)

[Adams 0.1 released](https://www.reddit.com/r/Common_Lisp/comments/fz5vq1/adams_01_released/)

[> Lisp is not and has never been the right tool for configuration management. W... | Hacker News](https://news.ycombinator.com/item?id=10399749)

[GitHub - atgreen/surrender](https://github.com/atgreen/surrender)

translate configs to ansible's yamls with lisp→yaml

[lisp to json/yaml/ini/toml/css/etc](lisp%20to%20js%20d0598.md)
https://habr.com/ru/post/549874/
https://sakecli.com/
https://arslan.io/2019/07/03/how-to-write-idempotent-bash-scripts/
https://github.com/eradman/rset
https://github.com/rollcat/judo
https://github.com/saltstack/salt/
https://pop.readthedocs.io/en/latest/
https://pop-book.readthedocs.io/en/latest/
https://gitlab.com/vmware/idem/idem
https://github.com/opsmop/opsmop
https://gist.github.com/Fizzadar/953ec565d10a1ef9e1d5
https://bundlewrap.org/
https://github.com/bearstech/nuka
https://www.paramiko.org/
https://www.fabfile.org/
https://pyinfra.com/
https://github.com/srevinsaju/togomak
https://github.com/fsouza/go-dockerclient
https://github.com/jsiebens/hashi-up
https://github.com/uber-go/mock/blob/main/ci/test.sh
https://github.com/opsmop/opsmop https://github.com/opsmop/opsmop-demo
https://blog.dave.tf/post/new-kubernetes/
https://github.com/basecamp/kamal https://evilmartians.com/chronicles/mrsk-hot-deployment-tool-or-total-game-changer
https://aosabook.org/en/500L/contingent-a-fully-dynamic-build-system.html
https://gianarb.it/blog/reactive-planning-and-reconciliation-in-go
https://github.com/ublue-os/fleek
https://github.com/viert/xc
https://github.com/bitfield/script
https://github.com/goyek/goyek
