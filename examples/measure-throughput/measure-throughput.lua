-- Fluent Bit Lua script to measure throughput and output in AWS CloudWatch EMF format
-- Enhanced version with TaskArn dimension support

-- Configuration from environment variables
local LOG_GROUP = os.getenv("LOG_GROUP") or "/aws/fluent-bit/analytics"
local NAMESPACE = os.getenv("METRIC_NAMESPACE") or "aws-for-fluent-bit/LogThroughput"

-- Variables to track throughput calculation
local last_throughput_time = 0
local total_bytes = 0
local last_total_bytes = 0
local task_arn = nil

-- Function to retrieve TaskArn from record field
function get_task_arn_from_record(record)
    -- Return cached value if already retrieved
    if task_arn then
        return task_arn
    end

    -- Extract TaskArn from record field ecs_task_arn
    if record and record.ecs_task_arn then
        task_arn = record.ecs_task_arn
        return task_arn
    end

    return "UNKNOWN"
end

function measure_throughput(tag, timestamp, record)
    -- Calculate the size of the current record
    local record_str = ""
    for k, v in pairs(record) do
        record_str = record_str .. tostring(k) .. ":" .. tostring(v) .. " "
    end
    local record_size = string.len(record_str)

    -- Update total bytes processed
    total_bytes = total_bytes + record_size

    -- Get current time
    local now = os.time()

    -- Initialize on first call
    if last_throughput_time == 0 then
        last_throughput_time = now
        last_total_bytes = total_bytes
        -- Pass through the original record
        return 1, timestamp, record
    end

    -- Check if 10 seconds have passed
    if now - last_throughput_time >= 10 then
        -- Calculate throughput
        local bytes_processed = total_bytes - last_total_bytes
        local time_elapsed = now - last_throughput_time
        local throughput_mbps = (bytes_processed / time_elapsed) / (1000 * 1000)

        -- Get current TaskArn from record
        local current_task_arn = get_task_arn_from_record(record)

        -- Create EMF formatted metric record as JSON string
        local emf_json = string.format('{"_aws":{"Timestamp":%d,"CloudWatchMetrics":[{"Namespace":"%s","Dimensions":[["LogGroup","TaskArn"]],"Metrics":[{"Name":"ThroughputMbps","Unit":"Megabytes/Second"}]}]},"ThroughputMbps":%f,"LogGroup":"%s","TaskArn":"%s"}',
            now * 1000,  -- EMF expects milliseconds
            NAMESPACE,
            throughput_mbps,
            LOG_GROUP,
            current_task_arn
        )

        -- Print EMF JSON to stdout
        print(emf_json)
        io.flush()

        -- Update tracking variables
        last_throughput_time = now
        last_total_bytes = total_bytes
    end

    -- Pass through the original record
    return 1, timestamp, record
end
