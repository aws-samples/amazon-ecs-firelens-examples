# ECS FireLens Log Deletion Example

This example will walk you through options to set up log file deletion/clean up in your ECS FireLens task.

This example is based upon the tail example in the [ecs-log-collection](https://github.com/aws-samples/amazon-ecs-firelens-examples/tree/mainline/examples/fluent-bit/ecs-log-collection) FireLens example. It adds log deletion.

## Scenario

For this example, consider that your app produces log files with paths like `/var/log/service-{timestamp}.log`. Each hour, your app write to a log file with a new timestamp. Thus, if your app runs for a very long time, you will quickly have lots of log files sitting around on disk. This is not ideal for two reasons:
1. Your disk may eventually fill up.
2. The large number of log files may slow down Fluent Bit. Fluent Bit continuously scans for files matched by your configured `Path` and runs a syscall on them to check if they were modified. If you have tons of log files that it must scan through, this will cause it to slow down, and in extreme cases, you may even lose logs. 

Fluent Bit will use this base configuration (based on the ECS Log Collection example) to tail the log files:

```
[INPUT]
    Name              tail
    Tag               service-log
    Path              /var/log/service*.log
    DB                /var/log/flb_service.db
    DB.locking        true
    Skip_Long_Lines   On
    Refresh_Interval  10
    Rotate_Wait       30
# Output for stdout logs
# Change the match if your container is not named 'app'
[OUTPUT]
    Name cloudwatch
    Match   app-firelens*
    region us-east-1
    log_group_name firelens-tutorial-$(ecs_cluster)
    log_stream_name /logs/app/$(ec2_instance_id)-$(ecs_task_id)
    auto_create_group true
    retry_limit 2
# Output for the file logs
[OUTPUT]
    Name cloudwatch
    Match   service-log
    region us-east-1
    log_group_name firelens-tutorial-$(ecs_cluster)
    log_stream_name /logs/service/$(ec2_instance_id)-$(ecs_task_id)
    auto_create_group true
    retry_limit 2
```

## Log Deletion Options

### Background

Make sure you are familiar with the [Fluent Bit documentation](https://docs.fluentbit.io/manual/) and with [AWS for Fluent Bit](https://github.com/aws/aws-for-fluent-bit) and its container image distribution.

### 1. Use Fluent Bit Exec Input to run deletion command

##### Warning: exec input issue in AWS for Fluent Bit <= 2.31.11

There is a [known issue](https://github.com/aws/aws-for-fluent-bit/issues/661#issuecomment-1569515241) in the exec input in AWS for Fluent Bit <= 2.31.11 which can cause it to crash, generally shortly after startup. 

This issue is resolved in [2.31.12](https://github.com/aws/aws-for-fluent-bit/releases).

##### Prerequisites:
* The volume mount for Fluent Bit must not be read-only
* Familiar with [exec input](https://docs.fluentbit.io/manual/pipeline/inputs/exec) of Fluent Bit

You can run a command that deletes log files by adding an exec input to the configuration. Here is the example:

```
[INPUT]
    Name          exec
    Tag           delete-old-files
    Command       find /var/log/  -type f -mmin +10080 -execdir  rm -f -- '{}' ';'
    Interval_Sec  600
    Interval_NSec 0
    Buf_Size      8mb
    Oneshot       false

[OUTPUT]
    Name          null
    Match         delete-old-files
```
Here is a breakdown of the individual components of the command:
* ```Interval_Sec```: This is the polling interval in seconds.
* ```Oneshot```: Only run once at startup.
* ```find /var/log/```: This is the ```find``` command, which searches for files in the directory ```/var/log/```.
* ```-type f```: This option tells find to search only for regular files, not directories or other types of files.
* ```-mmin +10080```: This option tells find to search for files that are older than 7 days (10080 minutes).
* ```-execdir rm -f -- '{}'';'```: This is the action to take on the files that match the search criteria. The ```-execdir``` option runs the ```rm``` command in the directory where each file is found. The ```-f``` option to ```rm``` tells it to force the deletion without prompting for confirmation. The ```--```is a separator that indicates the end of the command options and the start of the filenames, and ```{}``` is a placeholder that gets replaced with the name of each file found by find. Finally, the ```';'``` marks the end of the command that should be executed for each file.

### 2. Run deletion command with cron

##### Prerequisite:
* Familiar with [cron](https://en.wikipedia.org/wiki/Cron)
* Docker should be installed on your system

**1. Create a cron job file and write the command below inside the file**

* Create a file called "delete-log-files" in file ```/cron```
* Paste the following command into the file
```
0 0 * * * find /var/log/ -type f -mmin +10080 -execdir rm -f -- '{}' ';'
```
This is a cron job that runs every day at midnight (```0 0 * * *```). It uses the ```find``` command to locate files in the directory ```/var/log/``` that are older than 7 days (10080 minutes), and then deletes them using the ```rm -f``` command.
* ```0 0 * * *```: This is the cron schedule expression. It specifies that the job should run at 12:00 AM every day.
* Note: In the cron job file called "cron", an empty line is required at the end of the file for a valid cron file.

**2. Build image with cron**

* Install Cron
```
RUN yum install -y cronie which findutils perl-core compat-libcap1 lsof tar shadow-utils setuptool procps && yum clean all
```

* Copy crontab that deletes logs
```
COPY cron/delete-log-files /etc/cron.d/delete-log-files
```

* Give execution rights on the cron job
```
RUN chmod 0644 /etc/cron.d/delete-log-files
```

Note: You can view the dockerfile example [here](Dockerfile) and run command ```docker build``` to build the image.

### 3. Use log4j delete action

You can use log4j delete action via public RollingRandomAccessFile appender with the following guidances:
* [delete-logs-on-rollover](https://howtodoinjava.com/java/delete-logs-on-rollover/)
* [Apache log4j docs: Delete on Rollover](https://logging.apache.org/log4j/log4j-2.7/manual/appenders.html#CustomDeleteOnRollover)

Following is the example of deleting old log files older than 1 day:
```
<DefaultRolloverStrategy>
    <!--Delete old log files older than 1 day-->
    <Delete basePath="${log-file-path}" maxDepth="2">
        <IfLastModified age="1d">
        </IfLastModified>
    </Delete>
</DefaultRolloverStrategy>
```