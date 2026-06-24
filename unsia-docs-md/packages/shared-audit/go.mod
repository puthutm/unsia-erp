module github.com/unsia-erp/shared-audit

go 1.22

require (
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/unsia-erp/shared-auth v0.0.0
)

replace github.com/unsia-erp/shared-auth => ../shared-auth
