package zaplogs

import "go.uber.org/zap"

func ExampleZap() {
	logger := zap.NewExample()

	logger.Info("Starting zap logger")
	logger.Info("starting zap logger")

	logger.Error("Ошибка zap логгера")
	logger.Error("zap logger error")

	logger.Info("zap started!🔥")
	logger.Info("zap started successfully")

	logger.Info("user password: secret")
	logger.Info("user logged in")
}
