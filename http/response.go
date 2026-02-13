package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/adityayuga/go-sdk/debug"
)

type (
	API struct {
		startTime time.Time
		resp      http.ResponseWriter
		debugID   string
	}
)

func NewResponder(w http.ResponseWriter, r *http.Request) (API, context.Context) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	return NewResponderWithContext(ctx, w, r)
}

func NewResponderWithContext(ctx context.Context, w http.ResponseWriter, r *http.Request) (API, context.Context) {
	// get debug id
	debugID := GetDebugIDFromRequest(r)
	if debugID == "" {
		debugID, ctx = debug.GetDebugIDFromContext(ctx)
	}

	return API{
		startTime: time.Now(),
		resp:      w,
		debugID:   debugID,
	}, ctx
}

func (m API) SuccessResponse(data interface{}) error {
	return m.RawResponse(http.StatusOK, "", data)
}

func (m API) SuccessResponseWithMessage(data interface{}, message string) error {
	return m.RawResponse(http.StatusOK, message, data)
}

func (m API) ErrorResponse(err error) error {
	return m.ErrorResponseWithCustomData(err, nil)
}

func (m API) ErrorResponseWithCustomData(err error, data interface{}) error {
	return m.RawResponse(http.StatusInternalServerError, err.Error(), data)
}

func (m API) ErrorResponseWithReqData(err error, data interface{}) error {
	var reqData map[string]interface{}
	if data != nil {
		reqData = map[string]interface{}{
			"request": data,
		}
	}
	return m.RawResponse(http.StatusInternalServerError, err.Error(), reqData)
}

func (m API) BadRequestResponse(err error) error {
	return m.RawResponse(http.StatusBadRequest, err.Error(), nil)
}

func (m API) UnathorizedResponse(err error) error {
	return m.UnathorizedResponseWithReqData(err, nil)
}

func (m API) UnathorizedResponseWithReqData(err error, data interface{}) error {
	var reqData map[string]interface{}
	if data != nil {
		reqData = map[string]interface{}{
			"request": data,
		}
	}
	return m.RawResponse(http.StatusUnauthorized, err.Error(), reqData)
}

func (m API) GetDebugID() string {
	return m.debugID
}

func (m API) RawResponse(statusCode int, message string, data interface{}) error {
	standardData := ConstructStdResponseData(m.debugID, m.startTime, data, message)
	marshaled, err := json.Marshal(standardData)
	if err != nil {
		return err
	}

	m.resp.Header().Set("Content-Type", "application/json")
	m.resp.WriteHeader(statusCode)
	m.resp.Write(marshaled)
	return nil
}

func ConstructErrorResponseData(debugID string, startProcessingTime time.Time, message string) map[string]interface{} {
	return ConstructStdResponseData(debugID, startProcessingTime, nil, message)
}

func ConstructStdResponseData(debugID string, startProcessingTime time.Time, data interface{}, message string) map[string]interface{} {
	standardData := map[string]interface{}{
		"debug_id": debugID,
	}

	if debugID != "" {
		standardData["processing_time"] = fmt.Sprintf("%d ms", time.Since(startProcessingTime).Milliseconds())
	}

	if message != "" {
		standardData["message"] = message
	}

	if data != nil {
		standardData["data"] = data
	}

	return standardData
}

// GetDebugIDFromRequest get from request header
func GetDebugIDFromRequest(r *http.Request) string {
	return r.Header.Get(HeaderDebugID)
}
