package env_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/m2tx/gofxtest/internal/env"
	"github.com/stretchr/testify/assert"
)

type testcase[T any] struct {
	name          string
	config        *T
	test          func(t *testing.T, tcase testcase[T])
	env           map[string]string
	expectedValue any
	expectedError error
}

func test[T any](t *testing.T, testcases []testcase[T]) {
	for i, tcase := range testcases {
		t.Run(fmt.Sprintf("CASE %d - %s", i, tcase.name), func(t *testing.T) {
			if len(tcase.env) > 0 {
				for k, v := range tcase.env {
					err := os.Setenv(k, v)
					if err != nil {
						assert.Fail(t, "error setenv")
					}
				}
				defer func() {
					for k := range tcase.env {
						err := os.Unsetenv(k)
						if err != nil {
							assert.Fail(t, "error unsetenv")
						}
					}
				}()
			}

			err := env.PopulateFromEnv(tcase.config)
			if tcase.expectedError != nil {
				assert.Error(t, err)
				assert.IsType(t, &env.EnvError{}, err)
				assert.Equal(t, tcase.expectedError.Error(), err.Error())
			}

			if tcase.test != nil {
				tcase.test(t, tcase)
			}
		})
	}
}

func TestEnv_New(t *testing.T) {
	type config struct {
		ValueString   string        `env:"VALUE_STRING" default:"DEFAULT_STRING"`
		ValueInt      int           `env:"VALUE_INT" default:"42"`
		ValueFloat    float64       `env:"VALUE_FLOAT" default:"3.14"`
		ValueBool     bool          `env:"VALUE_BOOL" default:"true"`
		ValueDuration time.Duration `env:"VALUE_DURATION" default:"1s"`
	}

	c, err := env.New[config]()
	assert.NoError(t, err)
	assert.Equal(t, "DEFAULT_STRING", c.ValueString)
	assert.Equal(t, 42, c.ValueInt)
	assert.Equal(t, 3.14, c.ValueFloat)
	assert.Equal(t, true, c.ValueBool)
	assert.Equal(t, time.Second, c.ValueDuration)
}

func TestEnv_PopulateFromEnv(t *testing.T) {
	t.Run("Populate with defaults and env vars", func(t *testing.T) {
		testPopulateFromEnvWithDefaultsAndEnvVars(t)
	})
	t.Run("Populate with required vars", func(t *testing.T) {
		testPopulateFromEnvWithRequiredVars(t)
	})
	t.Run("Populate with min and max", func(t *testing.T) {
		testPopulateFromEnvWithMinAndMaxInt(t)
		testPopulateFromEnvWithMinAndMaxFloat(t)
	})
	t.Run("Errors", func(t *testing.T) {
		testPopulateFromEnvErrors(t)
	})
	t.Run("Unsupported type", func(t *testing.T) {
		testUnsupportedType(t)
	})
	t.Run("Invalid Value Tag", func(t *testing.T) {
		testInvalidRequiredTag(t)
		testInvalidDefaultIntMinValueTag(t)
		testInvalidDefaultIntMaxValueTag(t)
		testInvalidDefaultFloatMinValueTag(t)
		testInvalidDefaultFloatMaxValueTag(t)
	})
}

