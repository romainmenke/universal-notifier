package trace

import (
	"bytes"
	"runtime/debug"
)

var goroutineHeader = []byte("goroutine")

func captureStack(skip int) []byte {
	stack := debug.Stack()
	out := stack[:0]
	skip += 2

	if bytes.HasPrefix(stack, goroutineHeader) {
		if idx := bytes.IndexByte(stack, '\n'); idx > 0 {
			out = append(out, stack[:idx+1]...)
			stack = stack[idx+1:]
		}
	}

	// skip
	for i := 0; i < skip; i++ {
		if idx := bytes.IndexByte(stack, '\n'); idx > 0 {
			stack = stack[idx+1:]
		}
		if idx := bytes.IndexByte(stack, '\n'); idx > 0 {
			stack = stack[idx+1:]
		}
	}

	out = append(out, stack...)

	return out
}
