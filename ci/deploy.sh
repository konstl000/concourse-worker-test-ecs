#!/bin/bash
: ${TEAM:="fme"}
fly -t ${TEAM} sp -p worker-test -c pipeline.yml
fly -t ${TEAM} unpause-pipeline -p worker-test
