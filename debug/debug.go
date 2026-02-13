package debug

import (
	"context"
	"crypto/rand"
	"io"
	"time"
)

const (
	ctxDebugID string = "x-debug-id"
)

func GetDebugIDFromContext(ctx context.Context) (string, context.Context) {
	debugIDIface := ctx.Value(ctxDebugID)
	if debugIDIface == nil {
		debugID := GenerateDebugID()
		return debugID, SetDebugIDOnContext(ctx, debugID)
	}

	return debugIDIface.(string), ctx
}

func SetDebugIDOnContext(ctx context.Context, debugID string) context.Context {
	return context.WithValue(ctx, ctxDebugID, debugID)
}

func GenerateDebugID() string {
	n := 8
	table := [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	for {
		b := make([]byte, n)
		l, err := io.ReadAtLeast(rand.Reader, b, n)
		if err != nil || l != n {
			continue
		}
		for i := 0; i < len(b); i++ {
			b[i] = table[int(b[i])%len(table)]
		}
		return string(b) + time.Now().Format("20060102150405")
	}
}
