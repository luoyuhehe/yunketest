package thinkgo

import (
	"context"
	"gitee.com/sahara-go/thinkgo/log"
	"runtime/debug"
)

// GoroutineAttachPanicHandle 协程panic处理
func GoroutineAttachPanicHandle(f func()) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("goroutine panic: %v, stacktrace:%v", err, string(debug.Stack()))
		}
	}()
	f()
}

func HandleCancelAndPanic(f func(ctx context.Context)) func(ctx context.Context) {
	return func(ctx context.Context) {
		ctx1 := CopyCtx(ctx)
		defer func() {
			if err := recover(); err != nil {
				log.Errorf("goroutine panic: %v, stacktrace:%v", err, string(debug.Stack()))
			}
		}()
		f(ctx1)
	}
}

func CopyCtx(ctx context.Context) context.Context {
	return ctx
}
