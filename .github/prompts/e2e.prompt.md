---
name: 'e2e'
description: 'End-to-End testing workflow for validating complete user journeys and system integrations with stability patterns and result analysis'
agent: 'agent'
argument-hint: 'User journey or integration scenario to test'
tools: ['execute', 'read', 'edit', 'search', 'agent', 'todo']
---

# End-to-End & Integration Testing Workflow

Execute comprehensive E2E and integration tests to validate complete user journeys, cross-component interactions, and system integrations with automated result analysis and stability tracking.

## Mission

Validate end-to-end user workflows and system integrations through realistic testing scenarios. Ensure critical paths work correctly, integrations function properly, and tests remain stable and reliable.

## Scope & Preconditions

### Prerequisites
- **Scenario description**: ${input:scenario:User journey or scenario name}
- **Test environment available**: Dev/staging environment accessible
- **Unit tests passing**: Complete `/tdd` workflow first
- **Context accessible**: `context/` and `requirements/` directories

### What This Workflow Covers
1. ✅ User journey testing (end-to-end flows)
2. ✅ Integration point validation
3. ✅ Cross-browser/platform testing
4. ✅ Performance validation
5. ✅ Failure analysis and classification
6. ✅ Test stability tracking

### What This Does NOT Cover
- ❌ Unit testing (use `/tdd` instead)
- ❌ Build issues (use `/build-fix`)
- ❌ Manual QA or exploratory testing

## Context Loading Strategy

**Use the `context-loader` skill for efficient context loading:**

```
@context-loader
```

This skill implements the index-first approach optimized for E2E testing:

```
1. ✅ Loads `context/README.md` (500 tokens) - lightweight index
2. ✅ Selectively loads E2E-relevant knowledge:
   - context/tech/testing-strategy.md (E2E测试策略)
   - context/tech/architecture.md (系统架构)
   - context/api/endpoints.md (API集成点)
3. ✅ Loads current requirements as needed:
   - requirements/INDEX.md (需求索引)
   - requirements/in-progress/当前需求.md
4. Load detailed context only when needed:
   - context/business/domain-model.md (用户旅程时)
   - context/experience/test-patterns.md (测试模式参考)
```

## Workflow

### Phase 1: Test Planning & Scenario Identification

**Objective**: Identify critical user journeys and integration points to test

**Steps**:

1. **Load Context**
   ```
   Read: context/work-phase/README.md (index)
   Read: context/work-phase/e2e-context.md
   Read: requirement/work-phase/user-journeys.md
   Read: requirement/work-phase/integration-map.md
   ```

2. **Identify Scenarios**
   
   Based on ${input:scenario}, determine test type:
   
   **User Journey Tests** (Happy Paths):
   - Registration → Verification → Login → Setup
   - Product search → Add to cart → Checkout → Payment
   - Document upload → Processing → Review → Download
   
   **Integration Tests** (Component Interactions):
   - Frontend ↔ Backend API
   - Service A ↔ Service B ↔ Database
   - Application ↔ Third-party API
   - Microservice communication
   
   **Cross-Platform Tests**:
   - Multiple browsers (Chrome, Firefox, Safari, Edge)
   - Mobile devices (iOS, Android)
   - Different screen sizes
   - Various network conditions

3. **Prioritize Tests**
   - Business impact (critical vs. nice-to-have)
   - User frequency (daily vs. rare)
   - Risk level (high-risk vs. low-risk)
   - Complexity (simple vs. complex)

4. **Design Test Cases**
   
   For each scenario, define:
   ```markdown
   ## Test Case: [Scenario Name]
   
   **Prerequisites**:
   - User account exists
   - Test data prepared
   - Services healthy
   
   **Steps**:
   1. Navigate to [page]
   2. Interact with [element]
   3. Verify [expected outcome]
   
   **Expected Outcomes**:
   - Element displays correct data
   - Navigation works as expected
   - No errors in console
   
   **Cleanup**:
   - Delete test data
   - Reset state
   - Close connections
   ```

**Output**: Test scenario matrix, prioritized test list

---

### Phase 2: Test Implementation

**Objective**: Create maintainable, reliable E2E tests using best practices

**Steps**:

1. **Load Standards**
   ```
   Read: requirement/work-phase/e2e-standards.md
   Read: requirement/work-phase/test-patterns.md
   ```