func testPopulateFromEnvWithMinAndMaxInt(t *testing.T) {
	type config struct {
		ValueIntMin    int `env:"VALUE_INT_MIN" default:"10" min:"5"`
		ValueIntMax    int `env:"VALUE_INT_MAX" default:"20" max:"25"`
		ValueIntMinMax int `env:"VALUE_INT_MIN_MAX" default:"15" min:"10" max:"20"`
	}

	testcases := []testcase[config]{
		{
			name:   "Default Int Values",
			config: &config{},
			test: func(t *testing.T, tcase testcase[config]) {
				assert.Equal(t, 10, tcase.config.ValueIntMin)
				assert.Equal(t, 20, tcase.config.ValueIntMax)
				assert.Equal(t, 15, tcase.config.ValueIntMinMax)
			},
		},
		{
			config: &config{},
			name:   "Int between Min and Max",
			env: map[string]string{
				"VALUE_INT_MIN_MAX": "9",
			},
			test: func(t *testing.T, tcase testcase[config]) {
				assert.Empty(t, tcase.config.ValueIntMinMax)
			},
			expectedError: &env.EnvError{
				Env: "VALUE_INT_MIN_MAX",
				Msg: "value 9 is less than minimum 10",
			},
		},
		{
			config: &config{},
			name:   "Int between Min and Max",
			env: map[string]string{
				"VALUE_INT_MIN_MAX": "21",
			},
			test: func(t *testing.T, tcase testcase[config]) {
				assert.Empty(t, tcase.config.ValueIntMinMax)
			},
			expectedError: &env.EnvError{
				Env: "VALUE_INT_MIN_MAX",
				Msg: "value 21 is greater than maximum 20",
			},
		},
		{
			config: &config{},
			name:   "Int between Min and Max",
			env: map[string]string{
				"VALUE_INT_MIN_MAX": "15",
			},
			expectedValue: 15,
			test: func(t *testing.T, tcase testcase[config]) {
				assert.Equal(t, tcase.expectedValue, tcase.config.ValueIntMinMax)
			},
		},
		{
			config: &config{},
			name:   "Int Min",
			env: map[string]string{
				"VALUE_INT_MIN": "3",
			},
			test: func(t *testing.T, tcase testcase[config]) {
				assert.Empty(t, tcase.config.ValueIntMin)
			},
			expectedError: &env.EnvError{
				Env: "VALUE_INT_MIN",
				Msg: "value 3 is less than minimum 5",
			},
		},
		{
			config: &config{},
			name:   "Int Min",
			env: map[string]string{
				"VALUE_INT_MIN": "6",
			},
			expectedValue: 6,
			test: func(t *testing.T, tcase testcase[config]) {
				assert.Equal(t, tcase.expectedValue, tcase.config.ValueIntMin)
			},
		},
		{
			config: &config{},
			name:   "Int Max",
			env: map[string]string{
				"VALUE_INT_MAX": "30",
			},
			test: func(t *testing.T, tcase testcase[config]) {
				assert.Empty(t, tcase.config.ValueIntMax)
			},
			expectedError: &env.EnvError{
				Env: "VALUE_INT_MAX",
				Msg: "value 30 is greater than maximum 25",
			},
		},
		{
			config: &config{},
			name:   "Int Max",
			env: map[string]string{
				"VALUE_INT_MAX": "25",
			},
			expectedValue: 25,
			test: func(t *testing.T, tcase testcase[config]) {
				assert.Equal(t, tcase.expectedValue, tcase.config.ValueIntMax)
			},
		},
	}

	test(t, testcases)
}

