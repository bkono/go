package slice

func Diff[T comparable](slice1, slice2 []T) ([]T, []T) {
	// Create a map to keep track of the elements from the first slice
	knownMap := make(map[T]struct{})

	// Iterate over the first slice, adding each element to the map
	for _, item := range slice1 {
		knownMap[item] = struct{}{}
	}

	// Create a slice diff2 to store the differences from the second slice
	var diff2 []T

	// Iterate over the second slice, updating the known map and building diff2
	for _, item := range slice2 {
		if _, ok := knownMap[item]; ok {
			// If the element is in the known map, remove it
			delete(knownMap, item)
		} else {
			// If the element is not in the known map, add it to diff2
			diff2 = append(diff2, item)
		}
	}

	// Convert the remaining elements in the known map to the diff1 output slice
	var diff1 []T
	for item := range knownMap {
		diff1 = append(diff1, item)
	}

	// Return the two difference slices diff1 and diff2
	return diff1, diff2
}
