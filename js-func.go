package deno

import (
	elutils "github.com/rosbit/go-embedding-utils"
	"github.com/rosbit/go-expect"
	"strings"
	"bytes"
	"fmt"
	"reflect"
	"encoding/json"
)

func (js *Deno) bindFunc(funcName string, funcVarPtr interface{}) (err error) {
	helper, e := elutils.NewEmbeddingFuncHelper(funcVarPtr)
	if e != nil {
		err = e
		return
	}
	helper.BindEmbeddingFunc(js.wrapFunc(funcName, helper))
	return
}

func (js *Deno) wrapFunc(funcName string, helper *elutils.EmbeddingFuncHelper) elutils.FnGoFunc {
	return func(args []reflect.Value) (results []reflect.Value) {
		var jsArgs []interface{}

		// make js args
		itArgs := helper.MakeGoFuncArgs(args)
		for arg := range itArgs {
			jsArgs = append(jsArgs, arg)
		}

		// call JS function
		goVal, err := js.CallFunc(funcName, jsArgs...)
		isArray := false
		if err == nil && goVal != nil {
			_, isArray = goVal.([]interface{})
		}

		// convert result to golang
		results = helper.ToGolangResults(goVal, isArray, err)
		return
	}
}

func makeGoal(funcName string, args ...interface{}) (goal string, err error) {
	argc := len(args)
	if argc == 0 {
		goal = fmt.Sprintf("%s()", funcName)
		return
	}

	argv := make([]string, argc)
	for i, arg := range args {
		if arg == nil {
			argv[i] = "undefined"
		} else if jArg, e := json.Marshal(arg); e != nil {
			err = e
			return
		} else {
			argv[i] = string(jArg)
		}
	}
	goal = fmt.Sprintf("%s(%s)", funcName, strings.Join(argv, ","))
	return
}

func (js *Deno) call(goal string) (res interface{}, err error) {
	js.e.Send(fmt.Sprintf("%s\n", goal))

	for {
		_, _, e := js.e.ExpectCases(
			&expect.Case{Exp: promptRE, MatchedOnly: true, ExpMatched: func(_ []byte) expect.Action{
				return expect.Break
			}},
			&expect.Case{Exp: goalRE/*, SkipTill: '\n'*/, ExpMatched: func(_ []byte) expect.Action{
				return expect.Continue
			}},
			&expect.Case{Exp: errRE, ExpMatched: func(m []byte) expect.Action{
				err = fmt.Errorf("%s", m)
				return expect.Continue
			}},
			&expect.Case{Exp: errAtRE, ExpMatched: func(m []byte) expect.Action{
				if err != nil {
					err = fmt.Errorf("%s%s", err.Error(), m)
				} else {
					err = fmt.Errorf("%s", m)
				}
				return expect.Continue
			}},
			// &expect.Case{Exp: blankRE, SkipTill: '\n'},
			&expect.Case{Exp: nonResRE, SkipTill: '\n'},
			&expect.Case{Exp: resultRE, ExpMatched: func(m []byte) expect.Action{
				l := len(m)
				switch {
				case l == 0:
					return expect.Continue
				case m[l-1] == '\r':
					return expect.Continue
				case bytes.HasPrefix(m, functionId):
					return expect.Continue
				}
				if r, e1 := fromJSValue(m); e1 != nil {
					// fmt.Printf(">>>%s<<<, len: %d, %x\n", m, len(m), m)
					fmt.Printf("%s", m)
				} else {
					res = r
				}
				return expect.Continue
			}},
		)

		if e != nil {
			if e == expect.TimedOut || e == expect.NotFound {
				continue
			}
			err = e
		}
		break
	}

	return
}
