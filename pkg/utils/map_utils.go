package utils

import "sort"

func SortedIntMap(m map[int]string) map[int]string {
	sortedKeys := make([]int, 0, len(m))
	for k := range m {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Ints(sortedKeys)

	sortedMap := make(map[int]string)
	for _, k := range sortedKeys {
		sortedMap[k] = m[k]
	}
	return sortedMap
}

func IntMapWithOffset(allSpecialities map[int]string, offset int) map[int]string {
	sortedKeys := make([]int, 0, len(allSpecialities))
	for k := range allSpecialities {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Ints(sortedKeys)
	offsetMap := make(map[int]string)
	if offset < len(sortedKeys) {
		for _, k := range sortedKeys[offset:] {
			offsetMap[k] = allSpecialities[k]
		}
	}
	return offsetMap
}
