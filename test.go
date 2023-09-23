package cryptapi

import (
	"testing"

	"github.com/joey1123455/go-crypt-api/utils"
)

const callBackUrl = "https://webhook.site/fc6e4031-66e7-45ef-9b9c-b6101958478dsa1"

func checkErr(t *testing.T, err error) {
	if err != nil {
		t.Log(err)
		t.Error("fail")
	}
	t.Log("success")
}

// test get coins util
func TestGetCoins(t *testing.T) {
	_, err := utils.GetCoins()
	checkErr(t, err)
}

// test gen address method
func TestGenPaymentAdress(t *testing.T) {
	var ci = InitCryptWrapper("polygon_usdt", "0xA6B78B56ee062185E405a1DDDD18cE8fcBC4395d", callBackUrl, "", map[string]any{
		"order_id":    12345,
		"multi_chain": 1,
	}, map[string]any{})
	_, err := ci.GenPaymentAdress()
	checkErr(t, err)
}

// test get logs
func TestCheckLogs(t *testing.T) {
	var ci = InitCryptWrapper("polygon_matic", "0xA6B78B56ee062185E405a1DDDD18cE8fcBC4395d", callBackUrl, "", map[string]any{
		"order_id":    12345,
		"convert":     1,
		"multi_chain": 1,
	}, map[string]any{})
	_, err := ci.CheckLogs()
	checkErr(t, err)
}

// test getting QR code
func TestGenQR(t *testing.T) {
	var ci = InitCryptWrapper("polygon_matic", "0xA6B78B56ee062185E405a1DDDD18cE8fcBC4395d", callBackUrl, "", map[string]any{
		"convert":     1,
		"multi_chain": 1,
	}, map[string]any{})
	ci.GenPaymentAdress()
	_, err := ci.GenQR("1", 300)
	checkErr(t, err)
}

// test getting estimate
func TestEstTransactionFee(t *testing.T) {
	_, err := EstTransactionFee("polygon_matic", 1, "default")
	checkErr(t, err)
}

// test convert
func TestConvert(t *testing.T) {
	_, err := Convert("polygon_matic", 1, "default")
	checkErr(t, err)
}
