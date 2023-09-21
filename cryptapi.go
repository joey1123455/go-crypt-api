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
 * CryptWrapper - an interface defining the crypt api library
 * @GenPaymentAdress - returns a payment wallet address
 * @CheckLogs - checks payment logs for requets
 * @GenQR - generates a qr code for payment
 * @EstTransactionFee - estimates the fee of transaction
 * @convertRates - converts the rates from fiat to crypto or crypto to fiat
 * @GetCoins - returns a list of supported coins
 */
type CryptWrapper interface {
}

/*
 * GetCoins - info on the supported coins
 * @w - ptr to wrapper instance (reciever arg)
 * returns - map of supported coins
 */
func (w *Crypt) GetCoins() (map[string]any, error) {
	// TODO: make this request async
	info, err := utils.Info("")
	if err != nil {
		return nil, nil
	}
	if info == nil {
		return nil, nil
	}
	delete(info, "fee_tiers")
	coins := make(map[string]any)
	for chain, data := range info {
		_, isBaseCoin := data.(map[string]interface{})["ticker"]
		if isBaseCoin {
			coins[chain] = data
		} else {
			baseTicker := chain + "_"
			for token, subData := range data.(map[string]interface{}) {
				coins[baseTicker+token] = subData
			}
		}
	}
	return coins, nil
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

	// TODO: make concurrent
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
