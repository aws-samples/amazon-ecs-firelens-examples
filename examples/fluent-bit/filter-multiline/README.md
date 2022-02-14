### FireLens Example: Concatenate Multiline or Stack trace log messages

Multiline Filter is available on `aws-for-fluent-bit` >= v2.22.0. It helps to concatenate messages that originally belong to one context but were split across multiple records or log lines. Common examples are stack traces or applications that print logs in multiple lines. More information could be found in [Fluent Bit official doc](https://docs.fluentbit.io/manual/pipeline/filters/multiline-stacktrace).

With this filter, you are able to use Fluent Bit built-in parses with auto detection and multi format support on:
- go
- python
- ruby
- java

As those are built-in, you can directly specify them in a field called `multiline.parser` in `[FILTER]` section.

Let's go through an example shows how to use the multiline filter:

1. In order to parse logs, you can either create a parser file for custom stacktrace or choose the built-in parsers. A parser file `parsers_multiline.conf` could be like as below:
```
[MULTILINE_PARSER]
    name          multiline-regex-test
    type          regex
    flush_timeout 1000
    #
    # Regex rules for multiline parsing
    # ---------------------------------
    #
    # configuration hints:
    #
    #  - first state always has the name: start_state
    #  - every field in the rule must be inside double quotes
    #
    # rules |   state name  | regex pattern                  | next state
    # ------|---------------|--------------------------------------------
    rule      "start_state"   "/(Dec \d+ \d+\:\d+\:\d+)(.*)/"  "cont"
    rule      "cont"          "/^\s+at.*/"                     "cont"
```
2. Add your own custom config to `extra.conf`. In this config, you need to specify the above parser file in `[SERVICE]` section and have another `[FILTER]` section to add parsers. Here, You can also directly add a built-in parser like `go`.
3. Build a custom Fluent Bit image using the provided Docker file (which simply copies these two customized files into the AWS for Fluent Bit image) by `docker build .`. You can place this extra config file anywhere in the Docker image *except* `/fluent-bit/etc/fluent-bit.conf`. That config file path is the path used by FireLens. Push this custom image to Amazon ECR with `ecs-cli push`.
4. Reference the new custom Fluent Bit image ECR repo in the container definition for the FireLens container. In this example we have also included a sample application container, which works with the example configuration and parser. The logs produced by the app will be concatenated by Fluent Bit.
5. Reference the config file path in the FireLens configuration:
```
"firelensConfiguration": {
    "type": "fluentbit",
    "options": {
        "config-file-type": "file",
        "config-file-value": "/extra.conf"
    }
}
```
*Note*: 
- `permissions.json` shows the permissions we should include in our IAM to send logs to CloudWatch.
- the path to the parser file in `extra.conf` should be the absolute path in the image.
- If your logs go to files instead of stdout, it is recommended to use Multiline in [Tail](https://docs.fluentbit.io/manual/pipeline/inputs/tail#multiline-support) direclty instead of the filter.

