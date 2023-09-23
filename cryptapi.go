/*
 * go lang wrapper over cryptapi endpoints
 * @author - Joseph Folayan
 * @email - folayanjoey@gmail.com
 * @github - github.com\joey1123455
 */
package cryptapi

import (
	"errors"
	"net/url"

	"github.com/joey1123455/go-crypt-api/utils"
)

// the base url for crypt apis
// const BaseUrl = "https://api.cryptapi.io"

/*
 * Crypt - an instance representing the connection to the crypt api
 * @Coin - the currency to transact in
 * @OwnAddress - the wallet to recieve payment
 * @CallBack - url to send the request status
 * @Params - url querry parrams
 * @CaParams - querry params
 * @PaymentAddrs - the wallet to send payment to
 */
type Crypt struct {
	Coin         string
	OwnAddress   string
	CallBack     string
	Params       map[string]any
	CaParams     map[string]any
	PaymentAddrs string
}

/*
 * InitCryptWrapper - creates the crypt request instance
 * @coin - the currency to transact in
 * @ownAddress - the wallet to recieve payment
 * @callBack - url to send the request status
 * @paymentAddrs - the wallet to send payment to
 * @params - url querry parrams
 * @caParams - querry params
 * returns - ptr to crypt instance
 */
func InitCryptWrapper(coin, ownAddress, callBack, paymentAddrs string, params, caParams map[string]any) *Crypt {
	// TODO: check if coin is valid
	return &Crypt{
		Coin:         coin,
		OwnAddress:   ownAddress,
		CallBack:     callBack,
		PaymentAddrs: paymentAddrs,
		Params:       params,
		CaParams:     caParams,
	}
}

/*
 * EstTransactionFee - returns the estimated value for a crypto transaction
 * @coin - the coin being transacted in
 * @adress - no of address being credited
 * @priority - credit priority settings
 * returns - response data or error
 */
func EstTransactionFee(coin string, address int, priority string) (map[string]any, error) {
	res, err := utils.Request(coin, "estimate", map[string]any{
		"address":  address,
		"priority": priority,
	})
	if err != nil {
		return nil, err
	}
	if res["status"] == "success" {
		return res, nil
	}
	return nil, errors.New("failed to collect estimate")
}

/*
 * Convert - This method allows you to easily convert prices from FIAT to Crypto or even between cryptocurrencies
 * @coin - currency to convert to
 * @value -  the anount to convert
 * @from - currency to convert from
 * @returns - json response and error
 */
func Convert(coin string, value int, from string) (map[string]any, error) {
	param := map[string]any{
		"value": value,
		"from":  from,
	}
	res, err := utils.Request(coin, "convert", param)
	if err != nil {
		return nil, err
	}
	if res["status"] == "success" {
		return res, nil
	}
	return nil, errors.New("filed to convert currency")
}

/*
 * CryptWrapper - an interface defining the crypt api library
 * @GenPaymentAdress - returns a payment wallet address
 * @CheckLogs - checks payment logs for requets
 * @GenQR - generates a qr code for payment
 * @EstTransactionFee - estimates the fee of transaction
 */
type CryptWrapper interface {
}

/*
 * GenPaymentAdress - creates the address for customer to pay too
 * @w - crypt instance (reciever method)
 * returns - payment wallet or error
 */
func (w *Crypt) GenPaymentAdress() (string, error) {
	if w.Coin == "" || w.CallBack == "" || w.OwnAddress == "" {
		return "", errors.New("incomplte information")
	}
	callBackUrl, err := url.Parse(w.CallBack)
	query := callBackUrl.Query()
	if err != nil {
		return "", err
	}
	for key, val := range w.Params {
		query.Add(key, val.(string))
	}
	params := make(map[string]any)
	for key, val := range w.CaParams {
		params[key] = val
	}
	params["callback"] = url.QueryEscape(callBackUrl.String())
	params["address"] = w.OwnAddress

	res, err := utils.Request(w.Coin, "create", params)
	if err != nil {
		return "", err
	}
	if res["status"] == "success" {
		add := res["address_in"].(string)
		w.PaymentAddrs = add
		return add, nil
	}
	return "", errors.New("failed to gen wallet address")
}

/*
 * CheckLogs - provides logs for transactions sent to a payment wallet
 * @w - ptr to crypt instance (reciever method)
 * returns - logs or error
 */
func (w *Crypt) CheckLogs() (map[string]any, error) {
	if w.Coin == "" || w.CallBack == "" {
		return nil, errors.New("incomplete data")
	}
	callBackUrl, err := url.Parse(w.CallBack)
	query := callBackUrl.Query()
	if err != nil {
		return nil, err
	}
	for key, val := range w.Params {
		query.Add(key, val.(string))
	}
	callBack := url.QueryEscape(callBackUrl.String())
	res, err := utils.Request(w.Coin, "logs", map[string]any{
		"callback": callBack,
	})
	if err != nil {
		return nil, err
	}
	if res["status"] == "success" {
		return res, nil
	}
	return nil, errors.New("error while checking logs")
}

/*
 * GenQR - generates the qr code for the transaction
 * value - qr value
 * size - qr code size
 * returns qr code or error
 */
func (w *Crypt) GenQR(value string, size int) (map[string]any, error) {
	if size == 0 {
		size = 512
	}
	params := map[string]any{
		"address": w.PaymentAddrs,
	}
	if value != "" {
		params["value"] = value
	}
	params["size"] = size
	res, err := utils.Request(w.Coin, "qrcode", params)
	if err != nil {
		return nil, err
	}
	if res["status"] == "success" {
		return res, nil
	}
	return nil, errors.New("failed to generate qr code")
}
