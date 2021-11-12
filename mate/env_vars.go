package mate

import (
	"bytes"
	"sort"
	"strconv"
	"strings"
)

func NewEnvVars(prefix string) *EnvVars {
	return &EnvVars{Prefix: prefix}
}

func NewEnvVarsFromEnviron(prefix string, environ []string) *EnvVars {
	evs := NewEnvVars(prefix)

	for _, item := range environ {
		kv := strings.SplitN(item, "=", 2)
		if len(kv) != 2 {
			continue
		}

		k := strings.ToUpper(kv[0])
		if !strings.HasPrefix(k, prefix) {
			continue
		}

		evs.Set(&EnvVar{
			Key:   strings.TrimPrefix(k, prefix),
			Value: kv[1],
		})
	}

	return evs
}

type EnvVars struct {
	Prefix string
	Dict   map[string]*EnvVar
}

type EnvVar struct {
	Key   string
	Value string
}

func (ev *EnvVar) RealKey(prefix string) string {
	return prefix + ev.Key
}

func (evs *EnvVars) Set(envVar *EnvVar) {
	if evs.Dict == nil {
		evs.Dict = map[string]*EnvVar{}
	}

	evs.Dict[strings.ToUpper(envVar.Key)] = envVar
}

func (evs *EnvVars) Len(key string) int {
	maxIdx := -1

	for _, envVar := range evs.Dict {
		keyPath := strings.ToUpper(envVar.Key)
		k := strings.ToUpper(key)

		if strings.HasPrefix(keyPath, k) {
			v := strings.TrimLeft(keyPath, k+"_")
			parts := strings.Split(v, "_")

			i, err := strconv.ParseInt(parts[0], 10, 64)
			if err == nil {
				if int(i) > maxIdx {
					maxIdx = int(i)
				}
			}
		}
	}

	return maxIdx + 1
}

func (evs *EnvVars) Get(key string) *EnvVar {
	if evs.Dict == nil {
		return nil
	}

	return evs.Dict[strings.ToUpper(key)]
}

func (evs *EnvVars) View() string {
	dict := map[string]string{}
	for _, ev := range evs.Dict {
		dict[ev.RealKey(evs.Prefix)] = ev.Value
	}

	keys := make([]string, 0)
	for key := range dict {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	buf := new(bytes.Buffer)

	for _, key := range keys {
		buf.WriteString(key)
		buf.WriteByte('=')
		buf.WriteString(dict[key])
		buf.WriteByte('\n')
	}

	return buf.String()
}
