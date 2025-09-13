package xscan

import "testing"

func TestClean(t *testing.T) {
	var toolOut = "C:/Users/admin/Downloads/xscan-spider.json"
	var output = "C:/Users/admin/Downloads/xscan-spider.html"
	Clean(toolOut, output)
}
