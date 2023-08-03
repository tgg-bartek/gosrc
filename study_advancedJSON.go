// Advanced JSON Handling in Go
// https://youtu.be/vsN11YAEJHY

package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/davecgh/go-spew/spew"
)

func main2() {
	// --- Marshaling JSON
	// Marshal returns the JSON encoding of v
	//marshalJson1()
	//marshalJson2()

	// --- Unmarshalling JSON
	// Unmarshal parses the JSON-encoded data and stores the result
	// in the value pointed to by v
	//unmarshalJson1()
	//unmarshalJson2()

	// --- Unmarshalling JSON: Int OR String
	//unmarshalIntOrString()
	//unmarshalSliceOrString()

	// ---- Unmarshalling JSON: Unknown types
	//unmarshalUnknownTypes()
	unmarshalUnknownTypes2()

}

// ---- Dealing with cases when we do NOT know data types before
// we unmarshall it.
// Cases on "unknown input"
// a. int vs string: 32 vs "32"
// b. object or array of objects: {...} vs [{...}, {...}]
// c. input may be success or an error: {"success": true, "results": [...]} vs {"success": false, "error": "..."}

/////////////////////////////////////////////////
// --- Dealing with case C: Unknown types

// Example #1 : No fields overlap
// {"results": [...]}
// {"error": "not found", "reason": "The requested object..."}

// First, create our distinct success and error types

type Success struct {
	Results []string `json:"results"`
}

type Error struct {
	Error  string `json:"error"`
	Reason string `json:"reason"`
}

// In this simple example, because no fields are shared between the types,
// we can simply embed both types in a wrapper.

type Response struct {
	Success
	Error
}

/*
([]main.Response) (len=2 cap=4) {
 (main.Response) {
  Success: (main.Success) {
   Results: ([]string) (len=2 cap=4) {
    (string) (len=3) "one",
    (string) (len=3) "two"
   }
  },
  Error: (main.Error) {
   Error: (string) "",
   Reason: (string) ""
  }
 },
 (main.Response) {
  Success: (main.Success) {
   Results: ([]string) <nil>
  },
  Error: (main.Error) {
   Error: (string) (len=9) "not found",
   Reason: (string) (len=35) "The requested object does not exist"
  }
 }
}
*/

func unmarshalUnknownTypes() {
	data := []byte(`[
		{"results": ["one", "two"]},
		{"error": "not found", "reason": "The requested object does not exist"}
	]`)
	x := make([]Response, 0)
	if err := json.Unmarshal(data, &x); err != nil {
		panic(err)
	}
	spew.Dump(x)
}

// Example #2. Conflicting / Overlapping fields
// Success: {"status": "ok", "results": [...]}
// Failure: {"status": "not found", "reason": "The requested..."}

type Success2 struct {
	Status  string   `json:"status"`
	Results []string `json:"results"`
}

type Error2 struct {
	Status string `json:"status"`
	Reason string `json:"reason"`
}

type Response2 struct {
	Success2
	Error2
}

// Because fields overlap we need to create custom unmarshall func
// or else the unmarshal does not know where the status field goes

func (r *Response2) UnmarshalJSON(d []byte) error {
	// First assume success
	if err := json.Unmarshal(d, &r.Success2); err != nil {
		return err
	}
	if r.Success2.Status == "ok" {
		return nil
	}
	// if failure make success to its zero value (empty success)
	r.Success2 = Success2{}
	return json.Unmarshal(d, &r.Error2)
}

/*
([]main.Response) (len=2 cap=4) {
 (main.Response) {
  Success: (main.Success) {
   Results: ([]string) (len=2 cap=4) {
    (string) (len=3) "one",
    (string) (len=3) "two"
   }
  },
  Error: (main.Error) {
   Error: (string) "",
   Reason: (string) ""
  }
 },
 (main.Response) {
  Success: (main.Success) {
   Results: ([]string) <nil>
  },
  Error: (main.Error) {
   Error: (string) "",
   Reason: (string) (len=35) "The requested object does not exist"
  }
 }
}

*/
func unmarshalUnknownTypes2() {
	data := []byte(`[
		{"status": "ok", "results": ["one", "two"]},
		{"status": "not found", "reason": "The requested object does not exist"}
	]`)
	x := make([]Response, 0)
	if err := json.Unmarshal(data, &x); err != nil {
		panic(err)
	}
	spew.Dump(x)
}

/////////////////////////////////////////////////
// --- Dealing with case B: Array or single element

// SliceOrString is a slice of strings that may be unmarshalled from either
// a JSON array of strings, or a single JSON string

type SliceOrString []string

func (s *SliceOrString) UnmarshalJSON(d []byte) error {
	if d[0] == '"' {
		var v string
		err := json.Unmarshal(d, &v)
		fmt.Println("String ->", v)
		*s = SliceOrString{v} // makes v an elem of slice SliceOrString (not to confuse with type conversion ())
		return err
	}
	var v []string
	err := json.Unmarshal(d, &v)
	*s = v // or SliceOrString(v) <- () converts v to SliceOrString but is not needed here
	fmt.Println("Array ->", v)
	return err
}

func unmarshalSliceOrString() {
	data := []byte(`["one", ["two", "elements"]]`)
	x := make([]SliceOrString, 0)
	if err := json.Unmarshal(data, &x); err != nil {
		panic(err)
	}
	fmt.Println(x) // eg x[0]
}

/////////////////////////////////////////////////
// --- Dealing with case A: Number vs String

// IntOrString is an int that may be unmarshalled from either a JSON number
// literal, or a JSON string.

type IntOrString int

// You create custom unmarshaler by defining "UnmarshallJSON" func

func (i *IntOrString) UnmarshalJSON(d []byte) error {
	var v int
	err := json.Unmarshal(bytes.Trim(d, `"`), &v)
	*i = IntOrString(v)
	return err
}

func unmarshalIntOrString() {
	data := []byte(`[123, "321"]`)
	x := make([]IntOrString, 0)
	if err := json.Unmarshal(data, &x); err != nil {
		panic(err)
	}
	fmt.Println(x)
}

/////////////////////////////////////////////////

func unmarshalJson2() {
	type person struct {
		Name        string `json:"name"`
		Age         int    `json:"age"`
		Description string `json:"desc,omitempty"`
		secret      string // Unexported fields are never (un)marshalled
	}
	data := []byte(`{"name": "Jeo", "age": 32, "desc": "Python Programmer"}`)
	var x person
	_ = json.Unmarshal(data, &x)
	//fmt.Println(reflect.TypeOf(x))
	fmt.Println(x)

}

func unmarshalJson1() {
	data := []byte(`{"foo": "bar"}`)
	var x interface{}
	_ = json.Unmarshal(data, &x)
	dict := x.(map[string]interface{})
	fmt.Println(dict["foo"])

}

func marshalJson2() {
	type person struct {
		Name        string `json:"name"`
		Age         int    `json:"age"`
		Description string `json:"desc,omitempty"`
		secret      string // Unexported fields are never (un)marshaled
	}
	x := person{
		Name:   "Bob",
		Age:    32,
		secret: "Shhh!",
	}
	data, _ := json.Marshal(x)
	fmt.Println(string(data))

}

func marshalJson1() {
	x := map[string]string{
		"foo": "bar",
	}
	data, _ := json.Marshal(x)
	fmt.Println(string(data))
}
