# isucon7 振り返り

## 点数遷移
初期 4000 ~ 5500
最終 110000 ~ 150000程度


## やったこと
- alp導入
- messageJsonifyで発生しているN+1を修正
- dbにある画像をファイルの読み出しに変更，またnginx側でキャッシュするようにした
- historyのN+1を解消
- indexをはった
- 画像をgo-cacheでキャッシュするようにした(不採用)
- deploy scriptを書いた

## やり損ねたこと
- redisを使う

## alp install
- sudo apt isntall -y wget unzip
- wget https://github.com/tkuchiki/alp/releases/download/v1.0.3/alp_linux_amd64.zip
- unzip alp_linux_amd64.zip
- sudo mv alp /usr/local/bin/

- nginx.conf に以下を実行

```
log_format ltsv "time:$time_local"
                    "\thost:$remote_addr"
                    "\tforwardedfor:$http_x_forwarded_for"
                    "\treq:$request"
                    "\tstatus:$status"
                    "\tmethod:$request_method"
                    "\turi:$request_uri"
                    "\tsize:$body_bytes_sent"
                    "\treferer:$http_referer"
                    "\tua:$http_user_agent"
                    "\treqtime:$request_time"
                    "\tcache:$upstream_http_x_cache"
                    "\truntime:$upstream_http_x_runtime"
                    "\tapptime:$upstream_response_time"
                    "\tvhost:$host";

access_log  /var/log/nginx/access.log ltsv;
error_log /var/log/nginx/error.log;

# logの停止
#access_log off;
#error_log /dev/null crit;
```

## 実行
`$ sudo alp ltsv --file /var/log/nginx/access.log -r --sum | head -n 30`


## slowQuery
- /etc/my.cnf or /etc/mysql/mysqld.confの[mysqld]に以下を追記

```
slow_query_log=ON
long_query_time = 0.0001
slow_query_log_file = /var/log/mysql/slow.log
```

`/var/log/mysql/slow.log` のパーミッションに注意, 777でいい

### 確認
```
mysql> SHOW variables LIKE '%slow%';
mysql> SHOW variables LIKE '%long%';
```

