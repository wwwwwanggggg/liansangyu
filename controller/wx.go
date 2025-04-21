package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"liansangyu/config"
	"net/http"
)

var online = false

type WxJSON struct {
	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`
	Unionid    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

func code2session(code string) (string, error) {

	if !online {
		return code, nil
	}

	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", config.Config.WxAppID, config.Config.WxAppSecret, code)

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return "", errors.New("建立请求出错")
	}

	client := http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return "", errors.New("发送请求出错")
	}

	defer resp.Body.Close()

	n := make([]byte, 1024)

	_, err = resp.Body.Read(n)
	if err != nil {
		return "", errors.New("读取响应失败")
	}

	var info WxJSON
	if err := json.Unmarshal(n, &info); err != nil {
		return "", errors.New("解析json失败")
	}

	if info.ErrCode != 0 {
		return "", errors.New(info.ErrMsg)
	}

	return info.Openid, nil

}
