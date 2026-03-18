# Go Web API 工程架构说明（Gin）

本文档对应当前项目代码实现，目标是让你从 0 到 1 快速建立一套可持续迭代的 Go API 服务架构。

## 1. 设计目标

- 分层清晰：路由、控制器、业务、存储职责分离。
- 低耦合：`service` 依赖 `repository` 接口，便于后续替换为 Mongo/MySQL。
- 可观测：统一 `request_id`、日志、错误响应。
- 先跑通再扩展：先用内存仓储实现 CRUD，业务和接口稳定后再接数据库。

## 2. 当前目录结构

```text
go-server/
├─ cmd/
│  └─ api/
│     └─ main.go                 # 标准启动入口
├─ internal/
│  ├─ bootstrap/
│  │  └─ app.go                  # 应用装配、HTTP Server 启动
│  ├─ config/
│  │  └─ config.go               # 环境变量配置读取
│  ├─ middleware/
│  │  ├─ request_id.go           # 请求 ID 注入
│  │  ├─ logger.go               # 访问日志
│  │  ├─ recovery.go             # panic 恢复
│  │  └─ error_handler.go        # 统一错误出参
│  ├─ response/
│  │  └─ response.go             # 统一响应结构 + 业务错误类型
│  ├─ router/
│  │  └─ router.go               # 路由注册与模块装配
│  └─ modules/
│     ├─ health/
│     │  └─ handler.go           # 健康检查模块
│     └─ user/
│        ├─ dto.go               # 入参/出参 DTO
│        ├─ model.go             # 领域模型
│        ├─ repository.go        # 仓储接口 + 内存实现
│        ├─ service.go           # 业务逻辑
│        └─ handler.go           # HTTP 控制器
└─ docs/
   └─ ARCHITECTURE.md            # 本文档
```

## 3. 分层职责

### 3.1 `cmd`（入口层）

- 只做启动，不写业务逻辑。
- 调用 `bootstrap.Run()`。

### 3.2 `bootstrap`（应用装配层）

- 读取配置。
- 创建 gin engine。
- 创建 `http.Server` 并启动。

### 3.3 `router`（路由层）

- 注册全局中间件。
- 以模块为单位注册路由。
- 完成依赖注入：`repo -> service -> handler`。

### 3.4 `modules/*/handler`（接口层）

- 解析参数与校验（`ShouldBindJSON` + binding tag）。
- 调用 service。
- 负责 HTTP 状态码与响应格式。

### 3.5 `modules/*/service`（业务层）

- 承载业务规则。
- 定义错误语义（冲突、未找到等）。
- 不依赖 Gin，上下文只使用 `context.Context`。

### 3.6 `modules/*/repository`（数据访问层）

- 定义仓储接口。
- 当前实现为内存版本，便于 0->1 验证接口。
- 后续可平滑替换为 Mongo/MySQL 实现。

### 3.7 `response`（统一响应层）

- 成功响应：`code/message/data/request_id`。
- 业务错误：`AppError{status, code, message}`。
- 降低 handler/service 重复代码。

### 3.8 `middleware`（横切关注点）

- `RequestID`：生成并透传 `X-Request-ID`。
- `Logger`：记录状态码、路径、耗时、IP。
- `Recovery`：兜底 panic。
- `ErrorHandler`：统一将 error 转成标准 JSON。

## 4. 请求生命周期（链路）

1. 请求进入 Gin。
2. `RequestID` 生成请求追踪 ID。
3. `Logger` 记录开始时间。
4. 路由进入对应模块 handler。
5. handler 解析参数并调用 service。
6. service 执行业务并调用 repository。
7. 返回结果后由 response 输出统一 JSON。
8. 若出现错误，交给 `ErrorHandler` 统一处理。
9. `Logger` 输出最终日志。

## 5. 已提供 API（示例）

前缀：`/api/v1`

- `GET /health`
  - 服务探活。
- `POST /users`
  - 创建用户。
  - body: `{ "name": "Tom", "email": "tom@example.com" }`
- `GET /users`
  - 用户列表。
- `GET /users/:id`
  - 获取单个用户。
- `DELETE /users/:id`
  - 删除用户。

## 6. 响应格式规范

成功示例：

```json
{
  "code": "OK",
  "message": "success",
  "data": {
    "id": "..."
  },
  "request_id": "20260318120000-xxxx"
}
```

失败示例：

```json
{
  "code": "INVALID_ARGUMENT",
  "message": "invalid request body",
  "request_id": "20260318120000-xxxx"
}
```

## 7. 配置项（环境变量）

- `APP_NAME`，默认 `go-server`
- `APP_ENV`，默认 `local`
- `HTTP_HOST`，默认 `0.0.0.0`
- `HTTP_PORT`，默认 `8080`
- `HTTP_READ_TIMEOUT`，默认 `5s`
- `HTTP_WRITE_TIMEOUT`，默认 `10s`

## 8. 启动方式

```bash
go run ./cmd/api
```

## 9. 从 0 到 1 的迭代建议

### 阶段 A：接口与业务稳定

- 用内存仓储把核心 API 套路跑通。
- 先确定 DTO、错误码、响应结构。

### 阶段 B：接入数据库

- 新建 `repository_mongo.go` 或 `repository_mysql.go`。
- 保持 `Repository` 接口不变。
- 在 `router`（或 bootstrap）里切换具体实现。

### 阶段 C：工程化能力补全

- 配置中心（Viper / envconfig）。
- 结构化日志（zap / slog）。
- 鉴权（JWT + RBAC）。
- 参数校验国际化。
- OpenAPI/Swagger 文档自动生成。
- 单元测试 + 集成测试 + CI。

## 10. 为什么这套结构适合你当前阶段

- 学习曲线平稳：每层代码量可控，易读。
- 扩展成本低：新增模块按同样模板复制即可。
- 替换技术栈风险低：数据库和中间件可插拔。

如果你下一步要接 Mongo，我可以在这个架构上继续帮你把 `user repository` 切换为 `mongo-driver/v2` 版本，并补上索引初始化与连接管理。
