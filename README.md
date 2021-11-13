# golang-test-task
golang test task for dstack.ai is a CLI client for creating Docker containers.

Dependencies 
The dependencies for this package are the official AWS SDK for Go and official Docker SDK for Go.

Assumptions
* This module assumes that the  

To Run:

```
golang-test-task --docker-image python --bash-command $'echo "Hello World"' --cloudwatch-group golang-test-task-group-1 --cloudwatch-stream golang-test-task-group-2 --aws-access-key-id ... --aws-secret-access-key ... --aws-region ap-southeast-1
```


