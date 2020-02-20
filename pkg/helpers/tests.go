package helpers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// UUIDFromString ...
func UUIDFromString(t *testing.T, uStr string) uuid.UUID {
	id, err := uuid.Parse(uStr)
	if err != nil {
		if t == nil {
			log.Fatalf("Unable to get UUID from string: %s", err.Error())
		} else {
			t.Fatalf("Unable to get UUID from string: %s", err.Error())
		}
	}

	return id
}

// UUIDFromStringPointer ...
func UUIDFromStringPointer(t *testing.T, uStr string) *uuid.UUID {
	id := UUIDFromString(t, uStr)

	return &id
}

// ObjectToByte ...
func ObjectToByte(t *testing.T, obj interface{}) *bytes.Reader {
	b, err := json.Marshal(obj)
	if err != nil {
		t.Fatalf("Failed to marshal body: %v", err.Error())
	}

	return bytes.NewReader(b)
}

// RequestTest ...
func RequestTest(method, path string, e *echo.Echo) (int, string) {
	req := httptest.NewRequest(method, path, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	return rec.Code, rec.Body.String()
}

// RequestWithBody ...
func RequestWithBody(method, path string, body io.Reader, e *echo.Echo) (*httptest.ResponseRecorder, echo.Context) {
	req := httptest.NewRequest(method, path, body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	return rec, e.NewContext(req, rec)
}

// RequestObjectWithBody ...
func RequestObjectWithBody(t *testing.T, method, path string, data interface{}, e *echo.Echo) (*httptest.ResponseRecorder, echo.Context) {
	req := httptest.NewRequest(method, path, ObjectToByte(t, data))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	return rec, e.NewContext(req, rec)
}
