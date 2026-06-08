# Changelog: @go/id

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
