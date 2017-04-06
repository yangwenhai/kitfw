package service

import (
	"fmt"

	"golang.org/x/net/context"
)

func (handler *SumHandler) doProcess(ctx context.Context) {
	handler.reply.RetCode = 0
	handler.reply.Val = handler.request.Num1 + handler.request.Num2
	handler.logger.Log("request", fmt.Sprintf("%s %s", handler.request.Num1, handler.request.Num2))
}
