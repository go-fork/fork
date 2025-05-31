package router

import (
	"regexp"
	"strings"
	"sync"
)

// TrieNode đại diện cho một node trong route trie
type TrieNode struct {
	// children lưu trữ các node con, key là segment path
	children map[string]*TrieNode

	// isParam xác định node này có phải là parameter không (:id)
	isParam bool

	// paramName tên parameter nếu đây là param node
	paramName string

	// isWildcard xác định node này có phải là wildcard không (*)
	isWildcard bool

	// isOptional xác định parameter có phải là optional không (?id)
	isOptional bool

	// regexPattern regex constraint cho parameter
	regexPattern string

	// handlers lưu trữ handlers theo HTTP method
	handlers map[string]HandlerFunc

	// isEndNode xác định đây có phải là node cuối của route không
	isEndNode bool

	// mu bảo vệ truy cập đồng thời
	mu sync.RWMutex
}

// RouteTrie structure for optimized route lookup
type RouteTrie struct {
	root *TrieNode
	mu   sync.RWMutex
}

// NewRouteTrie tạo một route trie mới
func NewRouteTrie() *RouteTrie {
	return &RouteTrie{
		root: &TrieNode{
			children: make(map[string]*TrieNode),
			handlers: make(map[string]HandlerFunc),
		},
	}
}

// Insert thêm route vào trie
func (rt *RouteTrie) Insert(method, path string, handler HandlerFunc) {
	rt.mu.Lock()
	defer rt.mu.Unlock()

	segments := rt.splitPath(path)
	current := rt.root

	for _, segment := range segments {
		current.mu.Lock()

		if current.children == nil {
			current.children = make(map[string]*TrieNode)
		}

		// Xử lý các loại segment khác nhau
		key, node := rt.processSegment(segment)

		if existingNode, exists := current.children[key]; exists {
			current.mu.Unlock()
			current = existingNode
		} else {
			current.children[key] = node
			current.mu.Unlock()
			current = node
		}
	}

	// Đánh dấu là end node và set handler
	current.mu.Lock()
	current.isEndNode = true
	if current.handlers == nil {
		current.handlers = make(map[string]HandlerFunc)
	}
	current.handlers[method] = handler
	current.mu.Unlock()
}

// Find tìm handler trong trie
func (rt *RouteTrie) Find(method, path string) HandlerFunc {
	rt.mu.RLock()
	defer rt.mu.RUnlock()

	segments := rt.splitPath(path)
	current := rt.root

	return rt.findRecursive(current, segments, method, 0)
}

// findRecursive tìm kiếm đệ quy trong trie
func (rt *RouteTrie) findRecursive(node *TrieNode, segments []string, method string, index int) HandlerFunc {
	if node == nil {
		return nil
	}

	node.mu.RLock()
	defer node.mu.RUnlock()

	// Nếu đã xử lý hết segments
	if index >= len(segments) {
		if node.isEndNode {
			if handler, exists := node.handlers[method]; exists {
				return handler
			}
		}
		return nil
	}

	currentSegment := segments[index]

	// 1. Tìm exact match trước
	if child, exists := node.children[currentSegment]; exists {
		if result := rt.findRecursive(child, segments, method, index+1); result != nil {
			return result
		}
	}

	// 2. Tìm parameter match
	for _, child := range node.children {
		if child.isParam {
			// Kiểm tra regex constraint nếu có
			if child.regexPattern != "" {
				if regex, err := compileRegex(child.regexPattern); err == nil {
					if !regex.MatchString(currentSegment) {
						continue
					}
				} else {
					continue
				}
			}

			if result := rt.findRecursive(child, segments, method, index+1); result != nil {
				return result
			}

			// Xử lý optional parameter
			if child.isOptional {
				if result := rt.findRecursive(child, segments, method, index); result != nil {
					return result
				}
			}
		}
	}

	// 3. Tìm wildcard match cuối cùng
	for _, child := range node.children {
		if child.isWildcard {
			// Wildcard match với tất cả segments còn lại
			if child.isEndNode {
				if handler, exists := child.handlers[method]; exists {
					return handler
				}
			}
		}
	}

	return nil
}

