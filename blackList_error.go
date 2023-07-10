package errortool

type errorBlackList struct {
	Group map[errorGroupCode]struct{}
	Code  map[errorCode]struct{}
}

func SetErrorBlackList(group map[errorGroupCode]struct{}, code map[errorCode]struct{}) *errorBlackList {
	return &errorBlackList{
		Group: group,
		Code:  code,
	}
}

func AddGroup(group ...*errorGroup) map[errorGroupCode]struct{} {
	temp := make(map[errorGroupCode]struct{})
	for _, g := range group {
		temp[g.groupCode] = struct{}{}
	}
	return temp
}

func AddCode(code ...error) map[errorCode]struct{} {
	temp := make(map[errorCode]struct{})
	for _, c := range code {
		errStr, ok := parse(c)
		if !ok {
			continue
		}

		temp[errStr.code] = struct{}{}
	}

	return temp
}

func (ebl *errorBlackList) IsBlack(err *errorString) bool {
	if _, isExist := ebl.Group[err.groupCode]; isExist {
		return true
	}

	if _, isExist := ebl.Code[err.code]; isExist {
		return true
	}

	return false
}
