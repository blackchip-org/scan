package scan

// Repeat calls function fn in a loop for n times.
func Repeat(fn func(), n int) {
	for i := 0; i < n; i++ {
		fn()
	}
}

// While calls function fn while this rune is a member of class c.
func While(s *Scanner, c Class, fn func()) {
	for c(s.This) {
		fn()
	}
}
