format:
	gci-go write -s standard -s default -s "prefix(VK)" .


mock:
	mockgen -package=mocks  -destination=./internal/rpc/mocks/contract_mocks.go -source=./internal/rpc/contract.go
