### Fluentd Rewrite Tag Example

This folder contains the Fluentd rewrite tag filter example configuration.

### Pre-requisites

Build a Fluentd image with the [rewrite tag filter plugin](https://github.com/fluent/fluent-plugin-rewrite-tag-filter) and the [CloudWatch Logs plugin](https://github.com/fluent-plugins-nursery/fluent-plugin-cloudwatch-logs).

This repository is a good starting place for building Fluentd images; it's examples are frequently updated: https://github.com/fluent/fluentd-kubernetes-daemonset/tree/master/docker-image

### Use with FireLens

Upload the [config file](rewrite_tag.conf) to S3 and use the example [Task Definition](task_definition.json).

### Use outside of FireLens

Use the provided config file as a starting point. Edit it to add an input configuration for your logs; or `@include` it in an existing Fluentd config for your setup. 
