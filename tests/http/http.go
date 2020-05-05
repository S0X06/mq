package main

import (
	"encoding/json"
	"fmt"
	"mq/handler"

	"mq/model"
	"mq/utils"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Data struct {
	OrderSn string
	Data    string
}

var (
	appId  = "123"
	notify = "http://localhost:8888/"

	method = "POST"
	tryUrl = "http://127.0.0.1:8080/try"

	ackUrl    = "http://127.0.0.1:8080/ack"
	ackMethod = "PUT"
)

func RequstCleint() {

	for {

		orderSn := utils.GenValidateCode(15)

		data := Data{
			OrderSn: orderSn,
			Data:    utils.GenValidateCode(5),
		}
		jsonData, err := json.Marshal(data)
		if err != nil {
			continue
		}

		params := &map[string]interface{}{
			"app_id":   appId,
			"order_sn": orderSn,
			"data":     jsonData,
			"notify":   notify,
		}
		body, err := utils.RequstCleint(method, tryUrl, params)
		if err != nil {
			fmt.Println(err)
			continue
		}

		resp := &handler.Response{}
		err = json.Unmarshal(body, resp)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if resp.Code != 0 {
			continue
		}

		ackParams := &map[string]interface{}{
			"ack":      2,
			"order_sn": orderSn,
		}

		ackBody, err := utils.RequstCleint(ackMethod, ackUrl, ackParams)
		if err != nil {
			continue
		}
		ackResp := &handler.Response{}
		err = json.Unmarshal(ackBody, ackResp)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Printf("code: %v,message:%v\n", ackResp.Code, ackResp.Message)

	}
}

func Resp(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		// w.Write()
	}
	order_sn := r.Form["order_sn"][0]
	fmt.Println(order_sn)
	ack := model.Ack{
		OrderSn: order_sn,
		Ack:     2,
	}
	if err != nil {
		ack.Ack = 1
	}

	params, err := json.Marshal(ack)

	w.Write(params)

}

func main() {

	go RequstCleint()

	http.HandleFunc("/", Resp)

	http.ListenAndServe(":8888", nil)
	//平滑重启
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}

}
