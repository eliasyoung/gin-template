.PHONY: gen-docs
gen-docs:
	@~/go/bin/swag init -g ./cmd/main/main.go --parseInternal --parseDepth 5 && ~/go/bin/swag fmt
