package cryptapi

import (
	"testing"

	"github.com/joey1123455/go-crypt-api/utils"
)

const callBackUrl = "https://webhook.site/fc6e4031-66e7-45ef-9b9c-b6101958478dsa1"

// test get coins util
func TestGetCoins(t *testing.T) {
	_, err := utils.GetCoins()
	if err != nil {
		t.Log(err)
		t.Error("fail")
	}
	t.Log("success")
}

// test gen address method
func TestGenPaymentAdress(t *testing.T) {
	var ci = InitCryptWrapper("polygon_usdt", "0xA6B78B56ee062185E405a1DDDD18cE8fcBC4395d", callBackUrl, "", map[string]any{
		"order_id":    12345,
		"multi_chain": 1,
	}, map[string]any{})
	_, err := ci.GenPaymentAdress()
	if err != nil {
		t.Log(err)
		t.Error("fail")
	}
	t.Log("success")
}

// test get logs
func TestCheckLogs(t *testing.T) {
	var ci = InitCryptWrapper("polygon_matic", "0xA6B78B56ee062185E405a1DDDD18cE8fcBC4395d", callBackUrl, "", map[string]any{
		"order_id":    12345,
		"multi_chain": 1,
	}, map[string]any{})
	_, err := ci.CheckLogs()
	if err != nil {
		t.Log(err)
		t.Error("fail")
	}
	t.Log("success")
}
