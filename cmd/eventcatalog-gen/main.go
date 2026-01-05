package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// AsyncAPI spec structures
type AsyncAPISpec struct {
	Info       Info               `yaml:"info"`
	Channels   map[string]Channel `yaml:"channels"`
	Components Components         `yaml:"components"`
}

type Info struct {
	Title       string `yaml:"title"`
	Version     string `yaml:"version"`
	Description string `yaml:"description"`
}

type Channel struct {
	Subscribe *Operation `yaml:"subscribe"`
	Publish   *Operation `yaml:"publish"`
}

type Operation struct {
	Message MessageRef `yaml:"message"`
}

type MessageRef struct {
	Ref string `yaml:"$ref"`
}

type Components struct {
	Schemas  map[string]interface{} `yaml:"schemas"`
	Messages map[string]Message     `yaml:"messages"`
}

type Message struct {
	Name        string                 `yaml:"name"`
	Title       string                 `yaml:"title"`
	Summary     string                 `yaml:"summary"`
	Description string                 `yaml:"description"`
	ContentType string                 `yaml:"contentType"`
	Payload     map[string]interface{} `yaml:"payload"`
}

// EventCatalog structures
type EventFrontmatter struct {
	ID         string       `yaml:"id"`
	Name       string       `yaml:"name"`
	Version    string       `yaml:"version"`
	Summary    string       `yaml:"summary,omitempty"`
	Owners     []string     `yaml:"owners,omitempty"`
	Badges     []Badge      `yaml:"badges,omitempty"`
	SchemaPath string       `yaml:"schemaPath,omitempty"`
	Channels   []ChannelRef `yaml:"channels,omitempty"`
}

type Badge struct {
	Content         string `yaml:"content"`
	BackgroundColor string `yaml:"backgroundColor"`
	TextColor       string `yaml:"textColor"`
}

type ChannelRef struct {
	ID      string `yaml:"id"`
	Version string `yaml:"version"`
}

type ServiceFrontmatter struct {
	ID       string     `yaml:"id"`
	Name     string     `yaml:"name"`
	Version  string     `yaml:"version"`
	Summary  string     `yaml:"summary,omitempty"`
	Owners   []string   `yaml:"owners,omitempty"`
	Badges   []Badge    `yaml:"badges,omitempty"`
	Sends    []EventRef `yaml:"sends,omitempty"`
	Receives []EventRef `yaml:"receives,omitempty"`
}

type EventRef struct {
	ID      string `yaml:"id"`
	Version string `yaml:"version"`
}

type DomainFrontmatter struct {
	ID       string       `yaml:"id"`
	Name     string       `yaml:"name"`
	Version  string       `yaml:"version"`
	Summary  string       `yaml:"summary,omitempty"`
	Owners   []string     `yaml:"owners,omitempty"`
	Badges   []Badge      `yaml:"badges,omitempty"`
	Services []ServiceRef `yaml:"services,omitempty"`
}

type ServiceRef struct {
	ID      string `yaml:"id"`
	Version string `yaml:"version"`
}

type ChannelFrontmatter struct {
	ID        string   `yaml:"id"`
	Name      string   `yaml:"name"`
	Version   string   `yaml:"version"`
	Summary   string   `yaml:"summary,omitempty"`
	Address   string   `yaml:"address"`
	Protocols []string `yaml:"protocols,omitempty"`
}

// Domain definitions with colors
var domains = map[string]struct {
	Name        string
	Description string
	Color       string
}{
	"round": {
		Name:        "Round",
		Description: "Manages disc golf rounds including creation, scheduling, participant management, and scoring workflows.",
		Color:       "#10b981",
	},
	"score": {
		Name:        "Score",
		Description: "Handles score processing, validation, and storage for disc golf rounds.",
		Color:       "#f59e0b",
	},
	"user": {
		Name:        "User",
		Description: "Manages user accounts, profiles, and tag number assignments.",
		Color:       "#3b82f6",
	},
	"leaderboard": {
		Name:        "Leaderboard",
		Description: "Manages leaderboard updates and tag number swaps based on round results.",
		Color:       "#8b5cf6",
	},
	"guild": {
		Name:        "Guild",
		Description: "Handles Discord guild configuration and settings management.",
		Color:       "#ec4899",
	},
}

