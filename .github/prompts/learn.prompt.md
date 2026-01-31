---
name: 'learn'
description: 'Professional knowledge synthesis and continuous learning for Compound stage. Loads context intelligently, synthesizes experience into reusable skills, and maintains clean domain knowledge without pollution.'
agent: 'agent'
argument-hint: 'What did you learn? (patterns, bugs, solutions, workflows)'
tools: ['read', 'edit', 'search', 'todo']
---

# Professional Knowledge Synthesis & Continuous Learning

Transform project experience into structured, reusable knowledge assets. This workflow implements professional software engineering best practices for knowledge management, context engineering, and skill creation during the Compound stage.

## Mission

Systematically capture lessons learned, bug patterns, workflow improvements, and technical insights, then synthesize them into well-organized domain knowledge and reusable skills. Prevent knowledge loss and redundancy while building a compound learning system that makes the team smarter with each project.

## Core Principles of Effective Context Engineering

### 1. **Precision Over Volume**
- Load ONLY what's relevant to the current learning synthesis
- Read INDEX files first, then selectively load specific documents
- Avoid dumping entire directories into context
- Use semantic search to find relevant prior knowledge

### 2. **Incremental Updates**
- Update core sections only, never rewrite entire documents
- Append new learnings to existing sections
- Merge similar patterns before adding new ones
- Maintain document coherence without pollution

### 3. **Domain Separation**
- Keep technical knowledge separate from experience knowledge
- Organize by domain (testing, build, architecture, etc.)
- Create new domain files when existing ones don't fit
- Cross-reference related domains without duplication

### 4. **Verification Before Storage**
- Validate that the learning is repeatable and useful
- Ensure it's not already documented elsewhere
- Confirm it provides actionable guidance
- Test that examples are correct and complete

## Scope & Preconditions

### Prerequisites
- **Workspace exists**: ${workspaceFolder}
- **User provides input**: 
  - Learning (pattern, bug fix, workflow improvement, insight) OR
  - Request to bootstrap context structure for existing project

### Operational Modes

This prompt operates in **two modes** based on project state:

**Mode 1: Bootstrap Mode** (Auto-detected)
- **Trigger**: context/ or requirements/ directories don't exist
- **Action**: Initialize knowledge structure from existing project
- **Output**: Complete context/ and requirements/ structure with initial knowledge

**Mode 2: Learning Synthesis Mode** (Default)
- **Trigger**: Knowledge structure exists
- **Action**: Capture and synthesize new learnings
- **Output**: Incremental updates to existing knowledge base

The prompt **automatically detects** which mode to use and informs the user.

### Directory Structure (Auto-created)

```
${workspaceFolder}/
├── context/
│   ├── README.md                    # Lightweight index (prioritize)
│   ├── tech/
│   │   ├── README.md               # Tech domain index
│   │   ├── testing-strategy.md     # TDD/E2E patterns
│   │   ├── build-system.md         # Build configurations
│   │   ├── architecture.md         # System architecture decisions
│   │   └── frameworks.md           # Framework-specific knowledge
│   └── experience/
│       ├── README.md               # Experience domain index
│       ├── error-patterns.md       # Common errors and solutions
│       ├── test-patterns.md        # Effective test patterns
│       ├── debugging-tactics.md    # Debugging strategies
│       └── workflow-improvements.md # Process optimizations
│
└── requirements/
    ├── INDEX.md                     # Requirement index (prioritize)
    ├── in-progress/                 # Active requirements
    └── completed/                   # Finalized requirements
```
stages.


**Bootstrap Mode**:
1. ✅ Automatic detection of missing knowledge structure
2. ✅ Project analysis and tech stack detection
3. ✅ Knowledge base initialization with discovered patterns
4. ✅ Directory and file structure creation
5. ✅ Initial domain file generation with project-specific content

**Learning Synthesis Mode**:
1. ✅ Intelligent context loading (index-first strategy)
2. ✅ Learning analysis and categorization
3. ✅ Domain knowledge synthesis
4. ✅ Skill creation for complex workflows
5. ✅ Incremental document updates (no pollution)
6. ✅ Cross-referencing and deduplication
7. ✅ Quality validation before storage

