package weberrors

type DetailedError struct {
	Location string      `json:"location"`
	Param    string      `json:"param"`
	Value    interface{} `json:"value"`
	Msg      string      `json:"msg"`
}

type DetailedErrors struct {
	Errors []DetailedError `json:"errors"`
}
