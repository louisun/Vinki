package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

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
