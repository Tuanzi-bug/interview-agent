---
name: 'plan'
description: 'Professional software engineering planning workflow with structured requirement management. Integrates product analysis, architecture design, UX research, and TDD planning into a unified process with organized requirement documentation.'
agent: 'agent'
argument-hint: 'Describe your feature, project, or idea to plan'
tools: ['read', 'edit', 'search', 'github/*', 'agent', 'todo']
---

# Professional Software Engineering Plan

Execute a comprehensive planning workflow following industry-standard best practices. This workflow integrates multiple planning disciplines into a structured process that outputs organized, traceable requirements.

## Mission

Transform user ideas into actionable, well-documented requirements through systematic analysis. Generate structured requirement documents following professional software engineering standards, ensuring all aspects (product, architecture, UX, testing) are properly planned and tracked.

**CRITICAL**: Follow the [never-fabricate principle](../.github/instructions/never-fabricate.instructions.md) - stop and request clarification when uncertain rather than inventing details.

## Scope & Preconditions

### Prerequisites
- **Workspace must exist**: ${workspaceFolder}
- **User input required**: Feature description, project idea, or problem statement
- **Optional context**: Existing codebase, stakeholder inputs, or constraints

### What This Workflow Covers
1. ✅ Requirements discovery and analysis
2. ✅ Product value assessment (Jobs-to-be-Done)
3. ✅ User experience research and journey mapping
4. ✅ System architecture planning
5. ✅ Test-driven development planning (test cases first)
6. ✅ Structured requirement documentation
7. ✅ GitHub issue creation and linking
8. ✅ Requirement tracking and indexing

### What This Does NOT Cover
- ❌ Code implementation (use Work stage agents/prompts)
- ❌ Detailed UI design mockups (outputs research artifacts only)
- ❌ Actual test execution (defines test cases only)

## Context Loading Strategy

**Use the `context-loader` skill for efficient context loading:**

```
@context-loader
```

This skill implements the index-first strategy:
- Loads `context/README.md` and `requirement/INDEX.md` first
- Selectively loads architecture and design patterns based on planning needs
- Prevents context pollution by loading only relevant domain files
- Typical load: 3-5 files, ~2,500-4,000 tokens (vs 15+ files, 18,000+ tokens)

**For planning tasks, typically loads:**
- Architecture decisions and design patterns
- Existing requirements and dependencies
- Technical constraints and framework knowledge

See [context-loader skill](../.github/skills/context-loader/SKILL.md) for details.

---

## Workflow

Use the **todo list** to track progress through all planning stages.

### Phase 1: Discovery & Context Analysis

**Objective**: Understand what the user wants to build and why.

**First, load context efficiently:**
```
@context-loader
```

1. **Gather User Input**
   - Extract feature description from user request
   - Identify any explicit constraints (budget, timeline, tech stack)
   - Note any referenced existing issues, documents, or code

2. **Search Existing Context**
   - Search workspace for related features or documentation
   - Check for existing requirement files in `requirements/`
   - Review any existing GitHub issues related to the request
   - Identify potential conflicts or dependencies

3. **Clarify Through Questions** (Use SE: Product Manager perspective)
   - **Who's the user?** (Role, skill level, frequency of use)
   - **What problem are they solving?** (Current workflow, pain points, costs)
   - **What's the desired outcome?** (Success criteria, metrics)
   - **What are the constraints?** (Technical, business, timeline)
   - **What's out of scope?** (Explicitly excluded features)
   
   *Note: Apply never-fabricate principle throughout - request clarification for any gaps.*

4. **Document Initial Context**
   - Summarize user need in one sentence
   - List key stakeholders and their concerns
   - Identify primary success metrics
   - Note any critical assumptions

### Phase 2: Product Analysis (Jobs-to-be-Done)

**Objective**: Understand the job users are hiring this feature to do.

**Invoke SE: Product Manager agent for:**

1. **Jobs-to-be-Done Analysis**
   ```
   Job Statement: When [situation], I want to [motivation], 
                  so I can [expected outcome]
   ```
   - Identify the functional job
   - Identify emotional and social jobs
   - Map current alternatives (what users do today)

2. **Business Value Assessment**
   - Quantify impact (time saved, errors reduced, revenue increased)
   - Identify metrics to track success
   - Assess effort vs. value (prioritization)

3. **Risk & Dependency Analysis**
   - Technical risks
   - Resource constraints
   - Dependencies on other features or teams

### Phase 3: User Experience Research

**Objective**: Map user journeys and identify UX requirements.

**Invoke SE: UX Designer agent for:**

