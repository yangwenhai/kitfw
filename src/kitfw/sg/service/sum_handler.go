package service

import (
	"fmt"
	logger "kitfw/sg/log"

	"context"
)

func (handler *SumHandler) doProcess(ctx context.Context) {
	handler.reply.RetCode = 0
	handler.reply.Val = handler.request.Num1 + handler.request.Num2
	rlogger := ctx.Value("logger").(*logger.Logger)
	rlogger.Info("request", fmt.Sprintf("%s %s", handler.request.Num1, handler.request.Num2))
}
