# Contributing to Locus DataSource Implementations

Thank you for your interest in contributing a data source implementation!

## Getting Started

### 1. Check Existing Implementations

Before starting, check if someone is already working on a similar data source:
- Review [open pull requests](https://github.com/locus-search/datasource-implementations/pulls)
- Check [issues](https://github.com/locus-search/datasource-implementations/issues) for planned implementations
- Consider opening an issue to announce your intent

### 2. Use the Template

Start with the [datasource-template](https://github.com/locus-search/datasource-template):

```bash
git clone https://github.com/locus-search/datasource-template.git datasource-yourservice
cd datasource-yourservice
# Follow template customization instructions
```

### 3. Implement the Interface

Your data source must fully implement the [DataSource interface](https://github.com/locus-search/datasource-sdk):

```go
type DataSource interface {
    Init() error
    CheckAvailability() bool
    FetchTopics(count int, input NewQuestionInput) ([]DataSourceTopic, error)
    FetchData(count int, topicID int64) ([]DataSourceData, error)
}
```

## Submission Guidelines

### Code Requirements

Your implementation must:

#### Functionality
- ✅ Fully implement all four interface methods
- ✅ Return meaningful data from external API
- ✅ Handle pagination if API supports it
- ✅ Respect API rate limits
- ✅ Support at least basic text search

#### Code Quality
- ✅ Follow Go best practices and idioms
- ✅ Use `gofmt` for formatting
- ✅ Include meaningful comments
- ✅ Avoid code duplication
- ✅ Keep functions focused and testable

#### Error Handling
- ✅ Validate all inputs
- ✅ Return descriptive error messages
- ✅ Handle HTTP errors gracefully
- ✅ Handle API-specific error responses
- ✅ Use context with timeouts

#### Performance
- ✅ Complete requests within 8 seconds (normal case)
- ✅ Use appropriate HTTP client timeouts
- ✅ Implement connection pooling (via `http.Client`)
- ✅ Avoid memory leaks (`defer resp.Body.Close()`)

#### Testing
- ✅ Include unit tests for validation logic
- ✅ Provide integration test examples (can use mocks)
- ✅ All tests must pass: `go test ./...`
- ✅ Aim for >70% code coverage

### Documentation Requirements

#### README.md

Your data source directory must include a comprehensive README:

```markdown
# DataSource [Service Name]

Brief description of the service and what data it provides.

## Features

- List key features
- Note any special capabilities
- Mention limitations

## Setup

### Prerequisites
- List requirements (API keys, accounts, etc.)
- Link to registration pages

### Installation
\`\`\`bash
go get github.com/locus-search/datasource-implementations/yourservice
\`\`\`

### Configuration
- How to obtain API credentials
- Environment variables
- Configuration options

## Usage

### Basic Example
\`\`\`go
// Working code example
\`\`\`

### With Locus
\`\`\`go
// Integration example
\`\`\`

## API Details

- Rate limits
- Quotas
- Terms of service considerations
- Regional availability

## Testing

How to run tests, including any setup needed.

## License

State the license (MIT preferred).
```

#### Code Comments

- Document all exported types and functions (godoc style)
- Explain non-obvious logic
- Note any API quirks or workarounds

### File Structure

Your contribution should follow this structure:

```
yourservice/
├── datasource.go           # Main implementation
├── datasource_test.go      # Tests
├── README.md              # Documentation
├── LICENSE                # MIT, Apache 2.0, or BSD
└── [helper files]         # Any additional files
```

## Pull Request Process

### 1. Prepare Your Code

```bash
# Fork the repository
git clone https://github.com/yourusername/datasource-implementations.git
cd datasource-implementations

# Create your directory
mkdir yourservice
cd yourservice

# Add your implementation files
# ... develop and test ...

# Verify everything works
go test ./...
go mod tidy
```

### 2. Update Main README

Add your data source to the table in the main [README.md](README.md):

```markdown
| **Your Service** | Brief description | ✅ Stable | [README](yourservice/README.md) |
```

Or mark as:
- `🚧 Beta` - Still being refined
- `🧪 Experimental` - Early stage

### 3. Submit Pull Request

1. **Commit your changes**:
   ```bash
   git add yourservice/
   git commit -m "Add [ServiceName] data source implementation"
   ```

2. **Push to your fork**:
   ```bash
   git push origin main
   ```

3. **Open a Pull Request** with:
   - **Title**: "Add [ServiceName] data source"
   - **Description**:
     - What service you're integrating
     - Key features
     - Any special considerations
     - Link to API documentation
     - Testing performed

### 4. Review Process

Maintainers will review your PR for:

1. **Interface compliance** - Does it fully implement `DataSource`?
2. **Code quality** - Is it readable and well-structured?
3. **Testing** - Are there adequate tests?
4. **Documentation** - Is setup/usage clear?
5. **API best practices** - Does it respect API terms?
6. **Security** - No hardcoded credentials, proper input validation?

**Response time**: We aim to provide initial feedback within 1 week.

### 5. Address Feedback

- Respond to review comments
- Make requested changes
- Push updates to your PR branch
- Re-request review when ready

## Best Practices

### API Authentication

Never hardcode credentials:

```go
// ❌ Bad
const APIKey = "sk_live_1234567890"

// ✅ Good
type DataSourceYourService struct {
    APIKey string // Set by user
}
```

### Rate Limiting

Respect API limits:

```go
import "golang.org/x/time/rate"

type DataSourceYourService struct {
    rateLimiter *rate.Limiter
}

func New() *DataSourceYourService {
    return &DataSourceYourService{
        rateLimiter: rate.NewLimiter(rate.Limit(10), 20),
    }
}
```

### Timeouts

Always use timeouts:

```go
ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
defer cancel()

req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
```

### Error Messages

Be descriptive:

```go
// ❌ Bad
return nil, errors.New("error")

// ✅ Good
return nil, fmt.Errorf("failed to fetch topics from %s: %w", endpoint, err)
```

### Testing

Test edge cases:

```go
func TestFetchTopics(t *testing.T) {
    tests := []struct {
        name      string
        input     datasource.NewQuestionInput
        wantError bool
    }{
        {"valid input", datasource.NewQuestionInput{QuestionText: "test"}, false},
        {"empty input", datasource.NewQuestionInput{}, true},
        {"very long input", datasource.NewQuestionInput{QuestionText: strings.Repeat("a", 10000)}, true},
    }
    // ...
}
```

## Questions?

- 💬 [Open a discussion](https://github.com/locus-search/locus/discussions)
- 📧 Reach out to maintainers
- 📖 Review existing implementations for examples

## Code of Conduct

Be respectful, inclusive, and collaborative. We're all here to build something great together!

---

Thank you for contributing to Locus! 🎉
