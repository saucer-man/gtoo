# gtoo
penetration tools based golang

渗透测试中的脚手架工具。不断更新中

todo
- 编码部分
    - [x] base64编码/解码
    - [x] md5

- 漏洞扫描
    - [x] 端口探测

- 信息收集
    - [x] whois查询/备案信息查询
    - [ ] 子域名收集(非爆破)

- 内网
    - [ ] 存活探测
    - [ ] MS17010


```
$ gtoo                    
gtoo id a pentese tool, when you have this, you will got all

Usage:
  gtoo [command]

Available Commands:
  convert     Common encoding and decoding
  help        Help about any command
  portscan    scan open port of host
  version     version of gtoo

Flags:
  -h, --help   help for gtoo

Use "gtoo [command] --help" for more information about a command.
```