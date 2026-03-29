# go-server

基于 Gin 的 Go Web API 分层模板，适合从 0 到 1 快速起步。

## Quick Start

```bash
go run ./cmd/api
```

服务默认监听 `0.0.0.0:8080`。

## 数据库迁移

项目使用 `golang-migrate` 管理 MySQL 数据库结构变更，不再依赖运行时自动建表。

### 安装 migrate CLI

```bash
go install -tags "mysql" github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### 执行迁移
```bash
migrate -path ./migrations -database "mysql://root:root@tcp(127.0.0.1:3306)/go_server" up
```

### 回滚一步
```bash
migrate -path ./migrations -database "mysql://root:root@tcp(127.0.0.1:3306)/go_server" down 1
```

### 新建迁移文件
#### 迁移文件命名格式：
```bash
000001_create_users_table.up.sql
000001_create_users_table.down.sql
```
### 说明：
1. up.sql 表示向前升级数据库版本
2. down.sql 表示回滚当前版本
3. 所有数据库结构变更都应通过 migrations/ 目录维护

## 已实现

- 统一配置加载（环境变量）
- 统一响应结构与错误码
- 请求 ID + 日志 + panic 恢复
- `health` 探活接口
- `user` 示例业务（内存仓储 CRUD）

详细设计见：[docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)