func main() {
	// Read AsyncAPI spec
	specPath := "asyncapi/asyncapi.yaml"
	data, err := os.ReadFile(specPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading AsyncAPI spec: %v\n", err)
		os.Exit(1)
	}

	var spec AsyncAPISpec
	if err := yaml.Unmarshal(data, &spec); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing AsyncAPI spec: %v\n", err)
		os.Exit(1)
	}

	outputDir := "eventcatalog"

	// Clean existing generated content
	cleanDirectories(outputDir)

	// Generate team
	if err := generateTeam(outputDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating team: %v\n", err)
		os.Exit(1)
	}

	// Generate single NATS channel
	if err := generateNATSChannel(outputDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating channel: %v\n", err)
		os.Exit(1)
	}

	// Generate events (at root level, not nested under domains)
	if err := generateEvents(spec, outputDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating events: %v\n", err)
		os.Exit(1)
	}

	// Generate services
	if err := generateServices(spec, outputDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating services: %v\n", err)
		os.Exit(1)
	}

	// Generate domains
	if err := generateDomains(outputDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating domains: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("EventCatalog content generated successfully!")
}

func cleanDirectories(outputDir string) {
	dirs := []string{"events", "channels", "services", "domains", "teams"}
	for _, dir := range dirs {
		path := filepath.Join(outputDir, dir)
		os.RemoveAll(path)
	}
}

func generateTeam(outputDir string) error {
	teamsDir := filepath.Join(outputDir, "teams")
	if err := os.MkdirAll(teamsDir, 0755); err != nil {
		return fmt.Errorf("failed to create teams directory: %w", err)
	}

	content := `---
id: frolf-bot-team
name: Frolf Bot Team
summary: Team responsible for the Frolf Bot system
---

The Frolf Bot Team is responsible for maintaining and developing the Frolf Bot system.
`

	if err := os.WriteFile(filepath.Join(teamsDir, "frolf-bot-team.mdx"), []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write team file: %w", err)
	}

	fmt.Println("Generated team")
	return nil
}

func generateNATSChannel(outputDir string) error {
	channelsDir := filepath.Join(outputDir, "channels")
	channelPath := filepath.Join(channelsDir, "nats-jetstream")
	if err := os.MkdirAll(channelPath, 0755); err != nil {
		return fmt.Errorf("failed to create channel directory: %w", err)
	}

	frontmatter := ChannelFrontmatter{
		ID:        "nats-jetstream",
		Name:      "NATS JetStream",
		Version:   "1.0.0",
		Summary:   "Primary message broker for all Frolf Bot event communication",
		Address:   "{domain}.{event}.v1",
		Protocols: []string{"nats"},
	}

	content := generateMarkdown(frontmatter, `
NATS JetStream serves as the central message broker for the Frolf Bot event-driven architecture.

## Overview

All events in the system flow through NATS JetStream, providing:
- **Persistence**: Messages are stored durably
- **Replay**: Consumers can replay messages from any point
- **Acknowledgment**: Guaranteed delivery with ack/nack support
- **Consumer Groups**: Multiple instances can share workload

## Topic Naming Convention

Events follow the pattern: `+"`{domain}.{event-type}.v1`"+`

Examples:
- `+"`round.created.v1`"+`
- `+"`user.creation.requested.v1`"+`
- `+"`leaderboard.updated.v1`"+`

## Domains

Events are organized into the following domains:
- **round** - Round lifecycle management
- **score** - Score processing
- **user** - User management
- **leaderboard** - Leaderboard updates
- **guild** - Discord guild configuration
`)

	if err := os.WriteFile(filepath.Join(channelPath, "index.mdx"), []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write channel file: %w", err)
	}

	fmt.Println("Generated NATS JetStream channel")
	return nil
}

