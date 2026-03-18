# go-server

基于 Gin 的 Go Web API 分层模板，适合从 0 到 1 快速起步。

## Quick Start

```bash
go run ./cmd/api
```

服务默认监听 `0.0.0.0:8080`。

## 已实现

- 统一配置加载（环境变量）
- 统一响应结构与错误码
- 请求 ID + 日志 + panic 恢复
- `health` 探活接口
- `user` 示例业务（内存仓储 CRUD）

详细设计见：[docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)
