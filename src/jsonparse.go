package goson

import (
	"unicode/utf8"
	"strconv"
)

func ParseString(baseString string) (obj JSONObject, err error) {
		
	returnedObject, _, JSONError := parseObject(baseString)
	
	if err != nil {	
		err = error(JSONError)
	}
	
	return returnedObject, err
}

/*
 *	All functions from here on down assume that they are passed a string starting
 *	with the 'flag' character -- that is to say, the string function, which is 
 *	flagged to be called by \", expects its parameter to start with that \"
 * 
 *  BNF for an Object is:  Object := { Members }; {}
 *  BNF for Members is:  Members := Pair; Pair, Members
 */
func parseObject(baseString string) (object JSONObject, length int, err *JSONError) {
	//Everything's a value.
	baseLength := len(baseString)
	members := make(map[string]JSONNode)
	currentIndex := 0
	var name string
	var value JSONNode
	var memberlen int
	firstTimeThrough := true
	
	//Start with '{'
	currentRune, runeSize := utf8.DecodeRuneInString(baseString[currentIndex:])
	if currentRune != '{' {
		return JSONObject{nil}, -1, &JSONError{"Object does not begin with '}'", JSONErrorMalformedString}
	}
	currentIndex += runeSize
		
	for currentIndex < baseLength {
		addWhitespace, err := getDistanceToNextNonWhitespace(baseString[currentIndex:])
		
		if err != nil {
			return JSONObject{nil}, -1, err	
		}
		currentIndex += addWhitespace	
		
		currentRune, runeSize := utf8.DecodeRuneInString(baseString[currentIndex:])
		//Return if end of object.  Only check for '}' here on first time through.  Laaaaazy code.
        if currentRune == '}' {
        	if firstTimeThrough {
        		return JSONObject{members}, (currentIndex + runeSize), nil
        	} else {
				return JSONObject{nil}, -1, &JSONError{"Object terminates with ', }'", JSONErrorMalformedString}
        	}
        }
		firstTimeThrough = false
        
		
		//Parse value and add to object
		name, value, memberlen, err = parsePair(baseString[currentIndex:])
		
		if err != nil {
			return JSONObject{nil}, -1, err
		}
			
		members[name] = value
		
		currentIndex += memberlen
		addWhitespace, err = getDistanceToNextNonWhitespace(baseString[currentIndex:])
		
		if err!= nil {
			return JSONObject{nil}, -1, err
		}
		currentIndex += addWhitespace
		
		//End of object or a comma.
		currentRune, runeSize = utf8.DecodeRuneInString(baseString[currentIndex:])
		if currentRune == '}' {
			return JSONObject{members}, (currentIndex + runeSize), nil
		} else if currentRune == ',' {
			currentIndex += runeSize
		} else {
			return JSONObject{nil}, -1, &JSONError{"JSON object has 2 members unseparated by comma", JSONErrorMalformedString} 
		}	
	}
	return JSONObject{nil}, -1, &JSONError{"JSON object does not terminate with '}'", JSONErrorMalformedString}
}

//Starts at the \" of the pair's name.  BNF for a Pair is: Pair := String : Value
func parsePair(baseString string) (name string, value JSONNode, length int, err *JSONError) {
	var currentLength int
	var currentRune rune
	var whitespaceLength int
	var valueLength int
	var nameLength int
	var nameObj JSONString

	nameObj, nameLength, err = parseString(baseString)
	name = nameObj.Value
		
	if err != nil {
		return name, JSONNull{}, -1, err
	}
	
	currentLength = nameLength
	whitespaceLength, err = getDistanceToNextNonWhitespace(baseString[currentLength:])
	
	if err != nil {
		return "", JSONNull{}, -1, err
	}
	
	currentLength += whitespaceLength
		
	currentRune, _ = utf8.DecodeRuneInString(baseString[currentLength:])
	
	if currentRune != ':' {
		return "", JSONNull{}, -1, &JSONError{"Member name and value not separated.", JSONErrorMalformedString}
	}
		
	value, valueLength, err = parseValue(baseString[currentLength:])
	
	if err != nil {
		return "", JSONNull{}, -1, err
	}
	
	currentLength += valueLength
	return name, value, currentLength, nil
}

