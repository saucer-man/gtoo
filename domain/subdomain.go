package domain

import (
	"encoding/json"
	"fmt"
	"gtoo/utils"
	"net"
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
func MatchSubdomains(d, text string) ([]string, error) {
	// [A-Z] 表示一个区间
	// {n}	n 是一个非负整数。匹配确定的 n 次
	SUBRE := "(([a-zA-Z0-9]{1}|[_a-zA-Z0-9]{1}[_a-zA-Z0-9-]{0,61}[a-zA-Z0-9]{1})[.]{1})+"
	// SUBRE + strings.Replace(d, ".", "[.]", -1)
	log.Info(SUBRE + strings.Replace(d, ".", "[.]", -1))
	r := regexp.MustCompile(SUBRE + strings.Replace(d, ".", "[.]", -1))
	return r.FindAllString(text, -1), nil

}

//     result = re.findall(regexp, html, re.I)
//     if not result:
//         return set()
//     deal = map(lambda s: s.lower(), result)
//     if distinct:
//         return set(deal)
//     else:
//         return list(deal)
// else:
//     regexp = r'(?:\>|\"|\'|\=|\,)(?:http\:\/\/|https\:\/\/)?' \
//              r'(?:[a-z0-9](?:[a-z0-9\-]{0,61}[a-z0-9])?\.){0,}' \
//              + domain.replace('.', r'\.')
//     result = re.findall(regexp, html, re.I)
// if not result:
//     return set()
// regexp = r'(?:http://|https://)'
// deal = map(lambda s: re.sub(regexp, '', s[1:].lower()), result)
// if distinct:
//     return set(deal)
// else:
//     return list(deal)

type BufferoverResp struct {
	FDNSA []string `json:"FDNS_A"`
	RDNS  []string `json:"RDNS"`
}

func Bufferover(d string) error {
	resp, err := utils.Client.Get("https://dns.bufferover.run/dns?q=" + d)
	if err != nil {
		return fmt.Errorf("构造http请求失败：%v", err)
	}
	defer resp.Body.Close()
	var bufferoverResp BufferoverResp
	err = json.NewDecoder(resp.Body).Decode(&bufferoverResp)
	if err != nil {
		return fmt.Errorf("返回结果解析失败：%v", err)
	}
	log.Info("bufferover解析结果:")
	if len(bufferoverResp.FDNSA) > 0 {
		for _, l := range bufferoverResp.FDNSA {
			for _, t := range strings.Split(l, ",") {
				fmt.Printf("%s   ", t)
			}
			fmt.Printf("\n")
		}
	}
	if len(bufferoverResp.RDNS) > 0 {
		for _, l := range bufferoverResp.RDNS {
			for _, t := range strings.Split(l, ",") {
				fmt.Printf("%s   ", t)
			}
			fmt.Printf("\n")
		}
	}
	return nil
}

// // 证书透明度
// func Cetsh(d string) error {

// }
