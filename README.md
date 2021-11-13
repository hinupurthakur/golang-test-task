# golang-test-task
golang test task for dstack.ai is a CLI client for creating Docker containers.

## Problem Statement:

 ```
 Write a program called golang-test-task that accepts the following arguments:
Arguments
1. A name of a Docker image
2. A bash command (to run inside this Docker image)
3. A name of an AWS CloudWatch group
4. A name of an AWS CloudWatch stream
5. A name of an AWS region
Example: golang-test-task --docker-image python --bash-command $'pip install pip -U && pip
install tqdm && python -c \"import time\ncounter = 0\nwhile True:\n\tprint(counter)\n\tcounter
= counter + 1\n\ttime.sleep(0.1)"' --cloudwatch-group golang-test-task-group-1 --cloudwatch-stream golang-test-task-group-2 --aws-access-key-id ... --aws-secret-access-key ... --aws-region ...
Functionality
The program should create a Docker container using the given Docker image name,
and the given bash command
The program should handle the output logs of the container and send them to the
given AWS CloudWatch group/stream using the given AWS credentials. If the
corresponding AWS CloudWatch group or stream does not exist, it should create it
using the given AWS credentials.
Other requirements
The program should behave properly regardless of how much or what kind of logs
the container outputs
The program should gracefully handle errors and interruption

Golang Test Task 2
The credentials above have permissions to create AWS CloudWatch groups/streams
and to write and read log events
Source code:
Please feel free to share the final source code as a GitHub repository, and invite the
peterschmidt85 user to review it.

```

Dependencies 
The dependencies for this package are the official AWS SDK for Go and official Docker SDK for Go.

Assumptions
* This module assumes that the Access Key and the Secret Keys used with have the permissions to create AWS CloudWatch groups/streams and to write and read log events 
* This module assumes that the stream logs that are already present on the cloudwatch log stream should not be put again.

To Run:

```
./golang-test-task --docker-image python --bash-command $'echo "Hello World"' --cloudwatch-group golang-test-task-group-1 --cloudwatch-stream golang-test-task-group-2 --aws-access-key-id ... --aws-secret-access-key ... --aws-region ap-southeast-1
```

ToDo:
* Improve the log Messages with proper json format.

