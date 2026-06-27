# shared-serviceclient

This package provides a factory pattern for inter-service HTTP communication in the UNSIA microservices architecture.

## Features

- **Service Client Factory**: Centralized client management with caching
- **Circuit Breaker**: Built-in resilience using shared-httpclient
- **typed Responses**: Generic response handling with typed results
- **Environment Override**: Service URLs configurable via environment
- **Service-to-Service Auth**: Automatic token handling

## Usage

### Basic Client Creation

```go
import "github.com/BlackboxAI/unsia-docs-md/packages/shared-serviceclient"

// Create a factory with default config
factory := sharedserviceclient.NewFactory(sharedserviceclient.Config{
    Timeout:        30 * time.Second,
    MaxRetries:     3,
    FailureThreshold: 5,
    Cooldown:      30 * time.Second,
})

// Get Academic service client
academicClient, err := factory.GetClient(sharedserviceclient.ServiceAcademic)
if err != nil {
    log.Fatal(err)
}

// Make GET request
type StudentResponse struct {
    ID   string `json:"id"`
    Nim  string `json:"nim"`
    Name string `json:"name"`
}

var student StudentResponse
err = academicClient.GetTyped(ctx, "/api/v1/students/some-id", &student)
```

### Direct Client

```go
client := sharedserviceclient.New(sharedserviceclient.Config{
    ServiceName: sharedserviceclient.ServiceFinance,
    Timeout:     30 * time.Second,
})

var invoice Invoice
err = client.GetTyped(ctx, "/api/v1/invoices/123", &invoice)
```

### Service Names

```go
const (
    ServiceCore        // Core (Auth/SSO) - Port 8001
    ServiceReference  // Reference - Port 8002
    ServicePMB       // PMB - Port 8003
    ServiceAcademic  // Academic - Port 8004
    ServiceFinance   // Finance - Port 8005
    ServiceLMS      // LMS - Port 8006
    ServiceHRIS     // HRIS - Port 8008
    ServiceAssessment // Assessment - Port 8007
    ServiceCRM      // CRM - Port 8009
    ServicePortal   // Portal - Port 8010
)
```

### Environment Variables

Each service URL can be overridden:

```bash
ACADEMIC_SERVICE_URL=http://localhost:8004
FINANCE_SERVICE_URL=http://localhost:8005
SERVICE_TOKEN=your-service-token
```

## Cross-Service Communication Example

### Finance Service checking Academic Clearance

```go
// In Finance service handler
factory := sharedserviceclient.NewFactory(sharedserviceclient.Config{
    ServiceToken: os.Getenv("SERVICE_TOKEN"),
    Timeout:     10 * time.Second,
})

academicClient, _ := factory.GetClient(sharedserviceclient.ServiceAcademic)

// Check if student has no unpaid KRS
type ClearanceCheckRequest struct {
    StudentID string `json:"student_id"`
    Type      string `json:"type"` // "krs", "tuition"
}

type ClearanceCheckResponse struct {
    IsClear bool   `json:"is_clear"`
    Reason string `json:"reason,omitempty"`
}

var clearance ClearanceCheckResponse
err = academicClient.PostTyped(ctx, "/api/v1/clearance/check", 
    ClearanceCheckRequest{StudentID: studentID, Type: "krs"},
    &clearance)
if err != nil {
    return err
}

if !clearance.IsClear {
    return fmt.Errorf("student not clear: %s", clearance.Reason)
}
```

## Dependencies

- `shared-httpclient`: HTTP client with circuit breaker
- `shared-errorenvelope`: Error response format
