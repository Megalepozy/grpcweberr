package grpcweberr

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func New(errorID string, grpcCode codes.Code, httpCode int, userErrorMessage string) error {
	st := status.New(grpcCode, errorID)

	br := &errdetails.BadRequest{}
	addErrorID(errorID, br)
	addHTTPStatusCode(httpCode, br)
	addUserErrorMessage(userErrorMessage, br)

	return attachDetails(st, br).Err()
}

func AddLogTracingID(err error) error {
	logTracingID := uuid.New()

	st := status.Convert(err)
	for _, detail := range st.Details() {
		switch t := detail.(type) {
		case *errdetails.BadRequest:
			appendFieldViolation(t, "logTracingID", logTracingID.String())
			return attachDetails(st, t).Err()
		}
	}

	return err
}

func GetErrorID(err error) string {
	return getFieldViolationValue(err, "errorID")
}

func GetHTTPStatus(err error) int {
	code := getFieldViolationValue(err, "httpStatus")
	if code != "" {
		code, _ := strconv.Atoi(code)
		return code
	}

	return 500
}

func GetUserErrorMessage(err error) string {
	return getFieldViolationValue(err, "userErrorMessage")
}

func GetLogTracingID(err error) string {
	return getFieldViolationValue(err, "logTracingID")
}

func IsGotTracingID(err error) bool {
	logTracingID := getFieldViolationValue(err, "logTracingID")
	if logTracingID != "" {
		return true
	}

	return false
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
