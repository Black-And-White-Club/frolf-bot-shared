# filepath: /Users/jace/Documents/GitHub/frolf-bot-shared/Makefile

MOCKS_DIR=./mocks
EVENTBUS_FILE=./eventbus/eventbus.go
ERRORS_FILE=./errors/publish_errors.go
ROUND_TYPES=./types/round/types.go
USER_TYPES=./types/user/types.go
PACKAGE=github.com/Black-And-White-Club/frolf-bot-shared

mocks:
	mockgen -source=$(EVENTBUS_FILE) -destination=$(MOCKS_DIR)/eventbus_mock.go -package=mocks
	mockgen -source=$(ERRORS_FILE) -destination=$(MOCKS_DIR)/publish_errors_mock.go -package=mocks
	mockgen -source=$(ROUND_TYPES) -destination=$(MOCKS_DIR)/round_types_mock.go -package=mocks
	mockgen -source=$(USER_TYPES) -destination=$(MOCKS_DIR)/user_types_mock.go -package=mocks
.PHONY: mocks
