package ptr

func Of[T any](v T) *T {
	return &v
}

func True() *bool {
	return Of(true)
}

func False() *bool {
	return Of(false)
}

func Float64[T int | float64](t T) *float64 {
	return Of(float64(t))
}
