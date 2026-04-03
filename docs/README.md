# Documentation

This directory contains detailed documentation for the Go Foundation v2 template. Each guide focuses on a specific aspect of the template to help you understand and customize it for your service.

## 📚 Documentation Index

### [Example Service](example-service.md)
Understanding the included User Service that demonstrates the template's architecture and patterns.

**Topics covered:**
- Architecture overview and directory structure
- Key components (entry point, server, service, repository, model layers)
- Request flow through the system
- Code patterns (dependency injection, interfaces, error handling)
- Running the example service

**Start here if:** you want to understand how the template works before making changes.

---

### [Migration Guide](migration-guide.md)
Step-by-step instructions for transitioning from the example User Service to your custom service.

**Topics covered:**
- Three-phase migration approach
- Copying service structure and updating proto definitions
- Implementing domain models and business logic
- Configuration and deployment updates
- What to keep vs. what to replace

**Start here once:** you're ready to build your own service

---

### [Best Practices](best-practices.md)
Coding standards and architectural guidelines to maintain when building services.

**Topics covered:**
- Layered architecture principles
- Error handling patterns
- Configuration management
- Database patterns and migrations
- Testing guidelines
- Code organization and API design
- Logging, observability, and security

**Start here if:** you want to ensure your service follows established patterns and standards.

---

### [Build System](build-system.md)
Understanding the Makefile structure and Docker build process.

**Topics covered:**
- Makefile structure and organization
- Key variables to configure
- Essential Make targets (development, Docker, proto, utility)
- Environment-specific builds
- Extending the Makefile with custom targets
- Troubleshooting common build issues

**Start here if:** you need to understand or customize the build process.

---

### [Proto Management](proto-management.md)
Guide to managing Protocol Buffers with the centralized proto repository approach.

**Topics covered:**
- Proto modules configuration
- Complete proto workflow (fetch, generate, lint, refresh)
- Proto file structure
- Buf configuration and customization
- Central proto repository integration
- Proto best practices and troubleshooting

**Start here if:** you need to work with Protocol Buffers or integrate with the central proto repository.

---

## Quick Reference

### Getting Started Flow
1. Read [Example Service](example-service.md) to understand the architecture
2. Follow [Migration Guide](migration-guide.md) to create your service
3. Reference [Best Practices](best-practices.md) as you implement features
4. Use [Build System](build-system.md) and [Proto Management](proto-management.md) as needed

### Common Tasks
- **Running the example:** → [Example Service - Running the Example](example-service.md#running-the-example)
- **Creating your service:** → [Migration Guide - Phase 2](migration-guide.md#phase-2-gradual-replacement)
- **Updating proto files:** → [Proto Management - Proto Workflow](proto-management.md#proto-workflow)
- **Building and testing:** → [Build System - Essential Make Targets](build-system.md#essential-make-targets)
- **Following patterns:** → [Best Practices](best-practices.md)

## Additional Resources

- **Main README:** [../README.md](../README.md) - Quick start guide and setup instructions
- **Razorpay Go Style Guide:** https://github.com/razorpay/go-style-guide
- **Foundation Library:** https://github.com/razorpay/foundation
- **Proto Repository:** https://github.com/razorpay/proto

---

**Need help?** Reach out to the Developer Experience Engineering team on slack [#developer-experience](https://razorpay.enterprise.slack.com/archives/C08DS8AE7T8)