func testPopulateFromEnvWithMinAndMaxFloat(t *testing.T) {
	type config struct {
		ValueFloatMin    float64 `env:"VALUE_FLOAT_MIN" default:"10.5" min:"5.5"`
		ValueFloatMax    float64 `env:"VALUE_FLOAT_MAX" default:"20.5" max:"25.5"`
		ValueFloatMinMax float64 `env:"VALUE_FLOAT_MIN_MAX" default:"15.5" min:"10.5" max:"20.5"`
	}

	testcases := []testcase[config]{
		{
			name:   "Default Float Values",
			config: &config{},
			test: func(t *testing.T, tcase testcase[config]) {
				assert.Equal(t, 10.5, tcase.config.ValueFloatMin)
				assert.Equal(t, 20.5, tcase.config.ValueFloatMax)
				assert.Equal(t, 15.5, tcase.config.ValueFloatMinMax)
			},
		},
		{
			config: &config{},
			name:   "Float between Min and Max",
			env: map[string]string{
				"VALUE_FLOAT_MIN_MAX": "9",
			},
			test: func(t *testing.T, tcase testcase[config]) {
				assert.Empty(t, tcase.config.ValueFloatMinMax)
			},
			expectedError: &env.EnvError{
				Env: "VALUE_FLOAT_MIN_MAX",
				Msg: "value 9 is less than minimum 10.5",
			},
		},
		{
			config: &config{},
			name:   "Float between Min and Max",
			env: map[string]string{
				"VALUE_FLOAT_MIN_MAX": "21",
			},
			test: func(t *testing.T, tcase testcase[config]) {
				assert.Empty(t, tcase.config.ValueFloatMinMax)
			},
			expectedError: &env.EnvError{
				Env: "VALUE_FLOAT_MIN_MAX",
				Msg: "value 21 is greater than maximum 20.5",
			},
		},
		{
			config: &config{},
			name:   "Float between Min and Max",
			env: map[string]string{
				"VALUE_FLOAT_MIN_MAX": "15.5",
			},
			expectedValue: 15.5,
			test: func(t *testing.T, tcase testcase[config]) {
				assert.Equal(t, tcase.expectedValue, tcase.config.ValueFloatMinMax)
			},
		},
		{
			config: &config{},
			name:   "Float Min",
			env: map[string]string{
				"VALUE_FLOAT_MIN": "3",
			},
			test: func(t *testing.T, tcase testcase[config]) {
				assert.Empty(t, tcase.config.ValueFloatMin)
			},
			expectedError: &env.EnvError{
				Env: "VALUE_FLOAT_MIN",
				Msg: "value 3 is less than minimum 5.5",
			},
		},
		{
			config: &config{},
			name:   "Float Min",
			env: map[string]string{
				"VALUE_FLOAT_MIN": "6.5",
			},
			expectedValue: 6.5,
			test: func(t *testing.T, tcase testcase[config]) {
				assert.Equal(t, tcase.expectedValue, tcase.config.ValueFloatMin)
			},
		},
		{
			config: &config{},
			name:   "Float Max",
			env: map[string]string{
				"VALUE_FLOAT_MAX": "30.5",
			},
			test: func(t *testing.T, tcase testcase[config]) {
				assert.Empty(t, tcase.config.ValueFloatMax)
			},
			expectedError: &env.EnvError{
				Env: "VALUE_FLOAT_MAX",
				Msg: "value 30.5 is greater than maximum 25.5",
			},
		},
		{
			config: &config{},
			name:   "Float Max",
			env: map[string]string{
				"VALUE_FLOAT_MAX": "25.5",
			},
			expectedValue: 25.5,
			test: func(t *testing.T, tcase testcase[config]) {
				assert.Equal(t, tcase.expectedValue, tcase.config.ValueFloatMax)
			},
		},
	}

	test(t, testcases)
}

func testPopulateFromEnvWithRequiredVars(t *testing.T) {
	type config struct {
		ValueString   string        `env:"VALUE_STRING" required:"true"`
		ValueInt      int           `env:"VALUE_INT" required:"true"`
		ValueFloat    float64       `env:"VALUE_FLOAT" required:"true"`
		ValueBool     bool          `env:"VALUE_BOOL" required:"true"`
		ValueDuration time.Duration `env:"VALUE_DURATION" required:"true"`
	}

	testcases := []testcase[config]{
		{
			name:   "Error required String",
			config: &config{},
			env: map[string]string{
				"VALUE_INT":      "1",
				"VALUE_BOOL":     "t",
				"VALUE_FLOAT":    "3.14",
				"VALUE_DURATION": "1s",
			},
			expectedError: &env.EnvError{
				Env: "VALUE_STRING",
				Msg: "is required",
			},
		},
		{
			name:   "Error required Int",
			config: &config{},
			env: map[string]string{
				"VALUE_STRING":   "X",
				"VALUE_BOOL":     "t",
				"VALUE_FLOAT":    "3.14",
				"VALUE_DURATION": "1s",
			},
			expectedError: &env.EnvError{
				Env: "VALUE_INT",
				Msg: "is required",
			},
		},
		{
			name:   "Error required Bool",
			config: &config{},
			env: map[string]string{
				"VALUE_STRING":   "X",
				"VALUE_INT":      "1",
				"VALUE_FLOAT":    "3.14",
				"VALUE_DURATION": "1s",
			},
			expectedError: &env.EnvError{
				Env: "VALUE_BOOL",
				Msg: "is required",
			},
		},
		{
			name:   "Error required Float",
			config: &config{},
			env: map[string]string{
				"VALUE_STRING":   "X",
				"VALUE_INT":      "1",
				"VALUE_DURATION": "1s",
				"VALUE_BOOL":     "t",
			},
			expectedError: &env.EnvError{
				Env: "VALUE_FLOAT",
				Msg: "is required",
			},
		},
		{
			name:   "Error required Duration",
			config: &config{},
			env: map[string]string{
				"VALUE_STRING": "X",
				"VALUE_INT":    "1",
				"VALUE_BOOL":   "t",
				"VALUE_FLOAT":  "3.14",
			},
			expectedError: &env.EnvError{
				Env: "VALUE_DURATION",
				Msg: "is required",
			},
		},
		{
			name:   "Error required",
			config: &config{},
			env: map[string]string{
				"VALUE_STRING":   "X",
				"VALUE_INT":      "1",
				"VALUE_BOOL":     "t",
				"VALUE_FLOAT":    "3.14",
				"VALUE_DURATION": "1s",
			},
			test: func(t *testing.T, tcase testcase[config]) {
				assert.Equal(t, "X", tcase.config.ValueString)
				assert.Equal(t, 1, tcase.config.ValueInt)
				assert.Equal(t, true, tcase.config.ValueBool)
				assert.Equal(t, 3.14, tcase.config.ValueFloat)
				assert.Equal(t, time.Second, tcase.config.ValueDuration)
			},
		},
	}

	test(t, testcases)
}

