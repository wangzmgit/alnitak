package logger

import (
	"os"

	"github.com/natefinch/lumberjack"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 1 定义一下logger使用的常量
var (
	mode       string
	filename   string
	level      zapcore.Level
	maxSize    int
	maxAge     int
	maxBackups int
)

// 2 初始化Logger对象
func InitLogger() (err error) {
	// 读取配置
	mode = viper.GetString("log.mode")         //开发模式
	filename = viper.GetString("log.filename") // 日志存放路径
	//level    = viper.GetString("log.level")    // 日志级别
	level = zapcore.DebugLevel                   // 日志级别
	maxSize = viper.GetInt("log.max-size")       //最大存储大小
	maxAge = viper.GetInt("log.max-age")         //最大存储时间
	maxBackups = viper.GetInt("log.max-backups") //#备份数量

	// 创建Core三大件，进行初始化
	writeSyncer := getLogWriter(filename, maxSize, maxAge, maxBackups)
	encoder := getEncoder()
	// 创建核心-->如果是dev模式，就在控制台和文件都打印，否则就只写到文件中
	var core zapcore.Core
	if mode == "dev" {
		// 开发模式，日志输出到终端
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		// NewTee创建一个核心，将日志条目复制到两个或多个底层核心中。
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, writeSyncer, level),
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), level),
		)
	} else {
		core = zapcore.NewCore(encoder, writeSyncer, level)
	}

	//core := zapcore.NewCore(encoder, writeSyncer, level)
	// 创建 logger 对象
	log := zap.New(core, zap.AddCaller())
	// 替换全局的 logger, 后续在其他包中只需使用zap.L()调用即可
	zap.ReplaceGlobals(log)
	return
}

// 获取Encoder，给初始化logger使用的
func getEncoder() zapcore.Encoder {
	// 使用zap提供的 NewProductionEncoderConfig
	encoderConfig := zap.NewProductionEncoderConfig()
	// 设置时间格式
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// 时间的key
	encoderConfig.TimeKey = "time"
	// 级别
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// 显示调用者信息
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	// 返回json 格式的 日志编辑器
	return zapcore.NewJSONEncoder(encoderConfig)
}

// 获取切割的问题，给初始化logger使用的
func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	// 使用 lumberjack 归档切片日志
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
	}
	return zapcore.AddSync(lumberJackLogger)
}
