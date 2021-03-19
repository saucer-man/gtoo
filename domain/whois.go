package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"gtoo/utils"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type whoisResult struct {
	Code int       `json:"code"`
	Data WhoisData `json:"data"`
	Msg  string    `json:"msg"`
}

type WhoisData struct {
	Data   Data `json:"data"`
	Status int  `json:"status"`
}
type Data struct {
	Registrar                  string `json:"registrar"`
	Registrant                 string `json:"registrant"`
	RegistrarAbuseContactEmail string `json:"registrarAbuseContactEmail"`
	RegistrantContactEmail     string `json:"registrantContactEmail"`
	RegistrarWHOISServer       string `json:"registrarWHOISServer"`
	SponsoringRegistrar        string `json:"sponsoringRegistrar"`
	CreationDate               string `json:"creationDate"`
	RegistrationTime           string `json:"registrationTime"`
	RegistryExpiryDate         string `json:"registryExpiryDate"`
	ExpirationTime             string `json:"rxpirationTime"`
}

var myClient = &http.Client{Timeout: 10 * time.Second}

// Whois 查询whois
func Whois(domain string) error {

	u, err := utils.ParseURL(domain)
	if err != nil {
		return err
	}
	resp, err := myClient.Get(fmt.Sprintf("https://api.devopsclub.cn/api/whoisquery?domain=%s.%s&type=json", u.Domain, u.TLD))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var whoisres whoisResult
	err = json.NewDecoder(resp.Body).Decode(&whoisres)
	if err != nil {
		return fmt.Errorf("返回结果解析失败：%v", err)
	}
	if whoisres.Code != 0 {
		return errors.New(whoisres.Msg)
	}
	log.Info("whois查询结果:")
	if whoisres.Data.Data.Registrar != "" {
		fmt.Printf("注册人: %s\n", whoisres.Data.Data.Registrar)
	} else if whoisres.Data.Data.Registrant != "" {
		fmt.Printf("注册人: %s\n", whoisres.Data.Data.Registrant)
	} else {
		fmt.Println("注册人: N/A")
	}
	if whoisres.Data.Data.RegistrarAbuseContactEmail != "" {
		fmt.Printf("注册邮箱: %s\n", whoisres.Data.Data.RegistrarAbuseContactEmail)
	} else if whoisres.Data.Data.RegistrantContactEmail != "" {
		fmt.Printf("注册邮箱: %s\n", whoisres.Data.Data.RegistrantContactEmail)
	} else {
		fmt.Println("注册邮箱: N/A")
	}
	if whoisres.Data.Data.RegistrarWHOISServer != "" {
		fmt.Printf("注册商: %s\n", whoisres.Data.Data.RegistrarWHOISServer)
	} else if whoisres.Data.Data.SponsoringRegistrar != "" {
		fmt.Printf("注册商: %s\n", whoisres.Data.Data.SponsoringRegistrar)
	} else {
		fmt.Println("注册商: N/A")
	}
	if whoisres.Data.Data.CreationDate != "" {
		fmt.Printf("注册时间: %s\n", whoisres.Data.Data.CreationDate)
	} else if whoisres.Data.Data.RegistrationTime != "" {
		fmt.Printf("注册时间: %s\n", whoisres.Data.Data.RegistrationTime)
	} else {
		fmt.Println("注册时间: N/A")
	}
	if whoisres.Data.Data.RegistryExpiryDate != "" {
		fmt.Printf("到期时间: %s\n", whoisres.Data.Data.RegistryExpiryDate)
	} else if whoisres.Data.Data.ExpirationTime != "" {
		fmt.Printf("到期时间: %s\n", whoisres.Data.Data.ExpirationTime)
	} else {
		fmt.Println("到期时间: N/A")
	}
	return nil
}

type ipcResult struct {
	Success bool    `json:"success"`
	Domain  string  `json:"domain"`
	Info    ipcInfo `json:"info"`
	Message string  `json:"message"`
}

type ipcInfo struct {
	Name   string `json:"name"`
	Nature string `json:"nature"`
	Icp    string `json:"icp"`
	Title  string `json:"title"`
	Time   string `json:"time"`
}

// IPC备案查询
func Ipc(domain string) error {
	u, err := utils.ParseURL(domain)
	if err != nil {
		return err
	}
	// https://api.66mz8.com/api/icp.php?domain=example.com
	resp, err := myClient.Get(fmt.Sprintf("https://api.vvhan.com/api/icp?url=%s.%s", u.Domain, u.TLD))
	if err != nil {
		return err
	}
	var ipcres ipcResult
	err = json.NewDecoder(resp.Body).Decode(&ipcres)
	if err != nil {
		return fmt.Errorf("response解析失败：%v", err)
	}
	if !ipcres.Success {
		return errors.New("备案信息查找失败")
	}
	if ipcres.Message != "" {
		log.Infof("IPC查询结果:%s", ipcres.Message)
		return nil
	}
	log.Info("IPC查询结果:")
	fmt.Printf("域名: %s\n", ipcres.Domain)
	fmt.Printf("网站备案名称: %s\n", ipcres.Info.Title)
	fmt.Printf("主办单位性质: %s\n", ipcres.Info.Nature)
	fmt.Printf("主办单位性质: %s\n", ipcres.Info.Name)
	fmt.Printf("备案许可证号: %s\n", ipcres.Info.Icp)
	fmt.Printf("最新审核检测: %s\n", ipcres.Info.Time)
	return nil
}
