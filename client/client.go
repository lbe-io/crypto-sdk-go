package client

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// SDKConfig 用于存储SDK初始化参数
type SDKConfig struct {
	ServerHost string
	MerchantID string
	APIKey     string
	APISecret  string
	Timeout    time.Duration // 可选超时时间
}

// CryptoSDK 是SDK的主结构体
// 封装了所有对外暴露的方法
type CryptoClient struct {
	config     SDKConfig
	options    *clientOptions
	httpClient *http.Client
}

// NewCryptoSDK 初始化SDK
// 传入 server_host, merchant_id, apikey, apisecret
func NewClient(config SDKConfig, option ...time.Duration) (*CryptoClient, error) {
	if config.ServerHost == "" || config.MerchantID == "" || config.APIKey == "" || config.APISecret == "" {
		return nil, fmt.Errorf("所有参数均不能为空")
	}
	timeout := 10 * time.Second
	if len(option) > 0 {
		timeout = option[0]
	}
	config.Timeout = timeout

	client := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
		},
	}
	return &CryptoClient{config: config, httpClient: client}, nil
}

// 签名生成
func (sdk *CryptoClient) sign(data string) string {
	mac := hmac.New(sha256.New, []byte(sdk.config.APISecret))
	mac.Write([]byte(data))
	return hex.EncodeToString(mac.Sum(nil))
}

// 生成Authorization header
func (sdk *CryptoClient) buildAuthHeader(params map[string]interface{}) string {
	b, _ := json.Marshal(params)
	sign := sdk.sign(string(b))
	return fmt.Sprintf("%s:%s", sdk.config.APIKey, sign)
}

// 验签
func (sdk *CryptoClient) verifySign(data, signature string) bool {
	return sdk.sign(data) == signature
}

// 通用请求方法
func (sdk *CryptoClient) doRequest(method, path string, reqBody interface{}, respBody interface{}) error {
	url := fmt.Sprintf("%s%s", sdk.config.ServerHost, path)
	var body io.Reader
	var params map[string]interface{}
	if reqBody != nil {
		b, _ := json.Marshal(reqBody)
		body = bytes.NewReader(b)
		json.Unmarshal(b, &params)
	} else {
		params = map[string]interface{}{}
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", sdk.buildAuthHeader(params))
	resp, err := sdk.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	// 响应验签
	respSign := resp.Header.Get("Authorization")
	if !sdk.verifySign(string(respBytes), respSign) {
		return fmt.Errorf("响应验签失败")
	}
	if respBody != nil {
		return json.Unmarshal(respBytes, respBody)
	}
	return nil
}

// GET请求
func (sdk *CryptoClient) Get(path string, req, resp interface{}) error {
	return sdk.doRequest(http.MethodGet, path, req, resp)
}

func (sdk *CryptoClient) GetPage(path string, req, resp interface{}) error {
	return sdk.doRequest(http.MethodGet, path, req, resp)
}

// POST请求
func (sdk *CryptoClient) Post(path string, req, resp interface{}) error {
	return sdk.doRequest(http.MethodPost, path, req, resp)
}

func (sdk *CryptoClient) PostPage(path string, req, resp interface{}) error {
	return sdk.doRequest(http.MethodPost, path, req, resp)
}

// 实体定义示例
type ExampleRequest struct {
	Param1 string `json:"param1"`
	Param2 int    `json:"param2"`
}

type ExampleResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

// 回调处理
func (sdk *CryptoClient) HandleCallback(r *http.Request) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	sign := r.Header.Get("Authorization")
	if !sdk.verifySign(string(body), sign) {
		return nil, fmt.Errorf("回调验签失败")
	}
	return body, nil
}
