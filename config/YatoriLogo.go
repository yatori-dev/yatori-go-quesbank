package config

import _ "embed"

// 读取logo
//
//go:embed logo.txt
var logoStr []byte

func YatoriLogo() string {
	return string(logoStr)
}
