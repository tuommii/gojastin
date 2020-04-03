# gojastin
Proof of Concept Backend for https://github.com/hivehelsinki/remote-challs/tree/master/chall03

WIP.

## Run
`make`

## Notes

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
