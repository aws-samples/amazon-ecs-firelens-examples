### FireLens Example: Handling multiline logs with Fluentd

In use cases where your containers do not log json, but plain text, it might happen that the log messages include line breaks (e.g., Stacktraces in Java). You may want to handle these kind of logs in a way that lines are merged until a "real" new log line begins.
For this task, you can use Firelens with Fluentd to join these logs before sending them off.

You need a Docker image with Fluentd and the [concat plugin](https://github.com/fluent-plugins-nursery/fluent-plugin-concat), as well as additional plugins you may require to further process your logs or send them off (see [Dockerfile](Dockerfile) as an example).
In this example, the additional config ([extra.conf](extra.conf)) needs to be stored in S3, as referenced in `config-file-value` in the task definition:

```
"firelensConfiguration": {
    "type": "fluentd",
    "options": {
        "config-file-type": "s3",
        "config-file-value": "arn:aws:s3:::yourbucket/yourdirectory/extra.conf"
    }
},
```

This Fluentd configuration will concatinate lines until it encounters a new line matching the provided regex in [extra.conf](extra.conf). The Fluentd concat filter uses Ruby regular expressions to match log lines. **It is recommended that you use the [Rubular](https://rubular.com/) website to test your regular expressions.**

The example here would merge the lines:

```
2019-04-06 17:48:01.049 some log message with Stack
                        Stacktrace...
2019-05-06 17:48:01.049 further log messages with stack traces
                        Some other stacktrace..
2019-06-06 17:48:01.049 log lines without stacktrace
```

to merged lines:

```
2019-04-06 17:48:01.049 some log message with Stack \n Stacktrace...
2019-05-06 17:48:01.049 further log messages with stack traces \n Some other stacktrace..
2019-06-06 17:48:01.049 log lines without stacktrace
```

Additional configuration for the plugin can be found in the documentation of the [concat plugin](https://github.com/fluent-plugins-nursery/fluent-plugin-concat).

When forwarding these logs to CloudWatch Logs, you can use (CloudWatch Logs Insights)[https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/AnalyzingLogData.html] with the following query to see the merged log statements:

```
fields @timestamp, log
| sort @timestamp desc
```