func testPopulateFromEnvWithDefaultsAndEnvVars(t *testing.T) {
	type config struct {
		ValueString   string        `env:"VALUE_STRING" default:"DEFAULT_STRING"`
		ValueInt      int           `env:"VALUE_INT" default:"42"`
		ValueFloat    float64       `env:"VALUE_FLOAT" default:"3.14"`
		ValueBool     bool          `env:"VALUE_BOOL" default:"true"`
		ValueDuration time.Duration `env:"VALUE_DURATION" default:"1s"`
	}

	testcases := []testcase[config]{
		{
			name:   "Default",
			config: &config{},
			test: func(t *testing.T, tcase testcase[config]) {
				assert.Equal(t, "DEFAULT_STRING", tcase.config.ValueString)
				assert.Equal(t, 42, tcase.config.ValueInt)
				assert.Equal(t, true, tcase.config.ValueBool)
				assert.Equal(t, time.Second, tcase.config.ValueDuration)
			},
		},
		{
			name:   "No Default Values",
			config: &config{},
			env: map[string]string{
				"VALUE_STRING":   "Text",
				"VALUE_INT":      "123",
				"VALUE_FLOAT":    "6.28",
				"VALUE_BOOL":     "False",
				"VALUE_DURATION": "5s",
			},
			test: func(t *testing.T, tcase testcase[config]) {
				assert.Equal(t, "Text", tcase.config.ValueString)
				assert.Equal(t, 123, tcase.config.ValueInt)
				assert.Equal(t, 6.28, tcase.config.ValueFloat)
				assert.Equal(t, false, tcase.config.ValueBool)
				assert.Equal(t, time.Second*time.Duration(5), tcase.config.ValueDuration)
			},
		},
	}

	test(t, testcases)
}