func generateEvents(spec AsyncAPISpec, outputDir string) error {
	eventsDir := filepath.Join(outputDir, "events")
	eventCount := 0

	for channelName, channel := range spec.Channels {
		// Get the message reference
		var msgRef string
		if channel.Publish != nil {
			msgRef = channel.Publish.Message.Ref
		} else if channel.Subscribe != nil {
			msgRef = channel.Subscribe.Message.Ref
		}

		// Extract schema name from $ref (pointing to schemas, not messages)
		schemaName := extractRefName(msgRef)
		if schemaName == "" {
			continue
		}

		// Try to get schema - first check schemas, then check if it's in messages
		var schema interface{}
		var summary, description string

		// Check if there's a message with summary/description
		if msg, ok := spec.Components.Messages[schemaName]; ok {
			summary = msg.Summary
			description = msg.Description
			// Get the schema from the message's payload ref
			if msg.Payload != nil {
				if ref, ok := msg.Payload["$ref"].(string); ok {
					payloadSchemaName := extractRefName(ref)
					schema = spec.Components.Schemas[payloadSchemaName]
				} else {
					schema = msg.Payload
				}
			}
		}

		// If no schema from message, try direct schema lookup
		if schema == nil {
			schema = spec.Components.Schemas[schemaName]
		}

		// Parse domain from channel name
		parts := strings.Split(channelName, ".")
		domain := parts[0]

		// Create event directory at root level (not nested under domain)
		eventID := strings.ReplaceAll(channelName, ".", "-")
		eventPath := filepath.Join(eventsDir, eventID)
		if err := os.MkdirAll(eventPath, 0755); err != nil {
			return fmt.Errorf("failed to create event directory: %w", err)
		}

		// Generate friendly event name
		eventName := formatEventName(channelName)
		if summary != "" {
			eventName = summary
		}
		if description == "" {
			description = generateEventDescription(channelName, domain)
		}

		// Determine badge color based on domain
		badgeColor := getDomainColor(domain)

		frontmatter := EventFrontmatter{
			ID:         eventID,
			Name:       eventName,
			Version:    "1.0.0",
			Summary:    description,
			Owners:     []string{"frolf-bot-team"},
			SchemaPath: "schema.json",
			Badges: []Badge{
				{Content: domain, BackgroundColor: badgeColor, TextColor: "#ffffff"},
				{Content: "v1", BackgroundColor: "#6366f1", TextColor: "#ffffff"},
			},
			Channels: []ChannelRef{
				{ID: "nats-jetstream", Version: "1.0.0"},
			},
		}

		// Write schema.json file
		if schema != nil {
			schemaBytes, err := json.MarshalIndent(schema, "", "  ")
			if err == nil {
				schemaFile := filepath.Join(eventPath, "schema.json")
				if err := os.WriteFile(schemaFile, schemaBytes, 0644); err != nil {
					return fmt.Errorf("failed to write schema file: %w", err)
				}
			}
		}

		// Generate markdown content with SchemaViewer component
		content := generateMarkdown(frontmatter, fmt.Sprintf(`
%s

## Topic

This event is published on the `+"`%s`"+` topic via NATS JetStream.

## Schema

<SchemaViewer file="schema.json" />

## Producers & Consumers

Use the node graph above to see which services produce and consume this event.
`, description, channelName))

		if err := os.WriteFile(filepath.Join(eventPath, "index.mdx"), []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write event file: %w", err)
		}

		eventCount++
	}

	fmt.Printf("Generated %d events\n", eventCount)
	return nil
}

