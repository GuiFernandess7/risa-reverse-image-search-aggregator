package utils

func Try[T any](v T, err error) (T, bool, error) {
	if err != nil {
		return v, true, err
	}
	return v, false, nil
}
