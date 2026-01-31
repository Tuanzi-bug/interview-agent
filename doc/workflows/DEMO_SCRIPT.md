# `/req-dev` 工作流演示脚本

这个脚本展示如何使用 `/req-dev` 统一工作流从需求到部署。

## 演示 1: 完整工作流 - 添加导出功能

### 场景
产品经理提出需求："我们需要让用户能够导出他们的面试记录为 PDF"

### 命令
```bash
/req-dev "添加导出面试记录为PDF的功能"
```

### 预期流程

#### Phase 1: Discovery (SE: Product Manager)
```
📋 PHASE 1: Requirements Discovery
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Agent: SE: Product Manager

让我澄清一下需求：

1. **谁需要导出面试记录？**
   - 所有用户？还是只有付费用户？
   
2. **需要导出哪些信息？**
   - 面试问题和答案？
   - 评分和反馈？
   - 时间戳和统计？
   
3. **PDF 格式要求？**
   - 简单文本还是需要格式化？
   - 是否需要品牌 logo？
   
4. **有大小限制吗？**
   - 单次面试还是批量导出？
   
5. **成功指标是什么？**
   - 预期多少用户使用这个功能？
```

**用户回答**:
```
1. 所有用户都可以导出
2. 需要完整的面试问题、用户答案、AI评分和建议
3. 需要格式化的PDF，包含平台logo
4. 单次面试导出，大小控制在2MB以内
5. 目标是30%的活跃用户在完成面试后导出
```

**Agent 继续**:
```
✅ 明白了！创建 GitHub Epic...

GitHub Epic #123 已创建
- 标题: [Feature] 导出面试记录为PDF
- 标签: enhancement, size: medium, frontend, backend
- 子任务:
  - #124: 前端导出按钮UI
  - #125: 后端PDF生成API
  - #126: PDF模板设计

用户故事:
- 作为用户，我想导出面试记录为PDF，以便离线查看和存档
- 作为用户，我希望PDF格式美观，包含完整信息

验收标准:
- [ ] 用户可以点击"导出PDF"按钮
- [ ] PDF包含所有面试Q&A和评分
- [ ] PDF大小 < 2MB
- [ ] 生成时间 < 5秒
- [ ] 30%的用户在完成面试后使用导出功能

Quality Gates: ✓ Passed
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

#### Phase 2: UX Design (SE: UX Designer)
```
🎨 PHASE 2: User Experience Design
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Agent: SE: UX Designer

用户旅程分析:

场景: 用户完成面试后想要导出记录

Jobs-to-be-Done:
- 用户想要保存面试记录以便日后复习
- 用户想要分享给导师或朋友寻求建议
- 用户想要建立个人面试档案

用户旅程:
1. 用户完成面试 → 查看评分和反馈
2. 发现"导出PDF"按钮 → 点击
3. 系统生成PDF → 显示进度
4. 下载完成 → 打开PDF查看
5. (Optional) 分享给他人

交互设计:
- 位置: 面试结果页面右上角
- 按钮样式: 主要按钮 "导出PDF" with 下载图标
- 生成中状态: 按钮变为加载状态，显示进度
- 完成状态: 浏览器自动下载

移动端考虑:
- 按钮放在顶部工具栏
- 生成后提示用户选择保存位置

无障碍:
- [x] 键盘导航支持
- [x] 屏幕阅读器友好
- [x] 清晰的进度提示

用户旅程图: doc/ux/interview-export-journey.md

Quality Gates: ✓ Passed
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

