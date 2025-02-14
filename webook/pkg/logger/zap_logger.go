package logger

import "go.uber.org/zap"

// 使用适配器模式,将zap.Logger适配为Logger
// 适配器模式:将一个类的接口转换成客户希望的另外一个接口。
type ZapLogger struct {
	l *zap.Logger
}

func NewZapLogger(l *zap.Logger) *ZapLogger {
	return &ZapLogger{l: l}
}

func (z *ZapLogger) Debug(msg string, args ...Field) {
	z.l.Debug(msg, z.toArgs(args)...)
}

func (z *ZapLogger) Info(msg string, args ...Field) {
	z.l.Info(msg, z.toArgs(args)...)
}

func (z *ZapLogger) Warn(msg string, args ...Field) {
	z.l.Warn(msg, z.toArgs(args)...)
}

func (z *ZapLogger) Error(msg string, args ...Field) {
	z.l.Error(msg, z.toArgs(args)...)
}

// toArgs 将Field转换为zap.Field
// toArgs 将Field转换为zap.Field
// 参数:
//   - args: 可变参数,Field类型的参数列表
//
// 返回:
//   - []zap.Field: 转换后的zap.Field切片
func (z *ZapLogger) toArgs(args []Field) []zap.Field {
	fields := make([]zap.Field, 0, len(args))
	for _, arg := range args {
		fields = append(fields, zap.Any(arg.Key, arg.Value))
	}
	return fields
}
