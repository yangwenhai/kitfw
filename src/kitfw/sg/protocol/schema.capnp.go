package protocol

// AUTO GENERATED - DO NOT EDIT

import (
	"bufio"
	"bytes"
	"encoding/json"
	C "github.com/glycerine/go-capnproto"
	"io"
)

type ConcatReplyCapn C.Struct

func NewConcatReplyCapn(s *C.Segment) ConcatReplyCapn { return ConcatReplyCapn(s.NewStruct(8, 1)) }
func NewRootConcatReplyCapn(s *C.Segment) ConcatReplyCapn {
	return ConcatReplyCapn(s.NewRootStruct(8, 1))
}
func AutoNewConcatReplyCapn(s *C.Segment) ConcatReplyCapn { return ConcatReplyCapn(s.NewStructAR(8, 1)) }
func ReadRootConcatReplyCapn(s *C.Segment) ConcatReplyCapn {
	return ConcatReplyCapn(s.Root(0).ToStruct())
}
func (s ConcatReplyCapn) RetCode() int8     { return int8(C.Struct(s).Get8(0)) }
func (s ConcatReplyCapn) SetRetCode(v int8) { C.Struct(s).Set8(0, uint8(v)) }
func (s ConcatReplyCapn) Val() string       { return C.Struct(s).GetObject(0).ToText() }
func (s ConcatReplyCapn) ValBytes() []byte  { return C.Struct(s).GetObject(0).ToData() }
func (s ConcatReplyCapn) SetVal(v string)   { C.Struct(s).SetObject(0, s.Segment.NewText(v)) }
func (s ConcatReplyCapn) WriteJSON(w io.Writer) error {
	b := bufio.NewWriter(w)
	var err error
	var buf []byte
	_ = buf
	err = b.WriteByte('{')
	if err != nil {
		return err
	}
	_, err = b.WriteString("\"retCode\":")
	if err != nil {
		return err
	}
	{
		s := s.RetCode()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	err = b.WriteByte(',')
	if err != nil {
		return err
	}
	_, err = b.WriteString("\"val\":")
	if err != nil {
		return err
	}
	{
		s := s.Val()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	err = b.WriteByte('}')
	if err != nil {
		return err
	}
	err = b.Flush()
	return err
}
func (s ConcatReplyCapn) MarshalJSON() ([]byte, error) {
	b := bytes.Buffer{}
	err := s.WriteJSON(&b)
	return b.Bytes(), err
}
func (s ConcatReplyCapn) WriteCapLit(w io.Writer) error {
	b := bufio.NewWriter(w)
	var err error
	var buf []byte
	_ = buf
	err = b.WriteByte('(')
	if err != nil {
		return err
	}
	_, err = b.WriteString("retCode = ")
	if err != nil {
		return err
	}
	{
		s := s.RetCode()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	_, err = b.WriteString(", ")
	if err != nil {
		return err
	}
	_, err = b.WriteString("val = ")
	if err != nil {
		return err
	}
	{
		s := s.Val()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	err = b.WriteByte(')')
	if err != nil {
		return err
	}
	err = b.Flush()
	return err
}
func (s ConcatReplyCapn) MarshalCapLit() ([]byte, error) {
	b := bytes.Buffer{}
	err := s.WriteCapLit(&b)
	return b.Bytes(), err
}

type ConcatReplyCapn_List C.PointerList

func NewConcatReplyCapnList(s *C.Segment, sz int) ConcatReplyCapn_List {
	return ConcatReplyCapn_List(s.NewCompositeList(8, 1, sz))
}
func (s ConcatReplyCapn_List) Len() int { return C.PointerList(s).Len() }
func (s ConcatReplyCapn_List) At(i int) ConcatReplyCapn {
	return ConcatReplyCapn(C.PointerList(s).At(i).ToStruct())
}
func (s ConcatReplyCapn_List) ToArray() []ConcatReplyCapn {
	n := s.Len()
	a := make([]ConcatReplyCapn, n)
	for i := 0; i < n; i++ {
		a[i] = s.At(i)
	}
	return a
}
func (s ConcatReplyCapn_List) Set(i int, item ConcatReplyCapn) {
	C.PointerList(s).Set(i, C.Object(item))
}

type ConcatRequestCapn C.Struct

func NewConcatRequestCapn(s *C.Segment) ConcatRequestCapn { return ConcatRequestCapn(s.NewStruct(8, 2)) }
func NewRootConcatRequestCapn(s *C.Segment) ConcatRequestCapn {
	return ConcatRequestCapn(s.NewRootStruct(8, 2))
}
func AutoNewConcatRequestCapn(s *C.Segment) ConcatRequestCapn {
	return ConcatRequestCapn(s.NewStructAR(8, 2))
}
func ReadRootConcatRequestCapn(s *C.Segment) ConcatRequestCapn {
	return ConcatRequestCapn(s.Root(0).ToStruct())
}
func (s ConcatRequestCapn) UserId() int64     { return int64(C.Struct(s).Get64(0)) }
func (s ConcatRequestCapn) SetUserId(v int64) { C.Struct(s).Set64(0, uint64(v)) }
func (s ConcatRequestCapn) Str1() string      { return C.Struct(s).GetObject(0).ToText() }
func (s ConcatRequestCapn) Str1Bytes() []byte { return C.Struct(s).GetObject(0).ToData() }
func (s ConcatRequestCapn) SetStr1(v string)  { C.Struct(s).SetObject(0, s.Segment.NewText(v)) }
func (s ConcatRequestCapn) Str2() string      { return C.Struct(s).GetObject(1).ToText() }
func (s ConcatRequestCapn) Str2Bytes() []byte { return C.Struct(s).GetObject(1).ToData() }
func (s ConcatRequestCapn) SetStr2(v string)  { C.Struct(s).SetObject(1, s.Segment.NewText(v)) }
func (s ConcatRequestCapn) WriteJSON(w io.Writer) error {
	b := bufio.NewWriter(w)
	var err error
	var buf []byte
	_ = buf
	err = b.WriteByte('{')
	if err != nil {
		return err
	}
	_, err = b.WriteString("\"userId\":")
	if err != nil {
		return err
	}
	{
		s := s.UserId()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	err = b.WriteByte(',')
	if err != nil {
		return err
	}
	_, err = b.WriteString("\"str1\":")
	if err != nil {
		return err
	}
	{
		s := s.Str1()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	err = b.WriteByte(',')
	if err != nil {
		return err
	}
	_, err = b.WriteString("\"str2\":")
	if err != nil {
		return err
	}
	{
		s := s.Str2()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	err = b.WriteByte('}')
	if err != nil {
		return err
	}
	err = b.Flush()
	return err
}
func (s ConcatRequestCapn) MarshalJSON() ([]byte, error) {
	b := bytes.Buffer{}
	err := s.WriteJSON(&b)
	return b.Bytes(), err
}
func (s ConcatRequestCapn) WriteCapLit(w io.Writer) error {
	b := bufio.NewWriter(w)
	var err error
	var buf []byte
	_ = buf
	err = b.WriteByte('(')
	if err != nil {
		return err
	}
	_, err = b.WriteString("userId = ")
	if err != nil {
		return err
	}
	{
		s := s.UserId()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	_, err = b.WriteString(", ")
	if err != nil {
		return err
	}
	_, err = b.WriteString("str1 = ")
	if err != nil {
		return err
	}
	{
		s := s.Str1()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	_, err = b.WriteString(", ")
	if err != nil {
		return err
	}
	_, err = b.WriteString("str2 = ")
	if err != nil {
		return err
	}
	{
		s := s.Str2()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	err = b.WriteByte(')')
	if err != nil {
		return err
	}
	err = b.Flush()
	return err
}
func (s ConcatRequestCapn) MarshalCapLit() ([]byte, error) {
	b := bytes.Buffer{}
	err := s.WriteCapLit(&b)
	return b.Bytes(), err
}

type ConcatRequestCapn_List C.PointerList

func NewConcatRequestCapnList(s *C.Segment, sz int) ConcatRequestCapn_List {
	return ConcatRequestCapn_List(s.NewCompositeList(8, 2, sz))
}
func (s ConcatRequestCapn_List) Len() int { return C.PointerList(s).Len() }
func (s ConcatRequestCapn_List) At(i int) ConcatRequestCapn {
	return ConcatRequestCapn(C.PointerList(s).At(i).ToStruct())
}
func (s ConcatRequestCapn_List) ToArray() []ConcatRequestCapn {
	n := s.Len()
	a := make([]ConcatRequestCapn, n)
	for i := 0; i < n; i++ {
		a[i] = s.At(i)
	}
	return a
}
func (s ConcatRequestCapn_List) Set(i int, item ConcatRequestCapn) {
	C.PointerList(s).Set(i, C.Object(item))
}

type SumReplyCapn C.Struct

func NewSumReplyCapn(s *C.Segment) SumReplyCapn      { return SumReplyCapn(s.NewStruct(16, 0)) }
func NewRootSumReplyCapn(s *C.Segment) SumReplyCapn  { return SumReplyCapn(s.NewRootStruct(16, 0)) }
func AutoNewSumReplyCapn(s *C.Segment) SumReplyCapn  { return SumReplyCapn(s.NewStructAR(16, 0)) }
func ReadRootSumReplyCapn(s *C.Segment) SumReplyCapn { return SumReplyCapn(s.Root(0).ToStruct()) }
func (s SumReplyCapn) RetCode() int8                 { return int8(C.Struct(s).Get8(0)) }
func (s SumReplyCapn) SetRetCode(v int8)             { C.Struct(s).Set8(0, uint8(v)) }
func (s SumReplyCapn) Val() int64                    { return int64(C.Struct(s).Get64(8)) }
func (s SumReplyCapn) SetVal(v int64)                { C.Struct(s).Set64(8, uint64(v)) }
func (s SumReplyCapn) WriteJSON(w io.Writer) error {
	b := bufio.NewWriter(w)
	var err error
	var buf []byte
	_ = buf
	err = b.WriteByte('{')
	if err != nil {
		return err
	}
	_, err = b.WriteString("\"retCode\":")
	if err != nil {
		return err
	}
	{
		s := s.RetCode()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	err = b.WriteByte(',')
	if err != nil {
		return err
	}
	_, err = b.WriteString("\"val\":")
	if err != nil {
		return err
	}
	{
		s := s.Val()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	err = b.WriteByte('}')
	if err != nil {
		return err
	}
	err = b.Flush()
	return err
}
func (s SumReplyCapn) MarshalJSON() ([]byte, error) {
	b := bytes.Buffer{}
	err := s.WriteJSON(&b)
	return b.Bytes(), err
}
func (s SumReplyCapn) WriteCapLit(w io.Writer) error {
	b := bufio.NewWriter(w)
	var err error
	var buf []byte
	_ = buf
	err = b.WriteByte('(')
	if err != nil {
		return err
	}
	_, err = b.WriteString("retCode = ")
	if err != nil {
		return err
	}
	{
		s := s.RetCode()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	_, err = b.WriteString(", ")
	if err != nil {
		return err
	}
	_, err = b.WriteString("val = ")
	if err != nil {
		return err
	}
	{
		s := s.Val()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	err = b.WriteByte(')')
	if err != nil {
		return err
	}
	err = b.Flush()
	return err
}
func (s SumReplyCapn) MarshalCapLit() ([]byte, error) {
	b := bytes.Buffer{}
	err := s.WriteCapLit(&b)
	return b.Bytes(), err
}

type SumReplyCapn_List C.PointerList

func NewSumReplyCapnList(s *C.Segment, sz int) SumReplyCapn_List {
	return SumReplyCapn_List(s.NewCompositeList(16, 0, sz))
}
func (s SumReplyCapn_List) Len() int { return C.PointerList(s).Len() }
func (s SumReplyCapn_List) At(i int) SumReplyCapn {
	return SumReplyCapn(C.PointerList(s).At(i).ToStruct())
}
func (s SumReplyCapn_List) ToArray() []SumReplyCapn {
	n := s.Len()
	a := make([]SumReplyCapn, n)
	for i := 0; i < n; i++ {
		a[i] = s.At(i)
	}
	return a
}
func (s SumReplyCapn_List) Set(i int, item SumReplyCapn) { C.PointerList(s).Set(i, C.Object(item)) }

type SumRequestCapn C.Struct

func NewSumRequestCapn(s *C.Segment) SumRequestCapn      { return SumRequestCapn(s.NewStruct(24, 0)) }
func NewRootSumRequestCapn(s *C.Segment) SumRequestCapn  { return SumRequestCapn(s.NewRootStruct(24, 0)) }
func AutoNewSumRequestCapn(s *C.Segment) SumRequestCapn  { return SumRequestCapn(s.NewStructAR(24, 0)) }
func ReadRootSumRequestCapn(s *C.Segment) SumRequestCapn { return SumRequestCapn(s.Root(0).ToStruct()) }
func (s SumRequestCapn) UserId() int64                   { return int64(C.Struct(s).Get64(0)) }
func (s SumRequestCapn) SetUserId(v int64)               { C.Struct(s).Set64(0, uint64(v)) }
func (s SumRequestCapn) Num1() int64                     { return int64(C.Struct(s).Get64(8)) }
func (s SumRequestCapn) SetNum1(v int64)                 { C.Struct(s).Set64(8, uint64(v)) }
func (s SumRequestCapn) Num2() int64                     { return int64(C.Struct(s).Get64(16)) }
func (s SumRequestCapn) SetNum2(v int64)                 { C.Struct(s).Set64(16, uint64(v)) }
func (s SumRequestCapn) WriteJSON(w io.Writer) error {
	b := bufio.NewWriter(w)
	var err error
	var buf []byte
	_ = buf
	err = b.WriteByte('{')
	if err != nil {
		return err
	}
	_, err = b.WriteString("\"userId\":")
	if err != nil {
		return err
	}
	{
		s := s.UserId()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	err = b.WriteByte(',')
	if err != nil {
		return err
	}
	_, err = b.WriteString("\"num1\":")
	if err != nil {
		return err
	}
	{
		s := s.Num1()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	err = b.WriteByte(',')
	if err != nil {
		return err
	}
	_, err = b.WriteString("\"num2\":")
	if err != nil {
		return err
	}
	{
		s := s.Num2()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	err = b.WriteByte('}')
	if err != nil {
		return err
	}
	err = b.Flush()
	return err
}
func (s SumRequestCapn) MarshalJSON() ([]byte, error) {
	b := bytes.Buffer{}
	err := s.WriteJSON(&b)
	return b.Bytes(), err
}
func (s SumRequestCapn) WriteCapLit(w io.Writer) error {
	b := bufio.NewWriter(w)
	var err error
	var buf []byte
	_ = buf
	err = b.WriteByte('(')
	if err != nil {
		return err
	}
	_, err = b.WriteString("userId = ")
	if err != nil {
		return err
	}
	{
		s := s.UserId()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	_, err = b.WriteString(", ")
	if err != nil {
		return err
	}
	_, err = b.WriteString("num1 = ")
	if err != nil {
		return err
	}
	{
		s := s.Num1()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	_, err = b.WriteString(", ")
	if err != nil {
		return err
	}
	_, err = b.WriteString("num2 = ")
	if err != nil {
		return err
	}
	{
		s := s.Num2()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	err = b.WriteByte(')')
	if err != nil {
		return err
	}
	err = b.Flush()
	return err
}
func (s SumRequestCapn) MarshalCapLit() ([]byte, error) {
	b := bytes.Buffer{}
	err := s.WriteCapLit(&b)
	return b.Bytes(), err
}

type SumRequestCapn_List C.PointerList

func NewSumRequestCapnList(s *C.Segment, sz int) SumRequestCapn_List {
	return SumRequestCapn_List(s.NewCompositeList(24, 0, sz))
}
func (s SumRequestCapn_List) Len() int { return C.PointerList(s).Len() }
func (s SumRequestCapn_List) At(i int) SumRequestCapn {
	return SumRequestCapn(C.PointerList(s).At(i).ToStruct())
}
func (s SumRequestCapn_List) ToArray() []SumRequestCapn {
	n := s.Len()
	a := make([]SumRequestCapn, n)
	for i := 0; i < n; i++ {
		a[i] = s.At(i)
	}
	return a
}
func (s SumRequestCapn_List) Set(i int, item SumRequestCapn) { C.PointerList(s).Set(i, C.Object(item)) }
