package logs

import "log/slog"

func ExampleSlog() {
	slog.Info("Starting server on port 8080")                              // want "log message should start with a lowercase English letter"
	slog.Info("starting server on port 8080")                              // ok
	slog.Error("Failed to connect to database")                            // want "log message should start with a lowercase English letter"
	slog.Error("failed to connect to database")                            // ok

	slog.Info("запуск сервера")                                            // want "log message contains non-ASCII characters" "log message should contain only English letters"
	slog.Error("ошибка подключения к базе данных")                         // want "log message contains non-ASCII characters" "log message should contain only English letters"
	slog.Info("starting server")                                           // ok
	slog.Error("failed to connect to database")                            // ok

	slog.Info("server started!🚀")                                         // want "log message contains disallowed special characters"
	slog.Error("connection failed!!!")                                     // want "log message contains disallowed special characters"
	slog.Warn("warning: something went wrong...")                          // want "log message contains disallowed special characters"
	slog.Info("server started")                                            // ok
	slog.Error("connection failed")                                        // ok
	slog.Warn("something went wrong")                                      // ok

	slog.Info("user password: 123")                                        // want "log message contains disallowed special characters" "log message contains potentially sensitive data"
	slog.Debug("api_key=secret")                                           // want "log message contains disallowed special characters" "log message contains potentially sensitive data"
	slog.Info("token: abc")                                                // want "log message contains disallowed special characters" "log message contains potentially sensitive data"
	slog.Info("user authenticated successfully")                           // ok
	slog.Debug("api request completed")                                    // ok
	slog.Info("token validated")                                           // want "log message contains potentially sensitive data"
}

