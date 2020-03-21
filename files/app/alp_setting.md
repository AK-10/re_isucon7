# alpの設定(アクセスログ解析ツール)
## install

```bash
$ wget https://github.com/tkuchiki/alp/releases/download/v1.0.3/alp_linux_amd64.zip 
$ unzip alp_linux_amd64.zip
$ sudo install alp /usr/local/bin/alp
$ which alp # 確認
```

## setting
nginxのhttpディレクティブに以下を設定

```nginx
http {

	log_format ltsv "time:$time_local"
	    "\thost:$remote_addr"
	    "\tforwardedfor:$http_x_forwarded_for"
	    "\treq:$request"
	    "\tmethod:$request_method"
	    "\turi:$request_uri"
	    "\tstatus:$status"
	    "\tsize:$body_bytes_sent"
	    "\treferer:$http_referer"
	    "\tua:$http_user_agent"
	    "\treqtime:$request_time"
	    "\truntime:$upstream_http_x_runtime"
	    "\tapptime:$upstream_response_time"
	    "\tcache:$upstream_http_x_cache"
	    "\tvhost:$host";
	access_log  /var/log/nginx/access.log ltsv;
}
```

アクセスログの消去,nginx再起動

```bash
$ sudo rm /var/log/nginx/access.log
$ sudo systemctl restart nginx
```


## アクセスログ解析
```bash
$ alp -f /var/log/nignx/access.log
```

