#!/bin/bash
set -x

govendor update github.com/qnib/qframe-collector-docker-events/lib \
                github.com/qnib/qframe-collector-docker-stats/lib \
                github.com/qnib/qframe-filter-docker-stats/lib \
                github.com/qnib/qframe-handler-influxdb/lib \
                github.com/qnib/qframe-types \
                github.com/qnib/qframe-utils
