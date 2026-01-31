---
name: 'code-review'
description: 'Comprehensive code review workflow combining security, quality, architecture, and responsible AI reviews with context-aware loading and focused documentation'
agent: 'agent'
argument-hint: 'Files/PR to review or review focus area (security, quality, architecture, ai-ethics)'
tools: ['execute', 'read', 'edit', 'search', 'agent', 'todo']
---

# Code Review Workflow

Execute professional code review following industry best practices with multi-dimensional analysis: security vulnerabilities, code quality, architecture soundness, and responsible AI principles.

## Mission

Perform systematic, comprehensive code review that prevents production issues through layered analysis: OWASP security checks, SOLID principles validation, Well-Architected framework assessment, and bias/accessibility verification. Load only relevant context and update core decision logs without documentation pollution.

## Scope & Preconditions

### Prerequisites
- **Code to review**: ${input:target:Files, PR number, or commit hash to review}
- **Review focus** (optional): ${input:focus:security|quality|architecture|ai-ethics|all}
- **Context accessible**: `context/` and `requirements/` directories
- **Tests passing**: Ideally complete `/tdd` workflow first

### What This Workflow Covers
1. ✅ Security vulnerability detection (OWASP Top 10, LLM security)
2. ✅ Code quality assessment (SOLID, design patterns, complexity)
3. ✅ Architecture validation (Well-Architected frameworks)
4. ✅ Responsible AI review (bias, accessibility, privacy)
5. ✅ Performance and scalability analysis
6. ✅ Documentation and maintainability check

### What This Does NOT Cover
- ❌ Functional testing (use `/tdd` or `/e2e`)
- ❌ Build configuration fixes (use `/build-fix`)
- ❌ Business logic design (use `/plan`)

## Context Loading Strategy

**Use the `context-loader` skill for intelligent context loading:**

```
@context-loader
```

This skill automatically implements the index-first strategy for code review:

1. ✅ Loads `context/README.md` (~500 tokens) - lightweight index
2. ✅ Selectively loads review-relevant knowledge:
   - `context/tech/architecture.md` (架构审查时)
   - `context/experience/workflow-improvements.md` (审查最佳实践)
   - `context/experience/error-patterns.md` (已知问题库)
3. ✅ Loads current requirements as needed:
   - `requirements/INDEX.md` (需求索引)
   - `requirements/in-progress/当前需求.md` (if related)
4. ✅ Loads security/compliance context when needed:
   - `context/tech/security-standards.md` (安全标准)
   - `context/business/compliance.md` (合规要求)

**Key Benefits**: 
- 80-85% reduction in context token usage
- Only loads relevant knowledge for the specific review type
- Prevents context pollution while maintaining quality

See [context-loader skill](../.github/skills/context-loader/SKILL.md) for details.

## Workflow

### Phase 0: Review Initialization & Planning

**Objective**: Understand review scope and create targeted review plan

**Steps**:

1. **Load Context Index**
   ```
   Read: context/README.md (project overview)
   Read: requirements/INDEX.md (if exists)
   ```

2. **Analyze Review Target**
   
   Identify what you're reviewing:
   
   **By Code Type**:
   - 🌐 **Web API/Backend** → Security (OWASP), performance, scalability
   - 🤖 **AI/ML Integration** → OWASP LLM Top 10, bias, explainability
   - 🎨 **Frontend/UI** → Accessibility, responsive design, UX
   - 📊 **Data Processing** → Privacy, efficiency, error handling
   - 🔐 **Authentication** → Crypto, access control, session management
   - ⚙️ **Infrastructure** → Configuration security, secrets management
   
   **By Risk Level**:
   - 🔴 **High Risk**: Payment, authentication, AI decisions, admin functions, PII
   - 🟡 **Medium Risk**: User data, external APIs, reporting
   - 🟢 **Low Risk**: UI components, utilities, display logic