#### Phase 3: PRD (prd skill)
```
📄 PHASE 3: Product Requirements Document
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Skill: prd

生成 PRD: doc/requirements/interview-export-pdf.md

执行摘要:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
问题: 用户无法保存和分享他们的面试记录
解决方案: 提供PDF导出功能，包含完整的面试内容和反馈
成功指标:
- 30% 用户使用导出功能（30天内）
- 95% PDF 生成成功率
- <5秒 PDF 生成时间（p95）
- 4.5/5 用户满意度评分

技术规格:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
前端:
- React 组件: ExportButton
- 状态管理: 导出进度和状态
- 下载处理: 使用 Blob 和 URL.createObjectURL

后端:
- API: POST /api/interviews/:id/export
- PDF 生成: 使用 wkhtmltopdf 或类似库
- 存储: 临时文件在 S3，TTL 1小时
- 响应: 返回预签名 URL

数据模型:
- Interview (现有)
- ExportLog (新增，记录导出历史)

API 设计:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
POST /api/interviews/:id/export
Request: { format: "pdf", includeDetails: true }
Response: { downloadUrl: "https://...", expiresAt: "..." }

非目标 (v2):
- 批量导出
- Word 格式
- 自定义模板
- Email 发送

Quality Gates: ✓ Passed
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

#### Phase 4: Architecture (SE: Architect)
```
🏗️ PHASE 4: Architecture Review
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Agent: SE: Architect

系统分析: Web 应用，中等规模

架构设计:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

组件:
1. 前端 (React)
   - ExportButton 组件
   - 下载管理器
   
2. 后端 API (Go/Hertz)
   - Export Handler
   - PDF 生成服务
   
3. PDF 生成器
   - HTML 模板引擎
   - PDF 渲染库 (wkhtmltopdf)
   
4. 存储 (S3)
   - 临时 PDF 文件存储
   - 预签名 URL 生成

数据流:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
1. 用户点击"导出" → 前端发送请求
2. 后端验证权限 → 查询面试数据
3. 渲染 HTML 模板 → 生成 PDF
4. 上传到 S3 → 生成预签名 URL
5. 返回 URL → 前端下载

可扩展性考虑:
- PDF 生成异步化（Redis 队列）
- 如果并发高，使用专用 PDF 生成服务
- S3 CDN 加速下载

性能要求:
- 生成时间: <5秒 (p95)
- 并发处理: 100 req/min
- PDF 大小: <2MB

成本估算:
- S3 存储: ~$5/月 (100k exports, 1h TTL)
- PDF 生成: 计算资源，可忽略

ADR 创建: doc/architecture/adr-interview-pdf-export.md

Quality Gates: ✓ Passed
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

#### Phase 5: Security (SE: Security)
```
🔒 PHASE 5: Security Review
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Agent: SE: Security

风险等级: Medium (用户数据，文件生成)

安全检查:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✅ OWASP A01 - 访问控制
   - 验证: 用户只能导出自己的面试记录
   - JWT 验证必需
   - Interview ID 所有权检查
   
✅ OWASP A03 - 注入攻击
   - HTML 模板: 转义所有用户输入
   - PDF 生成: 使用安全的库
   - 文件名: 严格验证，防止路径遍历
   
✅ OWASP A04 - 不安全设计
   - 速率限制: 5次导出/小时/用户
   - 文件大小限制: 2MB
   - S3 预签名 URL: 1小时过期
   
✅ OWASP A05 - 安全配置错误
   - S3 bucket: 非公开访问
   - CORS: 限制为应用域名
   - 临时文件: 自动清理
   
✅ A07 - 识别和认证失败
   - 需要有效 JWT token
   - 无匿名导出

额外安全措施:
- 水印: 在 PDF 添加用户标识水印
- 审计日志: 记录所有导出操作
- PII 保护: 确保不泄露敏感信息

安全测试用例:
- 尝试导出他人的面试记录（应失败）
- 尝试路径遍历攻击（应被阻止）
- 速率限制测试（超出应返回 429）

Quality Gates: ✓ Passed
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

#### Phase 6: Implementation
```
💻 PHASE 6: Implementation
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

任务分解:
- Task 6.1: 前端 ExportButton 组件
- Task 6.2: 后端 Export API
- Task 6.3: PDF 生成服务
- Task 6.4: S3 集成
- Task 6.5: 测试

