package logrus_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	this "github.com/visenze/logrus"
)

const (
	testProject   = "AS"
	testComponent = "visenze-logrus-formatter"
)

func TestNewFormatter(t *testing.T) {
	Convey("When a formatter is initialized with no component name", t, func() {
		_, err := this.NewOSSFormatter(testProject, "")
		Convey("There must be an error", func() {
			So(err, ShouldNotBeNil)
		})
	})
	Convey("When a formatter is initialized with no project name", t, func() {
		_, err := this.NewOSSFormatter("", testComponent)
		Convey("There must be an error", func() {
			So(err, ShouldNotBeNil)
		})
	})
	Convey("When a formatter is initialized with both project and component names", t, func() {
		f, err := this.NewOSSFormatter("as", testComponent)
		Convey("There must be no error", func() {
			So(err, ShouldBeNil)
			So(f, ShouldNotBeNil)
		})
	})
}

func TestFormatter_Format(t *testing.T) {
	Convey("Given a formatter", t, func() {
		f, err := this.NewOSSFormatter(testProject, testComponent)
		So(err, ShouldBeNil)
		Convey("When an message is formatted", func() {
			now := time.Now()
			res, err := f.Format(&logrus.Entry{
				Time:    now,
				Level:   logrus.InfoLevel,
				Message: "hello",
				Data: logrus.Fields{
					"total_time":   56,
					"tag_group":    "qrcode",
					"rpc_error":    nil,
					"remote":       "remote-qrcode-2-shared.recognition-worker:9010",
					"grpc.time_ms": 200,
				}},
			)
			Convey("There should be no error", func() {
				So(err, ShouldBeNil)
			})
			Convey("The message must begin with {project_name:component_name}", func() {
				So(string(res), ShouldStartWith, fmt.Sprintf("%s:%s", testProject, testComponent))
			})
			Convey("The message must end with a new line", func() {
				So(string(res), ShouldEndWith, "\n")
			})
			Convey("The message must contain the timestamp", func() {
				So(string(res), ShouldContainSubstring, now.UTC().Format(this.TimestampFormat))
			})
			Convey("The message must contain the log_level", func() {
				So(string(res), ShouldContainSubstring, `log_level=INFO`)
			})
			Convey("The message must contain the message, triple-single-quoted", func() {
				So(string(res), ShouldContainSubstring, `msg='''hello'''`)
			})
			Convey("The message must contain the total_time", func() {
				So(string(res), ShouldContainSubstring, "total_time=56")
			})
			Convey("The message must contain the tag_group", func() {
				So(string(res), ShouldContainSubstring, "tag_group=qrcode")
			})
			Convey("The message must contain the rpc_error", func() {
				So(string(res), ShouldContainSubstring, "rpc_error=<nil>")
			})
			Convey("The message must contain the remote quoted", func() {
				So(string(res), ShouldContainSubstring, "remote=\"remote-qrcode-2-shared.recognition-worker:9010\"")
			})
			Convey("The message must contain the grpc_time_ms (replaced the . with a _)", func() {
				So(string(res), ShouldContainSubstring, "grpc_time_ms=200")
			})
		})
	})
}