3. **Create Targeted Review Plan**
   
   Based on ${input:focus} and code analysis, select 3-5 most relevant areas:
   
   ```markdown
   ## Review Plan for ${input:target}
   
   **Code Type**: [API/AI/Frontend/Data/Auth/Infra]
   **Risk Level**: [High/Medium/Low]
   **Primary Concerns**: [Security/Quality/Architecture/AI-Ethics]
   
   **Selected Review Areas** (3-5 max):
   - [ ] Security: [OWASP categories relevant to this code]
   - [ ] Quality: [Specific quality concerns]
   - [ ] Architecture: [Architecture aspects to validate]
   - [ ] Responsible AI: [If AI/ML code]
   - [ ] Performance: [If performance-critical]
   ```

4. **Load Domain-Specific Context**
   ```
   Based on review plan, load ONLY relevant files:
   - Security review → context/tech/security-standards.md
   - Architecture review → context/tech/architecture.md
   - AI review → context/experience/ai-safety-checklist.md
   ```

### Phase 1: Security Review

**Objective**: Identify security vulnerabilities using OWASP standards

**Agent**: Use `SE: Security` subagent for deep security analysis

**Steps**:

1. **Load Security Context** (if not already loaded)
   ```
   Read: context/tech/security-standards.md (if exists)
   Read: context/experience/security-incidents.md (if exists)
   ```

2. **OWASP Top 10 Checks** (Standard Applications)
   
   **A01 - Broken Access Control**:
   ```
   ❌ Check for:
   - Missing authorization checks
   - Insecure direct object references
   - Path traversal vulnerabilities
   - Privilege escalation risks
   
   ✅ Verify:
   - Authentication required for protected resources
   - Authorization checks on every sensitive operation
   - User can only access their own data
   - Admin functions properly protected
   ```
   
   **A02 - Cryptographic Failures**:
   ```
   ❌ Check for:
   - Weak hashing algorithms (MD5, SHA1)
   - Hard-coded secrets/passwords
   - Unencrypted sensitive data
   - Weak random number generation
   
   ✅ Verify:
   - Strong algorithms (bcrypt, scrypt, Argon2)
   - Secrets in environment variables/key vault
   - Encryption for sensitive data at rest/transit
   - Cryptographically secure random generation
   ```
   
   **A03 - Injection**:
   ```
   ❌ Check for:
   - SQL concatenation instead of parameterization
   - Command injection via user input
   - XSS vulnerabilities
   - Template injection
   
   ✅ Verify:
   - Parameterized queries or ORM
   - Input validation and sanitization
   - Output encoding
   - No eval() or unsafe exec()
   ```
   
   Continue through all OWASP categories relevant to code type.

3. **OWASP LLM Top 10 Checks** (AI/ML Systems)
   
   **LLM01 - Prompt Injection**:
   ```
   ❌ Check for:
   - User input directly interpolated into prompts
   - No prompt delimiters or separators
   - System instructions exposed to users
   
   ✅ Verify:
   - User input sanitized before prompt construction
   - Clear delimiters between system/user content
   - Prompt hardening techniques applied
   ```
   
   **LLM02 - Insecure Output Handling**:
   ```
   ❌ Check for:
   - LLM output used directly in SQL/commands
   - No validation of model responses
   - Trusting model output as safe
   
   ✅ Verify:
   - Output validation before use
   - Sanitization for downstream systems
   - Rate limiting and monitoring
   ```
   
   Apply remaining LLM security checks as relevant.

4. **Document Security Findings**
   ```
   Create: review-results/security-findings.md
   
   Structure:
   - 🔴 Critical: [Immediate action required]
   - 🟡 High: [Fix before merge]
   - 🟢 Medium: [Fix in follow-up]
   - ℹ️ Low: [Technical debt/best practice]
   
   For each finding:
   - Location (file:line)
   - Vulnerability type
   - Risk assessment
   - Recommended fix (with code example)
   - References (CWE, OWASP)
   ```

### Phase 2: Code Quality Review

**Objective**: Assess code maintainability, readability, and adherence to best practices

**Agent**: Use `TDD Refactor Phase - Improve Quality & Security` subagent for refactoring suggestions

**Steps**:

