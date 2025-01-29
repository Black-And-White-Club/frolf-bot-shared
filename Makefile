# filepath: /Users/jace/Documents/GitHub/frolf-bot-shared/Makefile

MOCKS_DIR=./mocks
EVENTBUS_FILE=./eventbus/eventbus.go
PACKAGE=github.com/Black-And-White-Club/frolf-bot-shared

mocks:
	mockgen -source=$(EVENTBUS_FILE) -destination=$(MOCKS_DIR)/eventbus_mock.go -package=mocks
.PHONY: mocks
