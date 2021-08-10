package domain

import (
	"encoding/json"
	"fmt"
	"gtoo/utils"
	"net/http"
	"strings"
)

type EntInfoResult struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data EntInfoData `json:"data"`
}

type EntInfoData struct {
	// IdStr   string `label:"网站备案名称"  json:"idStr"`
	CompanyName       string            `label:"企业名称" json:"companyName"`
	BusinessInfoModel BusinessInfoModel `label:"企业工商信息" json:"businessInfoModel"`
}
type BusinessInfoModel struct {
	Corporation        string `label:"法定代表人"  json:"corporation"`
	RegisteredCapital  string `label:"注册资本"  json:"registeredCapital"`
	RegistrationTime   string `label:"注册时间"  json:"registrationTime"`
	Type               string `label:"公司类型"  json:"type"`
	State              string `label:"公司状态"  json:"state"`
	RegistrationNumber string `label:"工商注册号"  json:"registrationNumber"`
	Industry           string `label:"所属行业"  json:"industry"`
	ApprovalDate       string `label:"核准日期"  json:"approvalDate"`
	RegisteredAddress  string `label:"注册地址"  json:"registeredAddress"`
	BusinessScope      string `label:"经营范围"  json:"businessScope"`
}

// 获取企业信息
func EntInfo(d string) (*EntInfoData, error) {

	var data = strings.NewReader("Kw=" + d)

	req, err := http.NewRequest("POST", "https://icp.chinaz.com/Home/QiYeData", data)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Host", "icp.chinaz.com")
	req.Header.Set("Connection", "close")
	req.Header.Set("Content-Length", "17")
	req.Header.Set("sec-ch-ua", "\" Not;A Brand\";v=\"99\", \"Google Chrome\";v=\"91\", \"Chromium\";v=\"91\"")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Origin", "https://icp.chinaz.com")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://icp.chinaz.com")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	resp, err := utils.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request请求失败：%v", err)
	}
	var entInfoResult EntInfoResult
	err = json.NewDecoder(resp.Body).Decode(&entInfoResult)
	if err != nil {
		return nil, fmt.Errorf("response解析失败：%v", err)
	}
	if entInfoResult.Data == (EntInfoData{}) {
		return nil, fmt.Errorf("未找到企业信息:%v", entInfoResult.Msg)
	}

	return &entInfoResult.Data, nil
}

type CompanyDomainResult struct {
	Amount       int                 `json:"amount"`
	PageNo       int                 `json:"pageNo"`
	PageSize     int                 `json:"pageSize"`
	PageLimit    int                 `json:"pageLimit"`
	PageAuthData string              `json:"pageAuthData"`
	ExportURL    string              `json:"exportUrl"`
	Code         int                 `json:"code"`
	Msg          string              `json:"msg"`
	Data         []CompanyDomainData `json:"data"`
}
type CompanyDomainData struct {
	Host       string `label:"域名" json:"host"`
	WebName    string `label:"公司名" json:"webName"`
	Owner      string `label:"所有人" json:"owner"`
	Permit     string `label:"备案号" json:"permit"`
	Typ        string `label:"类型" json:"typ"`
	VerifyTime string `label:"验证时间" json:"verifyTime"`
}

func GetCompanyDomain(companyName string) ([]CompanyDomainData, error) {
	var res []CompanyDomainData
	pageNo := 1
	for {
		var data = strings.NewReader(fmt.Sprintf("pageNo=%d&pageSize=20&Kw=%s", pageNo, companyName))
		req, err := http.NewRequest("POST", "http://icp.chinaz.com/Home/PageData", data)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
		req.Header.Set("X-Requested-With", "XMLHttpRequest")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		req.Header.Set("Origin", "http://icp.chinaz.com")
		req.Header.Set("Referer", "http://icp.chinaz.com/tuya.com")
		req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
		req.Header.Set("Cookie", "UM_distinctid=17a3cd265e5c6a-0d501c25d6ff1a-34647600-113a00-17a3cd265e6b3d; toolbox_words=103.202.147.199; qHistory=aHR0cDovL3Rvb2wuY2hpbmF6LmNvbS90b29scy91bmljb2RlLmFzcHhfVW5pY29kZee8lueggei9rOaNonxzYW1lL1/lkIxJUOe9keermeafpeivonxodHRwOi8vaXAudG9vbC5jaGluYXouY29tL3NpdGVpcC9fSVDmiYDlnKjlnLDmibnph4/mn6Xor6J8aHR0cDovL2lwLnRvb2wuY2hpbmF6LmNvbS9pcGJhdGNoL19JUOaJuemHj+afpeivonxodHRwOi8vaXAudG9vbC5jaGluYXouY29tX0lQL0lQdjbmn6Xor6LvvIzmnI3liqHlmajlnLDlnYDmn6Xor6J8aHR0cDovL3Rvb2wuY2hpbmF6LmNvbS9pcHdob2lzL19JUCBXSE9JU+afpeivonxodHRwOi8vdG9vbC5jaGluYXouY29tX+ermemVv+W3peWFt3xodHRwOi8vd2hvaXMuY2hpbmF6LmNvbS9yZXZlcnNlP2RkbFNlYXJjaE1vZGU9Ml/ms6jlhozkurrlj43mn6V8aHR0cDovL3dob2lzLmNoaW5hei5jb20vX1dob2lz5p+l6K+ifGh0dHA6Ly90b29sLmNoaW5hei5jb20vdG9vbHMvaW1ndG9iYXNlL19iYXNlNjTlm77niYflnKjnur/ovazmjaLlt6Xlhbc=; Hm_lvt_aecc9715b0f5d5f7f34fba48a3c511d6=1628596403; Hm_lpvt_aecc9715b0f5d5f7f34fba48a3c511d6=1628596403; Hm_lvt_ca96c3507ee04e182fb6d097cb2a1a4c=1626233319,1626688418,1627527262,1628596918; CNZZDATA5082706=cnzz_eid%3D585403542-1626683320-https%253A%252F%252Fwww.google.com%252F%26ntime%3D1628594875; Hm_lpvt_ca96c3507ee04e182fb6d097cb2a1a4c=1628597612")
		resp, err := utils.Client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		var companyDomainResult CompanyDomainResult
		err = json.NewDecoder(resp.Body).Decode(&companyDomainResult)
		if err != nil {
			break
		}
		if len(companyDomainResult.Data) == 0 {
			break
		}
		res = append(res, companyDomainResult.Data...)
		pageNo = pageNo + 1
	}
	return res, nil
}
