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

func IsUintInSlice(nums []uint, target uint) bool {
	for _, num := range nums {
		if num == target {
			return true
		}
	}
	return false
}
