## vpn
https://habr.com/ru/post/678458/
https://habr.com/ru/post/474250/
https://openvpn.net/community-resources/how-to/
[ad block](https://pi-hole.net/)
https://www.wireguard.com/quickstart/
https://www.digitalocean.com/community/tutorials/how-to-set-up-and-configure-an-openvpn-server-on-ubuntu-20-04-ru
https://www.digitalocean.com/community/tutorials/how-to-set-up-wireguard-on-ubuntu-20-04
https://pronomad.ru/blog/vpn-options/
https://habr.com/en/post/153855/
https://blog.cloudflare.com/warp-for-desktop/
https://github.com/notthebee/ansible-easy-vpn
https://github.com/angristan/openvpn-install
https://github.com/pritunl/pritunl?tab=readme-ov-file#install-from-source
[How To Set Up WireGuard on Debian 11 | DigitalOcean](https://www.digitalocean.com/community/tutorials/how-to-set-up-wireguard-on-debian-11)
https://www.qovery.com/blog/build-your-own-network-with-linux-and-wireguard
https://byurrer.ru/vpn-networking
https://gist.github.com/rprtr258/dc858d8c2c802be53ac011740e9c1f4d

## proxy
https://www.envoyproxy.io/
https://caddyserver.com/
https://openlitespeed.org/#features
[Introduction to modern network load balancing and proxying](https://blog.envoyproxy.io/introduction-to-modern-network-load-balancing-and-proxying-a57f6ff80236) 
[Безопасный HTTPS-прокси менее чем за 10 минут](https://habr.com/ru/post/687512/)
https://hysteria.network/
https://github.com/francisbesset/docker-tinyproxy
https://serverfault.com/questions/219740/redirecting-some-url-requests-to-a-lan-proxy
https://github.com/Dreamacro/clash
https://github.com/nadoo/glider

### haproxy
https://www.google.com/search?q=haproxy

### squid
https://github.com/sameersbn/docker-squid
https://www.opennet.ru/base/net/squid_inst.txt.html
https://www.google.com/search?q=squid+реклама+ACL
https://habr.com/ru/articles/733016/

### hand crafted forward proxy
find how to do something (maybe this [https://www.alibabacloud.com/blog/how-to-use-nginx-as-an-https-forward-proxy-server_595799](https://www.alibabacloud.com/blog/how-to-use-nginx-as-an-https-forward-proxy-server_595799)) instead of that
https://github.com/rprtr258/forward-proxy

### reverse proxy
https://blog.container-solutions.com/wtf-is-a-service-mesh
https://blog.container-solutions.com/wtf-is-istio
https://www.envoyproxy.io/
[fast reverse proxy](https://github.com/fatedier/frp)
    https://github.com/rprtr258/frp-test
    https://github.com/anydasa/frp-example/
https://nginxproxymanager.com
https://fabiolb.net/

#### nginx
https://eax.me/nginx/
https://www.freecodecamp.org/news/docker-nginx-letsencrypt-easy-secure-reverse-proxy-40165ba3aee2/
https://nginx-playground.wizardzines.com
https://www.youtube.com/watch?v=7VAI73roXaY
https://doka.guide/tools/nginx-web-server/
https://gongled.ru/blog/nginx-dynamic-ssl/
https://byurrer.ru/nginx-config-multisite

#### traefik
https://github.com/Satont/dotfiles/tree/42e025ea2eaeeaf4b3feebedc0af4f8d0cb27e66/Tools/traefik

#### consul
[Consul начало приключения](https://habr.com/ru/post/693128/)
https://www.consul.io
https://habr.com/ru/company/oleg-bunin/blog/486842/
https://developer.hashicorp.com/nomad/tutorials/integrate-consul/consul-service-mesh
https://medium.com/velotio-perspectives/a-practical-guide-to-hashicorp-consul-part-1-5ee778a7fcf4
https://github.com/mehdihadeli/awesome-software-architecture/blob/main/docs/service-discovery/consul.md
https://byurrer.ru/consul
https://github.com/ansible-community/ansible-consul
https://developer.hashicorp.com/nomad/docs/integrations/consul-connect
https://www.hashicorp.com/blog/consul-connect-integration-in-hashicorp-nomad

### SSL/TLS
https://doka.guide/tools/ssl-certificates/
https://doka.guide/recipes/lets-encrypt-nginx/
https://byurrer.ru/wildcard-lets-encrypt-certbot
https://blog.cloudflare.com/introducing-cfssl/

## dns
https://coredns.io/
  [CoreDNS — DNS-сервер для мира cloud native и Service Discovery для Kubernetes](https://habr.com/ru/company/flant/blog/331872/)
https://dev.to/chrisachard/dns-record-crash-course-for-web-developers-35hn
https://habr.com/ru/company/ozontech/blog/722042/
https://github.com/zmap/zdns
провайдер глобальных dns записей https://www.reg.ru/user/account/#/card/84737675/nss
dns сервер https://coredns.io/

## iptables
![](https://iximiuz.com/laymans-iptables-101/iptables-stages-white.png)
`chain` is sequence of `rule`s which are `match filter`+`target` expressions deciding what to do with packet. If all `match filter`s passed but no `target` found, `chain` `policy` is looked up to find default `target` (which in `policy` must be `ACCEPT` or `DROP`). For user-defined `chain`s `policy` is always `RETURN`. Some `chain`s are predefined:
- `PREROUTING` - packet from outside ingested by network interface
- `POSTROUTING` - packet is going to be sent to network interface
- `FORWARD` - packet is being filtered
- `INPUT` - packet is sent to some process
- `OUTPUT` - packet is sent from some process
Example `target`s:
- `DROP` - drop packet and don't send it anywhere
- `ACCEPT` - accept filter, stop processing `chain` and let packet pass further
- `LOG` - log packet
- `RETURN` - jump back to the caller `chain`
- `<CHAIN NAME>` - send to some another `chain`
Sample `iptables` commands:
```bash
# list existing rules, -S = --list-rules
iptables --list-rules
# -A = --append <CHAIN> <PARAMETERS...> - `A`ppends to the end of the chain, so lowest priority
# -I = --insert <CHAIN> [rulenum] <PARAMETERS...> - `I`nserts rule to rulenum position, 1 if not specified which is highest priority
# -R = --append <CHAIN> <rulenum> <PARAMETERS...> - `R`eplace rule
# -D = --delete <CHAIN> [rulenum] - delete rule
# -j = --jump <TARGET>
# -i = --input <INTERFACE>
# -o = --output <INTERFACE>
iptables -A FORWARD -o docker0 -j DOCKER
iptables -A FORWARD -i docker0 -o docker0 -j ACCEPT
## -p = --proto <tcp,udp>
## -m = --module <EXTENSION MODULE> used for custom filters:
##    - multiport - match multiple tcp/udp ports: --dports destination ports, --sports source ports
##    - conntrack - match connection state:
##        - NEW - connection is establishing, no reply packets were observed
##        - ESTABLISHED - connection established, both incoming and outcoming packets observed
##        - RELATED - connection established, (but ?, presumably stale)
iptables -A INPUT -p tcp -m multiport --dports 80,443 -m conntrack --ctstate NEW,ESTABLISHED -j ACCEPT
iptables -A OUTPUT -p tcp -m multiport --dports 80,443 -m conntrack --ctstate ESTABLISHED -j ACCEPT
# -P = --policy <CHAIN> <TARGET>
iptables -P FORWARD DROP
# -N = --new-chain <CHAIN>
# -X = --delete-chain <CHAIN>
# -E = --rename-chain <OLD-CHAIN> <NEW-CHAIN>
iptables -N DOCKER
```
`table`s in turn are collection of `chain`s with various purposes. Default `table` is `filter`. The sequence is following for app and router
![](https://www.frozentux.net/iptables-tutorial/images/tables_traverse.jpg)
(NOT SURE AT ALL ON CORRECTNESS):
![](https://iximiuz.com/laymans-iptables-101/tables-precedence.png)
![](https://iximiuz.com/laymans-iptables-101/tables-precedence-route.png)
`table` is specified using `-t <TABLE>` flag

https://en.wikipedia.org/wiki/Nftables iptables on steroids https://wiki.nftables.org/wiki-nftables/index.php/Quick_reference-nftables_in_10_minutes
https://iximiuz.com/en/posts/laymans-iptables-101/
https://www.digitalocean.com/community/tutorials/iptables-essentials-common-firewall-rules-and-commands
https://byurrer.ru/iptables
https://linuxforum.ru/viewtopic.php?id=228
https://www.digitalocean.com/community/tutorials/how-to-list-and-delete-iptables-firewall-rules
https://www.frozentux.net/iptables-tutorial/iptables-tutorial.html
http://xgu.ru/wiki/iptables

## ssh
https://iximiuz.com/en/posts/ssh-tunnels/
https://goteleport.com/blog/ssh-tunneling-explained/
https://github.com/francoismichel/ssh3

---



## other

http://linux-ip.net/html/
https://opensource.com/business/16/8/introduction-linux-network-routing
https://www.actualtechmedia.com/wp-content/uploads/2017/12/CUMULUS-NETWORKS-Linux101.pdf
https://events.static.linuxfound.org/sites/events/files/slides/2016%20-%20Linux%20Networking%20explained_0.pdf
https://blog.packagecloud.io/monitoring-tuning-linux-networking-stack-receiving-data/
https://blog.packagecloud.io/monitoring-tuning-linux-networking-stack-sending-data/
https://itsfoss.com/basic-linux-networking-commands/
[fast_inet](https://smirn0v.notion.site/86d54b31c0044e7d96696f7e45e05235)
[high-performance networking](https://hpbn.co/)
https://habr.com/ru/company/flant/blog/329830/

- [ ] study wg-quick commands
```go
$ wg-quick up wg0
[#] ip link add wg0 type wireguard
[#] wg setconf wg0 /dev/fd/63
[#] ip -4 address add 10.8.0.3/24 dev wg0
[#] ip link set mtu 1420 up dev wg0
[#] resolvconf -a tun.wg0 -m 0 -x
[#] wg set wg0 fwmark 51820
[#] ip -6 route add ::/0 dev wg0 table 51820
[#] ip -6 rule add not fwmark 51820 table 51820
[#] ip -6 rule add table main suppress_prefixlength 0
[#] nft -f /dev/fd/63
[#] ip -4 route add 0.0.0.0/0 dev wg0 table 51820
[#] ip -4 rule add not fwmark 51820 table 51820
[#] ip -4 rule add table main suppress_prefixlength 0
[#] sysctl -q net.ipv4.conf.all.src_valid_mark=1
[#] nft -f /dev/fd/63

$ wg-quick down wg0
[#] ip -4 rule delete table 51820
[#] ip -4 rule delete table main suppress_prefixlength 0
[#] ip -6 rule delete table 51820
[#] ip -6 rule delete table main suppress_prefixlength 0
[#] ip link delete dev wg0
[#] resolvconf -d tun.wg0 -f
[#] nft -f /dev/fd/63
```
https://alexgallacher.com/how-to-setup-nginx-ssl-on-docker/
https://habr.com/ru/post/583814/
https://systemweakness.com/rustscan-is-way-faster-than-nmap-heres-how-to-install-and-use-i-f74c895defef
https://habr.com/ru/company/oleg-bunin/blog/722778/
https://habr.com/ru/company/oleg-bunin/blog/723092/
https://habr.com/ru/post/725144/
https://habr.com/ru/companies/otus/articles/727662/
[[projects/networking/Network_Programming_with_Go_Learn_to_Code_Secure_and_Reliable_Network_Services.pdf]]
https://protohackers.com/problems
https://habr.com/ru/company/southbridge/blog/718534/
https://blog.container-solutions.com/wtf-is-cilium
https://habr.com/ru/articles/727868/
https://github.com/mr-karan/doggo/
https://github.com/cilium/pwru
https://book.systemsapproach.org/index.html
https://beej.us/guide/bgnet/
https://fabiolb.net/
https://habr.com/ru/articles/741246/
https://habr.com/ru/companies/ruvds/articles/580648/
https://t.me/roskomsvoboda/10924
https://iximiuz.com/en/series/computer-networking-fundamentals/
https://linkerd.io/
[[data/ctf|ctf]]
https://doka.guide/tools/tcp-udp-protocols/
https://habr.com/ru/articles/747616/
https://github.com/anderspitman/SirTunnel ngrok alternative script
https://github.com/anderspitman/awesome-tunneling
https://github.com/s0rg/decompose find connections between docker containers
https://github.com/kevwan/tproxy tool to monitor socket connections
https://linkmeup.ru/blog/1188/

subfinder (https://github.com/projectdiscovery/subfinder) — поиск поддоменов
nuclei (https://github.com/projectdiscovery/nuclei) — сканирование уязвимостей
aix (https://github.com/projectdiscovery/aix) — взаимодействие с API больших языковых моделей
alterx (https://github.com/projectdiscovery/alterx) — генерация словарей
asnmap (https://github.com/projectdiscovery/asnmap) — сопоставление диапазонов сетей организации с использованием ASN 
cdncheck (https://github.com/projectdiscovery/cdncheck) — обнаружение технологий по заданному IP-адресу
chaos-client (https://github.com/projectdiscovery/chaos-client) — взаимодействие с API Chaos DB
cloudlist (https://github.com/projectdiscovery/cloudlist) — получение активов от облачных провайдеров
dnsx (https://github.com/projectdiscovery/dnsx) — dig/host/nslookup на стероидах
httpx (https://github.com/projectdiscovery/httpx) — многоцелевой набор HTTP-инструментов
katana (https://github.com/projectdiscovery/katana) — сканирование веб-приложений и поиск информации — как паук в Burp Suite, только из командной строки
mapcidr (https://github.com/projectdiscovery/mapcidr) — получение информации для заданной подсети/диапазона CIDR
naabu (https://github.com/projectdiscovery/naabu) — сканер портов
https://habr.com/ru/articles/761798/
https://github.com/nicocha30/ligolo-ng
https://amarchenko.dev/translate/2023-10-02-network/
https://developer.hashicorp.com/nomad/docs/networking
[OSI Deprogrammer](https://docs.google.com/document/d/1iL0fYmMmariFoSvLd9U5nPVH1uFKC7bvVasUcYq78So/preview)
https://mrkaran.dev/posts/nomad-networking-explained/
[serve files from git](https://github.com/saitho/git-file-webserver)
[ngrok in rust](https://github.com/ekzhang/bore)
https://wiki.nikiv.dev/networking/
https://iximiuz.com/en/posts/service-proxy-pod-sidecar-oh-my/
https://www.cloudskillsboost.google/focuses/1743?parent=catalog
https://mt165.co.uk/blog/software-networking-and-interfaces-recording/
https://blog.container-solutions.com/wtf-is-a-service-mesh
[CNI SPEC](https://github.com/containernetworking/cni/blob/main/SPEC.md)
https://developer.hashicorp.com/nomad/docs/networking/cni
https://labs.iximiuz.com/tutorials/container-networking-from-scratch
[ip command guide](https://linuxopsys.com/topics/linux-ip-command)
https://linuxjourney.com/
https://www.brendangregg.com/chaosreader.html
https://www.brendangregg.com/Solaris/network.html
https://developers.redhat.com/blog/2019/05/17/an-introduction-to-linux-virtual-interfaces-tunnels
https://docs.strongswan.org/docs/5.9/features/routeBasedVpn.html
https://jvns.ca/tcpdump-zine.pdf
https://wizardzines.com/networking-tools-poster.pdf
https://github.com/milosgajdos/tenus
https://jvns.ca/debugging-zine.pdf
https://xgu.ru/wiki/IPsec_в_Linux
http://unixwiz.net/techtips/iguide-ipsec.html
http://xgu.ru/wiki/Таблица_маршрутизации
https://lartc.org/howto/
https://plantuml.com/nwdiag
https://www.youtube.com/watch?v=WL0ZTcfSvB4
https://gist.github.com/miglen/70765e663c48ae0544da08c07006791f
https://docs.netbox.dev/en/stable/
https://www.reddit.com/r/sysadmin/comments/4ux166/do_you_know_a_site_to_learn_more_about_networking/
https://github.com/leandromoreira/linux-network-performance-parameters
https://en.wikipedia.org/wiki/Iproute2
http://intronetworks.cs.luc.edu/current/html/index.html
https://en.wikipedia.org/wiki/Layer_2_Tunneling_Protocol
https://how-did-i-get-here.net
https://learn.cantrill.io/p/tech-fundamentals
https://www.google.com/search?q=setup+ipsec
https://github.com/hwdsl2/setup-ipsec-vpn
https://learnk8s.io/kubernetes-network-packets#intercepting-and-rewriting-traffic-with-netfilter-and-iptables
https://developers.redhat.com/blog/2018/10/22/introduction-to-linux-interfaces-for-virtual-networking
https://hechao.li/2017/12/13/linux-bridge-part1/ https://hechao.li/2018/01/31/linux-bridge-part2/
https://jvns.ca/blog/2016/12/22/container-networking/
https://www.kernel.org/doc/Documentation/networking/vxlan.txt
https://blog.scottlowe.org/2013/09/04/introducing-linux-network-namespaces/
https://www.youtube.com/watch?v=6v_BDHIgOY8
[socat usage examples](http://www.dest-unreach.org/socat/doc/socat.html#EXAMPLES)
https://latest.gost.run/
https://github.com/sumup-oss/gocat
