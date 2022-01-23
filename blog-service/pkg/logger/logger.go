package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"runtime"
	"time"
)

//日志分级文件

type Level int8

type Fields map[string]interface{}

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
	LevelPanic
)

//日志实例化
type Logger struct {
	newLogger *log.Logger
	ctx       context.Context
	level     Level
	fields    Fields
	callers   []string
}

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	case LevelFatal:
		return "fatal"
	case LevelPanic:
		return "panic"
	}
	return ""
}

func NewLogger(w io.Writer, prefix string, flag int) *Logger {
	l := log.New(w, prefix, flag)
	return &Logger{newLogger: l}
}

func (l *Logger) clone() *Logger {
	nl := *l
	return &nl
}

//设置日志等级
func (l *Logger) WithLevel(lvl Level) *Logger {
	l1 := l.clone()
	l1.level = lvl
	return l1
}

//设置日志公共的字段
func (l *Logger) WithFields(f Fields) *Logger {
	l1 := l.clone()
	if l1.fields == nil {
		l1.fields = make(Fields)
	}
	for k, v := range f {
		l1.fields[k] = v
	}
	return l1
}

//设置上下文属性
func (l *Logger) WithContext(ctx context.Context) *Logger {
	l1 := l.clone()
	l1.ctx = ctx
	return l1
}

//设置当前某一层调用栈的信息
func (l *Logger) WithCaller(skip int) *Logger {
	l1 := l.clone()
	pc, file, line, ok := runtime.Caller(skip)
	if ok {
		f := runtime.FuncForPC(pc)
		l1.callers = []string{fmt.Sprintf("%s: %d %s", file, line, f.Name())}
	}
	return l1
}

//设置当前的整个调用栈信息
func (l *Logger) WithCallersFrames() *Logger {
	maxCallerDepth := 25
	minCallerDepth := 1
	callers := []string{}
	pcs := make([]uintptr, maxCallerDepth)
	depth := runtime.Callers(minCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])
	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		callers = append(callers, fmt.Sprintf("%s: %d %s", frame.File, frame.Line, frame.Function))
		if !more {
			break
		}
	}

	l1 := l.clone()
	l1.callers = callers
	return l1
}

// //链路追踪
// func (l *Logger) WithTrace() *Logger {
// 	ginCtx, ok := l.ctx.(*gin.Context)
// 	if ok {
// 		return l.WithFields(Fields{
// 			"trace_id": ginCtx.MustGet("X-Trace-ID"),
// 			"span_id":  ginCtx.MustGet("X-Span-ID"),
// 		})
// 	}
// 	return l
// }

//日志内容的格式化
func (l *Logger) JSONFormat(message string) map[string]interface{} {
	data := make(Fields, len(l.fields)+4)
	data["level"] = l.level.String()
	data["time"] = time.Now().Local().UnixNano()
	data["message"] = message
	data["callers"] = l.callers
	if len(l.fields) > 0 {
		for k, v := range l.fields {
			if _, ok := data[k]; !ok {
				data[k] = v
			}
		}
	}
	return data
}

//日志内容输出
func (l *Logger) Output(message string) {
	body, _ := json.Marshal(l.JSONFormat(message))
	content := string(body)
	switch l.level {
	case LevelDebug:
		l.newLogger.Print(content)
	case LevelInfo:
		l.newLogger.Print(content)
	case LevelWarn:
		l.newLogger.Print(content)
	case LevelError:
		l.newLogger.Print(content)
	case LevelFatal:
		l.newLogger.Print(content)
	case LevelPanic:
		l.newLogger.Print(content)
	}
}

//日志分级输出,根据6个等级来编写
func (l *Logger) Debug(ctx context.Context, v ...interface{}) {
	l.WithLevel(LevelDebug).Output(fmt.Sprint(v...))
}

func (l *Logger) Debugf(ctx context.Context, format string, v ...interface{}) {
	l.WithLevel(LevelDebug).Output(fmt.Sprintf(format, v...))
}

func (l *Logger) Info(ctx context.Context, v ...interface{}) {
	l.WithLevel(LevelInfo).Output(fmt.Sprint(v...))
}

func (l *Logger) Infof(ctx context.Context, format string, v ...interface{}) {
	l.WithLevel(LevelInfo).Output(fmt.Sprintf(format, v...))
}

func (l *Logger) Warn(ctx context.Context, v ...interface{}) {
	l.WithLevel(LevelWarn).Output(fmt.Sprint(v...))
}

func (l *Logger) Warnf(ctx context.Context, format string, v ...interface{}) {
	l.WithLevel(LevelWarn).Output(fmt.Sprintf(format, v...))
}

func (l *Logger) Error(ctx context.Context, v ...interface{}) {
	l.WithLevel(LevelError).Output(fmt.Sprint(v...))
}

func (l *Logger) Errorf(ctx context.Context, format string, v ...interface{}) {
	l.WithLevel(LevelError).Output(fmt.Sprintf(format, v...))
}

func (l *Logger) Fatal(ctx context.Context, v ...interface{}) {
	l.WithLevel(LevelFatal).Output(fmt.Sprint(v...))
}

func (l *Logger) Fatalf(ctx context.Context, format string, v ...interface{}) {
	l.WithLevel(LevelFatal).Output(fmt.Sprintf(format, v...))
}

func (l *Logger) Panic(ctx context.Context, v ...interface{}) {
	l.WithLevel(LevelPanic).Output(fmt.Sprint(v...))
}

func (l *Logger) Panicf(ctx context.Context, format string, v ...interface{}) {
	l.WithLevel(LevelPanic).Output(fmt.Sprintf(format, v...))
}
