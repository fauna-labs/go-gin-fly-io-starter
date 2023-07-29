// Copyright Fauna, Inc.
// SPDX-License-Identifier: MIT-0

package utils

import (
	"net/http"
	"os"

	"github.com/fauna/fauna-go"
)

func GenerateResponse(res *fauna.QuerySuccess) (r map[string]any) {
	var d any
	if _, ok := res.Data.(*fauna.Page); ok {
		d = res.Data.(*fauna.Page).Data
	} else {
		d = res.Data
	}
	s := res.Stats
	r = map[string]any{
		"data":       d,
		"faunaStats": s,
		"flyStats": map[string]string{
			"usingFlyRegion": os.Getenv("FLY_REGION"),
		},
	}
	return
}

func GetErrorResponseStatusCode(err error) (s int) {
	s = http.StatusBadRequest

	if _, ok := err.(*fauna.ErrAbort); ok {
		return
	} else if _, ok := err.(*fauna.ErrAuthentication); ok {
		s = err.(*fauna.ErrAuthentication).StatusCode
	} else if _, ok := err.(*fauna.ErrAuthorization); ok {
		s = err.(*fauna.ErrAuthorization).StatusCode
	}
	return
}