func generateServices(spec AsyncAPISpec, outputDir string) error {
	servicesDir := filepath.Join(outputDir, "services")

	// Categorize events by type (request vs response)
	backendSends := []string{}
	backendReceives := []string{}
	discordSends := []string{}
	discordReceives := []string{}

	for channelName := range spec.Channels {
		eventID := strings.ReplaceAll(channelName, ".", "-")

		// Classify based on event name patterns
		if isRequestEvent(channelName) {
			// Request events are sent by Discord bot, received by backend
			discordSends = append(discordSends, eventID)
			backendReceives = append(backendReceives, eventID)
		} else {
			// Response/notification events are sent by backend, received by Discord bot
			backendSends = append(backendSends, eventID)
			discordReceives = append(discordReceives, eventID)
		}
	}

	// Service definitions
	services := []struct {
		ID          string
		Name        string
		Description string
		Sends       []string
		Receives    []string
		Color       string
	}{
		{
			ID:          "frolf-bot",
			Name:        "Frolf Bot Backend",
			Description: "Core backend service that handles business logic for rounds, scores, users, leaderboards, and guilds. Built with Go and Watermill, it processes requests from the Discord bot and emits domain events.",
			Sends:       backendSends,
			Receives:    backendReceives,
			Color:       "#059669",
		},
		{
			ID:          "discord-frolf-bot",
			Name:        "Discord Frolf Bot",
			Description: "Discord bot that provides the user interface for interacting with the Frolf Bot system. Handles slash commands, button interactions, and displays results to users in Discord channels.",
			Sends:       discordSends,
			Receives:    discordReceives,
			Color:       "#5865F2",
		},
	}

	for _, svc := range services {
		servicePath := filepath.Join(servicesDir, svc.ID)
		if err := os.MkdirAll(servicePath, 0755); err != nil {
			return fmt.Errorf("failed to create service directory: %w", err)
		}

		sends := make([]EventRef, len(svc.Sends))
		for i, e := range svc.Sends {
			sends[i] = EventRef{ID: e, Version: "1.0.0"}
		}

		receives := make([]EventRef, len(svc.Receives))
		for i, e := range svc.Receives {
			receives[i] = EventRef{ID: e, Version: "1.0.0"}
		}

		frontmatter := ServiceFrontmatter{
			ID:      svc.ID,
			Name:    svc.Name,
			Version: "1.0.0",
			Summary: svc.Description,
			Owners:  []string{"frolf-bot-team"},
			Badges: []Badge{
				{Content: "Go", BackgroundColor: "#00ADD8", TextColor: "#ffffff"},
				{Content: "Watermill", BackgroundColor: "#4A90D9", TextColor: "#ffffff"},
			},
			Sends:    sends,
			Receives: receives,
		}

		content := generateMarkdown(frontmatter, fmt.Sprintf(`
%s

## Architecture

<NodeGraph />

## Technology Stack

| Component | Technology |
|-----------|------------|
| Language | Go |
| Message Broker | NATS JetStream |
| Event Library | Watermill |
| Database | PostgreSQL |

## Event Summary

| Direction | Count |
|-----------|-------|
| **Publishes** | %d events |
| **Subscribes** | %d events |

## Integration

This service communicates exclusively through events via NATS JetStream. It does not expose HTTP APIs directly.
`, svc.Description, len(svc.Sends), len(svc.Receives)))

		if err := os.WriteFile(filepath.Join(servicePath, "index.mdx"), []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write service file: %w", err)
		}
	}

	fmt.Printf("Generated %d services\n", len(services))
	return nil
}

