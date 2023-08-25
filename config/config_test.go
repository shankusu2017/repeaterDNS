package config

import (
	"testing"
)

func TestGetDnsServerIP(t *testing.T) {
	InitConfig("./localDomain.conf")
	if GetDNSServerIP("www.shanks.link") == "223.5.5.5" {
		t.Fatalf("FATAL d13c01b3")
	}

	if GetDNSServerIP("www.huawei.com") != "223.5.5.5" {
		t.Fatalf("FATAL 6e33604c")
	}
	if GetDNSServerIP("14.cc") != "223.5.5.5" {
		t.Fatalf("FATAL b6bf3670")
	}
}
