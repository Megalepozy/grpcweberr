package grpcweberr

import (
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/grpc/codes"
)

func TestNew(t *testing.T) {
	Convey("Setup of valid values for error creation", t, func() {
		errorID := "short unique err description"
		grpcStatus := codes.Aborted
		httpStatus := http.StatusBadRequest
		userErrorMessage := "Oh no, something really bad happened"

		Convey("Passing of non existing http status code result in http status code 500", func() {
			expectedHTTPStatus := http.StatusInternalServerError
			httpStatus = 823764

			err := New(errorID, grpcStatus, httpStatus, userErrorMessage)

			So(err, ShouldNotBeNil)
			So(GetHTTPStatus(err), ShouldEqual, expectedHTTPStatus)
		})

		Convey("Passing of empty userErrorMessage result in generic userErrorMessage", func() {
			expectedUserErrorMessage := "An unexpected error has occurred, please try again later"
			userErrorMessage = ""

			err := New(errorID, grpcStatus, httpStatus, userErrorMessage)

			So(err, ShouldNotBeNil)
			So(GetUserErrorMessage(err), ShouldEqual, expectedUserErrorMessage)
		})

		Convey("Creation of error with valid values", func() {
			err := New(errorID, grpcStatus, httpStatus, userErrorMessage)

			So(err, ShouldNotBeNil)

			Convey("Proper err.Error()", func() {
				So(err.Error(), ShouldEqual, "rpc error: code = Aborted desc = "+errorID)
			})
			Convey("Getting of http status", func() {
				So(GetHTTPStatus(err), ShouldEqual, httpStatus)
			})
			Convey("Getting of user error userErrorMessage", func() {
				So(GetUserErrorMessage(err), ShouldEqual, userErrorMessage)
			})

			Convey("No log tracing ID", func() {
				So(IsGotTracingID(err), ShouldBeFalse)

				Convey("Returned log tracing ID is an empty string", func() {
					So(GetLogTracingID(err), ShouldEqual, "")
				})
			})

			Convey("Add log tracing ID", func() {
				err = AddLogTracingID(err)

				So(IsGotTracingID(err), ShouldBeTrue)

				Convey("Returned log tracing ID", func() {
					So(GetLogTracingID(err), ShouldNotBeEmpty)
				})
			})
		})
	})
}
