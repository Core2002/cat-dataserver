# 动作处理集成指南

## 架构概述

本项目采用**中间件驱动**的架构，将动作处理逻辑整合到现有业务流程中，而不是创建独立的上报路由。

## 核心组件

### 1. ActionProcessor (动作处理器)

位置：`middleware/action_processor.go`

**职责**：
- 处理所有 CatAction 创建请求
- 自动记录动作到数据库
- 根据动作类型自动更新状态机 (FSM)
- 解析动作详情中的数据（如体温、体重）

**主要方法**：

```go
// ProcessAction 处理动作并更新状态机
func (p *ActionProcessor) ProcessAction(action *model.CatAction) (*model.CatFSM, error)
```

### 2. CatActionController (动作控制器)

已集成 ActionProcessor，在创建动作时自动触发状态机更新。

## 使用方式

### 基本用法

创建任何 CatAction 都会自动触发状态机更新：

```bash
curl -X POST http://localhost:8080/cat-actions \
  -H "Content-Type: application/json" \
  -d '{
    "cat_id": 1,
    "site_id": 1,
    "user_id": 1,
    "action_type": "测体温",
    "action_detail": "39.5"
  }'
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
    "action_detail": "39.5",
    "created_at": "2026-03-22T23:45:00Z"
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

### 支持的动作类型

#### 更新状态机的动作

| 动作类型 | ActionDetail 格式 | FSM 字段 | 示例 |
|---------|-----------------|---------|------|
| 测体温 | `"39.5"` | `TemperatureC` | 测量体温后自动更新 |
| 修剪指甲 | `"修剪指甲"` | `TrimNailsTime` | 自动记录为当前时间 |
| 体检 | `"5.2"` | `WeightKG` | 体重记录自动更新 |

#### 仅记录的动作

| 动作类型 | 说明 |
|---------|------|
| 喂食 | 仅记录，不更新 FSM |
| 喂水 | 仅记录，不更新 FSM |
| 逗猫 | 仅记录，不更新 FSM |
| 绝育 | 仅记录，不更新 FSM |
| 驱虫 | 仅记录，不更新 FSM |
| 疫苗 | 仅记录，不更新 FSM |
| 清理猫砂 | 仅记录，不更新 FSM |
| 环境消毒 | 仅记录，不更新 FSM |
| 洗脚 | 仅记录，不更新 FSM |

### 数据解析格式

#### 体温解析

支持以下格式（自动提取数字）：
- `"39.5"`
- `"体温:39.5"`
- `"温度 39.5℃"`
- `"测体温 39.5"`

#### 体重解析

支持以下格式（自动提取数字）：
- `"5.2"`
- `"体重:5.2"`
- `"体重 5.2kg"`
- `"体检 体重5.2"`

## 集成流程

```
用户请求 → CatActionController.CreateCatAction
         ↓
    ActionProcessor.ProcessAction
         ↓
    1. 记录动作到数据库
         ↓
    2. 根据动作类型处理
         ↓
    - 测体温 → 解析体温 → 更新 FSM.TemperatureC
    - 修剪指甲 → 更新 FSM.TrimNailsTime
    - 体检 → 解析体重 → 更新 FSM.WeightKG
    - 其他 → 仅记录动作
         ↓
    返回动作和更新后的状态机
```

## 扩展指南

### 添加新的状态机更新逻辑

1. 在 `middleware/action_processor.go` 中添加新的处理方法

```go
func (p *ActionProcessor) updateNewField(action *model.CatAction, fsm *model.CatFSM) (*model.CatFSM, error) {
    // 解析数据
    value := parseNewValue(action.ActionDetail)
    if value == 0 {
        return fsm, nil
    }

    // 更新数据库
    if err := p.fsmRepo.UpdateNewField(action.CatID, value); err != nil {
        return nil, err
    }

    // 更新内存对象
    fsm.NewField = value
    return fsm, nil
}
```

2. 在 `updateFSM` 方法中添加 case

```go
case model.CatActionNewType:
    return p.updateNewField(action, fsm)
```

3. 如果需要新的数据库更新方法，在 `repository/cat_fsm_repository.go` 中添加

```go
func (r *CatFSMRepository) UpdateNewField(catID uint, value float32) error {
    return database.DB.Model(&model.CatFSM{}).Where("cat_id = ?", catID).Update("new_field", value).Error
}
```

## 测试

### 单元测试

```bash
# 测试动作处理器
go test ./middleware -v

# 测试控制器
go test ./controller -v
```

### 验证状态机更新

所有测试都包含对状态机更新的验证：

```go
// 创建动作后验证 FSM 是否更新
freshFSM, err := fsmRepo.FindByID(1)
if freshFSM.TemperatureC != expectedTemp {
    t.Errorf("Expected temperature %v, got %v", expectedTemp, freshFSM.TemperatureC)
}
```

## 与旧架构的对比

### 旧架构（已删除）

- 独立的上报路由 (`/api/actions`, `/api/cats/:cat_id/temperature` 等)
- 独立的 Service 层
- 专门的 ReportController

### 新架构（当前）

- **集成式**：动作处理整合到现有业务逻辑中
- **自动化**：通过 ActionProcessor 中间件自动处理
- **统一接口**：使用标准的 CatAction 创建接口
- **灵活性**：易于扩展，支持更多动作类型

## 最佳实践

### 1. 使用标准接口

创建动作使用统一的 `POST /cat-actions` 接口，而不是特定数据的接口。

### 2. 数据格式

将关键数据（如体温、体重）直接放在 `action_detail` 中，便于自动解析。

### 3. 错误处理

即使状态机更新失败，动作记录也会保存，确保数据不丢失。

### 4. 日志记录

所有动作处理都会记录日志，便于追踪和调试。

## 总结

新架构的优势：

1. **更简洁**：删除了不必要的独立路由和服务层
2. **更集成**：动作处理完全整合到现有业务流程
3. **更易维护**：逻辑集中在 ActionProcessor 中
4. **更易扩展**：添加新动作类型只需修改一处
5. **测试完整**：包含对状态机更新的完整测试

这种架构符合 Go 的简洁设计哲学，同时提供了强大的自动化处理能力。
