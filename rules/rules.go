package rules

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/hikingpig/twitter-stream/request"
)

const (
	rulesUrl = "https://api.twitter.com/2/tweets/search/stream/rules"
)

type Rule struct {
	Value string `json:"value"`
	Tag   string `json:"tag"`
}

type Rules struct {
	Add []*Rule `json:"add"`
}

func (r *Rules) AddRule(value, tag string) {
	r.Add = append(r.Add, &Rule{value, tag})
}

type RulesResponse struct {
	Data   []RespData
	Meta   []RespMeta
	Errors []RespError
}

type RespData struct {
	Value string `json:"Value"`
	Tag   string `json:"Tag"`
	Id    string `json:"id"`
}

type RespMeta struct {
	Sent    string      `json:"sent"`
	Summary MetaSummary `json:"summary"`
}

type MetaSummary struct {
	Created    uint `json:"created"`
	NotCreated uint `json:"not_created"`
}

type RespError struct {
	Value string `json:"Value"`
	Id    string `json:"id"`
	Title string `json:"title"`
	Type  string `json:"type"`
}

func Create(rules *Rules) (*RulesResponse, error) {
	data, err := json.Marshal(rules)
	if err != nil {
		return nil, err
	}
	res, err := request.Request(http.MethodPost, rulesUrl, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	rulesResp := &RulesResponse{}
	err = json.NewDecoder(res.Body).Decode(rulesResp)
	if err != nil {
		return nil, err
	}
	return rulesResp, err
}
