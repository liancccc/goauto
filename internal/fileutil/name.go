package fileutil

import (
	"fmt"
	"net/url"
	"time"
)

func GetUrlFileName(urlRaw string) string {
	urlParse, err := url.Parse(urlRaw)
	if err != nil {
		return fmt.Sprintf("%s", time.Now().Unix())
	}
	return urlParse.Hostname()
}
