package utils

import (
	"sync"

	"go.uber.org/zap"
)

var (
	logger     *zap.SugaredLogger
	loggerOnce sync.Once
)

// InitLogger tạo logger dùng zap; gọi nhiều lần vẫn trả về cùng instance.
func InitLogger() *zap.SugaredLogger {
	loggerOnce.Do(func() {
		base, err := zap.NewProduction()
		if err != nil {
			base = zap.NewExample()
		}
		logger = base.Sugar()
	})
	return logger
}

// L trả về logger đã init (lazy).
func L() *zap.SugaredLogger {
	if logger == nil {
		return InitLogger()
	}
	return logger
}
