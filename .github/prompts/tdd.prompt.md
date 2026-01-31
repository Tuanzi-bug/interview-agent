---
name: 'tdd'
description: 'Test-Driven Development workflow orchestrating Red-Green-Refactor cycle with context-aware loading and incremental quality improvement'
agent: 'agent'
argument-hint: 'Feature name or bug to develop with TDD'
tools: ['execute', 'read', 'edit', 'search', 'agent', 'todo']
---

# Test-Driven Development Workflow

Execute complete TDD cycle following Red-Green-Refactor methodology with context loading, requirement verification, and documentation updates.

## Mission

Guide test-driven development through systematic phases: write failing tests (Red), implement minimal code (Green), refactor for quality (Refactor). Load relevant context before each phase and update core documentation without pollution.

## Scope & Preconditions

### Prerequisites
- **Feature/bug description**: ${input:feature:Feature name or bug to develop}
- **Test framework configured**: Unit testing setup exists
- **Context accessible**: `context/` and `requirements/` directories

### What This Workflow Covers
1. ✅ Complete Red-Green-Refactor cycle
2. ✅ Context-aware requirement loading
3. ✅ Test-first development
4. ✅ Security and quality refactoring
5. ✅ Core documentation updates

### What This Does NOT Cover
- ❌ E2E testing (use `/e2e` instead)
- ❌ Build configuration (use `/build-fix` for build issues)
- ❌ Manual testing or QA processes

## Context Loading Strategy

**Use the `context-loader` skill for efficient context loading:**

```
@context-loader
```

This skill automatically implements the index-first strategy:

1. ✅ Loads `context/README.md` (500 tokens) - lightweight index
2. ✅ Selectively loads testing domain knowledge:
   - `context/tech/testing-strategy.md` (测试策略)
   - `context/experience/test-patterns.md` (测试模式)
3. ✅ Loads current requirements as needed:
   - `requirements/INDEX.md` (需求索引)
   - `requirements/in-progress/当前需求.md`
4. ✅ Loads detailed context only when needed:
   - `context/tech/architecture.md` (架构变更时)
   - `context/business/domain-model.md` (业务逻辑时)

**Efficiency**: Reduces context load by 80-85% compared to loading all files.

See [context-loader skill](../.github/skills/context-loader/SKILL.md) for details.

## Workflow

### Phase 1: RED - Write Failing Tests First

**Objective**: Create tests that define desired behavior before implementation exists

**Agent**: Use `TDD Red Phase - Write Failing Tests First` subagent

**Steps**:

1. **Load Context**
   ```
   Read: context/README.md (index)
   Read: context/tech/testing-strategy.md (if exists)
   Read: requirements/INDEX.md
   Read: requirements/in-progress/当前需求.md (current requirement)
   ```

2. **Analyze Requirements**
   - Parse feature description: ${input:feature}
   - Identify acceptance criteria from GitHub issue (if exists)
   - Break down into testable units
   - List edge cases and error scenarios

3. **Write Comprehensive Tests**
   - Create test files following project conventions
   - Use descriptive test names (describe what should happen)
   - Structure: Arrange-Act-Assert or Given-When-Then
   - Include:
     - Happy path scenarios
     - Edge cases (null, empty, boundary values)
     - Error scenarios and validation
     - Security requirements

4. **Verify Red State**
   ```bash
   # Run tests to confirm they fail
   npm test [test-file] || pytest [test-file] || go test
   ```
   - Ensure tests fail for the right reason
   - Document why each test fails
   - Capture failure messages

5. **Document Red Phase**
   ```
   Update context/work-phase/tdd-context.md:
   - Status: IDLE → RED
   - Feature: ${input:feature}
   - Tests written: [count]
   - Expected failures: [list]
   ```

**Output**: Test files with failing tests, clear failure messages

---

### Phase 2: GREEN - Make Tests Pass Quickly

**Objective**: Implement simplest code to make all tests pass without over-engineering

**Agent**: Use `TDD Green Phase - Make Tests Pass Quickly` subagent

**Steps**:

1. **Load Previous Context**
   ```
   Read: context/work-phase/tdd-context.md (current state)
   Review: Failing tests from Red phase
   ```

2. **Implement Solution**
   - Write simplest code that makes tests pass
   - Focus on functionality, not optimization
   - Avoid premature abstraction
   - Keep implementation straightforward
   - Maintain test isolation

3. **Verify Green State**
   ```bash
   # Run all tests
   npm test || pytest || go test ./...
   ```
   - Confirm 100% pass rate
   - Check for false positives
   - Verify tests are actually testing what they claim

4. **Check Coverage**
   ```bash
   # Generate coverage report
   npm run coverage || pytest --cov || go test -cover
   ```
   - Ensure minimum threshold met (typically 80%)
   - Identify any gaps

5. **Update Core Status**
   ```
   Update context/work-phase/tdd-context.md:
   - Status: RED → GREEN
   - Tests passing: [count]/[count]
   - Coverage: [percentage]%
   - Implementation: [brief summary]
   ```

**Output**: Working implementation, all tests passing, coverage report

---

### Phase 3: REFACTOR - Improve Quality & Security

