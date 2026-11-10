# mlpipe

A Go service where ML teams upload raw datasets, and workers validate, transform, and version them through RabbitMQ-driven pipelines. Teams retrieve versioned, ML-ready datasets through a REST API.

See [TODO.md](./TODO.md) for the staged development plan.

## Scope

- **Modalities:** images (resize, crop, convert, perceptual hash), CSV, and text (JSONL and `.txt`)
- **API:** REST only. No gRPC, no frontend.
- **Auth:** simple API keys
- **Size limits:** none. Large files are streamed, so bigger datasets just take longer.

## Core Flow

1. A client requests presigned URLs and uploads files directly to MinIO. Metadata is saved in PostgreSQL.
2. Completing an upload writes an event to an outbox table in the same transaction. A relay publishes it to RabbitMQ.
3. Independent Go workers consume events and run the fixed stage sequence: validate → hash → metadata → transform.
4. Each stage's status and errors are recorded. Transient failures retry with backoff, and permanent failures go to a dead-letter queue.
5. A dataset version is built from processed files, with a manifest and a deterministic train/validation/test split.

## Stack

- Go, REST
- RabbitMQ for events and work queues
- PostgreSQL for metadata, stage state, outbox, and versions
- MinIO for files (S3-compatible API)
- Docker Compose

## What This Project Demonstrates

- Acks, redelivery after worker crashes, and graceful shutdown
- Transactional outbox and idempotent consumers
- Retries with backoff, dead-letter queues, and processing timeouts
- Backpressure through worker concurrency limits
- Multi-instance workers and containerized deployment

## Getting Started

_To be added once Stage 1 is complete: setup, Compose commands, and a curl quick start._