1. **Load Quality Context** (if not already loaded)
   ```
   Read: context/experience/code-review-checklist.md (if exists)
   Read: context/design/patterns.md (if exists)
   ```

2. **SOLID Principles Check**
   
   **Single Responsibility**:
   ```
   ❌ Red flags:
   - Class doing too many things (>5 responsibilities)
   - Methods over 50 lines
   - Classes over 300 lines
   
   ✅ Good signs:
   - Clear, focused purpose
   - One reason to change
   - Easy to name and describe
   ```
   
   **Open/Closed**:
   ```
   ❌ Red flags:
   - Modification required for new features
   - Long if/else or switch statements
   - Hard-coded types or behaviors
   
   ✅ Good signs:
   - Extension through inheritance/composition
   - Strategy pattern for variations
   - Plugin architecture
   ```
   
   **Liskov Substitution**:
   ```
   ❌ Red flags:
   - Subclass changing parent behavior unexpectedly
   - NotImplementedException in override
   - Type checking before casting
   
   ✅ Good signs:
   - Substitutable implementations
   - Interface segregation
   - Behavioral consistency
   ```
   
   Apply remaining SOLID principles.

3. **Code Complexity Analysis**
   
   **Cyclomatic Complexity**:
   ```
   Target: <10 per method
   Warning: >15 (refactor recommended)
   Critical: >20 (refactor required)
   
   Check for:
   - Deeply nested conditions (>3 levels)
   - Long methods (>50 lines)
   - Multiple exit points
   ```
   
   **Cognitive Complexity**:
   ```
   Evaluate readability:
   - Can junior developer understand?
   - Are variable names self-documenting?
   - Is control flow clear?
   ```

4. **Code Smell Detection**
   
   **Common Smells**:
   - 💩 **Duplicate Code**: Same logic in multiple places
   - 💩 **Long Method**: >50 lines, doing too much
   - 💩 **Large Class**: >300 lines, god object
   - 💩 **Long Parameter List**: >4 parameters
   - 💩 **Primitive Obsession**: Should use domain objects
   - 💩 **Switch Statements**: Should use polymorphism
   - 💩 **Temporary Field**: Fields only used sometimes
   - 💩 **Magic Numbers**: Hard-coded values without meaning

5. **Design Patterns Review**
   ```
   Identify opportunities for:
   - Factory: Complex object creation
   - Strategy: Variable algorithms
   - Repository: Data access abstraction
   - Dependency Injection: Loose coupling
   - Observer: Event handling
   - Decorator: Adding behavior dynamically
   ```

6. **Document Quality Findings**
   ```
   Update: review-results/quality-findings.md
   
   Categories:
   - Architecture Violations
   - Code Smells
   - Complexity Issues
   - Pattern Opportunities
   - Refactoring Suggestions
   ```

### Phase 3: Architecture Review

**Objective**: Validate architectural soundness against Well-Architected principles

**Agent**: Use `SE: Architect` subagent for architecture analysis

**Steps**:

1. **Load Architecture Context**
   ```
   Read: context/tech/architecture.md (if exists)
   Read: context/design/architecture.md (if exists)
   ```

2. **Well-Architected Framework Check**
   
   **Operational Excellence**:
   ```
   ✅ Check for:
   - Structured logging with correlation IDs
   - Health check endpoints
   - Metrics and monitoring hooks
   - Feature flags for gradual rollout
   - Operational runbooks referenced
   ```
   
   **Security** (already covered in Phase 1):
   ```
   ✅ Verify again at architecture level:
   - Defense in depth
   - Zero trust principles
   - Least privilege access
   - Secrets management
   ```
   
   **Reliability**:
   ```
   ✅ Check for:
   - Error handling and retry logic
   - Circuit breaker pattern for external calls
   - Graceful degradation
   - Idempotency for critical operations
   - Timeout configuration
   ```
   
   **Performance Efficiency**:
   ```
   ✅ Check for:
   - Efficient database queries (N+1 problems)
   - Caching strategy for frequently accessed data
   - Async/await for I/O operations
   - Resource pooling (connections, threads)
   - Pagination for large datasets
   ```
   
   **Cost Optimization**:
   ```
   ✅ Check for:
   - Resource cleanup (disposal patterns)
   - Efficient algorithms (time/space complexity)
   - Appropriate data structures
   - Unnecessary allocations
   ```

