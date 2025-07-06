# LLM Prompt Templates

This document defines prompt templates that can be used with the `--prompt` option of the `list-codes` tool.

## Usage

```bash
./list-codes --prompt explain
./list-codes --prompt find-bugs
./list-codes --prompt refactor
# Other prompt names...
```

## Prompt Template List

### 1. explain (Project Overview)
```
Please analyze the following codebase and explain the following aspects:

1. **Project purpose and main features**
2. **Architecture and design patterns**
3. **Key components and their roles**
4. **Technology stack and dependencies**
5. **Code structure and organization**

Please provide clear and concise explanations that can be understood by non-technical people as well.
```

### 2. find-bugs (Bug Detection)
```
Please analyze the following codebase in detail and identify potential bugs, errors, and issues:

1. **Logic errors and implementation problems**
2. **Memory leaks and performance issues**
3. **Inadequate error handling**
4. **Type safety issues**
5. **Boundary value and edge case handling problems**
6. **Race conditions and thread safety issues**

For each issue, please provide the location, detailed description of the problem, and suggested fixes.
```

### 3. refactor (Refactoring Suggestions)
```
Please analyze the following codebase and identify refactoring opportunities:

1. **Code duplication elimination**
2. **Function and class responsibility separation**
3. **Naming improvements**
4. **Breaking down overly complex functions**
5. **Design pattern application opportunities**
6. **Performance improvements**

For each suggestion, please indicate the current problems, benefits after improvement, and implementation priority.
```

### 4. security (Security Audit)
```
Please analyze the following codebase from a security perspective and identify vulnerabilities and security risks:

1. **Input validation deficiencies**
2. **SQL injection and XSS possibilities**
3. **Authentication and authorization issues**
4. **Exposure of sensitive information**
5. **Inappropriate permission settings**
6. **Encryption and hashing deficiencies**

For each issue, please provide risk level, impact scope, and countermeasures.
```

### 5. optimize (Performance Optimization)
```
Please analyze the following codebase and identify performance optimization opportunities:

1. **Computational complexity improvement possibilities**
2. **Memory usage optimization**
3. **I/O operation efficiency**
4. **Cache utilization opportunities**
5. **Parallel processing introduction possibilities**
6. **Potential bottlenecks**

For each optimization, please indicate current problems, improvement proposals, and expected effects.
```

### 6. test (Testing Suggestions)
```
Please analyze the following codebase and provide testing improvement suggestions:

1. **Areas lacking test coverage**
2. **Insufficient edge case testing**
3. **Need for integration tests**
4. **Test code quality improvements**
5. **Mock and stub utilization opportunities**
6. **Test automation improvements**

Please also provide specific test case examples.
```

### 7. document (Documentation Improvement)
```
Please analyze the following codebase and provide documentation improvement suggestions:

1. **Missing specifications and design documents**
2. **Complex processes lacking comments**
3. **API specification completeness**
4. **README.md improvement points**
5. **Code documentation quality**
6. **Need for usage examples and tutorials**

Please provide suggestions from both user and developer perspectives.
```

### 8. migrate (Technology Migration Suggestions)
```
Please analyze the following codebase and provide suggestions for technology stack migration or updates:

1. **Updates for outdated libraries and frameworks**
2. **Suggestions for more appropriate technology choices**
3. **Language version updates**
4. **Architecture modernization**
5. **Risks and benefits of migration**
6. **Phased migration plan**

Please evaluate the complexity and benefits of migration.
```

### 9. scale (Scalability Analysis)
```
Please analyze the following codebase from a scalability perspective:

1. **Limitations of current architecture**
2. **Potential bottlenecks**
3. **Horizontal and vertical scaling support**
4. **Database design scalability**
5. **Microservices architecture possibilities**
6. **Load balancing mechanisms**

Please provide improvement suggestions for large-scale operations.
```

### 10. maintain (Maintainability Improvement)
```
Please analyze the maintainability of the following codebase and provide improvement suggestions:

1. **Code readability improvements**
2. **Inter-module dependency organization**
3. **Configuration management improvements**
4. **Standardization of logging and error handling**
5. **Development workflow improvements**
6. **Technical debt identification and resolution**

Please provide suggestions from a long-term maintainability perspective.
```

### 11. api-design (API Design Review)
```
Please analyze the API design of the following codebase and provide improvement suggestions:

1. **RESTful design validity**
2. **Endpoint design consistency**
3. **Request and response formats**
4. **Error response standardization**
5. **Versioning strategy**
6. **API documentation completeness**

Please provide suggestions aimed at usable and consistent API design.
```

### 12. patterns (Design Pattern Application)
```
Please analyze the following codebase and suggest applicable design patterns:

1. **Problems with current code structure**
2. **Applicable GoF patterns**
3. **Architectural pattern utilization**
4. **Functional programming patterns**
5. **Concurrency patterns**
6. **Error handling patterns**

Please indicate application points and expected effects for each pattern.
```

### 13. review (Code Review)
```
Please conduct a comprehensive code review of the following codebase:

1. **Compliance with coding conventions**
2. **Application of best practices**
3. **Code quality and consistency**
4. **Potential improvement points**
5. **Issues in team development**
6. **Technologies and techniques to learn**

Please provide constructive and specific feedback.
```

### 14. architecture (Architecture Analysis)
```
Please analyze the architecture of the following codebase and provide evaluation and suggestions:

1. **Characteristics and evaluation of current architecture**
2. **Validity of layered structure**
3. **Direction of dependencies**
4. **Appropriateness of module division**
5. **Application status of design principles (SOLID, etc.)**
6. **Support for future extensibility**

Please provide specific suggestions for better architecture.
```

### 15. deploy (Deployment and Operations Improvement)
```
Please analyze the deployment and operational aspects of the following codebase and provide improvement suggestions:

1. **CI/CD pipeline improvements**
2. **Deployment strategy optimization**
3. **Monitoring and logging enhancements**
4. **Incident response and rollback**
5. **Environment management improvements**
6. **Operations automation opportunities**

Please provide practical improvement suggestions from a DevOps perspective.
```

## Prompt Customization

The above templates are basic forms. It is recommended to customize them according to the characteristics of your project.

### Customization Examples

```bash
# Analysis focused on specific aspects
./list-codes --prompt "Please analyze the usage of goroutines in the following Go code and suggest improvements for concurrent processing:"

# Educational analysis
./list-codes --prompt "Please explain the design patterns and technical elements that beginners can learn from this code:"
```