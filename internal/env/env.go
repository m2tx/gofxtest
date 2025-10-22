package env

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	TagEnv      = "env"
	TagDefault  = "default"
	TagRequired = "required"

	TagMin = "min"
	TagMax = "max"

	errorMsgUnsupportedType = "unsupported type %s"
	errorMsgRequired        = "is required"
	errorFieldEnv           = "field %s: %s"
	errorMsgMinValue        = "value %v is less than minimum %v"
	errorMsgMaxValue        = "value %v is greater than maximum %v"
	errorEnv                = "environment variable %s: %s"

	bitSizeFloat = 64
)

type EnvError struct {
	Env string
	Msg string
}

func (e *EnvError) Error() string {
	return fmt.Sprintf(errorEnv, e.Env, e.Msg)
}

func New[T any]() (T, error) {
	var t T
	err := PopulateFromEnv(&t)
	return t, err
}

func PopulateFromEnv[T any](structEnv *T) error {
	fields := reflect.VisibleFields(reflect.TypeOf(*structEnv))
	for _, f := range fields {
		err := setEnv(f, structEnv)
		if err != nil {
			return err
		}
	}

	return nil
}

func setEnv[T any](f reflect.StructField, structEnv *T) error {
	tagEnv := f.Tag.Get(TagEnv)
	if tagEnv != "" {
		tagEnvDefault := f.Tag.Get(TagDefault)
		isRequired, err := tagToBool(f.Tag, TagRequired, tagEnv)
		if err != nil {
			return err
		}

		fv := reflect.ValueOf(structEnv).Elem().FieldByName(f.Name)
		if fv.CanSet() {
			switch fv.Interface().(type) {
			case time.Duration:
				vDuration, err := getEnvDuration(tagEnv, tagEnvDefault, isRequired)
				if err != nil {
					return err
				}

				fv.Set(reflect.ValueOf(vDuration))
			case string:
				vString, err := getEnvString(tagEnv, tagEnvDefault, isRequired)
				if err != nil {
					return err
				}

				fv.SetString(vString)
			case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
				vInt, err := getEnvInt(tagEnv, tagEnvDefault, isRequired, f.Tag)
				if err != nil {
					return err
				}

				fv.SetInt(int64(vInt))
			case float32, float64:
				vFloat, err := getEnvFloat(tagEnv, tagEnvDefault, isRequired, f.Tag)
				if err != nil {
					return err
				}

				fv.SetFloat(vFloat)
			case bool:
				vBool, err := getEnvBool(tagEnv, tagEnvDefault, isRequired)
				if err != nil {
					return err
				}

				fv.SetBool(vBool)
			default:
				return &EnvError{
					Env: tagEnv,
					Msg: fmt.Sprintf(errorMsgUnsupportedType, fv.Kind().String()),
				}
			}
		}
	}
	return nil
}

func tagToInt(tag reflect.StructTag, tagKey string, envKey string) (int, bool, error) {
	v := tag.Get(tagKey)
	if v == "" {
		return 0, false, nil
	}

	vInt, err := strconv.Atoi(v)
	if err != nil {
		return 0, true, &EnvError{
			Env: envKey,
			Msg: fmt.Sprintf(errorFieldEnv, tagKey, err.Error()),
		}
	}

	return vInt, true, nil
}

func tagToFloat(tag reflect.StructTag, tagKey string, envKey string) (float64, bool, error) {
	v := tag.Get(tagKey)
	if v == "" {
		return 0, false, nil
	}

	vFloat, err := strconv.ParseFloat(v, bitSizeFloat)
	if err != nil {
		return 0, true, &EnvError{
			Env: envKey,
			Msg: fmt.Sprintf(errorFieldEnv, tagKey, err.Error()),
		}
	}

	return vFloat, true, nil
}

