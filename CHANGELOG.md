# Changelog: @go/id

## v1.5.8 (2026-09-05)
- **依赖对齐**: 升级 `cast v1.5.6`、`encoding v1.5.7`、`jsmod v1.5.4`、`rand v1.5.4`。

## v1.5.7 (2026-07-18)
- **依赖校验修复**: 重新整理模块图并补齐 checksum。

## v1.5.6 (2026-06-22)
- **接口精简与测试逆向对齐 (API Simplification & Decryption Alignment)**:
  - 弃用并删除了 `FillRand` 的调用，重构为直接调用 `encoding.FillInt`。
  - 同步修改 `UnhashInt` 单元测试以匹配 `HashInt` 的 2 轮非线性传播混淆，并在 `TestReversibility` 中完美验证了变换的双射性与可逆性。

## v1.5.5 (2026-06-21)
- **碰撞安全与重排修复 (Entropy Fold & Order Fix)**:
  - 引入模同余折叠：在非有序模式截断前调用 `intEncoder.FoldInt` 将超长字节以 Base62 同余模加异或折叠，确保高并发下的唯一性熵完全保留在截断范围内，消除了物理截断的局部高碰撞 Bug。
  - 修复 `ExchangeInt` 就地修改对齐：随着底层 `ExchangeInt` 修复为原地修改，重新激活了有序数据库优化模式（`GetForMysql`/`GetForPostgreSQL`）下对后半段序列号混淆的能力，使得首位分片散列（`hashToHead`）成功防范 B+ 树写入热点。
  - 自动化回归测试：在 `id_test.go` 中引入了 `TestReversibility`（双射可逆测试）、`TestCollisionRate`（并发碰撞测试）与 `TestCheckHashDiff`（差分诊断），通过严格的回归用例确保后续算法更新安全。

## v1.5.4 (2026-06-21)
- **重构与错误堆栈支持**:
  - 重构 `js_export.go`，将 `Make` 匿名包装闭包改为包级具名函数 `jsMake`。
  - 升级 `encoding` 依赖至 v1.5.4，`jsmod` 依赖至 v1.5.3，`rand` 依赖至 v1.5.3，`cast` 依赖至 v1.5.3。

## v1.5.3 (2026-06-11)
- **版本对齐**: 基础设施全局对齐 v1.5.3。

## v1.5.2 (2026-06-09)
- **版本对齐**: 基础设施全局对齐 v1.5.2。

## v1.5.1 (2026-06-08)
- **JS 极简 API**: 彻底废除 JS 侧冗长难记的 `get8Bytes...` 等方法，统一收敛为极简的 `Make(size int, forDB ...string)`。充分利用了底层 JS Bridge 的零值回退特性，在不填第二个参数时也能安全优雅地工作。

## [v1.3.2] - 2026-05-30

### Added
- **JSMOD 注册**：将语义化 ID 生成快捷函数注册到 `jsmod`。

## [v1.0.1] - 2026-05-01

### Changed
- 重构：将 IdMaker 重命名为 IDMaker，符合 Go 命名规范。
- 优化：使用 Epoch 常量替代硬编码时间戳，优化 defaultIncr 控制流。

## [v1.0.0] - 2026-04-22

### Added
- 分布式唯一 ID 生成器，支持多种数据库场景优化。
- 高性能并发支持，集成 `@go/rand` 与 `@go/encoding`。
