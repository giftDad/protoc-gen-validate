package examples

import (
	"testing"

	"dario.cat/mergo"
	"github.com/stretchr/testify/assert"
)

func TestNumerics(t *testing.T) {
	type _case struct {
		r        *NumericsReq
		equal    bool
		errField string
	}

	noErrReq := new(NumericsReq) // merge structs to avoid prefix error

	_cases := []_case{
		// test `@eq:1.23` true
		{r: &NumericsReq{A: 1.23}, equal: true, errField: "A"},
		// test `@eq:1.23` false
		{r: &NumericsReq{A: 1.22}, equal: false, errField: "A"},

		// test `@gt:10` true
		{r: &NumericsReq{B: 11}, equal: true, errField: "B"},
		// test `@gt:10` false
		{r: &NumericsReq{B: 10}, equal: false, errField: "B"},

		// test `@lt:20` true
		{r: &NumericsReq{B: 19}, equal: true, errField: "B"},
		// test `@lt:20` false
		{r: &NumericsReq{B: 20}, equal: false, errField: "B"},

		// test `@gte:10` true
		{r: &NumericsReq{C: 10}, equal: true, errField: "C"},
		// test `@gte:10` false
		{r: &NumericsReq{C: 9}, equal: false, errField: "C"},

		// test `@lte:20` true
		{r: &NumericsReq{C: 20}, equal: true, errField: "C"},
		// test `@lte:20` false
		{r: &NumericsReq{C: 21}, equal: false, errField: "C"},

		// test `@in:[1,2,3]` true
		{r: &NumericsReq{D: 2}, equal: true, errField: "D"},
		// test `@in:[1,2,3]` false
		{r: &NumericsReq{D: 4}, equal: false, errField: "D"},

		// test `@not_in:[1,2,3]` true
		{r: &NumericsReq{E: 4}, equal: true, errField: "E"},
		// test `@not_in:[1,2,3]` false
		{r: &NumericsReq{E: 2}, equal: false, errField: "E"},

		// test `@range:(1,5)` true
		{r: &NumericsReq{F: 4}, equal: true, errField: "F"},
		// test `@range:(1,5)` false
		{r: &NumericsReq{F: 1}, equal: false, errField: "F"},
		// test `@range:(1,5)` false
		{r: &NumericsReq{F: 5}, equal: false, errField: "F"},

		// test `@range:[1,5]` true
		{r: &NumericsReq{G: 1}, equal: true, errField: "G"},
		// test `@range:[1,5]` true
		{r: &NumericsReq{G: 5}, equal: true, errField: "G"},
		// test `@range:[1,5]` false
		{r: &NumericsReq{G: 6}, equal: false, errField: "G"},
	}

	for i, c := range _cases {
		req := new(NumericsReq)
		*req = *noErrReq // value copy
		err := mergo.Merge(req, *c.r, mergo.WithOverride)
		assert.Nil(t, err)

		err = req.Validate()
		if c.equal {
			if err != nil {
				assert.NotEqual(t, c.errField, err.(NumericsReqValidationError).field, i+1)
				noErrReq = req
			}
		} else {
			assert.Equal(t, c.errField, err.(NumericsReqValidationError).field, i+1)
		}
	}

}

