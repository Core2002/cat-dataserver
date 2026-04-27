# Cat DataServer

一个用于管理猫咪数据和设施管理的后端服务，支持自动状态机更新、动作记录和设施维护管理。

## 核心功能

### 1. 动作记录与自动状态机更新

通过统一的动作接口，自动记录所有操作并更新猫咪和设施状态机。

**特性**：
- 自动记录所有猫咪动作（测体温、称重、修剪指甲、绝育、驱虫、疫苗、洗澡等）
- 根据动作类型自动更新猫咪状态机（体温、体重、修剪指甲时间）
- 自动记录所有站点动作（消毒、喂食、喂水、逗猫、清理猫砂）
- 自动更新设施状态机（消毒时间、喂食时间、喂水时间、逗猫时间、清理猫砂时间）
- 智能数据解析（自动提取体温、体重等数据）
- 完整的错误处理和日志记录
- 状态机更新失败时自动回滚动作记录

### 2. 数据模型

#### 猫咪相关
- **Cat**: 猫咪基本信息（名称、照片、种类、性别、主人信息等）
- **CatEvent**: 猫咪事件（生病、受伤、怀孕、分娩、死亡、合同解除）
- **CatAction**: 猫咪操作记录（测体温、称重、修剪指甲、绝育、驱虫、疫苗、洗澡）
- **CatFSM**: 猫咪状态机（体温、体重、修剪指甲时间）

#### 设施相关
- **Site**: 设施信息（猫舍、站点等，包含名称、地址、管理员电话）
- **SiteAction**: 站点操作记录（消毒、喂食、喂水、逗猫、清理猫砂）
- **SiteFSM**: 设施状态机（消毒时间、喂食时间、喂水时间、逗猫时间、清理猫砂时间）

### 3. 数据验证

- 完整的数据校验机制
- 自定义验证器（事件类型、动作类型、站点动作类型、电话号码等）
- 友好的中文错误提示

## 快速开始

### 环境要求

- Go 1.25.0 或更高版本
- SQLite 数据库（项目使用文件型数据库，无需额外安装）

### 1. 克隆项目

```bash
git clone <repository-url>
cd cat-dataserver
```

### 2. 安装依赖

```bash
go mod download
```

### 3. 编译

```bash
# Windows
go build -o cat-dataserver.exe .

# Linux/Mac
go build -o cat-dataserver
```

### 4. 运行

```bash
# Windows
.\cat-dataserver.exe

# Linux/Mac
./cat-dataserver
```

服务将在 `http://localhost:5100` 启动。

首次运行时会自动创建 `cat-dataserver.db` 数据库文件。

### 5. 配置

编辑 `config/config.go` 文件可修改配置：

```go
const (
    // ServerPort 服务器端口
    ServerPort = ":5100"
    // DatabaseDSN 数据库连接字符串
    DatabaseDSN = "cat-dataserver.db"
)
```

## API 接口

### 健康检查

**GET** `/health`

检查服务是否正常运行。

### 猫咪管理

#### 获取猫咪列表（分页）

**GET** `/cats/page?page=1&pageSize=10`

#### 获取单个猫咪

**GET** `/cats/:id`

#### 创建猫咪

**POST** `/cats`

```json
{
  "cat_id": 1,
  "cat_name": "小白",
  "cat_photo_uri": "https://example.com/cat.jpg",
  "cat_type": "英短",
  "cat_gender": "公",
  "master_name": "张三",
  "master_phone_number": "13800138000"
}
```

#### 更新猫咪

**PUT** `/cats/:id`

#### 删除猫咪

**DELETE** `/cats/:id`

### 猫咪事件管理

#### 获取事件列表（分页）

**GET** `/cat-events/page?page=1&pageSize=10`

#### 获取单个事件

**GET** `/cat-events/:id`

#### 按猫咪ID查询事件

**GET** `/cat-events/cat/:cat_id`

#### 按站点ID查询事件

**GET** `/cat-events/site/:site_id`

#### 创建事件

**POST** `/cat-events`

```json
{
  "event_id": 1,
  "event_type": "生病",
  "site_id": 1,
  "user_id": 1,
  "cat_id": 1,
  "detail": "猫咪出现呕吐症状"
}
```

#### 更新事件

**PUT** `/cat-events/:id`

#### 删除事件

**DELETE** `/cat-events/:id`

### 猫咪动作记录（自动更新状态机）

#### 获取动作列表（分页）

**GET** `/cat-actions/page?page=1&pageSize=10`

#### 获取单个动作

**GET** `/cat-actions/:id`

