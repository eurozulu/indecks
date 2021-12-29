package indecks

func appendMap(m map[string]Property, ms map[string]Property) map[string]Property {
	if len(ms) == 0 {
		return m
	}
	if m == nil {
		m = map[string]Property{}
	}
	for k, v := range ms {
		m[k] = v
	}
	return m
}
