package utils

import (
	"errors"
	"net/http"
	"reflect"
	"strings"

	"io/ioutil"
	// "fmt"
	"time"

	"github.com/ajg/form"

	"fmt"
	"math/rand"
)

func RequstCleint(method, url string, data *map[string]interface{}) (body []byte, err error) {

	httpClient := &http.Client{
		Timeout: time.Second,
	}
	var formData string
	formData, err = form.EncodeToString(data)
	if err != nil {
		return body, err
	}
	fmt.Println(formData)
	var req *http.Request
	req, err = http.NewRequest(method, url, strings.NewReader(formData))
	if err != nil {
		return body, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8;")
	res, err := httpClient.Do(req)
	if err != nil {
		return body, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return body, errors.New("请求失败!\n")

	}
	body, _ = ioutil.ReadAll(res.Body)

	return body, err

}

// 生成随机数验证码
func GenValidateCode(len int) string {
	numbers := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < len; i++ {
		fmt.Fprintf(&sb, "%d", numbers[rand.Intn(10)])
	}
	return sb.String()
}

//解析结构体
func ParseReflect(v reflect.Type) reflect.Type {
	switch v.Kind() {
	case reflect.Ptr:
		return ParseReflect(v.Elem())
	default:
		return v
	}
}
