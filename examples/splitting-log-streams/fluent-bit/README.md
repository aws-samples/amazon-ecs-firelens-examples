### Streams Processing Example

This example contains the Fluent Bit streams processing example configuration.

### Use with FireLens

The setup for Fluent Bit for this example includes multiple files, so you'll need to build a [custom Fluent Bit image](custom-fluent-bit-image) with them. Because at the moment FireLens only supports pulling a single config file from S3.

You can then use the [example task definition](task_definition.json).

### Use outside of FireLens.

The example files [here]((custom-fluent-bit-image)) will work, but you should make one change- place the `fluent-bit.conf` at the path `/fluent-bit/etc/fluent-bit.conf` in your custom image or in your kubernetes config map. This is the default config path used by most Fluent Bit images.

You will also need to add an input definition to `fluent-bit.conf` to ingest your logs. This will depend on the details of your individual setup.

You can choose to [turn routing off](https://fluentbit.io/announcements/v1.2.0/) for your input- this means that logs from your input will not be able to be sent to destinations. They can only be sent to the stream processor (which can re-tag them and then they can be sent to destinations). 

For example:

```
[INPUT]
    Name     tail
    Path     /path/to/my/*.log
    Tag      logs.*
    Routable false
```