1. **User Context Analysis**
   - User persona(s) with skill levels and accessibility needs
   - Usage context (when, where, device, interruptions)
   - Emotional state during use

2. **User Journey Mapping**
   - Current state journey (pain points, workarounds)
   - Desired future state journey
   - Key touchpoints and decision moments
   - Potential friction points

3. **UX Requirements**
   - Information architecture needs
   - Interaction patterns required
   - Accessibility requirements (WCAG 2.1 AA minimum)
   - Responsive design considerations

### Phase 4: System Architecture Planning

**Objective**: Design system structure and validate architectural decisions.

**Invoke SE: Architect agent for:**

1. **Architecture Context**
   - System type (web app, microservice, AI system, etc.)
   - Scale requirements (users, requests, data volume)
   - Integration points with existing systems

2. **Architecture Design**
   - Component structure and responsibilities
   - Data flow and state management
   - API design (if applicable)
   - Security architecture
   - Scalability considerations

3. **Technology Decisions**
   - Framework/library choices with rationale
   - Infrastructure requirements
   - Performance targets
   - Monitoring and observability strategy

4. **Architecture Validation**
   - Apply relevant Well-Architected principles
   - Identify architectural risks
   - Plan for failure modes and recovery

### Phase 5: Test Planning (TDD Red Phase)

**Objective**: Define test cases before implementation exists.

**Invoke TDD Red Phase agent for:**

1. **Test Strategy**
   - Test levels needed (unit, integration, E2E)
   - Test coverage targets
   - Test data requirements
   - Test environment needs

2. **Acceptance Test Cases**
   - Convert requirements to test scenarios
   - Define Given-When-Then test cases
   - Identify edge cases and boundary conditions
   - Specify expected behaviors for error cases

3. **Test Structure Planning**
   - Test file organization
   - Test utilities/fixtures needed
   - Mocking strategy for external dependencies

### Phase 6: Requirement Documentation

**Objective**: Create structured, traceable requirement documents.

**Directory Structure** (create if not exists):
```
${workspaceFolder}/requirements/
├── INDEX.md                    # Master requirement index
├── in-progress/                # Active requirements being refined
│   └── REQ-{ID}-{slug}.md     # Individual requirement files
└── completed/                  # Finalized requirements
    └── REQ-{ID}-{slug}.md     # Completed requirement files
```

**For each new requirement:**

1. **Generate Requirement ID**
   - Format: `REQ-YYYY-NNN` (e.g., `REQ-2026-001`)
   - Read INDEX.md to determine next available ID
   - Create slug from feature title (e.g., `user-authentication`)

