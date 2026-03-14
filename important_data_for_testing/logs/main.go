package logs

import "log/slog"

func ExampleSlog() {
	slog.Info("Starting server on port 8080")
	slog.Info("starting server on port 8080")
	slog.Error("Failed to connect to database")
	slog.Error("failed to connect to database")

	slog.Info("запуск сервера")
	slog.Error("ошибка подключения к базе данных")
	slog.Info("starting server")
	slog.Error("failed to connect to database")

	slog.Info("server started!🚀")
	slog.Error("connection failed!!!")
	slog.Warn("warning: something went wrong...")
	slog.Info("server started")
	slog.Error("connection failed")
	slog.Warn("something went wrong")

	slog.Info("user password: 123")
	slog.Debug("api_key=secret")
	slog.Info("token: abc")
	slog.Info("user authenticated successfully")
	slog.Debug("api request completed")
	slog.Info("token validated")
}
