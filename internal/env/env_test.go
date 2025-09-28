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
					os.Setenv(k, v)

				}
				defer func() {
					for k := range tcase.env {
						os.Unsetenv(k)
					}
				}()
			}

			err := env.PopulateFromEnv(tcase.config)
			if tcase.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tcase.expectedError, err)
			}

			if tcase.test != nil {
				tcase.test(t, tcase)
			}
		})
	}
}

func TestEnv_New(t *testing.T) {
	type config struct {
		ValueString string `env:"VALUE_STRING" default:"DEFAULT_STRING"`
		ValueInt    int    `env:"VALUE_INT" default:"42"`
		ValueBool   bool   `env:"VALUE_BOOL" default:"true"`
	}

	c, err := env.New[config]()
	assert.NoError(t, err)
	assert.Equal(t, "DEFAULT_STRING", c.ValueString)
	assert.Equal(t, 42, c.ValueInt)
	assert.Equal(t, true, c.ValueBool)
}

func TestEnv_PopulateFromEnv(t *testing.T) {
	t.Run("Populate with defaults and env vars", func(t *testing.T) {
		testPopulateFromEnvWithDefaultsAndEnvVars(t)
	})
	t.Run("Populate with required vars", func(t *testing.T) {
		testPopulateFromEnvWithRequiredVars(t)
	})
	t.Run("Populate with min and max", func(t *testing.T) {
		testPopulateFromEnvWithMinAndMax(t)
	})
}

func testPopulateFromEnvWithMinAndMax(t *testing.T) {
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

func testPopulateFromEnvWithRequiredVars(t *testing.T) {
	type config struct {
		ValueString string `env:"VALUE_STRING" required:"true"`
		ValueInt    int    `env:"VALUE_INT" required:"true"`
		ValueBool   bool   `env:"VALUE_BOOL" required:"true"`
	}

	testcases := []testcase[config]{
		{
			name:   "Error required",
			config: &config{},
			env: map[string]string{
				"VALUE_INT":  "1",
				"VALUE_BOOL": "t",
			},
			expectedError: &env.EnvError{
				Env: "VALUE_STRING",
				Msg: "is required",
			},
		},
		{
			name:   "Error required",
			config: &config{},
			env: map[string]string{
				"VALUE_STRING": "X",
				"VALUE_BOOL":   "t",
			},
			expectedError: &env.EnvError{
				Env: "VALUE_INT",
				Msg: "is required",
			},
		},
		{
			name:   "Error required",
			config: &config{},
			env: map[string]string{
				"VALUE_STRING": "X",
				"VALUE_INT":    "1",
			},
			expectedError: &env.EnvError{
				Env: "VALUE_BOOL",
				Msg: "is required",
			},
		},
		{
			name:   "Error required",
			config: &config{},
			env: map[string]string{
				"VALUE_STRING": "X",
				"VALUE_INT":    "1",
				"VALUE_BOOL":   "t",
			},
			test: func(t *testing.T, tcase testcase[config]) {
				assert.Equal(t, "X", tcase.config.ValueString)
				assert.Equal(t, 1, tcase.config.ValueInt)
				assert.Equal(t, true, tcase.config.ValueBool)
			},
		},
	}

	test(t, testcases)
}

func testPopulateFromEnvWithDefaultsAndEnvVars(t *testing.T) {
	type config struct {
		ValueString   string        `env:"VALUE_STRING" default:"DEFAULT_STRING"`
		ValueInt      int           `env:"VALUE_INT" default:"42"`
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
				"VALUE_BOOL":     "False",
				"VALUE_DURATION": "5s",
			},
			test: func(t *testing.T, tcase testcase[config]) {
				assert.Equal(t, "Text", tcase.config.ValueString)
				assert.Equal(t, 123, tcase.config.ValueInt)
				assert.Equal(t, false, tcase.config.ValueBool)
				assert.Equal(t, time.Second*time.Duration(5), tcase.config.ValueDuration)
			},
		},
	}

	test(t, testcases)
}
