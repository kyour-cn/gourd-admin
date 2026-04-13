# AGENTS.md

## 项目概述

GourdAdmin 是基于 Go + Vue 3 的后台管理系统。前后端分离，各自独立运行。

- **后端**: `server/` — Go，模块名 `app`（不是 `gourd-admin`），使用 [go-gourd](https://github.com/go-gourd/gourd) 事件驱动框架 + chi 路由 + GORM/MySQL
- **前端**: `admin/` — Vue 3 + Vite + Element Plus + Pinia，pnpm 管理

## 任务执行规则

- 回复、注释、文档尽量使用中文
- 禁止主动编写任务文档文件、任何示例、测试代码和执行测试命令，除非用户要求
- 输出任务总结尽量简洁，重要对话中不应该包含无关紧要的信息
- 实现新的功能时尽可能参考现有代码，保证代码风格一致且保证代码复用性。
- 每次任务执行完成如果修改了代码，需要总结一个git message文本，只需要一句概要描述，如 fix(order): 修复订单查询异常。
- 优化和重构时不要将原有的注释精简，保持注释的可读性，并且新增代码中复杂的逻辑需要添加注释说明。

## 常用命令

```bash
# 初始化（安装依赖）
task init

# 开发（前后端并行）
task dev

# 仅后端开发
task dev-server

# 仅前端开发
task dev-admin

# 构建（根据当前OS自动选择）
task build
task build-windows
task build-linux
```

需要先安装 [task](https://taskfile.dev)（go-task）。

后端也可直接运行：
```bash
cd server && go run ./cmd/app/main.go
```

前端也可直接运行：
```bash
cd admin && pnpm run dev
```

## 关键架构信息

### Go 模块名是 `app`

所有 import 路径使用 `app/internal/...`，不是 `gourd-admin`。

### GORM gen 生成的代码不可手动编辑

`server/internal/orm/query/*.gen.go` 和 `server/internal/orm/model/*.gen.go` 由 gorm/gen 自动生成。修改数据库表结构后，需要重新生成：

```bash
cd server && go run ./cmd/gorm/main.go
```

此命令需要可连接的 MySQL 数据库。

### 数据库初始化

```bash
cd server && go run ./cmd/install/main.go
```

依次执行 `assets/migrations/schema.sql` 和 `assets/migrations/seed.sql`。

### 服务端入口和启动流程

- 入口: `server/cmd/app/main.go`
- 事件驱动初始化: `server/internal/event/app.go`
  - `app.boot` → 日志 + 数据库 + 缓存
  - `app.init` → 命令行参数解析
  - `app.start` → 定时任务 + HTTP 服务 + 异步任务

### 配置文件

TOML 格式，位于 `server/configs/`：

| 文件 | 用途 |
|------|------|
| `database.toml` | 数据库连接（目前仅支持 MySQL） |
| `http.toml` | HTTP 端口(默认8080)、静态资源路径(默认`./web`) |
| `jwt.toml` | JWT 认证配置 |
| `log.toml` | 日志配置 |

### API 路由

所有后台 API 挂载在 `/admin/` 下，见 `server/internal/http/admin/router/router.go`。路由结构：`/admin/auth/*`、`/admin/system/*`、`/admin/upload/*`、`/admin/user/*`、`/admin/site/*`。

### 前端

- 开发端口: 2800
- Hash 路由模式（`createWebHashHistory`）
- 动态路由从 API 菜单加载
- 生产环境配置覆盖: `admin/public/admin/config.js`（`APP_CONFIG` 变量）
- API 地址在 `admin/src/config/index.js` 的 `API_URL` 中配置
- 构建: `cd admin && pnpm run build`，输出到 `admin/dist/`

### 生产部署

构建时 Taskfile 将 `server/configs/` 和 `server/web/` 复制到 `dist/`。服务器通过 `http.toml` 的 `static` 配置指向 `./web` 目录提供前端静态文件。因此前端构建产物需放入 `server/web/`。

## 注意事项

- `server/internal/` 下的包遵循 Go 的 `internal` 约束，不可被外部模块导入
- `server/web/uploads/` 已被 gitignore（用户上传目录）
- 前端 `.gitignore` 中忽略了 `pnpm-lock.yaml`，安装时无需该文件
- 没有配置测试框架，`server/test/` 下的测试文件为手动运行
- 没有配置 linter 或 formatter
- JWT 认证中间件: `server/internal/http/common/middleware/auth_jwt.go`