2. **Create Requirement File**: `requirements/in-progress/REQ-{ID}-{slug}.md`

   Use this template:

   ````markdown
   # REQ-{ID}: {Feature Title}
   
   **Status**: In Progress  
   **Created**: {Date}  
   **Last Updated**: {Date}  
   **Priority**: {High|Medium|Low}  
   **Related Issues**: #{issue-number} (if exists)
   
   ## Executive Summary
   
   [One paragraph: what, why, who benefits]
   
   ## Jobs-to-be-Done
   
   **Job Statement**:  
   When [situation], I want to [motivation], so I can [outcome]
   
   **Current Alternatives**:  
   - [How users solve this today]
   
   **Success Metrics**:  
   - [Metric 1]: [Target value]
   - [Metric 2]: [Target value]
   
   ## User Context
   
   **Primary User Persona**:
   - Role: [Job title/role]
   - Skill Level: [Beginner|Intermediate|Expert]
   - Context: [When/where they use this]
   - Accessibility Needs: [Any specific needs]
   
   **User Journey**:
   
   | Stage | Current State | Desired State | Pain Points |
   |-------|---------------|---------------|-------------|
   | [Step 1] | [What they do now] | [What they'll do] | [Issues] |
   | [Step 2] | [What they do now] | [What they'll do] | [Issues] |
   
   ## Functional Requirements
   
   ### Core Features
   
   | ID | Requirement | Priority | Acceptance Criteria |
   |----|-------------|----------|---------------------|
   | FR-001 | [Requirement description] | [H/M/L] | [How to verify] |
   | FR-002 | [Requirement description] | [H/M/L] | [How to verify] |
   
   ### Non-Functional Requirements
   
   | ID | Type | Requirement | Target |
   |----|------|-------------|--------|
   | NFR-001 | Performance | [Requirement] | [Metric] |
   | NFR-002 | Security | [Requirement] | [Standard] |
   | NFR-003 | Accessibility | [Requirement] | [WCAG level] |
   
   ## Technical Architecture
   
   ### System Components
   
   ```
   [Component diagram or description]
   ```
   
   ### Data Model
   
   ```
   [Key entities and relationships]
   ```
   
   ### API Design (if applicable)
   
   | Endpoint | Method | Purpose | Request | Response |
   |----------|--------|---------|---------|----------|
   | [path] | [GET/POST] | [Description] | [Schema] | [Schema] |
   
   ### Technology Stack
   
   - **Frontend**: [Framework/libraries with rationale]
   - **Backend**: [Framework/libraries with rationale]
   - **Database**: [Type and rationale]
   - **Infrastructure**: [Hosting/deployment approach]
   
   ### Security Considerations
   
   - [Authentication/Authorization approach]
   - [Data protection measures]
   - [OWASP Top 10 mitigation strategies]
   
   ## Test Planning
   
   ### Test Strategy
   
   - **Unit Tests**: [Coverage target, key test areas]
   - **Integration Tests**: [What integrations to test]
   - **E2E Tests**: [Critical user flows]
   
   ### Acceptance Test Cases
   
   | Test ID | Scenario | Given | When | Then | Priority |
   |---------|----------|-------|------|------|----------|
   | TC-001 | [Description] | [Context] | [Action] | [Expected] | [H/M/L] |
   | TC-002 | [Description] | [Context] | [Action] | [Expected] | [H/M/L] |
   
   ### Edge Cases & Error Handling
   
   | Scenario | Expected Behavior |
   |----------|-------------------|
   | [Edge case 1] | [How system should respond] |
   | [Error condition 1] | [Error message and recovery] |
   
   ## Implementation Plan
   
   ### Phase Breakdown
   
   | Phase | Description | Deliverables | Estimated Effort |
   |-------|-------------|--------------|------------------|
   | 1 | [Phase name] | [What's delivered] | [Time estimate] |
   | 2 | [Phase name] | [What's delivered] | [Time estimate] |
   
   ### Dependencies
   
   - [Dependency 1]: [Description and impact]
   - [Dependency 2]: [Description and impact]
   
   ### Risks & Mitigations
   
   | Risk | Probability | Impact | Mitigation Strategy |
   |------|-------------|--------|---------------------|
   | [Risk 1] | [H/M/L] | [H/M/L] | [How to mitigate] |
   | [Risk 2] | [H/M/L] | [H/M/L] | [How to mitigate] |
   
   ## Out of Scope
   
   - [Explicitly excluded feature 1]
   - [Explicitly excluded feature 2]
   - [Future enhancement 1]
   
   ## Sign-off Checklist
   
   - [ ] Business value clearly articulated
   - [ ] User needs validated
   - [ ] Acceptance criteria defined
   - [ ] Architecture reviewed and approved
   - [ ] Security considerations addressed
   - [ ] Test strategy defined
   - [ ] Dependencies identified
   - [ ] Risks assessed and mitigated
   - [ ] Implementation plan agreed upon
   
   ## Related Documents
   
   - GitHub Issue: #{issue-number}
   - Architecture Decision Records: [Links if applicable]
   - Design Mockups: [Links if applicable]
   - API Documentation: [Links if applicable]
   
   ## Change Log
   
   | Date | Author | Change |
   |------|--------|--------|
   | {Date} | {Name} | Initial creation |
   
   ````

3. **Update INDEX.md**

   If INDEX.md doesn't exist, create it with this template:

   ````markdown
   # Requirements Index
   
   Master index of all requirements for ${workspaceFolder}
   
   **Last Updated**: {Date}
   
   ## In Progress
   
   Requirements currently being refined and planned.
   
   | Req ID | Title | Priority | Created | Owner | Status |
   |--------|-------|----------|---------|-------|--------|
   | [REQ-{ID}](in-progress/REQ-{ID}-{slug}.md) | [Title] | [H/M/L] | {Date} | [Name] | Planning |
   
   ## Completed
   
   Requirements that have been finalized and approved for implementation.
   
   | Req ID | Title | Priority | Completed | Owner | GitHub Issue |
   |--------|-------|----------|-----------|-------|--------------|
   | [REQ-{ID}](completed/REQ-{ID}-{slug}.md) | [Title] | [H/M/L] | {Date} | [Name] | #{issue} |
   
   ## Statistics
   
   - **Total Requirements**: {count}
   - **In Progress**: {count}
   - **Completed**: {count}
   - **High Priority**: {count}
   - **Medium Priority**: {count}
   - **Low Priority**: {count}
   
   ## Quick Links
   
   - [Requirements Directory](.)
   - [GitHub Issues](../../issues)
   - [Architecture Decisions](../docs/adr)
   - [Project README](../../README.md)
   ````

   If INDEX.md exists, append the new requirement to the "In Progress" section.

### Phase 7: GitHub Issue Creation (Optional)

**Objective**: Create traceable GitHub issue linked to requirement.

**Invoke github-issues skill to:**

1. **Create Issue with Requirement Context**
   ```
   Title: [REQ-{ID}] {Feature Title}
   
   Body:
   - Link to requirement document
   - Executive summary
   - Acceptance criteria
   - Implementation phases
   - Related requirements
   ```

2. **Apply Labels**
   - `requirement`
   - `planning`
   - Priority label: `priority:high|medium|low`
   - Type label: `feature|enhancement|bug`

3. **Update Requirement Document**
   - Add GitHub issue link to "Related Issues" field
   - Update INDEX.md with issue number

### Phase 8: Hand-off Preparation

**Objective**: Ensure smooth transition to implementation.

1. **Summary Generation**
   - Create concise summary of all planning outputs
   - List all generated documents and their locations
   - Highlight key decisions and rationale
   - Identify next steps for implementation

2. **Readiness Checklist**
   - [ ] Requirements are clear and testable
   - [ ] User needs are validated
   - [ ] Architecture is sound and reviewed
   - [ ] Test cases are defined
   - [ ] Dependencies are documented
   - [ ] Risks are identified with mitigations
   - [ ] GitHub issue created (if requested)
   - [ ] INDEX.md updated
   - [ ] Implementation team can proceed

3. **Hand-off to Work Stage**
   Provide clear next steps:
   - Reference TDD Green Phase agent for implementation
   - Point to Context7 agent for API documentation during coding
   - Mention relevant work stage agents (DevOps, Technical Writer)

## Output Expectations

### Primary Outputs

1. **Requirement Document**: `requirements/in-progress/REQ-{ID}-{slug}.md`
   - Complete requirement specification following template
   - All sections filled with relevant information
   - Test cases clearly defined
   - Architecture documented

2. **Updated INDEX.md**: `requirements/INDEX.md`
   - New requirement added to "In Progress" section
   - Statistics updated
   - Proper linking to requirement file

3. **GitHub Issue** (if requested)
   - Properly formatted issue with requirement context
   - Appropriate labels applied
   - Linked to requirement document

### Success Criteria

- ✅ All planning phases completed without gaps
- ✅ Requirement document is comprehensive and actionable
- ✅ Technical feasibility validated through architecture review
- ✅ Test cases defined for all acceptance criteria
- ✅ Dependencies and risks documented
- ✅ Clear hand-off to implementation team

### Failure Conditions

- ❌ User needs are unclear after discovery (STOP and ask more questions)
- ❌ Technical feasibility is questionable (STOP and escalate)
- ❌ Critical dependencies are blocked (STOP and document blockers)
- ❌ Requirements conflict with existing system (STOP and resolve conflicts)
- ❌ **Missing critical information** (STOP and request clarification - NEVER fabricate)
- ❌ **Uncertain about technical details** (STOP and ask for specifics - NEVER invent)

## Quality Assurance
Missing critical information (apply never-fabricate principle

After completing the workflow:

1. **Requirement Quality Check**
   ```bash
   # Verify file exists
   ls -la requirements/in-progress/REQ-*.md
   
   # Check file size (should be substantial)
   wc -l requirements/in-progress/REQ-*.md
   ```

2. **Completeness Check**
   - All template sections are filled (not just placeholders)
   - Acceptance criteria are specific and testable
   - Test cases include Given-When-Then format
   - Architecture diagram or description is clear

3. **Consistency Check**
   - INDEX.md references the new requirement
   - Requirement ID is unique and follows format
   - GitHub issue (if created) links to requirement
   - Related documents are cross-referenced

4. **Review Prompts for User**
   - "Does this capture your intended feature accurately?"
   - "Are the success metrics aligned with your goals?"
   - "Is anything missing from the requirement specification?"
   - "Are the test cases comprehensive enough?"

### Common Issues & Troubleshooting

| Issue | Cause | Solution |
|-------|-------|----------|
| INDEX.md missing | First requirement in project | Create INDEX.md using template |
| Requirement ID conflict | ID already exists | Read INDEX.md and use next available ID |
| Unclear user needs | Insufficient discovery | Return to Phase 1 and ask clarifying questions |
| Architecture complexity | Feature is too large | Break into multiple smaller requirements |
| Missing agent access | Agent not available | Use inline analysis following agent principles |
| **Uncertain about technical details** | **Information gap** | **STOP workflow, ask specific questions, never fabricate** |
| **Missing implementation specifics** | **Insufficient context** | **Request examples, documentation, or codebase references** |

## Integration with Other Stages
Uncertain details | Information gap | Apply never-fabricate principle: stop and request clarification

### To Work Stage (TDD Green)
- Hand off requirement document to implementation team
- Test cases guide TDD implementation
- Architecture design informs code structure

### To Review Stage
- Acceptance criteria guide code review
- Test cases validate implementation completeness
- Architecture design is reference for architecture review

### To Compound Stage
- Planning decisions feed into memory bank
- Successful patterns documented for future reuse
- Lessons learned update planning templates

## Examples

### Example Usage

**User Request**: "I want to add user authentication to my app"

**Planning Response**:
1. Discovery: Ask about user types, authentication methods needed, security requirements
2. Product Analysis: Identify job (secure access to personal data)
3. UX Research: Map login/registration journey, identify friction points
4. Architecture: Design auth service, token management, session handling
5. Test Planning: Define test cases for login, registration, password reset
6. Documentation: Create `requirements/in-progress/REQ-2026-001-user-authentication.md`
7. GitHub Issue: Create issue #42 linked to requirement
8. Hand-off: Ready for TDD Green implementation

**Generated Files**:
- `requirements/INDEX.md` (created/updated)
- `requirements/in-progress/REQ-2026-001-user-authentication.md`
- GitHub Issue #42 (if requested)

### Example Requirement File Snippet

```markdown
# REQ-2026-001: User Authentication System

**Status**: In Progress  
**Created**: 2026-01-30  
**Priority**: High  

## Executive Summary

Enable secure user authentication to protect user data and personalize experiences. 
Supports email/password and OAuth (Google, GitHub) authentication methods.

## Jobs-to-be-Done

**Job Statement**:  
When users want to access their personal data, they want to securely authenticate 
their identity, so they can trust their information is protected.

**Current Alternatives**:  
- No authentication (insecure)
- Third-party auth services only

**Success Metrics**:  
- Registration completion rate: >80%
- Login success rate: >95%
- Password reset completion: >70%
- Zero authentication-related security incidents

## Functional Requirements

| ID | Requirement | Priority | Acceptance Criteria |
|----|-------------|----------|---------------------|
| FR-001 | User registration with email/password | High | User can create account with valid email and secure password |
| FR-002 | Email verification | High | Verification email sent, account activated on click |
| FR-003 | User login with email/password | High | Valid credentials grant access, invalid credentials show error |
| FR-004 | Password reset flow | High | User can reset password via email link |
| FR-005 | OAuth login (Google, GitHub) | Medium | Users can authenticate via OAuth providers |

## Test Planning

### Acceptance Test Cases

| Test ID | Scenario | Given | When | Then | Priority |
|---------|----------|-------|------|------|----------|
| TC-001 | Successful registration | Valid email and password | User submits registration | Account created, verification email sent | H |
| TC-002 | Duplicate email registration | Email already exists | User tries to register | Error message shown, registration blocked | H |
| TC-003 | Weak password rejection | Password is too simple | User submits weak password | Error shown with password requirements | M |
| TC-004 | Successful login | Valid credentials | User submits login form | User authenticated and redirected to dashboard | H |
| TC-005 | Failed login attempt | Invalid credentials | User submits incorrect password | Error shown, account not locked (after 5 attempts, lock) | H |
```

## Related Resources

### Agents Used in This Workflow
- **SE: Product Manager** - Business value and JTBD analysis
- **SE: UX Designer** - User research and journey mapping
- **SE: Architect** - Architecture design and validation
- **TDD Red Phase** - Test planning and test case definition
- **Context7** - API documentation and best practices lookup

### Skills Used
- **prd** - PRD generation structure
- **github-issues** - Issue creation and management

### Next Stage Agents
- **TDD Green Phase** - Implementation following requirement
- **SE: DevOps/CI** - Deployment and pipeline setup
- **SE: Technical Writer** - Documentation writing

## Notes

- This prompt creates a **living document** in `requirements/` that evolves through the project lifecycle
- Move requirements from `in-progress/` to `completed/` when approved and ready for implementation
- Use `requirements/INDEX.md` as single source of truth for all requirements
- Link GitHub issues to requirement documents for traceability
- Update requirement documents as implementation progresses with new learnings

---

**Version**: 1.0  
**Last Updated**: 2026-01-30  
**Maintained By**: GitHub Copilot Code Repository