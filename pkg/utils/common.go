package utils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

// RandString 随机长度字符串
func RandString(n int) string {
	var letterRunes = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// PrettyPrint 以 JSON 格式打印对象结构
func PrettyPrint(i interface{}) {
	s, _ := json.MarshalIndent(i, "", "\t")
	fmt.Printf("%s\n", s)
}

// splitIDs 将 ID 列表字符串转换为数值列表
func SplitIDs(s string) []uint64 {
	l := strings.Split(s, ",")
	var ret []uint64
	for _, l := range l {
		i, _ := strconv.ParseUint(l, 10, 64)
		ret = append(ret, i)
	}
	return ret
}

// IsInList 判断项目是否在列表内
func IsInList(list []string, target string) bool {
	for _, item := range list {
		if item == target {
			return true
		}
	}
	return false
}
