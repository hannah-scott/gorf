package main

import (
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
	if s.top == -1 {
		s.err = "Stack underflow! "
		return -1
	}
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

func (s *Machine) LessThan() {
	a := s.Pop()
	if s.Pop() < a { s.Push(-1); return }
	s.Push(0); return
}

// Printing functions
func (s *Machine) Print() {
	if s.top == -1 {
		s.err = "Stack underflow! "
	} else {
		print(s.Pop(), " ")
	}
	return
}

func (s *Machine) Prints() {
	for _, v := range s.stack {
		print(v, " ")
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

func (s *Machine) Dup2() {
	b := s.Pop()
	a := s.Pop()
	s.Push(a)
	s.Push(b)
	s.Push(a)
	s.Push(b)
}

func (s *Machine) Swap() {
	b := s.Pop()
	a := s.Pop()
	s.Push(b)
	s.Push(a)
	return
}

func (s *Machine) Swap2() {
	d := s.Pop()
	c := s.Pop()
	b := s.Pop()
	a := s.Pop()
	s.Push(c)
	s.Push(d)
	s.Push(a)
	s.Push(b)
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
	if v == "ELSE" { b_alt = true; return cons, alt, true, b_alt }
	// THEN is the ending
	if v == "THEN" {  s.ExecuteConditional(cons, alt); return cons, alt, false, b_alt }

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

func (s *Machine) PrintStep(idx int, ss []string, ps string) (string, bool) {
	v := ss[idx]
	if v == "\"" {
		s.ExecutePrint(ps)
		return "", false
	} else {
		return ps + " " + v, true
	}
}

func (s *Machine) ExecutePrint(ps string) {
	print(strings.TrimLeft(ps, " ") + " ")
} 

func (s *Machine) LoopStep(idx int, ss []string, ls string) (string, bool) {
	// Get contents of loop
	v := ss[idx]
	if v == "DO" {
		s.ExecuteLoop(ls)
		return "", false
	} else {
		return ls + " " + v, true
	}
}

func (s *Machine) ExecuteLoop(ls string) {
	s.Execute("DUP 0=")
	for s.Pop() == 0 {
		s.Execute(ls)
		s.Execute("DUP 0=")
	}
}

func (s *Machine) Execute(input string) {
	ss := strings.Fields(strings.Trim(input, "\n"))
	var compile, comment, print, conditional, loop, b_alt bool
	var name, actions string
	var cons, alt, ps, ls string
	
	for idx, v := range ss {
		if compile {
			name, actions, compile = s.CompileStep(idx, ss, name, actions)
		} else if conditional {
			cons, alt, conditional, b_alt = s.ConditionalStep(idx, ss, cons, alt, b_alt)
		} else if loop {
			ls, loop = s.LoopStep(idx, ss, ls)
		} else if print {
			ps, print = s.PrintStep(idx, ss, ps)
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
				case "<": s.LessThan()
				
				// Machine fun!(ctions)
				case "DUP": s.Dup()
				case "DUP2": s.Dup2()
				case "SWAP": s.Swap()
				case "SWAP2": s.Swap2()
				case "DROP": s.Drop()
				case "ROT": s.Rot()
				case "OVER": s.Swap(); s.Dup(); s.Rot(); s.Rot()
				case ".S": s.Prints()

				// flow controls
				case ":": compile = true
				case ";": compile = false
				case "IF": conditional = true
				case "THEN": conditional = false
				case "WHILE": loop = true
				case "DO": loop = false
				case "(": comment = true
				case ")": comment = false
				case ".\"": print = true
				case "\"": print = false

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
				print(s.err + "\n")
				s.err = ""
			} else {
				print("ok\n")
			}
		}
		if err == io.EOF { os.Exit(0) }
	}
}