func generateDomains(outputDir string) error {
	domainsDir := filepath.Join(outputDir, "domains")

	for domainID, domain := range domains {
		domainPath := filepath.Join(domainsDir, domainID)
		if err := os.MkdirAll(domainPath, 0755); err != nil {
			return fmt.Errorf("failed to create domain directory: %w", err)
		}

		frontmatter := DomainFrontmatter{
			ID:      domainID,
			Name:    domain.Name,
			Version: "1.0.0",
			Summary: domain.Description,
			Owners:  []string{"frolf-bot-team"},
			Badges: []Badge{
				{Content: "DDD", BackgroundColor: domain.Color, TextColor: "#ffffff"},
			},
			Services: []ServiceRef{
				{ID: "frolf-bot", Version: "1.0.0"},
				{ID: "discord-frolf-bot", Version: "1.0.0"},
			},
		}

		content := generateMarkdown(frontmatter, fmt.Sprintf(`
%s

## Overview

The **%s** domain is a bounded context within the Frolf Bot system that encapsulates all business logic related to %s.

## Architecture

<NodeGraph />

## Event Patterns

All events in this domain follow:
- **Event Notification Pattern** (Martin Fowler) - Events contain complete data snapshots
- **V1 Versioning** - Topic names include `+"`"+`.v1`+"`"+` suffix for version compatibility

## Services

Both the backend service (frolf-bot) and Discord bot (discord-frolf-bot) participate in this domain's workflows.
`, domain.Description, domain.Name, strings.ToLower(domain.Name)))

		if err := os.WriteFile(filepath.Join(domainPath, "index.mdx"), []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write domain file: %w", err)
		}
	}

	fmt.Printf("Generated %d domains\n", len(domains))
	return nil
}

func generateMarkdown(frontmatter interface{}, body string) string {
	fm, _ := yaml.Marshal(frontmatter)
	return fmt.Sprintf("---\n%s---\n%s", string(fm), body)
}

func getDomainColor(domain string) string {
	if d, ok := domains[domain]; ok {
		return d.Color
	}
	return "#6b7280" // gray default
}

// extractRefName extracts the schema name from a $ref like "#/components/messages/RoundCreatedV1"
func extractRefName(ref string) string {
	if ref == "" {
		return ""
	}
	parts := strings.Split(ref, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}

// formatEventName converts "round.created.v1" to "Round Created"
func formatEventName(channelName string) string {
	// Remove the .v1 suffix
	name := strings.TrimSuffix(channelName, ".v1")

	// Split by dots
	parts := strings.Split(name, ".")

	// Capitalize each part
	result := make([]string, len(parts))
	for i, part := range parts {
		if len(part) > 0 {
			result[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}

	return strings.Join(result, " ")
}

// generateEventDescription creates a descriptive summary based on the event name
func generateEventDescription(channelName, domain string) string {
	// Remove domain prefix and .v1 suffix
	name := strings.TrimSuffix(channelName, ".v1")
	name = strings.TrimPrefix(name, domain+".")

	// Convert to readable format
	name = strings.ReplaceAll(name, ".", " ")
	name = strings.ReplaceAll(name, "-", " ")

	domainInfo := domains[domain]

	// Create description based on patterns
	switch {
	case strings.Contains(name, "requested"):
		return fmt.Sprintf("Command event requesting a %s operation. Typically published by the Discord bot and consumed by the backend service.", name)
	case strings.Contains(name, "created"):
		return fmt.Sprintf("Domain event indicating a new %s entity was successfully created.", domain)
	case strings.Contains(name, "updated"):
		return fmt.Sprintf("Domain event indicating a %s entity was modified.", domain)
	case strings.Contains(name, "deleted"):
		return fmt.Sprintf("Domain event indicating a %s entity was removed.", domain)
	case strings.Contains(name, "failed"):
		return fmt.Sprintf("Error event indicating a %s operation did not complete successfully.", domain)
	case strings.Contains(name, "validated"):
		return fmt.Sprintf("Internal event indicating validation passed for a %s operation.", domain)
	case strings.Contains(name, "response"):
		return fmt.Sprintf("Response event containing requested %s data.", domain)
	default:
		return fmt.Sprintf("Event in the %s domain. %s", domainInfo.Name, domainInfo.Description)
	}
}

// isRequestEvent determines if an event is a request (sent by Discord bot) or response (sent by backend)
func isRequestEvent(channelName string) bool {
	requestPatterns := []string{
		"requested",
		"request",
	}

	nameLower := strings.ToLower(channelName)
	for _, pattern := range requestPatterns {
		if strings.Contains(nameLower, pattern) {
			return true
		}
	}
	return false
}
