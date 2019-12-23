package password

import "github.com/alexandrevicenzi/unchained"

// CryTo 哈希密码
func CryTo(password string, saltSize int, hasher string) (string, error) {
	return unchained.MakePassword(password, unchained.GetRandomString(saltSize), hasher)
}

// Validate 有效
func Validate(password string, cryto string) (bool, error) {
	return unchained.CheckPassword(password, cryto)
}
