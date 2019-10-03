### FireLens Example: Enabling Debug Logging for Fluent Bit

Log level can be set in the [Service](https://docs.fluentbit.io/manual/service) section of the Fluent Bit configuration file. This section is not used by FireLens; you can set it yourself using an external configuration file.

To enable debug logging in the AWS Fluent Bit plugins; set the environment variable `FLB_LOG_LEVEL`. It can be set to `debug`, `info`, and `error`.
