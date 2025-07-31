package ageo

import (
	"github.com/fsnotify/fsnotify"
	"github.com/jpfluger/alibs-slim/alog"
	"github.com/oschwald/geoip2-golang"
	"net"
	"path/filepath"
	"sync"
	"time"
)

var (
	db       *geoip2.Reader
	dbPath   string
	dbLock   sync.RWMutex
	watcher  *fsnotify.Watcher
	stopChan chan struct{}

	initOnce sync.Once
)

// InitGeoInfo initializes the database and starts watching the DB file for changes.
func InitGeoInfo(path string, enableWatchDb bool) error {
	var initErr error
	initOnce.Do(func() {
		dbLock.Lock()
		defer dbLock.Unlock()

		dbPath = path
		db, initErr = geoip2.Open(path)
		if initErr != nil {
			return
		}

		watcher, initErr = fsnotify.NewWatcher()
		if initErr != nil {
			return
		}
		stopChan = make(chan struct{})

		if enableWatchDb {
			go watchDBFile()
		}

		dir := filepath.Dir(path)
		initErr = watcher.Add(dir)
	})
	return initErr
}

// CloseGeoInfo shuts down the watcher and closes the DB.
func CloseGeoInfo() {
	dbLock.Lock()
	defer dbLock.Unlock()

	if stopChan != nil {
		close(stopChan)
		stopChan = nil
	}

	if watcher != nil {
		_ = watcher.Close()
		watcher = nil
	}

	if db != nil {
		_ = db.Close()
		db = nil
	}

	initOnce = sync.Once{} // reset for future InitGeoInfo
}

// MustLookupGeoInfoForIP guarantees a non-nil GeoInfo.
// If the underlying lookup fails or returns nil, an empty GeoInfo is returned.
func MustLookupGeoInfoForIP(ipStr string) *GeoInfo {
	if info := LookupGeoInfoForIP(ipStr); info != nil {
		return info
	}
	return &GeoInfo{}
}

// LookupGeoInfoForIP returns a GeoInfo based on the IP address.
func LookupGeoInfoForIP(ipStr string) *GeoInfo {
	dbLock.RLock()
	defer dbLock.RUnlock()

	if db == nil {
		return nil
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil
	}

	record, err := db.City(ip)
	if err != nil {
		return nil
	}

	var region string
	if len(record.Subdivisions) > 0 {
		region = record.Subdivisions[0].Names["en"]
	}

	return &GeoInfo{
		CountryCode: record.Country.IsoCode,
		IsEU:        record.Country.IsInEuropeanUnion,
		City:        record.City.Names["en"],
		Region:      region,
		GISPoint: GISPoint{
			Latitude:  record.Location.Latitude,
			Longitude: record.Location.Longitude,
		},
		IPv4: ipStr,
	}
}

// watchDBFile watches the DB file's directory and reloads if it changes.
func watchDBFile() {
	dbFilename := filepath.Base(dbPath)

	for {
		select {
		case <-stopChan:
			return
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if filepath.Base(event.Name) == dbFilename && (event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create) {
				alog.LOGGER(alog.LOGGER_APP).Info().Msg("geo reloading GeoIP DB")
				reloadDB()
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			alog.LOGGER(alog.LOGGER_APP).Err(err).Msg("geo watcher error")
		}
	}
}

func reloadDB() {
	dbLock.Lock()
	defer dbLock.Unlock()

	time.Sleep(200 * time.Millisecond) // Slight delay for file write to settle

	newDB, err := geoip2.Open(dbPath)
	if err != nil {
		alog.LOGGER(alog.LOGGER_APP).Err(err).Msg("geo failed to reload GeoIP DB")
		return
	}

	if db != nil {
		_ = db.Close()
	}
	db = newDB
	alog.LOGGER(alog.LOGGER_APP).Info().Msg("geo successfully reloaded GeoIP DB")
}
