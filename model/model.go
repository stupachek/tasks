package model

type User struct {
	Name   string // uint8(length) + [length]byte
	Age    uint64 // 1 bit bool (active field) + 63 bit uint (age field)
	Active bool
	Mass   float64  // regular float64
	Books  []string // uint8(all books length) + [length]byte, all books come as a single comma-separated string
}

type UserFromBinary struct {
	Name      []byte
	ActiveAge uint64
	Mass      float64
}
