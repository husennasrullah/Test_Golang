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
	"strings"
	"time"
)

func main() {
	urlPath := "wss://10.10.11.223/nexsocket-notif2"

	loc, err := url.ParseRequestURI(urlPath)
	if err != nil {
		log.Fatal("gagal connect : ", err)
	}

	origin, err := url.ParseRequestURI("https://10.10.11.223")
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

	err, token, cookie := login()
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

	sub, err := sc.Subscribe("/user/nexsoft/notifsend", stomp.AckClient)
	if err != nil {
		log.Fatal("error subscribe : ", err)
	}

	sub2, err := sc.Subscribe("/user/nexsoft/status", stomp.AckClient)
	if err != nil {
		log.Fatal("error subscribe pong : ", err)
	}

	sub3, err := sc.Subscribe("/user/nexsoft/receive", stomp.AckClient)
	if err != nil {
		log.Fatal("error subscribe pong : ", err)
	}

	//buat message untuk dikirim ke socket server
	sendMsg := createMessage()
	a, _ := json.Marshal(sendMsg)

	//header untuk send message
	var sendOpt = []func(*frame.Frame) error{
		stomp.SendOpt.Header("X-Token-Nexsoft", token),
		stomp.SendOpt.Header("heart-beat", "25000,25000"),
		stomp.SendOpt.Header("Cookie", "SOCKETSERVER="+cookie),
	}

	//fungsi stomp untuk mengirim pesan ke socket server
	err = sc.Send("/nexsocket/senduser", "application/json", a, sendOpt...)
	if err != nil {
		log.Fatal("error kirim pesan : ", err)
	}

	for {
		select {
		case md := <-sub.C:
			if md.Err != nil {
				log.Fatal("receive greeting message caught error: %v", md.Err)
			}
			fmt.Printf("----> notif send: %v\n", string(md.Body))
		case md2 := <-sub2.C:
			if md2.Err != nil {
				log.Fatal("receive greeting message caught error: %v", md2.Err)
			}
			fmt.Printf("----> status message : %v\n", string(md2.Body))

		case md3 := <-sub3.C:
			if md3.Err != nil {
				log.Fatal("receive greeting message caught error: %v", md3.Err)
			}
			fmt.Printf("----> Received Message : %v\n", string(md3.Body))

			//unmarshal message ke model
			var message ReceiveMessage
			_ = json.Unmarshal(md3.Body, &message)

			//send delivered messaged to socket to inform that the message already received by nexchief
			var sendOpt = []func(*frame.Frame) error{
				stomp.SendOpt.Header("X-Token-Nexsoft", token),
				stomp.SendOpt.Header("heart-beat", "25000,25000"),
				stomp.SendOpt.Header("Cookie", "SOCKETSERVER="+cookie),
			}

			bodyMessage := DeliveredMessage{
				SmsgID: message.SmsgID,
			}

			sendBody, _ := json.Marshal(bodyMessage)

			err = sc.Send("/nexsocket/delivMsg", "application/json;charset=UTF-8", sendBody, sendOpt...)
			if err != nil {
				fmt.Println("error send to destination /nexsocket/delivMsg")
				//break loopReceive
			}

		}
	}

}

