package goson

import (
    "testing"
)

const (
	emptyObjectStringParse = "{}"
	emptyArrayStringParse = "{\"Array\" : []}"
		
)

func TestEmptyObjectParse(t *testing.T) {
	resultObj, err := ParseString(emptyObjectStringParse)		
	if err != nil {
		t.Logf("TestEmptyObjectParse failed with error message: %s", err.Error())
		t.FailNow()
	}
	
	if resultObj.GetBaseString() != emptyObjectStringParse {
		t.Logf("TestEmptyObjectParse failed with mismatching string: %s", resultObj.GetBaseString())
		t.Fail()
	}
}