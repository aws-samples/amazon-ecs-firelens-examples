### FireLens Example: Logging to S3 with Fluent Bit

For documentation on the Fluent Bit S3 plugin, see: [out_s3 in Fluent Bit upstream](https://docs.fluentbit.io/manual/pipeline/outputs/s3)

S3 Support was released in Fluent Bit 1.6.0 and AWS for Fluent Bit 2.8.0.

### Reliability

The S3 plugin is meant to have a unique, persistent disk location configured with the `store_dir` parameter, which means:
- Each instance of Fluent Bit should use a unique directory/volume to buffer data
- Data volumes should be persistent. Fluent Bit can be restarted and will recover unsent data in the local buffer directory

The `store_dir` is used for two purposes:
1. Storing chunks of data before uploading them. If you enable S3 put object, then the plugin will buffer the entire file on Disk before sending. By default the plugin uses multipart uploads, and will only buffer a single chunk of the upload on disk at any point in time.
2. Storing metadata about multipart uploads. Multipart uploads will not be visible in S3 until all parts/chunks are sent. Fluent Bit expects a persistent disk to store this data. If Fluent Bit is stopped unexpectedly, it can be restarted with the same disk and will complete any unfinished uploads.

If you set up unique host mount volumes, or docker volumes for each of your tasks, you can have persistent storage with Fluent Bit.

This task definition example shows a scenario where you do not have persistent disk storage set up. Fluent Bit is therefore configured with a small upload file size and a short upload timeout, so that very little data is buffered disk at any point in time.

Please see the section in the Fluent Bit S3 documentation on reliability.

If reliability is a significant concern, we recommend using Kinesis Data Firehose to as a reliable distributed buffer between Fluent Bit and S3.
