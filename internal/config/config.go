package config

type StoreMode int

const (
	REPLICA StoreMode = iota
	MASTER
)

var Mode StoreMode = MASTER

var ExternalPort string = ":8080"
var InternalPort string = ":8090"

var MaxItemSize int = 1024 * 1000 // 1 MB
var LogFile string = "./logs/log.txt"

var SnapshotPath string = "./snapshots/"
var StoragePath string = "./storage/"
