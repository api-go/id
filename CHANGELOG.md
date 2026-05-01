# Changelog: @go/id

## [v1.0.1] - 2026-05-01

### Changed
- 重构：将 IdMaker 重命名为 IDMaker，符合 Go 命名规范。
- 优化：使用 Epoch 常量替代硬编码时间戳，优化 defaultIncr 控制流。

## [v1.0.0] - 2026-04-22

### Added
- 分布式唯一 ID 生成器，支持多种数据库场景优化。
- 高性能并发支持，集成 `@go/rand` 与 `@go/encoding`。
