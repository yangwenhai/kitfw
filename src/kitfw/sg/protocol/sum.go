package protocol

type SumRequest struct {
	UserId int64 `capid:"0"`
	Num1   int64 `capid:"1"`
	Num2   int64 `capid:"2"`
}

type SumReply struct {
	RetCode int8  `capid:"0"`
	Val     int64 `capid:"1"`
}

type ConcatRequest struct {
	UserId int64  `capid:"0"`
	Str1   string `capid:"1"`
	Str2   string `capid:"2"`
}

type ConcatReply struct {
	RetCode int8   `capid:"0"`
	Val     string `capid:"1"`
}
