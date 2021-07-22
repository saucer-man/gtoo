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
