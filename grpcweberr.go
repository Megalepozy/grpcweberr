package grpcweberr

import (
	"fmt"
	"net/http"
	"strconv"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// New create a new error where the returned error is a simple error object (avoid later type inference)
//
// id string - just a string to easily identify which error occured
// code codes.Code - grpc status code, import "google.golang.org/grpc/codes"
// http int - the intended http status code to return
// msg string - intended message returned to client, get default msg if empty
func New(id string, code codes.Code, http int, msg string) error {
	st := status.New(code, id)

	br := &errdetails.BadRequest{}
	addErrorID(id, br)
	addHTTPStatusCode(http, br)
	addUserErrorMessage(msg, br)

	return attachDetails(st, br).Err()
}

// AddLogTracingID append logTracingID value to the error
func AddLogTracingID(logTracingID string, err error) error {
	st := status.Convert(err)
	for _, detail := range st.Details() {
		switch t := detail.(type) {
		case *errdetails.BadRequest:
			appendFieldViolation(t, "logTracingID", logTracingID)
			return attachDetails(st, t).Err()
		}
	}

	return err
}

// GetErrorID is a getter for the id value which was supplied at New(...)
func GetErrorID(err error) string {
	return getFieldViolationValue(err, "errorID")
}

// GetHTTPStatus is a getter for the http value which was supplied at New(...)
func GetHTTPStatus(err error) int {
	code := getFieldViolationValue(err, "httpStatus")
	if code != "" {
		code, _ := strconv.Atoi(code)
		return code
	}

	return 500
}

// GetUserErrorMessage is a getter for the msg value which was supplied at New(...)
func GetUserErrorMessage(err error) string {
	return getFieldViolationValue(err, "userErrorMessage")
}

// GetLogTracingID is a getter for logTracingId
func GetLogTracingID(err error) string {
	return getFieldViolationValue(err, "logTracingID")
}

func addErrorID(errorID string, br *errdetails.BadRequest) {
	appendFieldViolation(br, "errorID", errorID)
}

func addHTTPStatusCode(httpCode int, br *errdetails.BadRequest) {
	if http.StatusText(httpCode) != "" {
		appendFieldViolation(br, "httpStatus", strconv.Itoa(httpCode))
	}
}

func addUserErrorMessage(userErrorMessage string, br *errdetails.BadRequest) {
	if userErrorMessage == "" {
		userErrorMessage = "An unexpected error has occurred, please try again later"
	}

	appendFieldViolation(br, "userErrorMessage", userErrorMessage)
}

func appendFieldViolation(br *errdetails.BadRequest, name string, value string) *errdetails.BadRequest {
	fv := &errdetails.BadRequest_FieldViolation{
		Field:       name,
		Description: value,
	}

	br.FieldViolations = append(br.FieldViolations, fv)
	return br
}

func attachDetails(st *status.Status, br *errdetails.BadRequest) *status.Status {
	st, err := st.WithDetails(br)
	if err != nil {
		panic(fmt.Sprintf("Unexpected error attaching metadata: %v", err))
	}

	return st
}

func getFieldViolationValue(err error, fieldname string) string {
	st := status.Convert(err)

	for _, detail := range st.Details() {
		switch t := detail.(type) {
		case *errdetails.BadRequest:
			for _, violation := range t.GetFieldViolations() {
				if violation.GetField() == fieldname {
					return violation.GetDescription()
				}
			}
		}
	}

	return ""
}
