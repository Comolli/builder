package training

import (
	"context"
	"time"
)

type Json2Yaml struct {
	helper     *Helper
	options    *options
	cancelFunc context.CancelFunc
	context    context.Context
}

type JsonOperatrion interface {
	Do() *Helper
	Parse(ctx context.Context) error
}

func NewJson2Yaml(json interface{}, opts ...Option) JsonOperatrion {
	o := newOptions(opts...)

	if err := o.precheckFunc(json); err != nil {
		o.errHandler(err)
	}

	_, cancel := context.WithCancel(o.ctx)
	j := &Json2Yaml{
		helper: &Helper{
			self: json,
		},
		cancelFunc: cancel,
		context:    context.Background(),
		options:    o,
	}

	return j
}

func (j *Json2Yaml) Parse(ctx context.Context) error {
	if err := j.helper.GenerateJson(j.helper.self); err != nil {
		return err
	}

	return nil
}

func (j *Json2Yaml) Do() *Helper {
	ctx, cancel := context.WithTimeout(j.options.ctx, 2*time.Second)
	defer cancel()

	if err := j.Parse(ctx); err != nil {
		j.options.errHandler(err)
		return nil
	}

	select {
	case <-ctx.Done():
		j.options.errHandler(ctx.Err())
		return nil
	default:
		//avoid blocking
	}

	return j.helper
}
