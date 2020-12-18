### FireLens Example: Storing Configs in EFS

This example is similar to the 'file' 'config-file-type' in FireLens example. However, in this example, instead of creating a custom Fluent Bit image with the config file, we create an EFS drive with the configuration and then mount it into all of our tasks. One option is to store multiple configs in a single EFS, then mount that into all tasks, and use the FireLens configuration to select the config file for each task.

Let's go through it step by step.

1. Create the EFS drive and add your custom configuration files to it: https://docs.aws.amazon.com/AmazonECS/latest/developerguide/tutorial-efs-volumes.html
2. Mount the EFS drive in your task definition; here is an example:
```
"volumes": [
            {
                "name": "firelens-conf",
                "efsVolumeConfiguration": {
                    "fileSystemId": "fs-c3c9bcbb",
                    "rootDirectory": "/",
                    "transitEncryption": "DISABLED",
                    "authorizationConfig": {
                        "iam": "DISABLED"
                    }
                }
            }
        ],
```
3. Next, mount the drive into your Fluent Bit container:
```
"mountPoints": [
                    {
                        "sourceVolume": "firelens-conf",
                        "containerPath": "/configs"
                    }
                ],
```
3. Finally, reference the configuration file from the EFS drive in FireLens:
```
"firelensConfiguration": {
                    "type": "fluentbit",
                    "options": {
                        "config-file-type": "file",
                        "config-file-value": "/configs/extra.conf"
                    }
                }
```


See the task definition file in this folder for a full example.
