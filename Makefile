.PHONY: check-coverage
check-coverage:
	go test -shuffle on -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated and opening in browser..."
	@open coverage.html || xdg-open coverage.html || start coverage.html