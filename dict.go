package main

// Initialize dictionary
func (s *Machine) InitializeDict() {
	s.dict = Dictionary{
		contents: []Word{
			{name: "+", actions:  s.Plus }, 
			{name: "-", actions:  s.Minus }, 
			{name: "*", actions:  s.Mult }, 
			{name: "/", actions:  s.Div }, 
			{name: "MOD", actions: s.Mod },
			{name: ".", actions:  s.Print }, 
			{name: "=", actions:  s.Equal }, 
			{name: "0=", actions:  func() { s.SubRoutine(func() { s.Execute("0 =")}) }}, 
			{name: "<", actions:  s.LessThan },
			{name: "DUP", actions:  s.Dup }, 
			{name: "DUP2", actions:  s.Dup2 }, 
			{name: "SWAP", actions:  s.Swap }, 
			{name: "SWAP2", actions:  s.Swap2 }, 
			{name: "DROP", actions:  s.Drop }, 
			{name: "..", actions: s.DropStack},
			{name: "ROT", actions:  s.Rot }, 
			{name: "OVER", actions:  func() { s.SubRoutine(func() { s.Execute("SWAP DUP ROT ROT")} )}}, 
			{name: ".S", actions:  s.Prints }, 
			{name: "CR", actions:  func() {print("\n")} }, 
			{name: ":", actions: s.CompileWord },
			{name: ";",},
			{name: "IF", actions: s.If},
			{name: "ELSE",},
			{name: "THEN",},
			{name: "WHILE", actions: s.While},
			{name: "DO"},
			{name: "(", actions: s.Comment},
			{name: ")"},
			{name: ".\"", actions: s.PrintString},
			{name: "\""},
			{name: "VARIABLE", actions: s.CreateVariable},
			{name: "!", actions: s.Set },
			{name: "@", actions: s.Get },
		},
	}
}