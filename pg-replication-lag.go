package main

/**
# pg-replication-lag:

Reports on the lag between 2 postgresql database servers, set up for master/slave replication.

*/

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v2"

	_ "github.com/lib/pq"
)

type conf struct {
	DbHostMaster      string        `yaml:"db_host_master"` // Database Master
	DbHostReplication string        `yaml:"db_host_replica"`
	DbPortMaster      string        `yaml:"db_port_master"`
	DbPortReplica     string        `yaml:"db_port_replica"`
	DbUser            string        `yaml:"db_user"`
	DbPassword        string        `yaml:"db_password"`
	DbName            string        `yaml:"db_name"`
	MaxLagBeforeExit  time.Duration `yaml:"max_lag_before_exit"`
}

const (
	XlogCurrentLocation       = "SELECT pg_current_xlog_location();"
	XlogReplicaReplayLocation = "SELECT pg_last_xlog_replay_location();"
	XlogDiff                  = "SELECT pg_xlog_location_diff('%s','%s');"
)

var (
	logger        *log.Logger
	masterDb      *sql.DB
	replicationDb *sql.DB
	err           error
	verbose       bool
	configFile    string
	c             conf
)

func init() {
	logger = log.New(os.Stderr, "pg-replication-lag: ", log.LstdFlags)
	flag.BoolVar(&verbose, "verbose", false, "Add verbosity to the output")
	flag.StringVar(&configFile, "config", "./pg-replication-lag.yaml", "Path to YAML config file")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nDocumentation at https://github.com/hilli/pg-replication-lag\n")
	}
	flag.Parse()
	c.loadConfigFile(configFile)

}

func main() {
	var initialBytesBehind int64

	dbMasterInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		c.DbHostMaster, c.DbPortMaster, c.DbUser, c.DbPassword, c.DbName)
	dbReplicaInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		c.DbHostReplication, c.DbPortReplica, c.DbUser, c.DbPassword, c.DbName)

	masterDb, err = sql.Open("postgres", dbMasterInfo)
	checkErr(err)
	defer masterDb.Close()

	replicationDb, err := sql.Open("postgres", dbReplicaInfo)
	checkErr(err)
	defer replicationDb.Close()
	masterXlog := getXlogLocation(masterDb, XlogCurrentLocation)
	startTime := time.Now()

	done := false
	firstRun := true
	for ok := true; ok; ok = (done != true) {
		replicaXlog := getXlogLocation(replicationDb, XlogReplicaReplayLocation)
		xlogDiffBytes := getXlogDiff(replicationDb, masterXlog, replicaXlog)
		if firstRun {
			initialBytesBehind = xlogDiffBytes
			firstRun = false
		}
		if verbose {
			logger.Println(fmt.Sprintf("Running %v behind, missing %v bytes", time.Since(startTime), xlogDiffBytes))
		}
		if xlogDiffBytes > 0 {
			// Bail out if we have waited too long
			if time.Since(startTime) > (c.MaxLagBeforeExit * time.Second) {
				logger.Panicf("BAILING, waited too long (%s, %v bytes behind)", time.Since(startTime), xlogDiffBytes)
				done = true
			}
			// Sleep 100ms
			time.Sleep(100 * time.Millisecond)
		} else {
			// Report time
			fmt.Println(fmt.Sprintf("{ \"postgresql-replication-lag\": { \"time\": \"%s\", \"bytes\": \"%v\" } }", time.Now().Sub(startTime), initialBytesBehind))
			done = true
		}
	}
}

func getXlogDiff(db *sql.DB, start string, end string) int64 {
	xlogDiffQuery := fmt.Sprintf(XlogDiff, start, end)
	result := getXlogLocation(db, xlogDiffQuery)
	diff, _ := strconv.ParseInt(result, 10, 32)
	return diff
}

func getXlogLocation(db *sql.DB, query string) string {
	rows, err := db.Query(query)
	checkErr(err)
	var result string
	for rows.Next() {
		err = rows.Scan(&result)
	}
	return result
}

func checkErr(err error) {
	if err != nil {
		logger.Panicln(fmt.Sprintf("Database error: %s", err))
	}
}

func (c *conf) loadConfigFile(configfile string) *conf {
	if _, err := os.Stat(configFile); !os.IsNotExist(err) {
		// Config file exists, load it
		yamlFile, err := ioutil.ReadFile(configfile)
		if err != nil {
			log.Printf("yamlFile.Get err   #%v ", err)
		}
		err = yaml.Unmarshal(yamlFile, c)
		if err != nil {
			log.Fatalf("Unmarshal: %v", err)
		}
	} else {
		logger.Panic("No config file found, can not continue")
	}
	return c
}
