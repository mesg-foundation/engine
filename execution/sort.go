package execution

// ByBlockHeight implements sort.Interface for []*Execution based on the block height field.
type ByBlockHeight []*Execution

func (a ByBlockHeight) Len() int           { return len(a) }
func (a ByBlockHeight) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByBlockHeight) Less(i, j int) bool { return a[i].BlockHeight < a[j].BlockHeight }
