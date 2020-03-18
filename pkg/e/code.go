package e

const (
	ERROR              = 10001
	ERROR_BINDING_FORM = 40001

	ERROR_USER_EXIST = 40002

	ERROR_WRONG_QUERY  = 50001
	ERROR_FAIL_ENCRYPT = 50002
)

var codeMsg = map[int]string{
	ERROR:              "unexpected error",
	ERROR_BINDING_FORM: "missing some required field or format not match",
	ERROR_USER_EXIST:   "email or nickname already exists",
	ERROR_WRONG_QUERY:  "some query is wrong",
	ERROR_FAIL_ENCRYPT: "fail to encrypt password",
}

func GetCodeMsg(code int) string {
	if v, ok := codeMsg[code]; ok {
		return v
	}
	return codeMsg[ERROR]
}