// Passes value on to relevant parse function
func parseValue(baseString string) (value JSONNode, length int, err *JSONError) {
	currentRune, _ := utf8.DecodeRuneInString(baseString)
	switch currentRune {
		case '[' :
			return parseArray(baseString)
		case '{' :
			return parseObject(baseString)
		case '"' :
			return parseString(baseString)
		case 't':
			fallthrough
		case 'f':
			return parseBoolean(baseString)
		case 'n':
			return parseNull(baseString)
		case '1' :
			fallthrough
		case '2' :
			fallthrough
		case '3' :
			fallthrough
		case '4' :
			fallthrough
		case '5' :
			fallthrough
		case '6' :
			fallthrough
		case '7' :
			fallthrough
		case '8' :
			fallthrough
		case '9' :
			fallthrough
		case '0' :
			fallthrough
		case '-' :
			return parseNumber(baseString)
		default :
			return JSONNull{}, 0, &JSONError{"Value in JSON does not parse properly", JSONErrorMalformedString}
	}
}

func parseArray(baseString string) (value JSONArray, length int, err *JSONError) {
	currentRune, runeLength := utf8.DecodeRuneInString(baseString)
	baseLength := len(baseString)
	values := []JSONNode{}
	var whitespaceLen int
	
	if currentRune != '[' {
		return JSONArray{&values}, 0, &JSONError{"Array does not start with '['", JSONErrorMalformedString}
	}
	
	currentLength := runeLength
	needsComma, needsValue := false, false
	
	for currentLength < baseLength {
		whitespaceLen, err = getDistanceToNextNonWhitespace(baseString[currentLength:])
		if err != nil {
			return JSONArray{new([]JSONNode)}, 0, err
		}
		
		currentLength += whitespaceLen
		
		currentRune, runeLength := utf8.DecodeRuneInString(baseString[currentLength:])
		
		if currentRune == ']' {
			if needsValue {
				return JSONArray{new([]JSONNode)}, 0, &JSONError{"Array ends with comma, no value", JSONErrorMalformedString}
			} else {
				return JSONArray{&values}, currentLength + runeLength, nil
			}	
		} else if currentRune == ',' {
			if needsValue {
				return JSONArray{new([]JSONNode)}, 0, &JSONError{"Array has two commas in a row", JSONErrorMalformedString}
			} else {
				needsComma = false
				needsValue = true
			}
			currentLength += runeLength
		} else {
			if needsComma {
				return JSONArray{new([]JSONNode)}, 0, &JSONError{"Array has two values not separated by comma", JSONErrorMalformedString}
			}
			
			value, valueLength, err := parseValue(baseString[currentLength:])
		
			if err != nil {
				return JSONArray{new([]JSONNode)}, 0, err
			}
			
			needsComma = true
			needsValue = false
	
			values = append(values, value)
			currentLength += valueLength
		}
	}
	return JSONArray{new([]JSONNode)}, 0, &JSONError{"Array does not end before end of string", JSONErrorMalformedString}
}

// Starts at what SHOULD be the start of the string -- checks for \"
func parseString(baseString string) (value JSONString, length int, err *JSONError) {
	currentRune, runeLength := utf8.DecodeRuneInString(baseString)
	baseLength := len(baseString)
	if currentRune != '"' {
		return JSONString{""}, 0, &JSONError{"String value does not begin with quotation mark (how did you even get here, then?)", JSONErrorMalformedString}
	}
	
	currentLength := runeLength
	currentRune, runeLength = utf8.DecodeRuneInString(baseString[currentLength:])
	for currentLength < baseLength {
		currentLength += runeLength
		if currentRune == '"' {
			return JSONString{baseString[:currentLength]}, currentLength, nil
		}	
	}
		
	return JSONString{""}, 0, &JSONError{"String value never ends.", JSONErrorMalformedString}	
}

// Starts at what SHOULD be the start of the null
func parseNull(baseString string) (value JSONNull, length int, err *JSONError) {
	if baseString[:4] != "null" {
		return JSONNull{}, 0, &JSONError{"Value for NULL does not have 4 characters == 'null'", JSONErrorMalformedString}
	}
	return JSONNull{}, 4, nil
}