// processSegment xử lý một segment và trả về key và node tương ứng
func (rt *RouteTrie) processSegment(segment string) (string, *TrieNode) {
	node := &TrieNode{
		children: make(map[string]*TrieNode),
		handlers: make(map[string]HandlerFunc),
	}

	// Static segment
	if !strings.HasPrefix(segment, ":") && !strings.HasPrefix(segment, "*") {
		return segment, node
	}

	// Parameter segment (:id)
	if strings.HasPrefix(segment, ":") {
		paramName := segment[1:]
		key := ":param"

		// Optional parameter (:id?)
		if strings.HasSuffix(paramName, "?") {
			paramName = paramName[:len(paramName)-1]
			node.isOptional = true
			key = ":optional"
		}

		// Regex constraint (:id<\d+>)
		if idx := strings.Index(paramName, "<"); idx >= 0 && strings.HasSuffix(paramName, ">") {
			node.regexPattern = paramName[idx+1 : len(paramName)-1]
			paramName = paramName[:idx]
			key = ":regex:" + node.regexPattern
		}

		node.isParam = true
		node.paramName = paramName
		return key, node
	}

	// Wildcard segment (*filepath)
	if strings.HasPrefix(segment, "*") {
		node.isWildcard = true
		node.paramName = segment[1:]
		return "*", node
	}

	return segment, node
}

// splitPath chia path thành các segments
func (rt *RouteTrie) splitPath(path string) []string {
	if path == "/" || path == "" {
		return []string{}
	}

	// Loại bỏ leading slash
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}

	// Loại bỏ trailing slash
	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}

	if path == "" {
		return []string{}
	}

	return strings.Split(path, "/")
}

// compileRegex compile regex pattern với caching
func compileRegex(pattern string) (*regexp.Regexp, error) {
	regexCacheMu.RLock()
	if regex, found := regexCache[pattern]; found {
		regexCacheMu.RUnlock()
		return regex, nil
	}
	regexCacheMu.RUnlock()

	regex, err := regexp.Compile("^" + pattern + "$")
	if err != nil {
		return nil, err
	}

	regexCacheMu.Lock()
	if existingRegex, found := regexCache[pattern]; found {
		regexCacheMu.Unlock()
		return existingRegex, nil
	}
	regexCache[pattern] = regex
	regexCacheMu.Unlock()

	return regex, nil
}

// Clear clears all nodes and handlers from the trie to prevent memory leaks
func (rt *RouteTrie) Clear() {
	rt.mu.Lock()
	defer rt.mu.Unlock()

	if rt.root != nil {
		rt.clearNode(rt.root)
		rt.root = &TrieNode{
			children: make(map[string]*TrieNode),
			handlers: make(map[string]HandlerFunc),
		}
	}
}

// clearNode recursively clears a node and its children
func (rt *RouteTrie) clearNode(node *TrieNode) {
	if node == nil {
		return
	}

	node.mu.Lock()
	defer node.mu.Unlock()

	// Clear children recursively
	for _, child := range node.children {
		rt.clearNode(child)
	}

	// Clear maps
	node.children = nil
	node.handlers = nil
}

// GetNodeCount returns the total number of nodes in the trie for monitoring
func (rt *RouteTrie) GetNodeCount() int {
	rt.mu.RLock()
	defer rt.mu.RUnlock()

	return rt.countNodes(rt.root)
}

// countNodes recursively counts nodes in the trie
func (rt *RouteTrie) countNodes(node *TrieNode) int {
	if node == nil {
		return 0
	}

	node.mu.RLock()
	defer node.mu.RUnlock()

	count := 1
	for _, child := range node.children {
		count += rt.countNodes(child)
	}
	return count
}
