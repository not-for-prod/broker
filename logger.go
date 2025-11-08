package broker

type Logger interface {
	Error(msg string, keysValues ...any)
}
