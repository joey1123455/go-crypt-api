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
func Request(coin, endpoint string, params map[string]any) (map[string]any, error) {
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
			query.Add(key, value.(string))
		}
		url.RawQuery = query.Encode()
	}

	// make request
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
	var params = make(map[string]any)
	if coin == "" {
		params["price"] = 0
	}
	// make request : make this a goroutine
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
