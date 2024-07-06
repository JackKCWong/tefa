package main

import (
	"bytes"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTefa(t *testing.T) {
	Convey("It can generate with sequence number", t, func() {
		tefa, err := newTefa("", "No. {{ .Seq }}\n")
		So(err, ShouldBeNil)

		out := &bytes.Buffer{}
		tefa.Execute(out, 2)

		So(out.String(), ShouldEqual, "No. 1\nNo. 2\n")
	})

	Convey("It can use fake data", t, func() {
		tefa, err := newTefa("", `{{ .Name }}`)
		So(err, ShouldBeNil)

		out := &bytes.Buffer{}
		tefa.Execute(out, 1)

		So(out.String(), ShouldNotEqual, "")
	})

	Convey("It can use sprig functions", t, func() {
		tefa, err := newTefa("", `{{ "!" | repeat 3 }}`)
		So(err, ShouldBeNil)

		out := &bytes.Buffer{}
		tefa.Execute(out, 1)

		So(out.String(), ShouldEqual, "!!!")
	})
}

func TestPreTemplate(t *testing.T) {
	tefa, err := newTefa(`{{$i := 0}}`, `{{ $i = add1 $i }}{{ $i | printf "%d" | repeat 3 }}` + "\n")
	if err != nil {
		t.Fatal(err)
		t.FailNow()
	}

	out := &bytes.Buffer{}
	tefa.Execute(out, 3)

	if out.String() != "111\n222\n333\n" {
		t.Fatal(out.String())
		t.FailNow()
	}
}
