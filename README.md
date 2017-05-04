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
$ docker run -ti -v /var/run/docker.sock:/var/run/docker.sock qnib/qframe-handler-influxdb
> execute CMD 'qframe-handler-influxdb'
2017/05/04 16:57:44 [II] Dispatch broadcast for Back, Data and Tick
2017/05/04 16:57:44 [  INFO] container-stats >> Start docker-stats filter v0.1.0
2017/05/04 16:57:44 [  INFO] container-stats >> [docker-stats]
2017/05/04 16:57:44 [  INFO] influxdb >> Start log handler influxdbv0.0.2
2017/05/04 16:57:44 [  INFO] influxdb >> Established connection to 'http://172.17.0.1:8086
2017/05/04 16:57:44 [  INFO] docker-events >> Start docker-events collector v0.2.1
2017/05/04 16:57:44 [  INFO] docker-stats >> Connected to 'moby' / v'17.05.0-ce-rc1' (SWARM: active)
2017/05/04 16:57:44 [  INFO] docker-stats >> Currently running containers: 3
2017/05/04 16:57:44 [II] Start listener for: 'vigilant_agnesi' [4ce0c01b7997711f8aa3653fc9b1ca655c04786e325425560ef450e092aae25e]
2017/05/04 16:57:44 [II] Start listener for: 'gcollect_influxdb.1.ues4sdm7vzmhtzopkc08qb8fc' [9da2ac9a0db27f30bd913176c2177e588f0eae84c7cd885200a7391c9e6e6d72]
2017/05/04 16:57:44 [II] Start listener for: 'gcollect_frontend.1.17u6xcak5rogijggdzebeouiu' [ea2ff849cbd72164fd965caa944b8e0f125afccf8605bd260fa08b65d10daf2d]
2017/05/04 16:57:45 [  INFO] docker-events >> Connected to 'moby' / v'17.05.0-ce-rc1'
2017/05/04 16:57:47 [  INFO] influxdb >> Ticker: Write batch of 30
2017/05/04 16:57:50 [  INFO] influxdb >> Ticker: Write batch of 45
...
```