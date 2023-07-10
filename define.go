package errortool

import (
	"log"
	"sort"
)

const (
	groupCodeDB int = 999
	errCodeMax      = 999
)

func Define() *define {
	return &define{
		groups: newGroupRepository(),
		codes:  newCodeRepository(),
	}
}

type define struct {
	groups iGroupRepository
	codes  iCodeRepository
}

func (d *define) Group(group int) *errorGroup {
	group = genGroupCode(group)
	d.groups.Add(group)
	return &errorGroup{
		codes:     d.codes,
		groups:    d.groups,
		groupCode: errorGroupCode(group),
	}
}

func genGroupCode(groupBase int) int {
	return groupBase * 1000
}

func (d *define) Plugin(f func(groups iGroupRepository, codes iCodeRepository) interface{}) interface{} {
	return f(d.groups, d.codes)
}

func (d *define) List() []errorString {
	keys := d.codes.Keys()
	sort.SliceStable(keys,
		func(i, j int) bool {
			return keys[i] < keys[j]
		})

	res := make([]errorString, len(keys))
	for i, v := range keys {
		if val, ok := d.codes.Get(v); ok {
			res[i] = *val
		} else {
			res[i] = errorString{}
		}
	}

	return res
}

type errorGroup struct {
	codes     iCodeRepository
	groups    iGroupRepository
	groupCode errorGroupCode
}

func (e *errorGroup) Error(code int, message string) error {
	if code > errCodeMax {
		log.Panicf("errorGroup error: code max than 999, code: %d", code)
	}

	errCode := e.makeErrorCode(e.groups.Get(int(e.groupCode)), code)
	err := &errorString{
		code:      errCode,
		groupCode: errorGroupCode(e.groupCode),
		baseCode:  errorBaseCode(code),
		message:   message,
	}
	e.codes.Add(errCode, err)
	return err
}

func (e *errorGroup) makeErrorCode(groupCode, code int) errorCode {
	return errorCode(groupCode + code)
}
