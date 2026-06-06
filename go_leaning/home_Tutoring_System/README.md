# Go 家教系统

这是一个使用 Go 重构的家教管理系统后端项目，原项目为 Python/Django 版本。当前 Go 版本使用 Gin 作为 Web 框架，Gorm 作为 ORM，SQLite 作为本地开发数据库。

项目还处于学习和重构阶段，目前主要完成了用户模块和订单模块的基础接口骨架。

## 技术栈

- Go
- Gin
- Gorm
- SQLite
- JWT
- bcrypt

## 项目结构

```text
home_Tutoring_System/
├── config/          # 项目配置，如端口、数据库路径、JWT 密钥
├── database/        # 数据库初始化和自动迁移
├── handler/         # 具体接口处理逻辑
├── middleware/      # JWT 登录认证中间件
├── model/           # Gorm 数据模型
├── router/          # 路由注册
├── main.go          # 项目入口
├── go.mod
└── go.sum
```

## 启动项目

进入项目目录：

```bash
cd go_leaning/home_Tutoring_System
```

安装依赖：

```bash
go mod tidy
```

启动服务：

```bash
go run .
```

默认服务地址：

```text
http://localhost:8080
```

健康检查：

```http
GET /ping
```

返回示例：

```json
{
  "message": "ping"
}
```

## 数据库

当前使用 SQLite，数据库文件配置在 `config/config.go`：

```go
const DBPath = "jiajiao.db"
```

项目启动时会自动执行 Gorm 迁移：

```go
database.Migrate()
```

当前自动迁移的模型包括：

- User
- TeacherProfile
- ChatMessage
- Match
- TeacherFavorite
- Order

## 已有接口

### 用户注册

```http
POST /api/user/register
```

学生注册请求示例：

```json
{
  "username": "student01",
  "phone": "13800138000",
  "password": "123456",
  "role": "student"
}
```

老师注册请求示例：

```json
{
  "username": "teacher01",
  "phone": "13900139000",
  "password": "123456",
  "role": "teacher",
  "subject": "数学",
  "teaching_years": 3
}
```

成功返回示例：

```json
{
  "message": "注册成功",
  "data": {
    "id": 1,
    "username": "student01",
    "phone": "13800138000",
    "role": "student"
  }
}
```

### 用户登录

```http
POST /api/user/login
```

请求示例：

```json
{
  "username": "student01",
  "password": "123456"
}
```

或者使用手机号：

```json
{
  "phone": "13800138000",
  "password": "123456"
}
```

注意：当前登录接口已经完成账号查询和 bcrypt 密码校验，但还没有返回 JWT token。后续需要补充 token 签发逻辑，否则需要登录认证的接口无法正常使用。

### 查看个人信息

```http
GET /api/user/me
Authorization: Bearer <token>
```

该接口通过 JWT 中间件获取当前登录用户 ID，然后查询用户信息。

成功返回示例：

```json
{
  "message": "获取个人身份信息成功",
  "data": {
    "id": 1,
    "username": "student01",
    "phone": "13800138000",
    "role": "student"
  }
}
```

### 修改个人信息

```http
POST /api/user/me
Authorization: Bearer <token>
```

请求示例：

```json
{
  "new_name": "new_student_name"
}
```

成功返回示例：

```json
{
  "message": "修改成功"
}
```

## 订单接口

订单相关接口需要携带 JWT：

```http
Authorization: Bearer <token>
```

### 创建订单

```http
POST /api/order/create
```

请求示例：

```json
{
  "teacher_id": 1,
  "subject": "数学",
  "scheduled_time": "2026-06-06T19:30:00+08:00",
  "duration": 120,
  "price": 300,
  "address": "北京市海淀区",
  "remarks": "希望重点讲函数和导数"
}
```

说明：

- `student_id` 不由前端传入，而是从 JWT token 中读取当前用户 ID。
- `teacher_id` 对应 `TeacherProfile` 的 ID。
- `scheduled_time` 使用 RFC3339 时间格式。

成功返回示例：

```json
{
  "message": "创建成功"
}
```

### 查看订单详情

```http
GET /api/order/:id
Authorization: Bearer <token>
```

示例：

```http
GET /api/order/1
```

这里的 `1` 是订单 ID，由前端从订单列表或页面状态中取得后拼到 URL 中。

成功返回示例：

```json
{
  "message": "查询成功",
  "data": {
    "student_username": "student01",
    "teacher_username": "teacher01",
    "subject": "数学",
    "scheduledTime": "2026-06-06T19:30:00+08:00",
    "duration": 120,
    "price": 300,
    "address": "北京市海淀区",
    "remarks": "希望重点讲函数和导数"
  }
}
```

## JWT 认证流程

需要登录的接口会经过 `middleware.AuthRequired()` 中间件。

前端请求头格式：

```http
Authorization: Bearer <token>
```

中间件验证通过后，会把用户信息写入 Gin 的 Context：

```go
c.Set("userID", claims.UserID)
c.Set("username", claims.Username)
c.Set("role", claims.Role)
```

handler 中可以这样读取：

```go
userID := c.GetUint("userID")
username := c.GetString("username")
role := c.GetString("role")
```

## 当前开发状态

已完成：

- 用户注册
- 用户登录的账号密码校验
- JWT 认证中间件
- 查看和修改个人信息
- 创建订单
- 查看订单详情

待完善：

- 登录成功后签发并返回 JWT token
- 退出登录接口
- 订单列表接口
- 订单状态修改接口
- 老师信息列表和详情接口
- 权限校验，例如只有学生可以创建订单
- 更完整的错误处理
- 单元测试或接口测试

## 学习重点

这个项目适合用来练习：

- Gin 路由分组
- JSON 请求体绑定
- Gorm 增删改查
- Gorm 模型关联和 `Preload`
- JWT 中间件
- bcrypt 密码加密和校验
- 从 Django/DRF 思维迁移到 Go/Gin 思维

