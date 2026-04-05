# 测试文档

本项目包含完整的测试套件，覆盖了 controller、repository、model 和 router 层。

## 测试结构

```
cat-dataserver/
├── controller/
│   ├── cat_controller_test.go
│   ├── cat_event_controller_test.go
│   ├── cat_action_controller_test.go
│   ├── cat_fsm_controller_test.go
│   ├── site_controller_test.go
│   ├── site_action_controller_test.go
│   └── site_fsm_controller_test.go
├── repository/
│   ├── cat_repository_test.go
│   ├── cat_event_repository_test.go
│   ├── cat_action_repository_test.go
│   ├── cat_fsm_repository_test.go
│   ├── site_repository_test.go
│   ├── site_action_repository_test.go
│   └── site_fsm_repository_test.go
├── model/
│   └── pagination_test.go
├── middleware/
│   ├── action_processor_test.go
│   └── validation_test.go
└── router/
    └── router_test.go
```

## 运行测试

### 运行所有测试

```bash
go test ./...
```

### 运行特定包的测试

```bash
# 测试 controller 层
go test ./controller/...

# 测试 repository 层
go test ./repository/...

# 测试 model 层
go test ./model/...

# 测试 router 层
go test ./router/...
```

### 运行特定测试函数

```bash
# 运行单个测试
go test -v ./controller/... -run TestCreateCat

# 运行匹配模式的测试
go test -v ./repository/... -run TestCatRepository.*
```

### 详细输出模式

```bash
go test -v ./...
```

### 覆盖率报告

```bash
# 生成覆盖率报告
go test ./... -cover

# 生成覆盖率文件
go test ./... -coverprofile=coverage.out

# 查看覆盖率详情
go tool cover -html=coverage.out
```

## 测试覆盖范围

### Controller 层测试
- **CatController**: 测试猫的 CRUD 操作、分页查询
- **CatEventController**: 测试猫事件的 CRUD 操作、按猫ID和站点ID查询
- **CatActionController**: 测试猫操作的 CRUD 操作、按猫ID、站点ID和用户ID查询，验证 FSM 自动更新
- **CatFSMController**: 测试猫状态机的只读操作
- **SiteController**: 测试站点的 CRUD 操作
- **SiteActionController**: 测试站点动作的 CRUD 操作、验证 FSM 自动更新
- **SiteFSMController**: 测试站点状态机的只读操作

### Repository 层测试
- **CatRepository**: 测试数据库 CRUD 操作、分页查询
- **CatEventRepository**: 测试事件数据库操作、按ID、猫ID、站点ID查询
- **CatActionRepository**: 测试操作数据库操作、按ID、猫ID、站点ID、用户ID查询
- **CatFSMRepository**: 测试状态数据库操作、更新体温、体重、修剪指甲时间
- **SiteRepository**: 测试站点数据库操作
- **SiteActionRepository**: 测试站点动作数据库操作
- **SiteFSMRepository**: 测试站点状态数据库操作、更新各类时间字段

### Middleware 层测试
- **ActionProcessor**: 测试猫咪动作处理器、状态机自动更新
- **Validators**: 测试自定义验证器

### Model 层测试
- **PaginationRequest**: 测试分页参数处理、默认值设置
- **PaginationResponse**: 测试分页响应生成
- **CatEventType**: 测试事件类型常量
- **CatActionType**: 测试操作类型常量

### Router 层测试
- **路由配置**: 测试所有 API 路由是否正确注册
- **健康检查**: 测试 /health 端点
- **CRUD 端点**: 测试各资源的创建、读取、更新、删除端点
- **查询端点**: 测试按条件查询的端点

## 测试特点

1. **隔离性**: 每个测试使用内存数据库 (`:memory:`)，确保测试之间互不影响
2. **全面性**: 覆盖正常流程、边界情况和错误场景
3. **快速执行**: 使用测试模式和内存数据库，执行速度快
4. **清晰的断言**: 使用 Go 标准测试框架，断言清晰明确

## 测试数据

测试使用模拟数据进行操作，主要包括：
- 猫的基本信息（名称、照片、种类、性别等）
- 猫事件（生病、受伤、怀孕等）
- 猫动作（测体温、称重、修剪指甲、绝育、驱虫、疫苗、洗澡）
- 猫状态（体温、体重、修剪指甲时间等）
- 站点信息（名称、地址、管理员电话等）
- 站点动作（消毒、喂食、喂水、逗猫、清理猫砂）
- 站点状态（消毒时间、喂食时间、喂水时间、逗猫时间、清理猫砂时间）

## 注意事项

1. 测试使用内存数据库，不会影响实际的 `test.db` 文件
2. 某些测试可能因为数据库记录未找到而返回 404，这是预期的行为
3. 运行完整测试套件建议使用 `-timeout` 参数限制执行时间
4. Controller 层测试需要初始化数据库，执行时间相对较长

## 持续集成建议

建议在 CI/CD 流程中运行测试：

```yaml
# GitHub Actions 示例
- name: Run tests
  run: go test ./... -v -cover

- name: Generate coverage
  run: |
    go test ./... -coverprofile=coverage.out
    go tool cover -html=coverage.out -o coverage.html
```

## 扩展测试

添加新测试时：

1. 在相应的包中创建 `*_test.go` 文件
2. 以 `Test` 开头命名测试函数
3. 使用 `testing.T` 参数进行断言
4. 为测试函数添加清晰的注释说明测试目的
5. 保持测试的独立性和可重复性
