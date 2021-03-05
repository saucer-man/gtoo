package domain

import (
	"fmt"
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
	fmt.Println(string(body))
	return nil
}
