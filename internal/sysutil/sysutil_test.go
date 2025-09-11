package sysutil

import "testing"

func TestName(t *testing.T) {
	processItems, err := GetChidrenByPID(38004)
	if err != nil {
		t.Error(err)
	}
	for _, processItem := range processItems {
		t.Log(processItem)
	}
}
