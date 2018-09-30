# pg-replication-lag

A command to let you get the _timed_ lag between a PostgreSQL master/slave replica set. PostgreSQL has queries to let you get the bytes difference between the master and the replica - This will let you get the amout of time there is between the masters writes and the replica has save the data to disk.

## Installation

```
go get github.com/hilli/pg-replication-lag
```

Or checkout the source and do a `build` (for your current platform) or optionally `make build-linux` (to get a Linux binary):

```
git clone https://github.com/hilli/pg-replication-lag.git
make deps
make build
```

## Getting started
Copy `pg-replication-lag.yaml.sample` to `pg-replication-lag.yaml`. 

```
curl https://raw.githubusercontent.com/hilli/pg-replication-lag/master/pg-replication-lag.yaml.sample > pg-replication-lag.yaml
```

Fill out `pg-replication-lag.yaml` with some sane values. The values should be self explanatory but it is probably worth noting that, since it is a replica, the username and password is the same. `max_lag_before_exit` is the time to maximal wait before exiting - If your replication is broken, we won't wait forever.

 Run `./pg-replication-lag` or use on of the flags as well:

```
$ ./pg-replication-lag --help
Usage of ./pg-replication-lag:
  -config string
    	Path to YAML config file (default "./pg-replication-lag.yaml")
  -verbose
    	Add verbosity to the output
```

It will look in the _current_ directory for a config file. If thats not where you keep it, specify it with `-config` flag.

### If all works out nicely

```
 $ ./pg-replication-lag_linux_amd64
{ "postgresql-replication-lag": { "time": "23.346179ms", "bytes": "0" } }
```

## Error cases
### If replication fails

```
$ ./pg-replication-lag_linux_amd64
pg-replication-lag: 2018/09/30 16:10:01 BAILING, waited too long (1m0.089542009s, 26716448 bytes behind)
panic: BAILING, waited too long (1m0.089542009s, 26716448 bytes behind)
```

### If DB connection fails

```
$ ./pg-replication-lag_linux_amd64
pg-replication-lag: 2018/09/30 16:23:16 Database error: dial tcp 127.0.0.1:5432: connect: connection refused
panic: Database error: dial tcp 127.0.0.1:5432: connect: connection refused
```

exits with code 2 in all cases.

