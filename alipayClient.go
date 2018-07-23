package alipay

import (
	"crypto"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/worldlove/sec"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func NewAliClient(appID, gateway, signType string, privateKey, publicKey []byte, aliPubKey string) *AliClient {
	var hash crypto.Hash
	if signType == "RSA2" {
		hash = crypto.SHA256
	} else if signType == "RSA" {
		hash = crypto.SHA1
	} else {
		log.Panic("SignTypeError")
	}
	rsaClient, err := sec.NewRSADefault(privateKey, publicKey)
	if err != nil {
		log.Panic("PrivateKeyOrPublicKeyError")
	}
	rsaClient.SetHash(hash)
	var alipub = sec.ParsePublicKey(aliPubKey)
	err = rsaClient.SetBusinessPubKey(alipub)
	if err != nil {
		log.Panic("AliPublicKeyError")
	}

	return &AliClient{
		appID:     appID,
		gateway:   gateway,
		signType:  signType,
		RSAClient: rsaClient,
	}
}

type AliClient struct {
	appID     string
	gateway   string
	signType  string
	RSAClient sec.RSACipher // rsa 加密客户端
	AESClient sec.AESCipher // aes 加密客户端
}

type Response map[string]string

func (this Response) Get(key string) string {
	return this[key]
}
func (this Response) Set(key, value string) {
	this[key] = value
}

func (this Response) Parse(raw []byte) ([]byte, []byte) {
	var reg = regexp.MustCompile(`{".*":({.*}),"sign":"(.*)"}`)
	var matchs = reg.FindAllSubmatch(raw, -1)
	if len(matchs) != 1 {
		return nil, nil
	}
	match := matchs[0]
	if len(match) != 3 {
		return nil, nil
	}
	json.Unmarshal(match[1], &this)

	sign, _ := base64.StdEncoding.DecodeString(string(match[2]))
	return match[1], sign
}

func (this *AliClient) Sign(params url.Values) {
	params.Set("app_id", this.appID)
	params.Set("sign_type", this.signType)
	params.Del("sign")
	var toSign, _ = url.QueryUnescape(params.Encode())
	var sign, _ = this.RSAClient.Sign([]byte(toSign))
	params.Add("sign", base64.StdEncoding.EncodeToString(sign))
}

func (this *AliClient) SignToUrl(params url.Values) string {
	this.Sign(params)
	return this.gateway + "?" + params.Encode()
}

func (this *AliClient) SignToJSON(params url.Values) ([]byte, error) {
	this.Sign(params)
	var res = copyValuesToMap(params)
	res["gateway"] = this.gateway
	return json.Marshal(res)
}

func (this *AliClient) SignAndRequest(params url.Values, method string, result Response) error {
	this.Sign(params)
	log.Println(params)
	log.Println(params.Encode())
	client := &http.Client{}
	req, err := http.NewRequest(method, this.gateway, strings.NewReader(params.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Func:SignAndRequest Error:%v", err)
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	log.Println("body:", string(body))
	//	ioutil.WriteFile("a.html", body, 0640)
	waitSign, sign := result.Parse(body)
	if waitSign == nil || sign == nil {
		//panic(errors.New("ResponseParseError"))
		return errors.New("ResponseParseError")
	}
	if err := this.RSAClient.VerifyBusiness(waitSign, sign); err != nil {
		log.Printf("Func:SignAndRequest Error:%v", err)
		return err
	}

	return nil
}
func (this *AliClient) SignAndRequestBody(params url.Values, method string) ([]byte, error) {
	this.Sign(params)
	log.Println(params)
	log.Println(params.Encode())
	client := &http.Client{}
	req, err := http.NewRequest(method, this.gateway, strings.NewReader(params.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Func:SignAndRequest Error:%v", err)
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (this *AliClient) VerifyResponseURL(params url.Values) error {
	if !this.CheckAppType(params) {
		return errors.New("AppTypeError")
	}
	sign := params.Get("sign")
	params.Del("sign")
	params.Del("sign_type")
	waitSign, _ := url.QueryUnescape(params.Encode())
	signBytes, _ := base64.StdEncoding.DecodeString(sign)
	return this.RSAClient.VerifyBusiness([]byte(waitSign), signBytes)
}

func (this *AliClient) CheckAppType(params url.Values) bool {
	return params.Get("sign_type") == this.signType && params.Get("app_id") == this.appID
}
