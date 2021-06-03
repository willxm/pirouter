package pirouter

type (
	// Tree records node
	Tree struct {
		root   *Node
		routes map[string]*Node
	}

	// Node records any URL params, and executes an end handler.
	Node struct {
		key string
		// path records a request path
		path   string
		handle interface{}
		// depth records Node's depth
		depth int
		// children records Node's children node
		children map[string]*Node
		// isPattern flag
		isPattern   bool
		middlewares []Middleware
	}
)

// NewNode returns a newly initialized Node object that implements the Node
func newNode(key string, depth int) *Node {
	return &Node{
		key:      key,
		depth:    depth,
		children: make(map[string]*Node),
	}
}

// NewTree returns a newly initialized Tree object that implements the Tree
func NewTree() *Tree {
	return &Tree{
		root:   newNode("/", 1),
		routes: make(map[string]*Node),
	}
}

// Add use `pattern` 、handle、middleware stack as node register to tree
func (t *Tree) Add(pattern string, handle interface{}) {
	var currentNode = t.root

	if pattern != currentNode.key {
		pattern = TrimPathPrefix(pattern)
		res := SplitPattern(pattern)
		for _, key := range res {
			node, ok := currentNode.children[key]
			if !ok {
				node = newNode(key, currentNode.depth+1)
				currentNode.children[key] = node
			}
			currentNode = node
		}
	}

	currentNode.handle = handle
	currentNode.isPattern = true
	currentNode.path = pattern

}

// Find returns nodes that the request match the route pattern
func (t *Tree) Find(pattern string) (nodes *Node, handlers []Middleware) {
	node := t.root

	if pattern == node.path {
		nodes = node
		return
	}

	pattern = TrimPathPrefix(pattern)

	res := SplitPattern(pattern)

	for _, key := range res {
		child, ok := node.children[key]
		if !ok {
			return nil, nil
		}
		handlers = append(handlers, child.middlewares...)
		if pattern == child.path {
			return child, handlers
		}
		node = child
	}
	return
}
