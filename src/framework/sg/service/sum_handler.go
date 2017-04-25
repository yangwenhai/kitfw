package service

import (
	"fmt"
	logger "framework/sg/log"

	"context"
)

func (handler *SumHandler) doProcess(ctx context.Context) {
	handler.reply.RetCode = 0
	handler.reply.Val = handler.request.Num1 + handler.request.Num2
	rlogger := ctx.Value("logger").(*logger.Logger)
	rlogger.Info("request", fmt.Sprintf("%d %d", handler.request.Num1, handler.request.Num2))
}
