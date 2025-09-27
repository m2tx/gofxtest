package env_test

import (
	"os"
	"testing"
	"time"

	"github.com/m2tx/gofxtest/internal/env"
	"github.com/stretchr/testify/assert"
)

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

func TestEnv_Get(t *testing.T) {
	t.Run("Get with defaults and env vars", func(t *testing.T) {
		testEnvGetWithDefaultsAndEnvVars(t)
	})
	t.Run("Get with required vars", func(t *testing.T) {
		testEnvGetWithRequiredVars(t)
	})
	t.Run("Get with min and max", func(t *testing.T) {
		testEnvGetWithMinAndMax(t)
	})
}

func testEnvGetWithMinAndMax(t *testing.T) {
	type config struct {
		ValueIntMin    int `env:"VALUE_INT_MIN" default:"10" min:"5"`
		ValueIntMax    int `env:"VALUE_INT_MAX" default:"20" max:"25"`
		ValueIntMinMax int `env:"VALUE_INT_MIN_MAX" default:"15" min:"10" max:"20"`
	}

	c, err := env.Get(&config{})
	assert.NoError(t, err)
	assert.Equal(t, 10, c.ValueIntMin)
	assert.Equal(t, 20, c.ValueIntMax)
	assert.Equal(t, 15, c.ValueIntMinMax)

	os.Setenv("VALUE_INT_MIN_MAX", "9")
	_, err = env.Get(&config{})
	assert.Error(t, err)

	os.Setenv("VALUE_INT_MIN_MAX", "21")
	_, err = env.Get(&config{})
	assert.Error(t, err)

	os.Setenv("VALUE_INT_MIN_MAX", "15")
	c, err = env.Get(&config{})
	assert.NoError(t, err)
	assert.Equal(t, 15, c.ValueIntMinMax)

	os.Setenv("VALUE_INT_MIN", "3")
	_, err = env.Get(&config{})
	assert.Error(t, err)

	os.Setenv("VALUE_INT_MIN", "6")
	c, err = env.Get(&config{})
	assert.NoError(t, err)
	assert.Equal(t, 6, c.ValueIntMin)

	os.Setenv("VALUE_INT_MAX", "30")
	_, err = env.Get(&config{})
	assert.Error(t, err)

	os.Setenv("VALUE_INT_MAX", "24")
	c, err = env.Get(&config{})
	assert.NoError(t, err)
	assert.Equal(t, 24, c.ValueIntMax)

	os.Setenv("VALUE_INT_MAX", "25")
	c, err = env.Get(&config{})
	assert.NoError(t, err)
	assert.Equal(t, 25, c.ValueIntMax)

	os.Setenv("VALUE_INT_MAX", "20")
	c, err = env.Get(&config{})
	assert.NoError(t, err)
	assert.Equal(t, 20, c.ValueIntMax)

	os.Setenv("VALUE_INT_MIN", "")
	os.Setenv("VALUE_INT_MAX", "")
	c, err = env.Get(&config{})
	assert.NoError(t, err)
	assert.Equal(t, 10, c.ValueIntMin)
	assert.Equal(t, 20, c.ValueIntMax)
}

func testEnvGetWithRequiredVars(t *testing.T) {
	type config struct {
		ValueString string `env:"VALUE_STRING" required:"true"`
		ValueInt    int    `env:"VALUE_INT" required:"true"`
		ValueBool   bool   `env:"VALUE_BOOL" required:"true"`
	}

	_, err := env.Get(&config{})
	assert.Error(t, err)

	os.Setenv("VALUE_STRING", "STRING")
	_, err = env.Get(&config{})
	assert.Error(t, err)

	os.Setenv("VALUE_INT", "123")
	_, err = env.Get(&config{})
	assert.Error(t, err)

	os.Setenv("VALUE_BOOL", "T")
	c, err := env.Get(&config{})
	assert.NoError(t, err)

	assert.Equal(t, "STRING", c.ValueString)
	assert.Equal(t, 123, c.ValueInt)
	assert.Equal(t, true, c.ValueBool)
}

func testEnvGetWithDefaultsAndEnvVars(t *testing.T) {
	type config struct {
		ValueString   string        `env:"VALUE_STRING" default:"DEFAULT_STRING"`
		ValueInt      int           `env:"VALUE_INT" default:"42"`
		ValueBool     bool          `env:"VALUE_BOOL" default:"true"`
		ValueDuration time.Duration `env:"VALUE_DURATION" default:"1s"`
	}

	c, err := env.Get(&config{})
	assert.NoError(t, err)
	assert.Equal(t, "DEFAULT_STRING", c.ValueString)
	assert.Equal(t, 42, c.ValueInt)
	assert.Equal(t, true, c.ValueBool)
	assert.Equal(t, time.Second, c.ValueDuration)

	os.Setenv("VALUE_STRING", "STRING")
	c, err = env.Get(&config{})
	assert.NoError(t, err)
	assert.Equal(t, "STRING", c.ValueString)

	os.Setenv("VALUE_INT", "123")
	c, err = env.Get(&config{})
	assert.NoError(t, err)
	assert.Equal(t, 123, c.ValueInt)
	os.Setenv("VALUE_BOOL", "T")
	c, err = env.Get(&config{})
	assert.NoError(t, err)
	assert.Equal(t, true, c.ValueBool)
	os.Setenv("VALUE_BOOL", "F")
	c, err = env.Get(&config{})
	assert.NoError(t, err)
	assert.Equal(t, false, c.ValueBool)
	os.Setenv("VALUE_BOOL", "1")
	c, err = env.Get(&config{})
	assert.NoError(t, err)
	assert.Equal(t, true, c.ValueBool)
	os.Setenv("VALUE_BOOL", "0")
	c, err = env.Get(&config{})
	assert.NoError(t, err)
	assert.Equal(t, false, c.ValueBool)

	os.Setenv("VALUE_DURATION", "2m")
	c, err = env.Get(&config{})
	assert.NoError(t, err)
	assert.Equal(t, 2*time.Minute, c.ValueDuration)

	os.Setenv("VALUE_BOOL", "X")
	c, err = env.Get(&config{})
	assert.Error(t, err)
	assert.Equal(t, false, c.ValueBool)
	os.Setenv("VALUE_BOOL", "")

	os.Setenv("VALUE_INT", "X")
	c, err = env.Get(&config{})
	assert.Error(t, err)
	assert.Equal(t, 0, c.ValueInt)
	os.Setenv("VALUE_INT", "")

	os.Setenv("VALUE_DURATION", "X")
	c, err = env.Get(&config{})
	assert.Error(t, err)
	assert.Equal(t, time.Duration(0), c.ValueDuration)
	os.Setenv("VALUE_DURATION", "")

	os.Setenv("VALUE_STRING", "")
	os.Setenv("VALUE_INT", "")
	os.Setenv("VALUE_BOOL", "")
	os.Setenv("VALUE_DURATION", "")
	c, err = env.Get(&config{})
	assert.NoError(t, err)
	assert.Equal(t, "DEFAULT_STRING", c.ValueString)
	assert.Equal(t, 42, c.ValueInt)
	assert.Equal(t, true, c.ValueBool)
	assert.Equal(t, time.Second, c.ValueDuration)
}
