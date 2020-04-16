package main

import (
	"encoding/json"
	"net/rpc"
	"strings"
)

type Codec struct {
	codec   rpc.ServerCodec
	request *rpc.Request
}

func (c *Codec) ReadRequestHeader(r *rpc.Request) error {
	c.request = r
	return c.codec.ReadRequestHeader(r)
}

func (c *Codec) ReadRequestBody(x interface{}) error {
	b, _ := json.Marshal(x)
	log.Debug("request", "->", c.request.ServiceMethod, "-", strings.TrimSpace(string(b)))
	return c.codec.ReadRequestBody(x)
}

func (c *Codec) WriteResponse(r *rpc.Response, x interface{}) error {
	if r.Error == "" {
		b, _ := json.Marshal(x)
		log.Debug("response", "->", r.ServiceMethod, "-", strings.TrimSpace(string(b)))
	} else {
		log.Debug("response with error", "->", r.ServiceMethod, "-", r.Error)
	}
	return c.codec.WriteResponse(r, x)
}

func (c *Codec) Close() error {
	return c.codec.Close()
}
