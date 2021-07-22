package domain

import (
    "bufio"
    "fmt"
    "gtoo/utils"
    "io/ioutil"
    "net"
    "net/http"
    "os"
    "regexp"
    "strings"

    log "github.com/sirupsen/logrus"
)

var SubdomainApi = map[string]string{
    // 证书透明度
    "crtsh":       "https://crt.sh/?output=json&q=%s",
    "certspotter": "https://api.certspotter.com/v1/issuances?domain=%s&include_subdomains=true&expand=dns_names",
    "googlecert":  "https://transparencyreport.google.com/transparencyreport/api/v3/httpsreport/ct/certsearch?include_expired=true&include_subdomains=true&domain=%s",

    // DNS数据集
    "mnemonic":  "https://api.mnemonic.no/pdns/v3/%s?limit=1000",
    "sublist3r": "https://api.sublist3r.com/search.php?domain=%s",
    "chaziyu":   "https://chaziyu.com/%s/",
    "chinaz":    "https://alexa.chinaz.com/%s",
    "rapiddns":  "https://rapiddns.io/subdomain/%s?full=1",
    "riddler":   "https://riddler.io/search?q=pld:%s",
    // "sitedossier": "http://www.sitedossier.com/parentdomain/%s",
    "bufferover":  "https://dns.bufferover.run/dns?q=%s",
    "threatminer": "https://api.threatminer.org/v2/domain.php?q=%s",
    "dnsgrep":     "https://www.dnsgrep.cn/subdomain/%s",
}

// IsWildCard 检测是否是泛解析域名
func IsWildCard(d string) bool {
    ranges := [2]int{}
    for _, _ = range ranges {
        subdomain := utils.RandomStr(6) + "." + d
        _, err := net.LookupIP(subdomain)
        if err != nil {
            continue
        }
        return true
    }
    return false
}

// 目前不解析根域名
func MatchSubdomainsFromText(d, text string) []string {
    // [A-Z] 表示一个区间
    // {n}    n 是一个非负整数。匹配确定的 n 次
    SUBRE := "(([a-zA-Z0-9]{1}|[_a-zA-Z0-9]{1}[_a-zA-Z0-9-]{0,61}[a-zA-Z0-9]{1})[.]{1})+" + strings.Replace(d, ".", "[.]", -1)
    // log.Info(SUBRE)

    r := regexp.MustCompile(SUBRE)
    return utils.RemoveRep(r.FindAllString(text, -1))

}

// 访问url，获取subdomain
func GetSubdomainFromUrl(d, url string) ([]string, error) {
    req, _ := http.NewRequest("GET", url, nil)
    req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
    req.Close = true
    req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
    req.Header.Set("Accept-Encoding", "gzip, deflate")
    req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")

    resp, err := utils.Client.Do(req)
    if err != nil {
        return []string{}, fmt.Errorf("发送http请求失败：%v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return []string{}, fmt.Errorf("response code：%d", resp.StatusCode)
    }
    bodyBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return []string{}, fmt.Errorf("读取http res失败：%v", err)
    }
    // log.Info(string(bodyBytes))
    return MatchSubdomainsFromText(d, string(bodyBytes)), nil
}

func ApiScan(domains []string, outputPath string) error {
    taskCreater := make(chan int, 1)   //一个任务生产者
    taskConsumer := make(chan int, 20) // 20个任务消费者
    resConsumer := make(chan int, 1)   //一个结果消费者，
    // 扫描任务 channel、每个元素是一个map
    task := make(chan map[string]string, 1000)

    // 扫描结果
    var resSet = make(map[string]bool)
    var resChannel = make(chan string, 10000)

    // 生产任务
    go func(task chan<- map[string]string, taskCreater chan<- int) {
        for _, d := range domains {
            for api := range SubdomainApi {
                task <- map[string]string{
                    "domain": d,
                    "api":    api,
                }
            }
        }
        log.Debug("生产者关闭了chanel")
        close(task)
        taskCreater <- 1
    }(task, taskCreater)

    // 消费任务、生产结果
    for i := 0; i < cap(taskConsumer); i++ {
        go func(tasks <-chan map[string]string, done chan<- int, i int) {
            for task := range tasks {
                apiName := task["api"]
                apiUrl := SubdomainApi[apiName]
                url := fmt.Sprintf(apiUrl, task["domain"])
                subdomain, err := GetSubdomainFromUrl(task["domain"], url)
                if err != nil {
                    log.Warnf("[%s]发生错误:%v", apiName, err)
                }
                // 将结果保存在resSet、resChannel中，并且检查是否重复
                for _, d := range subdomain {
                    if d == "" {
                        continue
                    }
                    _, ok := resSet[d]
                    if ok {
                        continue
                    }
                    resSet[d] = true
                    resChannel <- d
                }
                // fmt.Printf("[%s] %v\n", apiName, subdomain)
            }
            log.Debugf("消费者%d完成了任务", i)
            done <- 1
        }(task, taskConsumer, i)
    }
    // 消费结果，保存到文件中
    go func(resConsumer chan<- int) {
        f, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
        if err != nil {
            log.Fatalf("%s文件打开失败: %v", outputPath, err)
        }
        defer f.Close()
        //写入文件时，使用带缓存的 *Writer
        write := bufio.NewWriter(f)
        for res := range resChannel {
            write.WriteString(res + "\n")
        }
        //Flush将缓存的文件真正写入到文件中
        write.Flush()
        resConsumer <- 1
    }(resConsumer)
    // 等待任务生产完成
    <-taskCreater
    // 等待任务消费完成
    for i := 0; i < cap(taskConsumer); i++ {
        <-taskConsumer
    }
    close(resChannel)

    // 等待结果消费完成
    <-resConsumer

    return nil
}