3. **Scalability Assessment**
   ```
   Evaluate:
   - Horizontal scalability (can add more instances?)
   - Statelessness (or proper state management)
   - Database connection handling
   - Message queue integration for async work
   - Rate limiting for API endpoints
   ```

4. **Dependency Analysis**
   ```
   Check:
   - Coupling between components
   - Circular dependencies
   - Dependency injection usage
   - Third-party library versions
   - Unused dependencies
   ```

5. **Document Architecture Findings**
   ```
   Update: review-results/architecture-findings.md
   
   Sections:
   - Well-Architected Pillar Assessment
   - Scalability Concerns
   - Dependency Issues
   - Technical Debt
   ```

### Phase 4: Responsible AI Review (If Applicable)

**Objective**: Ensure AI systems are ethical, fair, and accessible

**Agent**: Use `SE: Responsible AI` subagent for bias and ethics analysis

**Condition**: Only execute if code involves AI/ML or impacts diverse users

**Steps**:

1. **Load AI Ethics Context**
   ```
   Read: context/experience/ai-safety-checklist.md (if exists)
   Read: context/business/compliance.md (if GDPR/privacy relevant)
   ```

2. **Bias & Fairness Check**
   
   **Test with Diverse Inputs**:
   ```python
   # Names from different cultures
   test_names = [
       "John Smith", "José García", "Lakshmi Patel", 
       "Ahmed Hassan", "李明", "O'Brien", "X Æ A-12"
   ]
   
   # Age diversity
   test_ages = [18, 25, 45, 65, 75]
   
   # Edge cases
   test_edge = ["", "  ", "NULL", "<script>", "'; DROP TABLE--"]
   ```
   
   **Red Flags**:
   - Different outcomes for same qualifications + different names
   - Age discrimination (unless legally required)
   - System fails with non-English characters
   - No explanation for AI decisions
   - Training data bias not addressed

3. **Accessibility Quick Check** (User-Facing Code)
   
   **Keyboard Navigation**:
   ```html
   ✅ Good:
   <button onClick={submit}>Submit</button>
   
   ❌ Bad:
   <div onClick={submit}>Submit</div>  <!-- Can't tab to it -->
   ```
   
   **Screen Reader Support**:
   ```html
   ✅ Good:
   <input aria-label="Search products" placeholder="Search..."/>
   <img src="chart.jpg" alt="Sales increased 25% in Q3"/>
   
   ❌ Bad:
   <input placeholder="Search..."/>  <!-- No context when empty -->
   <img src="chart.jpg"/>  <!-- No description -->
   ```
   
   **Visual Accessibility**:
   - Text contrast ratio ≥4.5:1 (normal text), ≥3:1 (large text)
   - Color not sole indicator (use icons + color)
   - Zoom to 200% without breaking layout
   - Focus indicators visible

4. **Privacy & Data Protection**
   ```
   ✅ Check for:
   - Minimal data collection (collect only what's needed)
   - Anonymization where possible
   - Clear data retention policies
   - User consent for data processing
   - Right to deletion implementation
   - Transparent data usage
   ```

5. **Explainability & Transparency**
   ```
   For AI decisions:
   - Can the system explain why it made a decision?
   - Are confidence scores exposed?
   - Is there a human appeal process?
   - Are limitations documented?
   ```

6. **Document Responsible AI Findings**
   ```
   Update: review-results/responsible-ai-findings.md
   
   Sections:
   - Bias & Fairness Issues
   - Accessibility Violations
   - Privacy Concerns
   - Explainability Gaps
   ```

### Phase 5: Performance & Best Practices

**Objective**: Identify performance issues and language-specific best practices

**Steps**:

