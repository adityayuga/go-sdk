package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/adityayuga/go-sdk/debug"
	"github.com/adityayuga/go-sdk/errors"
)

type (
	Request struct {
		URL           string
		Path          string
		Headers       map[string]string
		URLParams     map[string]string
		BodyParams    interface{}
		Timeout       time.Duration
		SkipErrorCode []int
	}
)

func Post(ctx context.Context, req Request) (respBody []byte, respCode int, err error) {
	respBody, respCode, err = doRequest(ctx, http.MethodPost, req)
	return
}

func Get(ctx context.Context, req Request) (respBody []byte, respCode int, err error) {
	respBody, respCode, err = doRequest(ctx, http.MethodGet, req)
	return
}

func doRequest(ctx context.Context, reqMethod string, req Request) (respBody []byte, respCode int, err error) {
	// debug id
	debugID, ctx := debug.GetDebugIDFromContext(ctx)

	// setup default timeout
	if req.Timeout == 0 {
		req.Timeout = 3 * time.Second
	}

	// create http client
	httpClient := &http.Client{
		Timeout: req.Timeout,
	}

	// setup url & path
	url, errParse := url.Parse(req.URL)
	if errParse != nil {
		err = errors.New(fmt.Sprintf("[%s] %s", debugID, errParse.Error()))
		return
	}

	if req.Path != "" {
		url.Path = req.Path
	}

	// setup url params
	if len(req.URLParams) > 0 {
		q := url.Query()
		for k, v := range req.URLParams {
			q.Set(k, v)
		}
		url.RawQuery = q.Encode()
	}

	// setup body params
	marshaledReqBodyParams, errMarshal := json.Marshal(req.BodyParams)
	if errMarshal != nil {
		err = errors.New(fmt.Sprintf("[%s] %s", debugID, errMarshal.Error()))
		return
	}

	bodyParams := bytes.NewBufferString(string(marshaledReqBodyParams))

	// setup http request
	httpReq, errSetupReq := http.NewRequestWithContext(ctx, reqMethod, url.String(), bodyParams)
	if errSetupReq != nil {
		err = errors.New(fmt.Sprintf("[%s] %s", debugID, errSetupReq.Error()))
		return
	}

	// setup headers
	if reqMethod == http.MethodPost {
		httpReq.Header.Set("Content-Type", "application/json")
	}
	httpReq.Header.Set(HeaderDebugID, debugID)
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	// inject headers from context
	if ctx != nil {
		// debug id
		var debugIDFromCtx string
		debugIDFromCtx, ctx = debug.GetDebugIDFromContext(ctx)
		httpReq.Header.Set(HeaderDebugID, debugIDFromCtx)

		// x-apps-id
		var appsIDFromCtx string
		appsIDFromCtx, _ = ctx.Value("x-apps-id").(string)
		if appsIDFromCtx != "" {
			httpReq.Header.Set("x-apps-id", appsIDFromCtx)
		}
	}

	// do request
	resp, errDoReq := httpClient.Do(httpReq)
	if errDoReq != nil {
		err = errors.New(fmt.Sprintf("[%s] %s", debugID, errDoReq.Error()))
		return
	}

	if resp == nil {
		err = errors.New(fmt.Sprintf("[%s] response request to %s is nil", debugID, req.URL))
		return
	}

	defer resp.Body.Close()

	respCode = resp.StatusCode

	// check status code
	skipErrCodeMap := map[int]bool{
		http.StatusOK:                   true,
		http.StatusCreated:              true,
		http.StatusAccepted:             true,
		http.StatusNonAuthoritativeInfo: true,
		http.StatusNoContent:            true,
		http.StatusResetContent:         true,
		http.StatusPartialContent:       true,
		http.StatusMultiStatus:          true,
		http.StatusAlreadyReported:      true,
		http.StatusIMUsed:               true,
	}
	for _, code := range req.SkipErrorCode {
		skipErrCodeMap[code] = true
	}

	if !skipErrCodeMap[respCode] {
		err = errors.New(fmt.Sprintf("[%s] status code: %d", debugID, respCode))
		return
	}

	// get response body
	respBody, errParseResp := io.ReadAll(resp.Body)
	if errParseResp != nil {
		err = errors.New(fmt.Sprintf("[%s] status code: %d, %s", debugID, respCode, errParseResp.Error()))
		return
	}

	return
}