func login() (error, string, string) {
	address := url.URL{
		Scheme: "https",
		Host:   "10.10.11.223",
		Path:   "/nexsocket/login",
	}

	bodyRequest := SocketLogin{
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

func createMessage() (result msg) {
	msgs := SocketMessage{NexsoftMessage: NexsoftMessage{
		Header: Header{
			MessageID:   "123",
			UserID:      "nexpos.dev",
			Password:    "",
			Version:     "1",
			PrincipalID: "no",
			Timestamp:   "2020/09/19",
			Action:      Action{},
		},
		Payload: Payload{
			Header: payloadHeader{
				Status:    0,
				Size:      0,
				Range:     rangeData{},
				ReqStatus: "SUCCESS",
				Message:   "SUCCESS",
			},
			Data:
			[]ProductFromSocket{
				{
					CompanyID:           "5283512",
					ProductCode:         "SOCKET3",
					ParentCode:          "",
					ProductName:         "Biskuit Roma",
					ProductPackaging:    "Kardus",
					ProductDescription:  "Ini Biskuit Roma",
					PgLevel1ID:          "11",
					PgLevel1Name:        "",
					PgLevel2ID:          "1101",
					PgLevel2Name:        "",
					PgLevel3ID:          "110101",
					PgLevel3Name:        "",
					BrandID:             "BR0051",
					BrandName:           "testsocket",
					ProductCategoryID:   "STANDARD",
					ProductCategoryname: "",
					MarketLaunchDate:    "2014-05-15",
					ProductStatus:       "A",
					AdditionalMargin:    "",
					BuyingPriceUom1:     10000.0,
					SellingPriceUom1:    11000.0,
					Uom1:                "KRT",
					Uom2:                "KRT",
					Uom3:                "KRT",
					Uom4:                "KRT",
					Conversion1to4:      1,
					Conversion2to4:      1,
					Conversion3to4:      1,
					WeightPerUnit:       0.0,
					VolumePerUnit:       0.0,
					CaseWidth:           0.0,
					CaseHeight:          0.0,
					CaseDepth:           0.0,
					CaseWeight:          0.0,
					IsProductWithVarian: "N",
					IsProductWithBatch:  "N",
					IsInventoryItem:     "Y",
					IsFreeGoodItem:      "N",
					IsEmbalaceItem:      "Y",
					IsTax1Applied:       "N",
					IsTax2Applied:       "N",
					IsTax3Applied:       "N",
					ISTaxSubsidized:     "N",
					IsAlocatedItem:      "N",
					IsNoPurchaseReturn:  "N",
					IsNoSalesReturn:     "N",
					IsNoRegularDiscount: "N",
					AllowSellUnderCost:  "N",
					ProductClass:        "1",
					ProductSortSequence: 9000,
					MinQtyPerLines:      "N",
					ReportColumnHeader1: "",
					ReportColumnHeader2: "",
					BufferStockLevel:    10.0,
					WeekRPPOB:           13,
					BarcodeListForUom1:  "18993278320217",
					BarcodeListForUom2:  "",
					BarcodeListForUom3:  "",
					BarcodeListForUom4:  "",
					ProductCreated:      "2020-12-09 16:42:25",
					ProductModified:     "2020-12-09 16:42:25",
					ProductModifiedBy:   "import",
					ProductActivation:   "import",
				},
			},
		},
	}}

	stringMsg, err := json.Marshal(msgs)
	if err != nil {
		log.Fatal("Gagal Unmarshal Model --> ", err)
	}

	result = msg{
		Touid:     "nexchief.dev",
		SrcSender: "id.co.nexsoft.nexmile",
		Tipe:      "POST_PRODUCT_DATAAAA",
		MsgID:     "123456",
		MsgSeq:    "1",
		Msg:       string(stringMsg),
	}

	return
}

//STRUCK PENDUKUNG
type ReceiveMessage struct {
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

type DeliveredMessage struct {
	SmsgID string `json:"smsgID"`
}

type SocketLogin struct {
	User string `json:"uID"`
	Pass string `json:"pw"`
}

type msg struct {
	Touid     string      `json:"toUID"`
	SrcSender string      `json:"srcSender"`
	Tipe      string      `json:"type"`
	MsgID     string      `json:"msgID"`
	MsgSeq    string      `json:"msgSeq"`
	Msg       interface{} `json:"msg"`
}

type SocketMessage struct {
	NexsoftMessage NexsoftMessage `json:"nexsoft_message"`
}

type NexsoftMessage struct {
	Header  Header  `json:"header"`
	Payload Payload `json:"payload"`
}

type Header struct {
	MessageID   string `json:"message_id"`
	UserID      string `json:"user_id"`
	Password    string `json:"password"`
	Version     string `json:"version"`
	PrincipalID string `json:"principal_id"`
	Timestamp   string `json:"timestamp"`
	Action      Action `json:"action"`
}

type Action struct {
	ClassName string `json:"class_name"`
	TypeName  string `json:"type_name"`
}

type Payload struct {
	Header payloadHeader `json:"header"`
	Data   interface{}   `json:"data"`
	//[]SalesmanFromSocket `json:"data"`
}

type payloadHeader struct {
	Status    int       `json:"status"`
	Size      int       `json:"size"`
	Range     rangeData `json:"range"`
	ReqStatus string    `json:"reqStatus"`
	Message   string    `json:"message"`
}

type rangeData struct {
	From int `json:"from"`
	To   int `json:"to"`
}

type SalesmanFromSocket struct {
	CompanyID                string    `json:"companyID"`
	BranchID                 string    `json:"branchID"`
	SalesmanID               string    `json:"salesmanID"`
	SalesmanName             string    `json:"salesmanName"`
	SalesmanAddress1         string    `json:"salesmanAddress1"`
	SalesmanAddress2         string    `json:"salesmanAddress2"`
	SalesmanCity             string    `json:"salesmanCity"`
	SalesmanPhone            string    `json:"salesmanPhone"`
	SalesmanEmail            string    `json:"salesmanEmail"`
	SalesmanKTP              string    `json:"salesmanKTP"`
	SalesmanBirthDate        string    `json:"salesmanBirthDate"`
	SalesmanJoinDate         time.Time `json:"salesmanJoinDate"`
	SupervisorID             string    `json:"supervisorID"`
	SalescoID                string    `json:"salescoID"`
	SalesmanType             string    `json:"salesmanType"`
	SalesmanClass            string    `json:"salesman_class"`
	ZoneID                   string    `json:"zoneID"`
	SalesmanCategory         string    `json:"salesmanCategory"`
	InvoiceTradeToCustomerID string    `json:"invoiceTradeToCustomerID"`
	SalesmanModified         string    `json:"salesmanModified"`
}

type ProductFromSocket struct {
	CompanyID           string  `json:"companyID"`
	ProductCode         string  `json:"productCode"`
	ParentCode          string  `json:"parentCode"`
	ProductName         string  `json:"productName"`
	ProductPackaging    string  `json:"productPackaging"`
	ProductDescription  string  `json:"productDescription"`
	PgLevel1ID          string  `json:"pgLevel1ID"`
	PgLevel1Name        string  `json:"pgLevel1Name"`
	PgLevel2ID          string  `json:"pgLevel2ID"`
	PgLevel2Name        string  `json:"pgLevel2Name"`
	PgLevel3ID          string  `json:"pgLevel3ID"`
	PgLevel3Name        string  `json:"pgLevel3Name"`
	BrandID             string  `json:"brandID"`
	BrandName           string  `json:"brand_name"`
	ProductCategoryID   string  `json:"productCategoryID"`
	ProductCategoryname string  `json:"productCategoryName"`
	MarketLaunchDate    string  `json:"marketLaunchDate"`
	ProductStatus       string  `json:"productStatus"`
	AdditionalMargin    string  `json:"additionalMargin"`
	BuyingPriceUom1     float64 `json:"buyingPriceUom1"`
	SellingPriceUom1    float64 `json:"sellingPriceUom1"`
	Uom1                string  `json:"uom1"`
	Uom2                string  `json:"uom2"`
	Uom3                string  `json:"uom3"`
	Uom4                string  `json:"uom4"`
	Conversion1to4      int32   `json:"conversion1to4"`
	Conversion2to4      int32   `json:"conversion2to4"`
	Conversion3to4      int32   `json:"conversion3to4"`
	WeightPerUnit       float64 `json:"weightPerUnit"`
	VolumePerUnit       float64 `json:"volumePerUnit"`
	CaseWidth           float64 `json:"caseWidth"`
	CaseHeight          float64 `json:"caseHeight"`
	CaseDepth           float64 `json:"caseDepth"`
	CaseWeight          float64 `json:"caseWeight"`
	IsProductWithVarian string  `json:"isProductWithVarian"`
	IsProductWithBatch  string  `json:"isProductWithBatch"`
	IsInventoryItem     string  `json:"isInventoryItem"`
	IsFreeGoodItem      string  `json:"isFreeGoodItem"`
	IsEmbalaceItem      string  `json:"isEmbalaceItem"`
	IsTax1Applied       string  `json:"isTax1Applied"`
	IsTax2Applied       string  `json:"isTax2Applied"`
	IsTax3Applied       string  `json:"isTax3Applied"`
	ISTaxSubsidized     string  `json:"isTaxSubsidized"`
	IsAlocatedItem      string  `json:"isAlocatedItem"`
	IsNoPurchaseReturn  string  `json:"isNoPurchaseReturn"`
	IsNoSalesReturn     string  `json:"is_no_sales_return"`
	IsNoRegularDiscount string  `json:"isNoRegularDiscount"`
	AllowSellUnderCost  string  `json:"allowSellUnderCost"`
	ProductClass        string  `json:"productClass"`
	ProductSortSequence int64  `json:"productSortSequence"`
	MinQtyPerLines      string  `json:"minQtyPerLines"`
	ReportColumnHeader1 string  `json:"reportColumnHeader1"`
	ReportColumnHeader2 string  `json:"reportColumnHeader2"`
	BufferStockLevel    float64  `json:"bufferStockLevel"`
	WeekRPPOB           int64  `json:"weekRPPOB"`
	BarcodeListForUom1  string  `json:"barcodeListForUom1"`
	BarcodeListForUom2  string  `json:"barcodeListForUom2"`
	BarcodeListForUom3  string  `json:"barcodeListForUom3"`
	BarcodeListForUom4  string  `json:"barcodeListForUom4"`
	ProductCreated      string  `json:"productCreated"`
	ProductModified     string  `json:"productModified"`
	ProductModifiedBy   string  `json:"productModifiedBy"`
	ProductActivation   string  `json:"productActivation"`
}
