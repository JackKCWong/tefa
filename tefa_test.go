package main

import (
	"bytes"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTefa(t *testing.T) {
	Convey("It can generate with sequence number", t, func() {
		tefa, err := newTefa(`No. {{ seq }}`)
		So(err, ShouldBeNil)

		out := &bytes.Buffer{}
		tefa.Execute(out)

		So(out.String(), ShouldEqual, "No. 1")

		out.Reset()
		tefa.Execute(out)

		So(out.String(), ShouldEqual, "No. 2")
	})
}
