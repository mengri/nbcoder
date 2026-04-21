# 项目目录结构与规范

## 目录结构

```
nbcoder/
├── backend/         # 后端代码（DDD 规范）
│   ├── domain/          # 领域层
│   ├── application/     # 应用层
│   ├── infrastructure/  # 基础设施层
│   └── interfaces/      # 接口层
├── web/             # 前端代码根目录
└── agents.md        # 目录结构与开发规范说明
```

## 规范说明

- 后端 backend 采用 DDD（领域驱动设计）分层：
  - domain：核心业务、实体、聚合、值对象、领域服务、仓储接口
  - application：用例、应用服务、DTO、协调领域对象完成业务流程
  - infrastructure：与外部资源交互，如数据库、消息中间件、第三方服务
  - interfaces：对外暴露 API、UI、消息等，接收外部请求并调用应用层
- 前端 web 目录为前端项目根目录，结构和技术栈可根据实际需求扩展
- agents.md 用于记录项目结构、命名规范、开发约定等说明

## 命名规范

- 目录、文件、类、方法、变量命名统一采用小写+下划线（snake_case）或小驼峰（camelCase），保持风格一致。
- 领域层（domain）实体、值对象、服务等建议用英文单数名词。
- 应用层（application）用例、服务以“*Service”或“*UseCase”结尾。
- 基础设施层（infrastructure）实现类以“*Repository”、“*Adapter”等结尾。
- 接口层（interfaces）API 文件以“api_*.py”或“*_controller.js”等命名。

## 开发约定

- 每次开发按任务拆分，完成一个任务后需 commit 并 push。
- 代码需包含必要注释，重要业务逻辑需配合注释说明。
- 领域层禁止直接依赖基础设施层，依赖倒置通过接口实现。
- 前后端分离，接口通过 API 文档（如 OpenAPI/Swagger）约定。
- 统一使用 git 进行版本管理，分支命名建议 feature/xxx、fix/xxx、refactor/xxx。
- 推荐使用单元测试，测试代码与业务代码分离。

## 任务开发与分支管理规范

- 每次开发只能执行 prd/ddd/tasks/ 下的一个任务，需严格按任务拆分推进。
- 每个任务开发时，需从 main 分支创建独立 feature/xxx 分支（如 feature/A1.1-xxx）。
- 若任务有依赖，需在新分支主动合并依赖任务的 feature 分支。
- 每个任务完成后，需 commit 并 push 到远程对应分支。
- 分支最终合并（如合入 main/dev）由专人统一操作，开发者无需自行合并。
- 禁止在同一分支并行开发多个任务，确保每个 feature 分支只对应一个任务。
- 任务开发顺序、依赖关系需在 agents.md 或任务描述中明确。

## 前后端接口规范

### 统一响应格式

所有接口响应必须使用统一的格式：

```typescript
{
  code: number,      // 状态码：200 成功，400 请求错误，401 未授权，403 无权限，404 未找到，500 服务器错误
  message: string,   // 消息描述
  data: T | null     // 数据内容
}
```

