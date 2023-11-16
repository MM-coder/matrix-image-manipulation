package main

// Make2D makes a 2D slice of any type of the given width and height
func Make2D[Type any](n, m int) [][]Type {
	matrix := make([][]Type, n)
	rows := make([]Type, n*m)
	for i, startRow := 0, 0; i < n; i, startRow = i+1, startRow+m {
		endRow := startRow + m
		matrix[i] = rows[startRow:endRow:endRow]
	}
	return matrix
}
