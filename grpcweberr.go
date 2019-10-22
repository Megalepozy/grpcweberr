package grpcweberr

import (
	"fmt"
	"net/http"
	"strconv"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPCWebErr serve as a receiver in order to allow mocking
type GRPCWebErr struct{}

// New is a factory method for GRPCWebErr struct
func New() *GRPCWebErr {
	return &GRPCWebErr{}
}

// New (with a receiver pointer) create a new error with embedded data of status codes and message
//
// grpcStatusCode codes.Code - import "google.golang.org/grpc/codes"
// httpStatusCode int - http status code to return
// messageToUser string - intended message returned to client, get default msg if empty
func (*GRPCWebErr) New(grpcStatusCode codes.Code, httpStatusCode int, messageToUser string) error {
	if messageToUser == "" {
		messageToUser = "An unexpected error has occurred, please try again later"
	}

	st := status.New(grpcStatusCode, messageToUser)

	br := &errdetails.BadRequest{}
	addHTTPStatusCode(httpStatusCode, br)
	addMessageToUser(messageToUser, br)

	return attachDetails(st, br).Err()
}

// AddLogTracingID append logTracingID value to the error so it could be tracked as it flow through the services
func (*GRPCWebErr) AddLogTracingID(logTracingID string, err error) error {
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

// GetHTTPStatus is a getter for the http value which was supplied at New(...)
func (*GRPCWebErr) GetHTTPStatus(err error) int {
	code := getFieldViolationValue(err, "httpStatus")
	if code != "" {
		code, _ := strconv.Atoi(code)
		return code
	}

	return 500
}

// GetMessageToUser is a getter for the messageToUser which was supplied at New(...)
func (*GRPCWebErr) GetMessageToUser(err error) string {
	return getFieldViolationValue(err, "messageToUser")
}

// GetLogTracingID is a getter for logTracingId
func (*GRPCWebErr) GetLogTracingID(err error) string {
	return getFieldViolationValue(err, "logTracingID")
}

func addHTTPStatusCode(httpStatusCode int, br *errdetails.BadRequest) {
	if http.StatusText(httpStatusCode) != "" {
		appendFieldViolation(br, "httpStatus", strconv.Itoa(httpStatusCode))
	}
}

func addMessageToUser(messageToUser string, br *errdetails.BadRequest) {
	appendFieldViolation(br, "messageToUser", messageToUser)
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
