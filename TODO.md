# TODO

Staged development plan for mlpipe. Estimated total: about 103 hours.

## Stage 0: RabbitMQ Ramp-up (6 hrs)

**Concepts to learn first**
- Producers, exchanges (direct, fanout, topic), queues, bindings, and routing keys
- Message acknowledgements (ack, nack, requeue) and prefetch (QoS)
- Durable queues and persistent messages, and what each one guarantees
- Dead-letter exchanges and per-message TTL
- Publisher confirms

**Tasks**
- [x] Run RabbitMQ in Docker and use the management UI
- [ ] Write a small Go producer and consumer with the `amqp091-go` client
- [ ] Demonstrate manual ack, nack with requeue, and prefetch with two consumers
- [ ] Set up a dead-letter exchange and confirm rejected messages land in it

---

## Stage 1: Project Foundation (9 hrs)

**Concepts to learn first**
- Go project layout (`cmd/`, `internal/`) and one binary with subcommands
- Config from environment variables (12-factor style)
- Structured logging with `log/slog`
- Docker Compose networking, health checks, and `depends_on` with conditions
- Migrations with `golang-migrate`

**Tasks**
- [ ] Set up the repo layout and a single binary with `api`, `relay`, and `worker` subcommands
- [ ] Add config loading and structured logging with correlation-ID fields
- [ ] Add `/healthz` (liveness) and `/readyz` (checks Postgres, MinIO, and RabbitMQ)
- [ ] Write a Compose file with Postgres, MinIO, RabbitMQ, and the three services
- [ ] Create the first migration

---

## Stage 2: Datasets and Upload API (13 hrs)

**Concepts to learn first**
- Presigned URLs and how they keep bytes off the API server
- Content-addressed storage (SHA-256 as the object key) and why it makes uploads idempotent
- MinIO through the S3-compatible Go SDK
- Bulk upload: one request returns one presigned URL per file

**Tasks**
- [ ] Create the `datasets` table (id, name, modality, status, created_at)
- [ ] `POST /datasets` and `GET /datasets/{id}`, plus paginated listing
- [ ] `POST /datasets/{id}/uploads` for single or bulk files, returning presigned PUT URLs
- [ ] `POST /datasets/{id}/uploads/{fileId}/complete` to verify the object exists and record its SHA-256 and size
- [ ] Store file metadata in a `files` table with a unique `(dataset_id, sha256)` constraint
- [ ] Add a presigned GET download endpoint
- [ ] Add API-key middleware

---

## Stage 3: Transactional Outbox and Relay (9 hrs)

**Concepts to learn first**
- The dual-write problem: why saving to the database and publishing to a queue separately can lose events
- The transactional outbox pattern
- `SELECT ... FOR UPDATE SKIP LOCKED` for safe polling by multiple relays
- At-least-once delivery and why consumers must be idempotent

**Tasks**
- [ ] Create the `outbox` table (id, type, payload JSONB, created_at, published_at)
- [ ] Write the `file.completed` event in the same transaction as the file status update
- [ ] Build the relay: poll unpublished rows with `SKIP LOCKED`, publish with confirms, then mark them published
- [ ] Test by killing the relay between publish and mark. The event should be sent again, not lost.

---

## Stage 4: Modality Adapters (12 hrs)

**Concepts to learn first**
- Interface design in Go and a registry keyed by name
- Streaming readers (`io.Reader`) so large files aren't loaded into memory
- Image decoding and encoding with the Go `image` packages (or a library such as `disintegration/imaging`)
- Perceptual hashing (dHash) and how it differs from SHA-256
- CSV parsing edge cases: quoting, encodings, ragged rows

**Tasks**
- [ ] Define the `Modality` interface: `Validate`, `ExtractMetadata`, `Transforms`, `ApplyTransform`
- [ ] Image adapter: format and corruption checks, dimension and EXIF metadata, resize, center crop, PNG/JPEG/WebP conversion, dHash
- [ ] CSV adapter: header detection, row and column counts, null rates per column, type inference, ragged-row detection
- [ ] Text adapter (JSONL and `.txt`): UTF-8 validation, line count, length stats, malformed JSON lines reported
- [ ] Exact-duplicate detection by SHA-256 for every modality
- [ ] Validate each dataset's transform config against its modality's `Transforms()` at creation time