**示例：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": "xxx",
    "name": "项目名称"
  }
}
```

**错误响应示例：**
```json
{
  "code": 400,
  "message": "参数错误：name 字段不能为空",
  "data": null
}
```

### HTTP 状态码使用

- `200 OK`：请求成功
- `201 Created`：资源创建成功
- `204 No Content`：删除成功（无返回数据）
- `400 Bad Request`：请求参数错误
- `401 Unauthorized`：未认证或 token 失效
- `403 Forbidden`：无权限访问
- `404 Not Found`：资源不存在
- `409 Conflict`：资源冲突（如重复创建）
- `500 Internal Server Error`：服务器内部错误

### 字段命名规范

- **统一使用 camelCase**（驼峰命名法）
- 示例：`projectId`, `createdAt`, `updatedAt`, `userId`
- **重要**：后端 DTO 中的 JSON 标签也必须使用 camelCase，例如：
  ```go
  type ProjectResponse struct {
      ID          string `json:"id"`
      Name        string `json:"name"`
      Description string `json:"description"`
      ProjectID   string `json:"projectId"`   // 正确
      CreatedAt   string `json:"createdAt"`    // 正确
      // 不要使用 project_id, created_at 等下划线命名
  }
  ```

**时间字段：**
- `createdAt`：创建时间，ISO 8601 格式（如 `2024-01-01T00:00:00Z`）
- `updatedAt`：更新时间，ISO 8601 格式

### 路由规范

#### RESTful 风格

使用资源导向的 URL 设计，遵循 RESTful 规范：

```
GET    /api/v1/resources           # 获取资源列表
POST   /api/v1/resources           # 创建资源
GET    /api/v1/resources/:id       # 获取单个资源
PUT    /api/v1/resources/:id       # 更新资源
DELETE /api/v1/resources/:id       # 删除资源
```

#### 项目级路由

需要关联项目的资源使用 `/projects/:projectId/` 前缀：

```
GET    /api/v1/projects/:projectId/cards         # 获取项目的卡片列表
POST   /api/v1/projects/:projectId/cards         # 创建卡片
GET    /api/v1/projects/:projectId/cards/:id     # 获取单个卡片
PUT    /api/v1/projects/:projectId/cards/:id     # 更新卡片
DELETE /api/v1/projects/:projectId/cards/:id     # 删除卡片
```

#### 操作类路由

对资源的操作使用子路由：

```
POST   /api/v1/projects/:projectId/pipelines/:id/start     # 启动流水线
POST   /api/v1/projects/:projectId/pipelines/:id/pause     # 暂停流水线
POST   /api/v1/projects/:projectId/pipelines/:id/resume    # 恢复流水线
POST   /api/v1/projects/:projectId/pipelines/:id/cancel    # 取消流水线
```

### 分页规范

#### 请求参数

```typescript
{
  page: number,      // 页码，从 1 开始
  size: number,      // 每页数量，默认 20
  keyword?: string   // 搜索关键词（可选）
}
```

#### 响应格式

```typescript
{
  code: 200,
  message: "success",
  data: {
    items: T[],      // 数据列表
    total: number,   // 总数量
    page: number,    // 当前页码
    size: number,    // 每页数量
    totalPages: number // 总页数
  }
}
```

### 请求方法使用

- `GET`：查询数据，不修改服务器状态
- `POST`：创建资源
- `PUT`：完整更新资源（全量更新）
- `PATCH`：部分更新资源（仅更新提供的字段）
- `DELETE`：删除资源

### ID 格式

所有 ID 使用字符串类型（UUID 格式）：

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000"
}
```

### 枚举值

枚举值使用大写字符串：

```json
{
  "status": "ACTIVE",
  "type": "PIPELINE"
}
```

### 文件上传/下载

#### 文件上传

使用 `multipart/form-data` 格式：

```typescript
const formData = new FormData()
formData.append('file', file)
formData.append('projectId', projectId)

POST /api/v1/projects/:projectId/documents/upload
Content-Type: multipart/form-data
```

#### 文件下载

设置响应类型为 `blob`：

```typescript
GET /api/v1/projects/:projectId/documents/:id/download
Response-Type: blob
```

### 认证与授权

- 使用 Bearer Token 认证
- Token 在请求头中携带：`Authorization: Bearer <token>`
- Token 存储在 `localStorage` 中，key 为 `token`

### 错误处理

#### 前端错误处理

```typescript
request.interceptors.response.use(
  (response) => {
    const { code, message, data } = response.data
    if (code !== 200) {
      ElMessage.error(message)
      return Promise.reject(new Error(message))
    }
    return data
  },
  (error) => {
    const message = error.response?.data?.message || error.message || '请求失败'
    ElMessage.error(message)

    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      window.location.href = '/login'
    }

    return Promise.reject(error)
  }
)
```

#### 后端错误处理

```go
// 统一错误响应
c.JSON(http.StatusBadRequest, gin.H{
    "code": 400,
    "message": "参数错误：name 字段不能为空",
    "data": nil,
})

// 业务错误
c.JSON(http.StatusInternalServerError, gin.H{
    "code": 500,
    "message": "创建项目失败：" + err.Error(),
    "data": nil,
})
```

### 接口版本管理

当前版本：`/api/v1`

新接口必须包含版本号，如需升级，创建新版本路径：
- `/api/v2/...`

### 接口文档

使用 OpenAPI/Swagger 规范编写接口文档，文档应包含：
- 接口路径和方法
- 请求参数（path、query、body）
- 响应格式和示例
- 错误码说明
- 认证要求

## 其他说明

- 可根据实际业务扩展目录结构，如增加 docs、scripts、tests 等。
- agents.md 持续维护，团队成员有新约定需及时补充。
