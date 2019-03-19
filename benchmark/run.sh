#!/bin/bash 

NCPUS=2
PARALLEL_REQ=${PARALLEL_REQ:-200}
TOTAL_REQ=${TOTAL_REQ:-20000}

echo ">>>> Execute Task"
ghz -c $PARALLEL_REQ -n $TOTAL_REQ -insecure -cpus $NCPUS \
	-proto ./protobuf/coreapi/api.proto \
	-call api.Core.ExecuteTask \
	-d '{"serviceID": "benchmark-service", "taskKey": "foo", "inputData": "{}" }' \
	localhost:50052 


echo ">>>> Emit Event"
ghz -c $PARALLEL_REQ -n $TOTAL_REQ -insecure -cpus $NCPUS \
  -proto ./protobuf/serviceapi/api.proto \
  -call api.Service.EmitEvent \
  -d '{"token": "benchmark-service", "eventKey": "foo", "eventData": "{}" }' \
  localhost:50052

