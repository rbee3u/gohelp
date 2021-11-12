package mate

import (
	"encoding"
	"fmt"
	"go/ast"
	"reflect"
	"strconv"
)

func (evs *EnvVars) To(v interface{}) error {
	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(v)
	}

	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("invalid rv.Kind(): %s != %s", rv.Kind(), reflect.Ptr)
	}

	d := &decoder{evs: evs, pw: NewPathWalker()}

	return d.to(rv)
}

type decoder struct {
	evs *EnvVars
	pw  *PathWalker
}

func (d *decoder) to(rv reflect.Value) error {
	if rv.CanAddr() {
		rp := rv.Addr()

		if defaultsSetter, ok := rp.Interface().(interface{ SetDefaults() error }); ok {
			if err := defaultsSetter.SetDefaults(); err != nil {
				return fmt.Errorf("failed to set defaults: %w", err)
			}
		}

		if textUnmarshaler, ok := rp.Interface().(encoding.TextUnmarshaler); ok {
			ev := d.evs.Get(d.pw.String())
			if ev == nil {
				return nil
			}

			if err := textUnmarshaler.UnmarshalText([]byte(ev.Value)); err != nil {
				return fmt.Errorf("failed to unmarshal text: %w", err)
			}

			return nil
		}
	}

	switch rv.Kind() {
	case reflect.Ptr:
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}

		return d.to(rv.Elem())
	case reflect.Slice:
		if rv.IsNil() {
			n := d.evs.Len(d.pw.String())
			rv.Set(reflect.MakeSlice(rv.Type(), n, n))
		}

		fallthrough
	case reflect.Array:
		for i := 0; i < rv.Len(); i++ {
			d.pw.Enter(IntegerPath(i))

			err := d.to(rv.Index(i))

			d.pw.Exit()

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
				d.pw.Enter(StringPath(name))
			}

			err := d.to(rv.Field(i))

			if !squash {
				d.pw.Exit()
			}

			if err != nil {
				return err
			}
		}

		return nil
	}

	ev := d.evs.Get(d.pw.String())
	if ev == nil {
		return nil
	}

	switch rv.Kind() {
	case reflect.Bool:
		x, err := strconv.ParseBool(ev.Value)
		if err != nil {
			return fmt.Errorf("failed to parse bool: %w", err)
		}

		rv.SetBool(x)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		x, err := strconv.ParseInt(ev.Value, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse int: %w", err)
		}

		rv.SetInt(x)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		x, err := strconv.ParseUint(ev.Value, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse uint: %w", err)
		}

		rv.SetUint(x)
	case reflect.Float32, reflect.Float64:
		x, err := strconv.ParseFloat(ev.Value, 64)
		if err != nil {
			return fmt.Errorf("failed to parse float: %w", err)
		}

		rv.SetFloat(x)
	case reflect.String:
		rv.SetString(ev.Value)
	default:
		return fmt.Errorf("unsupported rv.Type(): %s", rv.Type())
	}

	return nil
}
