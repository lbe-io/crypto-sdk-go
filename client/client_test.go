package client

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	sdk, err := NewClient("host", "merchant", "key", "secret")
	if err != nil {
		t.Fatalf("初始化失败: %v", err)
	}
	if sdk.config.ServerHost != "host" || sdk.config.MerchantID != "merchant" || sdk.config.APIKey != "key" || sdk.config.APISecret != "secret" {
		t.Fatalf("参数未正确赋值")
	}
}

func TestNewClient_EmptyParams(t *testing.T) {
	_, err := NewClient("", "merchant", "key", "secret")
	if err == nil {
		t.Fatalf("缺少参数时应返回错误")
	}
}
