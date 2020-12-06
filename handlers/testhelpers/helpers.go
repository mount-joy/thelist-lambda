package testhelpers

// BoolToPointer takes a bool as an argument and returns a pointer to that bool
func BoolToPointer(input bool) *bool {
	return &input
}
