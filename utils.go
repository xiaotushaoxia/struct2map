package struct2map

// mergeMap merge ms to m1, return m1 to allow chaining
func mergeMap(m1 map[string]any, ms ...map[string]any) map[string]any {
	for _, m := range ms {
		for s, a := range m {
			m1[s] = a
		}
	}
	return m1
}
