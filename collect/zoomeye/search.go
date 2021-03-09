package zoomeye

// from https://github.com/gyyyy/ZoomEye-go

import (
	"log"
)

func Search(apikey, search string) error {
	zoom := NewWithKey(apikey)

	// 查询用户资源信息
	info, err := zoom.ResourcesInfo()
	if err != nil {
		log.Printf("Error when get id info, maybe because of incorrect api key: %s\n", apikey)
		return err
	}
	log.Printf("Successfully login zoomeye. role: %s, quota: %d.\n", info.Plan, info.Resources.Search)
	// // 搜索
	// result, _ := zoom.DorkSearch("port:80 nginx", 1, "host", "app,service,os")
	// // 多页搜索，5页（100条）以上会进行并发搜索，减少搜索耗时
	// // results, _ := zoom.MultiPageSearch("wordpress country:cn", 5, "web", "webapp,server,os")
	// // 多页搜索（结果合并）
	// // result, _ := zoom.MultiToOneSearch("wordpress country:cn", 5, "web", "webapp,server,os")

	// // 对搜索结果进行统计
	// stat := result.Statistics("app,service,os")

	// // 对搜索结果进行筛选
	// filt := result.Filter("app,ip,title")

	// // 设备历史搜索（需要高级用户或VIP用户权限，结果包含多少条记录就会扣多少额度，非土豪慎用）
	// history, _ := zoom.HistoryIP("1.2.3.4")
	return nil
}
