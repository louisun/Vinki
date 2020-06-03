package utils

import (
	"errors"

	"github.com/speps/go-hashids"
)

// ID类型
const (
	UserID = iota // 用户
)

var (
	// ErrTypeNotMatch ID类型不匹配
	ErrTypeNotMatch = errors.New("ID类型不匹配")
)

// HashEncode 对给定数据计算HashID
func HashEncode(values []int, salt string) (string, error) {
	hd := hashids.NewData()
	hd.Salt = salt

	h, err := hashids.NewWithData(hd)
	if err != nil {
		return "", err
	}

	id, err := h.Encode(values)
	if err != nil {
		return "", err
	}
	return id, nil
}

// HashDecode 对给定数据计算原始数据
func HashDecode(hashString string, salt string) ([]int, error) {
	hd := hashids.NewData()
	hd.Salt = salt

	h, err := hashids.NewWithData(hd)
	if err != nil {
		return []int{}, err
	}

	return h.DecodeWithError(hashString)

}

// GenerateHash 计算数据库内主键对应的HashID
func GenerateHash(id uint64, pkType int, salt string) string {
	v, _ := HashEncode([]int{int(id), pkType}, salt)
	return v
}

// GetOriginID 计算HashID对应的数据库ID
func GetOriginID(hashString string, pkType int, salt string) (uint, error) {
	v, _ := HashDecode(hashString, salt)
	if len(v) != 2 || v[1] != pkType {
		return 0, ErrTypeNotMatch
	}
	return uint(v[0]), nil
}
