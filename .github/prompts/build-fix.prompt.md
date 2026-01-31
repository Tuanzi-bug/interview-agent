---
name: 'build-fix'
description: 'Automated build error detection, diagnosis, and repair workflow with pattern matching, dependency resolution, and knowledge accumulation'
agent: 'agent'
argument-hint: 'Optional: specific error type (dependencies, typescript, linting)'
tools: ['execute', 'read', 'edit', 'search', 'agent', 'todo']
---

# Build & Fix Workflow

Systematically diagnose and fix build errors through automated classification, pattern matching, and incremental repair with knowledge base accumulation.

## Mission

Detect, analyze, and repair build errors efficiently by loading relevant error patterns, applying proven solutions, and maintaining an error knowledge base for future reference.

## Scope & Preconditions

### Prerequisites
- **Build system configured**: package.json, Makefile, or build config exists
- **Context accessible**: `context/` and `requirements/` directories
- **Optional**: Error type filter: ${input:errorType:Error type to focus on}

### What This Workflow Covers
1. ✅ Build error detection and classification
2. ✅ Root cause analysis
3. ✅ Pattern matching against known errors
4. ✅ Systematic fix application
5. ✅ Build verification
6. ✅ Error knowledge base updates

### What This Does NOT Cover
- ❌ Unit test failures (use `/tdd` for test-related issues)
- ❌ E2E test failures (use `/e2e` for integration issues)
- ❌ Production runtime errors (use debugging workflows)

## Context Loading Strategy

**Use the `context-loader` skill for intelligent error context loading:**

```
@context-loader
```

This skill implements the index-first approach optimized for build errors:

```
1. ✅ Loads `context/README.md` (500 tokens) - lightweight index
2. ✅ Searches and loads error-relevant knowledge:
   - context/tech/build-system.md (构建系统配置)
   - context/experience/error-patterns.md (已知错误模式库)
3. ✅ Loads current requirements as needed:
   - requirements/INDEX.md (需求索引)
   - requirements/in-progress/当前需求.md
4. ✅ Loads detailed context only when needed:
   - context/tech/tech-stack.md (技术栈依赖)
   - context/experience/lessons-learned.md (历史修复经验)
```

## Workflow

### Phase 1: Build Error Detection

**Objective**: Identify and classify build failures quickly

**Steps**:

1. **Load Context**
   ```
   Read: context/README.md (index)
   Read: context/tech/build-system.md (if exists)
   Read: requirements/INDEX.md
   ```

2. **Trigger Build**
   ```bash
   # Run build command
   npm run build
   # or
   make build
   # or
   go build ./...
   # or
   cargo build
   ```
   
   Capture:
   - Complete output
   - Exit code
   - Build time
   - Error messages

3. **Classify Errors**
   
   Categorize each error into type:
   
   **🔴 Compilation Errors**:
   - Syntax errors
   - Type errors (TypeScript, Go, Rust)
   - Import/module errors
   - Missing declarations
   
   **📦 Dependency Errors**:
   - Missing packages
   - Version conflicts
   - Peer dependency issues
   - Broken lockfile
   - Registry/network issues
   
   **⚙️ Configuration Errors**:
   - Invalid config files
   - Missing environment variables
   - Path resolution issues
   - Build tool configuration
   
   **🧪 Test Failures**:
   - Unit tests failing during build
   - Linting errors
   - Type checking failures
   
   **📏 Linting/Formatting Errors**:
   - ESLint violations
   - Prettier formatting
   - Code style issues
   
   **🏗️ Infrastructure Errors**:
   - Out of memory
   - Build cache corruption
   - File system issues
   - Permission problems

4. **Prioritize Issues**
   
   Rank by severity:
   - **BLOCKER**: Build won't complete
   - **CRITICAL**: Functionality broken
   - **MAJOR**: Significant issues
   - **MINOR**: Warnings, style issues

5. **Extract Error Information**
   ```markdown
   ## Error Summary
   
   **Error Count**: [total]
   **Blockers**: [count]
   
   ### Error #1
   - Type: [category]
   - Severity: [BLOCKER|CRITICAL|MAJOR|MINOR]
   - Location: [file:line]
   - Message: [error message]
   ```

**Output**: Categorized error list, priority ranking

---

### Phase 2: Root Cause Analysis

**Objective**: Understand why the build is failing

**Steps**:

1. **Load Diagnostic Context**
   ```
   Read: context/experience/error-patterns.md (if exists)
   Read: context/tech/tech-stack.md (if exists)
   ```

