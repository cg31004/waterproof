package errortool

var (
	Codes = Define()
	ErrDB = Codes.Plugin(dBErrorPackage).(dbError)
)
