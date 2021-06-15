package list

type List interface {
	Contains(value int) bool
	Insert(value int) bool
	Delete(value int) bool
	Range(f func(value int) bool)
	Len() int
}