// Starts at what SHOULD be the start of the string -- checks for \"
func parseBoolean(baseString string) (value JSONBoolean, length int, err *JSONError) {
	firstRune, _ := utf8.DecodeRuneInString(baseString)
	if firstRune == 't' {
		if baseString[:4] != "true" {
			return JSONBoolean{true}, 0, &JSONError{"Value for JSONBoolean True does not equal 'true'", JSONErrorMalformedString}
		}	
		return JSONBoolean{true}, 4, nil
	} else if firstRune == 'f' {
		if baseString[:5] != "false" {
			return JSONBoolean{false}, 0, &JSONError{"Value for JSONBoolean False does not equal 'false'", JSONErrorMalformedString}
		}
		return JSONBoolean{false}, 5, nil
	} else {
		return JSONBoolean{false}, 0, &JSONError{"Value passed to parseBoolean starts with neither 't' or 'f'", JSONErrorProgramFlowIssue}
	}
}

// Starts at what SHOULD be the start of the string -- checks for \"
func parseNumber(baseString string) (value JSONNumber, length int, err *JSONError) {
	currentLength := 0
	hasExponent, firstCharacter, hasDecimal := false, true, false
	
	currentRune, runeLength := utf8.DecodeRuneInString(baseString)
		
	for currentRune != ' ' && currentRune != ',' {
		currentLength += runeLength
		
		valid, isDecimal, isExponent, isNegativeSign := checkValidNumberCharacter(currentRune)
	
		if !valid {
			return JSONNumber{0}, 0, &JSONError{"Invalid character found in Number", JSONErrorMalformedString}
		}
		if isDecimal && hasDecimal {
			return JSONNumber{0}, 0, &JSONError{"Second Decimal found in Number", JSONErrorMalformedString}
		}
		if isExponent && hasExponent {
			return JSONNumber{0}, 0, &JSONError{"Second Exponent 'e' or 'E' found in Number", JSONErrorMalformedString}	
		}
		if isExponent && firstCharacter {
			return JSONNumber{0}, 0, &JSONError{"Number leads with 'e' or 'E'", JSONErrorMalformedString}
		}
		if isNegativeSign && !firstCharacter {
			return JSONNumber{0}, 0, &JSONError{"Negative Sign found in middle of number", JSONErrorMalformedString}	
		}
		
		if isExponent {		
			hasExponent = true
			firstCharacter = true
			hasDecimal = false
		} else {
			firstCharacter = false
		}
		
		if isDecimal {
			hasDecimal = true
		}
		
		currentRune, runeLength = utf8.DecodeRuneInString(baseString[currentLength:])
	}
	
	resultFloat, parseErr := strconv.ParseFloat(baseString[:currentLength], 64)
	
	if parseErr != nil {
		return JSONNumber{0}, 0, &JSONError{"Error parsing float after initial pass", JSONErrorMalformedString}
	}
		
	
	return JSONNumber{resultFloat}, currentLength, nil	
}

func checkValidNumberCharacter(toCheck rune) (valid bool, isDecimal bool, isExponent bool, isNegativeSign bool) {
	switch toCheck {
		case '1' :
			fallthrough
		case '2' :
			fallthrough
		case '3' :
			fallthrough
		case '4' :
			fallthrough
		case '5' :
			fallthrough
		case '6' :
			fallthrough
		case '7' :
			fallthrough
		case '8' :
			fallthrough
		case '9' :
			fallthrough
		case '0' :
			valid, isDecimal, isExponent, isNegativeSign = true, false, false, false
		case '-' :
			valid, isDecimal, isExponent, isNegativeSign = true, false, false, true
		case 'e' :
			fallthrough
		case 'E' :
			valid, isDecimal, isExponent, isNegativeSign = true, false, true, false
		case '.' :
			valid, isDecimal, isExponent, isNegativeSign = true, true, false, false
		default :
			valid, isDecimal, isExponent, isNegativeSign = false, false, false, false
	}
	return
}

func getDistanceToNextNonWhitespace(baseString string) (int, *JSONError) {
	currentIndex := 0
	strlen := len(baseString)
	currentRune, runeSize := utf8.DecodeRuneInString(baseString[currentIndex:])
	
	for currentRune == ' ' {
		currentIndex += runeSize
		if(currentIndex > strlen) {
			return -1, &JSONError{"JSON Object String does not end before whitespace ends", JSONErrorMalformedString}
		}
		currentRune, runeSize = utf8.DecodeRuneInString(baseString[currentIndex:])
	}		
	return currentIndex, nil
}
