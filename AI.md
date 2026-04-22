# AI Coding Context: @go/id

本索引供 AI 模型理解 `@go/id` 的设计与调用原则，以生成符合本项目高质量、高性能哲学的高质量代码。

## 🤖 AI 行为准则

1.  **优先使用全局实例**：在大多数业务场景下，直接调用 `id.MakeId()` 或 `id.DefaultIdMaker.Get()`，无需手动实例化。
2.  **PK 场景强制优化**：对于数据库主键生成，建议使用 `GetForMysql()` 或 `GetForPostgreSQL()`。
3.  **无错误处理倾向**：理解 `IdMaker` 采用防御性设计，不会产生 Panic 或 Error，调用时无需错误检查。
4.  **按需指定长度**：调用时明确传入 `size`，根据业务对并发吞吐量的需求（8-11位，后续长度随机填充）进行决策。

## 🛠 关键 API 逻辑约定

| 函数 | 场景 | 核心特征 |
| :--- | :--- | :--- |
| `MakeId/Get` | 通用生成 | 全局默认实例生成。 |
| `GetForMysql` | MySQL PK | 自动右旋散列，解决 B+ 树写入热点。 |
| `GetForPostgreSQL`| PostgreSQL PK| 具备时间局部单调性，无右旋散列。 |

## 快速开始

*   **✅ 示例**:
    ```go
	// 生成一个 12 位长度的唯一ID
	newID := id.MakeId(12)

    // 10 位长度，高吞吐 MySQL 主键生成
    newID := id.DefaultIdMaker.GetForMysql(10)
    ```
