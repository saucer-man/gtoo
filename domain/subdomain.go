package domain

import (
	"fmt"
	"gtoo/utils"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

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
func MatchSubdomains(d, text string) []string {
	// [A-Z] 表示一个区间
	// {n}	n 是一个非负整数。匹配确定的 n 次
	SUBRE := "(([a-zA-Z0-9]{1}|[_a-zA-Z0-9]{1}[_a-zA-Z0-9-]{0,61}[a-zA-Z0-9]{1})[.]{1})+"
	// log.Info(text)
	r := regexp.MustCompile(SUBRE + strings.Replace(d, ".", "[.]", -1))
	return utils.RemoveRep(r.FindAllString(text, -1))

}

// TODO 重写
func Bufferover(d string) error {
	url := "https://dns.bufferover.run/dns?q=" + d
	sub, err := GetSubdomain(d, url)
	if err != nil {
		return err
	}

	for _, u := range sub {
		log.Info(u)
	}
	return nil
}

func GetSubdomain(d, url string) ([]string, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Close = true

	resp, err := utils.Client.Do(req)
	if err != nil {
		return []string{}, fmt.Errorf("发送http请求失败：%v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []string{}, fmt.Errorf("response非200：%d", resp.StatusCode)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []string{}, fmt.Errorf("读取http res失败：%v", err)
	}
	return MatchSubdomains(d, string(bodyBytes)), nil
}

func Cersh(d string) error {
	url := "https://crt.sh/?output=json&q=" + d
	sub, err := GetSubdomain(d, url)
	if err != nil {
		return err
	}

	for _, u := range sub {
		fmt.Printf("[cer.sh] %s\n", u)
	}
	return nil
}

// TODO 重写
func Mnemonic(d string) error {

	url := "https://api.mnemonic.no/pdns/v3/" + d
	sub, err := GetSubdomain(d, url)
	if err != nil {
		return err
	}

	for _, u := range sub {
		fmt.Printf("[mnemonic] %s\n", u)
	}
	return nil
}

func Sublist3r(d string) error {

	url := "https://api.sublist3r.com/search.php?domain=" + d
	sub, err := GetSubdomain(d, url)
	if err != nil {
		return err
	}

	for _, u := range sub {
		fmt.Printf("[sublist3r] %s\n", u)
	}
	return nil
}

func Chaziyu(d string) error {

	url := "https://chaziyu.com/" + d + "/"
	sub, err := GetSubdomain(d, url)
	if err != nil {
		return err
	}

	for _, u := range sub {
		fmt.Printf("[chaziyu] %s\n", u)
	}
	return nil
}

func Chinaz(d string) error {

	url := "https://alexa.chinaz.com/" + d
	sub, err := GetSubdomain(d, url)
	if err != nil {
		return err
	}

	for _, u := range sub {
		fmt.Printf("[chinaz] %s\n", u)
	}
	return nil
}

func Rapiddns(d string) error {

	url := "https://rapiddns.io/subdomain/" + d + "?full=1"
	sub, err := GetSubdomain(d, url)
	if err != nil {
		return err
	}

	for _, u := range sub {
		fmt.Printf("[rapiddns] %s\n", u)
	}
	return nil
}
func Riddler(d string) error {

	url := "https://riddler.io/search?q=pld:" + d
	sub, err := GetSubdomain(d, url)
	if err != nil {
		return err
	}

	for _, u := range sub {
		fmt.Printf("[riddler] %s\n", u)
	}
	return nil
}
func Sitedossier(d string) error {
	url := "http://www.sitedossier.com/parentdomain/" + d
	sub, err := GetSubdomain(d, url)
	if err != nil {
		return err
	}

	for _, u := range sub {
		fmt.Printf("[sitedossier] %s\n", u)
	}
	return nil
}
func Threatminer(d string) error {

	url := "https://api.threatminer.org/v2/domain.php?q=" + d
	sub, err := GetSubdomain(d, url)
	if err != nil {
		return err
	}

	for _, u := range sub {
		fmt.Printf("[threatminer] %s\n", u)
	}
	return nil
}