#### 按猫咪ID查询动作

**GET** `/cat-actions/cat/:cat_id`

#### 按站点ID查询动作

**GET** `/cat-actions/site/:site_id`

#### 按用户ID查询动作

**GET** `/cat-actions/user/:user_id`

#### 创建动作（自动更新 FSM）

**POST** `/cat-actions`

请求头：
```
X-User-ID: 1
```

```json
{
  "cat_id": 1,
  "site_id": 1,
  "action_type": "测体温",
  "action_detail": "{\"temperature_c\": 39.5}"
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
    "action_detail": "{\"temperature_c\": 39.5}",
    "created_at": "2026-03-25T10:00:00Z"
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

> **注意**：CatAction 为只读 + 创建模式，不提供直接更新和删除接口。

### 猫咪状态机管理（只读，由动作驱动更新）

#### 获取状态机列表（分页）

**GET** `/cat-fsms/page?page=1&pageSize=10`

#### 获取单个状态机

**GET** `/cat-fsms/:cat_id`

#### 按站点ID查询状态机

**GET** `/cat-fsms/site/:site_id`

> **注意**：CatFSM 由 CatAction 自动驱动更新，不提供直接创建、更新和删除接口。

### 设施管理

#### 获取设施列表（分页）

**GET** `/sites/page?page=1&pageSize=10`

#### 获取单个设施

**GET** `/sites/:id`

#### 创建设施

**POST** `/sites`

```json
{
  "site_id": 1,
  "site_name": "阳光猫舍",
  "site_address": "北京市朝阳区",
  "site_admin_phone_number": "13800138000"
}
```

#### 更新设施

**PUT** `/sites/:id`

#### 删除设施

**DELETE** `/sites/:id`

### 设施状态机管理（只读，由动作驱动更新）

#### 获取设施状态机列表（分页）

**GET** `/site-fsms/page?page=1&pageSize=10`

#### 获取单个设施状态机

**GET** `/site-fsms/:site_id`

#### 按站点ID查询状态机

**GET** `/site-fsms/site/:site_id`

> **注意**：SiteFSM 由 SiteAction 自动驱动更新，不提供直接创建、更新和删除接口。

### 站点动作记录（自动更新状态机）

#### 获取动作列表（分页）

**GET** `/site-actions/page?page=1&pageSize=10`

#### 获取单个动作

**GET** `/site-actions/:action_id`

#### 按站点ID查询动作

**GET** `/site-actions/site/:site_id`

#### 按用户ID查询动作

**GET** `/site-actions/user/:user_id`

#### 创建动作（自动更新 FSM）

**POST** `/site-actions`

请求头：
```
X-User-ID: 1
```

```json
{
  "site_id": 1,
  "action_type": "喂食",
  "action_detail": "{\"food_type\": \"猫粮\", \"amount\": \"100g\"}"
}
```

**响应**：
```json
{
  "action": {
    "action_id": 1,
    "site_id": 1,
    "user_id": 1,
    "action_type": "喂食",
    "action_detail": "{\"food_type\": \"猫粮\", \"amount\": \"100g\"}",
    "created_at": "2026-04-05T10:00:00Z"
  },
  "fsm": {
    "site_id": 1,
    "last_disinfect_time": "2026-04-04T08:00:00Z",
    "last_feed_time": "2026-04-05T10:00:00Z",
    "last_give_water_time": "2026-04-05T09:30:00Z",
    "last_play_time": "2026-04-05T09:00:00Z",
    "last_clean_litter_time": "2026-04-05T08:00:00Z"
  }
}
```

## 支持的动作类型

### 猫咪动作类型

| 动作类型 | action_type | action_detail 格式 (JSON) | 更新的 FSM 字段 |
|---------|-------------|-------------------------|----------------|
| 测体温 | "测体温" | `{"temperature_c": 39.5}` | TemperatureC |
| 称重 | "称重" | `{"weight_kg": 5.2}` | WeightKG |
| 修剪指甲 | "修剪指甲" | `{"notes": "修剪完成"}` | TrimNailsTime |
| 绝育 | "绝育" | `{"notes": "手术完成"}` | - |
| 驱虫 | "驱虫" | `{"drug_name": "xxx", "dosage": "xxx"}` | - |
| 疫苗 | "疫苗" | `{"vaccine_name": "xxx", "batch_no": "xxx"}` | - |
| 洗澡 | "洗澡" | `{"notes": "洗澡完成"}` | - |

### 站点动作类型

| 动作类型 | action_type | action_detail 格式 (JSON) | 更新的 FSM 字段 |
|---------|-------------|-------------------------|----------------|
| 消毒 | "消毒" | `{"disinfectant": "xxx", "notes": "xxx"}` | LastDisinfectTime |
| 喂食 | "喂食" | `{"food_type": "猫粮", "amount": "100g"}` | LastFeedTime |
| 喂水 | "喂水" | `{"water_type": "纯净水"}` | LastGiveWaterTime |
| 逗猫 | "逗猫" | `{"duration": 30, "notes": "xxx"}` | LastPlayTime |
| 清理猫砂 | "清理猫砂" | `{"litter_type": "xxx"}` | LastCleanLitter |

## 数据格式规范

### 测体温动作详情

```json
{
  "temperature_c": 39.5
}
```

**字段说明**：
- `temperature_c`：体温值，单位：摄氏度
- 范围：0°C - 50°C
- 类型：浮点数

### 称重动作详情

```json
{
  "weight_kg": 5.2
}
```

**字段说明**：
- `weight_kg`：体重值，单位：千克
- 范围：0.1kg - 25kg
- 类型：浮点数

### 驱虫动作详情

```json
{
  "drug_name": "福来恩",
  "dosage": "一支"
}
```

### 疫苗动作详情

```json
{
  "vaccine_name": "猫三联",
  "batch_no": "20260401001"
}
```

### 喂食动作详情

```json
{
  "food_type": "猫粮",
  "amount": "100g",
  "notes": "备注信息"
}
```

### 逗猫动作详情

```json
{
  "duration": 30,
  "notes": "玩了逗猫棒"
}
```

## 数据验证

系统会对所有输入数据进行严格验证：

- **必填字段检查**: 确保关键字段不为空
- **数据类型验证**: 确保数据类型正确
- **范围限制**: 
  - 体温：0°C - 50°C
  - 体重：0.1kg - 25kg
- **枚举值验证**: 
  - 事件类型必须是预定义值之一
  - 动作类型必须是预定义值之一
  - 站点动作类型必须是预定义值之一
- **格式验证**: 
  - 电话号码：中国大陆 11 位手机号
  - 照片 URL：有效的 URL 格式
- **长度限制**: 
  - 名称：1-100 字符
  - 详情：1-1000 字符

验证失败会返回详细的中文错误信息。

## 架构设计

### 分层架构

```
Controller 层  →  Middleware 层  →  Repository 层
(请求处理)      (业务逻辑)         (数据访问)
```

### 核心组件

#### 1. Controller 层
- 处理 HTTP 请求和响应
- 参数绑定和验证
- 调用业务逻辑层

#### 2. Middleware 层
- **ActionProcessor**: 猫咪动作处理器，负责自动更新猫咪状态机
- **SiteActionProcessor**: 站点动作处理器，负责自动更新站点状态机
- **Validators**: 自定义验证器，确保数据有效性
- CORS 中间件：支持跨域请求

#### 3. Repository 层
- 数据库操作封装
- CRUD 基础方法
- 分页查询支持

#### 4. Model 层
- 数据模型定义
- 数据库映射
- 分页模型

### 工作流程

1. 客户端发送请求到 Controller
2. Controller 进行参数绑定和验证
3. 验证通过后调用 Repository 或 ActionProcessor
4. ActionProcessor 处理业务逻辑并更新状态机
5. Repository 执行数据库操作
6. 返回结果给客户端

## 测试

本项目包含完整的测试套件，详见 [TESTING.md](TESTING.md)。

### 运行所有测试

```bash
go test ./...
```

### 运行特定包的测试

```bash
go test ./controller -v
go test ./repository -v
go test ./middleware -v
```

### 生成覆盖率报告

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## 技术栈

- **Go 1.25.0**: 编程语言
- **Gin 1.10.1**: Web 框架
- **GORM 1.25.12**: ORM 框架
- **SQLite 3**: 数据库（纯 Go 实现）
- **gin-contrib/cors**: CORS 中间件
- **go-playground/validator**: 数据验证库

## 项目结构

```
cat-dataserver/
├── config/              # 配置文件
│   └── config.go        # 服务器端口和数据库配置
├── controller/          # 控制器层
│   ├── cat_controller.go
│   ├── cat_event_controller.go
│   ├── cat_action_controller.go
│   ├── cat_fsm_controller.go
│   ├── site_controller.go
│   ├── site_action_controller.go
│   └── site_fsm_controller.go
├── database/            # 数据库初始化
│   └── database.go      # 数据库连接和迁移
├── middleware/          # 中间件层
│   ├── action_processor.go    # 猫咪动作处理器
│   ├── site_action_processor.go # 站点动作处理器
│   └── validation.go         # 数据验证
├── model/               # 数据模型
│   ├── cat.go           # 猫模型（含 CatProfile）
│   ├── cat_event.go     # 事件模型（含 CatAction）
│   ├── cat_fsm.go       # 猫状态机
│   ├── site.go          # 站点模型
│   ├── site_action.go   # 站点动作模型
│   ├── site_fsm.go      # 站点状态机
│   ├── pagination.go    # 分页模型
│   └── time_helper.go   # 时间格式辅助
├── repository/          # 数据访问层
│   ├── cat_repository.go
│   ├── cat_event_repository.go
│   ├── cat_action_repository.go
│   ├── cat_fsm_repository.go
│   ├── site_repository.go
│   ├── site_action_repository.go
│   └── site_fsm_repository.go
├── router/              # 路由配置
│   └── router.go        # 路由注册和中间件配置
├── main.go              # 程序入口
├── go.mod               # Go 模块定义
├── go.sum               # 依赖版本锁定
├── README.md            # 项目说明
├── TESTING.md           # 测试文档
├── INTEGRATION_GUIDE.md # 集成指南
└── REFACTORING_SUMMARY.md # 重构总结
```

## 文档

- **[INTEGRATION_GUIDE.md](INTEGRATION_GUIDE.md)**: 动作处理集成指南
- **[TESTING.md](TESTING.md)**: 测试文档和测试指南
- **[REFACTORING_SUMMARY.md](REFACTORING_SUMMARY.md)**: 架构重构总结

## 最佳实践

### 1. 使用统一接口

创建动作使用统一的 `POST /cat-actions` 和 `POST /site-actions` 接口，系统会自动处理状态机更新。

### 2. 数据格式

将关键数据（如体温、体重）使用 JSON 格式放在 `action_detail` 中，便于自动解析。

### 3. 错误处理

状态机更新失败时会自动回滚动作记录，确保数据一致性。

### 4. 分页查询

使用 `/page` 端点进行分页查询，避免一次返回大量数据。

### 5. 用户认证

创建猫咪动作和站点动作时都需要在请求头中提供 `X-User-ID`。

## 扩展开发

### 添加新的猫咪动作类型

1. 在 `model/cat_event.go` 中添加新的 `CatActionType` 常量
2. 在 `middleware/validation.go` 的 `validateCatActionType` 中添加验证
3. 在 `middleware/action_processor.go` 的 `updateFSM` 中添加处理逻辑
4. 如需更新 FSM，添加对应的 repository 方法
5. 编写测试验证功能

### 添加新的站点动作类型

1. 在 `model/site_action.go` 中添加新的 `SiteActionType` 常量
2. 在 `middleware/validation.go` 的 `validateSiteActionType` 中添加验证
3. 在 `middleware/site_action_processor.go` 的 `updateFSM` 中添加处理逻辑
4. 如需更新 FSM，添加对应的 repository 方法
5. 编写测试验证功能

### 添加新的状态字段

1. 在 `model/cat_fsm.go` 或 `model/site_fsm.go` 中添加新字段
2. 在对应的 repository 中添加更新方法
3. 在 `middleware/action_processor.go` 或 `middleware/site_action_processor.go` 中添加处理逻辑
4. 更新测试验证功能

### 添加新的查询接口

1. 在 repository 层添加查询方法
2. 在 controller 层添加处理函数
3. 在 router 层注册新路由
4. 编写测试验证功能

## CORS 配置

服务默认允许来自以下源的跨域请求：
- `http://localhost:5000`
- `http://localhost:5200`
- `https://tls.internal`

如需修改，请编辑 `router/router.go` 中的 CORS 配置。

## 许可证

MIT License

Copyright (c) 2026 Core2002


## 贡献

欢迎提交 Issue 和 Pull Request！

提交 PR 时请注意：
1. 确保代码通过所有测试
2. 遵循现有的代码风格
3. 添加必要的测试用例
4. 更新相关文档

## 联系方式

如有问题，请提交 Issue。

## 更新日志

### v1.1.0
- 新增站点动作（SiteAction）功能
- 新增站点状态机自动更新
- 猫咪和站点状态机改为只读，由动作驱动更新
- 新增清理猫砂动作类型
- CORS 配置更新

### v1.0.0
- 初始版本发布
- 支持猫咪和设施的基本 CRUD 操作
- 实现动作记录和自动状态机更新
- 完整的数据验证机制
- 分页查询支持
- 全面的测试覆盖
