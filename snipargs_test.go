package main

import "testing"

//import "fmt"

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

func benchGetArguments(str string, b *testing.B) {

	for n := 0; n < b.N; n++ {
		_ = getArguments(str)
	}

}

func BenchmarkGetArgumentsMulti(b *testing.B) {
	benchGetArguments("test spaces {:value 3!with spaces:} and with ip {:ip!:6.6.6.6:}", b)
}

func BenchmarkGetArgumentsSingle(b *testing.B) {
	benchGetArguments("test spaces {:value 3!with spaces:} and with ip                ", b)
}

func BenchmarkGetArgumentsNone(b *testing.B) {
	benchGetArguments("test spaces   value 3 with spaces   and with ip  ip  6.6.6.123 ", b)
}

func BenchmarkGetArgumentsIP(b *testing.B) {
	benchGetArguments("ping -t5 {:ip!8.8.8.8:}", b)
}

func TestGetArgumentList(t *testing.T) {
	in := "ping {:ip!8.8.8.8:} -t 5"
	expect := "{:ip!8.8.8.8:}"

	out, ok := getArgumentList(in)
	if ok != true {
		t.Errorf("getArgumentsList returned false")
	}
	if out[0] != expect {
		t.Errorf("getArguments return incorrect arg default, got: %s, want: %s.", out[0], expect)
	}
}

func BenchmarkGetArgumentList(b *testing.B) {
	in := "ping {:ip!8.8.8.8:}"

	for n := 0; n < b.N; n++ {
		_, _ = getArgumentList(in)
	}

}

func TestArgumentReplace(t *testing.T) {
	//argumentReplace(vars []snipArgs, code string) string
	//setup snipArgs
	args := []snipArgs{
		snipArgs{
			Name:  "ip",
			Value: "1.2.3.4",
		},
	}

	expected := "ping -t5 1.2.3.4"
	in := "ping -t5 {:ip!8.8.8.8:}"
	out := argumentReplace(args, in)
	if out != expected {
		t.Errorf("argumentReplace return incorrect value, got: %s, want: %s.", out, expected)
	}

}

func BenchmarkArgumentReplace(b *testing.B) {
	args := []snipArgs{
		snipArgs{
			Name:  "ip",
			Value: "1.2.3.4",
		},
	}

	in := "ping -t5 {:ip!8.8.8.8:}"
	for n := 0; n < b.N; n++ {
		_ = argumentReplace(args, in)
	}

}
