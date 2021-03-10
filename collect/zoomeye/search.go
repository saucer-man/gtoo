package zoomeye

// from https://github.com/gyyyy/ZoomEye-go

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Search 根据dork查询目标
// 其中resource代表type host/web
func Search(apikey, dork, resource string, num int, filename string) error {
	zoom := NewWithKey(apikey)

	// 查询用户资源信息
	info, err := zoom.ResourcesInfo()
	if err != nil {
		log.Errorf("Error when get id info, maybe because of incorrect api key: %s", apikey)
		return err
	}
	log.Infof("Successfully login zoomeye. role: %s, quota: %d.", info.Plan, info.Resources.Search)

	// 下面开始查询
	if num <= 0 {
		num = 20
	}
	log.Infof("Set search num %d.", num)
	maxPage := num / 20
	if num%20 > 0 {
		maxPage++
	}
	log.Infof("Set search max page %d.", maxPage)
	if resource = strings.ToLower(resource); resource != "web" {
		resource = "host"
	}
	log.Infof("Set search type %s.", resource)
	// 多页搜索（结果合并）
	result, err := zoom.MultiToOneSearch(dork, maxPage, resource, "")
	if err != nil {
		log.Errorf("Error when search target")
		return err
	}
	// 对搜索结果进行筛选并且保存
	f, _ := os.Create(filename)
	defer f.Close()

	filt := result.Filter("ip")
	for _, v := range filt {
		ip := v["ip"]
		switch ip.(type) {
		case string:
			f.WriteString(ip.(string) + "\n")
		default:
			log.Warnf("Unknow type of result :%v", ip)
		}
	}
	log.Infof("Save result in %s", filename)
	return nil
}
