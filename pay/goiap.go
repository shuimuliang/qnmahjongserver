package pay

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type Receipt struct {
	OriginalPurchaseDatePst string `json:"original_purchase_date_pst"`
	PurchaseDateMs          string `json:"purchase_date_ms"`
	UniqueIdentifier        string `json:"unique_identifier"`
	OriginalTransactionId   string `json:"original_transaction_id"`
	Bvrs                    string `json:"bvrs"`
	TransactionId           string `json:"transaction_id"`
	Quantity                string `json:"quantity"`
	UniqueVendorIdentifier  string `json:"unique_vendor_identifier"`
	ItemId                  string `json:"item_id"`
	ProductId               string `json:"product_id"`
	PurchaseDate            string `json:"purchase_date"`
	OriginalPurchaseDate    string `json:"original_purchase_date"`
	PurchaseDatePst         string `json:"purchase_date_pst"`
	Bid                     string `json:"bid"`
	OriginalPurchaseDateMs  string `json:"original_purchase_date_ms"`
}

type receiptRequestData struct {
	Receiptdata string `json:"receipt-data"`
}

type receiptresponseData struct {
	ReceiptContent *Receipt `json:"receipt"`
	Status         int32    `json:"status"`
}

const (
	appleSandboxURL    = "https://sandbox.itunes.apple.com/verifyReceipt"
	appleProductionURL = "https://buy.itunes.apple.com/verifyReceipt"
)

// Simple interface to get the original error code from the error object
type ErrorWithCode interface {
	Code() int32
}

type Error struct {
	error
	errCode int32
}

// Simple method to get the original error code from the error object
func (e *Error) Code() int32 {
	return e.errCode
}

// Given receiptData (base64 encoded) it tries to connect to either the sandbox (useSandbox true) or
// apples ordinary service (useSandbox false) to validate the receipt. Returns either a receipt struct or an error.
func VerifyReceipt(receiptData string, useSandbox bool) (*Receipt, error) {
	receipt, err := sendReceiptToApple(receiptData, verificationURL(useSandbox))
	return receipt, err
}

// Selects the proper url to use when talking to apple based on if we should use the sandbox environment or not
func verificationURL(useSandbox bool) string {

	if useSandbox {
		return appleSandboxURL
	}
	return appleProductionURL
}

// Sends the receipt to apple, returns the receipt or an error upon completion
func sendReceiptToApple(receiptData, url string) (*Receipt, error) {
	requestData, err := json.Marshal(receiptRequestData{receiptData})

	if err != nil {
		return nil, err
	}

	toSend := bytes.NewBuffer(requestData)

	resp, err := http.Post(url, "application/json", toSend)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	var responseData receiptresponseData
	err = json.Unmarshal(body, &responseData)

	if err != nil {
		return nil, err
	}

	if responseData.Status != 0 {
		return nil, verificationError(responseData.Status)
	}

	return responseData.ReceiptContent, nil
}

// Error codes as they returned by the App Store
const (
	UnreadableJSON       = 21000
	MalformedData        = 21002
	AuthenticationError  = 21003
	UnmatchedSecret      = 21004
	ServerUnavailable    = 21005
	SubscriptionExpired  = 21006
	SandboxReceiptOnProd = 21007
	ProdReceiptOnSandbox = 21008
)

// Generates the correct error based on a status error code
func verificationError(errCode int32) error {
	var errorMessage string

	switch errCode {
	case UnreadableJSON:
		errorMessage = "The App Store could not read the JSON object you provided."
		break
	case MalformedData:
		errorMessage = "The data in the receipt-data property was malformed."
		break

	case AuthenticationError:
		errorMessage = "The receipt could not be authenticated."
		break

	case UnmatchedSecret:
		errorMessage = "The shared secret you provided does not match the shared secret on file for your account."
		break

	case ServerUnavailable:
		errorMessage = "The receipt server is not currently available."
		break
	case SubscriptionExpired:
		errorMessage = "This receipt is valid but the subscription has expired. When this status code is returned to your server, " +
			"the receipt data is also decoded and returned as part of the response."
		break
	case SandboxReceiptOnProd:
		errorMessage = "This receipt is a sandbox receipt, but it was sent to the production service for verification."
		break
	case ProdReceiptOnSandbox:
		errorMessage = "This receipt is a production receipt, but it was sent to the sandbox service for verification."
		break
	default:
		errorMessage = "An unknown error ocurred"
		break
	}

	return &Error{errors.New(errorMessage), errCode}
}
