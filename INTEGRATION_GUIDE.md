# 动作处理集成指南

## 架构概述

本项目采用**中间件驱动**的架构，将动作处理逻辑整合到现有业务流程中，而不是创建独立的上报路由。系统包含两个独立的动作处理器：猫咪动作处理器和站点动作处理器。

## 核心组件

### 1. ActionProcessor (猫咪动作处理器)

位置：`middleware/action_processor.go`

**职责**：
- 处理所有 CatAction 创建请求
- 自动记录动作到数据库
- 根据动作类型自动更新猫咪状态机 (FSM)
- 解析动作详情中的数据（如体温、体重）

**主要方法**：

```go
// ProcessAction 处理动作并更新状态机
func (p *ActionProcessor) ProcessAction(action *model.CatAction) (*model.CatFSM, error)
```

### 2. SiteActionProcessor (站点动作处理器)

位置：`middleware/site_action_processor.go`

**职责**：
- 处理所有 SiteAction 创建请求
- 自动记录动作到数据库
- 根据动作类型自动更新站点状态机 (FSM)

**主要方法**：

```go
// ProcessAction 处理动作并更新状态机
func (p *SiteActionProcessor) ProcessAction(action *model.SiteAction) (*model.SiteFSM, error)
```

### 3. CatActionController (猫咪动作控制器)

已集成 ActionProcessor，在创建动作时自动触发状态机更新。

### 4. SiteActionController (站点动作控制器)

已集成 SiteActionProcessor，在创建动作时自动触发状态机更新。

## 使用方式

### 猫咪动作

创建任何 CatAction 都会自动触发状态机更新：

```bash
curl -X POST http://localhost:5100/cat-actions \
  -H "Content-Type: application/json" \
  -d '{
    "cat_id": 1,
    "site_id": 1,
    "user_id": 1,
    "action_type": "测体温",
    "action_detail": "{\"temperature_c\": 39.5}"
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
    "action_detail": "{\"temperature_c\": 39.5}",
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

### 站点动作

创建任何 SiteAction 都会自动触发站点状态机更新：

```bash
curl -X POST http://localhost:5100/site-actions \
  -H "Content-Type: application/json" \
  -H "X-User-ID: 1" \
  -d '{
    "site_id": 1,
    "action_type": "喂食",
    "action_detail": "{\"food_type\": \"猫粮\", \"amount\": \"100g\"}"
  }'
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

### 支持的动作类型

#### 猫咪动作

| 动作类型 | ActionDetail 格式 | FSM 字段 | 说明 |
|---------|-----------------|---------|------|
| 测体温 | `{"temperature_c": 39.5}` | `TemperatureC` | 测量体温后自动更新 |
| 称重 | `{"weight_kg": 5.2}` | `WeightKG` | 体重记录自动更新 |
| 修剪指甲 | `{"notes": "xxx"}` | `TrimNailsTime` | 自动记录为当前时间 |
| 绝育 | `{"notes": "xxx"}` | - | 仅记录 |
| 驱虫 | `{"drug_name": "xxx", "dosage": "xxx"}` | - | 仅记录 |
| 疫苗 | `{"vaccine_name": "xxx", "batch_no": "xxx"}` | - | 仅记录 |
| 洗澡 | `{"notes": "xxx"}` | - | 仅记录 |

#### 站点动作

| 动作类型 | ActionDetail 格式 | FSM 字段 | 说明 |
|---------|-----------------|---------|------|
| 消毒 | `{"disinfectant": "xxx", "notes": "xxx"}` | `LastDisinfectTime` | 自动记录为当前时间 |
| 喂食 | `{"food_type": "xxx", "amount": "xxx"}` | `LastFeedTime` | 自动记录为当前时间 |
| 喂水 | `{"water_type": "xxx"}` | `LastGiveWaterTime` | 自动记录为当前时间 |
| 逗猫 | `{"duration": 30, "notes": "xxx"}` | `LastPlayTime` | 自动记录为当前时间 |
| 清理猫砂 | `{"litter_type": "xxx"}` | `LastCleanLitter` | 自动记录为当前时间 |

## 集成流程

### 猫咪动作流程

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
    - 称重 → 解析体重 → 更新 FSM.WeightKG
    - 其他 → 仅记录动作
         ↓
    返回动作和更新后的状态机
```

### 站点动作流程

```
用户请求 → SiteActionController.CreateSiteAction
         ↓
    SiteActionProcessor.ProcessAction
         ↓
    1. 记录动作到数据库
         ↓
    2. 根据动作类型处理
         ↓
    - 消毒 → 更新 FSM.LastDisinfectTime
    - 喂食 → 更新 FSM.LastFeedTime
    - 喂水 → 更新 FSM.LastGiveWaterTime
    - 逗猫 → 更新 FSM.LastPlayTime
    - 清理猫砂 → 更新 FSM.LastCleanLitter
         ↓
    返回动作和更新后的状态机
```

## 扩展指南

### 添加新的猫咪动作类型

1. 在 `model/cat_event.go` 中添加新的 `CatActionType` 常量
2. 在 `middleware/validation.go` 的 `validateCatActionType` 中添加验证
3. 在 `middleware/action_processor.go` 中添加处理逻辑
4. 如需更新 FSM，添加对应的 repository 方法
5. 编写测试验证功能

### 添加新的站点动作类型

1. 在 `model/site_action.go` 中添加新的 `SiteActionType` 常量
2. 在 `middleware/validation.go` 的 `validateSiteActionType` 中添加验证
3. 在 `middleware/site_action_processor.go` 中添加处理逻辑
4. 如需更新 FSM，添加对应的 repository 方法
5. 编写测试验证功能

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

## 最佳实践

### 1. 使用标准接口

创建动作使用统一的 `POST /cat-actions` 和 `POST /site-actions` 接口，而不是特定数据的接口。

### 2. 数据格式

将关键数据（如体温、体重）使用 JSON 格式放在 `action_detail` 中，便于自动解析。

### 3. 错误处理

状态机更新失败时会自动回滚动作记录，确保数据一致性。

### 4. 日志记录

所有动作处理都会记录日志，便于追踪和调试。

## 总结

新架构的优势：

1. **更简洁**：删除了不必要的独立路由和服务层
2. **更集成**：动作处理完全整合到现有业务流程
3. **更易维护**：逻辑集中在 ActionProcessor 和 SiteActionProcessor 中
4. **更易扩展**：添加新动作类型只需修改一处
5. **测试完整**：包含对状态机更新的完整测试

这种架构符合 Go 的简洁设计哲学，同时提供了强大的自动化处理能力。
