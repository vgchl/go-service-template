package app

func memoize[T any](f func() T) func() T {
	var instance *T
	return func() T {
		if instance == nil {
			i := f()
			instance = &i
		}
		return *instance
	}
}
