# go-suricatasc
Unix Socket interaction with Suricata made in Go

## Usage
```
[go-suricatacs]# ./go-suricatacs -h
Usage of ./go-suricatacs:
  -interactive
    	Opens an interactive session to send commands to the socket
  -socket string
    	Full path to the suricata unix socket (default "/var/run/suricata/suricata-command.socket")

[go-suricatacs]# ./go-suricatacs -interactive=true
>> Entering Interactive Mode <<
>> Valid Commands:
{
    "commands": [
        "shutdown",
        "command-list",
        "help",
        "version",
        "uptime",
        "running-mode",
        "capture-mode",
        "conf-get",
        "dump-counters",
        "reload-rules",
        "register-tenant-handler",
        "unregister-tenant-handler",
        "register-tenant",
        "reload-tenant",
        "unregister-tenant",
        "pcap-file",
        "pcap-file-number",
        "pcap-file-list",
        "pcap-current"
    ],
    "count": 19
}
>>

[go-suricatacs]# ./go-suricatacs version
"3.0dev019f856"
```
