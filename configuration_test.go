package main_test

import "testing"
import (
	prCheck "github.com/lumasepa/prometheus_checkman"
	"fmt"
)

func TestLoadConf(t *testing.T) {
	checks, err := prCheck.ParseConf("./checkman.yml")
	if err != nil {
		t.Fatal("error parsing conf file", err)
	}
	fmt.Print("#v", checks)
}
