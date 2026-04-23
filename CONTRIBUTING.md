# Contributing

## Testing

Run the full test suite before opening a PR:

```bash
go test ./...
```

## Test Comment Convention

When writing new tests, add a short intent comment above each test function.
Keep it to 2-3 lines and explain what the test verifies, not how the code works.

Template:

```go
// Verifies <behavior/contract under test>.
// Confirms <expected outcome, edge case, or scope boundary>.
// Notes <important setup constraint> when helpful.
func TestExample(t *testing.T) {
	// test body
}
```

Guidelines:

- Focus on externally visible behavior (inputs, outputs, side effects).
- Mention boundaries such as normalization, error propagation, or precedence rules.
- Avoid restating implementation details that can drift during refactors.
