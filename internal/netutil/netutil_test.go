package netutil

import "testing"

func TestName(t *testing.T) {
	var rawUrl = "https://edu.xazlsec.com"
	t.Log(GetUrlHostname(rawUrl))
	t.Log(GetUrlMainDomain(rawUrl))
}
