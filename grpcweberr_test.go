package grpcweberr

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/grpc/codes"
)

func TestNew(t *testing.T) {
	Convey("Setup of valid values for error creation", t, func() {
		id := "short unique err description"
		code := codes.Aborted
		httpStatus := http.StatusBadRequest
		msg := "Oh no, something really bad happened"

		Convey("Passing of non existing http status code result in http status code 500", func() {
			expectedHTTPStatus := http.StatusInternalServerError
			httpStatus = 823764

			err := New(id, code, httpStatus, msg)

			So(err, ShouldNotBeNil)
			So(GetHTTPStatus(err), ShouldEqual, expectedHTTPStatus)
		})

		Convey("Passing of empty msg result in generic msg", func() {
			expectedUserErrorMessage := "An unexpected error has occurred, please try again later"
			msg = ""

			err := New(id, code, httpStatus, msg)

			So(err, ShouldNotBeNil)
			So(GetUserErrorMessage(err), ShouldEqual, expectedUserErrorMessage)
		})

		Convey("Creation of error with valid values", func() {
			err := New(id, code, httpStatus, msg)

			So(err, ShouldNotBeNil)

			Convey("Proper err.Error()", func() {
				So(err.Error(), ShouldEqual, "rpc error: code = Aborted desc = "+id)
			})
			Convey("Getting of id", func() {
				So(GetErrorID(err), ShouldEqual, id)
			})
			Convey("Getting of http status", func() {
				So(GetHTTPStatus(err), ShouldEqual, httpStatus)
			})
			Convey("Getting of user error msg", func() {
				So(GetUserErrorMessage(err), ShouldEqual, msg)
			})

			Convey("No log tracing ID", func() {
				Convey("Returned log tracing ID is an empty string", func() {
					So(GetLogTracingID(err), ShouldEqual, "")
				})
			})

			Convey("Add log tracing ID", func() {
				expectedLogTracingID := uuid.New().String()

				err = AddLogTracingID(expectedLogTracingID, err)

				Convey("Returned log tracing ID", func() {
					So(GetLogTracingID(err), ShouldEqual, expectedLogTracingID)
				})
			})
		})
	})
}
