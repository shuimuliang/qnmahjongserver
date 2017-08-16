package util

import (
	"crypto/sha1"
	"fmt"
	"strings"
)

func GetSqlString(sql string, params ...interface{}) string {
	for _, param := range params {
		sql = strings.Replace(sql, "?", "'"+fmt.Sprint(param)+"'", 1)
	}
	return sql
}

func Sha1Password(password string) string {
	saltedBytes := []byte("salted" + password)
	sha1Bytes := sha1.Sum(saltedBytes)
	return fmt.Sprintf("%x", sha1Bytes)
}
