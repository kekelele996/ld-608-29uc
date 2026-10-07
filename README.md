# 航空地勤周转保障平台（ground-turn）

面向机场地勤团队的**航班过站保障、资源调度、异常延误与任务签收**全栈系统。核心规则：**航班延误登记后，任务超时按「未关闭延误」顺延后的当前生效截止时间判定**，看板超时卡片与任务列表共用同一套口径。

## 快速启动（推荐）

```bash
cp .env.example .env && docker compose up -d
```

- 前端：<http://localhost:20108>
- 后端健康检查：<http://localhost:21108/health>
- 演示账号（口令均为 `demo123`）：`dispatcher`（地勤调度）、`team`（班组）、`resource`（资源管理员）、`supervisor`（运行督导）。登录弹窗可点选切换角色，验证按钮显隐与后端 RBAC。

## 超时/顺延判定口径（本项目核心，页面内也有同样说明）

> 设计决策：**保留原计划截止时间，不改写；另算一个“当前生效截止时间”。**

1. 任务的 `deadline`（原计划截止时间）**始终保留不变**，登记或关闭延误都不会回写该列。
2. 另算**当前生效截止时间**：

   ```
   effective_deadline = deadline + Σ 该航班上「尚未关闭」延误事件的分钟数
   ```

   延误事件 `resolved_at` 为空 = 未关闭，才参与顺延；**关闭过的延误不再顺延**。
3. 任务超时一律按 `effective_deadline` 判定：
   - 未签收 / 已签收进行中 / 阻塞任务：与**当前时间**比较；
   - 已签收**完成**的任务：与 **`actual_finish` 实际完成时间**比较。
4. **同一份计算结果**同时供给看板“超时任务”卡片和任务列表，保证“超时条数”与“顺延分钟数”对得上。

口径在前后端各有一份同构实现（唯一真源由后端返回，前端用于离线兜底）：

- 后端：`backend/src/utils/delayPolicy.go`（`SummarizeOpenDelays` / `EffectiveDeadline` / `EvaluateOverdue`），由 `services/ViewService.go` 统一装配，`constructors/GroundTask.go::BuildGroundTaskView` 落字段。
- 前端：`frontend/src/utils/delayPolicy.ts`（看板 `OverdueTaskCard`、任务页 `GroundTaskTable`、本地 mock 全部调用它）。
- 页面说明组件：`frontend/src/components/common/DeadlinePolicyBanner.tsx`（看板与任务页都渲染）。
- 任务页支持对**未签收**任务一键签收并显示**签收时间**（`accepted_at`），完成动作记录 `actual_finish` 并按其判定超时。

## 访问地址与接口示例

```bash
# 登录拿 token
curl -s -X POST http://localhost:21108/api/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"dispatcher","password":"demo123"}'

# 带 token 查看看板（含 overdue_tasks 与 open_delay_minutes）
curl -s http://localhost:21108/api/dashboard -H "Authorization: Bearer <token>"

# 登记一笔未关闭延误（立即顺延，不改 deadline）
curl -s -X POST http://localhost:21108/api/delay-events -H "Authorization: Bearer <token>" \
  -H 'Content-Type: application/json' \
  -d '{"turnaround_id":1,"delay_type":"WEATHER","minutes":20}'

# 关闭延误（该笔分钟停止顺延）
curl -s -X POST http://localhost:21108/api/delay-events/3/resolve -H "Authorization: Bearer <token>"

# 班组签收 / 完成任务
curl -s -X POST http://localhost:21108/api/ground-tasks/1/accept   -H "Authorization: Bearer <team-token>"
curl -s -X POST http://localhost:21108/api/ground-tasks/1/complete -H "Authorization: Bearer <team-token>"
```

## 本地开发

- 前端：`cd frontend && npm install && npm run dev`（端口 20108，请求统一走 `/api`）。
- 后端：`cd backend && go run .`（Go 1.22；需可连 MySQL，默认 `db:3306`，可通过环境变量覆盖）。
- 后端不可达时，前端自动切换到**相对当前时间**的本地种子数据（`frontend/src/mocks/seedData.ts`），离线也能演示顺延/超时。

## 技术栈

| 层 | 技术 |
|---|---|
| 前端 | React 18 + TypeScript + Vite + Ant Design 风格组件 + Zustand（Redux 风格 store） |
| 后端 | Go 1.22 + Gin + GORM |
| 数据库 | MySQL 8.0（`database/init.sql` 建库建表 + GORM AutoMigrate） |
| 认证鉴权 | JWT（HS256，自包含工具）+ RBAC 四角色 |
| 部署 | Docker Compose（frontend / backend / db 三服务） |

## 项目目录结构

```text
backend/src/
  config/        环境配置（PORT/DB_*/JWT_* 分散读取）
  constants/     枚举、错误码、错误消息、日志模板、状态文案
  models/        GORM 模型（GroundTask.Deadline 不改写，AcceptedAt/ActualFinish 可空）
  repositories/  数据访问层（含延误按航班批量汇总，避免 N+1）
  services/      ViewService(顺延/超时统一装配) + 各实体服务
  controllers/   控制器（读分组/写分组分开挂 RBAC）
  middlewares/   authMiddleware / rbacMiddleware / auditLogMiddleware / errorHandlerMiddleware / rateLimitMiddleware
  constructors/  默认模型 + 响应 DTO 构造（BuildGroundTaskView 落生效截止与超时）
  types/         请求/响应 DTO
  utils/         delayPolicy（唯一口径）、jwt、errors、response、formatters
  seed/          相对时间的本地种子
frontend/src/
  api/ stores/ types/ constants/ constructors/ components/common/
  hooks/ pages/ router/ utils/ mocks/
```

