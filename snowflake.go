package snowflake

import (
	"errors"
	"sync"
	"time"
)

var (
	epoch          int64 = 1732971921601       // 开始运行时间,不可更改
	nodeBits       uint8 = 10                  // 节点 ID 的位数
	seqBits        uint8 = 12                  // 1毫秒内可生成的 id 序号的二进制位数
	nodeMax        int64 = (1 << nodeBits) - 1 // 节点 ID 的最大值，用于防止溢出
	seqMax         int64 = (1 << seqBits) - 1  // 用来表示生成 id 序号的最大值
	timeShift            = nodeBits + seqBits  // 时间戳向左的偏移量
	nodeShift            = seqBits             // 节点 ID 向左的偏移量
	GetCurrentTime       = func() int64 {
		return time.Now().UnixMilli()
	}
)

type Node struct {
	// 并发安全
	mu sync.Mutex
	// 时间戳
	timestamp int64
	// 该节点的 ID
	NodeId int64
	// 当前毫秒已经生成的 id 序列号(从0开始累加) 1毫秒内最多生成4096个 ID
	seq int64
}

func NewNode(NodeId int64) (*Node, error) {
	// 检测 NodeId 是否在定义的范围内
	if NodeId < 0 || NodeId > nodeMax {
		return nil, errors.New("node ID excess of quantity")
	}

	return &Node{
		timestamp: 0,
		NodeId:    NodeId,
		seq:       0,
	}, nil
}

func (n *Node) GenerateCode() int64 {
	// 加锁
	n.mu.Lock()
	defer n.mu.Unlock()

	// 获取生成时的时间戳
	now := GetCurrentTime()
	if n.timestamp == now {
		// 对当前已生成的 ID 数递增
		// 采用的是位运算，确保在 seqMax 的范围内循环
		// 且位运算优于取余
		n.seq = (n.seq + 1) & seqMax

		// 判断当前工作节点是否在1毫秒内已经生成 seqMax 个 ID
		if n.seq == 0 {
			// 如果当前工作节点在1毫秒内生成的 ID 已经超过上限 需要等待1毫秒再继续生成
			for now <= n.timestamp {
				now = GetCurrentTime()
			}
		}
	} else {
		// 如果当前时间与工作节点上一次生成 ID 的时间不一致 则需要重置工作节点生成ID的序号
		n.seq = 0
	}

	// 将机器上一次生成 ID 的时间更新为当前时间
	n.timestamp = now

	// 第一段 now - epoch 为该算法目前已经奔跑了xxx毫秒
	// 如果在程序跑了一段时间修改了epoch这个值 可能会导致生成相同的ID
	ID := (now-epoch)<<timeShift | (n.NodeId << nodeShift) | (n.seq)

	return ID
}
