package main

import (
	"strconv"
	"strings"
	// "fmt"
)

type Machine struct {
	stack []int
	rstack []int
	heap	[]int
	top 	int
	dict	Dictionary
	err		string
	addr 	int
	ss		[]string
	tmp		[]string
}

type Word struct {
	name		string
	actions	func()
	ptr			int
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

func (s *Machine) Mod() {
	b, _ := s.PopInt()
	a, _ := s.PopInt()
	s.Push(a % b)
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

func (s *Machine) DropStack() {
	s.stack = []int{}; return
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

func (s *Machine) Compile(delim string) (string, string) {
	s.addr++
	a := ""
	w := s.ss[s.addr]
	for w != delim {
		a += w + " "
		s.addr++
		w = s.ss[s.addr]
	}
	sl := strings.SplitN(a, " ", 2)
	return sl[0], sl[1]
}

func (s *Machine) CompileWord() {
	n, a := s.Compile(";")
	// Special case where I don't want to reset the addr
	s.RPop()
	s.RPush(s.addr)
	s.AddWord(Word{name: n, actions: func(){ s.SubRoutine(func() { s.Execute(a)}) },})
}

func (s *Machine) If() {
	n, a := s.Compile("THEN")
	s.RPop()
	s.RPush(s.addr)

	t := n + " " + a
	ops := strings.Split(t, " ELSE ")

	if (s.Pop() != 0) {
		s.SubRoutine(func() { s.Execute(ops[0]) } )
	} else {
		if len(ops) > 1 {
			s.SubRoutine(func() { s.Execute(ops[1]) } )
		}
	}

	return
}

func (s *Machine) While() {
	n, a := s.Compile("DO")
	s.RPop()
	s.RPush(s.addr)

	t := n + " " + a

	
	for (s.Pop() != 0 ) {
		// fmt.Println(t)
		s.Execute(t)
		// fmt.Println(s.rstack)
		// fmt.Println(s.stack)
	}
	s.addr = s.RPop()

	return
}

func (s *Machine) Comment() {
	_, _ = s.Compile(")")
	s.RPop()
	s.RPush(s.addr)
	return
}

func (s *Machine) PrintString() {
	n, a := s.Compile("\"")
	s.RPop()
	s.RPush(s.addr)
	print(n + " " + a)
	return
}

// Variables

func (s *Machine) CreateVariable() {
	n := s.ss[s.addr + 1] // get name of variable
	s.heap = append(s.heap, -1)
	ptr := len(s.heap) - 1
	s.AddWord(Word{ name: n, ptr: ptr })
	s.RPop()
	s.RPush(s.addr + 1)
}

func (s *Machine) Fetch(ptr int) {
	s.Push(s.heap[ptr])
}

func (s *Machine) Set() {
	n := s.ss[s.addr + 1] // get whatever you're setting
	var entry Word
	for _, m := range s.dict.contents {
		if n == m.name {
			entry = m
		}
	}
	s.heap[entry.ptr] = s.Pop()
	s.RPop()
	s.RPush(s.addr + 1)
}

func (s *Machine) Get() {
	n := s.ss[s.addr + 1] // get the var name
	var ptr int
	for _, m := range s.dict.contents {
		if n == m.name {
			ptr = m.ptr
		}
	}
	s.Push(s.heap[ptr])
	s.RPop()
	s.RPush(s.addr + 1)
}


// 

func (s *Machine) AddWord(word Word) {
	s.dict.contents = append(s.dict.contents, word)
}

func (s *Machine) RPop() int {
	top := len(s.rstack) - 1
	if top < 0 { return s.addr }
	addr := s.rstack[top]
	s.rstack = s.rstack[:top]
	return addr
}

func (s *Machine) RPush(addr int) {
	s.rstack = append(s.rstack, addr)
	return
}

func (s *Machine) SubRoutine(routine func()) {
	s.RPush(s.addr)
	s.tmp = s.ss
	routine()
	s.addr = s.RPop()
	s.ss = s.tmp
}

func (s *Machine) Execute(input string) {
	// print("***\nDEBUG: input = " + input + "\n")
	s.ss = strings.Fields(strings.Trim(input, "\n"))
	s.addr = 0
	for s.addr < len(s.ss) {
		s.ExecuteWord()
		s.ss = strings.Fields(strings.Trim(input, "\n"))
		s.tmp = s.ss
	}
}

func (s *Machine) ExecuteWord() {
		// Get the word
		w := s.ss[s.addr]
		// check if it's in the dictionary
		indict := false // assume it isn't
		for _, n := range s.dict.contents {
			if w == n.name {
				indict = true
				s.SubRoutine(n.actions)
			}
		}
		if !indict { s.TryToPush(w) }
		s.addr++
}