实现中...

✅ Task 6.1: 前端组件 (15 min)
   文件: frontend/src/components/ExportButton.tsx
   - 按钮组件
   - 加载状态
   - 错误处理
   测试: 10/10 passed

✅ Task 6.2: 后端 API (20 min)
   文件: backend/api/handler/interview/export.go
   - Export handler
   - 权限验证
   - 速率限制
   测试: 12/12 passed

✅ Task 6.3: PDF 生成 (25 min)
   文件: backend/internal/service/pdf_generator.go
   - HTML 模板渲染
   - PDF 生成
   - 错误处理
   测试: 15/15 passed

✅ Task 6.4: S3 集成 (15 min)
   文件: backend/internal/service/storage.go
   - S3 上传
   - 预签名 URL 生成
   - TTL 管理
   测试: 8/8 passed

✅ Task 6.5: 集成测试 (10 min)
   - End-to-end 测试
   - 性能测试
   - 边界测试
   测试: 20/20 passed

代码统计:
- 新增文件: 8
- 代码行数: 1,456
- 测试覆盖率: 89%

Quality Gates: ✓ Passed
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

#### Phase 7: Documentation (SE: Tech Writer)
```
📝 PHASE 7: Documentation
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Agent: SE: Tech Writer

创建文档...

✅ API 文档 (10 min)
   文件: doc/api/interview-export.md
   
   # 导出面试记录 API
   
   ## Endpoint
   POST /api/interviews/:id/export
   
   ## Authentication
   需要有效的 JWT token
   
   ## Request
   ```json
   {
     "format": "pdf",
     "includeDetails": true
   }
   ```
   
   ## Response
   ```json
   {
     "downloadUrl": "https://s3.../interview-123.pdf",
     "expiresAt": "2026-01-30T15:00:00Z",
     "fileSize": 1234567
   }
   ```
   
   ## 错误码
   - 401: 未授权
   - 403: 无权限导出此面试
   - 429: 超出速率限制
   - 500: PDF 生成失败

✅ 用户指南 (10 min)
   文件: doc/user-guide/export-interview.md
   
   # 如何导出面试记录
   
   1. 完成面试后，进入面试结果页面
   2. 点击右上角的"导出PDF"按钮
   3. 等待几秒钟，系统生成 PDF
   4. PDF 自动下载到你的设备
   5. 打开 PDF 查看完整的面试记录
   
   ## 注意事项
   - 每小时最多导出 5 次
   - PDF 下载链接 1 小时后过期
   - 文件大小通常在 500KB - 2MB 之间

✅ 开发者文档 (10 min)
   文件: doc/dev/pdf-export-architecture.md
   
   # PDF 导出架构
   
   ## 组件
   - ExportButton: React 组件
   - ExportHandler: Go API handler
   - PDFGenerator: PDF 生成服务
   - StorageService: S3 集成
   
   ## 数据流
   [流程图]
   
   ## 本地开发
   1. 安装 wkhtmltopdf
   2. 配置 S3 凭证
   3. 运行测试

✅ README 更新 (3 min)
   在功能列表中添加 PDF 导出

Quality Gates: ✓ Passed
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

#### Phase 8: CI/CD (SE: DevOps/CI)
```
🚀 PHASE 8: CI/CD & Deployment
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Agent: SE: DevOps/CI

配置 CI/CD...

✅ CI Pipeline (15 min)
   文件: .github/workflows/interview-export.yml
   
   - Build: ✓ (2 min)
   - Unit tests: ✓ (3 min)
   - Integration tests: ✓ (5 min)
   - Security scan: ✓ (4 min)
   - Linting: ✓ (1 min)

✅ Staging Deployment (10 min)
   - 部署到 staging 环境
   - 运行 smoke tests
   - 手动测试导出功能
   结果: ✓ Passed

