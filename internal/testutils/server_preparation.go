package testutils

import (
	"fmt"
	influx "github.com/influxdata/influxdb/client/v2"
	"github.com/jackc/pgx"
	"github.com/timescale/outflux/internal/connections"
	"github.com/timescale/outflux/internal/schemamanagement/influx/influxqueries"
)

// PrepareServersForITest Creates a database with the same name on the default influx server and default timescale server
func PrepareServersForITest(db string) {
	CreateInfluxDB(db)
	CreateTimescaleDb(db)
}

// ClearServersAfterITest Deletes a database on both the default influx and timescale servers
func ClearServersAfterITest(db string) {
	DeleteInfluxDb(db)
	DeleteTimescaleDb(db)
}

// CreateInfluxDB creates a new influx database to the default influx server. Used for integration tests
func CreateInfluxDB(db string) {
	connService := connections.NewInfluxConnectionService()
	queryService := influxqueries.NewInfluxQueryService()
	connParams := &connections.InfluxConnectionParams{Server: InfluxHost}
	client, err := connService.NewConnection(connParams)
	panicOnErr(err)
	_, err = queryService.ExecuteQuery(client, db, "CREATE DATABASE "+db)
	panicOnErr(err)
}

// DeleteInfluxDb deletes a influx database on the default influx server. Used for integration tests
func DeleteInfluxDb(db string) {
	connService := connections.NewInfluxConnectionService()
	queryService := influxqueries.NewInfluxQueryService()
	connParams := &connections.InfluxConnectionParams{Server: InfluxHost}
	client, err := connService.NewConnection(connParams)
	panicOnErr(err)
	_, err = queryService.ExecuteQuery(client, db, "DROP DATABASE "+db)
	panicOnErr(err)
}

// CreateInfluxMeasure creates a measure with the specified name. For each point the tags and field values are given
// as maps
func CreateInfluxMeasure(db, measure string, tags []*map[string]string, values []*map[string]interface{}) {
	connService := connections.NewInfluxConnectionService()
	connParams := &connections.InfluxConnectionParams{Server: InfluxHost}
	client, err := connService.NewConnection(connParams)
	panicOnErr(err)

	bp, _ := influx.NewBatchPoints(influx.BatchPointsConfig{Database: db})

	for i, tagSet := range tags {
		point, _ := influx.NewPoint(
			measure,
			*tagSet,
			*values[i],
		)
		bp.AddPoint(point)

	}

	err = client.Write(bp)
	panicOnErr(err)
	client.Close()
}

// CreateTimescaleDb creates a new databas on the default server and then creates the extension on it
func CreateTimescaleDb(db string) {
	dbConn := OpenTSConn(defaultPgDb)
	defer dbConn.Close()
	_, err := dbConn.Exec("CREATE DATABASE " + db)
	panicOnErr(err)
}

// OpenTSConn opens a connection to a TimescaleDB
func OpenTSConn(db string) *pgx.Conn {
	connString := fmt.Sprintf(TsConnStringTemplate, db)
	connConfig, _ := pgx.ParseConnectionString(connString)
	c, _ := pgx.Connect(connConfig)
	return c
}

// DeleteTimescaleDb drops a databass on the default server
func DeleteTimescaleDb(db string) {
	dbConn := OpenTSConn(defaultPgDb)
	defer dbConn.Close()
	_, err := dbConn.Exec("DROP DATABASE " + db)
	panicOnErr(err)
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
