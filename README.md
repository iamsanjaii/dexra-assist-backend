# Dexra Assist Backend

The backend of the Dexra Assist platform, built with Go and the Gin framework. It is designed to be highly secure and production-ready.

## Security Features

To protect the application and its underlying infrastructure, we've implemented several critical security measures:

- **Authentication Middleware**: Ensures all protected routes require valid authorization.
- **Security Headers**: Standard HTTP headers (e.g., CORS, Content-Security-Policy) to mitigate common web vulnerabilities.
- **Prompt Injection Defense**: Evaluates and safely handles user prompts before passing them to the AI models.
- **Retrieval Grounding**: Ensures LLM responses are exclusively anchored to authorized contextual knowledge.
- **Rate Limiting**: Implemented lightweight API rate limiting (using `golang.org/x/time/rate`) for AI query endpoints and authentication to mitigate abuse and excessive LLM usage.
- **Audit Logging**: Comprehensive activity tracking across all administrative actions and backend services.

## Rate Limiting Specs
The `/chat/query` and `/auth/login` endpoints enforce a strict limit of 5 requests per second, with a burst tolerance of 10. Any excess requests will yield a standard `429 Too Many Requests` response.
