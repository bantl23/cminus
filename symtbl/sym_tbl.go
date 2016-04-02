package symtbl

import (
	"fmt"
)

type Value struct {
	MemLoc  MemLoc
	SymType SymbolType
	Lines   []int
}

type SymTbl map[string]*Value

func (s *SymTbl) PrintTable() {
	fmt.Printf("    Variable Name Type Memory Location Lines\n")
	fmt.Printf("    ============= ==== =============== =====\n")
	for i, e := range *s {
		fmt.Printf("    %13s %4s 0x%013x %+v\n", i, e.SymType, e.MemLoc, e.Lines)
	}
	fmt.Printf("\n")
}
