package permutation

func Permutation(arr []int, processorFunc func([]int) error) {
	_ = permutationRecursive(arr, processorFunc, 0, len(arr)-1)
}

func permutationRecursive(arr []int, processorFunc func([]int) error, l, r int) error {
	if l == r {
		err := processorFunc(arr)
		if err != nil {
			return err
		}
	} else {
		for i := l; i <= r; i++ {
			arr[l], arr[i] = arr[i], arr[l]

			err := permutationRecursive(arr, processorFunc, l+1, r)
			if err != nil {
				return err
			}

			arr[l], arr[i] = arr[i], arr[l]
		}
	}

	return nil
}
