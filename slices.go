package brbundle

func mergeStringSlices(sorted ...[]string) []string {
	length := 0
	sourceSlices := make([][]string, 0, len(sorted))
	for _, src := range sorted {
		if len(src) > 0 {
			length += len(src)
			sourceSlices = append(sourceSlices, src)
		}
	}
	if length == 0 {
		return nil
	} else if length == 1 {
		return sourceSlices[0]
	}
	result := make([]string, length)
	sourceSliceCount := len(sourceSlices)
	indexes := make([]int, sourceSliceCount)
	index := 0
	for {
		minSlice := 0
		minItem := sourceSlices[0][indexes[0]]
		for i := 1; i < sourceSliceCount; i++ {
			if sourceSlices[i][indexes[i]] < minItem {
				minSlice = i
				minItem = sourceSlices[i][indexes[i]]
			}
		}
		result[index] = minItem
		index++
		indexes[minSlice]++
		if indexes[minSlice] == len(sourceSlices[minSlice]) {
			sourceSlices = append(sourceSlices[:minSlice], sourceSlices[minSlice+1:]...)
			indexes = append(indexes[:minSlice], indexes[minSlice+1:]...)
			sourceSliceCount--
			if len(sourceSlices) == 1 {
				copy(result[index:], sourceSlices[0][indexes[0]:])
				return result
			}
		}
	}
}

func iterateOverStringSlices(callback func(item string, srcIndex int), sorted ...[]string) {
	sourceSlices := make([][]string, 0, len(sorted))
	for _, src := range sorted {
		if len(src) > 0 {
			sourceSlices = append(sourceSlices, src)
		}
	}
	sourceSliceCount := len(sourceSlices)
	if sourceSliceCount == 0 {
		return
	} else if sourceSliceCount == 1 {
		for i, value := range sourceSlices[0] {
			callback(value, i)
		}
		return
	}
	indexes := make([]int, sourceSliceCount)
	sliceIndex := make([]int, sourceSliceCount)
	for i := range sourceSlices {
		sliceIndex[i] = i
	}
	index := 0
	for {
		minSlice := 0
		minItem := sourceSlices[0][indexes[0]]
		for i := 1; i < sourceSliceCount; i++ {
			if sourceSlices[i][indexes[i]] < minItem {
				minSlice = i
				minItem = sourceSlices[i][indexes[i]]
			}
		}
		callback(minItem, sliceIndex[minSlice])
		index++
		indexes[minSlice]++
		if indexes[minSlice] == len(sourceSlices[minSlice]) {
			sourceSlices = append(sourceSlices[:minSlice], sourceSlices[minSlice+1:]...)
			indexes = append(indexes[:minSlice], indexes[minSlice+1:]...)
			sliceIndex = append(sliceIndex[:minSlice], sliceIndex[minSlice+1:]...)
			sourceSliceCount--
			if len(sourceSlices) == 1 {
				slice := sourceSlices[0]
				for i := indexes[0]; i < len(slice); i++ {
					callback(slice[i], sliceIndex[0])
				}
				return
			}
		}
	}
}
