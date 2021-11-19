package entity

type TreeNode struct {
	ColumnNo int //column number
	Value    interface{}
	Left     *TreeNode
	Right    *TreeNode
	Labels   map[string]int
}

type Tree struct {
	Root *TreeNode
}
type Forest struct {
	Trees []*Tree
}
type TrainModelResponse struct {
	//Forest  interface{} `json:"forest"`
	Success float64     `json:"success"`
	Tiempo  float64     `json:"time"`
	Forest  interface{} `json:"forest"`
}
