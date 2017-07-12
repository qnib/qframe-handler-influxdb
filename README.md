# qframe-handler-influxdb
Influxdb handler for qframe ETL-framework 

**Depreciated!** Moved to [qframe/handler-influxdb](https://github.com/qframe/handler-influxdb)

## Start DEV environment

The `docker-compose.yml` file expects to find some golang libs, so you better prepare.

```bash
$ go get github.com/qnib/qframe-collector-docker-events \
         github.com/qnib/qframe-collector-docker-stats \
         github.com/qnib/qframe-filter-docker-stats \
         github.com/qnib/qframe-handler-influxdb \
         github.com/qnib/qframe-types \
         github.com/qnib/qframe-utils
```

Fire up the services.

```bash
$ docker stack deploy -c docker-compose.yml qframe                                                                                                                                                  git:(master|✚1…
Creating service qframe_influxdb
Creating service qframe_frontend
Creating service qframe_qframe-dev
Creating service qframe_qframe
```

The `qframe-dev` task will run two scripts to update and fetch missing libraries, depending on your internet connection it might take a while.

```bash
$ docker logs -f $(docker ps -qlf label=com.docker.swarm.service.name=qframe_qframe-dev)                                                                                                                git:(master|✚2…
[II] qnib/init-plain script v0.4.19
> execute entrypoint '/opt/qnib/entry/00-govendor-update.sh'
+ govendor update github.com/qnib/qframe-collector-docker-events/lib github.com/qnib/qframe-collector-docker-stats/lib github.com/qnib/qframe-filter-docker-stats/lib github.com/qnib/qframe-handler-influxdb/lib github.com/qnib/qframe-types github.com/qnib/qframe-utils
> execute entrypoint '/opt/qnib/entry/10-govendor-fetch.sh'
+ govendor fetch +m
```
once it shows `> execute CMD 'wait.sh'` it is up and running.

```bash
$ docker exec -ti $(docker ps -qlf label=com.docker.swarm.service.name=qframe_qframe-dev) bash                                                                                                          git:(master|✚2…
root@383a9131f99b:/usr/local/src/github.com/qnib/qframe-handler-influxdb# go run main.go
2017/05/04 19:47:48 [II] Dispatch broadcast for Back, Data and Tick
2017/05/04 19:47:48 [  INFO] influxdb >> Start log handler influxdbv0.0.2
2017/05/04 19:47:48 [  INFO] influxdb >> Established connection to 'http://172.17.0.1:8086
2017/05/04 19:47:48 [  INFO] container-stats >> Start docker-stats filter v0.1.0
2017/05/04 19:47:48 [  INFO] container-stats >> [docker-stats]
2017/05/04 19:47:48 [  INFO] docker-events >> Start docker-events collector v0.2.1
2017/05/04 19:47:48 [  INFO] docker-events >> Connected to 'moby' / v'17.05.0-ce-rc1'
2017/05/04 19:47:49 [  INFO] docker-stats >> Connected to 'moby' / v'17.05.0-ce-rc1' (SWARM: active)
2017/05/04 19:47:49 [  INFO] docker-stats >> Currently running containers: 3
2017/05/04 19:47:49 [II] Start listener for: 'qframe_qframe.1.rpelk522jkavcez87qss4lckc' [383a9131f99bbe1872aa112d045b4d64a0990a8d588a3164e44c920ead39330c]
2017/05/04 19:47:49 [II] Start listener for: 'qframe_frontend.1.2xuwmen72eor9gids5iguuahc' [c128e1acfda824058feb2227368c2d612eb58ab2c4370b340f875cd6394d9138]
2017/05/04 19:47:49 [II] Start listener for: 'qframe_influxdb.1.u3lnhbnqrzukzzk7v1a27us87' [231e6c50eb683153afcd48bfeeb3111aef1b5b4d6d867d088d84de66bee3b301]
2017/05/04 19:47:51 [  INFO] influxdb >> Ticker: Write batch of 45
2017/05/04 19:47:54 [  INFO] influxdb >> Ticker: Write batch of 45
...
```

### Develop

To develop on one of the dependencies or the plugin itself, just change the code and update govendor.

```bash
$ 
```

# Run prepackaged container

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
