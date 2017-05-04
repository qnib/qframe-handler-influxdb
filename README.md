# qframe-handler-influxdb
Influxdb handler for qframe ETL-framework 


```bash
$ docker run -ti --name qframe-filter-docker-stats --rm -e SKIP_ENTRYPOINTS=1 \
             -v ${GOPATH}/src/github.com/qnib/qframe-collector-docker-events:/usr/local/src/github.com/qnib/qframe-collector-docker-events \
             -v ${GOPATH}/src/github.com/qnib/qframe-collector-docker-stats:/usr/local/src/github.com/qnib/qframe-collector-docker-stats \
             -v ${GOPATH}/src/github.com/qnib/qframe-filter-docker-stats:/usr/local/src/github.com/qnib/qframe-filter-docker-stats \
             -v ${GOPATH}/src/github.com/qnib/qframe-handler-influxdb:/usr/local/src/github.com/qnib/qframe-handler-influxdb \
             -v ${GOPATH}/src/github.com/qnib/qframe-types:/usr/local/src/github.com/qnib/qframe-types \
             -v ${GOPATH}/src/github.com/qnib/qframe-utils:/usr/local/src/github.com/qnib/qframe-utils \
             -v /var/run/docker.sock:/var/run/docker.sock \
             -w /usr/local/src/github.com/qnib/qframe-handler-influxdb \
              qnib/uplain-golang bash
$ govendor update github.com/qnib/qframe-collector-docker-events/lib \
                  github.com/qnib/qframe-collector-docker-stats/lib \
                  github.com/qnib/qframe-filter-docker-stats/lib \
                  github.com/qnib/qframe-handler-influxdb/lib \
                  github.com/qnib/qframe-types \
                  github.com/qnib/qframe-utils
$ govendor fetch +m
```

```bash
$ 

```