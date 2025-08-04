package mapper

import "github.com/mitchellh/mapstructure"

func NewDecoder[R interface{}](result R, fs ...mapstructure.DecodeHookFunc) (*mapstructure.Decoder, error) {
	return mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(fs...),
		Result:     result,
	})
}

func Decode[I any, R any](input I, result R, fs ...mapstructure.DecodeHookFunc) error {
	decoder, err := NewDecoder(result, fs...)
	if err != nil {
		return err
	}
	return decoder.Decode(input)
}