1. **Performance Checks**
   
   **Database Queries**:
   ```
   ❌ Check for:
   - N+1 query problems
   - Missing indexes
   - SELECT * instead of specific columns
   - No query pagination
   
   ✅ Verify:
   - Eager loading for related data
   - Proper indexing strategy
   - Query optimization
   - Connection pooling
   ```
   
   **Memory Management**:
   ```
   ❌ Check for:
   - Memory leaks (event listeners not removed)
   - Large object allocations in loops
   - Unnecessary object cloning
   - No disposal of disposable resources
   
   ✅ Verify:
   - Proper using statements (C#) or try-finally
   - Object pooling for frequently created objects
   - Efficient data structures
   ```
   
   **Async/Await Usage**:
   ```
   ❌ Check for:
   - Blocking calls in async methods (.Result, .Wait())
   - Async void (except event handlers)
   - Missing ConfigureAwait(false) in libraries
   
   ✅ Verify:
   - True async all the way
   - Proper cancellation token usage
   - Task.WhenAll for parallel operations
   ```

2. **Language-Specific Best Practices**
   
   **For C#/.NET**:
   - Nullable reference types enabled
   - Modern C# features (pattern matching, records)
   - IOptions pattern for configuration
   - Structured logging (Serilog, NLog)
   
   **For Python**:
   - Type hints for public APIs
   - Context managers for resources
   - List comprehensions over loops
   - f-strings for formatting
   
   **For TypeScript**:
   - Strict mode enabled
   - No `any` types (use `unknown` or proper types)
   - Discriminated unions for type safety
   - Proper error handling
   
   **For Go**:
   - Error handling at every step
   - Defer for cleanup
   - Proper goroutine management
   - Context for cancellation

3. **Documentation & Maintainability**
   ```
   ✅ Check for:
   - Public APIs documented
   - Complex logic explained
   - TODO/FIXME with issue references
   - README updated if new features
   - Breaking changes documented
   ```

### Phase 6: Consolidate & Report

**Objective**: Synthesize findings and create actionable review report

**Steps**:

1. **Consolidate All Findings**
   ```
   Merge:
   - review-results/security-findings.md
   - review-results/quality-findings.md
   - review-results/architecture-findings.md
   - review-results/responsible-ai-findings.md (if applicable)
   ```

2. **Prioritize Issues**
   
   **Priority Matrix**:
   ```
   🔴 P0 - Critical (Block Merge):
   - Security vulnerabilities (High/Critical)
   - Data loss risks
   - Production outages
   
   🟡 P1 - High (Fix Before Merge):
   - Security issues (Medium)
   - Accessibility violations (WCAG A/AA)
   - Major code quality issues
   - Bias/discrimination risks
   
   🟢 P2 - Medium (Fix in Follow-up):
   - Code smells
   - Performance optimizations
   - Best practice violations
   - Technical debt
   
   ℹ️ P3 - Low (Nice to Have):
   - Refactoring suggestions
   - Documentation improvements
   - Minor optimizations
   ```

3. **Create Review Summary**
   ```markdown
   # Code Review Summary: ${input:target}
   
   **Reviewed By**: GitHub Copilot `/code-review`
   **Review Date**: [Current Date]
   **Review Focus**: ${input:focus}
   
   ## Executive Summary
   
   [2-3 sentence overview of overall code quality]
   
   ## Key Metrics
   
   | Category | Critical | High | Medium | Low |
   |----------|----------|------|--------|-----|
   | Security | X | X | X | X |
   | Quality | X | X | X | X |
   | Architecture | X | X | X | X |
   | Responsible AI | X | X | X | X |
   
   **Overall Recommendation**: [Approve / Request Changes / Block]
   
   ## Critical Issues (Must Fix Before Merge)
   
   1. [Issue with location and fix]
   2. [Issue with location and fix]
   
   ## High Priority Issues
   
   [List with actionable recommendations]
   
   ## Medium/Low Priority Suggestions
   
   [Collapsed section with improvement opportunities]
   
   ## Positive Observations
   
   [Highlight good practices found in the code]
   
   ## Next Steps
   
   - [ ] Fix P0 issues
   - [ ] Fix P1 issues
   - [ ] Create follow-up tickets for P2/P3
   - [ ] Request re-review after fixes
   ```

