# Reducing the Likelihood of Out of Memory Exceptions in FireLens

In normal/happy case conditions, the memory usage of Fluent Bit will often stay below 100 MB ([see the 'Resource Usage' section in this link](https://aws.amazon.com/blogs/containers/under-the-hood-firelens-for-amazon-ecs-tasks/)). 

But if there is a spike in log output from your application that Fluent Bit cannot handle or your log destination goes down, then Fluent Bit will have to buffer lots of logs. By default, Fluent Bit buffers logs in memory, and thus, in these scenarios, the Fluent Bit container's memory usage may increase significantly. In these conditions, Fluent Bit's memory usage can often spike to up to a Gigabyte or even more. 

The following are actions that we recommend considering. 

1. [Retry_Limit](#1-retry_limit)
2. [Control Fluent Bit buffer memory](#2-control-fluent-bit-buffer-memory)
    - [Case 1: Memory Buffering Only: default or "storage.type memory"](#case-1-memory-buffering-only-default-or-storagetype-memory)
        - [Estimating Total Memory Usage of the Fluent Bit Process](#storagetype-memory-estimating-total-memory-usage-of-the-fluent-bit-process)
        - [What happens when the memory limit is reached?](#storagetype-memory-what-happens-when-the-memory-limit-is-reached)
    - [Case 2: Filesystem and Memory Buffering: "storage.type filesystem"](#case-2-filesystem-and-memory-buffering-storagetype-filesystem)
        - [What happens when the memory limit is reached?](#storagetype-filesystem-what-happens-when-the-memory-limit-is-reached)
        - [Estimating Total Memory Usage of the Fluent Bit Process](#storagetype-filesystem-estimating-total-memory-usage-of-the-fluent-bit-process)
        - [Controlling Disk Usage](#controlling-disk-usage)
        - [What happens when the disk limit is reached?](#what-happens-when-the-disk-limit-is-reached)
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

Of course, there are downsides to `Mem_Buf_Limit`. If an input runs out of memory buffer, it stop ingesting logs and emit messages like the following:

```
[input] forward.1 paused (mem buf overlimit)
```

This means you will lose logs. *New logs will not be ingested, older logs in the buffer will remain.* And other components may experience errors, for example, if you use the [log4j TCP appender to write logs to Fluent Bit](https://github.com/aws-samples/amazon-ecs-firelens-examples/tree/mainline/examples/fluent-bit/ecs-log-collection), log4j will experience a connection failure. 

Unfortunately, in FireLens, the [log inputs for stdout/stderr error logs from your application container(s) is auto-generated by ECS](https://aws.amazon.com/blogs/containers/under-the-hood-firelens-for-amazon-ecs-tasks/), and has no `Mem_Buf_Limit` set. Thus, you must use the workaround in this blog to set it: [How to set Fluentd and Fluent Bit input parameters in FireLens](https://aws.amazon.com/blogs/containers/how-to-set-fluentd-and-fluent-bit-input-parameters-in-firelens/).

Here is a full configuration example with only memory buffering:

```
[INPUT]
    Name                              tcp
    Listen                            0.0.0.0
    Port                              5170
    Chunk_Size                        32
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

Unfortunately, as noted in the documentation link above, the filesystem buffering must be enabled on the input definition, therefore, for the [stdout/stderr input generated by FireLens](https://aws.amazon.com/blogs/containers/under-the-hood-firelens-for-amazon-ecs-tasks/), the same workaround is required: [How to set Fluentd and Fluent Bit input parameters in FireLens](https://aws.amazon.com/blogs/containers/how-to-set-fluentd-and-fluent-bit-input-parameters-in-firelens/).

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
    Chunk_Size                        32
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


## 3. Monitor Fluent Bit

You can monitoring Fluent Bit's plugin metrics to track the number of retries and also track the storage usage. 
We have relevant FireLens examples: 
- [Send Fluent Bit internal metrics to CloudWatch](https://github.com/aws-samples/amazon-ecs-firelens-examples/tree/mainline/examples/fluent-bit/send-fb-internal-metrics-to-cw).
- [CPU, Disk, and Memory Usage Monitoring with ADOT](https://github.com/aws-samples/amazon-ecs-firelens-examples/tree/mainline/examples/fluent-bit/adot-resource-monitoring)