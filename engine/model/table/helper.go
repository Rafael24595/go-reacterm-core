package table

func RowCount(headers []string, cols map[string][]string) uint16 {
	maxRows := uint16(0)
	for _, h := range headers {
		maxRows = max(
			maxRows, uint16(len(cols[h])),
		)
	}
	return maxRows
}
