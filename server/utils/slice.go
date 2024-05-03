package utils

func SlicePagingStr(s []string, page, pageSize int) []string {
	startIdx := (page - 1) * pageSize
	endIdx := startIdx + pageSize
	if startIdx >= len(s) {
		return nil
	}
	if endIdx > len(s) {
		endIdx = len(s)
	}

	return s[startIdx:endIdx]
}
