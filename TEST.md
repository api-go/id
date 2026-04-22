# Test Report: @go/id

## 📋 测试概览
- **测试时间**: 2026-04-22
- **测试环境**: darwin/amd64 (Intel i9-9980HK)
- **Go 版本**: 1.25.0

## ✅ 功能测试 (Functional Tests)
| 场景 | 状态 | 描述 |
| :--- | :--- | :--- |
| `TestMakeId` | PASS | 生成长度符合预期的通用 ID。 |
| `TestGetForMysql` | PASS | 生成 MySQL 友好主键，包含右旋散列逻辑。 |
| `TestGetForPostgreSQL`| PASS | 生成 PostgreSQL 友好主键，无右旋散列。 |

## ⚡ 性能基准 (Benchmarks)
| 函数 | 平均耗时 | 内存分配 |
| :--- | :--- | :--- |
| `MakeId-10` | **1523 ns/op** | 912 B/op (9 allocs/op) |
| `GetForMysql-10` | **1505 ns/op** | 896 B/op (9 allocs/op) |

*   集成 `@go/rand` 的 `FastInt` 优化，在高并发下彻底规避了锁竞争。
*   ID 生成性能稳定，满足高性能主键生成场景。

## 🛡️ 鲁棒性防御 (Robustness)
- **并发安全**：核心计数器使用 `sync.Mutex` 保护，生成器随机数部分使用高性能无锁池。
- **碰撞防御**：秒级重置机制配合随机偏移初始化，极大地降低碰撞概率。
- **扩展性**：支持自定义钩子，轻松对接 Redis 等分布式协调服务。
