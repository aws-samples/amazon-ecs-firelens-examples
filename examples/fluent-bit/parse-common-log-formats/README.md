### FireLens Example: Parse common log formats

Often, your application will output logs in a common format (ex: apache, nginx), and you will want to turn them into JSON.

For example, an nginx log line like the following:
```
172.17.0.1 - - [03/Oct/2019:00:06:20 +0000] "GET / HTTP/1.1" 200 612 "-" "curl/7.54.0" "-"
```

Could look like the following when it is sent to your log destination:
```
{
  "remote": "172.17.0.1",
  "host": "-",
  "user": "-",
  "method": "GET",
  "path": "/",
  "code": "200",
  "size": "612",
  "referer": "-",
  "agent": "curl/7.54.0"
}
```

In order to parse logs, you will need a [parsers file](https://docs.fluentbit.io/manual/parser), which is specified in the [service section](https://docs.fluentbit.io/manual/service):
```
[SERVICE]
    Parsers_File /fluent-bit/parsers/parsers.conf
```

The `amazon/aws-for-fluent-bit` image contains a number of parsers files under `/fluent-bit/parsers/`. These parsers are copied directly from the official Fluent Bit Docker image. See their [repository for the parsers](https://github.com/fluent/fluent-bit-docker-image/tree/1.3).

You may need to create your own custom parser. To do this, you will need to create a custom image with the file.

Example Dockerfile:
```
FROM amazon/aws-for-fluent-bit:latest
ADD custom_parser.conf /fluent-bit/parsers/custom_parser.conf
```

Service Section in external Fluent Bit config file:
```
[SERVICE]
    Parsers_File /fluent-bit/parsers/custom_parser.conf
```
