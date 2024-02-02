package common

import (
	"fmt"
	"os"
)

const TTL = 30

var InvalidCode = DangerForeground.Render("TOTP secret is invalid")

func Check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}
