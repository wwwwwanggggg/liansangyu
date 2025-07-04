package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const configFilePath = "./testconfig.json"
const BASEURL string = "http://0.0.0.0:8088"

const (
	COLOR_RED    = "\033[31m"
	COLOR_RESET  = "\033[0m"
	COLOR_GREEN  = "\033[32m"
	COLOR_YELLOW = "\033[33m"
)

type reqInput struct {
	Headers map[string]string
	Body    interface{}
}

func LoadConfigJSON() ConfigJSON {
	if strings.HasSuffix(configFilePath, ".apifox.json") {
		return LoadApifoxJSON()
	}
	// 从配置文件加载JSON数据
	file, err := os.OpenFile(configFilePath, os.O_RDONLY, 0644)
	if err != nil {
		panic("配置文件加载失败: " + err.Error())
	}
	defer file.Close()

	var config ConfigJSON
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		panic("配置文件解析失败: " + err.Error())
	}

	return config
}

func LoadApifoxJSON() ConfigJSON {
	file, err := os.OpenFile(configFilePath, os.O_RDONLY, 0644)
	if err != nil {
		panic("配置文件加载失败: " + err.Error())
	}
	defer file.Close()

	var config ConfigJSON
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		panic("配置文件解析失败: " + err.Error())
	}

	return config
}

// 定义测试上下文结构
type TestContext struct {
	Request *http.Request
	Cookies []*http.Cookie
}

func testFunction(tc *TestContext, method string, path string, input reqInput, wanted interface{}, wantedCode int, t *testing.T) *TestContext {
	// 确保路径格式正确
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	URL := BASEURL + path

	// 准备请求体
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(input.Body); err != nil {
		t.Fatalf("请求体编码失败: %v", err)
	}

	// 创建或复用请求
	var req *http.Request
	var err error

	if tc != nil && tc.Request != nil {
		// 复用现有请求
		req = tc.Request.Clone(context.Background())
		req.Method = method
		req.URL, err = url.Parse(URL)
		if err != nil {
			t.Fatalf("URL解析失败: %v", err)
		}
		req.Body = io.NopCloser(&buf)
		req.ContentLength = int64(buf.Len())
	} else {
		// 创建新请求
		req, err = http.NewRequest(method, URL, &buf)
		if err != nil {
			t.Fatalf("请求创建失败: %v", err)
		}
	}

	// 设置请求头
	req.Header.Set("Allow-Origin", "*")
	for k, v := range input.Headers {
		req.Header.Set(k, v)
	}

	// 添加之前保存的Cookies
	if tc != nil && len(tc.Cookies) > 0 {
		for _, cookie := range tc.Cookies {
			req.AddCookie(cookie)
		}
	}

	// 执行请求
	w := httptest.NewRecorder()
	Router.ServeHTTP(w, req)

	// 验证响应
	wantedBytes, err := json.Marshal(wanted)
	if err != nil {
		t.Fatalf("期望结果编码失败: %v", err)
	}

	assert.Equal(t, wantedCode, w.Code, "响应状态码不匹配")
	assert.JSONEq(t, string(wantedBytes), w.Body.String(), "响应体不匹配")

	// 返回更新后的测试上下文
	return &TestContext{
		Request: req,
		Cookies: w.Result().Cookies(), // 保存新的Cookies
	}
}

func TestMain(t *testing.T) {
	// 拿到测试的config
	config := LoadConfigJSON()
	gin.SetMode(gin.TestMode)
	tc := &TestContext{}

	// 开测
	for _, handler := range config.Handlers {
		for index, exp := range handler.Examples {
			in := strconv.Itoa(index + 1)
			t.Run(handler.Title+"_"+handler.Method+" "+handler.Path+"     "+in, func(t *testing.T) {
				input := reqInput{
					Headers: make(map[string]string),
					Body:    exp.Body,
				}

				if header, ok := exp.Header.(map[string]interface{}); ok {
					for k, v := range header {
						input.Headers[k] = v.(string)
					}
				}
				tc = testFunction(tc, handler.Method, handler.Path, input, exp.Wanted, exp.WantedCode, t)
			})
		}
	}

}
