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
UTILS_METADATA=utils/metadata.go
MESSAGES_UTILS=utils/messages.go
# Observability interfaces
TEMPO_FILE=observability/otel/tracing/tempo.go
USER_METRICS=observability/otel/metrics/user/interface.go
SCORE_METRICS=observability/otel/metrics/score/interface.go
LEADERBOARD_METRICS=observability/otel/metrics/leaderboard/interface.go
ROUND_METRICS=observability/otel/metrics/round/interface.go
DISCORD_METRICS=observability/otel/metrics/discord/interface.go

mocks: generate-mocks

generate-mocks:
	@echo "Generating mocks..."
	mockgen -source=$(EVENTBUS_FILE) -destination=$(EVENTBUS_MOCKS_DIR)/eventbus_mock.go -package=mocks
	mockgen -source=$(ERRORS_FILE) -destination=$(MOCKS_DIR)/publish_errors_mock.go -package=mocks
	mockgen -source=$(USER_TYPES) -destination=$(MOCKS_DIR)/user_types_mock.go -package=mocks
	mockgen -source=$(USER_METRICS) -destination=$(OBSERVABILITY_MOCKS_DIR)/usermetrics_mock.go -package=mocks
	mockgen -source=$(SCORE_METRICS) -destination=$(OBSERVABILITY_MOCKS_DIR)/scoremetrics_mock.go -package=mocks
	mockgen -source=$(ROUND_METRICS) -destination=$(OBSERVABILITY_MOCKS_DIR)/roundmetrics_mock.go -package=mocks
	mockgen -source=$(DISCORD_METRICS) -destination=$(OBSERVABILITY_MOCKS_DIR)/discordmetrics_mock.go -package=mocks
	mockgen -source=$(LEADERBOARD_METRICS) -destination=$(OBSERVABILITY_MOCKS_DIR)/leaderboard_mock.go -package=mocks
	mockgen -source=$(TEMPO_FILE) -destination=$(OBSERVABILITY_MOCKS_DIR)/tracer_mock.go -package=mocks
	mockgen -source=$(UTILS_METADATA) -destination=$(MOCKS_DIR)/metadata_mock.go -package=mocks
	mockgen -source=$(MESSAGES_UTILS) -destination=$(MOCKS_DIR)/messages_mock.go -package=mocks
	@echo "Mocks generated successfully."

.PHONY: mocks generate-mocks

clean:
	@echo "Cleaning mock files..."
	rm -f $(MOCKS_DIR)/*_mock.go
	@echo "Mocks cleaned successfully."
.PHONY: clean


