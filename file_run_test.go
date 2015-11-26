package main

import (
	. "github.com/smartystreets/goconvey/convey"
	run "github.com/thriqon/involucro/steps/run"
	"testing"
)

func TestRunTaskDefinition(t *testing.T) {
	Convey("Given an empty runtime environment", t, func() {
		inv := InstantiateRuntimeEnv()
		env := inv.duk

		runCode := func(s string) func() {
			return func() {
				env.EvalString(s)
			}
		}

		Convey("Defining a Task with using/run succeeds", func() {
			So(runCode(`inv.task('test').using('blah').run('test')`), ShouldNotPanic)
		})

		Convey("Defining a Task with a number as ID panics", func() {
			So(runCode(`inv.task(5)`), ShouldPanic)
		})

		Convey("Calling using without a parameter panics", func() {
			So(runCode(`inv.task('test').using()`), ShouldPanic)
		})

		Convey("Calling run without a parameter panics", func() {
			So(runCode(`inv.task('test').using('blah').run()`), ShouldPanic)
		})

		Convey("When defining a task with one using/run", func() {
			env.EvalString(`inv.task('test').using('blah').run('test', '123')`)

			Convey("Then the task map has not an entry for another task", func() {
				_, ok := inv.tasks["another_task"]
				So(ok, ShouldBeFalse)
			})

			Convey("Then the task map has an entry for that task", func() {
				_, ok := inv.tasks["test"]
				So(ok, ShouldBeTrue)
			})

			Convey("Then the task map entry is an ExecuteImage struct", func() {
				So(inv.tasks["test"][0], ShouldHaveSameTypeAs, run.ExecuteImage{})
			})

			Convey("Then the task map entry has the given command set", func() {
				ei := inv.tasks["test"][0].(run.ExecuteImage)
				So(ei.Config.Cmd, ShouldResemble, []string{"test", "123"})
			})
		})
		Convey("Passing arguments to run works with different lengths", func() {
			env.EvalString(`inv.task('test1').using('blah').run('test')`)
			env.EvalString(`inv.task('test2').using('blah').run('test', 'asd')`)
			env.EvalString(`inv.task('test3').using('blah').run('test', 'asd', '123')`)
			env.EvalString(`inv.task('test4').using('blah').run('test', 'asd', '123', 'dsa')`)

			So(inv.tasks["test1"][0].(run.ExecuteImage).Config.Cmd, ShouldResemble, []string{"test"})
			So(inv.tasks["test2"][0].(run.ExecuteImage).Config.Cmd, ShouldResemble, []string{"test", "asd"})
			So(inv.tasks["test3"][0].(run.ExecuteImage).Config.Cmd, ShouldResemble, []string{"test", "asd", "123"})
			So(inv.tasks["test4"][0].(run.ExecuteImage).Config.Cmd, ShouldResemble, []string{"test", "asd", "123", "dsa"})
		})

		Convey("Executing multiple run tasks after each other works", func() {
			env.EvalString(`inv.task('test').using('blah').run('test').run('test2').using('asd').run('2')`)
			So(inv.tasks["test"][0].(run.ExecuteImage).Config.Image, ShouldEqual, "blah")
			So(inv.tasks["test"][0].(run.ExecuteImage).Config.Cmd, ShouldResemble, []string{"test"})

			So(inv.tasks["test"][1].(run.ExecuteImage).Config.Image, ShouldEqual, "blah")
			So(inv.tasks["test"][1].(run.ExecuteImage).Config.Cmd, ShouldResemble, []string{"test2"})

			So(inv.tasks["test"][2].(run.ExecuteImage).Config.Image, ShouldEqual, "asd")
			So(inv.tasks["test"][2].(run.ExecuteImage).Config.Cmd, ShouldResemble, []string{"2"})
		})
	})
}