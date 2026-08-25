# TextTool 文本处理服务

纯 Go 标准库实现的文本处理后端服务，提供文档管理、文本转换、哈希计算、正则规则等功能。

## 运行说明

```bash
cd origin/
go run ./cmd/server
```

默认监听 `:8080`，可通过环境变量 `PORT` 或 `ADDR` 修改。

## API 列表

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/documents | 创建文本文档 |
| GET | /api/documents | 列取文本文档（分页） |
| GET | /api/documents/{id} | 获取文本文档 |
| PUT | /api/documents/{id} | 更新文本文档 |
| DELETE | /api/documents/{id} | 删除文本文档 |
| GET | /api/documents/{id}/stats | 文本统计（字符/单词/行数/词频） |
| GET | /api/documents/{id}/dedup | 按行去重 |
| POST | /api/tasks | 创建转换任务 |
| GET | /api/tasks | 列取转换任务（分页） |
| GET | /api/tasks/{id} | 获取转换任务 |
| POST | /api/tasks/{id}/execute | 执行转换任务 |
| DELETE | /api/tasks/{id} | 删除转换任务 |
| POST | /api/hashes | 创建哈希记录 |
| GET | /api/hashes | 列取哈希记录（分页） |
| GET | /api/hashes/{id} | 获取哈希记录 |
| POST | /api/hashes/compute | 计算文档哈希 |
| POST | /api/hashes/verify | 验证文档哈希 |
| DELETE | /api/hashes/{id} | 删除哈希记录 |
| POST | /api/rules | 创建正则规则 |
| GET | /api/rules | 列取正则规则（分页） |
| GET | /api/rules/{id} | 获取正则规则 |
| PUT | /api/rules/{id} | 更新正则规则 |
| DELETE | /api/rules/{id} | 删除正则规则 |
| POST | /api/rules/{id}/match | 使用规则匹配文本 |
| POST | /api/rules/{id}/replace | 使用规则替换文本 |
| POST | /api/rules/match | 直接使用模式匹配 |
| POST | /api/rules/replace | 直接使用模式替换 |
