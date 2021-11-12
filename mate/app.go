package mate

import (
	"fmt"
	"reflect"
	"strings"
)

type App struct {
	prefix string
}

func NewApp(prefix string) *App {
	prefix = strings.ReplaceAll(prefix, "-", "_")

	return &App{prefix: strings.ToUpper(prefix)}
}

func (app *App) ScanDefaults(v interface{}) (*EnvVars, error) {
	evs := NewEnvVars(app.prefix)

	if err := evs.To(v); err != nil {
		return nil, fmt.Errorf("failed to v: %w", err)
	}

	if err := evs.From(v); err != nil {
		return nil, fmt.Errorf("failed from v: %w", err)
	}

	return evs, nil
}

func (app *App) Unmarshal(environ []string, v interface{}) error {
	evs := NewEnvVarsFromEnviron(app.prefix, environ)

	if err := evs.To(v); err != nil {
		return fmt.Errorf("failed to v: %w", err)
	}

	if err := triggerInitials(v); err != nil {
		return fmt.Errorf("failed to trigger initials: %w", err)
	}

	return nil
}

func triggerInitials(v interface{}) error {
	rv := reflect.ValueOf(v)
	rv = reflect.Indirect(rv)

	if rv.Kind() == reflect.Struct {
		for i := 0; i < rv.NumField(); i++ {
			value := rv.Field(i)

			if conf, ok := value.Interface().(interface{ Initialize() error }); ok {
				if err := conf.Initialize(); err != nil {
					return fmt.Errorf("failed to initialize: %w", err)
				}
			}
		}
	} else {
		if conf, ok := rv.Interface().(interface{ Initialize() error }); ok {
			if err := conf.Initialize(); err != nil {
				return fmt.Errorf("failed to initialize: %w", err)
			}
		}
	}

	return nil
}

func (app *App) View(v interface{}) (string, error) {
	evs := NewEnvVars(app.prefix)
	if err := evs.From(v); err != nil {
		return "", fmt.Errorf("failed from v: %w", err)
	}

	return evs.View(), nil
}
