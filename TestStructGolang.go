package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func main() {
	msgs := SocketMessage1{NexsoftMessage: NexsoftMessage1{
		Header: Header1{
			MessageID:   "12897183709jkabskjbas",
			UserID:      "nexchief.dev",
			Password:    "",
			Version:     "1",
			PrincipalID: "NDI",
			Timestamp:   "2022/08/23",
			Action:      Action1{
				ClassName: "",
				TypeName:  "",
			},
		},
		Payload: Payload1{
			Header: payloadHeader1{
				Status:    0,
				Size:      0,
				Range:     rangeData1{
					From: 0,
					To:   0,
				},
				ReqStatus: "SUCCESS",
				Message:   "SUCCESS",
			},
			Data: datakirim{
				CompanyID: "NS6044050001548",
				BranchID:  "1522910433080",
				Principal: "NDI",
			},
		},
	}}

	stringMsg, err := json.Marshal(msgs)
	if err != nil {
		log.Fatal("Gagal Unmarshal Model --> ", err)
	}

	fmt.Println(string(stringMsg))
}

type SocketMessage1 struct {
	NexsoftMessage NexsoftMessage1 `json:"nexsoft_message"`
}

type NexsoftMessage1 struct {
	Header  Header1  `json:"header"`
	Payload Payload1`json:"payload"`
}

type Header1 struct {
	MessageID   string `json:"message_id"`
	UserID      string `json:"user_id"`
	Password    string `json:"password"`
	Version     string `json:"version"`
	PrincipalID string `json:"principal_id"`
	Timestamp   string `json:"timestamp"`
	Action      Action1 `json:"action"`
}

type Action1 struct {
	ClassName string `json:"class_name"`
	TypeName  string `json:"type_name"`
}

type Payload1 struct {
	Header payloadHeader1 `json:"header"`
	Data   interface{}   `json:"data"`
	//[]SalesmanFromSocket `json:"data"`
}

type payloadHeader1 struct {
	Status    int       `json:"status"`
	Size      int       `json:"size"`
	Range     rangeData1 `json:"range"`
	ReqStatus string    `json:"reqStatus"`
	Message   string    `json:"message"`
}

type rangeData1 struct {
	From int `json:"from"`
	To   int `json:"to"`
}

type datakirim struct {
	CompanyID string `json:"companyId"`
	BranchID string `json:"branchId"`
	Principal string `json:"principalId"`
}