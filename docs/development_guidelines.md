# InstaPoll Development Guidelines

This document outlines the best practices, coding standards, and workflows to follow when developing the InstaPoll application. Adhering to these guidelines ensures code consistency, maintainability, and collaboration efficiency.

## 1. Code Style & Formatting

* **Go (Backend):**
    * Follow standard Go formatting (`gofmt`). Most IDEs can be configured to do this automatically on save.
    * Use `go vet` and linters (like `golangci-lint`) to catch potential issues early.
    * Keep functions concise and focused on a single responsibility.
    * Comment code where the logic is non-obvious. Explain *why*, not just *what*.
* **JavaScript/TypeScript (Frontend):**
    * Use a consistent code formatter (e.g., Prettier).
    * Follow established style guides (e.g., Airbnb JavaScript Style Guide, TypeScript recommendations).
    * Use linters (e.g., ESLint) to enforce style and catch errors.
* **General:**
    * Use meaningful names for variables, functions, and types.
    * Avoid magic numbers and strings; use constants or configuration values.

## 2. Testing

* **Test-Driven Development (TDD):**
    * **Always strive to use TDD where practical, especially for backend logic.**
    * Write tests *before* writing the implementation code.
    * Start with a failing test, write the minimum code to make it pass, then refactor.
* **Unit Tests:**
    * All core logic, utility functions, and non-trivial components should have unit tests.
    * Aim for good test coverage, focusing on critical paths and edge cases.
    * Use Go's built-in `testing` package for the backend.
    * Use standard testing libraries (e.g., Jest, React Testing Library) for the frontend.
* **Integration Tests:**
    * Write integration tests for interactions between components (e.g., API handler interacting with the database layer).
    * These may require setting up test databases or mocking external services.
* **End-to-End (E2E) Tests:** (Future Goal)
    * Implement E2E tests to simulate user flows through the entire application.

## 3. Database (MongoDB)

* **Schema Design:** Define clear schemas/structs for your MongoDB documents (`models/poll.go`).
* **Driver Usage:** Use the official `go.mongodb.org/mongo-driver` for backend interactions.
* **Connection Management:** Handle database connections properly (e.g., establish connection once at startup).
* **Error Handling:** Implement robust error handling for database operations.
* **Indexes:** Define necessary indexes to optimize query performance.

## 4. API Design (Backend)

* Follow RESTful principles where applicable.
* Use clear and consistent endpoint naming.
* Use standard HTTP verbs (GET, POST, PUT, DELETE, PATCH) appropriately.
* Return standard HTTP status codes.
* Structure request and response bodies consistently (JSON).
* Implement validation for all incoming request data.

## 5. Dependencies

* Minimize external dependencies.
* Keep dependencies up-to-date.
* Use Go modules (`go.mod`, `go.sum`) for backend dependency management.
* Use `npm` or `yarn` (`package.json`, `package-lock.json`/`yarn.lock`) for frontend dependency management.

## 6. LLM Collaboration Guidelines (If Applicable)

* **Specify Intent:** Clearly state the goal or task (e.g., "Refactor this function for clarity", "Write unit tests for this handler using TDD", "Generate boilerplate code for a new API endpoint").
* **Provide Context:** Include relevant code snippets, file names, error messages, and existing guidelines (like this document!).
* **Iterate:** Review generated code carefully. Provide feedback for corrections or improvements. Don't blindly accept suggestions.
* **Testing is Key:** Always run tests against code generated or modified with LLM assistance.
* **Prioritize Guidelines:** Remind the LLM of specific guidelines when necessary (e.g., "Remember to follow TDD principles", "Ensure comments explain the 'why'").

---

*This document is a living guide. Please suggest improvements or updates via pull requests or discussion.*
