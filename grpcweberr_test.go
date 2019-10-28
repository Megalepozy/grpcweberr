package grpcweberr

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/grpc/codes"
)

func TestNew(t *testing.T) {
	Convey("Calling New return GRPCWebErr (factory method)", t, func() {
		gwe := New()

		So(gwe, ShouldNotBeNil)
		So(gwe, ShouldHaveSameTypeAs, &GRPCWebErr{})

		Convey("Setup of valid values for error initialization", func() {
			code := codes.Aborted
			httpStatus := http.StatusBadRequest
			msg := "Oh no, something really bad happened"

			Convey("Passing of non existing http status code result in http status code 500", func() {
				expectedHTTPStatus := http.StatusInternalServerError
				httpStatus = 823764

				err := gwe.New(code, httpStatus, msg)

				So(err, ShouldNotBeNil)
				So(gwe.GetHTTPStatus(err), ShouldEqual, expectedHTTPStatus)
			})

			Convey("Passing of empty message to user result in generic message", func() {
				expectedMessageToUser := "An unexpected error has occurred, please try again later"
				msg = ""

				err := gwe.New(code, httpStatus, msg)

				So(err, ShouldNotBeNil)
				So(gwe.GetMessageToUser(err), ShouldEqual, expectedMessageToUser)
			})

			Convey("Creation of error with valid values", func() {
				err := gwe.New(code, httpStatus, msg)

				So(err, ShouldNotBeNil)

				Convey("Proper err.Error()", func() {
					So(err.Error(), ShouldEqual, "rpc error: code = Aborted desc = "+msg)
				})
				Convey("Getting of http status", func() {
					So(gwe.GetHTTPStatus(err), ShouldEqual, httpStatus)
				})
				Convey("Getting of message to user", func() {
					So(gwe.GetMessageToUser(err), ShouldEqual, msg)
				})

				Convey("No log tracing ID", func() {
					Convey("Returned log tracing ID is an empty string", func() {
						So(gwe.GetLogTracingID(err), ShouldEqual, "")
					})
				})

				Convey("Add log tracing ID", func() {
					expectedLogTracingID := uuid.New().String()

					err = gwe.AddLogTracingID(expectedLogTracingID, err)

					Convey("Returned log tracing ID", func() {
						So(gwe.GetLogTracingID(err), ShouldEqual, expectedLogTracingID)
					})
				})
			})
		})
	})
}
