package onionext

import (
	"os"
	"strings"

	"gopkg.in/fzerorubigd/onion.v2"
)

type envLoader struct {
	loaded         bool
	prefix         string
	envDelimiter   string
	onionDelimiter string
	data           map[string]interface{}
}

func (el *envLoader) Load() (map[string]interface{}, error) {
	if el.loaded {
		return el.data, nil
	}

	for _, env := range os.Environ() {
		kv := strings.SplitN(env, "=", 2)
		k, v := kv[0], kv[1]

		if strings.HasPrefix(k, el.prefix) {
			k = strings.TrimPrefix(k, el.prefix)
			k = strings.ToLower(k)

			where := el.data
			parts := strings.Split(k, el.envDelimiter)
			last, parts := parts[len(parts)-1], parts[:len(parts)-1]
			for _, part := range parts {
				next, found := where[part]
				if !found {
					next = make(map[string]interface{})
					where[part] = next
				}

				switch next.(type) {
				case string:
					n := make(map[string]interface{})
					n[""] = next
					next = n
					where[part] = next
				}

				where = next.(map[string]interface{})
			}

			where[last] = v
		}
	}
	el.loaded = true
	return el.data, nil
}

// NewEnvLayer create a environment loader.
func NewEnvLayer(prefix string, envDelimiter string, onionDelimiter string) onion.Layer {
	prefix = strings.ToLower(prefix)
	if prefix != "" {
		prefix = prefix + envDelimiter
	}
	return &envLoader{false, prefix, envDelimiter, onionDelimiter, make(map[string]interface{})}
}