2. **Choose Test Framework**
   
   Recommended tools:
   - **Web UI**: Playwright (preferred), Cypress, Selenium
   - **API**: SuperTest, REST Assured, Postman
   - **Mobile**: Appium, Detox
   - **Performance**: k6, JMeter

3. **Implement Tests Using Patterns**
   
   **Page Object Pattern** (for UI tests):
   ```javascript
   // Good: Abstracted UI interactions
   class LoginPage {
     async login(email, password) {
       await this.fillEmail(email);
       await this.fillPassword(password);
       await this.clickSubmit();
       await this.waitForRedirect();
     }
   }
   ```
   
   **Smart Waits** (not hard-coded delays):
   ```javascript
   // Bad: Hard-coded wait
   await page.waitForTimeout(3000);
   
   // Good: Wait for condition
   await page.waitForSelector('[data-testid="success-message"]');
   await page.waitForResponse(res => res.status() === 200);
   ```
   
   **Test Isolation** (independent tests):
   ```javascript
   beforeEach(async () => {
     await setupTestData();
     await loginUser();
   });
   
   afterEach(async () => {
     await cleanupTestData();
     await logout();
   });
   ```

4. **Setup Test Data & Environment**
   - Create test database/state
   - Mock external services (if needed)
   - Configure test environment
   - Prepare fixtures

5. **Implement Cleanup**
   - Teardown test data
   - Reset environment state
   - Clear caches
   - Close connections

**Output**: E2E test suite, fixtures, environment config

---

### Phase 3: Test Execution & Analysis

**Objective**: Run tests, identify failures, provide actionable feedback

**Steps**:

1. **Load Execution Context**
   ```
   Read: context/work-phase/e2e-context.md
   Review: Previous test results (if any)
   ```

2. **Execute Tests**
   ```bash
   # Run E2E tests
   npx playwright test
   # or
   npm run test:e2e
   # or
   pytest tests/e2e/
   ```
   
   Execution configuration:
   - Use headless mode for CI
   - Use headed mode for debugging
   - Capture screenshots on failure
   - Record video for complex scenarios
   - Log detailed execution traces
   - Monitor performance metrics

3. **Capture Artifacts**
   - Screenshots (on failure)
   - Videos (for complex flows)
   - Console logs
   - Network traces
   - Performance metrics (response times, load times)

4. **Analyze Results**
   
   Classify failures into categories:
   
   **Product Bugs** (real issues):
   - Functionality broken
   - Business logic incorrect
   - Data validation failing
   
   **Test Flakiness** (non-deterministic):
   - Race conditions
   - Timing issues
   - Environment instability
   
   **Environment Issues**:
   - Service unavailable
   - Network problems
   - Configuration errors
   
   **Infrastructure Problems**:
   - Test framework bugs
   - Browser driver issues
   - Platform compatibility

5. **Generate Report**
   
   Create summary:
   ```markdown
   ## E2E Test Report
   
   **Run**: [timestamp]
   **Environment**: [dev/staging/prod]
   **Pass Rate**: X/Y (Z%)
   
   ### Passed Tests
   - ✅ user-registration (2.3s)
   - ✅ user-login (1.8s)
   - ✅ checkout-flow (5.2s)
   
   ### Failed Tests
   - ❌ payment-integration (Error: Timeout)
     - Cause: Payment gateway unavailable
     - Action: Retry with mock gateway
   
   ### Flaky Tests
   - ⚠️  product-search (3/5 runs passed)
     - Issue: Race condition in search results
     - Action: Add proper wait for results
   
   ### Performance
   - Avg response time: 120ms
   - P95 response time: 340ms
   - Slowest test: checkout-flow (5.2s)
   ```

6. **Update Core Documentation**
   ```
   Update context/work-phase/e2e-context.md:
   - Last run: [timestamp]
   - Environment: [dev/staging/prod]
   - Pass rate: X/Y (Z%)
   - New failures: [critical ones only]
   - Flaky tests: [list with fix status]
   - Action items: [blockers only]
   ```

**Output**: Test report, artifacts (screenshots/videos), action items

---

## Test Types & Patterns

### 1. User Journey Tests

Validate complete user workflows from start to finish.

**Example**:
```
Registration → Email verification → Login → Profile setup → Dashboard
```

**Best Practices**:
- Test one complete journey per test
- Use realistic test data
- Verify all intermediate states
- Check side effects (emails, notifications, logs)
- Validate final outcomes

