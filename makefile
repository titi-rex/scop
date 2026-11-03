build_dir	= build

target			= scop
junit_report	= $(build_dir)/report.xml
cover_report	= $(build_dir)/coverage.out


build:
	@echo $P building $(notdir $(target))
	@go build -o $(target) ./internal/

test:
	@echo $P running test
	@go tool gotestsum \
		--format=dots \
		--junitfile $(junit_report) \
		-- \
		-coverprofile $(cover_report) \
		 ./...


run: build
	@echo $P run $(notdir $(target))
	@./$(target) ./resources/42.obj

tidy:
	@echo $P Tidying up!
	@go mod tidy

vendor:
	@echo $P Vendoring things
	@go mod vendor

# @gum spin --spinner meter --title "Vedoring" -- go mod vendor

clean:
	@echo $P cleanup
	@rm -rf $(build_dir)

re: clean build

e-%:
	@echo $*


P	= "\033[0;35m▶\033[0m"
.PHONY: build run test tidy vendor
