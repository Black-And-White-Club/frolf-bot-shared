package observability

import (
	"log/slog"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"

	discordmetrics "github.com/Black-And-White-Club/frolf-bot-shared/observability/otel/metrics/discord"
	eventbusmetrics "github.com/Black-And-White-Club/frolf-bot-shared/observability/otel/metrics/eventbus"
	guildmetrics "github.com/Black-And-White-Club/frolf-bot-shared/observability/otel/metrics/guild"
	leaderboardmetrics "github.com/Black-And-White-Club/frolf-bot-shared/observability/otel/metrics/leaderboard"
	roundmetrics "github.com/Black-And-White-Club/frolf-bot-shared/observability/otel/metrics/round"
	scoremetrics "github.com/Black-And-White-Club/frolf-bot-shared/observability/otel/metrics/score"
	usermetrics "github.com/Black-And-White-Club/frolf-bot-shared/observability/otel/metrics/user"
)

type Registry struct {
	Logger *slog.Logger
	Tracer trace.Tracer
	Meter  metric.Meter

	// Metrics
	UserMetrics        usermetrics.UserMetrics
	ScoreMetrics       scoremetrics.ScoreMetrics
	RoundMetrics       roundmetrics.RoundMetrics
	LeaderboardMetrics leaderboardmetrics.LeaderboardMetrics
	EventBusMetrics    eventbusmetrics.EventBusMetrics
	DiscordMetrics     discordmetrics.DiscordMetrics
	GuildMetrics       guildmetrics.GuildMetrics
}

func NewRegistry(provider *Provider, cfg Config) *Registry {
	meter := provider.MeterProvider.Meter(cfg.ServiceName)
	tracer := provider.TracerProvider.Tracer(cfg.ServiceName)

	userMetrics := usermetrics.NewNoop() // fallback
	scoreMetrics := scoremetrics.NewNoop()
	roundMetrics := roundmetrics.NewNoop()
	leaderboardMetrics := leaderboardmetrics.NewNoop()
	eventbusMetrics := eventbusmetrics.NewNoop()
	discordMetrics := discordmetrics.NewNoop()
	guildMetrics := guildmetrics.NewNoop()

	if cfg.MetricsEnabled() {
		userMetrics, _ = usermetrics.NewUserMetrics(meter, cfg.ServiceName)
		scoreMetrics, _ = scoremetrics.NewScoreMetrics(meter, cfg.ServiceName)
		roundMetrics, _ = roundmetrics.NewRoundMetrics(meter, cfg.ServiceName)
		leaderboardMetrics, _ = leaderboardmetrics.NewLeaderboardMetrics(meter, cfg.ServiceName)
		eventbusMetrics, _ = eventbusmetrics.NewEventBusMetrics(meter, cfg.ServiceName)
		discordMetrics, _ = discordmetrics.NewDiscordMetrics(meter, cfg.ServiceName)
		guildMetrics, _ = guildmetrics.NewGuildMetrics(meter, cfg.ServiceName)
	}

	return &Registry{
		Logger: provider.Logger,
		Tracer: tracer,
		Meter:  meter,

		UserMetrics:        userMetrics,
		ScoreMetrics:       scoreMetrics,
		RoundMetrics:       roundMetrics,
		LeaderboardMetrics: leaderboardMetrics,
		EventBusMetrics:    eventbusMetrics,
		DiscordMetrics:     discordMetrics,
		GuildMetrics:       guildMetrics,
	}
}
