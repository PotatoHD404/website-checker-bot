package utils

func Contains(tables []string, table string) bool {
	for _, t := range tables {
		if t == table {
			return true
		}
	}
	return false
}
