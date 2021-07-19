package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"gtoo/utils"

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

// Whois 查询whois
func Whois(domain string) error {
	resp, err := utils.Client.Get(fmt.Sprintf("https://api.devopsclub.cn/api/whoisquery?domain=%s&type=json", domain))
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
