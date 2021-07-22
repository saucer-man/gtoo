package ip

import (
	"encoding/json"
	"fmt"
	"gtoo/utils"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/axgle/mahonia"
	log "github.com/sirupsen/logrus"
)

type ThreatBook struct {
	apiKey string
}

func NewThreatBook(apiKey string) *ThreatBook {
	return &ThreatBook{
		apiKey: apiKey,
	}
}

type ThreatBookIPResult struct {
	ResponseCode int                         `json:"response_code"`
	VerboseMsg   string                      `json:"verbose_msg"`
	Data         map[string]ThreatBookIPData `json:"data"`
}

type ThreatBookIPData struct {
	Scene           string             `json:"scene"`
	UpdateTime      string             `json:"update_time"`
	ConfidenceLevel string             `json:"confidence_level"`
	Severity        string             `json:"severity"`
	IsMalicious     bool               `json:"is_malicious"`
	Judgments       []string           `json:"judgments"`
	TagsClasses     []ThreatBookIPTags `json:"tags_classes"`
	Basic           ThreatBookIPBasic  `json:"basic"`
	Asn             ThreatBookIPAsn    `json:"asn"`
}
type ThreatBookIPTags struct {
	Tags     []string `json:"tags"`
	TagsType string   `json:"tags_type"`
}

type ThreatBookIPBasic struct {
	Location ThreatBookIPLocation `json:"location"`
	Carrier  string               `json:"carrier"`
}

type ThreatBookIPLocation struct {
	Country     string `json:"country"`
	Province    string `json:"province"`
	City        string `json:"city"`
	Lng         string `json:"lng"`
	Lat         string `json:"lat"`
	CountryCode string `json:"country_code"`
}

type ThreatBookIPAsn struct {
	Rank   int    `json:"rank"`
	Info   string `json:"info"`
	Number int    `json:"number"`
}

// IP信誉查询
func (t *ThreatBook) IpReputation(ip string) error {
	url := fmt.Sprintf("https://api.threatbook.cn/v3/scene/ip_reputation?apikey=%s&resource=%s&lang=zh", t.apiKey, ip)

	resp, err := utils.Client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var threatBookIPResult ThreatBookIPResult
	err = json.NewDecoder(resp.Body).Decode(&threatBookIPResult)
	if err != nil {
		return fmt.Errorf("返回结果解析失败：%v", err)
	}
	if threatBookIPResult.ResponseCode != 0 {
		if threatBookIPResult.VerboseMsg == "Beyond Daily Limitation" {
			log.Info("微步 API 已超出当日使用次数")
		} else {
			log.Infof("微步 API 调用失败，错误信息:%s", threatBookIPResult.VerboseMsg)
		}
		return nil
	}

	for k, v := range threatBookIPResult.Data {
		// 情报可信度
		fmt.Printf("IP: %s\n", k)
		fmt.Printf("情报可信度: %s\n", v.ConfidenceLevel)
		fmt.Printf("是否为恶意IP: %v\n", v.IsMalicious)
		fmt.Printf("IP危害等级: %s\n", v.Severity)
		fmt.Printf("IP威胁类型: %+v\n", v.Judgments)
		fmt.Printf("应用场景: %s\n", v.Scene)
		for i, tag := range v.TagsClasses {
			fmt.Printf("IP安全事件%d\n", i)
			fmt.Printf("    标签: %+v\n", tag.Tags)
			fmt.Printf("    类别: %+v\n", tag.TagsType)
		}

		fmt.Printf("地理位置: %s %s %s\n", v.Basic.Location.Country, v.Basic.Location.Province, v.Basic.Location.City)
		fmt.Printf("经纬度: lng：%s lat：%s\n", v.Basic.Location.Lng, v.Basic.Location.Lat)
		fmt.Printf("运营商: %s\n", v.Basic.Carrier)
		fmt.Printf("情报更新时间: %s\n", v.UpdateTime)
	}
	return nil
}

func Ipinfo(ip string) error {
	url := fmt.Sprintf("https://www.ip138.com/iplookup.asp?ip=%s&action=2", ip)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.164 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cookie", "PHPSESSID=50vhfi98p3bm0uorm9b2kh99dk; Hm_lvt_f4f76646cd877e538aa1fbbdf351c548=1624444428,1626868959; ASPSESSIONIDCARCCCBT=BNANEOGCGDELKMMOJJJMPGMN; ASPSESSIONIDQATRTSDA=EJEOPIGCCIALGHNDGLPPKEMA; Hm_lpvt_f4f76646cd877e538aa1fbbdf351c548=1626869400")
	resp, err := utils.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	// 解决gbk编码 中文乱码
	bodystr := mahonia.NewDecoder("gbk").ConvertString(string(body))

	flysnowRegexp := regexp.MustCompile(`var ip_result = ([\s\S]*?);[\s\S]*?var ip_begin`)

	params := flysnowRegexp.FindStringSubmatch(bodystr)

	m := make(map[string]string)
	json.Unmarshal([]byte(params[1]), &m)
	for k, v := range m {
		if v != "" {
			fmt.Printf("%v: %+v\n", k, v)
		}
	}
	return nil

}
