# orchestrated-order-saga

Distributed order processing system implementing the **Orchestration-based Saga pattern** with **Transactional Outbox**, built with Go and Ruby, coordinated via Kafka.

Companion project to [distributed-order-saga](https://github.com/joaodddev/distributed-order-saga), which implements the same domain using the **Choreography-based** approach. Built to demonstrate, side by side, the trade-offs between the two saga strategies.

## Architecture

A central **saga-orchestrator** owns the full state machine of the saga and drives it forward by issuing **commands** to each service and reacting to their **replies**. Unlike the choreographed version, the executor services (`order`, `payment`, `inventory`) have no knowledge that they're part of a saga — they simply execute a command and report success or failure back.

| Service | Stack | Role |
|---|---|---|
| `saga-orchestrator` | Go + Gin, MySQL | Owns saga state, issues commands, decides next step or compensation |
| `order-service` | Go, MySQL | Executes `order.confirm` / `order.cancel` commands |
| `payment-service` | Ruby + Racecar, PostgreSQL | Executes `payment.reserve` / `payment.refund` commands |
| `inventory-service` | Go, MySQL | Executes `inventory.reserve` commands |

All services use **Clean Architecture** (`domain` / `application` / `infrastructure`). The orchestrator uses **Transactional Outbox** to guarantee command delivery even if Kafka is temporarily unavailable.

## Choreography vs. Orchestration — what this project demonstrates

| | Choreography | Orchestration |
|---|---|---|
| Coupling | Low — services only know events | Higher — orchestrator knows every service |
| Single point of failure | None | The orchestrator |
| Traceability | Harder — logic spread across services | Easier — one place owns the full flow |
| Best for | Few services, simple flows | Complex flows, many steps, need for central visibility |