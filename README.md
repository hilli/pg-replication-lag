# pg-replication-lag

### Getting started
Copy `pg-replication-lag.yaml.sample` to `pg-replication-lag.yaml`. Fill out `pg-replication-lag.yaml` with some sane values. Run `./pg-replication-lag` or use on of the flags as well:

```
$ ./pg-replication-lag --help
Usage of ./pg-replication-lag:
  -config string
    	Path to YAML config file (default "./pg-replication-lag.yaml")
  -verbose
    	Add verbosity to the output
```

It will look in the _current_ directory for a config file. If thats not where you keep it, specify it with `-config` flag.