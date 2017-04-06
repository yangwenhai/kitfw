package service

import (
	"fmt"

	"golang.org/x/net/context"
)

func (handler *ConcatHandler) doProcess(ctx context.Context) {
	handler.reply.RetCode = 0
	handler.reply.Val = handler.request.Str1 + handler.request.Str2
	handler.logger.Log("request", fmt.Sprintf("%s %s", handler.request.Str1, handler.request.Str2))
}
