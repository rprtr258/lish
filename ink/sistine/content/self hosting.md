#project #sre

- [ ] перенести coredns в nomad
- [x] сделать ansible playbook для поднятия nomad
- [ ] сделать ansible playbook для поднятия consul (если разберусь как его юзать блять)
- [ ] настроить секьюрити для consul

- links on deploy management
    - actual dashboards
        - [nomad](http://45.87.153.219:4646/ui/)
        - [consul](http://45.87.153.219:8500/ui/)
    - experimenting on yandex vms
        - [nomad](https://158.160.44.108:4646/ui/)
        - [vms](https://console.cloud.yandex.ru/folders/b1gcjejr56jpj3kms7d9/compute/instances)

[music server](/note/host%20music%20server)
https://blog.tjll.net/distributed-homelab-cluster/#exposing-applications https://blog.tjll.net/too-simple-to-fail-nomad-caddy-wireguard/
https://www.aapanel.com/new/index.html
https://hub.docker.com/_/adminer
[Как перестать велосипедить или 4 self-hosted сервиса для начинающего СТО](https://habr.com/ru/post/693198/)
https://github.com/awesome-selfhosted/awesome-selfhosted https://awesome-selfhosted.net/index.html
[email](https://poste.io/) https://maddy.email/
https://habr.com/ru/companies/ruvds/articles/729742/
https://medium.com/@alexfoleydevops/terraform-fastapi-encryption-975f1348f69b
https://medium.com/@alexfoleydevops/localhost-love-with-packer-ansible-vagrant-75767dc78a87
https://awstip.com/one-dollar-devops-terraform-docker-nginx-c72443ac0918
[Gitea + Drone + Nginx + Portainer. Пошаговое руководство по деплою аналога github на своём железе](https://habr.com/ru/post/703408/)
https://woodpecker-ci.org/
https://github.com/catppuccin/gitea
https://gitness.com/
https://dopos.github.io/dcape/baseapps/
https://github.com/concourse/concourse
[Домашнее облако](https://habr.com/ru/post/692008/)
https://github.com/mikeroyal/Self-Hosting-Guide
https://github.com/apitable/apitable
mail https://github.com/mjl-/mox
file listing https://github.com/alist-org/alist
https://fmnx.su/core/infr
https://github.com/filebrowser/filebrowser

```embed
title: 'Админка для Private Docker Registry (Registry Admin)'
image: 'https://habrastorage.org/getpro/habr/upload_files/82b/df5/ad9/82bdf5ad906dc2602322d08d6b9644f2.png'
description: 'Концепция контейнеризации на базе Docker, и ему подобных технологий, для многих разработчиков стала незаменимым инструментом доставки своих продуктов конечным пользователям в виде полностью...'
url: 'https://habr.com/ru/post/709988/'
```

```embed
title: 'Docker - Coder v2 Docs'
image: 'https://coder.com/og-image.png'
description: 'Setup Coder with Docker'
url: 'https://coder.com/docs/v2/latest/quickstart/docker'
```

docker-compose:
```yaml
# https://fleet.linuxserver.io/
#   https://demo.hedgedoc.org/
#   https://github.com/linuxserver/docker-airsonic-advanced
#   ? https://github.com/rembo10/headphones
#   ? https://github.com/healthchecks/healthchecks
#   ? https://hub.docker.com/r/linuxserver/heimdall
#   ? https://nextcloud.com/
#   ? http://piwigo.org/
# https://stackoverflow.com/questions/60064792/docker-container-connect-volume-to-remote-server-with-ssh-keys
# https://habr.com/ru/post/699374/
services:
  tinyproxy:
    image: francisbesset/tinyproxy # or 3proxy ?
    ports:
      - 1111:8888
    # args: 185.102.11.146
  mstream:
    image: lscr.io/linuxserver/mstream:latest
    environment:
      - PUID=1000
      - PGID=1000
      - TZ=Europe/London
    volumes:
      - type: volume
        source: /mnt/hdd/music
        target: /music
        read_only: true
      - /home/rprtr258/t:/config
    ports:
      - 3000:3000
    restart: unless-stopped
https://www.reddit.com/r/selfhosted/comments/xxb1p3/quarterly_post_sharing_your_favorite_tools_a/
```

https://nerdyarticles.com/a-clutter-free-life-with-paperless-ngx/
https://github.com/Dullage/flatnotes

monitor pages e.g. stolyarov posts/[guestbook](http://stolyarov.info/guestbook) e.g. using https://github.com/e1abrador/web.Monitor
[setup nomad guide](https://thekevinwang.com/2022/11/20/nomad-cluster)
https://mrkaran.dev/posts/home-server-nomad/
https://synergeticlabs.com/email-alchemy/ how to self host email