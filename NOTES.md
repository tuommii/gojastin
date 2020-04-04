# NOTES

## Git
`git tag -a v0.1.0 -m "my version 0.1.0"`

`git push origin v0.1.0`

## Server Notes
Log: `goaccess /var/log/nginx/access.log -c`

Open port: `sudo ufw allow $PORT-WANTED`

```
server {
	listen 80;
	listen [::]:80;

	server_name subdomain.domain.com;
	location / {
		proxy_set_header X-Forwarded-For $remote_addr;
      		proxy_set_header Host            $http_host;
     		proxy_pass http://sub.domain.com:PORT/;
			#try_files $uri $uri/ =404;
	}
}

Test nginx congig: `sudo nginx -t`

```
### Create service
`sudo vim /lib/systemd/system/chall03.service`
```
[Unit]
Description=PoC

[Service]
Type=simple
Restart=always
RestartSec=5s
ExecStart=/home/chall03

[Install]
WantedBy=multi-user.target
```
