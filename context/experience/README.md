# 经验知识库

常见错误、有效的测试模式、调试技巧和工作流优化

**最后更新**: 2026-02-01

## 领域概览

| 领域 | 描述 | 关键内容 |
|------|------|---------|
| [错误模式](error-patterns.md) | 常见错误和解决方案 | Go 错误、LLM 集成、数据库问题 |
| [测试模式](test-patterns.md) | 有效的测试模式 | Mock 策略、断言技巧、测试组织 |
| [调试技巧](debugging-tactics.md) | 调试策略和工具 | 日志分析、性能分析、追踪 |
| [工作流优化](workflow-improvements.md) | 流程优化和最佳实践 | 开发流程、协作方式 |

## 快速导航

### 遇到问题？
- **程序崩溃或异常** → [错误模式](error-patterns.md)
- **测试失败** → [测试模式](test-patterns.md) 和 [错误模式](error-patterns.md)
- **性能问题** → [调试技巧](debugging-tactics.md)
- **工作效率问题** → [工作流优化](workflow-improvements.md)

### 想改进？
- **提升测试质量** → [测试模式](test-patterns.md)
- **加快开发速度** → [工作流优化](workflow-improvements.md)
- **优化性能** → [调试技巧](debugging-tactics.md)

## 常见问题速查

| 问题 | 位置 |
|------|------|
| Hertz 路由返回 404 | [错误模式 - 路由问题](error-patterns.md#路由问题) |
| 数据库连接超时 | [错误模式 - 数据库问题](error-patterns.md#数据库问题) |
| LLM 请求超时 | [错误模式 - AI 集成问题](error-patterns.md#ai集成问题) |
| 单元测试 Mock 不生效 | [测试模式 - Mock 策略](test-patterns.md#mock策略) |
| 内存泄漏 | [调试技巧 - 性能分析](debugging-tactics.md#性能分析) |
| Git 工作流混乱 | [工作流优化 - 版本控制](workflow-improvements.md#版本控制) |

## 使用示例

### 查找错误解决方案

**问题**: "context deadline exceeded"

**查找步骤**:
1. 查看 [错误模式 - 超时问题](error-patterns.md#超时问题)
2. 识别可能的原因
3. 应用解决方案

### 改进测试

**目标**: 提升单元测试质量

**步骤**:
1. 查看 [测试模式 - 测试组织](test-patterns.md#测试组织)
2. 查看 [测试模式 - Mock 策略](test-patterns.md#mock策略)
3. 应用到现有测试中

---

**维护者**: GitHub Copilot  
**版本**: 1.0  
**最后更新**: 2026-02-01
