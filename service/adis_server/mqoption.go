package adis_server

import (
	"go_server/pkg/util"
	"golang.org/x/net/context"
)

type mqOptions struct {
	addr             string
	groupName        string
	retry            int
	consumerInstance string
	ctx              context.Context
}

func defaultOptions() mqOptions {
	opts := mqOptions{
		addr:             "0.0.0.0:9876",
		groupName:        "default",
		retry:            1,
		consumerInstance: util.UUIDShort(),
		ctx:              context.Background(),
	}
	return opts
}

// Option set option by close func
type Option func(*mqOptions)

func WithAddr(addr string) Option {
	return func(options *mqOptions) {
		if addr == "" {
			return
		}
		options.addr = addr
	}
}

func WithGroupName(g string) Option {
	return func(options *mqOptions) {
		if g == "" {
			return
		}
		options.groupName = g
	}
}

func WithRetry(t int) Option {
	return func(options *mqOptions) {
		if t < 0 {
			return
		}
		options.retry = t
	}
}

func WithContext(ctx context.Context) Option {
	return func(opts *mqOptions) {
		if ctx == nil {
			return
		}
		opts.ctx = ctx
	}
}

func WithInstanceName(s string) Option {
	return func(opts *mqOptions) {
		if s == "" {
			return
		}
		opts.consumerInstance = s
	}
}

type ReceiveCallBack func(topic, msgid string, body []byte)
