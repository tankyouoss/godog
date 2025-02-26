package models_test

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/cucumber/messages-go/v10"

	"github.com/tankyouoss/godog"
	"github.com/tankyouoss/godog/formatters"
	"github.com/tankyouoss/godog/internal/models"
)

func TestShouldSupportEmptyHandlerReturn(t *testing.T) {
	fn := func(a int64, b int32, c int16, d int8) {}

	def := &models.StepDefinition{
		StepDefinition: formatters.StepDefinition{
			Handler: fn,
		},
		HandlerValue: reflect.ValueOf(fn),
	}

	def.Args = []interface{}{"1", "1", "1", "1"}
	if err := def.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	def.Args = []interface{}{"1", "1", "1", strings.Repeat("1", 9)}
	if err := def.Run(); err == nil {
		t.Fatalf("expected convertion fail for int8, but got none")
	}
}

func TestShouldSupportIntTypes(t *testing.T) {
	fn := func(a int64, b int32, c int16, d int8) error { return nil }

	def := &models.StepDefinition{
		StepDefinition: formatters.StepDefinition{
			Handler: fn,
		},
		HandlerValue: reflect.ValueOf(fn),
	}

	def.Args = []interface{}{"1", "1", "1", "1"}
	if err := def.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	def.Args = []interface{}{"1", "1", "1", strings.Repeat("1", 9)}
	if err := def.Run(); err == nil {
		t.Fatalf("expected convertion fail for int8, but got none")
	}
}

func TestShouldSupportFloatTypes(t *testing.T) {
	fn := func(a float64, b float32) error { return nil }

	def := &models.StepDefinition{
		StepDefinition: formatters.StepDefinition{
			Handler: fn,
		},
		HandlerValue: reflect.ValueOf(fn),
	}

	def.Args = []interface{}{"1.1", "1.09"}
	if err := def.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	def.Args = []interface{}{"1.08", strings.Repeat("1", 65) + ".67"}
	if err := def.Run(); err == nil {
		t.Fatalf("expected convertion fail for float32, but got none")
	}
}

