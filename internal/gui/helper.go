package gui

func toPointer[T any](val T) *T {
	return &val
}

// Converts slice []T to []E
func mapNewSlice[T any, E any](items []T, callback func(item T) E) []E {
	if items == nil {
		return nil
	}

	output := make([]E, len(items))
	for i, item := range items {
		output[i] = callback(item)
	}

	return output
}
