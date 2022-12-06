# Test repo to simulate ECS tasks on any host with docker running
## Build
```
./build.sh
```
## Use
```
./bin/runner -c ${CONTAINER_CONFIG_FILE} -h ${HOST_CONFIG_FILE} -n ${CONTAINER_NAME}
```
For instance,
```
./bin/runner -c container.json -h host.json -n ccw
```
To connect to your concourse, you would need to provide the valid concourse keys located in /opt/concourse/keys as well as the valid value of `CONCOURSE_TSA_HOST` in the `container.json` file.

## Concourse worker tests
Once the concourse keys are provided and the values of TSA host are set properly in `container.json` and `container2.json` (with concourse 7.8.3 for comparison), you can deploy the pipeline in `ci/pipeline.yml`. You would also need to provide the valid target as `${TEAM}` to use the `deploy.sh` script.
```
./bin/runner -c container.json -h host.json -n ccw
./bin/runner -c container2.json -h host.json -n ccw783
cd ci
./deploy.sh
```
You should be able to observe the following:
  - the job fails for the latest concourse container image
  - the job works for the concourse container image with the 7.8.3 tag



