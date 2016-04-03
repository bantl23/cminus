package symtbl

import (
	"fmt"
)

type Value struct {
	MemLoc  MemLoc
	SymType SymbolType
	Args    []SymbolType
	Lines   []int
}

type SymTbl map[string]*Value

func (s *SymTbl) PrintTable() {
	fmt.Printf("    Variable Name Memory Location Type Args Lines\n")
	fmt.Printf("    ============= =============== ==== ==== =====\n")
	for i, e := range *s {
		fmt.Printf("    %13s 0x%013x %4s %+v %+v\n", i, e.MemLoc, e.SymType, e.Args, e.Lines)
	}
	fmt.Printf("\n")
}
