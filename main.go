package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
	"io"
	"strconv"
)

type Machine struct {
	stack []int
	top 	int
	dict	Dictionary
	err		string
}

type Word struct {
	name		string
	actions	string
}

type Dictionary struct {
	contents 	[]Word
}

func (s *Machine) Push(value int) {
	s.stack = append(s.stack, value)
	s.top++
	return
}

func (s *Machine) Pop() int {
	popped := s.stack[s.top]
	s.stack = s.stack[:s.top]
	s.top--
	return popped
}

func (s *Machine) TryToPush(in string) {
	v, err := strconv.Atoi(in)
	if err != nil {
		s.err = in + "? "
		return
	}
	s.Push(v)
	return
}

// math(s)
func (s *Machine) PopInt() (int, error) {
	return s.Pop(), nil
}

func (s *Machine) Plus() {
	b, _ := s.PopInt()
	a, _ := s.PopInt()
	s.Push(a + b)
	return
}

func (s *Machine) Minus() {
	b, _ := s.PopInt()
	a, _ := s.PopInt()
	s.Push(a - b)
	return
}

func (s *Machine) Mult() {
	b, _ := s.PopInt()
	a, _ := s.PopInt()
	s.Push(a * b)
	return
}

func (s *Machine) Div() {
	b, _ := s.PopInt()
	a, _ := s.PopInt()
	s.Push(a / b)
	return
}

func (s *Machine) Equal() {
	a, _ := s.PopInt()
	b, _ := s.PopInt()
	// Return -1 for true, otherwise 0
	if a == b { s.Push(-1); return }
	s.Push(0); return
}

// Printing functions
func (s *Machine) Print() {
	if s.top == -1 {
		s.err = "Stack underflow! "
	} else {
		fmt.Print(s.Pop(), " ")
	}
	return
}

func (s *Machine) Prints() {
	for _, v := range s.stack {
		fmt.Print(v, " ")
	}
	return
}

// Machine functions (i know it's all Machine functions stop)
func (s *Machine) Dup() {
	a := s.Pop()
	s.Push(a)
	s.Push(a)
	return
}

func (s *Machine) Swap() {
	b := s.Pop()
	a := s.Pop()
	s.Push(b)
	s.Push(a)
	return
}

func (s *Machine) Drop() {
	s.Pop(); return
}

func (s *Machine) Rot() {
	c := s.Pop()
	b := s.Pop()
	a := s.Pop()
	s.Push(b)
	s.Push(c)
	s.Push(a)
	return
}

// Control states are weirdddd
func (s *Machine) CompileStep(idx int, ss []string, n string, a string) (string, string, bool) {
	// get the word you're looking at
	v := ss[idx]
	if v == ";" { 
		// COMPILE TIME IS OVER, put it in the dict
		word := Word{
			name: n,
			actions: a,
		}
		s.dict.contents = append(s.dict.contents, word)
		return "", "", false
	} else {
		if ss[idx - 1] == ":" {
			// Get the name if you haven't already
			n = v
			a = ""
		} else {
			// oops all strings
			a += v + " "
		}
		return n, a, true
	}
}

func (s *Machine) ConditionalStep(idx int, ss[]string, cons string, alt string, b_alt bool) (string, string, bool, bool) {
	// For control words, return things as they are
	v := ss[idx]
	if (v == "IF") { b_alt = false } 
	if v == "ELSE" { b_alt = true }
	// THEN is the ending
	if v == "THEN" { return cons, alt, false, b_alt }

	if !b_alt {
		// build the consequent
		cons += v + " "
	} else {
		alt += v + " "
	}

	return cons, alt, true, b_alt
}

func (s *Machine) ExecuteConditional(cons string, alt string) {
	if (s.Pop() != 0) {
		s.Execute(cons)
	} else {
		s.Execute(alt)
	}
}

func (s *Machine) Execute(input string) {
	ss := strings.Fields(strings.Trim(input, "\n"))
	var compile, comment bool
	var conditional, b_alt bool
	var name, actions string
	var cons, alt string
	
	for idx, v := range ss {
		if compile {
			name, actions, compile = s.CompileStep(idx, ss, name, actions)
		} else if conditional {
			cons, alt, conditional, b_alt = s.ConditionalStep(idx, ss, cons, alt, b_alt)
		} else if comment {
			// Do absolutely nothing
		} else {
			indict := false
			for _, n := range s.dict.contents {
				if v == n.name {
					indict = true
					s.Execute(n.actions)
				}
			}
			if !indict {
				switch v {
				// maths
				case "+": s.Plus()
				case "-": s.Minus()
				case "*": s.Mult()
				case "/": s.Div()
				case ".": s.Print()
				
				// bools
				case "=": s.Equal()
				case "0=": s.Push(0); s.Equal()
				
				// Machine fun!(ctions)
				case "DUP": s.Dup()
				case "SWAP": s.Swap()
				case "DROP": s.Drop()
				case "ROT": s.Rot()
				case "OVER": s.Swap(); s.Dup(); s.Rot(); s.Rot()
				case ".S": s.Prints()

				// flow controls
				case ":": compile = true
				case ";": compile = false
				case "IF": conditional = true
				case "THEN": conditional = false; s.ExecuteConditional(cons, alt)
				case "(": comment = true
				case ")": comment = false

				// *sigh* put it in the Machine with the others
				default: s.TryToPush(v)
				}			
			}
		}
	}
}

func main () {
	reader := bufio.NewReader(os.Stdin)
	s := Machine{
		top: -1,
		dict : Dictionary{
			[]Word{
			},
		},
		err: "",
	}

	for {
		input, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			os.Exit(1)
		}
		if len(strings.Fields(input)) > 0 {
			s.Execute(strings.ToUpper(input))
			if s.err != "" {
				fmt.Println(s.err)
				s.err = ""
			} else {
				fmt.Println("ok")
			}
		}
		if err == io.EOF { os.Exit(0) }
	}
}