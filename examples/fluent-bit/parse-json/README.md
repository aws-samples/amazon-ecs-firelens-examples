### FireLens Example: Parsing container stdout logs that are serialized JSON

The external Fluent Bit config file in this example will parse any logs that are JSON.
For example, if the log emitted from your container looks like this:

```
"{\"requestID\": \"b5d716fca19a4252ad90e7b8ec7cc8d2\", \"requestInfo\": {\"ipAddress\": \"204.16.5.19\", \"path\": \"/activate\", \"user\": \"TheDoctor\"}}"
```

Then the log message sent to your destination will look like this:

```
{
  "requestID": "b5d716fca19a4252ad90e7b8ec7cc8d2",
  "requestInfo": {
    "ipAddress": "204.16.5.19",
    "path": "/activate",
    "user": "TheDoctor"
  }
}
```

Please note that the JSON parser will remove the ECS Log Metadata fields. See this issue on their repository for more information: https://github.com/fluent/fluent-bit/issues/1605
