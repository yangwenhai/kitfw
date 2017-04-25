package service

import (
	"context"
	protocol "framework/sg/protocol"
)

type BaseHandler interface {
	Process(context.Context, []byte) ([]byte, error)
	doProcess(ctx context.Context)
}

type NewHandlerFunc func() BaseHandler

var HandlerMap = map[int32]NewHandlerFunc{
	protocol.PROTOCOL_SUM_REQUEST:    func() BaseHandler { return NewSumHandler() },
	protocol.PROTOCOL_CONCAT_REQUEST: func() BaseHandler { return NewConcatHandler() },
}

func GetHandler(ProtoId int32) BaseHandler {
	if HandlerMap[ProtoId] != nil {
		return HandlerMap[ProtoId]()
	}
	return nil
}

type SumHandler struct {
	request *protocol.SumRequest
	reply   *protocol.SumReply
}

func NewSumHandler() *SumHandler {
	request := &protocol.SumRequest{}
	reply := &protocol.SumReply{}
	return &SumHandler{request, reply}
}

func (handler *SumHandler) Process(ctx context.Context, payload []byte) ([]byte, error) {
	err := protocol.Decode(handler.request, payload)
	if err != nil {
		return nil, err
	}
	handler.doProcess(ctx)
	ret, err := protocol.Encode(handler.reply)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

type ConcatHandler struct {
	request *protocol.ConcatRequest
	reply   *protocol.ConcatReply
}

func NewConcatHandler() *ConcatHandler {
	request := &protocol.ConcatRequest{}
	reply := &protocol.ConcatReply{}
	return &ConcatHandler{request, reply}
}

func (handler *ConcatHandler) Process(ctx context.Context, payload []byte) ([]byte, error) {
	err := protocol.Decode(handler.request, payload)
	if err != nil {
		return nil, err
	}
	handler.doProcess(ctx)
	ret, err := protocol.Encode(handler.reply)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
