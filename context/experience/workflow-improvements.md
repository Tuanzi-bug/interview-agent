# 工作流优化

**最后更新**: 2026-02-01  
**维护者**: 技术团队

## 概述

开发流程优化建议和最佳实践。

## 版本控制

### Git 工作流

```bash
# 功能分支开发
git checkout -b feature/resume-analyzer
# 进行开发...
git add .
git commit -m "feat: implement resume analyzer agent"

# 推送并创建 PR
git push origin feature/resume-analyzer

# 合并到主分支
git checkout main
git pull origin main
git merge feature/resume-analyzer
git push origin main
```

### Commit 规范

```
<type>(<scope>): <subject>

<body>

<footer>
```

**类型**:
- `feat`: 新功能
- `fix`: 错误修复
- `docs`: 文档
- `style`: 格式
- `refactor`: 重构
- `perf`: 性能优化
- `test`: 测试

**示例**:
```
feat(interview): implement question generation agent

Implement question generation using Eino framework.
- Support multiple interview types (comprehensive, specialized)
- Integrate with resume context
- Add difficulty levels

Closes #123
```

---

## 开发环境

### 快速启动

```bash
# 一键启动开发环境
make run

# 查看日志
docker-compose logs -f backend

# 清理
make clean
```

### 数据库迁移

```bash
# 初始化数据库
make db-migrate

# 查看当前状态
mysql -h localhost -u root -proot -e "show databases;"
```

---

## 代码质量

### 代码审查检查清单

审查代码时检查以下方面：

**功能性**:
- [ ] 代码是否实现了需求功能
- [ ] 是否有明显的逻辑错误
- [ ] 是否处理了边界情况

**测试**:
- [ ] 是否包含相关测试
- [ ] 测试覆盖率是否足够
- [ ] 测试是否有意义（不只是覆盖率）

**性能**:
- [ ] 是否有 N+1 查询问题
- [ ] 是否有不必要的循环
- [ ] 是否正确使用了缓存

**安全性**:
- [ ] 是否进行了输入验证
- [ ] 是否有 SQL 注入风险（应使用 GORM）
- [ ] 认证/授权是否正确

**代码风格**:
- [ ] 是否遵循项目编码风格
- [ ] 变量名是否清晰
- [ ] 函数是否太长（>50 行）

---

## 文档

### 更新进度时的文档

每当实现新功能或发现新的最佳实践时，更新相关文档：

```bash
# 发现新的错误模式
/learn We found that goroutines leak when not properly canceled with context

# 发现新的测试模式
/learn Table-driven tests work really well for multiple scenarios in Go

# 实现新功能
/learn Implemented resume analyzer agent using Eino framework with RAG support
```

---

## 团队协作

### 异步沟通

- 使用清晰的 commit 消息
- 在 PR 中详细描述变更
- 更新相关文档

### 同步讨论

- 问题澄清在 Issue 中记录
- 决策记录在 ADR (Architecture Decision Record) 中
- 定期同步会议讨论架构和规划

---

## 发布流程

### 版本标记

```bash
# 创建版本标签
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0

# 列出标签
git tag -l
```

### 发布清单

- [ ] 所有测试通过
- [ ] 代码审查完毕
- [ ] 文档已更新
- [ ] 版本号已更新
- [ ] CHANGELOG 已更新
- [ ] 标签已推送

---

## 常见优化建议

### 减少 DB 查询

```go
// ❌ 问题：N+1 查询
for _, interview := range interviews {
    user := db.First(&User{}, interview.UserID)  // 每个面试都查一次
    fmt.Println(user.Email)
}

// ✅ 解决：预加载
var interviews []Interview
db.Preload("User").Find(&interviews)
for _, interview := range interviews {
    fmt.Println(interview.User.Email)  // 一条查询
}
```

### 优化缓存

```go
// 使用缓存避免重复计算
cache := NewLRUCache(1000)

func GetUserProfile(userID string) *Profile {
    if cached, ok := cache.Get(userID); ok {
        return cached.(*Profile)
    }
    
    profile := computeProfile(userID)
    cache.Set(userID, profile)
    return profile
}
```

---

**维护者**: 技术团队  
**版本**: 1.0  
**最后更新**: 2026-02-01
