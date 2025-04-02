package db

import (
	"database/sql"
	"fmt"
	"main/pkg/config"
	"main/pkg/model"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() error {
	dbConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("error loading config: %v", err)
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName)

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return fmt.Errorf("error connecting to the database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("error pinging the database: %v", err)
	}

	_, err = db.Exec(`DROP TABLE IF EXISTS os_and_osquery`)
	if err != nil {
		return fmt.Errorf("error dropping table: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE os_and_osquery (
			id SERIAL PRIMARY KEY,
			pid INTEGER,
			uuid TEXT,
			instance_id TEXT,
			version TEXT,
			config_hash TEXT,
			config_valid TEXT,
			extensions TEXT,
			build_platform TEXT,
			build_distro TEXT,
			start_time INTEGER,
			watcher INTEGER,
			platform_mask INTEGER,
			name VARCHAR(255),
			major VARCHAR(50),
			minor VARCHAR(50),
			patch VARCHAR(50),
			build VARCHAR(255),
			platform VARCHAR(255),
			platform_like VARCHAR(255),
			codename VARCHAR(255),
			arch VARCHAR(50),
			extra TEXT
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating table: %v", err)
	}

	return nil
}

func SaveOSVersion(osVersion model.OSVersion) error {
	_, err := db.Exec(`
		INSERT INTO os_and_osquery (
			name, version, major, minor, patch, build, 
			platform, platform_like, codename, arch, extra
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		)`,
		osVersion.Name,
		osVersion.Version,
		osVersion.Major,
		osVersion.Minor,
		osVersion.Patch,
		osVersion.Build,
		osVersion.Platform,
		osVersion.PlatformLike,
		osVersion.Codename,
		osVersion.Arch,
		osVersion.Extra,
	)
	return err
}

func SaveOSQueryVersion(osQueryVersion model.OSQueryVersion) error {
	_, err := db.Exec(`
		INSERT INTO os_and_osquery (
			pid, uuid, instance_id, version, config_hash, config_valid,
			extensions, build_platform, build_distro, start_time,
			watcher, platform_mask
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
		)`,
		osQueryVersion.PID,
		osQueryVersion.UUID,
		osQueryVersion.InstanceID,
		osQueryVersion.Version,
		osQueryVersion.ConfigHash,
		osQueryVersion.ConfigValid,
		osQueryVersion.Extensions,
		osQueryVersion.BuildPlatform,
		osQueryVersion.BuildDistro,
		osQueryVersion.StartTime,
		osQueryVersion.Watcher,
		osQueryVersion.PlatformMask,
	)
	return err
}

func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
