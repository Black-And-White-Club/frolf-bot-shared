# filepath: /Users/jace/Documents/GitHub/frolf-bot-shared/Makefile

MOCKS_DIR=./mocks
OBSERVABILITY_MOCKS_DIR=observability/mocks
EVENTBUS_MOCKS_DIR=eventbus/mocks

# Define variables for each interface file.
# These paths are RELATIVE to the project root (where the Makefile is).
EVENTBUS_FILE=eventbus/eventbus.go
ERRORS_FILE=errors/publish_errors.go
ROUND_TYPES=types/round/types.go
USER_TYPES=types/user/types.go
EVENTS_METADATA=events/metadata.go

# Observability interfaces
LOKI_FILE=observability/loki.go
PROMETHEUS_FILE=observability/prometheus.go
TEMPO_FILE=observability/tempo.go

mocks: generate-mocks

generate-mocks:
	@echo "Generating mocks..."
	mockgen -source=$(EVENTBUS_FILE) -destination=$(EVENTBUS_MOCKS_DIR)/eventbus_mock.go -package=mocks
	mockgen -source=$(ERRORS_FILE) -destination=$(MOCKS_DIR)/publish_errors_mock.go -package=mocks
	mockgen -source=$(USER_TYPES) -destination=$(MOCKS_DIR)/user_types_mock.go -package=mocks
	mockgen -source=$(LOKI_FILE) -destination=$(OBSERVABILITY_MOCKS_DIR)/loki_mock.go -package=mocks
	mockgen -source=$(PROMETHEUS_FILE) -destination=$(OBSERVABILITY_MOCKS_DIR)/prometheus_mock.go -package=mocks
	mockgen -source=$(TEMPO_FILE) -destination=$(OBSERVABILITY_MOCKS_DIR)/tempo_mock.go -package=mocks
	mockgen -source=$(EVENTS_METADATA) -destination=$(MOCKS_DIR)/metadata_mock.go -package=mocks
	@echo "Mocks generated successfully."

.PHONY: mocks generate-mocks

clean:
	@echo "Cleaning mock files..."
	rm -f $(MOCKS_DIR)/*_mock.go
	@echo "Mocks cleaned successfully."
.PHONY: clean
