package main

import "testing"

//import "fmt"

func TestValidCmdType(t *testing.T) {

	out, sep := validCmdType(" bASh ")
	if out != "bash" {
		t.Errorf("validCmdType string return incorrect value, got: %s, want: %s.", out, "bash")
	}
	if sep != " \\" {
		t.Errorf("validCmdType seperator return incorrect value, got: %s, want: %s.", out, " \\")
	}
}

func TestCleanString(t *testing.T) {
	str := "Test.String_with-dash"
	strName := cleanString(str, "[^A-Za-z_.-]")
	if strName != str {
		t.Errorf("cleanString string return incorrect value, got: %s, want: %s.", strName, str)
	}

}

func TestGetArguments(t *testing.T) {
	//single arg
	values := []struct {
		cmd      string
		argName  string
		argValue string
	}{
		{"test string", "value 1", "default"},
		{"test values", "value 2", "69"},
		{"test spaces", "value 3", "with spaces"},
		{"test ip addr", "value4", "1.2.3.4"},
		{"test no default", "value5", ""},
	}

	//multi arg
	multivalue := "multi command argument {:arg1!val1:} with {:arg 2!1 answer:} and without {:arg3:} default values and ip addresses {:ip!123.234.34.4:}."
	multians := []struct {
		argName  string
		argValue string
	}{
		{"arg1", "val1"},
		{"arg 2", "1 answer"},
		{"arg3", ""},
		{"ip", "123.234.34.4"},
	}

	//test single argument string
	for _, item := range values {
		var v string = ""
		if len(item.argValue) > 0 {
			v = "!" + item.argValue
		}
		str := item.cmd + " {:" + item.argName + v + ":}"
		out := getArguments(str)
		if out[0].Name != item.argName {
			t.Errorf("getArguments return incorrect arg Name, got: %s, want: %s.", out[0].Name, item.argName)
		}
		if out[0].Value != item.argValue {
			t.Errorf("getArguments return incorrect arg default, got: %s, want: %s.", out[0].Value, item.argValue)
		}
	}

	//now check multi argument string
	out := getArguments(multivalue)
	//fmt.Println(multivalue)
	//fmt.Println(out)
	for i, item := range multians {
		if out[i].Name != item.argName {
			t.Errorf("getArguments return incorrect arg Name, got: %s, want: %s.", out[i].Name, item.argName)
		}
		if out[i].Value != item.argValue {
			t.Errorf("getArguments return incorrect arg default, got: %s, want: %s.", out[i].Value, item.argValue)
		}
	}

}
