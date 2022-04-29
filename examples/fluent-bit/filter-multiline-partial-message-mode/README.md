### FireLens Example: Concatenate Partial/Split Container Logs

The `partial_message` `mode` for the Multiline Filter is available on `aws-for-fluent-bit` >= 2.24.0

FireLens uses the [fluentd log driver](https://docs.docker.com/config/containers/logging/fluentd/) in its [implementation](https://aws.amazon.com/blogs/containers/under-the-hood-firelens-for-amazon-ecs-tasks/). Stdout/stderr logs pass from the container runtime to the log driver, to Fluent Bit. The container runtime denotes the stream of bytes from stdout and stderr as a series of log events by splitting the stream on bytes, or at 16KB. When logs are split at 16KB they become partial messages. For the sake of the example, assume the runtime's max buffer size is 4 bytes instead of 16KB; Fluent Bit would get the following series of messages:

```
{"source": "stdout", "log": "Sher", "partial_message": "true", "partial_id": "dc37eb08b4242c41757d4cd995d983d1cdda4589193755a22fcf47a638317da0", "partial_ordinal": "1", "partial_last": "false", "container_id": "a96998303938eab6087a7f8487ca40350f2c252559bc6047569a0b11b936f0f2", "container_name": "/hopeful_taussig"}]
{"partial_last": "false", "container_id": "a96998303938eab6087a7f8487ca40350f2c252559bc6047569a0b11b936f0f2", "container_name": "/hopeful_taussig", "source": "stdout", "log": "lock", "partial_message": "true", "partial_id": "dc37eb08b4242c41757d4cd995d983d1cdda4589193755a22fcf47a638317da0", "partial_ordinal": "2"}]
{"log": " Hol", "partial_message": "true", "partial_id": "dc37eb08b4242c41757d4cd995d983d1cdda4589193755a22fcf47a638317da0", "partial_ordinal": "3", "partial_last": "false", "container_id": "a96998303938eab6087a7f8487ca40350f2c252559bc6047569a0b11b936f0f2", "container_name": "/hopeful_taussig", "source": "stdout"}]
{"container_id": "a96998303938eab6087a7f8487ca40350f2c252559bc6047569a0b11b936f0f2", "container_name": "/hopeful_taussig", "source": "stdout", "log": "mes", "partial_message": "true", "partial_id": "dc37eb08b4242c41757d4cd995d983d1cdda4589193755a22fcf47a638317da0", "partial_ordinal": "4", "partial_last": "true"}]
```

None of these messages were split on a newline, all were split by the max buffer size. Fluent Bit can re-combine these messages so that you get one message:
```
{"container_id": "a96998303938eab6087a7f8487ca40350f2c252559bc6047569a0b11b936f0f2", "container_name": "/hopeful_taussig", "source": "stdout", "log": "Sherlock Holmes"}]
```

This can be accomplished with a simple filter definition:

```
[FILTER]
    name                  multiline
    match                 *
    multiline.key_content log
    # partial_message mode is incompatible with option multiline.parser
    mode                  partial_message
```

The files included in this example allow you to build a custom Fluent Bit image with this config- please see the base [config-file-type](https://github.com/aws-samples/amazon-ecs-firelens-examples/tree/mainline/examples/fluent-bit/config-file-type-file) example to understand how to build an image with a custom config. 

The included `multiline-app` is just a sample app for demonstration purposes. It prints the entire text of a Sherlock Holmes novel to stdout and stderr, which leads to many partial messages. Please note that CloudWatch Logs has a max event size of 256KB and Fluent Bit will truncate single messages that are over this size. This limit cannot be changed. 

*Note*: 
- `permissions.json` shows the permissions we should include in the IAM role to send logs to CloudWatch.
- the container runtime buffer size for logs is fixed at 16KB and can not be changed


