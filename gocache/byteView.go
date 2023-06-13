package gocache

type ByteView struct {
	b []byte
	s string
}

func (view *ByteView) Len() int {
	if view.b != nil {
		return len(view.b)
	}
	return len(view.s)
}

// ByteSlice return a copy of data which storage in ByteView through the way in byte array
func (view *ByteView) ByteSlice() []byte {
	if view.b != nil {
		// return a copy of b
		return cloneBytes(view.b)
	}
	return ([]byte)(view.s)
}

func (view *ByteView) String() string {
	if view.b != nil {
		return string(view.b)
	}
	return view.s
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
func (view *ByteView) SetString(s string) {
	view.s = s
}
func (view *ByteView) SetBytes(b []byte) {
	view.b = b
}

func (view *ByteView) At() {

}
