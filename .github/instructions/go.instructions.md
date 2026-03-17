---
name: 'Go Standards'
description: 'Coding conventions for Go files'
applyTo: '**/*.go'
---
Go Version: 1.25.6

1. Use camelCase for variable and function names, and PascalCase for exported names.
2. Keep comments in code under 80 characters for better readability.
3. Use tabs for indentation, not spaces.
4. Always include error handling for functions that return an error.
5. Use descriptive names for variables and functions to improve code clarity.
6. Order files as follows: package declaration, imports, constants, variables, types, exported functions, unexported functions.
7. Order functions alphabetically within their respective sections (exported vs unexported).

## Testing Conventions

1. Tests should be placed in a separate file with the suffix `_test.go`.
2. Tests should be organized using `t.Run` to group related test cases together.
3. Use the `assert` package for assertions in tests to improve readability and maintainability.
4. Test functions naming convention should be `Test_<FunctionName>`, `Test_<StructName>_<FunctionName>` for clarity.
5. Put a comment for each test function specifying what function it is testing, e.g. `// Tests for [CreateZipFileWithRandomFiles] function.`
6. Test names in `t.Run` should be descriptive and all lowercase, e.g. `t.Run("create 3 temp files", func(t *testing.T) { ... })`.