2. **Analyze Each Error**
   
   For each error, investigate:
   - Extract error message and code
   - Identify affected files/modules
   - Check recent changes (git log, git diff)
   - Review dependency tree
   - Check environment configuration

3. **Pattern Matching**
   
   Search error knowledge base for similar errors:
   ```
   Search: context/work-phase/error-patterns.md
   Match: Error signature against known patterns
   ```
   
   If match found:
   - Load proven solution
   - Estimate fix complexity
   - Assess risk of change
   
   If no match:
   - Search external resources (Stack Overflow, docs)
   - Formulate hypothesis
   - Plan investigation approach

4. **Hypothesis Formation**
   
   For each error, list:
   - **Likely Causes**: [ranked by probability]
   - **Potential Solutions**: [ordered by safety]
   - **Fix Complexity**: [simple|moderate|complex]
   - **Risk Assessment**: [low|medium|high]

5. **Create Fix Plan**
   ```markdown
   ## Fix Plan
   
   ### Phase 1: Quick Wins (Low-Risk)
   1. Clear build cache
   2. Reinstall dependencies
   3. Fix obvious typos
   
   ### Phase 2: Targeted Fixes (Moderate-Risk)
   4. Update conflicting dependencies
   5. Fix type errors
   6. Add missing config
   
   ### Phase 3: Complex Changes (High-Risk)
   7. Refactor problematic code
   8. Update architecture
   ```

**Output**: Root cause analysis, fix proposals, risk assessment

---

### Phase 3: Fix Implementation

**Objective**: Apply fixes systematically and verify success

**Steps**:

1. **Load Fix Guidelines**
   ```
   Read: context/experience/lessons-learned.md (if exists)
   Read: context/experience/error-patterns.md (for proven solutions)
   ```

2. **Apply Fixes Incrementally**
   
   **Principle**: Fix one category at a time, verify after each
   
   **Order of Operations**:
   ```
   1. Environment/config fixes (fastest)
   2. Dependency fixes (low risk)
   3. Compilation fixes (one at a time)
   4. Linting fixes (can often auto-fix)
   5. Complex refactoring (highest risk, last)
   ```
   
   **For Each Fix**:
   ```bash
   # Apply change
   [make the fix]
   
   # Verify build
   npm run build
   
   # If success, commit
   git add .
   git commit -m "fix: [description]"
   
   # If failure, analyze new error
   [repeat process]
   ```

3. **Common Fix Patterns**
   
   **Dependency Errors**:
   ```bash
   # Missing package
   npm install [package]
   
   # Version conflict
   npm install [package]@[compatible-version]
   
   # Broken lockfile
   rm -rf node_modules package-lock.json
   npm install
   
   # Peer dependency issue
   npm install --legacy-peer-deps
   ```
   
   **TypeScript Errors**:
   ```typescript
   // Add type annotations
   function process(data: UserData): Result { }
   
   // Use type guards
   if (typeof value === 'string') { }
   
   // Fix type mismatches
   const result: ExpectedType = transform(input);
   ```
   
   **Configuration Errors**:
   ```bash
   # Missing environment variable
   echo "JWT_SECRET=your-secret" >> .env
   
   # Fix path in config
   # Update tsconfig.json, webpack.config.js, etc.
   ```
   
   **Build Tool Errors**:
   ```bash
   # Out of memory
   export NODE_OPTIONS="--max-old-space-size=4096"
   
   # Cache corruption
   npm run clean
   rm -rf .cache dist
   npm run build
   ```

4. **Verify Fixes**
   
   After each fix:
   ```bash
   # Run build
   npm run build
   
   # Run tests
   npm test
   
   # Check for new errors
   # Compare error count: before vs after
   ```

5. **Handle Fix Failures**
   
   If fix doesn't work:
   - Revert change (`git checkout -- .`)
   - Document why it failed
   - Try alternative approach
   - Escalate if stuck after 3 attempts

6. **Update Documentation**
   
   **If New Error Pattern**:
   ```
   Update context/experience/error-patterns.md:
   
   ### Pattern: [Error Type]
   **Signature**: [error message pattern]
   **Cause**: [root cause]
   **Solution**: [fix approach]
   **Frequency**: [first occurrence]
   ```
   
   **Update Build Status** (轻量级更新):
   ```
   Update context/tech/build-system.md (if needed):
   - Last successful build: [timestamp]
   - Common issues: [brief notes]
   ```

**Output**: Working build, fix summary, updated error patterns

---

## Common Error Patterns & Solutions

### Dependency Errors