---

## Stage 5: Workers and Reliability (15 hrs)

**Concepts to learn first**
- Consumer prefetch and manual acks, and what happens to unacked messages when a worker dies
- Worker pools with bounded concurrency (backpressure)
- Retries with backoff using the TTL plus dead-letter exchange pattern
- Dead-letter queues and poison-message handling
- `context` deadlines for processing timeouts
- Graceful shutdown: stop consuming, finish in-flight work, then exit

**Tasks**
- [ ] Create the `processing_stages` table with a unique `(file_id, stage)` constraint so duplicate messages are harmless
- [ ] Worker consumes `file.completed` and runs the fixed stage sequence: validate → hash → metadata → transform
- [ ] Set the prefetch count and a configurable concurrency limit
- [ ] Retry transient errors with exponential backoff and jitter, using a delay queue
- [ ] Route messages that fail after N attempts to a DLQ and record the error on the stage row
- [ ] Enforce per-stage timeouts with `context`
- [ ] Handle SIGTERM: stop consuming, wait for in-flight work, then exit
- [ ] Run at least two worker containers and confirm they share the load

---

## Stage 6: Dataset Versioning and Splits (11 hrs)

**Concepts to learn first**
- Immutable versions and lineage (source file → processed output → version)
- Deterministic hashing for partitioning (`hash(fileID + seed)`)
- Designing a reproducible manifest format

**Tasks**
- [ ] Create `dataset_versions` and `version_files` tables. Versions reference files and store the transform config.
- [ ] `POST /datasets/{id}/versions` builds a version once its files finish processing
- [ ] Deterministic train/val/test split from a stored seed and ratios
- [ ] `GET /datasets/{id}/versions/{v}/manifest` returns files, hashes, and splits as JSON
- [ ] Publish a `dataset.version.ready` event through the outbox

---

## Stage 7: Processing Requests and Status APIs (7 hrs)

**Concepts to learn first**
- Request-to-job modeling: one request creates a pipeline run with many file-level stages
- Status rollups from child rows to a parent run

**Tasks**
- [ ] `POST /datasets/{id}/processing` starts a run and emits a `processing.requested` event
- [ ] Pipeline run table with status and timestamps
- [ ] `GET /runs/{id}` shows overall status and per-stage counts
- [ ] `GET /runs/{id}/history` lists state transitions

---

## Stage 8: Observability (7 hrs)

**Concepts to learn first**
- Prometheus metric types: counter, gauge, histogram
- Correlation ID propagation through HTTP headers and RabbitMQ message headers
- Consumer lag and queue depth from the RabbitMQ management API

**Tasks**
- [ ] Expose `/metrics` on each service
- [ ] Metrics: uploads, events published, messages consumed, stage duration histogram, retries, DLQ count
- [ ] Carry the correlation ID from the upload request to every worker log line
- [ ] Add a Prometheus config and a dashboard screenshot for the README (Grafana optional)

---

## Stage 9: Testing and Benchmarks (8 hrs)

**Concepts to learn first**
- Integration tests with `testcontainers-go`
- Injecting failures: killing containers mid-job, restarting RabbitMQ

**Tasks**
- [ ] Unit tests for each modality adapter with good and corrupted samples
- [ ] Integration test: upload → outbox → relay → worker → version
- [ ] Duplicate-delivery test: send the same message twice and confirm the final state is unchanged
- [ ] Crash test: kill a worker mid-job and verify the message is redelivered and completed
- [ ] Benchmark: run 1, 2, and 4 workers over the same dataset and record throughput

---

## Stage 10: Documentation and Release (6 hrs)

**Tasks**
- [ ] Expand README with setup, a quick start using curl, and an architecture diagram
- [ ] Document the API and event schemas
- [ ] Document a failure-mode table (what happens on crash, duplicate, or outage)
- [ ] Add sample datasets (a few images, a CSV, a JSONL file) and an example config
- [ ] Publish on GitHub

---

**To trim if needed:** drop the Grafana part of Stage 8 first, then reduce the benchmark in Stage 9.

**Recommended order:** Finish Stages 0 through 3 before writing any modality code. The queue and outbox lessons are the main reason for this project.