✅ 监控设置 (10 min)
   Dashboard: Export Feature Metrics
   
   指标:
   - Export requests/min
   - Export success rate (target: >95%)
   - PDF generation time (p95 target: <5s)
   - Error rate by type
   - S3 storage usage
   
   告警:
   - Success rate < 90% → Critical
   - Generation time > 10s (p95) → Warning
   - Error rate > 5% → Warning

✅ Runbook (5 min)
   文件: doc/runbooks/interview-export-troubleshooting.md
   
   # PDF 导出故障排查
   
   ## 常见问题
   
   ### PDF 生成失败
   - 检查 wkhtmltopdf 服务状态
   - 查看错误日志
   - 验证模板文件
   
   ### S3 上传失败
   - 检查 AWS 凭证
   - 验证 S3 bucket 权限
   - 查看网络连接
   
   ### 性能问题
   - 检查并发量
   - 查看 PDF 大小
   - 考虑异步队列

Quality Gates: ✓ Passed
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

#### Phase 9: Final Review
```
✅ PHASE 9: Final Review & Launch
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Launch Checklist:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

✓ 所有测试通过
  - 单元测试: 65/65 passed (89% coverage)
  - 集成测试: 20/20 passed
  - E2E 测试: 5/5 passed

✓ 安全审查完成
  - OWASP Top 10: 全部通过
  - 权限验证: 已实现
  - 速率限制: 已配置

✓ 文档完整
  - API 文档: ✓
  - 用户指南: ✓
  - 开发者文档: ✓
  - Runbook: ✓

✓ 监控激活
  - Grafana Dashboard: ✓
  - 告警规则: ✓
  - 日志收集: ✓

✓ Staging 测试
  - 功能测试: ✓
  - 性能测试: ✓ (4.2s p95)
  - 负载测试: ✓ (150 req/min)

✓ Stakeholder 审批
  - 产品经理: ✓ Approved
  - 技术 Leader: ✓ Approved

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

🎉 准备生产部署！

总结:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
- GitHub Epic: #123 (Ready to deploy)
- PRD: doc/requirements/interview-export-pdf.md
- Architecture: doc/architecture/adr-interview-pdf-export.md
- Code: 8 files, 1,456 insertions
- Tests: 90/90 passed, 89% coverage
- Documentation: 4 docs
- CI/CD: ✓ Configured
- 预计部署时间: 10 分钟

部署计划:
1. 部署到生产环境
2. 监控 30 分钟
3. 逐步开放给 10% → 50% → 100% 用户
4. 收集用户反馈

继续部署? (y/n)
```

---

## 总时间估算

| 阶段 | 时间 |
|------|------|
| Phase 1: Discovery | 10 分钟 |
| Phase 2: UX Design | 15 分钟 |
| Phase 3: PRD | 10 分钟 |
| Phase 4: Architecture | 15 分钟 |
| Phase 5: Security | 10 分钟 |
| Phase 6: Implementation | 85 分钟 |
| Phase 7: Documentation | 33 分钟 |
| Phase 8: CI/CD | 40 分钟 |
| Phase 9: Final Review | 10 分钟 |
| **总计** | **228 分钟 (3.8 小时)** |

---

## 关键要点

1. **完整性**: 从需求到部署的所有步骤都被覆盖
2. **质量**: 每个阶段都有质量门控确保标准
3. **自动化**: Agent 自动协调，减少人工工作
4. **文档化**: 自动生成全面的文档
5. **可追溯**: GitHub Issue 链接所有工作
6. **安全**: 安全审查是必需阶段
7. **监控**: 自动配置监控和告警

---

## 下一步

这个演示展示了一个中等复杂度功能的完整流程。

**尝试其他场景**:
- 简单 UI 变更: `/req-dev --fast "修改按钮颜色"`
- 紧急修复: `/req-dev --mode=hotfix "修复登录 bug"`
- 大型功能: `/req-dev "构建 AI 聊天机器人"`

**查看更多示例**:
- [req-dev-examples.md](req-dev-examples.md)
