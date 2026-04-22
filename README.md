# go/id (@go/id)

## 关于本项目
本项目完全由 AI 维护。代码源自 github.com/ssgo/u 的重构。

这是一个高性能、纯 Go 实现的分布式唯一ID生成器，具备对标雪花算法的特性，同时针对数据库主键场景进行了深度优化。

## 设计哲学

*   **数据库友好**：通过右旋散列算法，完美避开 B+ 树索引的写入热点，同时保持时间局部单调性。
*   **极致紧凑**：支持 8 位到 20 位 Base62 编码，根据业务负载按需配置。
*   **防御性设计**：无锁并发生成，自动处理秒级重置与碰撞防御，零 Panic。
*   **语义直觉**：通过 `secTag` 机制，在同一套编码格式下平滑过渡 300 年的时间跨度（2314 年后自动激活6位扩展模式最多可支撑1800年时间无碰撞生成）。

## 性能与规模
ID 结构由 `secTag` (1位) + `sec` (5-6位) + `secIndex` (1-5位) + `padding` (随机填充) 组成。

| 输出长度 | secIndex 长度 | 并发支持能力 (每秒) |
| :--- | :--- | :--- |
| 8位 | 2位 | 3,844 个 |
| 9位 | 3位 | 9 万个 |
| 10位 | 4位 | 1,477 万个 |
| 11位 | 6位 | 9 亿个 |

## API Reference

### 核心生成器
- `func NewIdMaker(incr func(sec uint64) uint64) *IdMaker`：创建自定义步长的 ID 生成器。
- `var DefaultIdMaker`：默认全局 ID 生成器（单机模式）。

### 格式化生成 (推荐作为 PK)
- `func MakeId(size int) string`：DefaultIdMaker.Get(size) 的别名。
- `func (im *IdMaker) Get(size int) string`：生成随机唯一ID。
- `func (im *IdMaker) GetForMysql(size int) string`：MySQL 优化版，自动右旋散列，解决写入热点。
- `func (im *IdMaker) GetForPostgreSQL(size int) string`：PostgreSQL 优化版，保持时间局部单调。

## 快速开始

```go
import "apigo.cc/go/id"

// 生成一个 12 位长度的唯一ID
newID := id.MakeId(12)

// 生成一个 10 位长度的 MySQL 友好主键
newID := id.DefaultIdMaker.GetForMysql(10)
```

## 分布式集群扩展示例
通过实现 `Incr` 钩子集成 Redis (使用 `github.com/gomodule/redigo`)：

```go
import (
    "fmt"
    "sync"
    "apigo.cc/go/id"
    "github.com/gomodule/redigo/redis"
)

type RedisIdMaker struct {
    rd          redis.Conn
    secCurrent  uint64
    secIndexMax uint64
    secIndexNext uint64
    lock        sync.Mutex
    maker       *id.IdMaker
}

func NewRedisIdMaker(rd redis.Conn) *id.IdMaker {
    rim := &RedisIdMaker{rd: rd}
    return id.NewIdMaker(rim.makeSecIndex)
}

func (rim *RedisIdMaker) makeSecIndex(sec uint64) uint64 {
    rim.lock.Lock()
    defer rim.lock.Unlock()
    
    if rim.secCurrent == sec && rim.secIndexNext <= rim.secIndexMax {
        idx := rim.secIndexNext
        rim.secIndexNext++
        return idx
    }
    
    rim.secCurrent = sec
    key := fmt.Sprintf("_SecIdx_%d", sec)
    // 每次从 Redis 预取 100 个序列号
    max, _ := redis.Uint64(rim.rd.Do("INCRBY", key, 100))
    rim.secIndexMax = max
    rim.secIndexNext = max - 99
    if max <= 100 { rim.rd.Do("EXPIRE", key, 10) }
    return rim.secIndexNext
}
```
