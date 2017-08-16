package heap

// +gen container:"HeapBy"
type S1 struct {
	fld int
}

func S1Less(a, b S1) bool {
	return a.fld < b.fld
}

// +gen * container:"HeapBy"
type S2 struct {
	fld int
}

func S2Less(a, b *S2) bool {
	return a.fld < b.fld
}
