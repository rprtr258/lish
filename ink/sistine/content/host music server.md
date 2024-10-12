#in
[self hosting](/note/self%20hosting)

use exec [jobs](http://45.87.153.219:4646/ui/jobs) if possible
upload music files to [s3](https://console.cloud.yandex.ru/folders/b1gcjejr56jpj3kms7d9/storage/buckets) bucket
mount s3 using [goofys](https://github.com/kahing/goofys) or [geesefs](https://github.com/yandex-cloud/geesefs)
serve music using [navidrome](https://www.navidrome.org/docs/overview/) (replace with executable)
```bash
docker run --name navidrome --user $(id -u):$(id -g) -v /mnt/hdd/music:/music -v $(pwd):/data -p 4533:4533 deluan/navidrome:latest
```
serve player ui using  [feishin](https://github.com/jeffvli/feishin)
```bash
docker run --name feishin -p 9180:9180 ghcr.io/jeffvli/feishin:latest
```
download and save youtube playlist music somehow