# Panic Recovery - Production Safety

## 🛡️ Overview

Service đã được trang bị **panic recovery** ở nhiều tầng để đảm bảo service không bị crash trong production.

## 📍 Recovery Layers

### 1. **Global Recovery** (main.go)

```go
func main() {
    defer func() {
        if r := recover() {
            // Log và exit gracefully
        }
    }()
    // Application logic...
}
```

**Catches:** Panics ở top level của application

### 2. **gRPC Interceptors** (middleware)

```go
// Unary RPC Recovery
grpc.UnaryInterceptor(
    middleware.RecoveryUnaryInterceptor(logger)
)

// Stream RPC Recovery
grpc.StreamInterceptor(
    middleware.RecoveryStreamInterceptor(logger)
)
```

**Catches:** Panics trong gRPC handlers, trả về error thay vì crash

### 3. **Goroutine Recovery** (middleware)

```go
middleware.RecoverGoroutine(logger, "name", func() {
    // Goroutine logic
})
```

**Catches:** Panics trong goroutines

### 4. **Function Recovery** (middleware)

```go
err := middleware.RecoverFunc(logger, "FuncName", func() error {
    // Function logic
})
```

**Catches:** Panics trong critical functions

## 🎯 What Happens When Panic Occurs?

### gRPC Handler Panic

```
Client Request → Handler Panics
    ↓
Recovery Interceptor Catches
    ↓
Log: Error + Stack Trace
    ↓
Return: Internal Server Error to Client
    ↓
Service Continues Running ✅
```

### Goroutine Panic

```
Goroutine Starts → Code Panics
    ↓
Recovery Wrapper Catches
    ↓
Log: Error + Stack Trace
    ↓
Goroutine Terminates
    ↓
Service Continues Running ✅
```

### Main Function Panic

```
Main Execution → Panic
    ↓
Global Defer Catches
    ↓
Log: Fatal Error + Stack Trace
    ↓
Graceful Exit with Code 1 ⚠️
```

## 📊 Example Logs

### gRPC Handler Panic

```json
{
  "level": "error",
  "timestamp": "2025-10-25T10:30:45Z",
  "msg": "Panic recovered in gRPC handler",
  "method": "/iam.IAMService/Login",
  "panic": "runtime error: index out of range",
  "stack": "goroutine 123 [running]:\n..."
}
```

### Goroutine Panic

```json
{
  "level": "error",
  "timestamp": "2025-10-25T10:31:12Z",
  "msg": "Panic recovered in goroutine",
  "goroutine": "grpc-server",
  "panic": "nil pointer dereference",
  "stack": "goroutine 456 [running]:\n..."
}
```

## 🔧 Usage

### Wrap Goroutines

```go
// ❌ Before (unsafe)
go func() {
    // Code that might panic
}()

// ✅ After (safe)
middleware.RecoverGoroutine(logger, "worker-name", func() {
    // Code that might panic
})
```

### Wrap Critical Functions

```go
// ✅ Wrap initialization
err := middleware.RecoverFunc(logger, "Initialize", func() error {
    // Critical initialization code
    return nil
})
```

### gRPC Handlers

```go
// ✅ Automatically protected by interceptors
grpcServer := grpc.NewServer(
    grpc.ChainUnaryInterceptor(
        middleware.RecoveryUnaryInterceptor(logger),
    ),
)
```

## 🚨 When Recovery Happens

| Scenario | Recovery Layer | Service Status | Client Gets |
|----------|----------------|----------------|-------------|
| Handler panic | gRPC Interceptor | ✅ Running | Internal Error |
| Goroutine panic | Goroutine Wrapper | ✅ Running | N/A |
| Init panic | Function Wrapper | ✅ Running | Error returned |
| Main panic | Global Defer | ❌ Exit | N/A |

## 💡 Best Practices

### 1. **Always Wrap Goroutines**

```go
// Production code
middleware.RecoverGoroutine(logger, "background-task", func() {
    // Long-running task
})
```

### 2. **Log Panic Details**

All recovery includes:
- ✅ Panic value
- ✅ Full stack trace
- ✅ Context (method, goroutine name, etc.)

### 3. **Return Errors to Clients**

gRPC clients receive:
```
status: INTERNAL
message: "Internal server error: <panic message>"
```

### 4. **Monitor Recovery Logs**

Set up alerts for panic logs:
```bash
# Example: Alert on panic recovery
grep "Panic recovered" /var/log/iam-service.log
```

## 🧪 Testing Recovery

### Simulate Panic in Handler

```go
func (h *Handler) TestPanic(ctx context.Context, req *pb.Request) (*pb.Response, error) {
    panic("test panic")  // Service will NOT crash ✅
}
```

### Test Result

```bash
# Client sees error
status: INTERNAL
message: "Internal server error: test panic"

# Server logs panic but continues
{"level":"error","msg":"Panic recovered in gRPC handler"...}

# Service still healthy
curl localhost:8080/health
{"status":"ok"}
```

## 📈 Benefits

| Benefit | Impact |
|---------|--------|
| **High Availability** | Service doesn't crash on panic |
| **Better UX** | Clients get errors instead of connection loss |
| **Debugging** | Full stack traces in logs |
| **Production Ready** | Safe for production use |
| **Peace of Mind** | Sleep well at night 😴 |

## 🔍 Monitoring

### Key Metrics to Track

1. **Panic Rate**
   ```
   rate(panic_recovered_total[5m])
   ```

2. **Panic by Method**
   ```
   panic_recovered_total{method="/iam.IAMService/Login"}
   ```

3. **Recovery Success Rate**
   ```
   recovery_success_total / recovery_attempts_total
   ```

## 🎯 Production Checklist

- [x] Global panic recovery in main
- [x] gRPC unary interceptor
- [x] gRPC stream interceptor  
- [x] Goroutine wrappers
- [x] Function wrappers for critical code
- [x] Detailed logging with stack traces
- [x] Error returned to clients
- [x] Service continues running

## 📚 Files

- `internal/middleware/recovery.go` - Recovery middleware
- `internal/app/app.go` - Uses recovery in Initialize/Run
- `cmd/server/main.go` - Global recovery

---

**Status:** ✅ **PRODUCTION READY** - Service protected against panics!

