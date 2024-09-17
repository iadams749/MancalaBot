.PHONY: check-coverage
check-coverage:
	@$(MAKE) LOG MSG_TYPE=info LOG_MESSAGE="Running unit tests and generating coverage report..."
	@go test -shuffle on -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@$(MAKE) LOG MSG_TYPE=success LOG_MESSAGE="Unit tests passed! Coverage report generated at coverage.html"

.PHONY: view-coverage
view-coverage:
	@open coverage.html || xdg-open coverage.html || start coverage.html

# Define colors
WHITE   = \033[0;37m
CYAN    = \033[0;36m
GREEN   = \033[0;32m
RED     = \033[0;31m
NC      = \033[0m   # No Color

# Define the LOG target
LOG:
	@if [ "$(MSG_TYPE)" = "debug" ]; then \
        echo "$(WHITE)$(LOG_MESSAGE)$(NC)"; \
    elif [ "$(MSG_TYPE)" = "success" ]; then \
        echo "$(GREEN)$(LOG_MESSAGE)$(NC)"; \
    elif [ "$(MSG_TYPE)" = "failure" ]; then \
        echo "$(RED)$(LOG_MESSAGE)$(NC)"; \
	elif [ "$(MSG_TYPE)" = "info" ]; then \
        echo "$(CYAN)$(LOG_MESSAGE)$(NC)"; \
    else \
        echo "$(WHITE)$(LOG_MESSAGE)$(NC)"; \
    fi