# **CRYPTAPI GO LANG WRAPPER**
An idiomatic wrapper written over the core cryptapi endpoints

[cryptapi docs](https://docs.cryptapi.io/)

---

## **Installation**
To install the SDK

```console
$ go get github.com/joey1123455/go-crypt-api
```
---

## **Usage**
Initiate a crypt instance to keep track of transactions per orders, users or vendors define your preffered criteria and specify it using a base call back webhook and querry params eg: `"www.callback.com/?order_id=1&user_id=10"`

- **coin** - currency to accept payment with this instance.
- **ownAddress** - the address to recieve payments, multiple wallets can be specified using this format, `<percentage_1>@<address_1>|<percentage_2>@<address_2>` max of 20 wallets and all percentages should add up to 100%
- **callBackUrl** - the base url for the webhook.
- **callBackParams** - unique callback params.
- **urlParams** - request params.

*Returns* - a pointer to the payment instance.

[link](https://docs.cryptapi.io/) to endpoint doc

```go
package main

import (
  "fmt"

  "github.com/joey1123455/go-crypt-api"
)

func main() {
    // create an instance of the SDK wrapper
    coin := "btc" 
    ownAddress := "bc1qcpxcm8cf52uv865j6qugl82twgy0lz88sfrpk4" 
    callBackUrl := "https://www.callbackexampleurl.com" 
    // eg of call back params used to make callback url unique
    // define your own params
    callBackParams := map[string]string{
		"user_id":  "10",
		"order_id": "12345",
	} 
    // eg of request params used to send requests
    urlParams := map[string]string{
		"multi_chain": "1",
	} 

    ci := cryptapi.InitCryptWrapper(
        coin, // the coin used in the instance
        ownAddress, // address to recieve payment
        callBackUrl, // payment notification webhook
        urlParams, // url querry parameters
        callBackParams // parameters to add to 
        )
}
```
---

### **Generating a Payment Address**
Generating an address provides a way to recieve and track payment related to a unique callback url

*Returns* - an address and an error

[link](https://docs.cryptapi.io/#operation/create) to endpoint doc

```go
    // generate the payment addresss
    address, err := ci.GenPaymentAdress()
    if err != nil {
      fmt.Println("while generating payment wallet:", err)
    } else {
      fmt.Println("payment wallet:", address)
    }
```

---
### **Checking Callback Logs**

Logs returns the history of transaction related to a unique call back url. 

*Returns* - the logs and an error

Before checking logs ensure you have generated a payment address using `ci.GenPaymentAdress()`

[link](https://docs.cryptapi.io/#operation/logs) to endpoint doc

```go
    // check payment history
    history, err := ci.CheckLogs()
    if err != nil {
      fmt.Println("while fetching callback history:", err)
    } else {
      fmt.Println("callback history:", history)
    }
```
---

### **Generating a QR Code**

Generates a qr code for the specified transaction.

- **size** - The size of the qr code in pixels.
- **value** - The transaction price.

*Returns* the qr code and a error.

[link](https://docs.cryptapi.io/#operation/qrcode) to endpoint doc

Before generating a qr code ensure you have generated a payment address using `ci.GenPaymentAdress()`

```go
    // generate payment qr code
    size := "300"
    value := "100"
    qr, err := ci.GenQR(
      value, // the quantity being charged may be left empty 
      size // the pixel size of the qr code
      )
    if err != nil {
      fmt.Println("while generating qr code:", err)
    } else {
      fmt.Println("qr code:", history)
    }
```

---
### **Estimating Transaction Fees**
Retrieves an estimated price forthe crypto transaction.

- **coin** - currency being traded.
- **address count** - no of wallets recieving payment.
- **priority** - the priority to use when making payments

*Returns* - estimate and an error

[link](https://docs.cryptapi.io/#operation/estimate) to endpoint doc

```go
package main

import (
  "fmt"

  "github.com/joey1123455/go-crypt-api"
)

func main() {
  coin := "polygon_matic"
  priority := "default"
  add := "1"
  est, err := cryptapi.EstTransactionFee(
    coin, // currency being traded 
    add, // no of addresses to send to
    priority // the priority taken when disbursing payment
    )
  if err != nil {
      fmt.Println("while fetching callback history:", err)
    } else {
      fmt.Println("estimated transaction fee:", est)
    }
}
```

---
### **Converting Between Currencies**
Converts between fiat or crypto currencies.

- **coin** - the currency being converted.
- **value** - quantity to compare.
- **to** - the currency to convert to

*Returns* - the conversion rate and an error.

[link](https://docs.cryptapi.io/#operation/convert) to endpoint doc
```go
package main

import (
  "fmt"

  "github.com/joey1123455/go-crypt-api"
)

func main() {
  coin := "btc"
  value := "1"
  from := "USD"
  conv, err := cryptapi.Convert(
    coin, // default currency
    value, // amount
    from // currency to convert to
    )
  if err != nil {
      fmt.Println("while converting:", err)
    } else {
      fmt.Println("converted:", conv)
    }
}
```
---