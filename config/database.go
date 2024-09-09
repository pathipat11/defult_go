package config

import (
	pdb "app/app/provider/database"
	"log"

	"github.com/uptrace/bun"
)

func Database() {
	// Connect to database
	pdb.Register(
		&db,
		&pdb.DBOption{
			Host:     confString("DB_HOST", "127.0.0.1"),
			Port:     confInt64("DB_PORT", int64(5432)),
			Database: confString("DB_DATABASE", "GRIT-HRM"),
			Username: confString("DB_USER", "postgres"),
			Password: confString("DB_PASSWORD", ""),
			TimeZone: confString("TZ", "Asia/Bangkok"),
			SSLMode:  confString("DB_SSLMODE", "require"),
		},
	)
	log.Println("database connected success")

	// Connect to database2
	// pdb.Register(
	// 	&db2,
	// 	&pdb.DBOption{
	// 		Host:     confString("DB2_HOST", "127.0.0.1"),
	// 		Port:     confInt64("DB2_PORT", int64(5432)),
	// 		Database: confString("DB2_DATABASE", "GRIT-HRM"),
	// 		Username: confString("DB2_USER", "postgres"),
	// 		Password: confString("DB2_PASSWORD", ""),
	// 		TimeZone: confString("TZ", "Asia/Bangkok"),
	// 	},
	// )
	// log.Println("database2 connected success")
}

var (
	db *bun.DB
	// db2 *bun.DB
)

func GetDB() *bun.DB {
	return db
}

// func GetDB2() *bun.DB {
// 	return db2
// }
