#project

clear docker images/containers/apt periodically
```sh
docker system prune
sudo apt update
sudo apt upgrade
sudo apt autoremove
sudo snap refresh
flatpak refresh
```
deps remove unnecessary apt and other deps
gtd bot to cron task on vps
check ssh server
git pull all branches on backuped repos
purge docker images,containers

# spam twitch chat
```c
* * * * * printf "lizardPls" | TWITCH_OAUTH_TOKEN=$(cat /home/rprtr258/secure-waters-73337/oauth.txt) /home/rprtr258/twitch-bot-api/sender screamlark
```

# apt update
```c
* * * * * flock ~/cron/.errors echo "apt update" $(~/cron/apt-update.sh 2>&1) >>~/cron/.errors
```

# backup configs
```c
* 0 * * * /home/rprtr258/GTD/reference/linux_config/upd.sh
```
