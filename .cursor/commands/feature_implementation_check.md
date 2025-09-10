# Feature Implementation Check

## Overview
Systematically verify that a feature implementation meets requirements, follows best practices, and is production-ready.

## Instructions for AI
When the user provides a feature name and description, analyze the codebase to verify the implementation against this checklist. For each item, provide specific findings and recommendations.

## Analysis Process

### Phase 1: Understanding the Feature
1. **Parse User Input**
   - Extract feature name and requirements from user description
   - Identify key acceptance criteria mentioned
   - Note any specific concerns or focus areas

2. **Locate Implementation**
   - Search for relevant files using feature name/keywords
   - Identify main components, modules, or services
   - Map out the feature's code structure

### Phase 2: Implementation Review

#### 1. **Requirements Verification** ✅
   - [ ] All stated requirements are implemented
   - [ ] Feature behaves as described
   - [ ] Edge cases from requirements are handled
   - [ ] No missing functionality

#### 2. **Code Quality Analysis** 💻
   - [ ] **Structure**: Follows project architecture patterns
   - [ ] **Readability**: Clear naming, proper formatting
   - [ ] **Maintainability**: No code duplication, proper abstraction
   - [ ] **Error Handling**: Comprehensive error catching and logging
   - [ ] **Input Validation**: All inputs validated and sanitized
   - [ ] **Dependencies**: Minimal and necessary dependencies

#### 3. **Testing Coverage** 🧪
   - [ ] **Unit Tests**: Core logic has unit tests
   - [ ] **Integration Tests**: API/service interactions tested
   - [ ] **Edge Cases**: Boundary conditions tested
   - [ ] **Error Scenarios**: Failure cases tested
   - [ ] **Test Quality**: Tests are clear and maintainable

#### 4. **Documentation Status** 📚
   - [ ] **Code Comments**: Complex logic explained
   - [ ] **API Documentation**: Endpoints/methods documented
   - [ ] **README Updates**: Feature usage explained
   - [ ] **Configuration**: Settings documented

#### 5. **Security Check** 🔒
   - [ ] **Authentication**: Proper access controls
   - [ ] **Data Protection**: Sensitive data handled securely
   - [ ] **Input Security**: Protection against injection attacks
   - [ ] **Dependencies**: No known vulnerabilities

#### 6. **Performance Considerations** ⚡
   - [ ] **Efficiency**: No obvious performance bottlenecks
   - [ ] **Database**: Optimized queries, proper indexing
   - [ ] **Caching**: Implemented where beneficial
   - [ ] **Resource Usage**: Memory/CPU usage reasonable

#### 7. **Integration & Deployment** 🚀
   - [ ] **Backward Compatibility**: Existing functionality preserved
   - [ ] **Configuration**: Environment variables properly set
   - [ ] **Migrations**: Database changes handled correctly
   - [ ] **Monitoring**: Logs and metrics in place

### Phase 3: Report Generation

## Report Template

```markdown
# Feature Implementation Check Report

**Feature**: [Feature Name]
**Status**: [✅ Complete | ⚠️ Needs Attention | ❌ Critical Issues]

## Executive Summary
[Brief overview of findings]

## Detailed Findings

### ✅ What's Working Well
- [List positive findings]

### ⚠️ Areas for Improvement
- [List non-critical issues with recommendations]

### ❌ Critical Issues
- [List blocking issues that must be fixed]

## Checklist Results

| Category | Status | Notes |
|----------|--------|-------|
| Requirements | ✅/⚠️/❌ | [Specific findings] |
| Code Quality | ✅/⚠️/❌ | [Specific findings] |
| Testing | ✅/⚠️/❌ | [Specific findings] |
| Documentation | ✅/⚠️/❌ | [Specific findings] |
| Security | ✅/⚠️/❌ | [Specific findings] |
| Performance | ✅/⚠️/❌ | [Specific findings] |
| Deployment | ✅/⚠️/❌ | [Specific findings] |

## Recommended Actions

### Immediate (Before Deployment)
1. [Critical fixes]

### Short-term (Next Sprint)
1. [Important improvements]

### Long-term (Technical Debt)
1. [Nice-to-have enhancements]

## Code Examples
[Include specific code snippets showing issues and fixes]

## Conclusion
[Final assessment and go/no-go recommendation]
```

## AI Behavior Guidelines

1. **Be Specific**: Reference actual file names and line numbers
2. **Prioritize**: Focus on critical issues first
3. **Be Constructive**: Provide solutions, not just problems
4. **Consider Context**: Respect existing patterns in the codebase
5. **Be Concise**: Focus on actionable findings

## Example Usage

**User**: "Check the implementation of the user authentication feature that handles login with email/password and stores sessions in Redis"
