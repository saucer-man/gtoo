package domain

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// 获取企业信息
func EntInfo(d string) error {

	var data = strings.NewReader(`$Kw=saucer-man.com`)
	req, err := http.NewRequest("POST", "$https://icp.chinaz.com/Home/QiYeData", data)
	if err != nil {
		return err
	}
	req.Header.Set("$Host", "icp.chinaz.com")
	req.Header.Set("$Connection", "close")
	req.Header.Set("$Content-Length", "17")
	req.Header.Set("$sec-ch-ua", "\" Not;A Brand\";v=\"99\", \"Google Chrome\";v=\"91\", \"Chromium\";v=\"91\"")
	req.Header.Set("$Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("$X-Requested-With", "XMLHttpRequest")
	req.Header.Set("$sec-ch-ua-mobile", "?0")
	req.Header.Set("$User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36")
	req.Header.Set("$Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("$Origin", "https://icp.chinaz.com")
	req.Header.Set("$Sec-Fetch-Site", "same-origin")
	req.Header.Set("$Sec-Fetch-Mode", "cors")
	req.Header.Set("$Sec-Fetch-Dest", "empty")
	req.Header.Set("$Referer", "https://icp.chinaz.com/saucer-man.com")
	req.Header.Set("$Accept-Encoding", "gzip, deflate")
	req.Header.Set("$Accept-Language", "zh-CN,zh;q=0.9")
	resp, err := utils.client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", bodyText)
}
