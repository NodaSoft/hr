package repository

func cloneMap[K string, V any](m map[K]V) map[K]V {
	nm := make(map[K]V)
	for k, v := range m {
		nm[k] = v
	}
	return nm
}
