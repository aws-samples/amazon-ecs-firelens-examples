# Log Generator

Container image that generates JSON log entries at configurable rates and sizes.

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `LOG_RATE_PER_SECOND` | `10` | Log generation rate per second |
| `LOG_SIZE_KB` | `7` | Size of each log entry in KB |
| `LOG_SIZE_EXTRA_KB` | `0` | Random extra size added to LOG_SIZE_KB (0 to LOG_SIZE_EXTRA_KB, 0 disables extra size) |
| `LOG_SIZE_BURST_KB` | Same as `LOG_SIZE_KB` | Size of burst log entries in KB |
| `BURST_INTERVAL_SECONDS` | `15` | Interval between burst logs in seconds |
| `OUTPUT_FILE` | `STDOUT` | Output destination (file path, STDOUT, FORWARD, or UDP) |
| `FLUENT_HOST` | `localhost` | Fluent-bit host (only used when OUTPUT_FILE=FORWARD) |
| `FLUENT_PORT` | `24224` | Fluent-bit port (only used when OUTPUT_FILE=FORWARD) |
| `FLUENT_TAG` | `video.analytics` | Fluent-bit tag (only used when OUTPUT_FILE=FORWARD) |
| `UDP_HOST` | `localhost` | UDP host (only used when OUTPUT_FILE=UDP) |
| `UDP_PORT` | `5170` | UDP port (only used when OUTPUT_FILE=UDP) |
| `LOG_ROTATE_MAX_SIZE_MB` | `100` | Maximum size in MB before log rotation (0 disables rotation) |
| `LOG_ROTATE_MAX_BACKUPS` | `5` | Maximum number of backup log files to keep |
| `LOG_ROTATE_MAX_AGE` | `1` | Maximum age in days to keep log files (0 = no age limit) |

> If both `LOG_ROTATE_MAX_BACKUPS` and `LOG_ROTATE_MAX_AGE` are 0 then no old log files will be deleted.
