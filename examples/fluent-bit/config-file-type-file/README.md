### FireLens Example: Using the 'file' config type

This is the same config example as the "Adding custom keys to log events" example. In this though, we specifically look at how to use the 'file' 'config-file-type' in FireLens.

Let's go through it step by step.

1. Add your own custom config to extra.conf. Note that the name of that file is not special; you could name it anything.
2. Build a custom Fluent Bit image using the provided Docker file (which simply copies your extra config file into the AWS for Fluent Bit image). You can place this extra config file anywhere in the Docker image *except* `/fluent-bit/etc/fluent-bit.conf`. That config file path is the path used by FireLens.  Push this custom image to Amazon ECR.
3. Reference your new custom Fluent Bit image ECR repo in the container definition for the FireLens container.
4. Reference the config file path in the FireLens configuration:
```
"firelensConfiguration": {
    "type": "fluentbit",
    "options": {
        "config-file-type": "file",
        "config-file-value": "/extra.conf"
    }
}
```
