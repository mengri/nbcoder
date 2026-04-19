# 克隆池领域（ClonePool Domain）需求说明

## 领域职责
- 负责 Git 仓库克隆实例的管理、分配、回收。
- 支持多仓库、多实例、状态追踪。

## 主要实体
- CloneInstance（克隆实例）
- Repository（仓库）
- ClonePool（克隆池）

## 主要功能需求
1. 支持按仓库独立管理克隆池，动态扩容与回收。
2. 克隆实例状态：idle、busy、dirty。
3. 支持任务获取/归还克隆实例，自动 clean 校验。
4. 支持异常恢复与崩溃后状态标记。
5. 领域事件：CloneAcquired、CloneReleased、CloneBecameDirty。

## 领域事件
- CloneAcquired：克隆被占用
- CloneReleased：克隆被释放
- CloneBecameDirty：克隆产生未提交变更

## 限界上下文
- 仅管理克隆实例，不直接操作代码内容、任务调度。
- 与 Agent 领域通过端口和事件解耦。
