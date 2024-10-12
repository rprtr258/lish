#project

[sadservers](https://sadservers.com/)

## Если речь идёт о базе, "без которой никуда"
* ну, поставить незнакомый дистр линукса можете?
* А FreeBSD?
* А дальше с ними работать через сеть удалённо в терминальном режиме?
* А поднять на них веб-сервер?
* А почтовую систему поднять?
* А воткнуть в машину вторую сетёвку и превратить её в маршрутизатор?
* А сетку из десятка машин собрать?
* А NAT настроить?
* А порт статический пробросить (через NAT снаружи внутрь)?
* А бекапы автоматические заточить?

[51 задание для оттачивания навыков системного администрирования](https://proglib.io/p/become-sysadmin/)
https://github.com/dastergon/awesome-sre


## domain registration
https://www.reg.ru/user/account/
https://skt.ru/domain-registration/


## hostings
https://vercel.com/
https://fly.io/
https://www.clever-cloud.com/

## SRE
https://github.com/dastergon/awesome-sre
https://github.com/SquadcastHub/awesome-sre-tools
[Google - Site Reliability Engineering](https://sre.google/books/)
https://slurm.io/sre
[[reference/devops, sre, admin/enterprise-roadmap-to-sre.pdf]]
https://habr.com/ru/company/southbridge/blog/715762/
https://habr.com/ru/company/otus/blog/720030/
https://habr.com/ru/post/718796/
[cooperation notebooks for sre](https://fiberplane.com/)
https://medium.com/trendyol-tech/how-search-crashed-suddenly-every-sunday-the-path-of-sre-88ca9431d34d
https://habr.com/ru/companies/skbkontur/articles/739774/
https://yandex.ru/jobs/pages/devops\_cloud
https://www.youtube.com/playlist?list=PLjCCarnDJNstX36A6Cw\_YD28thNFev1op
https://sre.google/books/
https://habr.com/ru/articles/747618/
https://github.com/bregman-arie/sre-checklist
https://www.srepath.com/
https://github.com/upgundecha/howtheysre

```embed
title: 'Data driven SRE'
image: 'https://habrastorage.org/getpro/habr/upload_files/23a/0e1/af8/23a0e1af84766136b1960a428a8549d5.png'
description: 'Начнем эту увлекательную историю издалека. Во время первого локдауна, в начале 2020 года, сайт Леруа Мерлен испытал бóльшую нагрузку, чем когда-либо. Проводившие много времени дома и на даче наши...'
url: 'https://habr.com/ru/company/leroy_merlin/blog/712960/'
```

```embed
title: 'A Journey into Site Reliability Engineering'
image: 'https://images.thoughtbot.com/blog-images/social-share-default.png'
description: 'While Rails gained a lot of popularity among companies to develop products quickly, technical debt and scalability issues were challenges that also gained space in this context. Let&rsquo;s talk about some SRE fundamentals that can address those situations.'
url: 'https://thoughtbot.com/blog/a-journey-into-site-reliability-engineering'
```

```embed
title: 'SRE: управление инцидентами'
image: 'https://habrastorage.org/getpro/habr/upload_files/e45/339/a3e/e45339a3e31ef2ee3a20a07fb2141cd3.png'
description: 'Автор статьи: Рустем Галиев IBM Senior DevOps Engineer &amp; Integration Architect. Официальный DevOps ментор и коуч в IBM Привет Хабр! Не так давно общался с SRE в нашей команде и он рассказал мне о…'
url: 'https://habr.com/ru/company/otus/blog/722892/'
```

Программа
День 1: знакомство с теорией SRE, настройка мониторинга и алёртинга
В первый день вы познакомитесь с теорией SRE, научитесь настраивать мониторинг и алёртинг, а также объединитесь в команду с другими участниками интенсива.
Расскажем про метрики SLO, SLI, SLA и как они соотносятся с требованиями бизнеса. Поделимся Best Practices по настройке мониторинга и правилами для пожарной команды. Дадим первые практические кейсы.
Мониторинг
Зачем нужен мониторинг
Symptoms vs Causes
Black-Box vs White-Box Monitoring
Golden Signals
Перцентили
Alerting
Observability
Практика: Делаем базовый дашборд и настраиваем необходимые алерты
Теория SRE
SLO, SLI, SLA
Durability
Error budget
Практика: Добавляем на дашборд SLO/SLI + алерты
Практика: Первая нагрузка системы
Практика, решение кейса: зависимость downstream.
В большой системе существует много взаимозависимых сервисов, и не всегда они работают одинаково хорошо. Особенно обидно, когда с вашим сервисом порядок, а соседний, от которого вы зависите, периодически уходит в down.
Учебный проект окажется именно в таких условиях, а вы сделаете так, чтобы он все равно выдавал качество на максимально возможном уровне.
Управление инцидентами
Resiliencе Engineering
Как выстраивается пожарная бригада
Насколько ваша команда эффективна в инциденте
7 правил для лидера инцидента
5 правил для пожарного
HiPPO — highest paid person's opinion. Communications Leader
Практика, решение кейса: зависимость upstream.
Одно дело, когда вы зависите от сервиса с низким SLO. Другое дело, когда ваш сервис является таковым для других частей системы. Так бывает, если критерии оценки не согласованы: например, вы отвечаете на запрос в течение секунды и считаете это успехом, а зависимый сервис ждёт всего 500 мск и уходит с ошибкой.
В кейсе обсудим важность согласования метрик и научимся смотреть на качество глазами клиента.
День 2: решение проблем с окружением и архитектурой
Второй день практически полностью построен вокруг решения двух кейсов: проблемы с окружением (здесь будет подробный разбор Health Checking) и проблемы с архитектурой. Спикеры расскажут про работу с постмортерами (post mortem) и дадут шаблоны, которые вы сможете использовать в своей команде.
Health Checking
Health Check в Kubernetes
Жив ли наш сервис?
Exec probes
InitialDelaySeconds
Secondary Health Port
Sidecar Health Server
Headless Probe
Hardware Probe
Практика, решение кейса: проблема с окружением, билеты купить невозможно.
Задача Healthcheck — обнаружить неработающий сервис и заблокировать трафик к нему. И если вы думаете, что для этого достаточно сделать рутом запрос к сервису и получить ответ, то вы ошибаетесь: даже если сервис ответит, это не гарантирует его работоспособность — проблемы могут быть в окружении.
Через этот кейс вы научитесь настраивать корректный Healthcheck и не пускать трафик туда, где он не может быть обработан.
Практика работы с постмортемами
Практика: Пишем постмортем по предыдущему кейсу и разбираем его со спикерами.
Решение проблем с инфраструктурой
Мониторинг PostgreSQL
SLO/SLI для PostgreSQL
Anomaly detection
Практика, решение кейса: проблемы с базой данных.
База данных тоже может быть источником проблем. Например, если не следить за replication relay, то реплика устареет и приложение будет отдавать старые данные. Причём дебажить такие случаи особенно сложно: сейчас данные рассогласованы, а через несколько секунд уже нет, и в чём причина проблемы — непонятно.
Через кейс вы прочувствуете всю боль дебага и узнаете, как предотвращать подобные проблемы.
День 3: Traffic shielding и канареечные релизы
Тут два кейса про высокую доступность продакшена: traffic shielding и canary deployment. Вы узнаете об этих подходах и научитесь их применять. Хардкорной настройки руками не планируем, хотя кто знает.
Traffic shielding
Поведение графиков роста количества запросов и бизнес операций
Понятие saturation и capacity planning
Traffic shielding и внедрение rate limiting
Настройка sidecar с rate-limiting на 100 запросов в секунду
Практика, решение кейса: Traffic shielding, исследуем поведение провайдера под нагрузкой, которую он не в состоянии выдержать.
Когда падает прод? Например, когда мощности рассчитаны на 100 пользователей, а приходит 1000. Вы столкнётесь с подобным кейсом и научитесь делать так, чтобы система не падала целиком, а продолжала обслуживать то количество клиентов, на которое была рассчитана.
Блокируя избыточный трафик, вы сохраните возможность системы выполнять задачи для части пользователей.
Canary Deployment
Стратегии деплоя в k8s (RollingUpdate vs Recreate)
Canary и blue-green стратегии
Обзор инструментов для blue-gree/canary release в k8s
Настройка canary release в GitLab CI/CD
Пояснение схемы работы canary release
Внесение изменений в .gitlab-ci.yml
Практика, решение кейса: проблема с кодом.
Как бы хорошо новые фичи не работали на стейджинге, всегда есть вероятность, что в продакшене что-то пойдёт не так. Снизить потенциальный ущерб можно, если выкатить обновление только на часть пользователей и предусмотреть возможность быстрого отката назад.
Подобный подход называется Canary Deployment, и через практический кейс вы научитесь его применять.
Тема 9: SRE онбординг проекта
В крупных компаниях нередко формируют отдельную команду SRE, которая берёт на поддержку сервисы других отделов. Но не каждый сервис готов к тому, чтобы его можно было взять на поддержку. Расскажем, каким требованиям он должен отвечать.

## devops

https://roadmap.sh/devops
https://github.com/milanm/DevOps-Roadmap
https://www.techworld-with-nana.com/devops-roadmap
https://github.com/milanm/DevOps-Roadmap
https://github.com/ahmadalibagheri/devops-roadmap
https://aws.amazon.com/getting-started/learning-path-devops-engineer/
https://cloud.google.com/certification/cloud-devops-engineer
https://spacelift.io/blog/devops-tools
http://awesome-devops.xyz/
https://github.com/AcalephStorage/awesome-devops
https://xebialabs.com/periodic-table-of-devops-tools/
https://gist.github.com/emedvedev/d27c590280b9bf690793f3dd37212b5c
[От стесняшки до архитектора: какими бывают DevOps и как стать одним из них](https://habr.com/en/company/netologyru/blog/577180/)
[[reference/devops, sre, admin/Microservices Architecture Enables DevOps. An Experience Report on Migration to a Cloud-Native Architecture.pdf]]
https://habr.com/ru/companies/oleg-bunin/articles/728524/
https://habr.com/ru/companies/cloud\_mts/articles/740234/
https://habr.com/ru/articles/745532/
https://habr.com/ru/companies/nixys/articles/514098/
https://aws.plainenglish.io/20-devops-interview-questions-and-answers-with-detailed-explanation-b6cb09969c7c
https://habr.com/ru/articles/755970/
https://github.com/bregman-arie/devops-exercises
https://github.com/LearnWithTwoHeads/devops

```embed
title: 'Как я создавал homelab для учебы на DevOps-инженера'
image: 'https://habrastorage.org/getpro/habr/upload_files/47a/a65/041/47aa65041943742eccd1cf6c49790211.jpg'
description: 'В феврале 2022-го стало ясно, что надо приобретать профессию, востребованную за пределами России. На тот момент я жил в Москве и успешно практиковал как юрист. За плечами у меня была работа в…'
url: 'https://habr.com/ru/post/699372/'
```

```embed
title: '90DaysOfDevOps/2022.md at main · MichaelCade/90DaysOfDevOps'
image: 'https://repository-images.githubusercontent.com/441903012/5a26a00b-f955-437b-b10b-1d102b0f5c4e'
description: 'I am using this repository to document my journey learning about DevOps. I began this process on January 1, 2022, and plan to continue until March 31. I will be dedicating one hour each day, includ...'
url: 'https://github.com/MichaelCade/90DaysOfDevOps/blob/main/2022.md'
```

```embed
title: '90DaysOfDevOps/2023.md at main · MichaelCade/90DaysOfDevOps'
image: 'https://repository-images.githubusercontent.com/441903012/5a26a00b-f955-437b-b10b-1d102b0f5c4e'
description: 'I am using this repository to document my journey learning about DevOps. I began this process on January 1, 2022, and plan to continue until March 31. I will be dedicating one hour each day, includ...'
url: 'https://github.com/MichaelCade/90DaysOfDevOps/blob/main/2023.md'
```

[[reference/devops, sre, admin/The DevOps Handbook - How to Create World-Class Agility, Reliability, and Security in Technology Organizations.epub]]
[[reference/devops, sre, admin/DevOps Handbook 2.pdf]]
[[reference/devops, sre, admin/devops-in-practice.pdf]]
[[reference/devops, sre, admin/Effective\_DevOps.pdf]]
[[reference/devops, sre, admin/Проект-Феникс.pdf]]
[10 задач для девопса, когда уже нечем заняться](https://habr.com/ru/company/ruvds/blog/703872/)

## linux
[Learn Linux, 101: A roadmap for LPIC-1](https://developer.ibm.com/tutorials/l-lpic1-map/)
https://linuxopsys.com/topics/category/commands
https://linuxjourney.com/
https://ss64.com/bash/
https://debian-handbook.info/browse/stable/
https://makelinux.github.io/kernel/map/


## shell mastering
[tmux alternative](https://github.com/zellij-org/zellij)
https://github.com/jlevy/the-art-of-command-line
https://sharats.me/posts/shell-script-best-practices/
[bash linter](https://www.shellcheck.net/)
https://github.com/kellyjonbrazil/jc
http://blog.angel2s2.ru/2008/08/grep-r-grep-i-grep-w-grep-n-grep-x.html
[termbin.com - terminal pastebin](http://termbin.com/)


## virtual machines
https://www.vagrantup.com/
[Quick Start - QEMU documentation](https://www.qemu.org/docs/master/system/quickstart.html)
[Quick Start - OpenNebula 6.2.0 documentation](https://docs.opennebula.io/6.2/quick_start/index.html#qs)
https://www.youtube.com/watch?v=SgJf8uvpDCc kubevirt - kubernetes for virtual machines


## kubernetes
https://kwok.sigs.k8s.io/
https://www.youtube.com/@OldPythonKAA/videos
[Создание Kubernetes-кластера на пальцах или почему это не сложно](https://habr.com/ru/post/713520/)
[Kubernetes через грабли или внедрение в университете](https://habr.com/ru/post/711008/)
https://www.youtube.com/watch?v=sLQefhPfwWE
https://www.youtube.com/playlist?list=PL3SzV1\_k2H1VDePbSWUqERqlBXIk02wCQ
https://developers.redhat.com/developer-sandbox/activities/learn-kubernetes-using-red-hat-developer-sandbox-openshift
https://dev.to/jonatan5524/series/20062
https://kubernetes.io/docs/tutorials/kubernetes-basics/
https://www.youtube.com/watch?v=kTp5xUtcalw
https://habr.com/ru/company/kts/blog/593599/
https://habr.com/ru/post/656639/
https://habr.com/ru/company/domclick/blog/577964/
https://habr.com/ru/post/541118/
https://habr.com/ru/company/ozontech/blog/586308/
https://docs.flagger.app/
https://doc.traefik.io/traefik/
https://medium.com/devopsanswers/kubernetes-the-easy-way-19f69e57738d
https://github.com/kelseyhightower/kubernetes-the-hard-way
https://habr.com/ru/company/flant/blog/685616/
https://spacelift.io/blog/kubernetes-tutorial
https://github.com/stefanprodan/podinfo
https://www.youtube.com/watch?v=s\_o8dwzRlu4
https://www.youtube.com/watch?v=X48VuDVv0do
https://kubernetes.io/blog/2022/10/04/introducing-kueue/
https://developers.redhat.com/developer-sandbox/activities/learn-kubernetes-using-red-hat-developer-sandbox-openshift
https://k3s.io/
https://k3d.io/v5.4.9/
https://helm.sh/docs/intro/quickstart/
https://github.com/komodorio/helm-dashboard
https://github.com/werf/werf
https://itnext.io/kubernetes-i-multi-node-deployment-using-terraform-and-terragrunt-30c40a1238e8
https://spacelift.io/blog/kubernetes-secrets
https://octopus.com/blog/ssh-into-kubernetes-cluster
https://itnext.io/best-alternatives-for-the-top-kubernetes-ide-you-need-to-know-of-87a2cefe2daa
https://carlosbecker.com/posts/k8s-pod-shutdown-lifecycle/
https://habr.com/ru/company/yandex\_praktikum/blog/698626/
https://k0sproject.io/
https://www.youtube.com/playlist?list=PL8D2P0ruohOA4Y9LQoTttfSgsRwUGWpu6
https://kind.sigs.k8s.io/
https://krew.sigs.k8s.io/
https://learnkubernetes.withgoogle.com/
https://microk8s.io/#install-microk8s
https://please.build/codelabs/k8s/
https://github.com/100daysofkubernetes/100DaysOfKubernetes/
https://kubeshark.co/
https://itnext.io/kubernetes-in-a-box-7a146ba9f681
https://werf.io/guides/golang/100\_basic/10\_build.html
https://opensource.com/article/22/6/kubernetes-networking-fundamentals
https://itnext.io/inspecting-and-understanding-service-network-dfd8c16ff2c5
[[projects/devops, sre, admin/KUBERNETES\_A\_Simple\_Guide\_to\_Master\_Kubernetes\_for\_Beginners\_and\_Advanced\_Users\_2020\_Edition\_\_Brian\_Docker.epub]]
https://github.com/techiescamp/kubernetes-learning-path
https://github.com/doitintl/kube-no-trouble
https://habr.com/ru/company/otus/blog/717486/
https://habr.com/ru/post/481662/
https://github.com/derailed/k9s https://k9scli.io/
https://habr.com/ru/company/timeweb/blog/720510/
https://roadmap.sh/kubernetes
https://www.groundcover.com/
https://habr.com/ru/companies/southbridge/articles/727470/
https://github.com/labring/sealos
https://github.com/mr-karan/doggo/
https://chaos-mesh.org/
https://habr.com/ru/companies/southbridge/articles/729228/
https://www.portainer.io/
https://rancherdesktop.io/
https://mesos.apache.org/documentation/latest/
https://habr.com/ru/articles/734928/
https://www.youtube.com/playlist?list=PLz0t90fOInA5IyhoT96WhycPV8Km-WICj
https://github.com/hidetatz/kubecolor
https://betterprogramming.pub/10-antipatterns-for-kubernetes-deployments-e97ce1199f2d
https://david.coffee/why-and-how-i-use-k8s-for-personal-stuff/
https://github.com/octarinesec/kube-scan https://github.com/cyberark/KubiScan
https://github.com/2gis/k8s-handle
https://fluxcd.io/flux/
https://k8s-examples.container-solutions.com/ configs examples
https://awstip.com/tools-to-make-your-terminal-devops-and-kubernetes-friendly-64d27a35bd3f
https://habr.com/ru/companies/southbridge/articles/750264/
https://habr.com/ru/companies/cloud\_mts/articles/750560/ book on k8s
https://habr.com/ru/companies/otus/articles/751752/ k3s
https://gitlab.com/k11s-os/k8s-lessons
https://github.com/derailed/popeye
https://ru.werf.io/guides/golang/100\_basic.html k8s ci/cd
https://iximiuz.com/en/series/working-with-kubernetes-api/
https://github.com/sbstp/kubie alternative to kubectx kubens
https://nubenetes.com/images/learnk8s\_debug\_your\_pods.png
https://kubevela.io/ k8s ui
https://habr.com/ru/companies/T1Holding/articles/767056/
https://github.com/jamiehannaford/what-happens-when-k8s
https://github.com/omerbsezer/Fast-Kubernetes
[List of terminal commands for Kubernetes](https://awstip.com/list-of-terminal-commands-for-kubernetes-ffc63f0dcec0)
https://runbooks.prometheus-operator.dev/
https://sonobuoy.io/docs/v0.57.1/
[Aptakube - Kubernetes Desktop Client](https://aptakube.com/)
[Kustomize - Kubernetes native configuration management](https://kustomize.io/)
[How to Deploy Postgres on Kubernetes | Tutorial](https://www.containiq.com/post/deploy-postgres-on-kubernetes)

![[projects/devops, sre, admin/3T6xi7B.png]]
![[projects/devops, sre, admin/lUwYzzr.png]]

[Не куб, а кубик: kubernetes для не-highload](https://habr.com/ru/post/711440/)
[Как легко пройти собеседование по Kubernetes в 2023 году?](https://habr.com/ru/company/southbridge/blog/713884/)

```embed
title: 'Meshery The Kubernetes and Cloud Native Manager'
image: 'https://meshery.io/images/logos/meshery-gradient.png'
description: 'A CNCF project, Meshery, is the open source, Kubernetes and cloud native manager.'
url: 'https://meshery.io/'
```

```embed
title: 'Crossplane - The cloud-native control plane framework'
image: 'https://crossplane-admin.prod.unomena.io/media/original_images/crossplane-og.jpg'
description: 'Crossplane is a framework for building cloud native control planes without needing to write code. It has a highly extensible backend that enables you to build a control plane that can orchestrate applications and infrastructure no matter where they run, and a highly configurable frontend that puts y…'
url: 'https://www.crossplane.io/'
```

```embed
title: 'A guide to container orchestration with Kubernetes'
image: 'https://opensource.com/sites/default/files/lead-images/kenlon-music-conducting-orchestra.png'
description: 'To learn all about container orchestration with Kubernetes, download our new eBook.'
url: 'https://opensource.com/article/22/6/container-orchestration-kubernetes'
```

https://itnext.io/12-factor-microservice-applications-on-kubernetes-db913008b018

#### kubernetes guide

https://telegra.ph/Introduction-to-Kubernetes-05-18
https://telegra.ph/Cluster-Environment-Building-05-18
https://telegra.ph/Resource-Management-05-18
https://telegra.ph/How-to-Deploy-an-NGINX-Service-in-a-Kubernetes-Cluster-05-19
https://telegra.ph/Pod-Detailed-Explanation-1-05-19
https://telegra.ph/Pod-Detailed-Explanation-2-05-19
https://telegra.ph/In-depth-Explanation-of-Pod-Controllers-05-19
https://telegra.ph/A-Detailed-Explanation-of-Services-05-19
https://telegra.ph/Data-Storage-05-19

## backup
https://github.com/rclone/rclone
[backup postgres in s3 script](https://github.com/Satont/postgresql-backup-s3/blob/master/backup.sh#L70)
[gobackup](https://gobackup.github.io/) \- simple enough to use\, sample config below put in `~/.gobackup/gobackup.yml` or `/etc/gobackup/gobackup.yml` and collect backups with `gobackup perform`
TODO: add periodicity, adding to cron

```yaml
models:
  list:
    compress_with:
      type: tgz
    storages:
      s3:
        type: s3
        bucket: testtset
        region: ru-central1
        access_key_id: YCAJEyHso2njZ08YMNJgd0TTI
        secret_access_key: YCO8If6BEQIz91w2d68aDCxTgopYHzifnFMuKYb4
        endpoint: storage.yandexcloud.net
    databases:
      list-db:
        type: postgresql
        host: localhost
        port: 49155
        database: todos
        username: backend
        password: example
    notifiers:
      telegram:
        type: telegram
        chat_id: 310506774
        token: 5415014383:AAElsIuS3hqpF0PUayvhsIznbbDMSB4Ioh8
```

https://kopia.io/
[borgbackup](https://borgbackup.readthedocs.io/en/stable/) https://www.borgbackup.org/
[autorestic](https://autorestic.vercel.app) or [restic](https://restic.net/) itself
https://github.com/kopia/kopia
https://bup.github.io/

## infrastructure-management

https://spacelift.io/blog/devops-automation-tools
https://www.pulumi.com/blog/iac-recommended-practices-developer-stacks-git-branches/
https://saltproject.io/

```embed
title: 'Модели управления инфраструктурой'
image: 'https://habrastorage.org/getpro/habr/upload_files/a62/d95/6d3/a62d956d3f955fd31211930834d69451.png'
description: 'Управление инфраструктурой даже средней организации является непростой задачей. Большое количество серверов требует постоянного внимания. Установка обновлений и развертывание новых систем все это...'
url: 'https://habr.com/ru/company/otus/blog/709588/'
```

### terraform

https://habr.com/ru/companies/yandex\_praktikum/articles/738682/
https://habr.com/ru/companies/southbridge/articles/744564/
https://habr.com/ru/company/otus/blog/721166/
https://slack.engineering/how-we-use-terraform-at-slack/
https://www.youtube.com/watch?v=q12v5mbMnco
https://terragrunt.gruntwork.io/
https://github.com/terraform-docs/terraform-docs

### cloud guide

https://cloud.yandex.ru/docs/tutorials/infrastructure-management/terraform-quickstart
https://habr.com/ru/post/684964/
https://habr.com/ru/post/685062/
https://habr.com/ru/post/685520/

### pulumi

https://www.pulumi.com/
https://www.pulumi.com/blog/pulumi-kubernetes-new-2022/

### ansible
https://mitogen.networkgenomics.com/ansible_detailed.html
[Автоматизация управления с помощью Ansible](https://habr.com/ru/company/otus/blog/711136/)
[разворачивание k8s с помощью ansible](https://habr.com/ru/articles/751582/)
[Practical%20Network%20Automation](/static/old/someday_maybe/networking/Practical%20Network%20Automation.pdf)
https://habr.com/ru/company/southbridge/blog/691212/
https://habr.com/ru/company/southbridge/blog/688724/
https://github.com/weareinteractive/ansible-pm2
https://ru.hexlet.io/courses/production-basics/lessons/deploy/theory\_unit
https://github.com/apenella/go-ansible
https://habr.com/ru/post/704518/
https://habr.com/ru/post/724450/
https://github.com/geerlingguy/ansible-role-kubernetes https://github.com/geerlingguy/ansible-role-docker https://github.com/nginxinc/ansible-role-nginx https://github.com/geerlingguy/ansible-role-gitlab https://github.com/geerlingguy/ansible-role-postgresql https://github.com/geerlingguy/ansible-role-mysql https://github.com/geerlingguy/ansible-role-rabbitmq https://github.com/geerlingguy/ansible-role-securitya plugins
[Ansible для начинающих](https://habr.com/ru/company/southbridge/blog/714000/)
[[someday\_maybe/programming\_projects/Ansible alternative in common lisp]]
https://vaiti.io/kak-ispolzovat-ansible-dlya-prostyh-i-slozhnyh-zadach/
[Компонентный подход к Ansible или как навести порядок в инфраструктурном коде](https://habr.com/ru/companies/just_ai/articles/772382/)
[Ansible это вам не bash](https://habr.com/ru/articles/494738/)
[Фильтры Ansible: превращаем сложное в простое](https://habr.com/ru/articles/778206/)

## multiserver shells
[remote multiserver automation](https://github.com/capistrano/capistrano)
[distributed shell](https://easyengine.io/tutorials/linux/dsh/)
https://github.com/alajmo/sake
https://github.com/pressly/sup
ansible?
puppet
chef
https://pyinfra.com/

## ssh

[Наглядное руководство по SSH-туннелям](https://habr.com/ru/company/flant/blog/691388/)
https://iximiuz.com/en/posts/ssh-tunnels/
https://faun.pub/use-ssh-port-forwarding-to-connect-to-resources-221534e9037
[Архитектура SSH. Узел-бастион и принцип нулевого доверия](https://habr.com/ru/company/ruvds/blog/720244/)
[SSH tricks](https://matt.might.net/articles/ssh-hacks/)

```embed
title: 'GitHub - moul/assh: make your ssh client smarter'
image: 'https://repository-images.githubusercontent.com/3049190/52761c00-ca74-11e9-991f-1062c00708f1'
description: ':computer: make your ssh client smarter. Contribute to moul/assh development by creating an account on GitHub.'
url: 'https://github.com/moul/assh'
```

## secrets management

https://learn.hashicorp.com/vault

## alerting

https://github.com/megaease/easeprobe
scripted alerts https://balerter.com/
alerts https://github.com/keephq/keep

## observability

`observability = metrics + traces + logs`
[beautiful logs/monitoring/etc tools](https://betterstack.com/)
https://ozontech.github.io/file.d
[open source performance monitoring](https://github.com/SigNoz/signoz)
https://pyroscope.io/
[Practical Monitoring: Effective Strategies for the Real World](https://www.amazon.com/Practical-Monitoring-Effective-Strategies-World-ebook/dp/B076XZWQVW/ref=sr_1_1?keywords=monitoring&qid=1555329911&s=gateway&sr=8-1)
[sqshq/sampler](https://github.com/sqshq/sampler)
https://github.com/stefanprodan/dockprom
https://vector.dev/
https://habr.com/ru/company/ruvds/blog/701034/
https://grafana.github.io/grafonnet-lib/examples/
https://blog.creekorful.org/2020/12/how-to-setup-easily-elk-docker-swarm/
https://konstellationapp.com/
[прикольная штука для observability](https://coroot.com/) https://habr.com/ru/companies/flant/articles/742030/
[Как мы перешли с Elastic на Grafana stack и сократили расходы в несколько раз](https://habr.com/ru/company/m2tech/blog/693504/)
https://cronitor.io/status-pages
https://github.com/keyval-dev/odigos
https://opsverse.io
https://qryn.metrico.in/
https://signoz.io
https://flow.com/engineering-blogs/golang-services-improving-observability
https://www.datadoghq.com/
https://www.influxdata.com - influxdb
https://www.influxdata.com/time-series-platform/telegraf/
https://habr.com/ru/company/kts/blog/719938/
https://www.highlight.io/
dependencies grafana map https://github.com/groundcover-com/caretta
[Основные аспекты наблюдаемости систем](https://habr.com/ru/companies/ruvds/articles/727072/)
https://habr.com/ru/companies/monq/articles/727938/
https://habr.com/ru/companies/otus/articles/727556/
https://habr.com/ru/companies/kts/articles/723980/
https://www.parseable.io/docs
https://dev.to/bmf\_san/building-a-monitoring-infrastructure-starting-with-a-container-1efd
https://github.com/teletrace/teletrace
https://github.com/openobserve/openobserve
https://docs.victoriametrics.com/
https://habr.com/ru/companies/ruvds/articles/746086/
https://gethelios.dev/
https://grafana.com/blog/2023/07/13/how-to-monitor-kubernetes-network-and-security-events-with-hubble-and-grafana/
https://github.com/adriannovegil/awesome-observability
https://github.com/ccfos/nightingale
https://github.com/oklog/oklog
https://github.com/ductnn/domolo
https://github.com/grafana/agent
[Parca - Open Source infrastructure-wide continuous profiling](https://www.parca.dev/')
[logging in nomad](https://atodorov.me/2021/07/09/logging-on-nomad-and-log-aggregation-with-loki/)
[Мониторинг: смысл, цели и универсальные рецепты](https://habr.com/ru/company/web3_tech/blog/711816/)
https://brandur.org/canonical-log-lines

### prometheus

https://www.youtube.com/watch?v=8KaHRs93UJw
https://github.com/samber/awesome-prometheus-alerts
https://github.com/roaldnefs/awesome-prometheus
https://habr.com/ru/articles/747350/
https://m3db.io

```embed
title: 'Monitoring Linux instances with Prometheus and Grafana'
image: 'https://dev.to/social_previews/article/1336334.png'
description: 'If you want to know how is everything going in your servers, you must monitor them, yeah, sometimes...'
url: 'https://dev.to/aldorvv__/monitoring-linux-instances-with-prometheus-and-grafana-441i'
```

```embed
title: 'ТОП-10 экспортеров для Prometheus 2023'
image: 'https://habr.com/share/publication/711936/5c4bcc8479bbd31f4551e87eac6d9184/'
description: 'Статья Основы мониторинга (обзор Prometheus и Grafana) оборвалась на самом интересном месте. Автор предложил искать и использовать актуальные экспортеры, а читатель такой – окей, где референс? Что ж,...'
url: 'https://habr.com/ru/post/711936/'
```

## ci/cd

[CI/CD для проекта в GitHub с развертыванием на AWS EC2](https://habr.com/ru/post/536118/)
[Простой CI/CD на Ansible Semaphore](https://habr.com/ru/post/645927/)
[Develop Your CI/CD Pipelines as Code and Run Them Anywhere](https://dagger.io/)
https://bass-lang.org/guide.html
[self hosted git](https://gitea.io/en-us/)
https://blog.logrocket.com/creating-separate-monorepo-ci-cd-pipelines-github-actions/
[[projects/devops, sre, admin/Continuous\_Delivery.pdf]]
[[projects/devops, sre, admin/nepreryvnaya\_integratsiya\_uluchshenie\_kachestva\_programmnogo\_obespecheniya\_i\_snizhenie\_riska\_signature\_series\_3643739.pdf]]
[[projects/devops, sre, admin/nepreryvnoe\_razvertyvanie\_po\_avtomatizatsiya\_protsessov\_sborki\_testirovaniya\_i\_vnedreniya\_novykh\_versii\_programm\_signature\_series\_3643776.pdf]]

## letsencrypt clients

[https://doka.guide/tools/ssl-certificates/](https://doka.guide/tools/ssl-certificates/)
[https://letsencrypt.org/docs/client-options/](https://letsencrypt.org/docs/client-options/)

## docker

[jesseduffield/lazydocker](https://github.com/jesseduffield/lazydocker)
[Полное практическое руководство по Docker: с нуля до кластера на AWS](https://habr.com/ru/post/310460/)
https://training.play-with-docker.com/
https://docs.docker.com/storage/storagedriver/#container-and-layers
https://doka.guide/tools/docker-data-management/
https://github.com/docker/docker-bench-security
[[reference/devops, sre, admin/Securing Docker.pdf]]
https://infosecadalid.com/2021/08/30/containers-rootful-rootless-privileged-and-super-privileged/
dockerfile user

```Dockerfile
RUN addgroup --system javauser && \
adduser -S -s /usr/sbin/nologin -G javauser javauser
RUN chown -R javauser:javauser /opt/app
USER javauser
CMD ["java", "-jar", "/opt/app/app.jar"]


RUN addgroup -g 11211 memcache && adduser -D -u 11211 -G memcache memcache
USER memcache
```

### docker alternatives

https://highload.today/8-besplatnyh-alternativ-docker-na-2022-god/https://github.com/containers/podman-desktop
https://podman.io/
https://itnext.io/goodbye-docker-desktop-hello-minikube-3649f2a1c469

## virtual machines

```embed
title: 'qemantra - Control QEMU like magic!'
image: 'https://qemantra.pspiagicw.xyz/qemantra.svg'
description: 'qemantra is a cli application designed to create, run and manage virtual machines using QEMU/KVM and implemented in Go'
url: 'https://qemantra.pspiagicw.xyz/'
```

### unikernels

https://unzip.dev/0x005-unikernels/

```embed
title: 'OPS - Easily Build and Run Unikernels'
image: 'https://ops.city/dist/img/logo.png'
description: 'OPS is a unikernel compilation and orchestration tool. It is the only tool that allows instant building and running of raw linux binaries as unikernels.'
url: 'https://ops.city/'
```

## eBPF

https://github.com/iovisor/bpftrace
https://github.com/iovisor/bpftrace/blob/master/docs/tutorial\_one\_liners.md
https://brendangregg.com/bpf-performance-tools-book.html
https://github.com/iovisor/bcc
https://github.com/Gui774ume/ebpfkit
https://github.com/ebpfdev/explorer

## unsorted

no sudo
https://help.ubuntu.com/community/Sudoers
https://www.google.com/search?q=how+to+use+linux+without+sudo
Как жить без sudo
переключаетесь на текстовую консоль по Ctrl-Alt-F[1,2,3...], там входите в систему под root'ом и делаете, что нужно.
mount в повседневной жизни не требует полномочий root'а, просто соответствующие файловые системы (всякие флешки и cdrom'ы) должны быть внесены в /etc/fstab с флагом user. fdisk в повседневной жизни вообще не нужен, ради него можно и на текстовую консоль переключиться. Перезагрузка — если нужна именно перезагрузка — выполняется переключением на ту же консоль и там (то есть уже не в XWin) нажатием старого доброго Ctrl-Alt-Del. Если бы вы использовали Devuan или другую систему с классическим System V init (а не с systemd), я бы вам мог рассказать, как на Ctrl-Alt-Del вместо перезагрузки повесить останов системы a.k.a. выключение машины, но как это сделать с systemd — понятия не имею (и не хочу иметь, systemd не стоит времени, которое тратится на то, чтобы с ним разобраться).
Ну а чтобы выключить машину можно было, не входя под root'ом, можно на /sbin/halt навесить SetUid bit. Только тогда стоит, наверное, предпринять дополнительные меры — скажем, создать группу для этого или хотя бы воспользоваться уже существующей (например, wheel обычно уже есть). Один раз под root'ом сделать вот так:
chown root:wheel /sbin/halt
chmod 4750 /sbin/halt
потом ещё себя, любимого, добавить в группу wheel (я бы это сделал просто редактированием /etc/group, но можно и командой usermod, тут дело вкуса) — и всё, команду halt можно будет давать с правами вашего обычного пользователя.

https://matt.might.net/articles/how-to-emergency-web-scaling/
http://toly.github.io/blog/2016/04/20/simple-0-downtime-blue-green-deployments/
https://www.raphaelmichel.de/blog/2014/deploying-django.html
[[reference/backend/\_architecture/Systems.Performance.Enterprise.and.the.Cloud.pdf]]
[Как масштабировать дата-центры. Доклад Яндекса](https://habr.com/en/company/yandex/blog/476146/)
http://rus-linux.net/MyLDP/BOOKS/slackware/index.html
https://www.portainer.io/
[Деплой приложений в VM, Nomad и Kubernetes](https://habr.com/ru/company/lamoda/blog/451644/)
[Kubernetes, микросервисы, CI/CD и докер для ретроградов: советы по обучению](https://habr.com/ru/company/itsumma/blog/499102/)
https://0xax.gitbooks.io/linux-insides/content/Booting/
[Настройка LEMP сервера с помощью docker для простых проектов. Часть первая: База](https://habr.com/en/company/nixys/blog/661443/)
https://github.com/mr-karan/hydra
https://slurm.io/devops-tools-to-dev
[Управление учетными записями в Linux. Часть 3. Различные способы поднятия привилегий](https://habr.com/ru/company/otus/blog/691756/)
https://killercoda.com/playgrounds
https://github.com/Swordfish-Security/awesome-devsecops-russia

https://github.com/ripienaar/free-for-dev
https://github.com/minio/sidekick
https://github.com/score-spec/spec
https://goteleport.com/
https://github.com/Netflix/conductor
https://github.com/awesome-foss/awesome-sysadmin
https://github.com/unixorn/sysadmin-reading-list
https://containrrr.dev/watchtower/

```embed
title: 'LINSTOR — это как Kubernetes, но для блочных устройств (обзор и видео доклада)'
image: 'https://habrastorage.org/getpro/habr/upload_files/97a/3ab/82a/97a3ab82a41b50da48ad90a3a39e9893.png'
description: 'В июне я выступил на объединенной конференции DevOpsConf &amp; TechLead Conf 2022 . Доклад был посвящен LINSTOR — Open Source-хранилищу от компании LINBIT (разработчики DRBD). Основной идеей...'
url: 'https://habr.com/ru/company/flant/blog/680286/'
```

```embed
title: 'Сертификаты Let’s Encrypt и ACME вообще во внутренней сети'
image: 'https://habrastorage.org/getpro/habr/upload_files/c5d/969/419/c5d969419fb9f83832b72576ae67cd1a.jpg'
description: 'Обычно внутри корпоративной сети нынче полно всяких приложений, и хотелось бы чтобы они работали по SSL. Можно, конечно, поднять свой УЦ, раздать сертификаты, прописать пользователям свой корневой...'
url: 'https://habr.com/ru/post/708510/'
```

[[reference/devops, sre, admin/Accelerate - Building and Scaling High Performing Technology Organisations.pdf]]
[[projects/devops, sre, admin/Масштабирование\_приложений.pdf]]

```embed
title: 'Сбор и анализ логов в Linux'
image: 'https://habrastorage.org/getpro/habr/upload_files/cab/933/b8f/cab933b8f52d89aad90a74be77b5e280.png'
description: 'Журналирование событий, происходящих в&nbsp;системе является неотъемлемой частью функционала любого серьезного программного обеспечения. Операционная система или&nbsp;приложение должны…'
url: 'https://habr.com/ru/company/otus/blog/714266/'
```

https://www.keycloak.org/
https://github.com/deepfence/ThreatMapper
https://landscape.cncf.io/
https://github.com/collabnix/dockerlabs
https://collabnix.github.io/kubelabs/
https://world.hey.com/dhh/introducing-mrsk-9330a267
https://shubhsharma19.hashnode.dev/advanced-cloud-concepts
https://medium.com/4th-coffee/10-new-devops-tools-to-watch-in-2023-e974dbb1f1bb
https://www.sonatype.com/products/nexus-repository

#### container builders

* https://github.com/containers/skopeo
* https://github.com/GoogleContainerTools/kaniko
* https://github.com/containers/buildah
* https://github.com/genuinetools/img

### how to setup ssh on server

* in `/etc/ssh/ssgd_config`: set

```yaml
PermitRootLogin: prohibit-password
PasswordAuthenticatioin: no
```

* `service ssh restart`
* setup `root` keys [[reference/devops, sre, admin/devops, sre, admin#how to add keys for ssh server]]

### how to add keys for ssh server

* locally generate key with `ssh-gen` and following instructions
* remotely in home directory of user add to `~/.ssh/authorized_keys` new row with public key

### how to add user

`adduser rprtr258`

https://github.com/mephux/envdb
https://osquery.io/
https://github.com/cncf/trailmap/blob/master/CNCF\_TrailMap\_latest.pdf
https://aws.amazon.com/ru/lightsail/
[Синий свет — зеленый свет: релизим без даунтаймов](https://habr.com/ru/company/oleg-bunin/blog/720986/)
https://roadmap.sh/best-practices/aws
https://longhorn.io/

```embed
title: 'Configuration options | Plausible docs'
image: 'https://user-images.githubusercontent.com/85956139/132954658-2d5bc2c3-22c2-4300-b9c6-cbe4f8f8987e.png'
description: 'The easiest way to get started with Plausible is with our official managed service in the Cloud. It takes 2 minutes to start counting your stats with a worldwide CDN, high availability, backups, security and maintenance all done for you by us. Our managed hosting can save a substantial amount of dev…'
url: 'https://plausible.io/docs/self-hosting-configuration'
```

```embed
title: 'GoatCounter – open source web analytics'
image: 'https://static.zgo.at/logo.svg'
description: 'Easy web analytics. No tracking of personal data.'
url: 'https://www.goatcounter.com/'
```

https://medium.com/@yoanante/the-easy-way-to-use-a-vps-and-deploy-your-personal-projects-2edd8b242f3b
https://www.ericsdevblog.com/posts/how-to-set-up-a-web-server/
https://habr.com/ru/company/yandex/blog/499534/
https://registry.terraform.io/providers/kreuzwerker/docker/latest/docs
https://developer.hashicorp.com/terraform/tutorials/docker-get-started/docker-change?optInFrom=learn
https://business-science.github.io/shiny-production-with-aws-book/
https://habr.com/ru/companies/netologyru/articles/729010/
https://habr.com/ru/companies/timeweb/articles/733058/
https://slurm.io/linux-admin-base
https://github.com/diego-treitos/linux-smart-enumeration
https://sensu.io/
https://tproger.ru/articles/21-zadacha-iz-opyta-devops-inzhenera/
https://editor.networkpolicy.io/?id=it4UBTYVii4t10gS
https://ko.build/
https://github.com/TwiN/gatus
https://habr.com/ru/companies/otus/articles/742040/
https://www.theforeman.org/
https://firecracker-microvm.github.io/
https://labs.iximiuz.com/main#courses
https://github.com/antonputra/tutorials/blob/main/docs/contents.md
check who uses port: `lsof -i :8080` or `fuser 8080/tcp`, file: `lsof <file>` or `fuser <file>`, kill process by port usage: `fuser -k 8080/tcp`
https://goteleport.com/blog/kubectl-cheatsheet/
https://kyverno.io/
https://fluxcd.io/
https://github.com/SnellerInc/sneller
https://github.com/kitabisa/teler Real-time HTTP Intrusion Detection
https://dzone.com/articles/platform-setup-how-to-structure-your-infrastructur
https://www.youtube.com/playlist?list=PL4\_hYwCyhAvZcOr5sJzuLmze2F6wPms-A
https://www.youtube.com/playlist?list=PL4\_hYwCyhAva5JvmFBIUc\_-JqZk4xfc5f
https://www.youtube.com/playlist?list=PL4\_hYwCyhAvZe\_MY9Lb64ObGN3l4s-84A
https://about.gitlab.com/blog/2023/06/27/efficient-devsecops-workflows-with-rules-for-conditional-pipelines/
https://skarnet.org/software/s6/
https://habr.com/ru/articles/496492/ LXD
https://github.com/getsops/sops
https://github.com/bcicen/ctop
настроить алерты на ошибки в логах
https://habr.com/ru/users/alitenicole/posts/
https://encore.dev/
https://github.com/ContainerSolutions/k8s-deployment-strategies
https://habr.com/ru/companies/nixys/articles/512766/ https://habr.com/ru/companies/nixys/articles/513118/ https://habr.com/ru/companies/nixys/articles/513578/ https://habr.com/ru/companies/nixys/articles/513918/
https://habr.com/ru/companies/badoo/articles/507718/ collect logs with loki
https://deepflow.io/community.html
ops articles https://habr.com/ru/articles/118475/ https://habr.com/ru/articles/118966/ https://habr.com/ru/articles/50008/ https://habr.com/ru/articles/91896/ https://habr.com/ru/articles/353762/ https://habr.com/ru/companies/ruvds/articles/702570/ https://habr.com/ru/companies/ruvds/articles/724676/
https://gist.github.com/rprtr258/b9e9f50ec9f8718f2c3debf001516efb 9 Platfrom Engineers’ Mistakes
https://github.com/statsd/statsd
https://graphite.readthedocs.io/en/stable/
https://devconnected.com/
https://iximiuz.com/en/series/implementing-container-manager/
https://cloud.beeline.ru/devopscloud/
https://github.com/Infisical/infisical secrets manager
https://habr.com/ru/articles/755624/
https://www.youtube.com/@KirillSemaev/videos
https://www.container-security.site/
https://biriukov.dev/docs/fd-pipe-session-terminal/1-file-descriptor-and-open-file-description/
https://github.com/AdminTurnedDevOps/100DaysOfContainersAndOrchestration
https://sysxplore.com/
https://mcs.mail.ru/cloud-native-diy/
https://devhands.io/ru/
https://github.com/topics/linux-administration
https://pboyd.io/posts/securing-a-linux-vm/
[collaborative ssh](https://sshx.io/)
[Подборка видео с последнего SREcon](https://habr.com/ru/articles/773448/)
https://mmonit.com/monit/
https://github.com/juju/juju
https://fly.io/dist-sys/ system design practical tasks
https://habr.com/ru/articles/775776/
https://rootly.com/blog/status-pages-101-how-to-create-a-status-page-you-and-your-customers-will-actually-want-to-use https://firehydrant.com/blog/your-guide-to-better-incident-status-pages/ status pages
[CSI SPEC](https://github.com/container-storage-interface/spec/blob/master/spec.md)
[cloud native infra book](https://github.com/CloudNativeInfra/cni/blob/master/en/SUMMARY.md)
https://cloudprober.org/
https://landontclipp.github.io/blog/2023/06/22/prefer-systemd-timers-over-cron/#benefits_1
https://gruntwork.io/devops-checklist
https://github.com/adnanh/webhook
https://brandur.org/alerting
https://flowpipe.io/
https://jvns.ca/strace-zine-v3.pdf
https://github.com/jonashaag/klaus
https://github.com/bregman-arie/devops-resources
https://linkedin.github.io/school-of-sre/
https://google.github.io/building-secure-and-reliable-systems/raw/toc.html
https://sre.google/workbook/table-of-contents/
https://www.serf.io/
https://github.com/Swfuse/devops-interview/blob/main/references.md
https://github.com/Swfuse/devops-interview/blob/main/interview.md