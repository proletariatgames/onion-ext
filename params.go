package onionext

import (
	"net/url"
	"time"

	"github.com/fzerorubigd/onion"
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
	init(set *ParamSet) bool
	parse(found bool, set *ParamSet) error
}

func (p *baseParam) init(set *ParamSet) bool {
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

type URLParam struct {
	baseParam
	def   string
	value *url.URL
}

type URLSliceParam struct {
	baseParam
	def   []string
	value []*url.URL
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
		if err := param.parse(param.init(p), p); err != nil {
			return err
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

func (p *ParamSet) URL(key string, def string) *URLParam {
	sp := &URLParam{baseParam{0, key}, def, nil}
	p.params = append(p.params, sp)
	return sp
}

func (p *ParamSet) URLSlice(key string, def []string) *URLSliceParam {
	sp := &URLSliceParam{baseParam{0, key}, def, nil}
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

func (sp *StringParam) parse(found bool, set *ParamSet) error {
	if found {
		sp.value = set.onion.GetString(sp.baseParam.key)
	}
	return nil
}

func (sp *StringSliceParam) parse(found bool, set *ParamSet) error {
	if found {
		sp.value = set.onion.GetStringSlice(sp.baseParam.key)
	}
	return nil
}

func parseURLS(strs []string) ([]*url.URL, error) {
	parsed := make([]*url.URL, len(strs))
	for i, str := range strs {
		if value, err := url.Parse(str); err != nil {
			return nil, err
		} else {
			parsed[i] = value
		}
	}
	return parsed, nil
}

func (sp *URLParam) parse(found bool, set *ParamSet) error {
	if parsed, err := parseURLS([]string{sp.def}); err != nil {
		return err
	} else {
		sp.value = parsed[0]
	}

	if found {
		if parsed, err := parseURLS([]string{set.onion.GetString(sp.key)}); err != nil {
			return err
		} else {
			sp.value = parsed[0]
		}
	}

	return nil
}

func (sp *URLSliceParam) parse(found bool, set *ParamSet) error {
	if parsed, err := parseURLS(sp.def); err != nil {
		return err
	} else {
		sp.value = parsed
	}

	if found {
		if parsed, err := parseURLS(set.onion.GetStringSlice(sp.key)); err != nil {
			return err
		} else {
			sp.value = parsed
		}
	}

	return nil
}

func (bp *BoolParam) parse(found bool, set *ParamSet) error {
	if found {
		bp.value = set.onion.GetBool(bp.baseParam.key)
	}
	return nil
}

func (ip *IntParam) parse(found bool, set *ParamSet) error {
	if found {
		ip.value = set.onion.GetInt(ip.baseParam.key)
	}
	return nil
}

func (ip *Int64Param) parse(found bool, set *ParamSet) error {
	if found {
		ip.value = set.onion.GetInt64(ip.baseParam.key)
	}
	return nil
}

func (fp *Float32Param) parse(found bool, set *ParamSet) error {
	if found {
		fp.value = set.onion.GetFloat32(fp.baseParam.key)
	}
	return nil
}

func (fp *Float64Param) parse(found bool, set *ParamSet) error {
	if found {
		fp.value = set.onion.GetFloat64(fp.baseParam.key)
	}
	return nil
}

func (dp *DurationParam) parse(found bool, set *ParamSet) error {
	if found {
		dp.value = set.onion.GetDuration(dp.baseParam.key)
	}
	return nil
}

func (sp *StringParam) Get() string          { return sp.value }
func (sp *StringSliceParam) Get() []string   { return sp.value }
func (up *URLParam) Get() *url.URL           { return up.value }
func (up *URLSliceParam) Get() []*url.URL    { return up.value }
func (bp *BoolParam) Get() bool              { return bp.value }
func (ip *IntParam) Get() int                { return ip.value }
func (ip *Int64Param) Get() int64            { return ip.value }
func (fp *Float32Param) Get() float32        { return fp.value }
func (fp *Float64Param) Get() float64        { return fp.value }
func (dp *DurationParam) Get() time.Duration { return dp.value }
