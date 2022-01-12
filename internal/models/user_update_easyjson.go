// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonFcc75609DecodeForumInternalModels(in *jlexer.Lexer, out *UserUpdate) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "fullname":
			out.FullName = string(in.String())
		case "about":
			out.About = string(in.String())
		case "email":
			out.Email = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonFcc75609EncodeForumInternalModels(out *jwriter.Writer, in UserUpdate) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"fullname\":"
		out.RawString(prefix[1:])
		out.String(string(in.FullName))
	}
	{
		const prefix string = ",\"about\":"
		out.RawString(prefix)
		out.String(string(in.About))
	}
	{
		const prefix string = ",\"email\":"
		out.RawString(prefix)
		out.String(string(in.Email))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UserUpdate) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonFcc75609EncodeForumInternalModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserUpdate) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonFcc75609EncodeForumInternalModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UserUpdate) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonFcc75609DecodeForumInternalModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserUpdate) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonFcc75609DecodeForumInternalModels(l, v)
}
