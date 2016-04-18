package gen

import (
	"fmt"
	"github.com/bantl23/cminus/log"
	"github.com/bantl23/cminus/symtbl"
	"github.com/bantl23/cminus/syntree"
	"os"
)

const (
	ofpFO  int = 0  // old frame pointer
	retFO  int = -1 // return address
	initFO int = -2 // param list
)

const (
	pc  int = 7 // program counter
	mp  int = 6 // memory pointer
	gp  int = 5 // global pointer
	ac1 int = 1 // accumlator
	ac  int = 0 // accumlator
)

type Gen struct {
	filename string
	file     *os.File
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
	g.loc = 0
	g.highLoc = 0
	return g
}

func (g *Gen) emit(out string) {
	g.file.WriteString(out)
	log.DstLog.Printf(out)
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

func (g *Gen) genCompound(node syntree.Node) {
	if len(node.Children()) > 0 {
		n0 := node.Children()[0]
		for n0 != nil {
			log.CodeLog.Printf("%+v", n0)
			if n0.IsVar() {
				if n0.IsArray() {
					log.CodeLog.Printf("+%d offset", n0.Value())
				} else {
					log.CodeLog.Printf("+1 offset")
				}
			}
			n0 = n0.Sibling()
		}
	}

	if len(node.Children()) > 1 {
		n1 := node.Children()[1]
		g.gen(n1)
	}
}

func (g *Gen) genFunction(node syntree.Node) {
	if len(node.Children()) > 0 {
		n0 := node.Children()[0]
		g.gen(n0)
	}

	if len(node.Children()) > 1 {
		n1 := node.Children()[1]
		g.gen(n1)
	}
}

func (g *Gen) genIteration(node syntree.Node) {
	if len(node.Children()) > 0 {
		n0 := node.Children()[0]
		g.gen(n0)
	}

	if len(node.Children()) > 1 {
		n1 := node.Children()[1]
		g.gen(n1)
	}
}

func (g *Gen) genReturn(node syntree.Node) {
	g.gen(node)
	g.emitRM("LD", pc, retFO, mp, "return to caller")
}

func (g *Gen) genSelection(node syntree.Node) {
	if len(node.Children()) > 0 {
		n0 := node.Children()[0]
		g.gen(n0)
	}

	if len(node.Children()) > 1 {
		n1 := node.Children()[0]
		g.gen(n1)
	}

	if len(node.Children()) > 2 {
		n2 := node.Children()[0]
		g.gen(n2)
	}
}

func (g *Gen) genAssign(node syntree.Node) {
	if len(node.Children()) > 0 {
		n0 := node.Children()[0]
		g.gen(n0)
	}

	if len(node.Children()) > 1 {
		n1 := node.Children()[0]
		g.gen(n1)
	}
}

func (g *Gen) genCall(node syntree.Node) {
	if len(node.Children()) > 0 {
		n0 := node.Children()[0]
		g.gen(n0)
	}
	if node.Name() == "input" {
		g.emitRO("IN", ac, 0, 0, "read from stdin into ac")
	} else if node.Name() == "output" {
		g.emitRM("LD", ac, 0, 0, "load into ac")
		g.emitRO("OUT", ac, 0, 0, "write to stdout with ac")
	} else {
		log.CodeLog.Printf("TODO %+v", node)
	}
}

func (g *Gen) genConst(node syntree.Node) {
	comment := fmt.Sprintf("load constant (%d) directly into ac", node.Value())
	g.emitRM("LDC", ac, node.Value(), 0, comment)
}

func (g *Gen) genOp(node syntree.Node) {
	if len(node.Children()) > 0 {
		n0 := node.Children()[0]
		g.gen(n0)
	}
	if len(node.Children()) > 1 {
		n1 := node.Children()[1]
		g.gen(n1)
	}

	switch node.TokType() {
	case syntree.PLUS:
		g.emitRO("ADD", ac, ac1, ac, "op +")
	case syntree.MINUS:
		g.emitRO("SUB", ac, ac1, ac, "op -")
	case syntree.TIMES:
		g.emitRO("MUL", ac, ac1, ac, "op *")
	case syntree.OVER:
		g.emitRO("DIV", ac, ac1, ac, "op -")
	case syntree.EQ:
		g.emitRO("SUB", ac, ac1, ac, "op substract")
		g.emitRM("JEQ", ac, 2, pc, "branch if true")
		g.emitRM("LDC", ac, 0, 0, "load constant 0 into ac (false)")
		g.emitRM("LDA", pc, 1, pc, "unconditional jump 1")
		g.emitRM("LDC", ac, 1, 0, "load constant 1 into ac (true)")
	case syntree.NEQ:
		g.emitRO("SUB", ac, ac1, ac, "op substract")
		g.emitRM("JNE", ac, 2, pc, "branch if true")
		g.emitRM("LDC", ac, 0, 0, "load constant 0 into ac (false)")
		g.emitRM("LDA", pc, 1, pc, "unconditional jump 1")
		g.emitRM("LDC", ac, 1, 0, "load constant 1 into ac (true)")
	case syntree.LT:
		g.emitRO("SUB", ac, ac1, ac, "op substract")
		g.emitRM("JLT", ac, 2, pc, "branch if true")
		g.emitRM("LDC", ac, 0, 0, "load constant 0 into ac (false)")
		g.emitRM("LDA", pc, 1, pc, "unconditional jump 1")
		g.emitRM("LDC", ac, 1, 0, "load constant 1 into ac (true)")
	case syntree.LTE:
		g.emitRO("SUB", ac, ac1, ac, "op substract")
		g.emitRM("JLE", ac, 2, pc, "branch if true")
		g.emitRM("LDC", ac, 0, 0, "load constant 0 into ac (false)")
		g.emitRM("LDA", pc, 1, pc, "unconditional jump 1")
		g.emitRM("LDC", ac, 1, 0, "load constant 1 into ac (true)")
	case syntree.GT:
		g.emitRO("SUB", ac, ac1, ac, "op substract")
		g.emitRM("JGT", ac, 2, pc, "branch if true")
		g.emitRM("LDC", ac, 0, 0, "load constant 0 into ac (false)")
		g.emitRM("LDA", pc, 1, pc, "unconditional jump 1")
		g.emitRM("LDC", ac, 1, 0, "load constant 1 into ac (true)")
	case syntree.GTE:
		g.emitRO("SUB", ac, ac1, ac, "op substract")
		g.emitRM("JGE", ac, 2, pc, "branch if true")
		g.emitRM("LDC", ac, 0, 0, "load constant 0 into ac (false)")
		g.emitRM("LDA", pc, 1, pc, "unconditional jump 1")
		g.emitRM("LDC", ac, 1, 0, "load constant 1 into ac (true)")
	default:
		log.ErrorLog.Printf("unknown operator type %s", node.TokType())
	}
}

func (g *Gen) genId(node syntree.Node) {
	/*
		comment := fmt.Sprintf("load %s with %d", node.Name(), node.Value())
		if symtbl.GlbSymTblMap[node.SymKey()].HasId(node.Name()) {
			memLoc := symtbl.GlbSymTblMap[node.SymKey()].GetMemLoc(node.Name())
			g.emitRM("LD", g.ac, int(memLoc), g.gp, comment)
		} else {
			log.ErrorLog.Printf(">>>>> Error %s not found.", node.Name())
		}
	*/

	log.CodeLog.Printf("searching for %s", node.Name())
	if symtbl.GlbSymTblMap[node.SymKey()].HasId(node.Name()) {
		memLoc := symtbl.GlbSymTblMap[node.SymKey()].GetMemLoc(node.Name())
		log.CodeLog.Printf("found %s %+v", node.Name(), memLoc)
	}

	if node.IsArray() {
		if len(node.Children()) > 0 {
			n0 := node.Children()[0]
			g.gen(n0)
		}
	} else {
	}
}

func (g *Gen) genParam(node syntree.Node) {
	if node.IsArray() {
	} else {
	}
}

func (g *Gen) genVar(node syntree.Node) {
	if node.IsArray() {
	} else {
	}
}

func (g *Gen) getPrelude() {
	g.emitComment("cminus compilation into tiny machine for " + g.filename)
	g.emitComment("prelude beg")
	g.emitRM("LD", gp, 0, ac, "load global pointer with maxaddress")
	g.emitRM("LDA", mp, 0, gp, "copy global pointer into memory pointer")
	g.emitRM("ST", ac, 0, ac, "clear location 0")
	g.emitComment("prelude end")
}

func (g *Gen) getHalt() {
	g.emitRO("HALT", 0, 0, 0, "halting program")
}

func (g *Gen) gen(node syntree.Node) {
	if node != nil {
		log.CodeLog.Printf("=> %+v", node)
		if node.IsCompound() {
			g.genCompound(node)
		} else if node.IsFunc() {
			g.genFunction(node)
		} else if node.IsIteration() {
			g.genIteration(node)
		} else if node.IsReturn() {
			g.genReturn(node)
		} else if node.IsSelection() {
			g.genSelection(node)
		} else if node.IsAssign() {
			g.genAssign(node)
		} else if node.IsCall() {
			g.genCall(node)
		} else if node.IsConst() {
			g.genConst(node)
		} else if node.IsOp() {
			g.genOp(node)
		} else if node.IsId() {
			g.genId(node)
		} else if node.IsParam() {
			g.genParam(node)
		} else if node.IsVar() {
			g.genVar(node)
		}
		log.CodeLog.Printf("<= %+v", node)
		g.gen(node.Sibling())
	}
}

func Generate(node syntree.Node, filename string) {
	g := NewGen(filename)
	if g != nil {
		g.getPrelude()
		g.gen(node)
		g.getHalt()
	} else {
		log.ErrorLog.Printf(">>>>> Error opening %s", filename)
	}
}
