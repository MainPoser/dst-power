package cryptox

import (
	"github.com/MainPoser/dst-power/pkg/util/gmd5"
)

// EncodeMD5 加密
func EncodeMD5(str string) (string, error) {
	// 第一次MD5加密
	str, err := gmd5.Encrypt(str)
	if err != nil {
		return "", err
	}
	// 第二次MD5加密
	str, err = gmd5.Encrypt(str)
	if err != nil {
		return "", err
	}
	return str, nil
}
