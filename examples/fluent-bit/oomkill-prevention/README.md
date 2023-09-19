# Reducing the Likelihood of Out of Memory Exceptions in FireLens

In normal/happy case conditions, the memory usage of Fluent Bit will often stay below 100 MB ([see the 'Resource Usage' section in this link](https://aws.amazon.com/blogs/containers/under-the-hood-firelens-for-amazon-ecs-tasks/)). 

But if there is a spike in log output from your application that Fluent Bit cannot handle or your log destination goes down, then Fluent Bit will have to buffer lots of logs. By default, Fluent Bit buffers logs in memory, and thus, in these scenarios, the Fluent Bit container's memory usage may increase significantly. In these conditions, Fluent Bit's memory usage can often spike to up to a Gigabyte or even more. 

The following are actions that we recommend considering. 

1. [Retry_Limit](#1-retry_limit)
2. [Control Fluent Bit buffer memory](#2-control-fluent-bit-buffer-memory)
    - [Background and Recommendations](#background-and-recommendations)
        - [When should I use memory or filesystem buffer?](#when-should-i-use-memory-or-filesystem-buffer)
        - [What is a 'chunk'?](#what-is-a-chunk)
    - [Case 1: Memory Buffering Only: default or "storage.type memory"](#case-1-memory-buffering-only-default-or-storagetype-memory)
        - [Estimating Total Memory Usage of the Fluent Bit Process](#storagetype-memory-estimating-total-memory-usage-of-the-fluent-bit-process)
        - [What happens when the memory limit is reached?](#storagetype-memory-what-happens-when-the-memory-limit-is-reached)
    - [Case 2: Filesystem and Memory Buffering: "storage.type filesystem"](#case-2-filesystem-and-memory-buffering-storagetype-filesystem)
        - [What happens when the memory limit is reached?](#storagetype-filesystem-what-happens-when-the-memory-limit-is-reached)
        - [Estimating Total Memory Usage of the Fluent Bit Process](#storagetype-filesystem-estimating-total-memory-usage-of-the-fluent-bit-process)
        - [Controlling Disk Usage](#controlling-disk-usage)
        - [What happens when the disk limit is reached?](#what-happens-when-the-disk-limit-is-reached)
    - [Full FireLens Configuration Examples](#full-firelens-configuration-examples)
3. [Monitor Fluent Bit](#3-monitor-fluent-bit)


## 1. Retry_Limit

Fluent Bit allows each output to [configure a `Retry_Limit`](https://docs.fluentbit.io/manual/administration/scheduling-and-retries) for the total of retries when there are failures. Retries for chunks use exponential backoff, so with more retries, the amount of time Fluent Bit will have to buffer a chunk of logs will increase exponetially. Thus, if Fluent Bit can not send logs, a high number of retries will quickly lead to the buffers increasing dramatically. 

If Fluent Bit is retrying, you will see messages like:

```
 [ warn] [engine] failed to flush chunk '1-1647467985.551788815.flb', retry in 9 seconds: task_id=0, input=forward.1 > output=cloudwatch_logs.0 (out_id=0)
```

If retries for a chunk have expired and it will be dropped, then you will see messages like:

```
[2022/02/16 20:11:36] [ warn] [engine] chunk '1-1645042288.260516436.flb' cannot be retried: task_id=0, input=tcp.3 > output=cloudwatch_logs.1
```

## 2. Control Fluent Bit buffer memory

### Background and Recommendations

#### When should I use memory or filesystem buffer?

The default buffer type is `memory`; `filesystem` buffering can optionally be enabled. Please read this full guide to understand both buffer types and how to tune their behavior. 

In general, AWS Fluent Bit team recommends configuring `filesystem` buffering for the following reasons. Please carefully consider these benefits, as they depend on your exact use case.
* **Larger buffer space**: If you have more free disk space available than free memory space (or if disk is cheaper than memory), then choosing filesystem buffering enables a larger buffer. During a spike in throughput or a momentary failure in the logging destination, Fluent Bit will be able to buffer logs longer and the risk of log loss will be reduced.
* **Buffer can be recovered after restart**: With a filesystem buffer, if Fluent Bit stops and restarts, it can pick up the existing buffer files on the filesystem. This consideration only matters if Fluent Bit will be restarted with the same disk and volume mounts; consider whether your use case involves restarting Fluent Bit with access to the same files. In Kubernetes/EKS, this is very common as Fluent Bit is typically runs as a daemonset. 
* **Total Fluent Bit memory usage can be more restricted**: With the filesystem buffer, Fluent Bit can be configured to use less memory than if the memory buffer is configured. Please note that the filesystem buffer does still use memory for buffering as well and the `storage.max_chunks_up` setting must be carefully chosen; read [storage.type filesystem: What happens when the memory limit is reached?](#storagetype-filesystem-what-happens-when-the-memory-limit-is-reached) 
* **Inputs do not need to be paused**: By default, when the memory buffer is configured and the limit is reached, the input plugin is paused. An [overlimit warning](https://github.com/aws/aws-for-fluent-bit/blob/mainline/troubleshooting/debugging.md#overlimit-warnings) will be logged and [all input plugins will lose data](https://github.com/aws/aws-for-fluent-bit/blob/mainline/troubleshooting/debugging.md#overlimit-warnings), except for [tail](https://docs.fluentbit.io/manual/pipeline/inputs/tail), which can save its file offset before it is paused. With filesystem buffering, by default the inputs will not be paused when limits are reached. They can keep ingesting new data. Please carefully read:
    * [storage.type memory: What happens when the memory limit is reached?](#storagetype-memory-what-happens-when-the-memory-limit-is-reached)
    * [storage.type filesystem: What happens when the memory limit is reached?](#storagetype-filesystem-what-happens-when-the-memory-limit-is-reached)
    * [What happens when the disk limit is reached?](#what-happens-when-the-disk-limit-is-reached)

There are cases where the default `memory` buffer is preferable:
* **Using tail input**: With the tail input, logs are already ingested from the filesystem. Storing the logs on disk twice, once in the original files and once in the Fluent Bit buffer, may not make sense and simply duplicates disk space usage. The tail input will store its file byte offset for each file to track its progress; with the `DB` option the offset is persisted in a local SQLite DB file. With `DB.locking false`, the SQLite DB can be read by other processes.
* **Limited available disk space**: If your available disk space is limited, then using filesystem buffering may not be ideal.
* **Limited disk IOPs**: If many other processes are reading and writing to the disk, then disk IOPs may be saturated. This is a [known issue](https://github.com/aws/aws-for-fluent-bit/blob/mainline/troubleshooting/debugging.md#filesystem-buffering-and-log-files-saturate-node-disk-iops) that prevents some users from enabling filesystem buffering. 

Finally, it should also be noted that enabled filesystem buffering adds a slight performance hit compared to memory buffering, since Fluent Bit makes disk IO calls. However, the difference in performance is generally trivial (but has never been benchmarked) and AWS for Fluent Bit team has never witnessed a customer case where the difference in performance was significant enough to matter.

#### What is a 'chunk'?

Throughput this guide, and the Fluent Bit documentation, there are many uses of the word 'chunk'. This term is overloaded; internally in Fluent Bit there are multiple possible meanings for 'chunk'.

* **internal buffer chunk**: When Fluent Bit ingests data, it stores it temporarily (after applying filters) before sending data to the outputs. This is the meaning of 'chunk' used throughput this guide on buffering. These internal buffer chunks are managed by the Fluent Bit engine; it accepts data from inputs, and appends it to chunks. Each chunk is targeted to be 2MB in size; *this size is not configurable*. *The 2MB size does not restrict Fluent Bit from ingesting logs larger than 2 MB*; the target size can be exceeded to store a single large log. For each unique tag value for each of the inputs that are actively ingesting data, there should be at least one active chunk to collect the data for each tag. When the Fluent Bit engine invokes an output plugin to send data, it sends single chunks at a time. There are only a few configuration options that affect internal buffer chunks. The settings discussed in this guide explain how to control where chunks are stored (memory vs filesystem) and the behavior of Fluent Bit when limits are reached. The [`Flush` parameter](https://docs.fluentbit.io/manual/administration/configuring-fluent-bit/classic-mode/configuration-file) controls the frequency at which Fluent Bit attempts to send buffer chunks. 
* **chunks in plugins**: In plugin documentation and examples, the word 'chunk' is used more generally to mean a unit of data. These usages of chunk are all different than the internal buffer chunks discussed above. Plugins do not have any control over internal buffer chunks. For example, the [S3 output documentation](https://docs.fluentbit.io/manual/pipeline/outputs/s3) uses the word 'chunk' to refer to a single part in a multipart upload with the `upload_chunk_size` parameter. The TCP input uses the word ['chunk' in a setting](https://github.com/aws/aws-for-fluent-bit/blob/mainline/troubleshooting/debugging.md#chunk_size-and-buffer_size-for-large-logs-in-tcp-input) that controls memory allocation. 

### Case 1: Memory Buffering Only: default or "storage.type memory"

*Full example config file:* [**memory.conf**](memory.conf)

**Note:** This guidance applies if you have chosen the default `memory` buffer type for `storage.type`. If you choose `filesystem` then skip to the next section which covers filesystem buffering.

To prevent high memory usage, Fluent Bit has [a setting called `Mem_Buf_Limit`](https://docs.fluentbit.io/manual/administration/backpressure) which can be set on each input definition. 

This setting restricts the total amount of memory that the input can use- Fluent BIt considers inputs to "own" the buffers for logs, so this setting controls the size of the memory buffer that each input can use. Therefore, to fully restrict the memory usage fully for Fluent Bit you must set `Mem_Buf_Limit` on every input definition. Though in practice, if you set `Mem_Buf_Limit` on only the inputs that experience the highest throughput of logs this should mostly contain its total memory usage. 

#### storage.type memory: Estimating Total Memory Usage of the Fluent Bit Process

This setting only restricts the input memory, not the memory used by filters and outputs. However, in practice, if the input can not buffer more than a certain amount of logs, this effectively controls the memory usage of all of the other plugins in the Fluent Bit pipeline, since inputs control log ingestion. (You can picture this with a simple scenario, an input ingests one 2 MB chunk, and then while that is sitting in memory, an output tries to send it and thus has to create a request buffer which copies the chunk into a new 2 MB memory block).

To be clear, in Fluent Bit there are the following types of memory:
1. **Input Plugin Buffer Memory**: Memory used by inputs to buffer incoming logs. *This can be directly limited and controlled with `Mem_Buf_Limit`*. 
2. **Filter Plugin Processing Memory**: Memory used by filter plugins to process your log records. When filters make modifications to your records, they must temporarily copy them and then alter the copy. *This memory can not be directly controlled*, however, controlling the input memory indirectly controls the amount of logs flowing through filters. 
3. **Output Plugin Processing Memory**: Memory used by output plugins to process and send your log records. When outputs send your log records, they must create request/send buffers. *This memory can not be directly controlled*, however, controlling the input memory indirectly controls the amount of logs flowing through outputs. 
4. **Internal Fluent Bit Engine Memory**: Memory used by the Fluent Bit engine and scheduler to manage all of the plugins. This memory usage is generally small, and does generally increase as log throughput increases. *This memory usage can not be directly controlled.*

AWS FireLens team has found that in general, the total memory usage of Fluent Bit will stay below two times the sum of the `Mem_Buf_Limit` for each input. This is a rough estimate for sum of all types of memory usage noted above:

```
Total Max Memory Usage <= 2 * SUM(Each input Mem_Buf_Limit)
```

Thus, if you give Fluent Bit a 200 MB memory limit in your container definition, then the sum of all `Mem_Buf_Limit` should not exceed 100 MB. 

*Please remember that this is an estimation method and not a guarantee*.

#### storage.type memory: What happens when the memory limit is reached?

Of course, there are downsides to `Mem_Buf_Limit`. If an input runs out of memory buffer, it stop ingesting logs and [emit messages](https://github.com/aws/aws-for-fluent-bit/blob/mainline/troubleshooting/debugging.md#overlimit-warnings) like the following:

```
[input] forward.1 paused (mem buf overlimit)
```

[All input plugins will lose data](https://github.com/aws/aws-for-fluent-bit/blob/mainline/troubleshooting/debugging.md#overlimit-warnings), except for [tail](https://docs.fluentbit.io/manual/pipeline/inputs/tail), which can save its file offset before it is paused. *New logs will not be ingested, older logs in the buffer will remain.* And other components may experience errors, for example, if you use the [log4j TCP appender to write logs to Fluent Bit](https://github.com/aws-samples/amazon-ecs-firelens-examples/tree/mainline/examples/fluent-bit/ecs-log-collection), log4j will experience a [connection failure](https://github.com/aws/aws-for-fluent-bit/blob/mainline/troubleshooting/debugging.md#log4j-tcp-appender-write-failure). 

Unfortunately, in FireLens, the [log inputs for stdout/stderr error logs from your application container(s) is auto-generated by ECS](https://aws.amazon.com/blogs/containers/under-the-hood-firelens-for-amazon-ecs-tasks/), and has no `Mem_Buf_Limit` set. Thus, you must use the workaround in this blog to set it: [How to set Fluentd and Fluent Bit input parameters in FireLens](https://aws.amazon.com/blogs/containers/how-to-set-fluentd-and-fluent-bit-input-parameters-in-firelens/). See [Full FireLens Configuration Examples](#full-firelens-configuration-examples) for an example that includes setting limits and overriding the stdout/stderr input generated by FireLens.

Here is a full configuration example with only memory buffering:

```
[INPUT]
    Name                              tcp
    Listen                            0.0.0.0
    Port                              5170
    Buffer_Size                       64
    Format                            none
    Tag                               tcp-logs
    Mem_Buf_Limit                     50M
```

### Case 2: Filesystem and Memory Buffering: "storage.type filesystem"

*Full example config file:* [**filesystem.conf**](filesystem.conf)

You can also tell the input to buffer on the filesystem. *If you choose `storage.type filesystem` then the `Mem_Buf_Limit` parameter explained in the previous section no longer limits memory. Instead, as explained below, use `storage.max_chunks_up`*.

By default, when you enable filesystem buffering, Fluent Bit buffers all chunks BOTH in memory and on the filesystem. The number of chunks that are "up" in memory can be restricted with the `storage.max_chunks_up` setting. When the maximum number of "up" chunks in memory is reached, Fluent Bit will start keeping some chunks only on the filesystem or "down". A chunk stores ingested log records and is restricted by Fluent Bit to be at most 2 MB. Thus, with the default `storage.max_chunks_up` of 128 Fluent Bit is limited to roughly 256 MB of memory per input. Please note that in contrast to the limit setting for `storage.type memory`- `Mem_Buf_Limit` is defined per input and applies per input. Instead, `storage.max_chunks_up` is a global `[SERVICE]` section level config that then applies to each input individually- each input gets its own "up" chunk limit. So with the default value, each input gets 128 "up" chunks.

#### storage.type filesystem: What happens when the memory limit is reached?

Fluent Bit can not append/ingest to a new chunk unless it is "up" in memory. Therefore, it is possible to have slightly more than `storage.max_chunks_up` up in memory if Fluent Bit is actively ingesting records. This can be controlled by the input level setting `storage.pause_on_chunks_overlimit`. This per-input setting is disabled ("off") by default and if turned "on" then the input will be paused when the number of "up" chunks for that input goes over the `max_chunks_up` limit. This means that the input will stop ingesting new data until more chunks have been sent and there is "room" add new data. With `storage.pause_on_chunks_overlimit` you will ensure memory usage stays below the `max_chunks_up` limit. However, with many inputs, you will also lose logs when the input is paused. With the tail input plugin, this is less of a concern, since pausing just means halting at the current file offset- if the file remains on disk, tail can resume reading it later. However, with other inputs, like the TCP input or the forward input FireLens uses to collect stdout & stderr logs, pausing the input will lead to immediate log loss. 

If an input has been paused due to the limit, then you will see a log message from Fluent Bit like:

```
[input] forward.1 paused (storage buf overlimit)
```

In summary:
* `storage.pause_on_chunks_overlimit On`: When the `max_chunks_up` memory limit is reached, the input will stop ingesting records. *New logs will not be ingested and older logs will remain buffered*.
* `storage.pause_on_chunks_overlimit Off`: When the `max_chunks_up` memory limit is reached, Fluent Bit will work to move chunks to be "down" or only on the filesystem. The input will continue to ingest new data. 

#### storage.type filesystem: Estimating Total Memory Usage of the Fluent Bit Process

The estimate in the previous section of max memory usage for Fluent Bit can be used here too. While `storage.max_chunks_up` only limits input buffer memory, the filters and outputs that process each input chunk should at most only double the total amount of memory used. (You can picture this with a simple scenario, an input ingests one 2 MB chunk, and then while that is sitting in memory, an output tries to send it and thus has to create a request buffer which copies the chunk into a new 2 MB memory block).

To be clear, in Fluent Bit there are the following types of memory:
1. **Input Plugin Buffer Memory**: Memory used by inputs to buffer incoming logs. *This can be directly limited and controlled with `Mem_Buf_Limit`*. 
2. **Filter Plugin Processing Memory**: Memory used by filter plugins to process your log records. When filters make modifications to your records, they must temporarily copy them and then alter the copy. *This memory can not be directly controlled*, however, controlling the input memory indirectly controls the amount of logs flowing through filters. 
3. **Output Plugin Processing Memory**: Memory used by output plugins to process and send your log records. When outputs send your log records, they must create request/send buffers. *This memory can not be directly controlled*, however, controlling the input memory indirectly controls the amount of logs flowing through outputs. 
4. **Internal Fluent Bit Engine Memory**: Memory used by the Fluent Bit engine and scheduler to manage all of the plugins. This memory usage is generally small, and does generally increase as log throughput increases. *This memory usage can not be directly controlled.*

AWS FireLens team has found that in general, the total memory usage of Fluent Bit will stay below two times the sum of the `Mem_Buf_Limit` for each input. This is a rough estimate for sum of all types of memory usage noted above:

```
Total Max Memory Usage <= 2 * # of input definitions * storage.max_chunks_up * 2 MB per chunk
```

For example, if you have the default 128 value for `max_chunks_up` and then two input definitions, even under high load Fluent Bit should not consume much more than:

```
2 * 2 inputs * 128 chunks per input * 2 MB per chunk = 1024 MB memory estimate for high load conditions
```

*Please remember that this is an estimation method and not a guarantee*.

Please read the [official Fluent Bit documentation](https://docs.fluentbit.io/manual/administration/buffering-and-storage) to understand the filesystem buffering settings.

Unfortunately, as noted in the documentation link above, the filesystem buffering must be enabled on the input definition, therefore, for the [stdout/stderr input generated by FireLens](https://aws.amazon.com/blogs/containers/under-the-hood-firelens-for-amazon-ecs-tasks/), the same workaround is required: [How to set Fluentd and Fluent Bit input parameters in FireLens](https://aws.amazon.com/blogs/containers/how-to-set-fluentd-and-fluent-bit-input-parameters-in-firelens/). See [Full FireLens Configuration Examples](#full-firelens-configuration-examples) for an example that includes setting limits and overriding the stdout/stderr input generated by FireLens.

Here is a full configuration example with disk buffering:

```
[SERVICE]
    flush                     1
    log_Level                 info
    storage.path              /var/log/flb-storage/
    storage.sync              normal
    storage.checksum          off
    storage.max_chunks_up     128
    storage.backlog.mem_limit 5M

[INPUT]
    Name                              tcp
    Listen                            0.0.0.0
    Port                              5170
    Buffer_Size                       64
    Format                            none
    Tag                               tcp-logs
    storage.type                      filesystem
    storage.pause_on_chunks_overlimit On

[OUTPUT]
    name                      cloudwatch_logs
    match                     *
    log_group_name            buffer-example
    log_stream_name           this-is-an-example
    region                    us-east-1
    storage.total_limit_size  1G
```

#### Controlling Disk Usage

If you choose to follow the above recommendations and [enable disk buffering](https://docs.fluentbit.io/manual/administration/buffering-and-storage), you may want to limit the total amount of disk space that Fluent Bit can use. This can be configured with the output configuration `storage.total_limit_size`:

```
[OUTPUT]
    name                      cloudwatch_logs
    match                     *
    log_group_name            buffer-example
    log_stream_name           this-is-an-example
    region                    us-east-1
    storage.total_limit_size  1G
```

A couple of key notes:
1. This setting only works and comes into affect if the input(s) sending logs to the output have `storage.type filesystem`.
2. The setting is per output, to fully control disk usage you need to set this for all outputs whose inputs have `storage.type filesystem`. Each input must send each chunk of log records and has its own retries. This is why the setting is per output- consider the case where one input sends to multiple outputs. One output could be done with a specific chunk after a successful send, while another output could be waiting on a retry for the same chunk/

#### What happens when the disk limit is reached?

When the `storage.total_limit_size` is reached, the oldest chunks will be deleted first to make room for new chunks. *Older buffered logs will be lost in favor of making room for newer logs.*

In AWS for Fluent Bit [2.31.9](https://github.com/aws/aws-for-fluent-bit/releases/tag/v2.31.9) there is an info level message logged when this limit is reached and chunks will be removed:

```
[info] [input chunk] remove chunk 1-1678829080.465983679.flb with size 2000000 bytes from input plugin forward.1 to output plugin cloudwatch_logs.0 to place the incoming data with size 1000000 bytes, total_limit_size=200000000
```

Also, there is an error message when new chunks can not be ingested at all due to the limit. If you set a limit that is not very small this should not happen in practice as Fluent Bit should always be able to delete older chunks to ingest new data.

```
[error] [input chunk] no available chunk
[error] [input chunk] chunk 1-1678829080.465983679.flb would exceed total limit size in plugin cloudwatch_logs.0
```

### Full FireLens Configuration Examples

* [Full Filesystem Storage Example](firelens-full-filesystem-example/)
* [Full Memory Buffer Example](firelens-full-memory-example/)

These examples show how to control the memory usage for all inputs in FireLens, including the FireLens stdout/stderr unix socket forward input. These examples use the technique in [How to set Fluentd and Fluent Bit input parameters in FireLens](https://aws.amazon.com/blogs/containers/how-to-set-fluentd-and-fluent-bit-input-parameters-in-firelens/). The entrypoint for the Fluent Bit container is customized to use an alternate configuration file path and ignore the FireLens generated input. 

Please note, this technique ignores & overrides the FireLens generated configuration and thus has the following downsides:
* All containers that use Fluent Bit for stdout and stderr logging should specify the `awsfirelens` log driver without any options. 
* All configuration must be present in your custom config file; the `logConfiguration` options map should be empty as its content goes to the FireLens generated configuration.
* The `enable-ecs-log-metadata` firelens setting (default true) should be set to `false` because the metadata goes to the FireLens generated configuration file. Your logs will not have metadata attached: [What will the logs collected by FireLens look like?](https://github.com/aws/aws-for-fluent-bit/blob/mainline/troubleshooting/debugging.md#what-will-the-logs-collected-by-fluent-bit-look-like). To add metadata to logs, consider using the [init tag with ECS metadata support in env vars](https://github.com/aws-samples/amazon-ecs-firelens-examples/tree/mainline/examples/fluent-bit/init-metadata).

Please review the full configuration examples carefully. To understand the technique used, please review [How to set Fluentd and Fluent Bit input parameters in FireLens](https://aws.amazon.com/blogs/containers/how-to-set-fluentd-and-fluent-bit-input-parameters-in-firelens/).

The following settings are of note and are **required** for these examples to work properly:
* The Dockerfile for Fluent Bit adds our custom config file and overrides the entrypoint to use it directly. This means we will not use the FireLens generated config file placed at `/fluent-bit/etc/fluent-bit.conf`.
```
FROM public.ecr.aws/aws-observability/aws-for-fluent-bit:stable
ADD extra.conf /fluent-bit/alt/fluent-bit.conf
CMD echo -n "AWS for Fluent Bit Container Image Version " && \
    cat /AWS_FOR_FLUENT_BIT_VERSION && echo "" && \
    exec /fluent-bit/bin/fluent-bit -e /fluent-bit/firehose.so -e /fluent-bit/cloudwatch.so -e /fluent-bit/kinesis.so -c /fluent-bit/alt/fluent-bit.conf
```
* The app container specifies `awsfirelens` with no options. FireLens uses the options to generate an output definition that goes in the generated config file at `/fluent-bit/etc/fluent-bit.conf`. Since we ignore that file, we do not set any options; please place all output configurations in your `/fluent-bit/alt/fluent-bit.conf` file. Simply specifying the `awsfirelens` log driver tells ECS to configure the container to send logs to the Fluentd log driver and send them to Fluent Bit over the `/var/run/fluent.sock` unix socket that we specified in our `/fluent-bit/alt/fluent-bit.conf` file. See our debugging guide for information on the [FireLens tag pattern used for container logs](https://github.com/aws/aws-for-fluent-bit/blob/mainline/troubleshooting/debugging.md#firelens-tag-and-match-pattern-and-generated-config). 
```
"logConfiguration": {
	"logDriver":"awsfirelens"
}
```
* The FireLens configuration object does not specify a config file to import. This is necessary because any config file specified would be imported (using an `@INCLUDE` statement) into the FireLens generated config file, which we ignore/override. Additionally, we explicitly set `"enable-ecs-log-metadata":"false"` since the metadata `record_modifier` filter would be added to generated config file, which we ignore/override. 
```
"firelensConfiguration": {
	"type": "fluentbit",
	"options":{
		"enable-ecs-log-metadata":"false"
	}
}
```
* The custom Fluent Bit image must be built and used in the task definition:
```
"image": "XXXXXXXXXXXX.dkr.ecr.us-east-1.amazonaws.com/customized-flb-with-entrypoint-overridden:latest"
```

To use the examples, simply:
1. Customize `extra.conf` with your desired configuration. Notice the `# TODO` comment lines.
2. Build the `fluent-bit-image` directory with `docker build`. 
3. Push the Fluent Bit image to Amazon ECR for use with your task.
4. Customize and register the `task-definition.json`. Remember to add your custom image.

## 3. Monitor Fluent Bit

You can monitoring Fluent Bit's plugin metrics to track the number of retries and also track the storage usage. 
We have relevant FireLens examples: 
- [Send Fluent Bit internal metrics to CloudWatch](https://github.com/aws-samples/amazon-ecs-firelens-examples/tree/mainline/examples/fluent-bit/send-fb-internal-metrics-to-cw).
- [CPU, Disk, and Memory Usage Monitoring with ADOT](https://github.com/aws-samples/amazon-ecs-firelens-examples/tree/mainline/examples/fluent-bit/adot-resource-monitoring)