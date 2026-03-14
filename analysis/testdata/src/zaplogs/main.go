package zaplogs

import "go.uber.org/zap"

func ExampleZap() {
	logger := zap.NewExample()

	logger.Info("Starting zap logger")             // want "log message should start with a lowercase English letter"
	logger.Info("starting zap logger")             // ok

	logger.Error("Ошибка zap логгера")             // want "log message contains non-ASCII characters" "log message should contain only English letters"
	logger.Error("zap logger error")               // ok

	logger.Info("zap started!🔥")                  // want "log message contains disallowed special characters"
	logger.Info("zap started successfully")        // ok

	logger.Info("user password: secret")           // want "log message contains disallowed special characters" "log message contains potentially sensitive data"
	logger.Info("user logged in")                  // ok
}