| Error | Cause | Solution |
|-------|-------|----------|
| `Cannot find module` | Missing package | `npm install [package]` |
| `ERESOLVE unable to resolve` | Version conflict | Update to compatible versions |
| `Integrity check failed` | Corrupted lockfile | Delete and regenerate lockfile |
| `peer dependency` | Missing peer dep | Install peer dependencies |

### Compilation Errors

| Error | Cause | Solution |
|-------|-------|----------|
| `Type 'X' is not assignable to 'Y'` | Type mismatch | Add type annotations or type guards |
| `Unexpected token` | Syntax error | Fix syntax at indicated line |
| `Cannot find name` | Undeclared variable | Add declaration or import |
| `Module not found` | Bad import path | Fix import statement |

### Configuration Errors

| Error | Cause | Solution |
|-------|-------|----------|
| `Invalid configuration object` | Bad config | Validate config syntax |
| `Environment variable not set` | Missing .env | Create .env with required vars |
| `Cannot resolve path` | Bad path in config | Update path to correct location |

### Infrastructure Errors

| Error | Cause | Solution |
|-------|-------|----------|
| `JavaScript heap out of memory` | Insufficient memory | Increase NODE_OPTIONS heap size |
| `Build cache corrupted` | Bad cache | Clear cache and rebuild |
| `ENOSPC: no space left` | Disk full | Free up disk space |

## Output Expectations

### Deliverables
- ✅ Successful build completion
- ✅ All tests passing
- ✅ No linting errors
- ✅ Fix summary report
- ✅ Updated error patterns (if new)

### Success Criteria
- [ ] Build completes successfully
- [ ] All tests pass
- [ ] No errors or warnings (or acceptable warnings documented)
- [ ] Build time acceptable
- [ ] Output artifacts correct
- [ ] Error patterns documented (if new)

### Documentation Updates (Core Only)

**DO UPDATE**:
- Build status changes
- Error counts and types
- New error patterns
- Fix approaches that worked

**DON'T UPDATE**:
- Complete build logs
- Detailed stack traces
- Temporary debugging output
- Routine successful builds

## Quality Assurance

### Pre-Completion Checklist

- [ ] Build runs successfully
- [ ] All tests pass
- [ ] Linting passes
- [ ] Environment variables set
- [ ] Dependencies installed correctly
- [ ] No new warnings introduced
- [ ] Build time reasonable
- [ ] Documentation updated

### Anti-Patterns to Avoid

1. **Fixing Multiple Issues Simultaneously**: Fix one type at a time
2. **Ignoring Warnings**: Warnings often become errors
3. **Cargo Cult Fixes**: Understand why a fix works
4. **Skipping Verification**: Always verify after each fix
5. **Polluting Docs**: Don't log every build, only updates
6. **Not Learning**: Record new patterns for future

## Integration with Other Workflows

**After TDD**: Build should pass with all tests

```
/tdd feature → tests pass → /build-fix → production build
```

**Before E2E**: Need working build for E2E testing

```
/build-fix → build success → /e2e user-journey
```

**Continuous**: Build on every commit

```
git push → CI runs /build-fix → report results
```

## Examples

### Example 1: Dependency Conflict

```
User: /build-fix

Error: ERESOLVE unable to resolve dependency tree

AI:
1. Classify: Dependency error (BLOCKER)
2. Analyze: react@18 conflicts with react-dom@17
3. Fix: Update react-dom to @18
4. Verify: Build succeeds
5. Document: No new pattern (known issue)
```

### Example 2: TypeScript Errors

```
User: /build-fix typescript

Errors: 15 type errors across 8 files

AI:
1. Prioritize: Group by file
2. Fix file by file:
   - Add type annotations
   - Fix type mismatches
   - Use type guards
3. Verify after each file
4. Result: All errors resolved, build passes
```

### Example 3: New Error Pattern

```
User: /build-fix

Error: Custom webpack plugin failure (new pattern)

AI:
1. Analyze: Plugin incompatible with webpack 5
2. Search: Find webpack 5 compatible version
3. Fix: Update plugin
4. Verify: Build succeeds
5. Document: Add new pattern to error-patterns.md
```

## References

- [Webpack Build Performance](https://webpack.js.org/guides/build-performance/)
- [npm Troubleshooting](https://docs.npmjs.com/troubleshooting)
- [TypeScript Compiler Options](https://www.typescriptlang.org/tsconfig)
- [Build Requirements](requirement/work-phase/build-requirements.md)
- [Fix Patterns](requirement/work-phase/fix-patterns.md)
- [Context Index](context/work-phase/README.md)

---

**Remember**:
- Fix systematically, one error type at a time
- Verify after each fix
- Document new error patterns
- Keep builds fast with caching
- Fail fast, fix faster
- Update context, not logs
