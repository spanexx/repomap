# Milestone 4.2: Quick Reference

## Summary

Authentication and authorization framework for API security. API keys, JWT, OAuth 2.0, RBAC, audit logging.

## Timeline

- **Duration:** 2–3 weeks
- **Tasks:** 25 sequential

## Key Metrics

| Metric | Target |
|--------|--------|
| Auth Latency | <100ms |
| Code Coverage | >80% |
| Security Review | Passed |

## Acceptance Criteria

- ✅ Multiple auth methods working
- ✅ Fine-grained authorization
- ✅ Audit trail functional
- ✅ <100ms auth latency
- ✅ Security review approved
- ✅ >80% code coverage

## Example Usage

```bash
# Create API key
agents-server admin add-key --user alice --role admin

# Use API key
curl -H "Authorization: ApiKey <key>" http://localhost:8080/api/v1/repomap

# JWT token auth
curl -H "Authorization: Bearer <jwt>" http://localhost:8080/api/v1/sandbox-run
```
