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

	_, err = db.Exec(`DROP TABLE IF EXISTS apps`)
	if err != nil {
		return fmt.Errorf("error dropping apps table: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE apps (
			id SERIAL PRIMARY KEY,
			name TEXT,
			path TEXT NOT NULL,
			bundle_executable TEXT,
			bundle_identifier TEXT NOT NULL,
			bundle_name TEXT,
			bundle_short_version TEXT,
			bundle_version TEXT,
			bundle_package_type TEXT,
			environment TEXT,
			element TEXT,
			compiler TEXT,
			development_region TEXT,
			display_name TEXT,
			info_string TEXT,
			minimum_system_version TEXT,
			category TEXT,
			applescript_enabled TEXT,
			copyright TEXT,
			last_opened_time TEXT,
			UNIQUE (bundle_identifier, path)
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating apps table: %v", err)
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

func SaveApplication(app model.Application) error {
	if app.BundleIdentifier == "" || app.Path == "" {
		return nil
	}

	_, err := db.Exec(`
		INSERT INTO apps (
			name, path, bundle_executable, bundle_identifier, bundle_name,
			bundle_short_version, bundle_version, bundle_package_type,
			environment, element, compiler, development_region,
			display_name, info_string, minimum_system_version,
			category, applescript_enabled, copyright, last_opened_time
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12,
			$13, $14, $15, $16, $17, $18, $19
		)
		ON CONFLICT (bundle_identifier, path) DO UPDATE SET
			name = EXCLUDED.name,
			bundle_executable = EXCLUDED.bundle_executable,
			bundle_name = EXCLUDED.bundle_name,
			bundle_short_version = EXCLUDED.bundle_short_version,
			bundle_version = EXCLUDED.bundle_version,
			bundle_package_type = EXCLUDED.bundle_package_type,
			environment = EXCLUDED.environment,
			element = EXCLUDED.element,
			compiler = EXCLUDED.compiler,
			development_region = EXCLUDED.development_region,
			display_name = EXCLUDED.display_name,
			info_string = EXCLUDED.info_string,
			minimum_system_version = EXCLUDED.minimum_system_version,
			category = EXCLUDED.category,
			applescript_enabled = EXCLUDED.applescript_enabled,
			copyright = EXCLUDED.copyright,
			last_opened_time = EXCLUDED.last_opened_time`,
		app.Name,
		app.Path,
		app.BundleExecutable,
		app.BundleIdentifier,
		app.BundleName,
		app.BundleShortVersion,
		app.BundleVersion,
		app.BundlePackageType,
		app.Environment,
		app.Element,
		app.Compiler,
		app.DevelopmentRegion,
		app.DisplayName,
		app.InfoString,
		app.MinimumSystemVersion,
		app.Category,
		app.ApplescriptEnabled,
		app.Copyright,
		app.LastOpenedTime,
	)
	return err
}

func SaveOSAndOSQueryInfo(osVersion model.OSVersion, osQueryVersion model.OSQueryVersion) error {
	_, err := db.Exec(`
		INSERT INTO os_and_osquery (
			name, version, major, minor, patch, build, 
			platform, platform_like, codename, arch, extra,
			pid, uuid, instance_id, config_hash, config_valid,
			extensions, build_platform, build_distro, start_time,
			watcher, platform_mask
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11,
			$12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22
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
		osQueryVersion.PID,
		osQueryVersion.UUID,
		osQueryVersion.InstanceID,
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
