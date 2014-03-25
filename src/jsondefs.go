package goson

import (
	"strings"
	"strconv"
)

const (
	JSONValueString = iota
	JSONValueNumber = iota
	JSONValueObject = iota
	JSONValueArray = iota
	JSONValueBoolean = iota
	JSONValueNull = iota
)

const (
	JSONErrorMalformedString = iota
	JSONErrorProgramFlowIssue = iota
)

type JSONNode interface {
	GetBaseString() string
	GetType() int
}

type JSONObject struct {
	Members map[string]JSONNode
}

type JSONArray struct {
	Elements *[]JSONNode
}

type JSONString struct {
	Value string
}

type JSONNumber struct {
	Value float64
}

type JSONBoolean struct {
	Value bool
}

type JSONNull struct {
}

type JSONError struct {
	errString string
	errCode int
}

func (jnull JSONNull) GetBaseString() string {
	return "null"
}

func (jnull JSONNull) GetType() int {
	return JSONValueNull
}

func (jbool JSONBoolean) GetBaseString() string {
	if jbool.Value {
		return "true"
	} else {
		return "false"
	}
}

func (jbool JSONBoolean) GetType() int {
	return JSONValueBoolean
}

func (jnum JSONNumber) GetBaseString() string {
	return strconv.FormatFloat(jnum.Value, 'e', -1, 64)
}

func (jnum JSONNumber) GetType() int {
	return JSONValueNumber
}

func (jstr JSONString) GetBaseString() string {
	return "\"" + jstr.Value + "\""
}

func (jstr JSONString) GetType() int {
	return JSONValueString
}

func (jarr JSONArray) GetBaseString() string {
	retVal := "["
	for _, value := range *jarr.Elements {
		retVal += value.GetBaseString() + ", "	
	}
	retVal = strings.TrimRight(retVal, ", ")
	retVal += "]"
	return retVal
}

func (jarr JSONArray) GetType() int {
	return JSONValueArray
}

func (jobj JSONObject) GetBaseString() string {
	retVal := "{"
	for name, value := range jobj.Members {
			retVal += "\"" + name + "\" : " + value.GetBaseString() + ", "
	}
	retVal = strings.TrimRight(retVal, ", ")
	retVal += "}"
	return retVal
}

func (jobj JSONObject) GetType() int {
	return JSONValueObject
}

func (jerr JSONError) Error() string {
	return jerr.errString + " JSON Error Code: " + strconv.FormatInt(int64(jerr.errCode), 10)
}