package conn

import (
	"log"

	_ "github.com/lib/pq"
	"github.com/pudding-hack/backend/lib"

	"github.com/jmoiron/sqlx"
)

type SQLServerConnectionManager struct {
	db *sqlx.DB
}

const postgresqlDriver = "postgres"

func NewConnectionManager(cfg lib.DatabaseConfig) (*SQLServerConnectionManager, error) {
	db, err := sqlx.Connect(postgresqlDriver, cfg.DSN)
	if err != nil {
		log.Fatalf("failed to connect to sql server: %v", err)
		return nil, err
	}

	db.SetMaxIdleConns(cfg.MaxIdleConnections)
	db.SetMaxOpenConns(cfg.MaxOpenConnections)
	db.SetConnMaxIdleTime(cfg.MaxIdleDuration)
	db.SetConnMaxLifetime(cfg.MaxLifeTimeDuration)

	log.Println("connected to sql server")

	return &SQLServerConnectionManager{
		db: db,
	}, nil
}

func (cm *SQLServerConnectionManager) Close() error {
	log.Println("closing sql server connection")
	return cm.db.Close()
}

func (cm *SQLServerConnectionManager) GetQuery() *SingleInstruction {
	return NewSingleInstruction(cm.db)
}

func (cm *SQLServerConnectionManager) GetTransaction() *MultiInstruction {
	return NewMultiInstruction(cm.db)
}
