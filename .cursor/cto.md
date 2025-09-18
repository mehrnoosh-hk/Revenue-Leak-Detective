## CTO/Product Strategy Mentor Prompt

### Role

Act as my CTO and product strategy mentor.  
Your role is to help me plan, prioritize, and align engineering work with business goals, just like a real CTO or senior PM would.

---

### Task

- Plan the next software development step based on the current state of the project.
- Break that plan into concrete, incremental, and testable tasks, with dependencies called out.

---

### Context

#### Idea

Revenue Leak Detective is an agent that hunts down money leaks in a SaaS: failed charges, paused subscriptions, coupon misuse, “trial forever” zombies, and quiet churn signals (no logins + downgrades).  
It triages issues, suggests fixes, drafts customer outreach, and files tasks automatically.

#### Current State

- Basic scaffolding is done.
- Backend: Golang server, sqlc, migrate, Postgres
- AI-Workers: Python, Langchain, LangGraph

---

### References

- [Road_Map.md]
- [Road_Map.txt]
- [Product_Idea.txt]
- Project repository (current state)

---

### Constraints

- Use provided references if you need more information to complete your task.
- We are in MVP stage.
- Core problem-solving features only.
- Leave “nice-to-haves” for later.
- Ensure scalability, maintainability, and observability from early stages.
- Do not overcomplicate with corporate jargon.
- Be pragmatic, focus on getting to market fast without painting us into a corner.
- If ambiguous, ask up to 5 targeted questions, then proceed with sensible defaults.
- Always tie advice to how a CTO or senior PM would think in the real world.
- Use lean startup principles: validate, iterate, scale.

---

### Priority

- Prioritize vertical slices (end-to-end usable features)
- Balance quick wins vs foundational work

---

### Output Format

- Plan
- Assumptions
- Highlight implications
- Tasks

---

### Output

- Plan for the next steps
- Assumptions and context (what you inferred about the product stage, customers, or goals)
- Highlight implications on scalability, maintainability, velocity
- Tasks with comprehensive details and evaluation rubric

---

### Checks

- Did I follow all constraints?
- Is the format correct?

---

## 🎯 **Cursor Command Template for Revenue Leak Detective**

### **Context Setting Commands**

```bash
<code_block_to_apply_changes_from>
```

### **Development Workflow Commands**

#### **Database & Schema Work**
```bash
# When working on database changes
/context I'm working on database schema for Revenue Leak Detective. Current tables: users, tenants, customers, leaks, actions. Using postgres with numbered migrations (001_, 002_, etc.) and SQLC for code generation. Check existing schema in services/api/migrations/ and generated code in services/api/internal/db/sqlc/

# Generate domain models after schema changes
Please review the current database schema in services/api/migrations/ and update the domain models in services/api/internal/domain/models/generated_models.go to match the latest schema. Follow Go conventions and include proper validation tags.

# Review migration files
Review all migration files in services/api/migrations/ and ensure they follow best practices for postgres schema design, proper indexing, and foreign key relationships. Check for any potential issues or improvements.
```

#### **API Development**
```bash
# When working on API endpoints
/context Working on Go API endpoints following Clean Architecture. Structure: handlers/ receive HTTP requests → call domain services in internal/domain/ → use repositories in internal/db/repository/ → SQLC generated queries. All responses should use proper HTTP status codes and structured JSON.

# Add new endpoint
I need to add a new API endpoint for [specific functionality]. Please follow the existing pattern: create handler in handlers/, add domain service if needed, update repository if database access required, and include proper error handling and tests.

# Review API structure
Please review the current API structure in services/api/ and suggest improvements for [specific area]. Focus on Clean Architecture principles, error handling, and Go best practices.
```

#### **Testing & Quality**
```bash
# For test development
/context This project emphasizes testing with comprehensive test coverage. API tests use Go testing package with mocks for repositories. Python tests use pytest. All new code should include unit tests following existing patterns in *_test.go files.

# Fix failing tests
There are failing tests in the codebase. Please identify and fix them while maintaining the existing test patterns and ensuring good coverage. Run make api-test to check status.

# Improve test coverage
Review the test coverage for [specific component] and add missing tests. Follow the existing patterns using interfaces and mocks for proper unit testing.
```

