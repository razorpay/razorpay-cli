# Protocol Buffers Management

Guide to managing Protocol Buffers with the centralized proto repository approach.

## Overview

This template uses a centralized proto repository to manage proto definitions across services. Proto files are fetched from a central repository and code is generated locally.

## Proto Modules Configuration

The `build/proto_modules` file specifies which proto modules to fetch:

```bash
# View current proto module
cat build/proto_modules
# Output: go_foundation_v2
```

### Updating for Your Service

After renaming your service, update the proto module:

```bash
echo "your_service_name" > build/proto_modules
```

## Proto Workflow

### 1. Fetch Proto Files

Download proto files from the central repository:

```bash
make proto-fetch
```

This fetches proto files from `https://github.com/razorpay/proto.git` into the `proto/` directory.

### 2. Generate Go Code

Generate Go code from proto files:

```bash
make proto-generate
```

This creates Go code in the `rpc/` directory using buf and protoc plugins.

### 3. Lint Proto Files

Validate proto files against best practices:

```bash
make proto-lint
```

### 4. Complete Refresh

Run the complete proto workflow:

```bash
make proto-refresh
```

This performs: clean → fetch → generate → lint

## Proto File Structure

After fetching and generating, your proto structure will look like:

```
proto/
└── your_service_name/
    └── your_domain/
        └── v1/
            └── your_service.proto

rpc/
└── your_service_name/
    └── your_domain/
        └── v1/
            ├── your_service.pb.go          # Generated proto messages
            ├── your_service_grpc.pb.go     # Generated gRPC services
            └── your_service.pb.validate.go # Generated validation
```

## Buf Configuration

The template uses [buf](https://buf.build/) for proto management:

- `buf.yaml`: Defines the proto module and lint rules
- `buf.gen.yaml`: Configures code generation plugins
- `buf.lock`: Locks dependencies for reproducible builds

### Updating Buf Configuration

Modify `buf.gen.yaml` to customize code generation:

```yaml
version: v1
plugins:
  - plugin: buf.build/protocolbuffers/go
    out: rpc
    opt: paths=source_relative
  - plugin: buf.build/grpc/go
    out: rpc
    opt: paths=source_relative
```

## Central Proto Repository

### Adding New Proto Definitions

1. Add your proto files to the central proto repository
2. Follow the standard directory structure: `your_service/domain/v1/`
3. Update `build/proto_modules` in your service repository
4. Run `make proto-refresh` to fetch and generate code

### Proto Best Practices

- Version your APIs (`v1`, `v2`, etc.)
- Use meaningful message and field names
- Add validation rules using `buf.build/bufbuild/protovalidate`
- Document your proto files with comments
- Keep proto files focused on a single domain

## Troubleshooting

### Missing Generated Code
```bash
# Regenerate proto code
make proto-generate
```

### Proto Lint Errors
```bash
# Check lint errors
make proto-lint

# Fix formatting issues in proto files
```

### Buf Dependencies Out of Date
```bash
# Update buf dependencies
buf mod update proto/
```

### Proto Fetch Failures
Ensure you have access to the central proto repository and your GitHub credentials are configured correctly.

