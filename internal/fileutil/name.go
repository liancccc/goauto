package fileutil

import (
	"fmt"
	"net/url"
	"time"
)

func GetUrlFileName(urlRaw string) string {
	urlParse, err := url.Parse(urlRaw)
	if err != nil {
		return GetUnixNmae()
	}
	return urlParse.Hostname()
}

func GetUnixNmae() string {
	return fmt.Sprintf("%v", time.Now().Unix())
}
