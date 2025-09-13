package netutil

import (
	"net/url"

	"golang.org/x/net/publicsuffix"
)

func GetUrlHostname(rawUrl string) string {
	urlParse, err := url.Parse(rawUrl)
	if err != nil {
		return ""
	}
	return urlParse.Hostname()
}

func GetUrlMainDomain(rawUrl string) string {
	var hostname = GetUrlHostname(rawUrl)
	domain, err := publicsuffix.EffectiveTLDPlusOne(hostname)
	if err != nil {
		return ""
	}
	return domain
}

func GetBaseUrl(rawUrl string) string {
	urlParse, err := url.Parse(rawUrl)
	if err != nil {
		return ""
	}
	urlParse.RawQuery = ""
	urlParse.Fragment = ""
	return urlParse.String()
}