func testPopulateFromEnvErrors(t *testing.T) {
	type config struct {
		ValueInt      int           `env:"VALUE_INT" default:"42"`
		ValueFloat    float64       `env:"VALUE_FLOAT" default:"3.14"`
		ValueBool     bool          `env:"VALUE_BOOL" default:"true"`
		ValueDuration time.Duration `env:"VALUE_DURATION" default:"1s"`
	}

	testcases := []testcase[config]{
		{
			name:   "Invalid Float",
			config: &config{},
			env: map[string]string{
				"VALUE_FLOAT": "NotAFloat",
			},
			expectedError: &env.EnvError{
				Env: "VALUE_FLOAT",
				Msg: "strconv.ParseFloat: parsing \"NotAFloat\": invalid syntax",
			},
		},
		{
			name:   "Invalid Int",
			config: &config{},
			env: map[string]string{
				"VALUE_INT": "NotAnInt",
			},
			expectedError: &env.EnvError{
				Env: "VALUE_INT",
				Msg: "strconv.Atoi: parsing \"NotAnInt\": invalid syntax",
			},
		},
		{
			name:   "Invalid Bool",
			config: &config{},
			env: map[string]string{
				"VALUE_BOOL": "NotABool",
			},
			expectedError: &env.EnvError{
				Env: "VALUE_BOOL",
				Msg: "strconv.ParseBool: parsing \"notabool\": invalid syntax",
			},
		},
		{
			name:   "Invalid Duration",
			config: &config{},
			env: map[string]string{
				"VALUE_DURATION": "NotADuration",
			},
			expectedError: &env.EnvError{
				Env: "VALUE_DURATION",
				Msg: "time: invalid duration \"NotADuration\"",
			},
		},
	}

	test(t, testcases)
}

func testUnsupportedType(t *testing.T) {
	type config struct {
		Value struct{} `env:"VALUE_OBJECT" default:"{}"`
	}

	testcases := []testcase[config]{
		{
			name:   "Unsupported Type Struct",
			config: &config{},
			expectedError: &env.EnvError{
				Env: "VALUE_OBJECT",
				Msg: "unsupported type struct",
			},
		},
	}

	test(t, testcases)
}

func testInvalidRequiredTag(t *testing.T) {
	type config struct {
		ValueString string `env:"VALUE_STRING" required:"yes"`
	}

	testcases := []testcase[config]{
		{
			name:   "Invalid required tag value",
			config: &config{},
			expectedError: &env.EnvError{
				Env: "VALUE_STRING",
				Msg: "field required: strconv.ParseBool: parsing \"yes\": invalid syntax",
			},
		},
	}

	test(t, testcases)
}

func testInvalidDefaultIntMinValueTag(t *testing.T) {
	type config struct {
		ValueInt int `env:"VALUE_INT" min:"NotAnInt" default:"1"`
	}

	testcases := []testcase[config]{
		{
			name:   "Invalid min int value tag",
			config: &config{},
			expectedError: &env.EnvError{
				Env: "VALUE_INT",
				Msg: "field min: strconv.Atoi: parsing \"NotAnInt\": invalid syntax",
			},
		},
	}

	test(t, testcases)
}

func testInvalidDefaultIntMaxValueTag(t *testing.T) {
	type config struct {
		ValueInt int `env:"VALUE_INT" max:"NotAnInt" default:"1"`
	}

	testcases := []testcase[config]{
		{
			name:   "Invalid max int value tag",
			config: &config{},
			expectedError: &env.EnvError{
				Env: "VALUE_INT",
				Msg: "field max: strconv.Atoi: parsing \"NotAnInt\": invalid syntax",
			},
		},
	}

	test(t, testcases)
}

func testInvalidDefaultFloatMinValueTag(t *testing.T) {
	type config struct {
		ValueInt float64 `env:"VALUE_FLOAT" min:"NotAnFloat" default:"1.0"`
	}

	testcases := []testcase[config]{
		{
			name:   "Invalid min float value tag",
			config: &config{},
			expectedError: &env.EnvError{
				Env: "VALUE_FLOAT",
				Msg: "field min: strconv.ParseFloat: parsing \"NotAnFloat\": invalid syntax",
			},
		},
	}

	test(t, testcases)
}

func testInvalidDefaultFloatMaxValueTag(t *testing.T) {
	type config struct {
		ValueInt float64 `env:"VALUE_FLOAT" max:"NotAnFloat" default:"1.0"`
	}

	testcases := []testcase[config]{
		{
			name:   "Invalid max float value tag",
			config: &config{},
			expectedError: &env.EnvError{
				Env: "VALUE_FLOAT",
				Msg: "field max: strconv.ParseFloat: parsing \"NotAnFloat\": invalid syntax",
			},
		},
	}

	test(t, testcases)
}
