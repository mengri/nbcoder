# NBCoder Copilot Skill 配置说明

本目录用于存放本项目自定义的 Copilot Skill，每个子目录为一个独立 Skill。

## 目录结构示例

```
.NBCoder/skills/
  ├── api-design/
  │   ├── SKILL.md
  │   ├── instructions.md
  │   └── references/
  └── db-design/
      ├── SKILL.md
      └── instructions.md
```

## 配置步骤
1. 每个 Skill 建立独立子目录（如 api-design/、db-design/ 等）。
2. 在每个 Skill 目录下，创建 SKILL.md，内容包括：
   - Skill 名称与描述
   - 输入/输出说明
   - 工具权限（如 read_file、write_file 的目录范围）
   - 推荐模型（如 primary: gpt-4o）
3. 可选：补充 instructions.md（详细指令）、references/（参考文档）等。
4. 在 Agent 配置或调用时，指定可用的 skill 路径和权限。

## 示例 SKILL.md

```markdown
# Skill: API 设计

## 描述
生成 RESTful API 设计与接口文档

## 输入
- 业务需求描述
- 数据模型定义

## 输出
- API 接口定义
- 请求/响应示例

## 工具权限
- read_file: knowledge-base/api/
- write_file: cards/{REQ-ID}/design.md

## 模型推荐
- primary: gpt-4o
```
