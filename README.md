# gojastin
Proof of Concept Backend for https://github.com/hivehelsinki/remote-challs/tree/master/chall03

Challenge was fun but server side seemed far more interesting so I decided try to do my own.

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
