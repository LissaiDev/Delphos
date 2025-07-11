package hermes

import (
	"io"
)

type Service string
type Method int

const (
	MethodGet Method = iota
	MethodPost
	MethodPut
	MethodDelete
	MethodPatch
	MethodOptions
)

var SERVICES = map[Service]string{
	Service("DISCORD"): "https://discord.com/api/webhooks",
}

type Response struct {
	Code    int    `json:"code"`
	Success bool   `json:"success"`
	Data    []byte `json:"data"`
}

type Request struct {
	Service      Service            `json:"service"`
	Url          string             `json:"url"`
	Method       Method             `json:"method"`
	Headers      *map[string]string `json:"headers"`
	Body         *map[string]any    `json:"body"`
	ResolvedBody io.Reader
}

type Fetcher interface {
	Fetch(req *Request) Response
	Get(service Service, url string, headers *map[string]string) Response
	Post(service Service, url string, body *map[string]any, headers *map[string]string) Response
	Put(service Service, url string, body *map[string]any, headers *map[string]string) Response
	Delete(service Service, url string, headers *map[string]string) Response
	Patch(service Service, url string, body *map[string]any, headers *map[string]string) Response
}
