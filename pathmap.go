package pathmap

import "strings"

type Map map[string]any

/* ---------- public API ---------- */

func (m Map) Set(path string, v any) {
	m.setPath(strings.Split(path, "."), v)
}

func (m Map) Get(path string) (any, bool) {
	return m.getPath(strings.Split(path, "."))
}

func (m Map) Delete(path string) bool {
	parts := strings.Split(path, ".")
	var stack []Map
	cur := m
	for len(parts) > 1 {
		n, ok := cur[parts[0]].(Map)
		if !ok {
			return false // path missing
		}
		stack = append(stack, cur)
		cur = n
		parts = parts[1:]
	}
	if _, ok := cur[parts[0]]; !ok {
		return false
	}
	delete(cur, parts[0])

	// prune empties
	for i := len(stack) - 1; i >= 0; i-- {
		if len(cur) != 0 {
			break
		}
		parent := stack[i]
		key := parts[0] // part from previous loop iteration
		delete(parent, key)
		cur = parent
		parts = parts[:len(parts)-1]
	}
	return true
}

func (m Map) FlattenedKeys() []string {
	// Accumulate here to avoid reallocating slices in recursion.
	keys := make([]string, 0, len(m))

	var walk func(node Map, prefix string)
	walk = func(node Map, prefix string) {
		for k, v := range node {
			full := k
			if prefix != "" {
				full = prefix + "." + k
			}
			if child, ok := v.(Map); ok {
				// Recurse into nested map.
				walk(child, full)
				continue
			}
			keys = append(keys, full)
		}
	}

	walk(m, "")
	return keys
}

/* ---------- internal helpers ---------- */

func (m Map) child(k string) Map {
	if n, ok := m[k].(Map); ok {
		return n
	}
	c := Map{}
	m[k] = c
	return c
}

func (m Map) setPath(parts []string, v any) {
	cur := m
	for len(parts) > 1 {
		cur = cur.child(parts[0])
		parts = parts[1:]
	}
	cur.merge(parts[0], v)
}

func (m Map) merge(k string, v any) {
	if ex, ok := m[k]; ok {
		switch ex := ex.(type) {
		case map[string]any:
			if nv, ok := v.(map[string]any); ok {
				for k2, v2 := range nv {
					ex[k2] = v2
				}
				return
			}
		case []any:
			if nv, ok := v.([]any); ok {
				m[k] = append(ex, nv...)
				return
			}
		}
	}
	m[k] = v
}

func (m Map) getPath(parts []string) (any, bool) {
	cur := m
	for len(parts) > 1 {
		n, ok := cur[parts[0]].(Map)
		if !ok {
			return nil, false
		}
		cur = n
		parts = parts[1:]
	}
	v, ok := cur[parts[0]]
	return v, ok
}
