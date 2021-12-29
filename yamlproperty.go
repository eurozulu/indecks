package indecks

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v3"
	"strings"
)

type YamlProperty struct {
	Names []string
}

func (y YamlProperty) ParseProperties(p string, by []byte) (map[string]Property, error) {
	m := map[string]interface{}{}
	if err := yaml.Unmarshal(by, &m); err != nil {
		return nil, fmt.Errorf("failed to parse %s as  %w", p, err)
	}

	pm := map[string]Property{}
	var mName string
	for _, n := range y.Names {
		ks := strings.Split(n, ".")
		last := len(ks) - 1
		// navigate down to map containting last elements key(s)
		if last > 0 {
			mName = strings.Join(ks[:last], ".")
			m = childMap(m, ks[:last])
			if m == nil {
				// ignore name as no key not found
				continue
			}
		}

		var kns []string
		if ks[last] == "*" {
			kns = Keynames(m)
		} else {
			kns = []string{ks[last]}
		}
		for _, k := range kns {
			if v, ok := m[k]; ok {
				kn := strings.Join([]string{mName, k}, ".")
				pm[kn] = bytes.NewBufferString(fmt.Sprintf("%v", v))
			}
		}
	}
	return pm, nil
}

func childMap(m map[string]interface{}, ks []string) map[string]interface{} {
	for len(ks) > 0 {
		v, ok := m[ks[0]]
		if !ok {
			return nil
		}
		m, ok = v.(map[string]interface{})
		if !ok {
			return nil
		}
		ks = ks[1:]
	}
	return m
}

func Keynames(m map[string]interface{}) []string {
	names := make([]string, len(m))
	var i int
	for k := range m {
		names[i] = k
		i++
	}
	return names
}
