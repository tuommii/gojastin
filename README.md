# gojastin
Proof of Concept Backend for https://github.com/hivehelsinki/remote-challs/tree/master/chall03

Challenge was fun but server side seemed far more interesting so I decided try to do my own.

WIP.

## Run
`make`

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
