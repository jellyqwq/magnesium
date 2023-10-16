package utils

import (
	"math/rand"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

// 生成一个长度为n的随机数字
func RandiInt(n int) string {
	// 这样每次调用这个随机数都会生成新的种子
	rand.New(rand.NewSource(time.Now().UnixNano()))
	str := ""
	for i:=0; i<n; i++ {
		str += strconv.FormatInt(rand.Int63n(10), 10)
	}
	log.Debug(str)
	return str
}
