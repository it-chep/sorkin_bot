package sms_dto

type BasicRequest struct {
	To    []string `json:"to"`
	From  string   `json:"from"`
	Text  string   `json:"text"`
	Parts int      `json:"parts"`
}

type BasicResponse struct {
	Accepted bool   `json:"accepted"`
	To       string `json:"to"`
	Id       string `json:"id"`
}
