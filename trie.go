package pirouter

import (
	"fmt"
	"reflect"
	"runtime"
)

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
		handle []HandlerFunc
		// depth records Node's depth
		depth int
		// children records Node's children node
		children map[string]*Node
		// isPattern flag
		isPattern bool
		// wildcard flag
		wildcard bool
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
func (t *Tree) Add(pattern string, handle ...HandlerFunc) {
	var currentNode = t.root

	if pattern != currentNode.key {
		pattern = TrimPathPrefix(pattern)
		res := SplitPattern(pattern)
		for _, key := range res {
			var wildcard bool
			if len(key) > 0 && key[0] == ':' {
				wildcard = true
				for _, node := range currentNode.children {
					if node.wildcard && node.key != key {
						panic(fmt.Sprintf("ambiguous route! can not register '%s', wildcard '%s' already exsit", pattern, node.path))
					}
				}
			}
			node, ok := currentNode.children[key]
			if !ok {
				node = newNode(key, currentNode.depth+1)
				node.wildcard = wildcard
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
func (t *Tree) Find(pattern string) (nodes []*Node) {
	var (
		node  = t.root
		queue []*Node
	)

	if pattern == node.path {
		nodes = append(nodes, node)
		return
	}

	pattern = TrimPathPrefix(pattern)

	res := SplitPattern(pattern)

	for i, key := range res {
		end := i == len(res)-1
		child, ok := node.children[key]

		if !ok {
			for _, v := range node.children {
				if v.wildcard {
					child = v
					break
				}
			}
			if child == nil {
				return
			}
		}

		if (key == child.key || child.wildcard) && child.isPattern && end {
			nodes = append(nodes, child)
			return
		}

		node = child
	}

	queue = append(queue, node)

	for len(queue) > 0 {
		var queueTemp []*Node
		for _, n := range queue {
			if n.isPattern {
				nodes = append(nodes, n)
			}

			for _, childNode := range n.children {
				queueTemp = append(queueTemp, childNode)
			}
		}

		queue = queueTemp
	}

	return
}

func (t *Tree) String() {
	if t.root != nil {
		t.root.String()
	}
}

func (n *Node) String() {
	// TODO beauty
	if n.isPattern {
		var funName string
		for _, handler := range n.handle {
			fn := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
			funName += " " + fn
		}
		fmt.Printf("\tKEY:%s\tPATH:%s\tHANDLER:%s\n", n.key, n.path, funName)
	}
	for _, child := range n.children {
		child.String()
	}
}
