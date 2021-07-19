package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"gtoo/utils"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type IcpVvhanResult struct {
	Success bool         `json:"success"`
	Domain  string       `json:"domain"`
	Info    IcpVvhanInfo `json:"info"`
	Message string       `json:"message"`
}

type IcpVvhanInfo struct {
	Name   string `label:"网站备案名称"  json:"name"`
	Nature string `label:"主办单位性质" json:"nature"`
	Icp    string `label:"备案许可证号" json:"icp"`
	Title  string `label:"主办单位名称" json:"title"`
	Time   string `label:"最新审核检测" json:"time"`
}

// IPC备案查询 by vvhan
func IcpVvhan(domain string) (*IcpVvhanInfo, error) {
	// https://api.66mz8.com/api/icp.php?domain=example.com
	resp, err := myClient.Get(fmt.Sprintf("https://api.vvhan.com/api/icp?url=%s", domain))
	if err != nil {
		return nil, err
	}
	var ipcres IcpVvhanResult
	err = json.NewDecoder(resp.Body).Decode(&ipcres)
	if err != nil {
		return nil, fmt.Errorf("response解析失败：%v", err)
	}
	if !ipcres.Success {
		return nil, errors.New("备案信息查找失败")
	}
	if ipcres.Message != "" {
		return nil, errors.New(ipcres.Message)
	}

	return &ipcres.Info, nil
}

type IcpChinazResult struct {
	IcpNumber string `label:"备案许可证号" json:"icp_number"`
	IcpName   string `label:"网站备案名称" json:"icp_name"`
	Attr      string `label:"主办单位性质" json:"attr"`
	Date      string `label:"审核时间" json:"date"`
}

func IcpChinaz(url string) (*IcpChinazResult, error) {
	url = "https://icp.chinaz.com/" + url
	icp := new(IcpChinazResult)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return icp, err
	}
	request.Header.Add("user-agent", utils.RandomUserAgent())

	resp, err := utils.Client.Do(request)
	if err != nil {
		return icp, err
	}
	defer resp.Body.Close()
	gp, err := goquery.NewDocumentFromReader(resp.Body)
	// gp, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return icp, err
	}

	gp.Find("#first > li").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			icp.IcpName = strings.TrimSpace(s.Find("p").Text())
		}

		if i == 1 {
			icp.Attr = strings.TrimSpace(s.Find("p").Text())
		}

		if i == 2 {
			icp.IcpNumber = strings.TrimSpace(s.Find("p > font").Text())
		}

		if i == 6 {
			icp.Date = strings.TrimSpace(s.Find("p").Text())
		}

	})

	if icp.IcpName == "" {
		return icp, fmt.Errorf("没有查询到备案信息")
	}
	return icp, nil
}
