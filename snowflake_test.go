package snowflake

import (
	"testing"
)

type MockTimeProvider struct {
	mockTime int64
}

func (mtp *MockTimeProvider) GetCurrentTime() int64 {
	return mtp.mockTime
}

func TestGenerateCode(t *testing.T) {
	node, err := NewNode(1)
	if err != nil {
		t.Fatalf("Failed to create node: %v", err)
	}

	// 使用 MockTimeProvider 来模拟时间
	mockTimeProvider := &MockTimeProvider{}
	// 替换全局的时间提供函数
	GetCurrentTime = mockTimeProvider.GetCurrentTime

	testCases := []struct {
		name       string
		mockTime   int64
		expectedID func(now int64, nodeId int64, seq int64) int64
	}{
		{
			name:     "正常生成",
			mockTime: epoch + 1000,
			expectedID: func(now int64, nodeId int64, seq int64) int64 {
				return (now-epoch)<<timeShift | (nodeId << nodeShift) | seq
			},
		},
		{
			name:     "跨时间生成 ID",
			mockTime: epoch + 2000,
			expectedID: func(now int64, nodeId int64, seq int64) int64 {
				return (now-epoch)<<timeShift | (nodeId << nodeShift) | seq
			},
		},
		{
			name:     "长时间生成能力",
			mockTime: epoch + 5000,
			expectedID: func(now int64, nodeId int64, seq int64) int64 {
				return (now-epoch)<<timeShift | (nodeId << nodeShift) | seq
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 设置模拟时间
			mockTimeProvider.mockTime = tc.mockTime

			// 重载序列号
			node.seq = 0

			// 调用 GenerateCode 生成 ID
			id := node.GenerateCode()

			// 获取期望的 ID
			expected := tc.expectedID(tc.mockTime, node.NodeId, node.seq)

			// 比较生成的 ID 和期望的 ID
			if id != expected {
				t.Errorf("For test case '%s', expected ID: %d, got ID: %d", tc.name, expected, id)
			}
		})
	}

}
