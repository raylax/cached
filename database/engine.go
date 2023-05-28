package database

import (
	"errors"
	"github.com/dgraph-io/badger/v3"
	"github.com/hashicorp/go-hclog"
)

type Transaction interface {
}

type Engine interface {
	Close() error
}

type Options struct {
	Path   string
	Logger Logger
}

func (opts *Options) WithPath(path string) *Options {
	opts.Path = path
	return opts
}

func (opts *Options) WithLogger(logger hclog.Logger) *Options {
	opts.Logger = newLoggerAdapter(logger)
	return opts
}

func NewOptions() *Options {
	return &Options{}
}

func NewDefaultOptions() *Options {
	return &Options{}
}

func Open(opts *Options) (Engine, error) {

	if opts.Path == "" {
		return nil, errors.New("path is required")
	}

	if opts.Logger == nil {
		return nil, errors.New("logger is required")
	}

	options := badger.DefaultOptions(opts.Path).
		WithLogger(opts.Logger)
	db, err := badger.Open(options)

	if err != nil {
		return nil, err
	}

	database := &engineImpl{
		badger: db,
	}

	return database, nil
}
