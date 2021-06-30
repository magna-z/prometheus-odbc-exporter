package odbc

import (
	"database/sql"

	"github.com/rs/zerolog/log"

	//_ "github.com/alexbrainman/odbc"
	_ "github.com/ClickHouse/clickhouse-go"
)

func Init(dsn string) {
	//dsn, err := url.Parse("tcp://ch.local:9000?username=USER&password=PASSWORD")
	//if err != nil {
	//	log.Fatal().
	//		Err(err).
	//		Str("metrics.dsn", *metricsDSN).
	//		Msg("Parse DSN string error")
	//}

	//dsnQuery := dsn.Query()
	//dsnQueryEncoded := dsnQuery.Encode()
	//log.Info().Msg(dsnQueryEncoded)
	////dsnQueryPass := dsnQuery.Get("password")
	////
	////dsnQuery.Set("password", url.QueryEscape(dsnQueryPass))
	//
	//odbc.Init(dsn.String())

	//db, err := sql.Open("odbc", "DRIVER={EXASolution Driver};EXAHOST=ex1..5.local:8563;UID=USER;PWD=PASSWORD")
	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("DB connection error")
	}

	var now sql.NullTime
	err = db.QueryRow("select now();").Scan(&now)
	if err != nil {
		log.Fatal().Err(err).Msg("DB query row error")
	}

	log.Info().Msgf("Query return %q", now)

	if db.Close() != nil {
		log.Fatal().Err(err).Msg("DB connection close error")
	}
}
