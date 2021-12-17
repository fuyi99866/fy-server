package robot_service

import "crypto/tls"

/**
Mqtt的连接设置
*/
func defaultOptions() mqOptions {
	opts := mqOptions{
		clientId: "",
		username: "",
		password: "",
		tls:      true,
		retry:    5000,
	}
	return opts
}

type mqOptions struct {
	url              string
	clientId         string
	username         string
	password         string
	tls              bool
	OnConnect        ConnectedHandler
	OnConnectionLost ConnectionLostHandler
	OnReconnecting   ReconnectHandler
	OnMessageComing  MessageHandler
	retry            int
	autoReconnect    bool
}

// Option set option by close func 通过闭合函数设置option
type Option func(*mqOptions)

// WithClientId set clientId name  to mq clientId 设置clientId
func WithClientId(c string) Option {
	return func(options *mqOptions) {
		if c == "" {
			return
		}
		options.username = c
	}
}

// WithUserName set username 设置 username
func WithUserName(u string) Option {
	return func(options *mqOptions) {
		if u == "" {
			return
		}
		options.username = u
	}
}

// WithPassword set password //设置密码
func WithPassword(p string) Option {
	return func(opts *mqOptions) {
		if p == "" {
			return
		}
		opts.password = p
	}
}

// WithTls set tls //设置tls
func WithTls(t bool) Option {
	return func(opts *mqOptions) {
		opts.tls = t
	}
}

// WithOnConnectHandler //处理连接
func WithConnectedHandler(onConn ConnectedHandler) Option {
	return func(opts *mqOptions) {
		opts.OnConnect = onConn
	}
}

// WithConnectionLostHandler //处理断开
func WithConnectionLostHandler(onLost ConnectionLostHandler) Option {
	return func(opts *mqOptions) {
		opts.OnConnectionLost = onLost
	}
}

// WithReconnectingHandler //处理重连
func WithReconnectingHandler(rc ReconnectHandler) Option {
	return func(opts *mqOptions) {
		opts.OnReconnecting = rc
	}
}

// WithMessageHandler //处理消息
func WithMessageHandler(mh MessageHandler) Option {
	return func(opts *mqOptions) {
		opts.OnMessageComing = mh
	}
}

// WithRetry set retry to mq group retry //设置重连次数
func WithRetry(t int) Option {
	return func(opts *mqOptions) {
		if t < 0 {
			return
		}
		opts.retry = t
	}
}

// WithAutoReconnect is reconnect wile true //设置是否自动重连
func WithAutoReconnect(a bool) Option {
	return func(opts *mqOptions) {
		opts.autoReconnect = a
	}
}

//TLS称为安全传输层协议
func DefaultTlsConfig() *tls.Config {
	return &tls.Config{
		Rand:                        nil,
		Time:                        nil,
		Certificates:                nil,
		NameToCertificate:           nil,
		GetCertificate:              nil,
		GetClientCertificate:        nil,
		GetConfigForClient:          nil,
		VerifyPeerCertificate:       nil,
		VerifyConnection:            nil,
		RootCAs:                     nil,
		NextProtos:                  nil,
		ServerName:                  "",
		ClientAuth:                  tls.NoClientCert,
		ClientCAs:                   nil,
		InsecureSkipVerify:          false,//skip verify the server cert
		CipherSuites:                nil,
		PreferServerCipherSuites:    false,
		SessionTicketsDisabled:      false,
		SessionTicketKey:            [32]byte{},
		ClientSessionCache:          nil,
		MinVersion:                  0,
		MaxVersion:                  0,
		CurvePreferences:            nil,
		DynamicRecordSizingDisabled: false,
		Renegotiation:               0,
		KeyLogWriter:                nil,
	}
}

// MessageHandler receive message with subscribe call back 处理收到的订阅消息
type MessageHandler func(topic string, data []byte)

// ConnectionLostHandler 断开连接处理
type ConnectionLostHandler func(error)

// OnConnectHandler 连接成功处理
type ConnectedHandler func()

// ReconnectHandler 重连处理
type ReconnectHandler func()
