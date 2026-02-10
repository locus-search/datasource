# Locus DataSource Implementations

Official collection of data source implementations for [Locus](https://github.com/locus-search/locus).

## Overview

This repository hosts production-ready implementations of the [DataSource interface](https://github.com/locus-search/datasource-sdk) that integrate external knowledge sources with Locus.

## Available Data Sources

### Official Implementations

| Data Source | Description | Status | Documentation |
|------------|-------------|--------|---------------|
| **Stack Exchange** | Search across 170+ Stack Exchange sites (Stack Overflow, Server Fault, etc.) | ✅ Stable | [README](stackexchange/README.md) |
| **Wikipedia** | Search Wikipedia articles and extracts | ✅ Stable | [README](wikipedia/README.md) |
| **DuckDuckGo** | Instant answers and web search results | ✅ Stable | [README](duckduckgo/README.md) |

### Community Contributions

We welcome community-contributed data sources! See [Contributing](#contributing) below.

## Quick Start

### Using a Data Source

```bash
# Install the data source you want
go get github.com/locus-search/datasource-implementations/stackexchange
```

```go
package main

import (
    "github.com/locus-search/locus/backend/internal/core"
    "github.com/locus-search/datasource-implementations/stackexchange"
)

func main() {
    // Create Locus core
    locusCore := core.New()
    
    // Create and register Stack Exchange data source
    se := stackexchange.New()
    se.Key = "your-stackapps-api-key" // Optional but recommended
    locusCore.RegisterDataSource(se)
    
    // Initialize
    if err := locusCore.InitDataSources(); err != nil {
        panic(err)
    }
    
    // Now Stack Exchange is available for queries
}
```

### Using Multiple Data Sources

```go
import (
    "github.com/locus-search/datasource-implementations/stackexchange"
    "github.com/locus-search/datasource-implementations/wikipedia"
    "github.com/locus-search/datasource-implementations/duckduckgo"
)

func setupDataSources(core *core.Core) {
    // Register multiple sources
    core.RegisterDataSource(stackexchange.New())
    core.RegisterDataSource(wikipedia.New())
    core.RegisterDataSource(duckduckgo.New())
    
    // Initialize all at once
    core.InitDataSources()
}
```

## Contributing

We welcome high-quality data source implementations! 

### Before Contributing

1. **Use the template**: Start with [datasource-template](https://github.com/locus-search/datasource-template)
2. **Follow the SDK**: Implement the [datasource-sdk](https://github.com/locus-search/datasource-sdk) interface
3. **Test thoroughly**: Include comprehensive tests
4. **Document well**: Provide clear setup and usage instructions

### Contribution Process

1. **Create your implementation** using the template
2. **Test it** with Locus to ensure it works correctly
3. **Fork this repository**
4. **Add your implementation** in a new directory: `yourservice/`
5. **Include**:
   - Complete source code
   - README.md with setup instructions
   - Tests
   - LICENSE (MIT preferred, but Apache 2.0 or BSD acceptable)
6. **Update** this main README to list your data source
7. **Submit a Pull Request**

### Quality Standards

Your implementation must:

- ✅ Fully implement the `DataSource` interface
- ✅ Handle errors gracefully
- ✅ Include timeouts (≤ 8 seconds for normal operations)
- ✅ Validate inputs
- ✅ Have passing tests (use `go test ./...`)
- ✅ Include documentation
- ✅ Respect API rate limits
- ✅ Use meaningful error messages
- ✅ Follow Go best practices

### Review Criteria

Pull requests are reviewed for:

1. **Code Quality**: Readable, idiomatic Go code
2. **Completeness**: Full implementation of interface
3. **Documentation**: Clear README with setup/usage
4. **Testing**: Good test coverage
5. **API Compliance**: Proper use of external APIs
6. **Error Handling**: Robust error handling
7. **Performance**: Reasonable response times

## Repository Structure

```
datasource-implementations/
├── README.md                    # This file
├── stackexchange/              # Stack Exchange implementation
│   ├── datasource.go
│   ├── init.go
│   ├── questions.go
│   ├── datasource_test.go
│   ├── README.md
│   └── LICENSE
├── wikipedia/                   # Wikipedia implementation
│   ├── datasource.go
│   ├── datasource_test.go
│   ├── README.md
│   └── LICENSE
├── duckduckgo/                  # DuckDuckGo implementation
│   ├── datasource.go
│   ├── datasource_test.go
│   ├── README.md
│   └── LICENSE
└── yourservice/                 # Your contribution
    ├── datasource.go
    ├── datasource_test.go
    ├── README.md
    └── LICENSE
```

## Development

### Running Tests

Test all implementations:
```bash
go test ./...
```

Test specific implementation:
```bash
go test ./stackexchange/...
```

With coverage:
```bash
go test -cover ./...
```

### Building

All implementations should be importable as Go modules. No build step required.

## FAQ

### How do I choose which data sources to use?

It depends on your use case:
- **Technical Q&A**: Stack Exchange
- **General knowledge**: Wikipedia
- **Quick facts**: DuckDuckGo
- **Multiple domains**: Use all three!

### Do I need API keys?

- **Stack Exchange**: Optional but recommended (increases rate limits)
- **Wikipedia**: No
- **DuckDuckGo**: No

### Can I contribute a data source for a commercial API?

Yes, as long as:
1. Users can obtain their own API keys
2. Your implementation doesn't violate the API's terms of service
3. You clearly document any costs or limitations

### What if my data source needs special setup?

Document it clearly in your README. Examples:
- OAuth2 flows
- API key registration
- Rate limit considerations
- Regional availability

### Can I update an existing implementation?

Yes! Submit a PR with improvements. We welcome:
- Performance optimizations
- Bug fixes
- Better error handling
- Additional features
- Documentation improvements

## Resources

- 📖 [DataSource SDK](https://github.com/locus-search/datasource-sdk) - Interface documentation
- 🎨 [DataSource Template](https://github.com/locus-search/datasource-template) - Starter template
- 💬 [Discussions](https://github.com/locus-search/locus/discussions) - Ask questions
- 🐛 [Issues](https://github.com/locus-search/datasource-implementations/issues) - Report bugs

## License

Each data source implementation has its own license (see individual directories). New contributions should use MIT, Apache 2.0, or BSD licenses for maximum compatibility.

The repository structure and documentation are released under the MIT License.

## Contact

Questions? Open an issue or discussion in the [main Locus repository](https://github.com/locus-search/locus).

---

**Happy integrating!** 🚀
