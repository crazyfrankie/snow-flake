package snowflake

import (
	"testing"
)

// 在一个短时间内（例如几百万次调用）生成 ID，检查是否有重复。
// 使用哈希表存储所有生成的 ID，如果有任何重复值，测试失败。
// TestUniqueID -- PASS
func TestUniqueID(t *testing.T) {
	node, err := NewNode(1)
	if err != nil {
		t.Fatalf("Error creating node: %v", err)
	}
	ids := make(map[int64]struct{})
	for i := 0; i < 1000000; i++ {
		id := node.GenerateCode()
		if _, exists := ids[id]; exists {
			t.Errorf("Duplicate ID found: %d", id)
		}
		ids[id] = struct{}{}
	}
}
