package deno

import (
	"github.com/rosbit/go-expect"
	"reflect"
	"bytes"
	"fmt"
	"io"
	"time"
	"runtime"
	"strings"
)

const (
	timeout = 1 * time.Second
)

type Deno struct {
	e *expect.Expect
}

func NewDeno(denoExe string, jsFile string) (*Deno, error) {
	if len(jsFile) == 0 {
		return nil, fmt.Errorf("jsFile expectd")
	}
	e, err := spawn(map[string]string{"NO_COLOR":"1"}, denoExe, "repl", fmt.Sprintf("--eval-file=%s", jsFile))
	if err != nil {
		return nil, err
	}
	e.SetTimeout(timeout)
	// e.RemoveColor()
	e.SetNotMatchedHandler(func(m []byte) (skipN int){
		// fmt.Printf("not matched: >>>%s<<<, len: %d, %x\n", m, len(m), m)
		return 0
	})

	for {
		_, _, e1 := e.ExpectCases(
			&expect.Case{Exp: errLoadRE, ExpMatched: func(m []byte) expect.Action {
				err = fmt.Errorf("%s", m)
				return expect.Continue
			}},
			&expect.Case{Exp: promptRE, MatchedOnly: true, ExpMatched: func(_ []byte) expect.Action{
				return expect.Break
			}},
			&expect.Case{Exp: ignoreRE, SkipTill: '\n'},
			&expect.Case{Exp: ignore2RE, SkipTill: '\n'},
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
			&expect.Case{Exp: blankRE, SkipTill: '\n'},
		)
		if e1 != nil {
			// timeout
			if e1 == expect.TimedOut || e1 == expect.NotFound {
				continue
			}
			if e1 == io.EOF {
				break
			}
			err = e1
		}
		break
	}
	if err != nil {
		return nil, err
	}

	t := &Deno{e: e}
	runtime.SetFinalizer(t, closeDeno)
	return t, nil
}

func closeDeno(js *Deno) {
	js.Quit()
}

func (js *Deno) Quit() {
	if js.e == nil {
		return
	}
	js.e.Send("Deno.exit()\n")
	time.Sleep(100*time.Millisecond)
	js.e.Close()
	js.e = nil
}

func (js *Deno) checkFuncName(funcName string) (err error) {
	funcName = strings.TrimSpace(funcName)
	if len(funcName) == 0 {
		return fmt.Errorf("funcName expected")
	}

	_, e := js.GetGlobal(funcName)
	if e == nil {
		err = fmt.Errorf("%s is not function", funcName)
		return
	}
	if bytes.HasPrefix([]byte(e.Error()), functionId) {
		return
	}
	err = e
	return
}

func (js *Deno) GetGlobal(varName string) (val interface{}, err error) {
	varName = strings.TrimSpace(varName)
	if len(varName) == 0 {
		err = fmt.Errorf("varName expected")
		return
	}
	js.e.Send(fmt.Sprintf("%s\n", varName))

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
				// fmt.Printf("resultRE matched in GetGlobal(): >>>%s<<<, %x\n", m, m)
				if bytes.HasPrefix(m, functionId) {
					err = fmt.Errorf("%s", m)
					return expect.Continue
				}
				val, err = fromJSValue(m)
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

func (js *Deno) CallFunc(funcName string, args ...interface{}) (res interface{}, err error) {
	if len(funcName) == 0 {
		err = fmt.Errorf("funcName expected")
		return
	}
	defer func() {
		if r := recover(); r != nil {
			if v, o := r.(error); o {
				err = v
				return
			}
			err = fmt.Errorf("%v", r)
			return
		}
	}()

	goal, e := makeGoal(funcName, args...)
	if e != nil {
		err = e
		return
	}
	res, err = js.call(goal)
	return
}

// bind a var of golang func with JS function name, so calling golang func // is just calling the related JS function.
// @param funcVarPtr  in format `var funcVar func(....) ...; funcVarPtr = &funcVar`
func (js *Deno) BindFunc(funcName string, funcVarPtr interface{}) (err error) {
	if funcVarPtr == nil {
		err = fmt.Errorf("funcVarPtr must be a non-nil poiter of func")
		return
	}
	t := reflect.TypeOf(funcVarPtr)
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Func {
		err = fmt.Errorf("funcVarPtr expected to be a pointer of func")
		return
	}

	if err = js.checkFuncName(funcName); err != nil {
		return
	}
	return js.bindFunc(funcName, funcVarPtr)
}
