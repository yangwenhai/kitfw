package service

import (
	"fmt"
	logger "framework/sg/log"

	"context"
)

func (handler *ConcatHandler) doProcess(ctx context.Context) {
	handler.reply.RetCode = 0
	handler.reply.Val = handler.request.Str1 + handler.request.Str2
	rlogger := ctx.Value("logger").(*logger.Logger)
	rlogger.Info("request", fmt.Sprintf("%s %s", handler.request.Str1, handler.request.Str2))
}