### What This Does NOT Cover
- ❌ Code implementation (use Work stage)
- ❌ Architecture design (use Plan stage)
- ❌ Code review (use Review stage)
- ❌ Initial project scaffolding (this bootstraps knowledge, not cod
fi
```

**2. If BOOTSTRAP MODE Required**:

Inform user:
```
🔍 Detected: This project lacks context/requirements structure.

I can bootstrap a complete knowledge base by:
1. Analyzing your existing codebase
2. Creating organized context/ and requirements/ directories
3. Extracting initial domain knowledge
4. Setting up index files for future learning

Would you like me to:
[A] Bootstrap full knowledge structure (recommended for existing projects)
[B] Create minimal structure only (just directories and READMEs)
[C] Skip bootstrap and proceed with learning synthesis

Please choose A, B, or C.
```

**STOP and wait for user choice.**

**3. Execute Bootstrap (if user chooses A or B)**

#### Option A: Full Bootstrap

**Step 1: Analyze Existing Project**

```bash
# Gather project context
- Read package.json / requirements.txt / go.mod (identify tech stack)
- Scan directory structure (understand organization)
- Check for existing docs (README, CONTRIBUTING, etc.)
- Identify test frameworks (Jest, pytest, Go test, etc.)
- Find build systems (webpack, vite, make, etc.)
```

**Step 2: Create Directory Structure**

```bash
mkdir -p context/tech
mkdir -p context/experience
mkdir -p requirements/in-progress
mkdir -p requirements/completed
```

**Step 3: Generate Initial Context Files**

Create `context/README.md`:
````markdown
# Project Knowledge Base

**Project**: ${workspaceFolder}  
**Initialized**: {Date}  
**Last Updated**: {Date}

## Overview

This knowledge base captures technical patterns, common issues, and project-specific experience to help team members work effectively.

## Structure

### Technical Knowledge (`tech/`)
System architecture, testing strategies, build configurations, and framework-specific patterns.

- [Testing Strategy](tech/testing-strategy.md) - TDD/E2E patterns
- [Build System](tech/build-system.md) - Build configurations and optimizations
- [Architecture](tech/architecture.md) - System design decisions
- [Frameworks](tech/frameworks.md) - Framework-specific knowledge

### Experience Knowledge (`experience/`)
Common errors, debugging tactics, workflow improvements, and hard-won lessons.

- [Error Patterns](experience/error-patterns.md) - Common errors and solutions
- [Test Patterns](experience/test-patterns.md) - Effective testing patterns
- [Debugging Tactics](experience/debugging-tactics.md) - Debugging strategies
- [Workflow Improvements](experience/workflow-improvements.md) - Process optimizations

## Quick Start

### Adding New Knowledge
Use `/learn` command to capture learnings:
```
/learn [your learning/pattern/solution]
```

### Finding Knowledge
1. Check this README for topic overview
2. Navigate to relevant domain file
3. Use search for specific patterns

## Statistics

- **Total Domain Files**: {count}
- **Technical Domains**: {count}
- **Experience Domains**: {count}
- **Last Learning Added**: {date}

---

**Maintenance**: Update this file when adding new domain areas or major changes.
````

Create `context/tech/README.md`:
````markdown
# Technical Knowledge

Framework configurations, architectural patterns, testing strategies, and build system knowledge.

## Domains

| Domain | Description | Key Patterns |
|--------|-------------|--------------|
| [Testing Strategy](testing-strategy.md) | TDD/E2E patterns, test organization | {summary from analysis} |
| [Build System](build-system.md) | Build configs, optimization | {summary from analysis} |
| [Architecture](architecture.md) | System design, patterns | {summary from analysis} |
| [Frameworks](frameworks.md) | Framework-specific knowledge | {detected: framework names} |

## Quick Reference

{Generate based on project analysis}

---
**Last Updated**: {Date}
````

Create `context/experience/README.md`:
````markdown
# Experience Knowledge

Lessons learned, common errors, debugging tactics, and workflow improvements discovered during development.

## Domains

| Domain | Description | Patterns |
|--------|-------------|----------|
| [Error Patterns](error-patterns.md) | Common errors and solutions | {count} patterns |
| [Test Patterns](test-patterns.md) | Effective testing patterns | {count} patterns |
| [Debugging Tactics](debugging-tactics.md) | Debugging strategies | {count} tactics |
| [Workflow Improvements](workflow-improvements.md) | Process optimizations | {count} improvements |

## Getting Started

Start with empty domains. As you encounter issues and discover solutions, use `/learn` to capture them.

---
**Last Updated**: {Date}
````

**Step 4: Create Initial Domain Files** (based on project analysis)

For each domain, create initial file with discovered patterns:

Example `context/tech/testing-strategy.md`:
````markdown
# Testing Strategy

**Last Updated**: {Date}

## Overview

Testing approach for this project based on {detected test framework}.

## Current Setup

**Test Framework**: {detected: Jest/pytest/etc.}  
**Test Organization**: {analyzed from project structure}  
**Coverage Tool**: {detected if present}

## Test Patterns

{If tests exist, extract 1-2 initial patterns}

### Pattern: {Discovered Pattern Name}

**Context**: {When this is used in the codebase}
**Implementation**: 
```code
{Example from actual codebase}
```
**Benefits**: {Why this works}

## Getting Started

This file will grow as you:
- Discover effective test patterns
- Encounter testing challenges
- Learn testing best practices

Use `/learn` to add new testing knowledge.

---
**Change Log**:
| Date | Change | Author |
|------|--------|--------|
| {Date} | Initial bootstrap from project analysis | System |
````

**Step 5: Generate Requirements Index**

Create `requirements/INDEX.md`:
````markdown
# Requirements Index

**Project**: ${workspaceFolder}  
**Last Updated**: {Date}

## Overview

This index tracks all requirements for the project. Use `/plan` to create new requirements.

## In Progress

{If GitHub issues exist, suggest importing them}

Currently tracking requirements in progress.

| Req ID | Title | Priority | Created | Owner | Status |
|--------|-------|----------|---------|-------|--------|
| - | - | - | - | - | - |

## Completed

Requirements that have been finalized and implemented.

| Req ID | Title | Priority | Completed | Owner | GitHub Issue |
|--------|-------|----------|-----------|-------|--------------|
| - | - | - | - | - | - |

## Statistics

- **Total Requirements**: 0
- **In Progress**: 0
- **Completed**: 0

## Getting Started

To create your first requirement, use:
```
/plan [describe your feature or project]
```

This will create a structured requirement document following professional software engineering practices.

## Integration

- **Planning**: Use `/plan` to create requirements
- **Learning**: Use `/learn` to capture implementation insights
- **Review**: Link code reviews back to requirements

---

**Next Steps**: 
1. Use `/plan` to document your first requirement
2. Link GitHub issues to requirement documents
3. Use `/learn` to capture implementation learnings
````

**Step 6: Generate Bootstrap Summary**

````markdown
## ✅ Knowledge Base Bootstrapped

### Created Structure

```
context/
├── README.md (✓ Created)
├── tech/
│   ├── README.md (✓ Created)
│   ├── testing-strategy.md (✓ Created with {n} initial patterns)
│   ├── build-system.md (✓ Created)
│   ├── architecture.md (✓ Created)
│   └── frameworks.md (✓ Created - {detected frameworks})
└── experience/
    ├── README.md (✓ Created)
    ├── error-patterns.md (✓ Created - ready for learnings)
    ├── test-patterns.md (✓ Created - ready for learnings)
    ├── debugging-tactics.md (✓ Created - ready for learnings)
    └── workflow-improvements.md (✓ Created - ready for learnings)

requirements/
├── INDEX.md (✓ Created)
├── in-progress/ (✓ Created)
└── completed/ (✓ Created)
```

### Discovered & Documented

- **Tech Stack**: {list detected technologies}
- **Test Framework**: {detected framework}
- **Build System**: {detected build tool}
- **Initial Patterns**: {count} patterns extracted from codebase

### Recommended Next Steps

1. **Review Initial Knowledge**: 
   - Check `context/tech/testing-strategy.md` for accuracy
   - Verify detected frameworks in `context/tech/frameworks.md`

2. **Add Your First Learning**:
   ```
   /learn [describe a pattern or issue you've encountered]
   ```

3. **Document Requirements**:
   ```
   /plan [describe a feature or project goal]
   ```

4. **Import Existing Issues** (if applicable):
   - Review GitHub issues
   - Convert to formal requirements using `/plan`

### Your Knowledge Base is Ready! 🎉

The structure is in place. Now use:
- `/learn` - Capture new patterns and solutions
- `/plan` - Document requirements
- Context files - Reference during development
````

#### Option B: Minimal Bootstrap

Create only directory structure and basic README files without analysis:

```bash
mkdir -p context/tech context/experience
mkdir -p requirements/in-progress requirements/completed

# Create minimal README files with templates
# Create empty domain files
```

**4. After Bootstrap Complete**

Set mode to `LEARNING_SYNTHESIS` and inform user:

```
Knowledge base initialized! You can now:
- Add learnings: /learn [pattern/solution]
- Create requirements: /plan [feature description]
- Browse knowledge: context/README.md
```
### What This Workflow Covers
1. ✅ Intelligent context loading (index-first strategy)
2. ✅ Learning analysis and categorization
3. ✅ Domain knowledge synthesis
4. ✅ Skill creation for complex workflows
5. ✅ Incremental document updates (no pollution)
6. ✅ Cross-referencing and deduplication
7. ✅ Quality validation before storage

### What This Does NOT Cover
- ❌ Code implementation (use Work stage)
- ❌ Architecture design (use Plan stage)
- ❌ Code review (use Review stage)

## Workflow

Use the **todo list** to track progress through all learning stages.

### Phase 1: Intelligent Context Loading

**Objective**: Load ONLY relevant context for the current learning synthesis.

**Index-First Strategy**:

1. **Load Primary Indexes** (read these first, always)
   ```
   context/README.md
   requirements/INDEX.md
   ```

2. **Analyze Learning Domain** (from user input)
   - Identify if it's technical or experiential
   - Determine specific domain (testing, build, debugging, etc.)
   - Note any requirement references

3. **Selectively Load Domain Context** (only what's needed)
   
   For technical learnings:
   ```
   context/tech/README.md
   context/tech/<relevant-domain>.md  # Only the specific domain
   ```
   
   For experiential learnings:
   ```
   context/experience/README.md
   context/experience/<relevant-domain>.md  # Only the specific domain
   ```

4. **Load Related Requirements** (if applicable)
   - Only if learning relates to specific requirements
   - Load only the specific requirement file, not all requirements
   - Use INDEX.md to locate the correct file

**Context Loading Rules**:
- ❌ Never load entire directories
- ❌ Never load documents "just in case"
- ✅ Load index files to understand structure
- ✅ Load specific domain files based on learning category
- ✅ Use grep/search to find specific patterns before loading files
- ✅ Load related requirements only when explicitly referenced

### Phase 2: Learning Analysis & Validation

**Objective**: Understand what was learned and validate its value.

1. **Extract the Core Learning**
   - What was the problem or challenge?
   - What was the solution or insight?
   - Why does this matter for future work?
   - What context is needed to understand it?

2. **Validate Repeatability**
   - Is this a one-time fix or a reusable pattern?
   - Will this apply to future situations?
   - Are there generalizable principles?
   - Does this contradict existing knowledge?

3. **Check for Existing Knowledge**
   - Search loaded context for similar patterns
   - Identify if this enhances or replaces existing guidance
   - Find opportunities to merge or consolidate
   - Avoid redundancy with existing documentation

4. **Determine Documentation Level**
   
   | Learning Type | Action Required |
   |--------------|----------------|
   | Simple pattern/tip | Add to existing domain knowledge file |
   | Complex workflow | Create new skill with scripts/templates |
   | Architecture decision | Update architecture documentation + create ADR |
   | Error pattern | Add to error-patterns.md with solution |
   | Process improvement | Update workflow documentation |

### Phase 3: Domain Knowledge Synthesis

**Objective**: Update domain knowledge incrementally without pollution.

**For Simple Learnings** (adding to existing files):

1. **Read the Target File**
   - Load the specific domain file (e.g., `error-patterns.md`)
   - Identify the appropriate section
   - Check for similar existing entries

2. **Craft the Update** (incremental, not rewrite)
   ```markdown
   ## [Category Name]
   
   ### Existing Pattern 1
   [Existing content - DO NOT MODIFY unless merging]
   
   ### Existing Pattern 2
   [Existing content - DO NOT MODIFY unless merging]
   
   ### [NEW] Your New Pattern  ← ADD THIS ONLY
   **Problem**: [Concise problem description]
   **Solution**: [Actionable solution]
   **Example**: 
   ```code
   [Working example]
   ```
   **Context**: [When this applies]
   **References**: [Links to relevant requirements/docs]
   ```

3. **Apply Surgical Update**
   - Use `replace_string_in_file` to insert new section
   - Include sufficient context (3-5 lines before/after)
   - Update table of contents if file has one
   - Add cross-references to related patterns

4. **Update Domain Index**
   - Update `context/tech/README.md` or `context/experience/README.md`
   - Add one-line summary of new knowledge
   - Maintain index structure

**For Complex Learnings** (creating new domain files):

1. **Create Structured Domain File**
   
   Template for new domain knowledge file:
   
   ````markdown
   # {Domain Name}
   
   **Last Updated**: {Date}
   **Maintained By**: {Team/Person}
   
   ## Overview
   
   [One paragraph: what this domain covers and why it matters]
   
   ## Table of Contents
   
   - [Quick Reference](#quick-reference)
   - [Patterns](#patterns)
   - [Anti-Patterns](#anti-patterns)
   - [Examples](#examples)
   - [Related Resources](#related-resources)
   
   ## Quick Reference
   
   | Scenario | Solution | Link |
   |----------|----------|------|
   | [Common scenario 1] | [Quick solution] | [Link to detail] |
   
   ## Patterns
   
   ### Pattern 1: {Name}
   
   **Context**: [When to use this]
   **Solution**: [How to implement]
   **Example**:
   ```code
   [Working example]
   ```
   **Benefits**: [Why this works]
   **Gotchas**: [What to watch out for]
   
   ## Anti-Patterns
   
   ### Anti-Pattern 1: {Name}
   
   **Problem**: [What's wrong with this approach]
   **Why It Fails**: [Explanation]
   **Better Alternative**: [Link to correct pattern]
   
   ## Examples
   
   ### Example 1: {Scenario}
   
   [Complete, working example with context]
   
   ## Related Resources
   
   - [Relate0: Bootstrap Existing Project (New Scenario)

**User Input**: 
```
/learn
```

**System Detection**:
```
🔍 Checking project structure...
❌ context/ directory not found
❌ requirements/ directory not found

MODE: BOOTSTRAP
```

**User Interaction**:
```
System: "This project lacks knowledge structure. Bootstrap now?"
User: "A" (full bootstrap)
```

**Actions Taken**:
1. ✅ Analyzed project structure
2. ✅ Detected: Node.js + TypeScript + Jest + Vite
3. ✅ Created context/ with 8 domain files
4. ✅ Extracted 3 initial testing patterns from existing tests
5. ✅ Created requirements/INDEX.md
6. ✅ Generated bootstrap summary

**Files Created**:
- context/README.md
- context/tech/README.md (with 4 domains)
- context/tech/testing-strategy.md (with 3 discovered patterns)
- context/tech/build-system.md (Vite configuration patterns)
- context/tech/frameworks.md (React + TypeScript patterns)
- context/experience/README.md (with 4 domains)
- context/experience/error-patterns.md (ready for learnings)
- requirements/INDEX.md

**Output**:
```
✅ Knowledge base bootstrapped!

Discovered:
- Tech Stack: Node.js, TypeScript, React
- Test Framework: Jest with React Testing Library
- Build System: Vite
- Initial Patterns: 3 testing patterns extracted

Next steps:
- Review context/tech/testing-strategy.md
- Use /learn to add new patterns
- Use /plan to document requirements
```

### Example d domain file 1](./other-domain.md)
   - [Relevant skill](../../.github/skills/skill-name/SKILL.md)
   - [External documentation](https://example.com)
   
   ---
   
   **Change Log**:
   | Date | Change | Author |
   |------|--------|--------|
   | {Date} | Initial creation | {Name} |
   ````

2. **Update Context README**
   - Add entry to appropriate section (tech or experience)
   - Link to new domain file
   - Update statistics

### Phase 4: Skill Creation (When Needed)

**Objective**: Create reusable skills for complex, repeatable workflows.

**When to Create a Skill** (not just documentation):

- ✅ Workflow requires multiple tools/commands
- ✅ Pattern includes scripts, templates, or resources
- ✅ Process has 3+ steps that need bundling
- ✅ Workflow is triggered by specific user requests
- ✅ Solution needs executable scripts or generators

**Skill Creation Process**:

1. **Read Skill Creation Instructions**
   ```
   .github/instructions/agent-skills.instructions.md
   .github/skills/skill-creator/SKILL.md
   ```

2. **Invoke Skill Creator**
   - Use skill-creator skill for guidance
   - Follow professional skill structure
   - Include proper frontmatter with triggers

3. **Skill Structure Requirements**
   
   ```
   .github/skills/{skill-name}/
   ├── SKILL.md                 # Main skill instructions
   ├── scripts/                 # Executable scripts (if needed)
   │   └── {script-name}.sh
   ├── templates/               # File templates (if needed)
   │   └── {template-name}.md
   ├── references/              # Reference materials (if needed)
   │   └── {reference-name}.md
   └── examples/                # Working examples (if needed)
       └── {example-name}/
   ```

4. **Skill Frontmatter Requirements**
   
   ```yaml
   ---
   name: {skill-name}
   description: '{What it does}. Use when {specific triggers/scenarios/keywords}. Handles {capabilities}'
   license: MIT
   ---
   ```
   
   **Critical**: Description MUST include:
   - What the skill does
   - When to use it (triggers)
   - Keywords users might mention
   - Specific capabilities

5. **Skill Body Structure**
   
   ```markdown
   # {Skill Name}
   
   ## Overview
   [One paragraph: purpose and value]
   
   ## When to Use
   [Specific scenarios, triggers, file types]
   
   ## Prerequisites
   [Required tools, files, context]
   
   ## Workflow
   
   ### Step 1: {Action}
   [Detailed instructions]
   
   ### Step 2: {Action}
   [Detailed instructions]
   
   ## Examples
   
   ### Example 1: {Scenario}
   ```bash
   # Command or code example
   ```
   
   ## Troubleshooting
   
   | Issue | Cause | Solution |
   |-------|-------|----------|
   | [Common issue] | [Root cause] | [How to fix] |
   
   ## Related Resources
   - [Domain knowledge](../../../context/tech/domain.md)
   - [Other skill](../other-skill/SKILL.md)
   ```

6. **Test the Skill**
   - Verify it activates on appropriate triggers
   - Test with example scenarios
   - Validate scripts execute correctly
   - Ensure templates generate properly

7. **Document in Context**
   - Reference the skill in relevant domain knowledge
   - Add to context README if significant
   - Cross-link from related documentation

### Phase 5: Requirement Linkage (When Applicable)

**Objective**: Link learnings to specific requirements for traceability.

**Only if learning relates to a requirement**:

1. **Identify Related Requirements**
   - Use requirements/INDEX.md to find relevant requirements
   - Load only the specific requirement files
   - Check if learning impacts requirement

2. **Update Requirement Documentation**
   
   Add learning to relevant requirement file section:
   
   ```markdown
   ## Implementation Learnings
   
   ### {Learning Title} - {Date}
   
   **Context**: [What was being implemented]
   **Discovery**: [What was learned]
   **Impact**: [How this affects the requirement]
   **Action Taken**: [Changes made to code/docs/tests]
   **Knowledge Captured**: [Link to domain knowledge or skill]
   ```

3. **Update Requirement Status** (if needed)
   - Mark as updated in INDEX.md
   - Add changelog entry
   - Update completion percentage if applicable

### Phase 6: Quality Validation & Cleanup

**Objective**: Ensure knowledge is high-quality and accessible.

**Validation Checklist**:

- [ ] **Clarity**: Is the learning clearly explained?
- [ ] **Actionability**: Can someone use this guidance directly?
- [ ] **Examples**: Are there working, tested examples?
- [ ] **Context**: Is it clear when/why to apply this?
- [ ] **Non-Redundancy**: Doesn't duplicate existing knowledge
- [ ] **Linkage**: Properly cross-referenced with related docs
- [ ] **Scannability**: Easy to find and quick to understand
- [ ] **Minimal Pollution**: Only added/updated what's necessary

**Cleanup Tasks**:

1. **Remove Redundancy**
   - If new learning overlaps existing knowledge, merge them
   - Consolidate similar patterns into one comprehensive pattern
   - Remove outdated information that's been superseded

2. **Update Cross-References**
   - Add bidirectional links between related knowledge
   - Update "Related Resources" sections
   - Link from index files to new content

3. **Verify File Integrity**
   - Check markdown formatting is correct
   - Verify code examples syntax
   - Test that links resolve correctly
   - Ensure tables are properly formatted

4. **Update Metrics** (in README files)
   - Increment pattern counts
   - Update "Last Updated" dates
   - Add to change logs

### Phase 7: Synthesis Summary

**Objective**: Provide clear summary of what was learned and where it's stored.

**Generate Summary**:

```markdown
## Learning Synthesis Complete

### What Was Learned
[Concise summary of the core learning]

### Where It's Stored
- **Domain Knowledge**: [Link to updated domain file(s)]
- **Skills Created**: [Links to new skills, if any]
- **Requirements Updated**: [Links to requirement files, if any]

### Key Takeaways
1. [Actionable takeaway 1]
2. [Actionable takeaway 2]
3. [Actionable takeaway 3]

### Recommended Next Steps
- [Suggestion for applying this learning]
- [Related patterns to explore]
- [Potential improvements to consider]

### Quick Access
Use these commands/searches to access this knowledge:
- Search: `{relevant keywords}`
- Skill: `{skill-name}` (if created)
- Document: `context/{domain}/{file}.md`
```

## Output Expectations

### Primary Outputs

1. **Updated Domain Knowledge Files**
   - Incremental updates to existing files
   - OR new domain files following template
   - Clean, non-polluted additions

2. **New Skills** (when warranted)
   - Complete skill directory with SKILL.md
   - Scripts, templates, references as needed
   - Proper frontmatter for discoverability

3. **Updated Index Files**
   - context/README.md
   - context/{tech|experience}/README.md
   - requirements/INDEX.md (if applicable)

4. **Synthesis Summary**
   - Clear documentation of what was learned
   - Locations of all updates
   - Recommended usage

### Success Criteria

- ✅ Learning is captured in appropriate domain
- ✅ Only relevant context was loaded (no waste)
- ✅ Updates are surgical, not wholesale rewrites
- ✅ Skills created for complex workflows
- ✅ Cross-references are bidirectional
- ✅ Examples are tested and working
- ✅ Documentation is scannable and actionable

### Failure Conditions

- ❌ Loaded too much unnecessary context
- ❌ Created redundant documentation
- ❌ Polluted existing files with full rewrites
- ❌ Failed to create skill for complex workflow
- ❌ Learning is too vague to be actionable
- ❌ Missing links to related resources

## Quality Assurance

### Context Loading Validation

```bash
# Verify minimal context loading
# Should load < 10 files for typical learning synthesis

# Check index files exist
ls -la context/README.md
ls -la requirements/INDEX.md

# Verify domain structure
tree context/
```

### Documentation Quality Check

**10/10 Quality Criteria**:

1. **Zero Knowledge Loss**: Every detail captured
2. **Maximum Scannability**: Clear headings, tables, bullets
3. **Minimal Redundancy**: No duplicate content
4. **Actionable Examples**: Working, tested code
5. **Clear Context**: When/why to use
6. **Proper Cross-refs**: Bidirectional links
7. **Surgical Updates**: No wholesale rewrites
8. **Domain Separation**: Right content in right file
9. **Progressive Loading**: Index → specific docs
10. **Verified Accuracy**: Examples tested

### Skill Quality Check (if created)

- [ ] Frontmatter has descriptive triggers
- [ ] Description includes what/when/keywords
- [ ] Workflow is step-by-step
- [ ] Examples are complete and tested
- [ ] Scripts execute successfully
- [ ] Templates generate correctly
- [ ] Activates on appropriate user prompts

## Integration with Other Stages

### From Review Stage
- Security findings → experience/error-patterns.md
- Refactoring learnings → tech/architecture.md
- Code quality insights → experience/test-patterns.md

### From Work Stage
- Build errors → experience/error-patterns.md
- Test strategies → tech/testing-strategy.md
- Framework patterns → tech/frameworks.md

### From Plan Stage
- Architecture decisions → tech/architecture.md
- Requirement learnings → link to specific requirements
- Design patterns → tech/architecture.md

### To Future Work
- Skills available for reuse
- Patterns prevent repeated mistakes
- Context enables faster onboarding
- Knowledge compounds over time

## Advanced Patterns

### Session Reflection Pattern

For end-of-session synthesis:

1. **Review session context**
2. **Identify 3-5 key learnings**
3. **Synthesize each learning individually**
4. **Generate session summary document** (optional)
   ```
   context/sessions/YYYY-MM-DD-session-summary.md
   ```

### Continuous Learning Loop

Integrate with session hooks:

1. **SessionStart Hook**: Load relevant prior learnings
2. **Work Phase**: Accumulate learnings during work
3. **Stop Hook**: Trigger synthesis at session end
4. **Next Session**: Prior learnings inform new work

### Knowledge Compounding

Build reusable patterns:

```
Learning → Domain Knowledge → Skill → Best Practice
   ↑                                      ↓
   └──────── Continuous Improvement ──────┘
```

## Examples

### Example 1: Error Pattern Learning

**User Input**: 
```
/learn We kept getting "EADDRINUSE" errors in tests because 
Jest wasn't cleaning up server instances between test files
```

**Context Loading**:
- ✅ Load: context/README.md
- ✅ Load: context/experience/README.md
- ✅ Load: context/experience/error-patterns.md
- ❌ Skip: All other files (not relevant)

**Update Applied** (surgical insertion):
```markdown
### EADDRINUSE - Port Already in Use

**Problem**: Tests fail with "EADDRINUSE" error
**Root Cause**: Server instances not cleaned up between test files
**Solution**: Add global teardown in Jest configuration

```javascript
// jest.config.js
module.exports = {
  globalTeardown: './tests/teardown.js'
}

// tests/teardown.js
module.exports = async () => {
  // Close all server instances
  if (global.__SERVER__) {
    await global.__SERVER__.close()
  }
}
```

**Context**: Affects integration tests that spawn servers
**Related**: [Testing Strategy](../tech/testing-strategy.md#integration-tests)
```

### Example 2: Complex Workflow → Skill Creation

**User Input**:
```
/learn We have a complex process for creating ADRs that involves 
templates, numbering, and linking. Should be automated.
```

**Actions Taken**:
1. Loaded context/README.md and context/tech/README.md
2. Determined this needs a skill (multi-step workflow)
3. Created `.github/skills/adr-creator/SKILL.md`
4. Added template in `scripts/adr-template.md`
5. Updated context/tech/architecture.md with reference
6. Tested skill activation with "create an ADR"

**Skill Created**:
```
.github/skills/adr-creator/
├── SKILL.md (workflow for ADR creation)
├── templates/
│   └── adr-template.md
└── scripts/
    └── generate-adr-number.sh
```

## Troubleshooting

| Issue | Cause | Solution |
|-------|-------|----------|
| Loaded too many files | Didn't check indexes first | Always load README/INDEX first, then selectively |
| Duplicate content | Didn't search existing knowledge | Always search before adding new content |
| Polluted documents | Full file rewrite instead of surgical update | Use replace_string_in_file with precise context |
| Skill not activating | Poor description triggers | Add specific keywords and scenarios to description |
| Learning too vague | Insufficient analysis | Ask clarifying questions in Phase 2 |
| Missing cross-references | Didn't check related docs | Review related domain files and add links |

## Related Resources

### Skills Referenced
- **skill-creator**: Guide for creating effective skills
- **make-skill-template**: Scaffold new skills

### Instructions Referenced
- **agent-skills.instructions.md**: Skill creation guidelines
- **memory-bank.instructions.md**: Memory management system

### Prompts Referenced
- **remember**: Quick memory capture (global/workspace)
- **memory-merger**: Consolidate mature learnings

---

**Version**: 1.0  
**Last Updated**: 2026-01-30  
**Maintained By**: GitHub Copilot Code Repository  
**Based On**: [The Longform Guide to Everything Claude Code](https://github.com/affaan-m/everything-claude-code/blob/main/the-longform-guide.md)
