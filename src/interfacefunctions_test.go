package goson

import (
    "testing"
)

const (
	emptyObjectString = "{}"
	emptyArrayString = "{\"Array\" : []}"
)

func TestEmptyObjectBaseString(t *testing.T) {
	objMap := make(map[string]JSONNode)
	testObj := JSONObject{objMap}
	if(testObj.GetBaseString() != emptyObjectString) {
		t.Logf("Empty object test failed with the following output: %s", testObj.GetBaseString())
		t.Fail()
	}
}

func TestEmptyArrayBaseString(t *testing.T) {
	objMap := make(map[string]JSONNode)
	testObj := JSONObject{objMap}
	
	objArr := new([]JSONNode)
	testArr := JSONArray{*objArr}
	
	testObj.Members["Array"] = testArr
	
	if(testObj.GetBaseString() != emptyArrayString) {
		t.Logf("Empty array test failed with the following output:  %s", testObj.GetBaseString())
		t.Fail()
	}
}