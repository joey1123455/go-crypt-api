package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

// the base url for crypt apis
const baseUrl = "https://api.cryptapi.io"

/*
 * Request - handles requests to url endpoints
 * @coin - coin
 * @endpoint - url to request
 * @params - request params
 * returns - unmarshalled json body and an error
 */
func Request(coin, endpoint string, params map[string]string) (map[string]any, error) {
	parsedCoin := strings.Replace(coin, "_", "/", -1)
	var url *url.URL
	var err error
	// if coin is specified make request to ticker endpoint else generall endpoint
	if coin != "" {
		toParse := baseUrl + "/" + parsedCoin + "/" + endpoint + "/"
		url, err = url.Parse(toParse)
		if err != nil {
			return nil, err
		}
	} else {
		toParse := baseUrl + "/" + endpoint + "/"
		url, err = url.Parse(toParse)
		if err != nil {
			return nil, err
		}
	}

	// extract querry parameters
	if params != nil {
		query := url.Query()
		for key, value := range params {
			query.Add(key, value)
		}
		url.RawQuery = query.Encode()
	}
	req, _ := http.NewRequest("GET", url.String(), nil)
	req.Header.Set("referer", baseUrl)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data := make(map[string]any)
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

/*
 * Info - get information related to currency
 * @coin - coin to retrieve
 * returns - unmarshaled json body or an error
 */
func Info(coin string) (map[string]any, error) {
	var err error
	var res map[string]any
	var params = make(map[string]string)
	if coin == "" {
		params["price"] = "0"
	}
	res, err = Request(coin, "info", params)
	if err != nil {
		return nil, err
	}
	if coin == "" || res["status"] == "success" {
		return res, nil
	}
	err = errors.New("no coin was provided")
	return nil, err
}

/*
 * GetCoins - info on the supported coins
 * @w - ptr to wrapper instance (reciever arg)
 * returns - map of supported coins
 */
func GetCoins() (map[string]any, error) {
	info, err := Info("")
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, errors.New("failed to get coin info")
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