**Objective**: Enhance code quality, apply security best practices, improve design while maintaining green tests

**Agent**: Use `TDD Refactor Phase - Improve Quality & Security` subagent

**Steps**:

1. **Load Quality Standards**
   ```
   Read: requirement/work-phase/quality-standards.md
   Read: requirement/work-phase/security-checklist.md
   ```

2. **Apply Improvements**
   
   **Code Quality**:
   - Extract duplicated code into functions
   - Apply appropriate design patterns
   - Improve naming clarity
   - Add documentation/comments
   - Simplify complex logic
   
   **Security**:
   - Input validation
   - Output encoding
   - Authentication/authorization checks
   - Secure data handling
   - Error handling (don't leak sensitive info)
   
   **Design**:
   - SOLID principles
   - Dependency injection
   - Separation of concerns
   - Proper error handling

3. **Maintain Test Coverage**
   ```bash
   # Run tests after EACH change
   npm test
   ```
   - Ensure tests remain green
   - Add tests for discovered edge cases
   - Refactor tests if needed for clarity

4. **Verify No Regression**
   - All original tests still pass
   - No new bugs introduced
   - Coverage maintained or improved

5. **Update Documentation (Core Only)**
   ```
   Update context/work-phase/tdd-context.md:
   - Status: GREEN → REFACTOR → COMPLETE
   - Refactoring applied: [patterns/techniques]
   - Security improvements: [specific fixes]
   - Code quality metrics: [improvements]
   - Next steps: Run E2E tests
   ```

**Output**: Clean, refactored code, tests passing, improved security

---

## Output Expectations

### Deliverables
- ✅ Comprehensive test suite (Red phase)
- ✅ Working implementation (Green phase)
- ✅ Refactored, secure code (Refactor phase)
- ✅ Updated context documentation (core sections only)
- ✅ Coverage report meeting threshold

### Documentation Updates (Core Only)

**DO UPDATE** in `context/work-phase/tdd-context.md`:
- Phase status changes
- Test counts and pass rates
- Coverage percentages
- Critical decisions
- Blocking issues

**DON'T UPDATE**:
- Detailed test output logs
- Complete stack traces
- Implementation details (use code comments)
- Temporary debugging notes

### Success Criteria
- [ ] All tests written and passing
- [ ] Test coverage ≥ 80% (or project threshold)
- [ ] No code smells detected
- [ ] Security best practices applied
- [ ] Documentation updated (core sections)
- [ ] Ready for E2E testing

## Quality Assurance

### Pre-completion Checklist

**Red Phase**:
- [ ] Tests fail with expected error messages
- [ ] Test names clearly describe expected behavior
- [ ] Tests are isolated and independent
- [ ] Edge cases identified
- [ ] Error scenarios covered

**Green Phase**:
- [ ] All tests pass (100%)
- [ ] Implementation is minimal (no over-engineering)
- [ ] Tests remain isolated
- [ ] No false positives

**Refactor Phase**:
- [ ] Code is DRY (Don't Repeat Yourself)
- [ ] SOLID principles applied
- [ ] Security vulnerabilities addressed
- [ ] Tests still pass (100%)
- [ ] Documentation updated

### Anti-Patterns to Avoid

1. **Writing Implementation Before Tests**: Always write tests first
2. **Over-Engineering in Green Phase**: Keep it simple, optimize in Refactor
3. **Skipping Refactor Phase**: Never skip quality improvements
4. **Polluting Documentation**: Only update core status/progress
5. **Breaking Tests During Refactor**: Run tests after each change

## Integration with Other Workflows

**Before E2E**: Complete TDD cycle ensures units work correctly

```
/tdd feature-name → tests pass → /e2e user-journey
```

**Before Build**: All tests must pass

```
/tdd feature-name → /build-fix → production build
```

**After Plan**: Convert requirements into tests first

```
/plan feature → /tdd feature → implementation
```

## Examples

### Example 1: New Feature

```
User: /tdd user-authentication

AI:
1. RED: Write authentication tests (registration, login, token validation)
2. GREEN: Implement auth service (minimal code)
3. REFACTOR: Add security (bcrypt, JWT), improve structure
4. Output: Auth service with 95% coverage, all tests passing
```

### Example 2: Bug Fix

```
User: /tdd fix-password-validation-bug

AI:
1. RED: Write test that reproduces the bug
2. GREEN: Fix the validation logic
3. REFACTOR: Improve validation, add edge case handling
4. Output: Bug fixed, tests pass, validation improved
```

## References

- [Kent Beck - Test-Driven Development](https://www.amazon.com/Test-Driven-Development-Kent-Beck/dp/0321146530)
- [Martin Fowler - TDD](https://martinfowler.com/bliki/TestDrivenDevelopment.html)
- [Uncle Bob - The Three Rules of TDD](http://butunclebob.com/ArticleS.UncleBob.TheThreeRulesOfTdd)
- [TDD Requirements](requirement/work-phase/tdd-requirements.md)
- [Context Index](context/work-phase/README.md)

---

**Remember**: 
- Tests first, always
- Green phase: simple works
- Refactor phase: quality matters
- Update context core, not logs
- TDD ensures quality from the start
