# Security Policy

## Reporting Security Vulnerabilities

**Please do not report security vulnerabilities through public GitHub issues.**

Instead, please report them via email to: **security@locus-search.org** (update with actual email)

Include:
- Description of the vulnerability
- Steps to reproduce
- Potential impact
- Suggested fix (if any)

We will respond within 48 hours and work with you to address the issue.

## Supported Versions

We support the latest version of each data source implementation with security updates.

## Security Best Practices

When implementing data sources:

### API Keys and Credentials
- ❌ Never hardcode API keys or secrets
- ❌ Never commit credentials to git
- ✅ Accept credentials via configuration/environment variables
- ✅ Document credential management in README

### Input Validation
- ✅ Validate all user inputs
- ✅ Sanitize data before sending to external APIs
- ✅ Check input length to prevent overflow
- ✅ Reject malformed requests early

### HTTP Security
- ✅ Always use HTTPS for API calls
- ✅ Validate TLS certificates
- ✅ Set appropriate timeouts
- ✅ Handle redirects carefully

### Data Handling
- ✅ Don't log sensitive data
- ✅ Sanitize error messages (no credential leaking)
- ✅ Close response bodies (`defer resp.Body.Close()`)
- ✅ Limit response sizes to prevent memory exhaustion

### Dependencies
- ✅ Keep dependencies up to date
- ✅ Review dependency security advisories
- ✅ Use `go mod tidy` regularly
- ✅ Minimize external dependencies

## Known Security Considerations

### API Rate Limiting
Implementations should respect API rate limits to prevent:
- Service disruption
- Account suspension
- IP blocking

### Data Privacy
When handling user queries:
- Be transparent about what data is sent to external APIs
- Respect user privacy
- Follow GDPR/privacy regulations if applicable

## Updates

We will update this policy as needed. Last updated: 2026-02-10
