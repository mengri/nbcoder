# Git 领域（Git Domain）需求说明

## 领域职责
- 负责仓库抽象、PR/MR 操作、分支策略管理。
- 支持分支、提交、PR/MR 创建与追溯。

## 主要实体
- PullRequest（PR/MR）
- Branch（分支）
- CommitHistory（提交历史）

## 主要功能需求
1. 支持分支策略配置、分支创建、命名规范。
2. 支持 PR/MR 创建、描述生成、变更摘要。
3. 支持提交历史追溯、增量 commit、squash 合并。
4. 不自动合并 PR/MR，人工审核后合并。

## 限界上下文
- 仅管理仓库与 PR/MR，不直接操作代码内容、任务调度。
- 与 Project、Agent 领域通过端口和事件解耦。
