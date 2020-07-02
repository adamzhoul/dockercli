package pkg

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
)

func EncryptionArithmetic(username, token string) string {
	var signature string
	str := fmt.Sprintf("%s-%s", username, token)
	h := md5.New()
	h.Write([]byte(str))
	signature = hex.EncodeToString(h.Sum(nil))
	signature = strings.ToLower(signature)
	return signature
}
