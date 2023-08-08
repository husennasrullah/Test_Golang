package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/go-stomp/stomp"
	"github.com/go-stomp/stomp/frame"
	"golang.org/x/net/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	urlPath := "ws://localhost:8091/nexsocket-notif2"

	loc, err := url.ParseRequestURI(urlPath)
	if err != nil {
		log.Fatal("gagal connect : ", err)
	}

	origin, err := url.ParseRequestURI("http://localhost:8091")
	if err != nil {
		log.Fatal("gagal connect : ", err)
	}

	wsConfig := new(websocket.Config)
	wsConfig.Location = loc
	wsConfig.Origin = origin
	wsConfig.Protocol = []string{"tcp"}
	wsConfig.Version = websocket.ProtocolVersionHybi13
	wsConfig.TlsConfig = &tls.Config{InsecureSkipVerify: true} //tlsConfig
	wsConfig.Header = make(map[string][]string)

	conn, err := websocket.DialConfig(wsConfig)
	if err != nil {
		log.Fatal("gagal connect : ", err)
	}

	err, token, cookie := login2()
	if err != nil {
		log.Fatal("gagal login : ", err)
	}

	var options = []func(*stomp.Conn) error{
		stomp.ConnOpt.Header("X-Token-Nexsoft", token),
		stomp.ConnOpt.Header("heart-beat", "25000,25000"),
		stomp.ConnOpt.Header("Cookie", "SOCKETSERVER="+cookie),
		stomp.ConnOpt.HeartBeatError(360 * time.Second),
	}

	sc, err := stomp.Connect(conn, options...)
	if err != nil {
		log.Fatal("error connect to socket server  ==> ", err)
	}

	sub, err := sc.Subscribe("/user/nexsoft/message", stomp.AckClient)
	if err != nil {
		log.Fatal("error subscribe : ", err)
	}

	for {
		select {
		case md := <-sub.C:
			if md.Err != nil {
				log.Fatal("receive greeting message caught error: %v", md.Err)
			}
			fmt.Printf("----> notif send: %v\n", string(md.Body))

			//unmarshal message ke model
			var message ReceiveMessage2
			_ = json.Unmarshal(md.Body, &message)

			switch message.Type {
			case "REQ_SYSADMIN_DATA":
				respSysadminData(sc, token, cookie)

			case "TAMBAH TYPE LAIN DISINI":
			}
		}
	}
}

func respSysadminData(sc *stomp.Conn, token string, cookie string) {
	// Open our jsonFile
	jsonFile, err := os.Open("respnd6.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened users.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	kirim := msg2{
		Touid:     "nexchief.dev",
		SrcSender: "nexpos.dev",
		Tipe:      "RESP_REQ_SYSADMIN_DATA",
		MsgID:     "672893104837",
		MsgSeq:    "1",
		Msg:       string(byteValue),
	}

	test, _ := json.Marshal(kirim)

	var sendOpt = []func(*frame.Frame) error{
		stomp.SendOpt.Header("X-Token-Nexsoft", token),
		stomp.SendOpt.Header("heart-beat", "25000,25000"),
		stomp.SendOpt.Header("Cookie", "SOCKETSERVER="+cookie),
	}

	time.Sleep(10 * time.Second)
	err = sc.Send("/nexsocket/senduser", "application/json", test, sendOpt...)
	if err != nil {
		log.Fatal("error kirim pesan : ", err)
	}
}

func login2() (error, string, string) {
	address := url.URL{
		Scheme: "http",
		Host:   "localhost:8091",
		Path:   "/nexsocket/login",
	}

	bodyRequest := SocketLogin2{
		User: "nexpos.dev",
		Pass: "Nexpos.dev19",
	}

	body, err := json.Marshal(bodyRequest)
	if err != nil {
		log.Fatal("gagal marshal body request login --> ", err)
	}

	headerRequest := make(map[string][]string)
	headerRequest["Content-Type"] = []string{"application/json"}

	request := &http.Request{
		Method: "POST",
		URL:    &address,
		Header: headerRequest,
		Body:   ioutil.NopCloser(strings.NewReader(string(body))),
	}

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	response, err := http.DefaultClient.Do(request)
	bodyResultByte, _ := ioutil.ReadAll(response.Body)

	defer func() {
		_ = response.Body.Close()
	}()

	var loginResponse map[string]interface{}
	_ = json.Unmarshal(bodyResultByte, &loginResponse)

	if response.StatusCode == 200 {
		return nil, loginResponse["token"].(string), loginResponse["cookie"].(string)
	} else {
		return err, "", ""
	}
}

//=====================================================================================

type ReceiveMessage2 struct {
	SndName      string `json:"sndName"`
	SndUID       string `json:"sndUID"`
	Msg          string `json:"msg"`
	MsgSeq       int32  `json:"msgSeq"`
	Type         string `json:"type"`
	SmsgID       string `json:"smsgID"`
	Source       string `json:"source"`
	MsgID        string `json:"msgID"`
	DeliveryDate string `json:"deliveryDate"`
}

type SocketLogin2 struct {
	User string `json:"uID"`
	Pass string `json:"pw"`
}

type msg2 struct {
	Touid     string      `json:"toUID"`
	SrcSender string      `json:"srcSender"`
	Tipe      string      `json:"type"`
	MsgID     string      `json:"msgID"`
	MsgSeq    string      `json:"msgSeq"`
	Msg       interface{} `json:"msg"`
}
