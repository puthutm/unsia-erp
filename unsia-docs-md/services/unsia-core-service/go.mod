module github.com/unsia-erp/unsia-core-service

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
	github.com/golang-jwt/jwt/v5 v5.2.1
	golang.org/x/crypto v0.24.0
	gorm.io/driver/postgres v1.5.9
	gorm.io/gorm v1.25.10
	github.com/google/uuid v1.6.0
	github.com/go-redis/redis/v8 v8.11.5
	github.com/hibiken/asynq v0.24.1
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
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
