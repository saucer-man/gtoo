package zoomeye


package zoomeye

import (
	"fmt"
	"bytes"
	"strconv"
	"io/ioutil"
	"net/http"
	"encoding/json"
)

const (
        defaultBaseURL = "https://api.zoomeye.org/"
        getAccessTokenURL = "user/login"
        searchHostURL = "host/search"
        searchWebURL = "web/search"
        resourcesInfoURL = "resources-info"
)

var AccessToken string
var userName string
var passWord string

type Auth struct {
	AccessToken  string `json:"access_token"`
}


type Portinfo struct {
	Port     int   `json:"port"`
}

type Matche struct {
	Ip       string     `json:"ip"`
	Portinfo Portinfo `json:"portinfo"`
}

type Matches struct {
	Matches  []Matche `json:"matches"`
	Total	 int      `json:"total"`
}



func Get(urlStr string) []byte {
	request, err := http.NewRequest("GET", urlStr, nil)
	request.Header.Set("Authorization", "JWT " + AccessToken)
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	return respBytes
}

func Post(urlStr string, data interface{}) []byte {
	bytesData, err := json.Marshal(data)
	body := bytes.NewReader(bytesData)
	request, err := http.NewRequest("POST", urlStr, body)
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	return respBytes
}

func GetAccessToken() string{
	data := make(map[string]interface{})
	data["username"] = userName
	data["password"] = passWord

	auth := Auth{}

	err :=  json.Unmarshal(Post(defaultBaseURL + getAccessTokenURL,data),&auth)
	if err!=nil{
		fmt.Println(err)
	}
	return auth.AccessToken
}

func Search(query string,page int,username string,password string) []string {
	userName = username
	passWord = password
	AccessToken = GetAccessToken()
	resp := Get(defaultBaseURL + searchHostURL + "?query=" + query + "&page=" + strconv.Itoa(page))
	matches := Matches{}
	json.Unmarshal(resp, &matches)
	fmt.Println("Total:",matches.Total)
	var result []string
	for _,line := range matches.Matches {
		result = append(result,line.Ip + ":" + strconv.Itoa(line.Portinfo.Port))
	}
	return result
}