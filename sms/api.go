package sms

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const (
	gateway   = "http://121.204.249.117:8000/api/"
	developer = ""
)

var (
	token         string
	loginError    = errors.New("未知错误")
	responseError = errors.New("返回结果错误")
)

func Login(username, password string) error {
	fileds, err := req(fmt.Sprintf("%ssign/username=%s&password=%s", gateway, username, password))
	if err != nil {
		return err
	}
	if len(fileds) != 2 {
		return loginError
	}
	token = fileds[1]
	return nil
}

//查询余额
func QueryMoney() (float32, error) {
	fileds, err := req(fmt.Sprintf("%syh_gx/token=%s", gateway, token))
	if err != nil {
		return 0, err
	}
	f, err := strconv.ParseFloat(fileds[1], 32)
	return float32(f), err
}

///获取手机号码
func GetPhone(projectId string) (string, error) {
	fileds, err := req(fmt.Sprintf("%syh_qh/id=%s&operator=0&Region=0&card=2&phone=&loop=1&filer=&token=%s", gateway, projectId, token))
	if err != nil {
		return "", err
	}
	return fileds[1], nil
}

//获取短信内容
//1.因部分短信可能延迟，所以建议该方法每5秒调用一次以免触发系统机制被做访问限制，调用300秒
func GetMessage(projectId, phone string) (string, error) {
	fileds, err := req(fmt.Sprintf("%syh_qm/id=%s&phone=%s&t=%s&token=%s", gateway, projectId, phone, developer, token))
	if err != nil {
		return "", err
	}
	return fileds[1], nil
}

//释放手机号码
func FreePhone(projectId, phone string) error {
	_, err := req(fmt.Sprintf("%syh_sf/id=%s&phone=%s&token=%s", gateway, projectId, phone, token))
	return err
}

func req(url string) ([]string, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	} else if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get statuscode:%d", response.StatusCode)
	}
	defer response.Body.Close()

	bs, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	fileds := strings.Split(string(bs), "|")
	if len(fileds) == 0 {
		return nil, responseError
	} else if fileds[0] == "0" {
		return nil, errors.New(fileds[1])
	}
	return fileds, nil
}
