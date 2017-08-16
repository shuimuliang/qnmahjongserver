package pay

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	APP_ID_IOS           = "3001852920"
	APP_ID_ANDROID       = "3001853586"
	IPAY_ORDER_URL       = "http://ipay.iapppay.com:9999/payapi/order"
	IPAY_QUERY_URL       = "http://ipay.iapppay.com:9999/payapi/queryresult"
	IPAY_H5_REDIRECT_URL = "https://web.iapppay.com/h5/gateway"
)

type IpayHelper struct {
	appId     string
	signKey   *rsa.PrivateKey
	verifyKey *rsa.PublicKey
}

func NewIpayHelper(appId string, privKeyBytes, pubKeyBytes []byte) (helper *IpayHelper, err error) {
	block, _ := pem.Decode(privKeyBytes)
	if block == nil {
		err = errors.New("no priv key block(s)")
		return
	}
	if block.Type != "RSA PRIVATE KEY" {
		err = errors.New("block type not RSA PRIVATE KEY")
		return
	}
	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return
	}
	block, _ = pem.Decode(pubKeyBytes)
	if block == nil {
		err = errors.New("no pub key block(s)")
		return
	}
	if block.Type != "PUBLIC KEY" {
		err = errors.New("block type not PUBLIC KEY")
		return
	}
	pubKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return
	}
	pubKey, ok := pubKeyInterface.(*rsa.PublicKey)
	if !ok {
		err = errors.New("pubkey is not rsa public key")
		return
	}
	helper = &IpayHelper{
		appId:     appId,
		signKey:   privKey,
		verifyKey: pubKey,
	}
	return
}

func NewIpayHelperWithPem(appId, privPemPath, pubPemPath string) (*IpayHelper, error) {
	privBytes, err := ioutil.ReadFile(privPemPath)
	if err != nil {
		return nil, err
	}
	pubBytes, err := ioutil.ReadFile(pubPemPath)
	if err != nil {
		return nil, err
	}
	return NewIpayHelper(appId, privBytes, pubBytes)
}

func (h *IpayHelper) request(ipayUrl string, datas interface{}) (respBody []byte, err error) {
	transdata, sign, err := h.ipaySign(datas)
	if err != nil {
		return
	}
	reqForm := url.Values{
		"transdata": {transdata},
		"sign":      {sign},
		"signtype":  {"RSA"},
	}
	resp, err := http.PostForm(ipayUrl, reqForm)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = errors.New("request ipay fail: " + resp.Status)
		return
	}
	return ioutil.ReadAll(resp.Body)
}

func (h *IpayHelper) parseResponse(respBody []byte) (rspTransInfo map[string]interface{}, err error) {
	rspForm, err := url.ParseQuery(string(respBody))
	if err != nil {
		return
	}
	rspTransdata := rspForm.Get("transdata")
	if rspTransdata == "" {
		err = errors.New("response no transdata:" + string(respBody))
		return
	}
	rspTransInfo = make(map[string]interface{})
	if err = json.Unmarshal([]byte(rspTransdata), &rspTransInfo); err != nil {
		return
	}
	if errMsg, exists := rspTransInfo["errmsg"]; exists {
		code := rspTransInfo["code"]
		err = errors.New(fmt.Sprintf("ipay response error: %v|%v", code, errMsg))
		return
	}
	rspSign := rspForm.Get("sign")
	if err = h.ipayVerify(rspTransdata, rspSign); err != nil {
		return
	}
	return
}

func (h *IpayHelper) ipaySign(datas interface{}) (transdata string, sign string, err error) {
	reqJson, err := json.Marshal(datas)
	if err != nil {
		return
	}
	transdata = string(reqJson)
	d := md5.Sum(reqJson)
	signBytes, err := rsa.SignPKCS1v15(rand.Reader, h.signKey, crypto.MD5, d[:])
	if err != nil {
		return
	}
	sign = base64.StdEncoding.EncodeToString(signBytes)
	return
}

func (h *IpayHelper) ipayVerify(transdata, sign string) error {
	tBytes := []byte(transdata)
	sBytes, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return err
	}
	d := md5.Sum(tBytes)
	return rsa.VerifyPKCS1v15(h.verifyKey, crypto.MD5, d[:], sBytes)
}

func (h *IpayHelper) CreateIpayOrder(
	waresId int,
	waresName string,
	cpOrderId string,
	price float32,
	appUserId string,
	cpPrivateInfo string,
	notifyUrl string,
) (transId string, err error) {
	datas := map[string]interface{}{
		"appid":     h.appId,
		"waresid":   waresId,
		"cporderid": cpOrderId,
		"currency":  "RMB",
		"appuserid": appUserId,
		"price":     price,
	}
	if waresName != "" {
		datas["waresname"] = waresName
	}
	if cpPrivateInfo != "" {
		datas["cpprivateinfo"] = cpPrivateInfo
	}
	if notifyUrl != "" {
		datas["notifyurl"] = notifyUrl
	}
	respBody, err := h.request(IPAY_ORDER_URL, datas)
	if err != nil {
		return
	}
	rspTransInfo, err := h.parseResponse(respBody)
	if err != nil {
		return
	}
	transIdInterface, exists := rspTransInfo["transid"]
	if !exists {
		err = errors.New("response no transid")
		return
	}
	transId = transIdInterface.(string)
	return
}

