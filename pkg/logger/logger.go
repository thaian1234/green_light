package logger

import (
	"sync"

	"github.com/thaian1234/green_light/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	instance *Adapter
	once     sync.Once
)

type Adapter struct {
	logger *zap.SugaredLogger
}

func NewAdapter(cfg *config.Logger) (*Adapter, error) {

	writerSyncer := getLogWritter(cfg)
	encoder := getEncoder()

	level, err := zapcore.ParseLevel(cfg.LogLevel)
	if err != nil {
		return nil, err
	}

	core := zapcore.NewCore(encoder, writerSyncer, level)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return &Adapter{
		logger: logger.Sugar(),
	}, nil
}

func Initialize(cfg *config.Logger) error {
	var err error
	once.Do(func() {
		instance, err = NewAdapter(cfg)
	})
	return err
}

func GetLogger() *Adapter {
	return instance
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWritter(cfg *config.Logger) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   cfg.LogPath,
		MaxSize:    cfg.LogMaxSize,  // megabytes
		MaxBackups: cfg.LogBackUps,  // number of backups
		MaxAge:     cfg.LogMaxAge,   // days
		Compress:   cfg.LogCompress, // compress the backups
	}

	return zapcore.AddSync(lumberJackLogger)
}

func Debug(msg string, fields ...interface{}) {
	GetLogger().logger.Debugw(msg, fields...)
}

func Info(msg string, fields ...interface{}) {
	GetLogger().logger.Infow(msg, fields...)
}

func Warn(msg string, fields ...interface{}) {
	GetLogger().logger.Warnw(msg, fields...)
}

func Error(msg string, fields ...interface{}) {
	GetLogger().logger.Errorw(msg, fields...)
}

func Fatal(msg string, fields ...interface{}) {
	GetLogger().logger.Fatalw(msg, fields...)
}
