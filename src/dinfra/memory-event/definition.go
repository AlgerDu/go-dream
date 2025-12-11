package memoryevent

import "github.com/AlgerDu/go-dream/src/dinfra"

type (
	HandlerRecord struct {
		ID      string              // 每个订阅都需要分配一个 ID，方便日志查找及取消订阅
		Handler dinfra.EventHandler // 事件的处理函数
	}

	// 做成一个简单的树结构，可以处理 + # 等特殊匹配
	Node struct {
		Key      string                    // 节点的订阅字符
		Topic    string                    // 对应的主题
		Handlers map[string]*HandlerRecord // 对应的 handler 记录
		Children map[string]*Node          // 子项
	}
)

func BuildeNode(key string) *Node {
	return &Node{
		Key:      key,
		Handlers: map[string]*HandlerRecord{},
		Children: map[string]*Node{},
	}
}

func FindChild(node *Node, keys []string) (*Node, []string) {
	if len(keys) == 0 {
		return node, []string{}
	}

	key := keys[0]
	child, exist := node.Children[key]
	if !exist {
		return node, keys
	}

	return FindChild(child, keys[1:])
}

func RemoveRecord(node *Node, subscribeID string) *HandlerRecord {
	record, exist := node.Handlers[subscribeID]
	if exist {
		delete(node.Handlers, subscribeID)
		return record
	}

	for _, child := range node.Children {
		record = RemoveRecord(child, subscribeID)
		if record != nil {
			return record
		}
	}

	return nil
}

func CopyHandlers(node *Node, records []*HandlerRecord) []*HandlerRecord {
	for _, handler := range node.Handlers {
		records = append(records, handler)
	}

	return records
}

func MacthHandlers(node *Node, keys []string) []*HandlerRecord {
	records := []*HandlerRecord{}
	if len(keys) == 0 {
		return CopyHandlers(node, records)
	}

	if node, exist := node.Children["#"]; exist {
		records = CopyHandlers(node, records)
	}

	if node, exist := node.Children["+"]; exist {
		records = append(records, MacthHandlers(node, keys[1:])...)
	}

	key := keys[0]
	if node, exist := node.Children[key]; exist {
		records = append(records, MacthHandlers(node, keys[1:])...)
	}

	return records
}
