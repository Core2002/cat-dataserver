# Cat DataServer

一个用于管理猫咪数据的后端服务，支持自动状态机更新和动作记录。

## 核心功能

### 1. 动作记录与自动状态机更新

通过统一的动作接口，自动记录所有操作并更新猫咪状态机。

**特性**：
- 自动记录所有动作（喂食、测体温、体检等）
- 根据动作类型自动更新状态机
- 智能数据解析（自动提取体温、体重等数据）
- 完整的错误处理和日志记录

### 2. 数据模型

- **Cat**: 猫咪基本信息
- **CatEvent**: 猫咪事件（生病、怀孕、死亡等）
- **CatAction**: 猫咪操作记录（喂食、测体温、体检等）
- **CatFSM**: 猫咪状态机（体温、体重、修剪指甲时间）
- **Site**: 设施信息（猫舍、站点等）

## 快速开始

### 1. 编译

```bash
go build -o cat-dataserver.exe .
```

### 2. 运行

```bash
.\cat-dataserver.exe
```

服务将在 `http://localhost:8080` 启动。

## API 接口

### 动作上报（自动更新状态机）

#### 创建动作（自动更新 FSM）

**POST** `/cat-actions`

```json
{
  "cat_id": 1,
  "site_id": 1,
  "user_id": 1,
  "action_type": "测体温",
  "action_detail": "39.5"
}
```

**响应**：
```json
{
  "action": {
    "action_id": 1,
    "cat_id": 1,
    "site_id": 1,
    "user_id": 1,
    "action_type": "测体温",
    "action_detail": "39.5"
  },
  "fsm": {
    "cat_id": 1,
    "site_id": 1,
    "temperature_c": 39.5,
    "weight_kg": 4.2,
    "trim_nails_time": "2026-03-22T23:45:00Z"
  }
}
```

### 其他 CRUD 接口

- **Cat**: `/cats` (GET, POST), `/cats/:id` (GET, PUT, DELETE)
- **CatEvent**: `/cat-events` (GET, POST), `/cat-events/:id` (GET, PUT, DELETE)
- **CatAction**: `/cat-actions` (GET, POST), `/cat-actions/:id` (GET, PUT, DELETE)
- **CatFSM**: `/cat-fsms` (GET, POST), `/cat-fsms/:id` (GET, PUT, DELETE)
- **Site**: `/sites` (GET, POST), `/sites/:id` (GET, PUT, DELETE)

## 支持的动作类型

### 更新状态机的动作

| 动作类型 | action_type | action_detail 格式 | 更新的 FSM 字段 |
|---------|-------------|-------------------|----------------|
| 测体温 | "测体温" | "39.5" | TemperatureC |
| 修剪指甲 | "修剪指甲" | "修剪指甲" | TrimNailsTime |
| 体检 | "体检" | "5.2" | WeightKG |

### 仅记录的动作

- "喂食"
- "喂水"
- "逗猫"
- "绝育"
- "驱虫"
- "疫苗"
- "清理猫砂"
- "环境消毒"
- "洗脚"

## 数据解析规则

### 体温解析

支持以下格式（自动提取数字）：
- "39.5"
- "体温:39.5"
- "温度 39.5℃"
- "测体温 39.5"

### 体重解析

支持以下格式（自动提取数字）：
- "5.2"
- "体重:5.2"
- "体重 5.2kg"
- "体检 体重5.2"

## 架构设计

### 分层架构

```
Controller 层  →  Middleware 层  →  Repository 层
(请求处理)      (业务逻辑)         (数据访问)
```

### 核心组件

- **ActionProcessor**: 动作处理器，负责自动更新状态机
- **Validators**: 自定义验证器，确保数据有效性
- **Repositories**: 数据访问层，封装数据库操作

## 测试

### 运行所有测试

```bash
go test ./...
```

### 运行特定包的测试

```bash
go test ./middleware -v
go test ./controller -v
```

## 技术栈

- **Go 1.x**: 编程语言
- **Gin**: Web 框架
- **GORM**: ORM 框架
- **SQLite**: 数据库

## 项目结构

```
cat-dataserver/
├── config/          # 配置文件
├── controller/      # 控制器层
├── database/        # 数据库初始化
├── middleware/     # 中间件（包含 ActionProcessor）
├── model/          # 数据模型
├── repository/     # 数据访问层
├── router/         # 路由配置
└── service/        # （预留，业务逻辑层）
```

## 文档

- **INTEGRATION_GUIDE.md**: 动作处理集成指南
- **REFACTORING_SUMMARY.md**: 架构重构总结

## 最佳实践

### 1. 使用统一接口

创建动作使用统一的 `POST /cat-actions` 接口，系统会自动处理状态机更新。

### 2. 数据格式

将关键数据（如体温、体重）直接放在 `action_detail` 中，便于自动解析。

### 3. 错误处理

即使状态机更新失败，动作记录也会保存，确保数据不丢失。

## 扩展开发

### 添加新的动作类型

1. 在 `model/cat_event.go` 中添加新的 `CatActionType` 常量
2. 在 `middleware/action_processor.go` 的 `updateFSM` 中添加处理逻辑
3. 如需更新 FSM，添加对应的 repository 方法
4. 编写测试验证功能

### 添加新的状态字段

1. 在 `model/cat_fsm.go` 中添加新字段
2. 在 `repository/cat_fsm_repository.go` 中添加更新方法
3. 在 `middleware/action_processor.go` 中添加处理逻辑
4. 更新测试验证功能

## 许可证

本项目仅供学习和研究使用。

## 贡献

欢迎提交 Issue 和 Pull Request！

## 联系方式

如有问题，请提交 Issue。
