## 雪花算法-SnowFlake

Snowflake，雪花算法是由 Twitter 开源的分布式ID生成算法，以划分命名空间的方式将 64-bit 位分割成多个部分，每个部分代表不同的含义。而 Go 中 64bit 的整数是 int64 类型，所以在 Go 中 SnowFlake 算法生成的 ID 就是 int64 来存储的。
- 第1位占用1bit，其值始终是0，可看做是符号位不使用。
- 第2位开始的41位是时间戳，41-bit位可表示2^41个数，每个数代表毫秒，那么雪花算法可用的时间年限是(1L<<41)/(1000L360024*365)=69 年的时间。
- 中间的10-bit位可表示机器数，即2^10 = 1024台机器，但是一般情况下我们不会部署这么台机器。如果我们对IDC（互联网数据中心）有需求，还可以将 10-bit 分 5-bit 给 IDC，分5-bit给工作机器。这样就可以表示32个IDC，每个IDC下可以有32台机器，具体的划分可以根据自身需求定义。
- 最后12-bit位是自增序列，可表示2^12 = 4096个数。

![image](https://github.com/user-attachments/assets/b7d21139-09f8-43eb-9322-376d36c806ad)

### 基本数据结构
```
var (
	epoch     int64 = 1732971921601       // 开始运行时间,不可更改
	nodeBits  uint8 = 10                  // 节点 ID 的位数
	seqBits   uint8 = 12                  // 1毫秒内可生成的 id 序号的二进制位数
	nodeMax   int64 = (1 << nodeBits) - 1 // 节点 ID 的最大值，用于防止溢出
	seqMax    int64 = (1 << seqBits) - 1  // 用来表示生成 id 序号的最大值
	timeShift       = nodeBits + seqBits  // 时间戳向左的偏移量
	nodeShift       = seqBits             // 节点 ID 向左的偏移量
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
```
### 使用示例
```
node, err := NewNode(1)
if err != nil {
    panic(err)
}
id := node.GenerateCode()
```
