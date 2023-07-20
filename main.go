package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
	"strconv"
	"io"
)

type Stack struct {
	items []string
	top 	int
	dict	Dictionary
}

type Word struct {
	name		string
	actions	string
}

type Dictionary struct {
	contents 	[]Word
}

func (s *Stack) Push(value string) {
	s.items = append(s.items, value)
	s.top++
	return
}

func (s *Stack) Pop() string {
	popped := s.items[s.top]
	s.items = s.items[:s.top]
	s.top--
	return popped
}

func typeOf(i string) string {
	return fmt.Sprintf("%T", i)
}

// math(s)
func (s *Stack) PopInt() (int, error) {
	return strconv.Atoi(s.Pop())
}

func (s *Stack) Plus() {
	b, _ := s.PopInt()
	a, _ := s.PopInt()
	s.Push(strconv.Itoa(a + b))
	return
}

func (s *Stack) Minus() {
	b, _ := s.PopInt()
	a, _ := s.PopInt()
	s.Push(strconv.Itoa(a - b))
	return
}

func (s *Stack) Mult() {
	b, _ := s.PopInt()
	a, _ := s.PopInt()
	s.Push(strconv.Itoa(a * b))
	return
}

func (s *Stack) Div() {
	b, _ := s.PopInt()
	a, _ := s.PopInt()
	s.Push(strconv.Itoa(a / b))
	return
}

func (s *Stack) Equal() {
	a, _ := s.PopInt()
	b, _ := s.PopInt()
	// Return -1 for true, otherwise 0
	if a == b { s.Push("-1"); return }
	s.Push("0"); return
}

// Printing functions
func (s *Stack) Print() {
	if s.top == -1 {
		fmt.Print("Stack underflow! ")
	} else {
		fmt.Print(s.Pop(), " ")
	}
	return
}

func (s *Stack) Prints() {
	for _, v := range s.items {
		fmt.Print(v, " ")
	}
	return
}

// Stack functions (i know it's all stack functions stop)
func (s *Stack) Dup() {
	a := s.Pop()
	s.Push(a)
	s.Push(a)
	return
}

func (s *Stack) Swap() {
	b := s.Pop()
	a := s.Pop()
	s.Push(b)
	s.Push(a)
	return
}

func (s *Stack) Drop() {
	s.Pop(); return
}

func (s *Stack) Rot() {
	c := s.Pop()
	b := s.Pop()
	a := s.Pop()
	s.Push(b)
	s.Push(c)
	s.Push(a)
	return
}

// Control states are weirdddd
func (s *Stack) CompileStep(idx int, ss []string, n string, a string) (string, string, bool) {
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


func (s *Stack) Execute(input string) {
	ss := strings.Fields(strings.Trim(input, "\n"))
	var compile, conditional, comment bool
	var name, actions string
	
	for idx, v := range ss {
		if compile {
			name, actions, compile = s.CompileStep(idx, ss, name, actions)
		} else if conditional {
			
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
				case "0=": s.Push("0"); s.Equal()
				
				// stack fun!(ctions)
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
				case "THEN": conditional = false
				case "(": comment = true
				case ")": comment = false

				// *sigh* put it in the stack with the others
				default: s.Push(v)
				}			
			}
		}
	}
}

func main () {
	reader := bufio.NewReader(os.Stdin)
	s := Stack{
		top: -1,
		dict : Dictionary{
			[]Word{
			},
		},
	}

	for {
		input, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			os.Exit(1)
		}
		if len(strings.Fields(input)) > 0 {
			s.Execute(strings.ToUpper(input))
			fmt.Println("ok")
		}
		if err == io.EOF { os.Exit(0) }
	}
}