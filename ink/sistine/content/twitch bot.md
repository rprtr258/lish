#someday_maybe #project

[S3 with db backups](https://console.cloud.yandex.ru/folders/b1gcjejr56jpj3kms7d9/storage/buckets/backups--twitch-bot--db)

## technical debt
- [x] настроить бекапы бд ✅ 2023-05-04
- [ ] choose and use some authorization/permissions system for commands (see backend.authorization)
    https://github.com/etcd-io/bbolt
    maybe use ldap
        https://www.digitalocean.com/community/tutorials/understanding-the-ldap-protocol-data-hierarchy-and-entry-components
    go server: https://github.com/glauth/glauth
    check different permission schemes https://stackoverflow.com/questions/3177361/modelling-a-permissions-system
- [ ] docker image user setup 🔽
- [ ] generate token нормально, fix `oauth token is missing`, refresh it(?)
    прочитать
        https://dev.twitch.tv/docs/authentication
        https://twitchtokengenerator.com/
        https://github.com/Satont/go-helix/blob/main/helix.go#L178
        https://github.com/gempir/go-twitch-irc/issues/189#issuecomment-1320209811
- [ ] pasta search
    - https://kekg.xyz/
    - https://pepeprikol.peepo.club/s/pepepast?before=472
- [ ] add monitoring, alerts, health checks 🔽
## bugfixes
- [x]  fix `http://` to `https://` in blab responses
- [x] записать про периодическую оплату серверов (для прокси и для бота)
- [ ] fix ddos with commands 🔽
- [ ] fix sakefile work_dir

[sic](https://git.suckless.org/sic/file/sic.c.html)
[ii](https://git.suckless.org/ii/file/FAQ.html)
https://eax.me/tcp-server/
https://eax.me/network-application-mistakes/
https://insobot.handmade.network/
https://beej.us/guide/bgnet/html/#intro
https://habr.com/ru/post/677170/

commands are template with lisp-like queries, like so:
```lisp
!followage = $(followage $(user))
где
followage $user = `follow юзера $user`
user = пользователь, введший команду, либо первый аргумент, если есть
```
команды могут делать действия(опционально) и давать ответ, который подставляется вместо команды
нульарные команды(без аргументов) могут так же служить в качестве констант или динамически вычисляемых переменных типо `$(stream-title)`

![[data/static/old/someday_maybe/programming_projects/twitch-bot/architecture.png]]