4. **Update Core Documentation** (Prevent Pollution!)
   
   **ONLY update if significant patterns discovered**:
   
   ```
   Update: context/experience/code-review-checklist.md
   
   Add ONLY:
   - New vulnerability patterns specific to this project
   - Recurring issues worth tracking
   - Project-specific review focus areas
   
   DO NOT add:
   - Detailed review logs
   - Individual finding descriptions
   - Temporary observations
   ```
   
   ```
   Update: context/experience/lessons-learned.md
   
   Add ONLY:
   - New anti-patterns discovered
   - Effective fixes for recurring problems
   - Architecture decisions that prevented issues
   
   DO NOT add:
   - Complete review transcripts
   - One-time issues
   - Generic best practices
   ```

5. **Output Review Results**
   
   Present summary to user and offer next steps:
   ```
   - Export detailed findings to file
   - Create GitHub issue for each critical finding
   - Generate automated fix suggestions
   - Schedule follow-up review
   ```

## Quality Assurance

### Review Checklist

Before completing review, verify:

- [ ] All selected review areas covered
- [ ] Critical issues have actionable fixes with code examples
- [ ] Findings prioritized by severity and impact
- [ ] Positive patterns highlighted (not just problems)
- [ ] Documentation updated only with core patterns (no pollution)
- [ ] Context files loaded were relevant (no over-loading)
- [ ] Review summary is concise and actionable

### Review Quality Metrics

Track review effectiveness:
```
- Issues found vs. issues in production (catch rate)
- False positives (issues that weren't real problems)
- Time to complete review
- Context tokens used
- Documentation updates made
```

## Integration with Other Workflows

### Before Code Review
- Complete `/tdd` workflow (unit tests passing)
- Run `/build-fix` if build errors exist
- Complete `/e2e` for integration validation

### After Code Review
- Run `/build-fix` to verify fixes don't break build
- Re-run `/tdd` to ensure tests still pass
- Use `/remember` to capture lessons learned
- Use `/memory-merger` to update project knowledge

## Best Practices

### Effective Context Loading
- **Start lightweight**: Load only index files initially
- **Load incrementally**: Add detailed context based on findings
- **Avoid pre-loading**: Don't load all possible context upfront
- **Use grep/search**: Find specific information instead of reading full files

### Prevent Documentation Pollution
- **Core only**: Update only architectural decisions, patterns, recurring issues
- **No logs**: Don't document detailed review transcripts
- **Reference**: Link to review files instead of copying content
- **Temporary**: Use `review-results/` for session-specific findings

### Token Optimization
- **Targeted agents**: Use specialized subagents for deep analysis
- **Parallel review**: If multiple files, review independently when possible
- **Pattern matching**: Check against known issue patterns first
- **Summary first**: Generate executive summary before details

### Continuous Improvement
- **Track patterns**: Identify recurring issues across reviews
- **Update checklists**: Add project-specific checks as patterns emerge
- **Share learnings**: Use `/remember` to capture valuable insights
- **Feedback loop**: Validate review effectiveness against production issues

## References

### Standards & Frameworks
- OWASP Top 10: https://owasp.org/www-project-top-ten/
- OWASP LLM Top 10: https://owasp.org/www-project-top-10-for-large-language-model-applications/
- WCAG 2.1: https://www.w3.org/WAI/WCAG21/quickref/
- AWS Well-Architected: https://aws.amazon.com/architecture/well-architected/
- Microsoft Cloud Adoption Framework: https://docs.microsoft.com/azure/cloud-adoption-framework/

### Related Agents
- `SE: Security` - Deep security analysis
- `SE: Architect` - Architecture validation  
- `SE: Responsible AI` - Bias and ethics review
- `TDD Refactor Phase` - Quality improvement suggestions

### Related Prompts
- `/tdd` - Test-driven development workflow
- `/e2e` - End-to-end testing
- `/build-fix` - Build error resolution
- `/remember` - Capture lessons learned
