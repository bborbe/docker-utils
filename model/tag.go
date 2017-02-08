package model

type Tag string

func (t Tag) String() string {
	return string(t)
}

type TagsByName []Tag

func (t TagsByName) Len() int {
	return len(t)
}

func (t TagsByName) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t TagsByName) Less(i, j int) bool {
	return t[i] < t[j]
}
