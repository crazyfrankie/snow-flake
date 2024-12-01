package snowflake

import (
	"sync"
	"testing"
	"time"
)

// 生成 ID
func (n *Node) Generatecode() int64 {
	n.mu.Lock()
	defer n.mu.Unlock()

	now := time.Now().UnixMilli()
	if n.timestamp == now {
		// 对当前已生成的 ID 数递增
		n.seq = (n.seq + 1) & seqMax

		// 判断当前工作节点是否在1毫秒内已经生成 seqMax 个 ID
		if n.seq == 0 {
			// 使用忙等待
			for now <= n.timestamp {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		// 如果当前时间与工作节点上一次生成 ID 的时间不一致 则需要重置工作节点生成ID的序号
		n.seq = 0
	}

	n.timestamp = now
	ID := (now-epoch)<<timeShift | (n.NodeId << nodeShift) | (n.seq)

	return ID
}

// 生成 ID 使用 time.NewTimer
func (n *Node) GenerateCodeWithTimer() int64 {
	n.mu.Lock()
	defer n.mu.Unlock()

	now := time.Now().UnixMilli()
	if n.timestamp == now {
		// 对当前已生成的 ID 数递增
		n.seq = (n.seq + 1) & seqMax

		// 判断当前工作节点是否在1毫秒内已经生成 seqMax 个 ID
		if n.seq == 0 {
			// 使用 time.NewTimer
			timer := time.NewTimer(time.Millisecond)
			<-timer.C // 等待 1 毫秒
		}
	} else {
		n.seq = 0
	}

	n.timestamp = now
	ID := (now-epoch)<<timeShift | (n.NodeId << nodeShift) | (n.seq)

	return ID
}

// 测试生成 ID 的性能，比较忙等待和 time.NewTimer
func BenchmarkGenerateCode(b *testing.B) {
	node, err := NewNode(1)
	if err != nil {
		b.Fatalf("Failed to create node: %v", err)
	}

	// 设置并发的数量
	concurrency := 100
	var wg sync.WaitGroup

	b.Run("BenchmarkBusyWait", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// 启动多个 goroutine 来模拟并发请求
			for j := 0; j < concurrency; j++ {
				wg.Add(1) // 每启动一个新的 goroutine 增加计数
				go func() {
					defer wg.Done()
					node.GenerateCode() // 使用忙等待生成 ID
				}()
			}
		}
		wg.Wait()
	})

	b.Run("BenchmarkTimeNewTimer", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// 启动多个 goroutine 来模拟并发请求
			for j := 0; j < concurrency; j++ {
				wg.Add(1) // 每启动一个新的 goroutine 增加计数
				go func() {
					defer wg.Done()
					node.GenerateCodeWithTimer() // 使用 time.NewTimer 生成 ID
				}()
			}
		}
		wg.Wait()
	})
}
