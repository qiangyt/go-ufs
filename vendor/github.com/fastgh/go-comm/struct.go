package comm

import (
	"github.com/pkg/errors"
	"github.com/tiaotiao/mapstruct"
)

func StructToMap(src any) map[string]any {
	return mapstruct.Struct2Map(src)
}

func Map2Struct(src map[string]any, dest any) {
	if err := mapstruct.Map2Struct(src, dest); err != nil {
		panic(errors.Wrap(err, "failed to convert map to struct"))
	}
}
