package domain

import (
	"fmt"

	"github.com/tidwall/pretty"

	"gtoo/utils"
	"io/ioutil"
	"net/http"
	"strings"
)

type whoisResult struct {
	Code    int
	Message string
	Data    map[string]interface{}
}

// Whois 查询whois
func Whois(domain string) error {
	if !strings.HasPrefix(domain, "http") {
		domain = "https://" + domain
	}
	u, err := utils.ParseURL(domain)
	if err != nil {
		return err
	}
	resp, err := http.Get(fmt.Sprintf("https://api.devopsclub.cn/api/whoisquery?domain=%s.%s&type=json", u.Domain, u.TLD))
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Printf("1. whois查询结果:\n%s\n\n", pretty.Color(pretty.Pretty(body), pretty.TerminalStyle))
	resp, err = http.Get(fmt.Sprintf("https://api.66mz8.com/api/icp.php?domain=%s.%s", u.Domain, u.TLD))
	if err != nil {
		return err
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("2. 备案信息查询结果:\n%s\n\n", pretty.Color(pretty.Pretty(body), pretty.TerminalStyle))
	return nil
}
