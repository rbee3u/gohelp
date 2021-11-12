package mate

import (
	"encoding"
	"fmt"
	"go/ast"
	"reflect"
	"strconv"
)

func (evs *EnvVars) From(v interface{}) error {
	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(v)
	}

	e := &encoder{envVars: evs, pw: NewPathWalker()}

	return e.from(rv)
}

type encoder struct {
	envVars *EnvVars
	pw      *PathWalker
}

func (e *encoder) from(rv reflect.Value) error {
	setValue := func(value string) {
		e.envVars.Set(&EnvVar{Key: e.pw.String(), Value: value})
	}

	if rv.CanInterface() {
		if textMarshaler, ok := rv.Interface().(encoding.TextMarshaler); ok {
			text, err := textMarshaler.MarshalText()
			if err != nil {
				return fmt.Errorf("failed to marshal text: %w", err)
			}

			setValue(string(text))

			return nil
		}
	}

	switch rv.Kind() {
	case reflect.Ptr:
		if rv.IsNil() {
			return nil
		}

		return e.from(rv.Elem())
	case reflect.Slice, reflect.Array:
		for i := 0; i < rv.Len(); i++ {
			e.pw.Enter(IntegerPath(i))

			err := e.from(rv.Index(i))

			e.pw.Exit()

			if err != nil {
				return err
			}
		}

		return nil
	case reflect.Struct:
		rt := rv.Type()
		for i := 0; i < rv.NumField(); i++ {
			field := rt.Field(i)
			name := field.Name

			if !ast.IsExported(name) {
				continue
			}

			var (
				tagValue string
				tagFlags map[string]bool
			)

			if tag, ok := field.Tag.Lookup("env"); ok {
				tagValue, tagFlags = parseTag(tag)
				if tagValue == "-" {
					continue
				}

				if len(tagValue) != 0 {
					name = tagValue
				}
			}

			ft := field.Type
			for ; ft.Kind() == reflect.Ptr; ft = ft.Elem() {
			}

			squash := ft.Kind() == reflect.Struct &&
				len(tagValue) == 0 &&
				(field.Anonymous || tagFlags["squash"])

			if !squash {
				e.pw.Enter(StringPath(name))
			}

			err := e.from(rv.Field(i))

			if !squash {
				e.pw.Exit()
			}

			if err != nil {
				return err
			}
		}

		return nil
	}

	switch rv.Kind() {
	case reflect.Bool:
		setValue(strconv.FormatBool(rv.Bool()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		setValue(strconv.FormatInt(rv.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		setValue(strconv.FormatUint(rv.Uint(), 10))
	case reflect.Float32:
		setValue(strconv.FormatFloat(rv.Float(), 'g', -1, 32))
	case reflect.Float64:
		setValue(strconv.FormatFloat(rv.Float(), 'g', -1, 64))
	case reflect.String:
		setValue(rv.String())
	default:
		return fmt.Errorf("unsupported rv.Type(): %s", rv.Type())
	}

	return nil
}
