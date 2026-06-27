module github.com/BlackboxAI/unsia-docs-md/packages/shared-serviceclient

go 1.22

require (
	github.com/BlackboxAI/unsia-docs-md/packages/shared-errorenvelope v0.0.0
	github.com/BlackboxAI/unsia-docs-md/packages/shared-httpclient v0.0.0
)

replace (
	github.com/BlackboxAI/unsia-docs-md/packages/shared-errorenvelope => ../shared-errorenvelope
	github.com/BlackboxAI/unsia-docs-md/packages/shared-httpclient => ../shared-httpclient
)
