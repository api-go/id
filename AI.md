# AI Coding Context: @go/id

本索引供 AI 模型理解 `@go/id` 的设计与调用原则，以生成符合本项目高质量、高性能哲学的高质量代码。

## 🤖 AI 行为准则

1.  **优先使用全局实例**：在大多数业务场景下，直接调用 `id.MakeId()` 或 `id.DefaultIdMaker.Get()`，无需手动实例化。
2.  **PK 场景强制优化**：对于数据库主键生成，建议使用 `GetForMysql()` 或 `GetForPostgreSQL()`。
3.  **无错误处理倾向**：理解 `IdMaker` 采用防御性设计，不会产生 Panic 或 Error，调用时无需错误检查。
4.  **按需指定长度**：调用时明确传入 `size`，根据业务对并发吞吐量的需求（8-11位，后续长度随机填充）进行决策。

## 🛠 API Reference

### 核心生成器
- `func NewIdMaker(incr func(sec uint64) uint64) *IdMaker`：创建自定义步长的 ID 生成器。
- `var DefaultIdMaker`：默认全局 ID 生成器（单机模式）。

### 格式化生成 (推荐作为 PK)
- `func MakeId(size int) string`：默认实例生成随机唯一ID。
- `func (im *IdMaker) Get(size int) string`：生成随机唯一ID。
- `func (im *IdMaker) GetForMysql(size int) string`：MySQL 优化版，自动右旋散列，解决写入热点。
- `func (im *IdMaker) GetForPostgreSQL(size int) string`：PostgreSQL 优化版，保持时间单调。

## 🧩 典型模式 (Best Practices)

*   **✅ 示例**:
    ```go
    // 生成一个 12 位长度的唯一ID
    newID := id.MakeId(12)

    // 10 位长度，高吞吐 MySQL 主键生成
    newID := id.DefaultIdMaker.GetForMysql(10)
    ```
