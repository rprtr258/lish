#someday_maybe #project

[link](https://github.com/rprtr258/cronus)
https://cloud.google.com/scheduler

## features
- [ ] implement workers: running docker containers: [pr](https://github.com/rprtr258/cronus/pull/3)
    - [ ] test manually
    - [ ] make sample docker containers for cron tasks
[nice feature descriptions](https://www.easycron.com)
[nice ui?](https://cronitor.io/cron-job-monitoring)
store artifacts temporarily in s3
alerts
    on task result differs from previous
        telegram
        email
ephemeral - remove/stop after first success

- написать docker-compose
- api versioning
- auth
- run containers somehow more natively
  https://docs.docker.com/engine/api/sdk/
  https://pkg.go.dev/github.com/opencontainers/runc@v1.1.4/libcontainer#section-readme
  https://github.com/containerd/containerd/blob/main/docs/getting-started.md
  https://pkg.go.dev/github.com/containerd/containerd
  https://github.com/fsouza/go-dockerclient
  https://github.com/opencontainers/runc
  https://github.com/bcicen/ctop
  https://github.com/containers/podman
  https://github.com/dankohn/libpod
  https://pkg.go.dev/github.com/lxc/lxd/client
  https://github.com/testcontainers/testcontainers-go
  https://github.com/acouvreur/sablier
- store logs, traces, metrics in sentry, prometheus
- cron format improvements
    only extend, not change, so regular cron syntax is valid
    make running every Friday 13th possible, or last Friday of month or similar
    timezones?
    allow day to be negative: count from end of month

```embed
title: 'Reliable Cron across the Planet - ACM Queue'
image: 'https://queue.acm.org/img/acmqueue_logo.gif'
description: ''
url: 'https://queue.acm.org/detail.cfm?id=2745840'
```
