# gtoo
penetration tools based golang

渗透测试中的脚手架工具，基本原则是无任何依赖，开箱即用。

```bash
$ gtoo                              
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
  -h, --help      help for gtoo
  -v, --verbose   verbose output

Use "gtoo [command] --help" for more information about a command.

```


### 1. 域名信息

```bash
$ gtoo domain info www.baidu.com     
2021-07-28 20:32:49 [INFO] 解析出域名:baidu.com

2021-07-28 20:32:50 [INFO] whois查询结果:
注册人: MarkMonitor Inc.
注册邮箱: abusecomplaints@markmonitor.com
注册商: whois.markmonitor.com
注册时间: 1999-10-11T11:05:17Z
到期时间: 2026-10-11T11:05:17Z
2021-07-28 20:32:50 [INFO] ICP备案结果:
备案许可证号: 京ICP证030173号-1
网站备案名称: 北京百度网讯科技有限公司
主办单位性质: 企业
审核时间: 2021-06-10
2021-07-28 20:32:50 [INFO] 企业工商信息：
企业名称: 北京百度网讯科技有限公司
法定代表人: 梁志祥
注册资本: 1,342,128万(元)
注册时间: 2001-06-05
公司类型: 有限责任公司(自然人投资或控股)
公司状态: 开业
工商注册号: 110108002734659
所属行业: 科技推广和应用服务业
核准日期: 2021-01-26
注册地址: 北京市海淀区上地十街10号百度大厦2层
经营范围: 技术转让、技术咨询、技术服务、技术培训、技术推广;设计、开发、销售计算机软件;经济信息咨询;利用www.baidu.com、www.hao123.com(www.hao222.net、www.hao222.com)网站发布广告;设计、制作、代理、发布广告;货物进出口、技术进出口、代理进出口;医疗软件技术开发;委托生产电子产品、玩具、照相器材;销售家用电器、机械设备、五金交电(不含电动自行车)、电子产品、文化用品、照相器材、计算机、软件及辅助设备、化妆品、卫生用品、体育用品、纺织品、服装、鞋帽、日用品、家具、首饰、避孕器具、工艺品、钟表、眼镜、玩具、汽车及摩托车配件、仪器仪表、塑料制品、花、草及观赏植物、建筑材料、通讯设备、汽车电子产品、器件和元件、自行开发后的产品;预防保健咨询;公园门票、文艺演出、体育赛事、展览会票务代理;翻译服务;通讯设备和电子产品的技术开发;计算机系统服务;车联网技术开发;汽车电子产品设计、研发、制造(北京市中心城区除外);演出经纪;人才中介服;经营电信业务;利用信息网络经营音乐娱乐产品、演出剧(节)目、动漫产品、游戏产品(含网络游戏虚拟货币发行)、表演、网络游戏技法展示或解说(网络文化经营许可证有效期至2020年04月17日);因特网信息服务业务(除出版、教育、医疗保健以外的内容);图书、电子出版物、音像制品批发、零售、网上销售。(市场主体依法自主选择经营项目,开展经营活动;演出经纪、人才中介服务、利用信息网络经营音乐娱乐产品、演出剧(节)目、动漫产品、游戏产品(含网络游戏虚拟货币发行)、表演、网络游戏技法展示或解说、经营电信业务以及依法须经批准的项目,经相关部门批准后依批准的内容开展经营活动;不得从事国家和本市产业政策禁止和限制类项目的经营活动。)
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
$ gtoo ip info 120.78.49.231 --threadbookapikey xxx
2021-07-28 20:33:57 [INFO] 微步在线查询IP信息:
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
情报更新时间: 2021-07-26 23:43:34
2021-07-28 20:33:57 [INFO] ip138查询ip信息:
ASN归属地: 广东省深圳市  阿里云 数据中心
iP段: 120.76.0.0 - 120.79.255.255
兼容IPv6地址: ::784E:31E7
映射IPv6地址: ::FFFF:784E:31E7
2021-07-28 20:33:58 [INFO] ip反查域名:
2021-07-28 20:33:59 [INFO] ip反查域名版本2:
lc.59cl.cn
saucer-man.com
www.furongdo.com
xiaogeng.top
```