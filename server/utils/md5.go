package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

func GenerateSaltedMD5(input string, salt string) string {
	saltedInput := input + salt
	hasher := md5.New()
	io.WriteString(hasher, saltedInput)
	return hex.EncodeToString(hasher.Sum(nil))
}
