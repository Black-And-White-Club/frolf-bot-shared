# filepath: /Users/jace/Documents/GitHub/frolf-bot-shared/Makefile

MOCKS_DIR=./mocks
EVENTBUS_FILE=./eventbus/eventbus.go
ERRORS_FILE=./errors/publish_errors.go
PACKAGE=github.com/Black-And-White-Club/frolf-bot-shared

mocks:
	mockgen -source=$(EVENTBUS_FILE) -destination=$(MOCKS_DIR)/eventbus_mock.go -package=mocks
	mockgen -source=$(ERRORS_FILE) -destination=$(MOCKS_DIR)/publish_errors_mock.go -package=mocks
.PHONY: mocks
