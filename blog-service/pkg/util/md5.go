package util

import (
	"crypto/md5"
	"encoding/hex"
)

//对传入的文件格式化   加密写入
func EncodeMD5(vaule string) string {
	m := md5.New()
	m.Write([]byte(vaule))

	return hex.EncodeToString(m.Sum(nil))
}
