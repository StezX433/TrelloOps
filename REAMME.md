# TrelloOps

TrelloOps is a Go backend automation service for Trello.

It monitors cards in a configured Trello list and automatically identifies
cards that have been inactive for a configurable amount of time.

The current automation rule is:

> If a card has been inactive in the `Doing` list for 5 or more days,
> apply the `At Risk` label.

---

## Features

- Trello REST API integration
- Board and list discovery
- Card retrieval
- Stale-card detection
- Configurable inactivity threshold
- Automatic `At Risk` label creation
- Automatic label assignment
- Duplicate-label protection
- Dry-run mode
- Unit tests
- Environment-based configuration
- No Trello credentials stored in source code

---

## Architecture

```text
Trello
   |
   v
Trello REST API
   |
   v
internal/trello
   |
   v
internal/service
   |
   +------> internal/rules
   |
   v
At Risk label