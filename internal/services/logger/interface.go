package logger

type AppLogger interface {
	Info(msg string, args map[string]interface{})
	Error(msg string, args map[string]interface{})
	Debug(msg string, args map[string]interface{})
	Fatal(msg string, args map[string]interface{})
}
