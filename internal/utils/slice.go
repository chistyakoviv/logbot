package utils

func Chunk[T any](items []T, size int) [][]T {
	if size <= 0 {
		return nil
	}

	var chunks [][]T
	for size < len(items) {
		items, chunks = items[size:], append(chunks, items[0:size:size])
	}
	chunks = append(chunks, items)
	return chunks
}
