package main

import (
	"encoding/json"
	"fmt"
	"github.com/corpix/uarand"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
	"time"
)

type Account struct {
	Proxy     string
	Cookie    string
	Pass      string
	UserAgent string
	Csrf      string
	Client    *fasthttp.Client
}

type BuyRequest struct {
	CategoryId   string `json:"categoryId"`
	DistributeId string `json:"distributeId"`
	Size         string `json:"size"`
}

type buyResp struct {
	Success bool `json:"success"`
}

func (a *Account) getCsrf() {
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()
	req.SetRequestURI("https://www.kucoin.com/_api/ucenter/user-info?lang=en_US")
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.Set("User-Agent", a.UserAgent)
	req.Header.Set("Cookie", a.Cookie)
	err := a.Client.Do(req, res)
	if err != nil {
		fmt.Println(err)
	}
	fasthttp.ReleaseRequest(req)
	type kucoinResp struct {
		Data struct {
			Csrf string `json:"csrf"`
		} `json:"data"`
	}
	s := &kucoinResp{}
	err = json.Unmarshal(res.Body(), s)
	if err != nil {
		fmt.Println(err)
	}
	fasthttp.ReleaseResponse(res)
	a.Csrf = s.Data.Csrf
}

// BuyNFT пытается купить нфт и возвращает true/false
func (a *Account) BuyNFT(n int) bool {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(fmt.Sprintf("https://www.kucoin.com/_api/spot-nft/buy/normal?c=%s", a.Csrf))
	req.Header.SetMethod(fasthttp.MethodPost)
	req.SetBodyRaw(buyBody)
	req.Header.Set("User-Agent", a.UserAgent)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", a.Cookie)

	res := fasthttp.AcquireResponse()
	err := a.Client.Do(req, res)
	if err != nil {
		fmt.Println(err)
		return false
	}
	fasthttp.ReleaseRequest(req)

	ok := &buyResp{}
	err = json.Unmarshal(res.Body(), ok)
	if err != nil {
		fmt.Println(err)
		return false
	}

	//fmt.Println(string(res.Body()))

	fasthttp.ReleaseResponse(res)
	if !ok.Success {
		return false
	}
	fmt.Println("Succes bought acc ", n)
	return true
}

func (a *Account) InitAcc(cook string, prox string, pass string) {
	a.Client = &fasthttp.Client{
		MaxConnsPerHost:          1000,
		NoDefaultUserAgentHeader: true,
		Dial:                     fasthttpproxy.FasthttpHTTPDialer(prox),
		MaxIdleConnDuration:      time.Second * 10,
		ReadTimeout:              time.Second * 2,
		WriteTimeout:             time.Second * 2,
	}
	a.Cookie = cook
	a.Pass = pass
	a.UserAgent = uarand.GetRandom()
	a.getCsrf()
}

func (a *Account) Validate() {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(fmt.Sprintf("https://www.kucoin.com/_api/ucenter/verify-validation-code?c=%s", a.Csrf))
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.Set("User-Agent", a.UserAgent)
	req.Header.Set("Cookie", a.Cookie)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	b := fmt.Sprintf("bizType=PURCHASE_SPOT_NFT&validations[withdraw_password]=%s", a.Pass)
	req.SetBodyString(b)

	res := fasthttp.AcquireResponse()
	err := a.Client.Do(req, res)
	if err != nil {
		fmt.Println(err)
	}
	fasthttp.ReleaseRequest(req)
	fasthttp.ReleaseResponse(res)
}

