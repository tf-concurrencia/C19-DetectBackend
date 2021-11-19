package store

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
type ForestModel struct {
	Forest interface{} `json:"forest"`
}

var Model *ForestModel
