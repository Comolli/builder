package training

import "context"

type (
	PrecheckFunc func(self interface{}) error
	ErrHandler   func(err error)
)

type Option func(o *options)

type options struct {
	ctx          context.Context
	precheckFunc PrecheckFunc
	errHandler   ErrHandler
}

func WithContext(ctx context.Context) Option {
	return func(o *options) {
		o.ctx = ctx
	}
}

func WithPrecheckFunc(precheckFunc PrecheckFunc) Option {
	return func(o *options) {
		o.precheckFunc = precheckFunc
	}
}

func WithErrHandler(errHandler ErrHandler) Option {
	return func(o *options) {
		o.errHandler = errHandler
	}
}

func newOptions(opts ...Option) *options {
	o := &options{
		ctx:        context.Background(),
		errHandler: func(err error) {},
		precheckFunc: func(self interface{}) error {
			return nil
		},
	}

	for _, opt := range opts {
		opt(o)
	}

	return o
}