### **Specific Domain Commands**

#### **Revenue Leak Detection Logic**
```bash
# Core business logic
/context The core business domain focuses on detecting revenue leaks: failed charges (Stripe integration), quiet churn (no logins + downgrades), coupon misuse, and trial zombies. The system should detect → suggest → draft tasks → human approval workflow.

# Implement leak detection
I need to implement the core leak detection logic for [specific leak type]. This should follow the domain-driven design pattern, include proper error handling, and integrate with the existing repository pattern.
```

#### **Integration Work**
```bash
# External integrations
/context Revenue Leak Detective integrates with external services: Stripe (billing), Slack (notifications), Linear (task management). Integration code should be in separate packages with proper error handling and retry logic.

# Add new integration
I need to add integration with [service name]. Please create a new integration package following the existing patterns, include proper configuration, error handling, and tests.
```

### **Build & Deployment Commands**

```bash
# Build issues
/context This project uses Make for build automation. Check Makefile and make/*.mk files for available targets. Common commands: make deps (install dependencies), make api-test (run tests), make sqlc-generate (regenerate database code), make migrate-up (run migrations).

# Docker work
/context The project uses Docker with separate Dockerfiles for API (Dockerfile.api) and workers (Dockerfile.workers). Uses multi-stage builds for optimization. Check deploy/docker/ directory for Docker configurations.

# Fix build errors
There are build errors in the project. Please identify and fix them, ensuring all dependencies are properly managed and the build process follows the existing Makefile patterns.
```

### **Code Review & Refactoring Commands**

```bash
# Code review
Please review the code in [specific file/directory] for adherence to Go best practices, Clean Architecture principles, and project conventions. Focus on error handling, separation of concerns, and maintainability.

# Refactoring
The code in [specific area] needs refactoring to better follow Clean Architecture principles. Please suggest improvements while maintaining backward compatibility and existing test coverage.

# Performance optimization
Review [specific component] for potential performance improvements, focusing on database queries, memory usage, and Go-specific optimizations.
```

### **Documentation Commands**

```bash
# API documentation
Please generate/update API documentation for the endpoints in handlers/. Include request/response examples, error codes, and follow OpenAPI standards if possible.

# Code documentation
Add proper Go documentation comments to [specific file/package] following Go conventions. Include package overview, function descriptions, and example usage where appropriate.
```

### **Quick Action Templates**

```bash
# Quick debugging
There's an issue with [specific functionality]. Please investigate the problem, identify the root cause, and suggest a fix while maintaining the existing architecture patterns.

# Feature implementation
I need to implement [specific feature] following the existing Clean Architecture pattern. Please provide a complete implementation including handlers, domain logic, repository layer, tests, and any necessary database changes.

# Configuration update
Please update the configuration system in services/api/config/ to support [new requirement]. Ensure proper validation, environment variable handling, and documentation.
```

### **Multi-Service Commands**

```bash
# Go + Python coordination
/context This is a multi-service project with Go API and Python workers. Changes often require coordination between services. API handles HTTP requests and database operations, workers handle background processing and external integrations.

# Cross-service feature
I need to implement [feature] that spans both Go API and Python workers. Please provide the implementation for both services, ensuring proper communication patterns and error handling.
```

---

## 🔧 **Usage Tips**

1. **Start with context**: Always begin with `/context` to set the stage
2. **Be specific**: Mention specific files, directories, or components when possible
3. **Reference existing patterns**: Ask Cursor to follow existing code patterns
4. **Include constraints**: Mention architecture principles, testing requirements, etc.
5. **Multi-step requests**: Break complex requests into focused sub-tasks

This template is designed specifically for your Revenue Leak Detective project structure and will help Cursor understand your codebase context and provide more relevant assistance.