### 2. Integration Tests

Validate interactions between components/services.

**Example**:
```
Frontend → API Gateway → Auth Service → User Service → Database
```

**Best Practices**:
- Test actual integrations (not mocks when possible)
- Verify data flow correctness
- Test error handling across boundaries
- Check retry mechanisms
- Validate timeouts and circuit breakers

### 3. Cross-Browser/Platform Tests

Ensure consistent behavior across environments.

**Coverage**:
- Browsers: Chrome, Firefox, Safari, Edge
- Mobile: iOS Safari, Android Chrome
- Screen sizes: Mobile, tablet, desktop
- Network: Fast 3G, 4G, WiFi

### 4. Performance Tests

Validate system behavior under load.

**Metrics**:
- Response times (avg, p95, p99)
- Throughput (requests/second)
- Resource utilization (CPU, memory)
- Error rates
- Scalability (load testing)

## Output Expectations

### Deliverables
- ✅ Comprehensive E2E test suite
- ✅ Test execution report
- ✅ Failure analysis with classification
- ✅ Screenshots/videos for failures
- ✅ Updated context documentation

### Success Criteria
- [ ] All critical paths tested and passing
- [ ] Integration points validated
- [ ] Test stability > 95% (< 5% flakiness)
- [ ] Performance metrics acceptable
- [ ] Security scenarios tested
- [ ] Documentation updated (core only)

### Documentation Updates (Core Only)

**DO UPDATE**:
- Test run status and pass rates
- New failure patterns
- Flaky test list
- Critical blockers

**DON'T UPDATE**:
- Complete console logs
- Full stack traces
- Verbose test output
- Temporary debugging details

## Quality Assurance

### Test Stability Checklist

- [ ] Tests use smart waits (not `sleep()`)
- [ ] Tests are isolated (no shared state)
- [ ] Selectors are stable (data-testid preferred)
- [ ] Tests handle async operations properly
- [ ] Cleanup always runs (even on failure)
- [ ] Tests can run in any order
- [ ] Tests can run in parallel

### Anti-Patterns to Avoid

1. **Testing Implementation Details**: Focus on user behavior, not internals
2. **Hard-Coded Waits**: Use smart waits for conditions
3. **Brittle Selectors**: Use semantic, stable locators
4. **No Test Isolation**: Each test must be independent
5. **Ignoring Flakiness**: Fix flaky tests immediately
6. **Verbose Logging**: Summarize results, don't dump logs

## Integration with Other Workflows

**After TDD**: Unit tests pass, ready for integration testing

```
/tdd feature → tests pass → /e2e user-journey
```

**Before Build**: Gate deployment on E2E results

```
/e2e regression → all pass → /build-fix → deploy
```

**Continuous**: Run on every commit/PR

```
git push → CI runs /e2e smoke → merge if pass
```

## Examples

### Example 1: User Registration Flow

```
User: /e2e user-registration-flow

AI:
1. Plan: Registration → Email verify → Login → Profile setup
2. Implement: Playwright tests with page objects
3. Execute: Run tests, capture results
4. Report: 100% pass, 2.5s avg time
```

### Example 2: Payment Integration

```
User: /e2e payment-integration

AI:
1. Plan: Add to cart → Checkout → Payment → Confirmation
2. Implement: API + UI tests with payment gateway mock
3. Execute: Test all payment methods
4. Report: 95% pass (1 flaky test identified and fixed)
```

### Example 3: Cross-Browser Testing

```
User: /e2e checkout --browsers=chrome,firefox,safari

AI:
1. Plan: Checkout flow across 3 browsers
2. Implement: Playwright multi-browser config
3. Execute: Parallel execution
4. Report: All browsers pass, Safari 15% slower
```

## References

- [Google Testing Blog - Test Sizes](https://testing.googleblog.com/2010/12/test-sizes.html)
- [Martin Fowler - TestPyramid](https://martinfowler.com/bliki/TestPyramid.html)
- [Playwright Best Practices](https://playwright.dev/docs/best-practices)
- [E2E Standards](requirement/work-phase/e2e-standards.md)
- [User Journeys](requirement/work-phase/user-journeys.md)
- [Context Index](context/work-phase/README.md)

---

**Remember**:
- Test user behavior, not implementation
- Keep tests fast and stable
- Run in realistic environments
- Fix flakiness immediately
- Maintain test code like production code
- Document failures with context, not logs
