package errortool

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

func dBErrorPackage(groups iGroupRepository, codes iCodeRepository) interface{} {
	groupCode := genGroupCode(groupCodeDB)
	groups.Add(groupCode)
	group := &errorGroup{
		codes:     codes,
		groups:    groups,
		groupCode: errorGroupCode(groupCode),
	}

	return dbError{
		DBGroup:              group,
		NoRow:                group.Error(0, "DB no row"),
		CannotCreateTable:    group.Error(1, "DB cannot create table"),
		CannotCreateDatabase: group.Error(2, "DB cannot create database"),
		DatabaseCreateExists: group.Error(3, "DB database create exists"),
		TooManyConns:         group.Error(4, "DB too many conns"),
		AccessDenied:         group.Error(5, "DB access denied"),
		UnknownTable:         group.Error(6, "DB unknown table"),
		DuplicateEntry:       group.Error(7, "DB duplicate entry"),
		NoDefaultForField:    group.Error(8, "DB no default for field"),
	}
}

type dbError struct {
	DBGroup              *errorGroup
	NoRow                error
	CannotCreateTable    error
	CannotCreateDatabase error
	DatabaseCreateExists error
	TooManyConns         error
	AccessDenied         error
	UnknownTable         error
	DuplicateEntry       error
	NoDefaultForField    error
}

var (
	dbErrorCode = map[int]error{
		1005: ErrDB.CannotCreateTable,
		1006: ErrDB.CannotCreateDatabase,
		1007: ErrDB.DatabaseCreateExists,
		1040: ErrDB.TooManyConns,
		1045: ErrDB.AccessDenied,
		1051: ErrDB.UnknownTable,
		1062: ErrDB.DuplicateEntry,
		1364: ErrDB.NoDefaultForField,
	}
)

func ConvertDB(err error) error {
	if err == sql.ErrNoRows || err == gorm.ErrRecordNotFound {
		return ErrDB.NoRow
	}

	if e, ok := parseDBError(err); ok {
		return e
	}

	return err
}

func parseDBError(err error) (error, bool) {
	s := strings.TrimSpace(err.Error())
	data := strings.Split(s, ":")
	if len(data) == 0 {
		return nil, false
	}

	numStr := strings.ToLower(data[0])
	numStr = strings.Replace(numStr, "error", "", -1)
	numStr = strings.TrimSpace(numStr)
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return nil, false
	}

	e, ok := dbErrorCode[num]
	if !ok {
		return nil, false
	}

	return e, true
}

func Errorf(err error, args ...interface{}) error {
	errString, ok := parse(err)
	if !ok {
		return err
	}
	tempMsg := errString.GetMessage()
	tempMsg = fmt.Sprintf(tempMsg, args...)
	errString.setMessage(tempMsg)

	return errString
}
