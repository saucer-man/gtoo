package utils

import (
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"runtime"
	"time"
)

func IsRoot() bool {
	if runtime.GOOS != "windows" {
		return os.Getuid() == 0
	}
	return false
}

// 产生随机字符串
func RandomStr(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz1234567890")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// 通过两重循环过滤重复元素 时间换空间
func removeRepByLoop(slc []string) []string {
	result := []string{} // 存放结果
	for i := range slc {
		flag := true
		for j := range result {
			if slc[i] == result[j] {
				flag = false // 存在重复元素，标识为false
				break
			}
		}
		if flag { // 标识为false，不添加进结果
			result = append(result, slc[i])
		}
	}
	return result
}

// 通过map主键唯一的特性过滤重复元素 空间换时间
func removeRepByMap(slc []string) []string {
	result := []string{}
	tempMap := map[string]byte{} // 存放不重复主键
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, e)
		}
	}
	return result
}

// 元素去重
func RemoveRep(slc []string) []string {
	if len(slc) < 1024 {
		// 切片长度小于1024的时候，循环来过滤
		return removeRepByLoop(slc)
	} else {
		// 大于的时候，通过map来过滤
		return removeRepByMap(slc)
	}
}

// 按照tag打印结构体
func PrintUseTag(ptr interface{}) error {

	// 获取入参的类型
	t := reflect.TypeOf(ptr)

	// 入参类型校验
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		fmt.Println("参数应该为结构体指针")
		fmt.Println(t.Kind())
		return fmt.Errorf("参数应该为结构体指针")
	}

	// 取指针指向的结构体变量
	v := reflect.ValueOf(ptr).Elem()

	// 解析字段
	for i := 0; i < v.NumField(); i++ {
		// fmt.Println(v.Field(i).Type().Name())
		if v.Field(i).Type().Kind() == reflect.Struct {
			PrintUseTag(v.Field(i).Addr().Interface())
			continue
		}
		// 取tag
		fieldInfo := v.Type().Field(i)
		tag := fieldInfo.Tag

		// 解析label tag
		label := tag.Get("label")
		if label == "" {
			label = fieldInfo.Name + ": "
		}

		// 取出value
		value := fmt.Sprintf("%v", v.Field(i))

		fmt.Println(label + ": " + value)
	}

	return nil
}