func TestString(t *testing.T) {
	type _case struct {
		r        *StringsReq
		equal    bool
		errField string
	}

	noErrReq := new(StringsReq) // merge structs to avoid prefix error

	_cases := []_case{
		// test `@contains:"bar"` true
		{r: &StringsReq{A: "bars"}, equal: true, errField: "A"},
		// test `@contains:"bar"` false
		{r: &StringsReq{A: "basr"}, equal: false, errField: "A"},

		// test `@not_contains:"bar"` true
		{r: &StringsReq{B: "basr"}, equal: true, errField: "B"},
		// test `@not_contains:"bar"` false
		{r: &StringsReq{B: "bars"}, equal: false, errField: "B"},

		// test `@eq:"bar"` true
		{r: &StringsReq{C: "bar"}, equal: true, errField: "C"},
		// test `@eq:"bar"` false
		{r: &StringsReq{C: "bars"}, equal: false, errField: "C"},

		// test `@in:["foo", "bar", "baz"]` true
		{r: &StringsReq{D: "bar"}, equal: true, errField: "D"},
		// test `@in:["foo", "bar", "baz"]` false
		{r: &StringsReq{D: "bars"}, equal: false, errField: "D"},

		// test `@not_in:["foo", "bar", "baz"]` true
		{r: &StringsReq{E: "bars"}, equal: true, errField: "E"},
		// test `@not_in:["foo", "bar", "baz"]` false
		{r: &StringsReq{E: "bar"}, equal: false, errField: "E"},

		// test `@len:5` true
		{r: &StringsReq{F: "abced"}, equal: true, errField: "F"},
		// test `@len:5` false
		{r: &StringsReq{F: "abc"}, equal: false, errField: "F"},

		// test `@min_len:5` true
		{r: &StringsReq{G: "abcdef"}, equal: true, errField: "G"},
		// test `@min_len:5` false
		{r: &StringsReq{G: "abcd"}, equal: false, errField: "G"},

		// test `@max_len:10` true
		{r: &StringsReq{G: "abcdef"}, equal: true, errField: "G"},
		// test `@max_len:10` false
		{r: &StringsReq{G: "abcdefghijk"}, equal: false, errField: "G"},

		// test `@pattern:"(?i)^[0-9a-f]+$"` true
		{r: &StringsReq{H: "aAfF09"}, equal: true, errField: "H"},
		// test `@pattern:"(?i)^[0-9a-f]+$"` false
		{r: &StringsReq{H: "a.b"}, equal: false, errField: "H"},

		// test `@prefix:"foo"` true
		{r: &StringsReq{I: "foozxc"}, equal: true, errField: "I"},
		// test `@prefix:"foo"` false
		{r: &StringsReq{I: "fozxco"}, equal: false, errField: "I"},

		// test `@suffix:"bar"` true
		{r: &StringsReq{J: "abcbar"}, equal: true, errField: "J"},
		// test `@suffix:"bar"` false
		{r: &StringsReq{J: "abcar"}, equal: false, errField: "J"},

		// test `@type:url` true
		{r: &StringsReq{K: "http://www.baidu.com"}, equal: true, errField: "K"},
		// test `@type:url` false
		{r: &StringsReq{K: "aaa"}, equal: false, errField: "K"},

		// test `@type:phone` true
		{r: &StringsReq{L: "15801812345"}, equal: true, errField: "L"},
		// test `@type:phone` false
		{r: &StringsReq{L: "12345677"}, equal: false, errField: "L"},

		// test `@type:email` true
		{r: &StringsReq{M: "12345@gmail.com"}, equal: true, errField: "M"},
		// test `@type:email` false
		{r: &StringsReq{M: "12345fsad.c"}, equal: false, errField: "M"},

		// test `@type:ip` true
		{r: &StringsReq{N: "127.0.0.1"}, equal: true, errField: "N"},
		// test `@type:ip` true
		{r: &StringsReq{N: "::ffff:192.0.2.1"}, equal: true, errField: "N"},
		// test `@type:ip` false
		{r: &StringsReq{N: "1230.0.1"}, equal: false, errField: "N"},
	}

	for i, c := range _cases {
		req := new(StringsReq)
		*req = *noErrReq // value copy
		err := mergo.Merge(req, *c.r, mergo.WithOverride)
		assert.Nil(t, err)

		err = req.Validate()
		if c.equal {
			if err != nil {
				assert.NotEqual(t, c.errField, err.(StringsReqValidationError).field, i+1)
				noErrReq = req
			}
		} else {
			assert.Equal(t, c.errField, err.(StringsReqValidationError).field, i+1)
		}
	}
}

func TestRepeated(t *testing.T) {
	type _case struct {
		r        *RepeatedReq
		equal    bool
		errField string
	}

	noErrReq := new(RepeatedReq) // merge structs to avoid prefix error

	_cases := []_case{
		// test `@min_items:1` true
		{r: &RepeatedReq{A: []int32{1}}, equal: true, errField: "A"},
		// test `@min_items:1` false
		{r: &RepeatedReq{A: []int32{}}, equal: false, errField: "A"},

		// test `@max_items:2` true
		{r: &RepeatedReq{A: []int32{1, 2}}, equal: true, errField: "A"},
		// test `@max_items:2` false
		{r: &RepeatedReq{A: []int32{1, 2, 3}}, equal: false, errField: "A"},

		// test `@unique:true` true
		{r: &RepeatedReq{A: []int32{1}, B: []int64{1, 2}}, equal: true, errField: "B"},
		// test `@unique:true` false
		{r: &RepeatedReq{A: []int32{1}, B: []int64{1, 1}}, equal: false, errField: "B"},

		// test `@unique:true` true
		{r: &RepeatedReq{A: []int32{1}, C: []string{"abb", "bcc"}}, equal: true, errField: "C"},
		// test `@unique:true` false
		{r: &RepeatedReq{A: []int32{1}, C: []string{"abc", "abc"}}, equal: false, errField: "C"},

		// test `embed` false
		{r: &RepeatedReq{A: []int32{1}, D: []*NumericsReq{{A: 1}}}, equal: false, errField: "D"},

		// test `every item @eq:1.23` true
		{r: &RepeatedReq{A: []int32{1}, E: []float32{1.23, 1.23}}, equal: true, errField: "E"},
		// test `every item @eq:1.23` false
		{r: &RepeatedReq{A: []int32{1}, E: []float32{1.23, 5}}, equal: false, errField: "E"},
	}

	for i, c := range _cases {
		req := new(RepeatedReq)
		*req = *noErrReq                                  // value copy
		err := mergo.Merge(req, *c.r, mergo.WithOverride) // repeated unable override
		assert.Nil(t, err)

		err = req.Validate()
		if c.equal {
			if err != nil {
				assert.NotEqual(t, c.errField, err.(RepeatedReqValidationError).field, i+1)
				noErrReq = req
			}
		} else {
			assert.Equal(t, c.errField, err.(RepeatedReqValidationError).field, i+1)
		}
	}
}

func TestRequired(t *testing.T) {
	type _case struct {
		r        *Required
		equal    bool
		errField string
	}

	_cases := []_case{
		// test `@required:true` true
		{r: &Required{A: &Foo{}}, equal: true, errField: "A"},
		// test `@required:true` false
		{r: &Required{A: nil}, equal: false, errField: "A"},
	}

	for i, c := range _cases {
		err := c.r.Validate()
		if c.equal {
			if err != nil {
				assert.NotEqual(t, c.errField, err.(RequiredValidationError).field, i+1)
			}
		} else {
			assert.Equal(t, c.errField, err.(RequiredValidationError).field, i+1)
		}
	}
}
