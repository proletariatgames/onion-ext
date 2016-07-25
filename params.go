package onionext

import (
	"time"

	"gopkg.in/fzerorubigd/onion.v2"
)

const (
	paramLoaded = iota
	paramIsSet  = iota << 1
)

type baseParam struct {
	state int
	key   string
}

type paramLoader interface {
	validate(set *ParamSet) bool
	parse(set *ParamSet)
}

func (p *baseParam) validate(set *ParamSet) bool {
	p.state |= paramLoaded

	if _, ok := set.onion.Get(p.key); ok {
		p.state |= paramIsSet
	}

	return p.IsSet()
}

func (p *baseParam) IsSet() bool {
	return (p.state & paramIsSet) != 0
}

type StringParam struct {
	baseParam
	value string
}

type StringSliceParam struct {
	baseParam
	value []string
}

type BoolParam struct {
	baseParam
	value bool
}

type IntParam struct {
	baseParam
	value int
}

type Int64Param struct {
	baseParam
	value int64
}

type Float32Param struct {
	baseParam
	value float32
}

type Float64Param struct {
	baseParam
	value float64
}

type DurationParam struct {
	baseParam
	value time.Duration
}

type ParamSet struct {
	onion  *onion.Onion
	params []paramLoader
}

func NewParamSet() *ParamSet {
	params := &ParamSet{onion.New(), nil}
	return params
}

func (p *ParamSet) Load(layers []onion.Layer) error {
	for _, layer := range layers {
		if err := p.onion.AddLayer(layer); err != nil {
			return err
		}
	}

	for _, param := range p.params {
		if param.validate(p) {
			param.parse(p)
		}
	}
	return nil
}

func (p *ParamSet) String(key string, def string) *StringParam {
	sp := &StringParam{baseParam{0, key}, def}
	p.params = append(p.params, sp)
	return sp
}

func (p *ParamSet) StringSlice(key string, def []string) *StringSliceParam {
	sp := &StringSliceParam{baseParam{0, key}, def}
	p.params = append(p.params, sp)
	return sp
}

func (p *ParamSet) Bool(key string, def bool) *BoolParam {
	bp := &BoolParam{baseParam{0, key}, def}
	p.params = append(p.params, bp)
	return bp
}

func (p *ParamSet) Int(key string, def int) *IntParam {
	ip := &IntParam{baseParam{0, key}, def}
	p.params = append(p.params, ip)
	return ip
}

func (p *ParamSet) Int64(key string, def int64) *Int64Param {
	ip := &Int64Param{baseParam{0, key}, def}
	p.params = append(p.params, ip)
	return ip
}

func (p *ParamSet) Float32(key string, def float32) *Float32Param {
	fp := &Float32Param{baseParam{0, key}, def}
	p.params = append(p.params, fp)
	return fp
}

func (p *ParamSet) Float64(key string, def float64) *Float64Param {
	fp := &Float64Param{baseParam{0, key}, def}
	p.params = append(p.params, fp)
	return fp
}

func (p *ParamSet) Duration(key string, def time.Duration) *DurationParam {
	dp := &DurationParam{baseParam{0, key}, def}
	p.params = append(p.params, dp)
	return dp
}

func (sp *StringParam) parse(set *ParamSet) { sp.value = set.onion.GetString(sp.baseParam.key) }
func (sp *StringSliceParam) parse(set *ParamSet) {
	sp.value = set.onion.GetStringSlice(sp.baseParam.key)
}
func (bp *BoolParam) parse(set *ParamSet)     { bp.value = set.onion.GetBool(bp.baseParam.key) }
func (ip *IntParam) parse(set *ParamSet)      { ip.value = set.onion.GetInt(ip.baseParam.key) }
func (ip *Int64Param) parse(set *ParamSet)    { ip.value = set.onion.GetInt64(ip.baseParam.key) }
func (fp *Float32Param) parse(set *ParamSet)  { fp.value = set.onion.GetFloat32(fp.baseParam.key) }
func (fp *Float64Param) parse(set *ParamSet)  { fp.value = set.onion.GetFloat64(fp.baseParam.key) }
func (dp *DurationParam) parse(set *ParamSet) { dp.value = set.onion.GetDuration(dp.baseParam.key) }

func (sp *StringParam) Get() string          { return sp.value }
func (sp *StringSliceParam) Get() []string   { return sp.value }
func (bp *BoolParam) Get() bool              { return bp.value }
func (ip *IntParam) Get() int                { return ip.value }
func (ip *Int64Param) Get() int64            { return ip.value }
func (fp *Float32Param) Get() float32        { return fp.value }
func (fp *Float64Param) Get() float64        { return fp.value }
func (dp *DurationParam) Get() time.Duration { return dp.value }
