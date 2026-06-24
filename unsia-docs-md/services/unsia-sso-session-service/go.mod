module github.com/unsia-erp/unsia-sso-session-service

go 1.22

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/joho/godotenv v1.5.1
	github.com/unsia-erp/shared-errorenvelope v0.0.0
	github.com/unsia-erp/shared-observability v0.0.0
	gorm.io/driver/postgres v1.5.4
	gorm.io/gorm v1.25.5
)

replace (
	github.com/unsia-erp/shared-errorenvelope => ../packages/shared-errorenvelope
	github.com/unsia-erp/shared-observability => ../packages/shared-observability
)
