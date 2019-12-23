package main

type HostsSource interface {
	Load() (records, error)
	Sync(records) error
}
