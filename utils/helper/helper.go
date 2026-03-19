package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rshby/go-event-ticketing/utils/singleton"
)

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
	~float32 | ~float64
}

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

func GenerateSonyFlakeID() (int64, error) {
	return singleton.SF.NextID()
}

func ExpectedNumber[T Number](input any) T {
	var zero T

	if input == nil {
		return zero
	}

	switch v := input.(type) {
	case int:
		return T(v)
	case int8:
		return T(v)
	case int16:
		return T(v)
	case int32:
		return T(v)
	case int64:
		return T(v)
	case uint:
		return T(v)
	case uint8:
		return T(v)
	case uint16:
		return T(v)
	case uint32:
		return T(v)
	case uint64:
		return T(v)
	case float32:
		return T(v)
	case float64:
		return T(v)
	case bool:
		if v {
			return T(1)
		}

		return zero
	case string:
		parsed, err := strconv.ParseFloat(v, 64)
		if err == nil {
			return T(parsed)
		}
		return zero
	default:
		return zero
	}
}

func ExpectedString(input any) string {
	if input == nil {
		return ""
	}

	switch v := input.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	case []byte:
		return string(v)
	case error:
		return v.Error()
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}
