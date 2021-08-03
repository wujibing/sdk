package kuaishibie

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

//快识别api
const (
	gateway = "http://api.kuaishibie.cn/"
)

var (
	_username, _password string
)

type (
	Request interface {
		ReqPre()
		Method() string
		GetResponse() interface{}
	}
	//图片识别
	Predict struct {
		Username  string `json:"username"`
		Password  string `json:"password"`
		TypeId    string `json:"typeid"`
		Angle     string `json:"angle"` //旋转角度：当typeid为14时旋转角度 默认90
		Image     string `json:"image"`
		ImageBack string `json:"imageback"` //(缺口识别2张图传背景图需要)
		Content   string `json:"content"`   //(快速点选需要)	标题内容 如：填写 "你好"中文请unicode编码以免造成错误。
		response  *PredictResponse
	}

	ReportError struct {
		Id string `json:"id"`
	}

	PredictResponse struct {
		Success bool        `json:"success"`
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    PredictData `json:"data"`
	}

	PredictData struct {
		Result string `json:"result"`
		Id     string `json:"id"`
	}

	Point struct {
		X int
		Y int
	}
)

func (p *PredictResponse) GetPoint() []Point {
	temps := strings.Split(p.Data.Result, "|")
	points := make([]Point, len(temps))
	for key, value := range temps {
		t := strings.Split(value, ",")
		if len(t) != 2 {
			continue
		}
		x, err := strconv.Atoi(t[0])
		y, err1 := strconv.Atoi(t[1])
		if err != nil || err1 != nil {
			continue
		}
		points[key] = Point{x, y}
	}
	return points
}

func (p *Predict) OpenImage(reader io.Reader) error {
	return p.readImage(reader, &p.Image)
}

func (p *Predict) OpenImageByUrl(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return p.readImage(resp.Body, &p.Image)
}

func (p *Predict) GetResponse() interface{} {
	if p.response == nil {
		p.response = &PredictResponse{}
	}
	return p.response
}

func (p *Predict) GetPredictResponse() *PredictResponse {
	return p.response
}

func (p *Predict) OpenImageBack(reader io.Reader) error {
	return p.readImage(reader, &p.ImageBack)
}

func (p *Predict) readImage(reader io.Reader, data *string) error {
	src, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	*data = base64.StdEncoding.EncodeToString(src)
	return nil
}

func (p *Predict) ReqPre() {
	p.Username = _username
	p.Password = _password
}

func (p *Predict) Method() string {
	return "predict"
}

func (r *ReportError) GetResponse() interface{} {
	return nil
}
func (r *ReportError) ReqPre() {
}

func (r *ReportError) Method() string {
	return "reporterror.json"
}

func SetUsername(username, password string) {
	_username = username
	_password = password
}

func Req(request Request) error {
	request.ReqPre()
	writer := new(bytes.Buffer)
	if err := json.NewEncoder(writer).Encode(request); err != nil {
		return err
	}
	resp, err := http.Post(fmt.Sprintf("%s%s", gateway, request.Method()), "application/json;charset=UTF-8", writer)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	response := request.GetResponse()
	if response == nil {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(response)
}
