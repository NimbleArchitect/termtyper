package main

import (
	"log"
	"regexp"
	"strings"
)

func cleanString(data string, regex string) string {
	logDebug("F:cleanString:start")
	logDebug("F:cleanString:data =", data)

	reg, err := regexp.Compile(regex)
	if err != nil {
		log.Fatal(err)
	}
	newstr := reg.ReplaceAllString(data, "")
	logDebug("F:cleanString:return -", newstr)
	return newstr
}

// returns a list of start and end positions of each argument, found in text
func getArgumentPos(text string) ([][]int, bool) {
	var ok bool
	var matches [][]int
	logDebug("F:getArgumentPos:start")
	if len(text) > 0 {
		regexstring := regexp.MustCompile(regexMatch)
		matches = regexstring.FindAllStringIndex(text, -1)
		ok = true
	} else {
		ok = false
	}
	if len(matches) <= 0 {
		ok = false
	}
	logDebug("F:getArgumentPos:return =", matches, ",", ok)
	return matches, ok
}

// return string array of arguments found in text
func getArgumentList(text string) ([]string, bool) {
	var ok bool
	var matches []string
	logDebug("F:getArgumentList:start =", text)

	if len(text) > 0 {
		regexstring := regexp.MustCompile(regexMatch)
		matches = regexstring.FindAllString(text, -1)
		ok = true
	} else {
		ok = false
	}
	if len(matches) <= 0 {
		ok = false
	}
	logDebug("F:getArgumentList:return =", matches, ",", ok)
	return matches, ok
}

//search text looking for arguments returns array of snipArgs
func getArguments(text string) []snipArgs {
	var namelist []snipArgs
	var varlist []string
	var varitem snipArgs

	logDebug("F:getArguments:start")
	varlist, ok := getArgumentList(text)
	if ok == false { //no arguments found in text
		logDebug("F:getArguments:return =", namelist)
		return namelist
	}

	for _, varpos := range varlist {
		//var pos is start and end locations in array
		vars := strings.Split(varpos, ":")     //arguments are enclosed in : so we remove those first
		varname := strings.Split(vars[1], "!") // default values for arguments can be found after !
		logDebug("F:getArguments:varname =", varname)
		logDebug("F:getArguments:len(varname) =", len(varname))
		strName := cleanString(varname[0], "[^A-Za-z0-9_. -]") //remove invalid chars from name
		namelen := len(varname)
		if namelen == 1 { // ! is optional so check if argument dosent have a default value
			varitem = snipArgs{
				Name:  strings.TrimSpace(strName),
				Value: "",
			}
		} else if namelen == 2 { //argument has a default value
			varitem = snipArgs{
				Name:  strings.TrimSpace(strName),
				Value: strings.TrimSpace(varname[1]),
			}
		} else {
			// multiple defaults values have been suppilied so write warning
			logWarn("multipule default values detected.")
		}

		namelist = append(namelist, varitem)
	}
	logDebug("F:getArguments:return =", namelist)
	return namelist
}

//search code look ing for arguments, replace with values from snipArgs
func argumentReplace(vars []snipArgs, code string) string {
	var newcode string
	var val string
	logDebug("F:argumentReplace:start")
	if len(code) <= 0 {
		return ""
	}
	itmarg := getArguments(code)      // get array of arguments from code
	argPos, _ := getArgumentPos(code) // get array or argument start/end positions
	//spin through all arguments and replace variables as needed
	itmlen := len(itmarg)
	varlen := len(vars)
	if varlen < 0 {
		varlen = 0
	}

	if varlen != itmlen { // itmlen is not the same length as varlen
		emptyarg := snipArgs{Name: "", Value: ""}
		for c := varlen; c <= itmlen; c++ {
			vars = append(vars, emptyarg) //so add enough empty values to our vars array this helps the for loop below
		}
	}
	newcode = code
	logDebug("V:itmlen =", itmlen)
	logDebug("V:varlen =", len(vars))

	startlen := itmlen - 1
	for i := (startlen); i >= 0; i-- {
		logDebug("V:i =", i)
		itm := itmarg[i]
		//logDebug("F:argumentReplace:vars[i] =", vars[i])
		logDebug("F:argumentReplace:itm =", itm)

		//make sure the incomming argument name matches the variable in the code
		if itm.Name != vars[i].Name {
			if vars[i].Name != "" { //vars name could of been added above so can be empty
				logDebug("F:argumentReplace:return = \"\"")
				return "" // refuse and exit function if name wasn't empty
			}
		}

		if len(vars[i].Value) > 0 {
			val = vars[i].Value //incomming value is valid so use that
		} else if len(itm.Value) > 0 {
			val = itm.Value //incoming value not valid but we have a default value so use it
		} else {
			val = "{" + itm.Name + "}" //nothing is valid so we default to the name in braces
		}

		itmpos := argPos[i] //start and end pos of txt to replace
		s := itmpos[0]
		e := itmpos[1]

		newcode = newcode[:s] + val + newcode[e:]
	}
	logDebug("F:argumentReplace:return =", newcode)
	return newcode
}