func tagToBool(tag reflect.StructTag, tagKey string, envKey string) (bool, error) {
	v := tag.Get(tagKey)
	if v == "" {
		return false, nil
	}

	vBool, err := strconv.ParseBool(v)
	if err != nil {
		return false, &EnvError{
			Env: envKey,
			Msg: fmt.Sprintf(errorFieldEnv, tagKey, err.Error()),
		}
	}

	return vBool, nil
}

func getEnvDuration(key string, defaultValue string, required bool) (time.Duration, error) {
	v := os.Getenv(key)
	if v == "" && defaultValue != "" {
		v = defaultValue
	}

	if v == "" && required {
		return 0, &EnvError{
			Env: key,
			Msg: errorMsgRequired,
		}
	}

	vDuration, err := time.ParseDuration(v)
	if err != nil {
		return 0, &EnvError{
			Env: key,
			Msg: err.Error(),
		}
	}

	return vDuration, nil
}

func getEnvInt(key string, defaultValue string, required bool, tag reflect.StructTag) (int, error) {
	v := os.Getenv(key)
	if v == "" && defaultValue != "" {
		v = defaultValue
	}

	if v == "" && required {
		return 0, &EnvError{
			Env: key,
			Msg: errorMsgRequired,
		}
	}

	vInt, err := strconv.Atoi(v)
	if err != nil {
		return 0, &EnvError{
			Env: key,
			Msg: err.Error(),
		}
	}

	min, ok, err := tagToInt(tag, TagMin, key)
	if err != nil {
		return 0, err
	}

	if vInt < min && ok {
		return 0, &EnvError{
			Env: key,
			Msg: fmt.Sprintf(errorMsgMinValue, vInt, min),
		}
	}

	max, ok, err := tagToInt(tag, TagMax, key)
	if err != nil {
		return 0, err
	}

	if vInt > max && ok {
		return 0, &EnvError{
			Env: key,
			Msg: fmt.Sprintf(errorMsgMaxValue, vInt, max),
		}
	}

	return vInt, nil
}

func getEnvString(key string, defaultValue string, required bool) (string, error) {
	v := os.Getenv(key)
	if v == "" && defaultValue != "" {
		v = defaultValue
	}

	if v == "" && required {
		return "", &EnvError{
			Env: key,
			Msg: errorMsgRequired,
		}
	}

	return v, nil
}

func getEnvBool(key string, defaultValue string, required bool) (bool, error) {
	v := os.Getenv(key)
	if v == "" && defaultValue != "" {
		v = defaultValue
	}

	if v == "" && required {
		return false, &EnvError{
			Env: key,
			Msg: errorMsgRequired,
		}
	}

	vBool, err := strconv.ParseBool(strings.ToLower(v))
	if err != nil {
		return false, &EnvError{
			Env: key,
			Msg: err.Error(),
		}
	}

	return vBool, nil
}

func getEnvFloat(key string, defaultValue string, required bool, tag reflect.StructTag) (float64, error) {
	v := os.Getenv(key)
	if v == "" && defaultValue != "" {
		v = defaultValue
	}

	if v == "" && required {
		return 0, &EnvError{
			Env: key,
			Msg: errorMsgRequired,
		}
	}

	vFloat, err := strconv.ParseFloat(v, bitSizeFloat)
	if err != nil {
		return 0, &EnvError{
			Env: key,
			Msg: err.Error(),
		}
	}

	min, ok, err := tagToFloat(tag, TagMin, key)
	if err != nil {
		return 0, err
	}

	if vFloat < min && ok {
		return 0, &EnvError{
			Env: key,
			Msg: fmt.Sprintf(errorMsgMinValue, vFloat, min),
		}
	}

	max, ok, err := tagToFloat(tag, TagMax, key)
	if err != nil {
		return 0, err
	}

	if vFloat > max && ok {
		return 0, &EnvError{
			Env: key,
			Msg: fmt.Sprintf(errorMsgMaxValue, vFloat, max),
		}
	}

	return vFloat, nil
}
