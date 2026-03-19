package helper

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
)

func Dump(v any) string {
	jsonByte, err := json.Marshal(v)
	if err == nil {
		return string(jsonByte)
	}

	return ""
}

func DumpIncomingContext(ctx context.Context) string {
	if ctx == nil {
		return "null"
	}

	dumpData := make(map[string]any)

	if ginCtx, ok := ctx.(*gin.Context); ok {
		if len(ginCtx.Keys) == 0 {
			return "{}"
		}
		for k, v := range ginCtx.Keys {
			dumpData[k.(string)] = v
		}
	} else {
		dumpData["raw_context"] = fmt.Sprintf("%v", ctx)
	}

	// Convert map ke JSON string
	bytesData, err := json.Marshal(dumpData)
	if err != nil {
		return fmt.Sprintf(`{"error": "marshal failed: %v"}`, err)
	}

	return string(bytesData)
}
