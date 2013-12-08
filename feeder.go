package marto

type Feeder interface {
	Feed(key string) string
}

type MapFeeder struct {
	data map[string]string
}

func NewMapFeeder(data map[string]string) *MapFeeder {
	return &MapFeeder{data}
}

func (mf *MapFeeder) Feed(key string) string {
	return mf.data[key]
}