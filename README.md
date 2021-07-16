# gtoo
penetration tools based golang

渗透测试中的脚手架工具，基本原则是无任何依赖，开箱即用。

```bash
gtoo id a pentese tool, when you have this, you will got all

Usage:
  gtoo [command]

Available Commands:
  convert     常用编码解码
  domain      域名相关信息
  help        Help about any command
  ip          ip相关信息
  version     version of gtoo

Flags:
  -h, --help   help for gtoo

Use "gtoo [command] --help" for more information about a command.
```


### 1. 域名信息

```bash
$ gtoo domain info baidu.com
2021-07-14 19:06:01 [INFO] whois查询结果:
注册人: MarkMonitor Inc.
注册邮箱: abusecomplaints@markmonitor.com
注册商: whois.markmonitor.com
注册时间: 1999-10-11T11:05:17Z
到期时间: 2026-10-11T11:05:17Z
2021-07-14 19:06:01 [INFO] IPC查询结果:
域名: baidu.com
网站备案名称: 百度
主办单位性质: 企业
主办单位性质: 北京百度网讯科技有限公司
备案许可证号: 京ICP证030173号-1
最新审核检测: 2021-07-14 18:03:43
```

### 2. 子域名扫描

```
$ gtoo domain subdomain -d baidu.com  
2021-07-16 20:09:49 [INFO] 设置输出文件: /Users/yanq/Documents/self/gtoo/result.txt
2021-07-16 20:09:49 [INFO] 下面开始api扫描: [googlecert] [sublist3r] [chaziyu] [rapiddns] [threatminer] [bufferover] [crtsh] [certspotter] [mnemonic] [chinaz] [riddler]
2021-07-16 20:10:09 [INFO] api扫描结束!
```

### 3. ip信息

```bash
$ gtoo ip info 120.78.49.231 --threadbookapikey xxxxx
2021-07-14 19:09:41 [INFO] 微步在线查询IP信息...
IP: 120.78.49.231
情报可信度: 低
是否为恶意IP: true
IP危害等级: 低
IP威胁类型: [IDC服务器 扫描]
应用场景: 数据中心
IP安全事件0
    标签: [阿里云主机]
    类别: public_info
地理位置: 中国 广东省 深圳市
经纬度: lng：114.014495 lat：22.542702
运营商: 阿里云
情报更新时间: 2021-07-12 23:04:45
```