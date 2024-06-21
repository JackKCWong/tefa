package main

import (
	"bytes"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTefa(t *testing.T) {
	Convey("It can generate with sequence number", t, func() {
		tefa, err := newTefa(`No. {{ .Seq }}`)
		So(err, ShouldBeNil)

		out := &bytes.Buffer{}
		tefa.Execute(out)

		So(out.String(), ShouldEqual, "No. 1")

		out.Reset()
		tefa.Execute(out)

		So(out.String(), ShouldEqual, "No. 2")
	})

	Convey("It can use fake data", t, func() {
		tefa, err := newTefa(`{{ .Name }}`)
		So(err, ShouldBeNil)

		out := &bytes.Buffer{}
		tefa.Execute(out)

		So(out.String(), ShouldNotEqual, "")
	})

	Convey("It can use sprig functions", t, func() {
		tefa, err := newTefa(`{{ "!" | repeat 3 }}`)
		So(err, ShouldBeNil)

		out := &bytes.Buffer{}
		tefa.Execute(out)

		So(out.String(), ShouldEqual, "!!!")
	})
}