## 环境变量说明

| 变量 | 默认 | 说明 |
|---|---|---|
| `COMPOSE_PROJECT_NAME` | `ground-turn` | Compose 项目名 / 容器名前缀 |
| `FRONTEND_PORT` | `20108` | 前端宿主机端口 |
| `BACKEND_PORT` | `21108` | 后端宿主机端口（容器内 3000） |
| `DB_PORT` | `33060` | MySQL 宿主机端口（容器内 3306） |
| `DB_NAME / DB_USER / DB_PASSWORD` | `app_db / app_user / app_password` | 数据库凭据 |
| `JWT_SECRET / JWT_TTL_HOURS` | `local-dev-secret / 12` | JWT 签名密钥与有效期 |

## Docker 部署说明

- 根 `docker-compose.yml`：顶层 `name: ground-turn`，不写 `version:`；所有 `container_name` 带 `${COMPOSE_PROJECT_NAME:-ground-turn}` 前缀。
- 端口：前端 `${FRONTEND_PORT:-20108}:80`，后端 `${BACKEND_PORT:-21108}:3000`。
- 数据库使用**命名卷** `db_data`（不绑定挂载，避免中文路径权限问题），配置 healthcheck；后端 `depends_on: condition: service_healthy`，前端再依赖后端健康。
- 前端 Nginx：`location /api/` 反代到 `http://backend:3000/api/`，其余 `try_files $uri $uri/ /index.html;`。
- 常见问题：端口占用改 `.env` 后重启；重置数据执行 `docker compose down -v`；后端启动时自动重试连接并写入种子（已存在用户则跳过）。

## 枚举/常量出现位置清单

- **GroundTaskType**（CLEANING/CATERING/BAGGAGE/REFUEL/WATER_SERVICE/PUSHBACK）
  - 后端：`constants/GroundTaskType.go`、`models/GroundTask.go`、`constructors/GroundTask.go`、`services/GroundTask.go`（校验）、`constants/logTemplates.go`、`constants/errorMessages.go`、`seed/seed.go`
  - 前端：`constants/GroundTaskType.ts`、`types/GroundTaskType.ts`、`constants/statusText.ts`、`constructors/GroundTaskConstructor.ts`、延误/任务筛选与展示（`GroundTaskTable`、`DelaysPage`、mocks）
- **TurnaroundStatus**（ARRIVING/ON_STAND/IN_SERVICE/READY/DEPARTED/DELAYED）
  - 后端：`constants/TurnaroundStatus.go`、`models/FlightTurnaround.go`、`constructors/FlightTurnaround.go`、`services/FlightTurnaround.go`、`constants/logTemplates.go`、`constants/statusText.go`
  - 前端：`constants/TurnaroundStatus.ts`、`types/TurnaroundStatus.ts`、`TurnaroundTimeline.tsx`、看板筛选/状态徽标、mocks
- **ResourceStatus**（AVAILABLE/BOOKED/MAINTENANCE/OFFLINE）
  - 后端：`constants/ResourceStatus.go`、`models/GroundResource.go`、`constructors/GroundResource.go`、`services/GroundResource.go`、`constants/logTemplates.go`
  - 前端：`constants/ResourceStatus.ts`、`types/ResourceStatus.ts`、`ResourcesPage.tsx`、`ResourceCalendar.tsx`
- 另有贯穿全栈的 **GroundTaskStatus**（PENDING/ACCEPTED/IN_PROGRESS/BLOCKED/COMPLETED）与 RBAC **Role**（DISPATCHER/TEAM/RESOURCE_MANAGER/SUPERVISOR），分别落在前后端 constants、types、store、按钮显隐、`rbacMiddleware`。

## 横切关注点

- **RBAC**：读接口任意登录角色；派工/航班/延误登记关闭限地勤调度与运行督导；签收/完成/阻塞限班组/调度/督导；资源与预约限资源管理员/督导。触达：`app_user` 表、`authMiddleware`/`rbacMiddleware`、路由分组、前端登录态与按钮显隐、错误码 `FORBIDDEN`。
- **操作日志**：派工、签收、完成、阻塞、延误登记/关闭、预约冲突/释放等写操作经 `auditLogMiddleware` 落 `audit_log`，模板集中在 `constants/logTemplates.go`（每实体 ≥4 条）。
- **全局错误处理**：service 抛 `utils.AppError`（带错误码），controller 二次包装，`errorHandlerMiddleware` 兜底；错误码在 `constants/errorCodes.go`、消息在 `constants/errorMessages.go`。
- **限流**：`rateLimitMiddleware` 按 IP 固定窗口（默认 300 次/分钟）。

## 为什么会牵一发动全身

顺延/超时口径被刻意拆到 `utils/delayPolicy`（后端）与 `delayPolicy.ts`（前端），再被 service/constructor/store/看板卡片/任务表格/mock 共同依赖；改一条判定规则必须同步类型 DTO（`effective_deadline/open_delay_minutes/overdue/compare_basis`）、构造器、日志模板、页面说明与种子。枚举、错误码、日志模板、状态文案同样分散在前后端 `constants/` 多层引用，新增一个状态值会触达模型、构造器、校验、筛选器、徽标展示与 README。

## 测试

- 后端：`cd backend && go test ./...`
  - `utils/delayPolicy_test.go`：未关闭汇总、生效截止、按当前/实际时间判定、关闭延误后超时条数回落。
  - `constructors/policy_integration_test.go`：DTO 装配层端到端对账（超时条数与顺延分钟一致）。

## License

MIT
