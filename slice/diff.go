package slice

// Diff takes two slices s1 and s2 of any comparable type T as input
// and returns three slices: s1Only, s2Only, and common. The s1Only slice contains
// elements that are unique to s1, the s2Only slice contains elements that are
// unique to s2, and the common slice contains elements that are common to both s1 and s2.
func Diff[T comparable](s1, s2 []T) (s1Only, s2Only, common []T) {
	// Create a map to keep track of the elements from the first slice
	knownMap := make(map[T]struct{})

	// Iterate over the first slice, adding each element to the map
	for _, item := range s1 {
		knownMap[item] = struct{}{}
	}

	// Iterate over the second slice, updating the known map and building s2Only and common
	for _, item := range s2 {
		if _, ok := knownMap[item]; ok {
			// If the element is in the known map, add it to common and remove it from knownMap
			common = append(common, item)
			delete(knownMap, item)
		} else {
			// If the element is not in the known map, add it to s2Only
			s2Only = append(s2Only, item)
		}
	}

	// Convert the remaining elements in the known map to the s1Only output slice
	for item := range knownMap {
		s1Only = append(s1Only, item)
	}

	// Return the three slices s1Only, s2Only, and common
	return s1Only, s2Only, common
}
