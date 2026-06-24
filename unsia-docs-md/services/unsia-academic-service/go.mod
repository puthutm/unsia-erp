module github.com/unsia-erp/unsia-academic-service

go 1.22

require (
	github.com/gin-gonic/gin v1.10.0
	github.com/joho/godotenv v1.5.1
	github.com/unsia-erp/shared-audit v0.0.0
	github.com/unsia-erp/shared-auth v0.0.0
	github.com/unsia-erp/shared-errorenvelope v0.0.0
	github.com/unsia-erp/shared-event v0.0.0
	github.com/unsia-erp/shared-idempotency v0.0.0
	github.com/unsia-erp/shared-observability v0.0.0
	github.com/unsia-erp/shared-rbac v0.0.0
	gorm.io/driver/postgres v1.5.9
	gorm.io/gorm v1.25.10
)

replace (
	github.com/unsia-erp/shared-audit => ../../packages/shared-audit
	github.com/unsia-erp/shared-auth => ../../packages/shared-auth
	github.com/unsia-erp/shared-errorenvelope => ../../packages/shared-errorenvelope
	github.com/unsia-erp/shared-event => ../../packages/shared-event
	github.com/unsia-erp/shared-idempotency => ../../packages/shared-idempotency
	github.com/unsia-erp/shared-observability => ../../packages/shared-observability
	github.com/unsia-erp/shared-rbac => ../../packages/shared-rbac
)
