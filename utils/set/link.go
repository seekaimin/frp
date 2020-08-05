package set

type LinkItem struct {
	pre   *LinkItem
	next  *LinkItem
	value interface{}
}

type Link struct {
	value *LinkItem
}
