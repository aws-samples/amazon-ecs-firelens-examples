### FireLens Example: Logging to SolarWinds Loggly with Fluent Bit

You can use Fluent Bit's [HTTP output plugin](https://docs.fluentbit.io/manual/pipeline/outputs/http) to send your FireLens log events to SolarWinds Loggly using the [HTTP/S Bulk Endpoint](https://documentation.solarwinds.com/en/Success_Center/loggly/Content/admin/http-bulk-endpoint.htm). The [Loggly HTTP/S Bulk Endpoint](https://documentation.solarwinds.com/en/Success_Center/loggly/Content/admin/http-bulk-endpoint.htm) is available at `/bulk/TOKEN/` and can be used with [tags](https://documentation.solarwinds.com/en/Success_Center/loggly/Content/admin/tags.htm) with `/bulk/TOKEN/tag/TAG/`.

AWS recommends that you store sensitive information (like the URI containing your Loggly [customer token](https://documentation.solarwinds.com/en/Success_Center/loggly/Content/admin/customer-token-authentication-token.htm?cshid=loggly_customer-token-authentication-token)) using [secretOptions](https://docs.aws.amazon.com/AmazonECS/latest/APIReference/API_Secret.html), as shown in the example [task definition](task-definition.json). This is optional - it is also valid to specify the token as part of the URI in the options map:

```
"logConfiguration": {
        "logDriver": "awsfirelens",
        "options": {
                "Name": "http",
                "Host": "logs-01.loggly.com",
                "URI": "/bulk/<token>/tag/<tag>/",
                "Port": "443",
                "tls": "on",
                "Format": "json_lines",
                "Json_date_key": "timestamp",
                "Json_date_format": "iso8601",
                "Retry_limit": "False"
        }
}
```
