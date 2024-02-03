package types

type LinkedList struct {
	head *nodePoint
	tail *nodePoint
	Len  int
}

func (l *LinkedList) Append(val any) {
	//TODO implement me
	panic("implement me")
}

func (l *LinkedList) Delete(val any) {
	//TODO implement me
	panic("implement me")
}

func (l LinkedList) Add(index int, val any) {

}

func (l *LinkedList) AddV1(index int, val any) {

}

type nodePoint struct {
}
