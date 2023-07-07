package util

func Map[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i := range ts {
			us[i] = f(ts[i])
	}
	return us
}

func Filter[T any](ss []T, test func(T) bool) (ret []T) {
	for _, s := range ss {
			if test(s) {
					ret = append(ret, s)
			}
	}
	return
}