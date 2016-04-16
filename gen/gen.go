package gen

import (
	"fmt"
	"github.com/bantl23/cminus/log"
	"github.com/bantl23/cminus/symtbl"
	"github.com/bantl23/cminus/syntree"
	"os"
)

type Gen struct {
	filename string
	file     *os.File
	pc       int
	mp       int
	gp       int
	ac       int
	ac1      int
	tmp      int
	loc      int
	highLoc  int
}

func NewGen(filename string) *Gen {
	file, err := os.Create(filename)
	if err != nil {
		return nil
	}
	g := new(Gen)
	g.filename = filename
	g.file = file
	g.pc = 7
	g.mp = 6
	g.gp = 5
	g.ac = 0
	g.ac1 = 1
	g.tmp = 1
	g.loc = 0
	g.highLoc = 0
	return g
}

func (g *Gen) emit(out string) {
	g.file.WriteString(out)
	log.CodeLog.Printf(out)
}

func (g *Gen) emitRO(opcode string, target int, source0 int, source1 int, comment string) {
	out := fmt.Sprintf("%3d: %5s %d,%d,%d\t* %s\n", g.loc, opcode, target, source0, source1, comment)
	g.emit(out)
	g.loc = g.loc + 1
	if g.highLoc < g.loc {
		g.highLoc = g.loc
	}
}

func (g *Gen) emitRM(opcode string, target int, offset int, base int, comment string) {
	out := fmt.Sprintf("%3d: %5s %d,%d(%d)\t* %s\n", g.loc, opcode, target, offset, base, comment)
	g.emit(out)
	g.loc = g.loc + 1
	if g.highLoc < g.loc {
		g.highLoc = g.loc
	}
}

func (g *Gen) emitRMAbs(opcode string, target int, abs int, comment string) {
	out := fmt.Sprintf("%3d: %5s %d,%d(%d)\t* %s\n", g.loc, opcode, target, abs-(g.loc+1), g.pc, comment)
	g.emit(out)
	if g.highLoc < g.loc {
		g.highLoc = g.loc
	}
}

func (g *Gen) emitComment(comment string) {
	out := fmt.Sprintf("* %s\n", comment)
	g.emit(out)
}

func (g *Gen) emitSkip(amount int) int {
	log.CodeLog.Printf("skipping %d amount\n", amount)
	i := g.loc
	g.loc = g.loc + amount
	if g.highLoc < g.loc {
		g.highLoc = g.loc
	}
	return i
}

func (g *Gen) emitBackup(loc int) {
	log.CodeLog.Printf("backing up to %d\n", loc)
	if loc > g.highLoc {
		log.ErrorLog.Printf(">>>>> Error in emitBackup\n")
	}
	g.loc = loc
}

func (g *Gen) emitRestore() {
	log.CodeLog.Printf("restoring to %d\n", g.highLoc)
	g.loc = g.highLoc
}

func (g *Gen) load() {
	g.emitComment("cminus compilation into tiny machine for " + g.filename)
	g.emitComment("prelude beg")
	g.emitRM("LD", g.mp, 0, g.ac, "load maxaddress from location 0")
	g.emitRM("ST", g.ac, 0, g.ac, "clear location 0")
	g.emitComment("prelude end")
}

func (g *Gen) halt() {
	g.emitRO("HALT", 0, 0, 0, "halting program")
}

func (g *Gen) gen(node syntree.Node) {
	if node != nil {
		log.CodeLog.Printf("%+v\n", node)
		if node.IsStmt() {
			g.genStmt(node)
		} else if node.IsExp() {
			g.genExp(node)
		}
		g.gen(node.Sibling())
	}
}

func (g *Gen) genStmt(node syntree.Node) {
	if node.IsCompound() {
	} else if node.IsFunc() {
	} else if node.IsIteration() {
	} else if node.IsReturn() {
	} else if node.IsSelection() {
	}
}

func (g *Gen) genExp(node syntree.Node) {
	if node.IsAssign() {
	} else if node.IsCall() {
	} else if node.IsConst() {
		g.genConst(node)
	} else if node.IsOp() {
	} else if node.IsId() {
		if node.IsArray() {
		} else {
		}
	} else if node.IsParam() {
		if node.IsArray() {
		} else {
		}
	} else if node.IsVar() {
		if node.IsArray() {
		} else {
		}
	}
}

func (g *Gen) genConst(node syntree.Node) {
	comment := fmt.Sprintf("load const with %d", node.Value())
	g.emitRM("LDC", g.ac, node.Value(), 0, comment)
}

func (g *Gen) getId(node syntree.Node) {
	comment := fmt.Sprintf("load %s with %d", node.Name(), node.Value())
	if symtbl.GlbSymTblMap[node.SymKey()].HasId(node.Name()) {
		memLoc := symtbl.GlbSymTblMap[node.SymKey()].GetMemLoc(node.Name())
		g.emitRM("LD", g.ac, int(memLoc), g.gp, comment)
	} else {
		log.ErrorLog.Printf(">>>>> Error %s not found.", node.Name())
	}
}

func Generate(node syntree.Node, filename string) {
	g := NewGen(filename)
	if g != nil {
		g.load()
		g.gen(node)
		g.halt()
	} else {
		log.ErrorLog.Printf(">>>>> Error opening %s", filename)
	}
}