func TestShouldNotSupportOtherPointerTypesThanGherkin(t *testing.T) {
	fn1 := func(a *int) error { return nil }
	fn2 := func(a *messages.PickleStepArgument_PickleDocString) error { return nil }
	fn3 := func(a *messages.PickleStepArgument_PickleTable) error { return nil }

	def1 := &models.StepDefinition{
		StepDefinition: formatters.StepDefinition{
			Handler: fn1,
		},
		HandlerValue: reflect.ValueOf(fn1),
		Args:         []interface{}{(*int)(nil)},
	}
	def2 := &models.StepDefinition{
		StepDefinition: formatters.StepDefinition{
			Handler: fn2,
		},
		HandlerValue: reflect.ValueOf(fn2),
		Args:         []interface{}{&messages.PickleStepArgument_PickleDocString{}},
	}
	def3 := &models.StepDefinition{
		StepDefinition: formatters.StepDefinition{
			Handler: fn3,
		},
		HandlerValue: reflect.ValueOf(fn3),
		Args:         []interface{}{(*messages.PickleStepArgument_PickleTable)(nil)},
	}

	if err := def1.Run(); err == nil {
		t.Fatalf("expected conversion error, but got none")
	}

	if err := def2.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := def3.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestShouldSupportOnlyByteSlice(t *testing.T) {
	fn1 := func(a []byte) error { return nil }
	fn2 := func(a []string) error { return nil }

	def1 := &models.StepDefinition{
		StepDefinition: formatters.StepDefinition{
			Handler: fn1,
		},
		HandlerValue: reflect.ValueOf(fn1),
		Args:         []interface{}{"str"},
	}
	def2 := &models.StepDefinition{
		StepDefinition: formatters.StepDefinition{
			Handler: fn2,
		},
		HandlerValue: reflect.ValueOf(fn2),
		Args:         []interface{}{[]string{}},
	}

	if err := def1.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := def2.Run(); err == nil {
		t.Fatalf("expected conversion error, but got none")
	}
}

func TestUnexpectedArguments(t *testing.T) {
	fn := func(a, b int) error { return nil }
	def := &models.StepDefinition{
		StepDefinition: formatters.StepDefinition{
			Handler: fn,
		},
		HandlerValue: reflect.ValueOf(fn),
	}

	def.Args = []interface{}{"1"}

	res := def.Run()
	if res == nil {
		t.Fatalf("expected an error due to wrong number of arguments, but got none")
	}

	err, ok := res.(error)
	if !ok {
		t.Fatalf("expected an error due to wrong number of arguments, but got %T instead", res)
	}

	if !errors.Is(err, models.ErrUnmatchedStepArgumentNumber) {
		t.Fatalf("expected an error due to wrong number of arguments, but got %v instead", err)
	}
}

func TestStepDefinition_Run_StepShouldBeString(t *testing.T) {
	test := func(t *testing.T, fn interface{}) {
		def := &models.StepDefinition{
			StepDefinition: formatters.StepDefinition{
				Handler: fn,
			},
			HandlerValue: reflect.ValueOf(fn),
		}

		def.Args = []interface{}{12}

		res := def.Run()
		if res == nil {
			t.Fatalf("expected a string convertion error, but got none")
		}

		err, ok := res.(error)
		if !ok {
			t.Fatalf("expected a string convertion error, but got %T instead", res)
		}

		if !errors.Is(err, models.ErrCannotConvert) {
			t.Fatalf("expected a string convertion error, but got '%v' instead", err)
		}
	}

	// Ensure step type error if step argument is not a string
	// for all supported types.
	test(t, func(a int) error { return nil })
	test(t, func(a int64) error { return nil })
	test(t, func(a int32) error { return nil })
	test(t, func(a int16) error { return nil })
	test(t, func(a int8) error { return nil })
	test(t, func(a string) error { return nil })
	test(t, func(a float64) error { return nil })
	test(t, func(a float32) error { return nil })
	test(t, func(a *godog.Table) error { return nil })
	test(t, func(a *godog.DocString) error { return nil })
	test(t, func(a []byte) error { return nil })

}

func TestStepDefinition_Run_InvalidHandlerParamConversion(t *testing.T) {
	test := func(t *testing.T, fn interface{}) {
		def := &models.StepDefinition{
			StepDefinition: formatters.StepDefinition{
				Handler: fn,
			},
			HandlerValue: reflect.ValueOf(fn),
		}

		def.Args = []interface{}{12}

		res := def.Run()
		if res == nil {
			t.Fatalf("expected an unsupported argument type error, but got none")
		}

		err, ok := res.(error)
		if !ok {
			t.Fatalf("expected an unsupported argument type error, but got %T instead", res)
		}

		if !errors.Is(err, models.ErrUnsupportedArgumentType) {
			t.Fatalf("expected an unsupported argument type error, but got '%v' instead", err)
		}
	}

	// Lists some unsupported argument types for step handler.

	// Pointers should work only for godog.Table/godog.DocString
	test(t, func(a *int) error { return nil })
	test(t, func(a *int64) error { return nil })
	test(t, func(a *int32) error { return nil })
	test(t, func(a *int16) error { return nil })
	test(t, func(a *int8) error { return nil })
	test(t, func(a *string) error { return nil })
	test(t, func(a *float64) error { return nil })
	test(t, func(a *float32) error { return nil })

	// I cannot pass structures
	test(t, func(a godog.Table) error { return nil })
	test(t, func(a godog.DocString) error { return nil })
	test(t, func(a testStruct) error { return nil })

	// I cannot use maps
	test(t, func(a map[string]interface{}) error { return nil })
	test(t, func(a map[string]int) error { return nil })

	// Slice works only for byte
	test(t, func(a []int) error { return nil })
	test(t, func(a []string) error { return nil })
	test(t, func(a []bool) error { return nil })

	// I cannot use bool
	test(t, func(a bool) error { return nil })

}

func TestStepDefinition_Run_StringConversionToFunctionType(t *testing.T) {
	test := func(t *testing.T, fn interface{}, args []interface{}) {
		def := &models.StepDefinition{
			StepDefinition: formatters.StepDefinition{
				Handler: fn,
			},
			HandlerValue: reflect.ValueOf(fn),
			Args:         args,
		}

		res := def.Run()
		if res == nil {
			t.Fatalf("expected a cannot convert argument type error, but got none")
		}

		err, ok := res.(error)
		if !ok {
			t.Fatalf("expected a cannot convert argument type error, but got %T instead", res)
		}

		if !errors.Is(err, models.ErrCannotConvert) {
			t.Fatalf("expected a cannot convert argument type error, but got '%v' instead", err)
		}
	}

	// Lists some unsupported argument types for step handler.

	// Cannot convert invalid int
	test(t, func(a int) error { return nil }, []interface{}{"a"})
	test(t, func(a int64) error { return nil }, []interface{}{"a"})
	test(t, func(a int32) error { return nil }, []interface{}{"a"})
	test(t, func(a int16) error { return nil }, []interface{}{"a"})
	test(t, func(a int8) error { return nil }, []interface{}{"a"})

	// Cannot convert invalid float
	test(t, func(a float32) error { return nil }, []interface{}{"a"})
	test(t, func(a float64) error { return nil }, []interface{}{"a"})

	// Cannot convert to DataArg
	test(t, func(a *godog.Table) error { return nil }, []interface{}{"194"})

	// Cannot convert to DocString ?
	test(t, func(a *godog.DocString) error { return nil }, []interface{}{"194"})

}

// @TODO maybe we should support duration
// fn2 := func(err time.Duration) error { return nil }
// def = &models.StepDefinition{Handler: fn2, HandlerValue: reflect.ValueOf(fn2)}

// def.Args = []interface{}{"1"}
// if err := def.Run(); err == nil {
// 	t.Fatalf("expected an error due to wrong argument type, but got none")
// }

type testStruct struct {
	a string
}
