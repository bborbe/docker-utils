package docker

type TagName string

func (t TagName) String() string {
	return string(t)
}

type TagsByName []TagName

func (t TagsByName) Len() int {
	return len(t)
}

func (t TagsByName) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t TagsByName) Less(i, j int) bool {
	return t[i] < t[j]
}
