package utils

import "go.uber.org/zap"

func ErrorLog(msg, module, err string) {
	zap.L().Error(msg, zap.String("module", module), zap.String("err", err))
}