type TransResult struct {
	TransType int
	CpOrderId string
	TransId   string
	AppUserId string
	AppId     string
	WaresId   int
	FeeType   int
	Money     float64
	Currency  string
	Result    int
	TransTime string
	CpPrivate string
	PayType   int
}

func readStringField(transInfo map[string]interface{}, key string, defaultValue string) string {
	if value, exists := transInfo[key]; exists {
		return value.(string)
	}
	return defaultValue
}

func readIntField(transInfo map[string]interface{}, key string, defaultValue int) int {
	if value, exists := transInfo[key]; exists {
		if iValue, ok := value.(int); ok {
			return iValue
		} else if fValue, ok := value.(float64); ok {
			return int(fValue)
		}
	}
	return defaultValue
}

func readFloatField(transInfo map[string]interface{}, key string, defaultValue float64) float64 {
	if value, exists := transInfo[key]; exists {
		return value.(float64)
	}
	return defaultValue
}

func (h *IpayHelper) QueryResult(cpOrderId string) (result *TransResult, err error) {
	datas := map[string]interface{}{
		"appid":     h.appId,
		"cporderid": cpOrderId,
	}
	respBody, err := h.request(IPAY_QUERY_URL, datas)
	if err != nil {
		return
	}
	rspTransInfo, err := h.parseResponse(respBody)
	if err != nil {
		return
	}
	result = &TransResult{
		TransType: -1,
		CpOrderId: readStringField(rspTransInfo, "cporderid", ""),
		TransId:   readStringField(rspTransInfo, "transid", ""),
		AppUserId: readStringField(rspTransInfo, "appuserid", ""),
		AppId:     readStringField(rspTransInfo, "appid", ""),
		WaresId:   readIntField(rspTransInfo, "waresid", 0),
		FeeType:   readIntField(rspTransInfo, "feetype", -1),
		Money:     readFloatField(rspTransInfo, "money", float64(0)),
		Currency:  readStringField(rspTransInfo, "currency", ""),
		Result:    readIntField(rspTransInfo, "result", -1),
		TransTime: readStringField(rspTransInfo, "transtime", ""),
		CpPrivate: readStringField(rspTransInfo, "cpprivate", ""),
		PayType:   readIntField(rspTransInfo, "paytype", -1),
	}
	return
}

func (h *IpayHelper) ParseNotifyInfo(postBytes []byte) (result *TransResult, err error) {
	rspTransInfo, err := h.parseResponse(postBytes)
	if err != nil {
		return
	}
	result = &TransResult{
		TransType: readIntField(rspTransInfo, "transtype", -1),
		CpOrderId: readStringField(rspTransInfo, "cporderid", ""),
		TransId:   readStringField(rspTransInfo, "transid", ""),
		AppUserId: readStringField(rspTransInfo, "appuserid", ""),
		AppId:     readStringField(rspTransInfo, "appid", ""),
		WaresId:   readIntField(rspTransInfo, "waresid", 0),
		FeeType:   readIntField(rspTransInfo, "feetype", -1),
		Money:     readFloatField(rspTransInfo, "money", float64(0)),
		Currency:  readStringField(rspTransInfo, "currency", ""),
		Result:    readIntField(rspTransInfo, "result", -1),
		TransTime: readStringField(rspTransInfo, "transtime", ""),
		CpPrivate: readStringField(rspTransInfo, "cpprivate", ""),
		PayType:   readIntField(rspTransInfo, "paytype", -1),
	}
	return
}

func (h *IpayHelper) GetHtml5RedirectUrl(transId, redirectUrl, cpUrl string) (targetUrl string, err error) {
	trans := map[string]string{
		"transid":     transId,
		"redirecturl": redirectUrl,
	}
	if cpUrl != "" {
		trans["cpurl"] = cpUrl
	}
	transdata, sign, err := h.ipaySign(trans)
	urlQuery := url.Values{
		"transdata": {transdata},
		"sign":      {sign},
		"signtype":  {"RSA"},
	}
	targetUrl = IPAY_H5_REDIRECT_URL + "?" + urlQuery.Encode()
	return
}

func (h *IpayHelper) GetNewHtml5RedirectUrl(tid, app, url_r, url_h string) (targetUrl string, err error) {
	trans := map[string]string{
		"tid":   tid,
		"app":   app,
		"url_r": url_r,
		"url_h": url_h,
	}
	transdata, sign, err := h.ipaySign(trans)
	urlQuery := url.Values{
		"data":      {transdata},
		"sign":      {sign},
		"sign_type": {"RSA"},
	}
	targetUrl = IPAY_H5_REDIRECT_URL + "?" + urlQuery.Encode()
	return
}
