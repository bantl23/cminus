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
	gp  int = 6 // global pointer
	fp  int = 5 // frame pointer
	sp  int = 4 // stack pointer
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

func (g *Gen) emitRMAbs(opcode string, target int, absolute int, comment string) {
	out := fmt.Sprintf("%3d: %5s %d,%d(%d)\t* %s\n", g.loc, opcode, target, absolute-(g.loc+1), pc, comment)
	g.emit(out)
	g.loc = g.loc + 1
	if g.highLoc < g.loc {
		g.highLoc = g.loc
	}
}

func (g *Gen) emitPush() {
	g.emitRM("LDA", sp, -1, sp, "push stack")
}

func (g *Gen) emitPushSize(size int) {
	g.emitRM("LDA", sp, -1*size, sp, "push stack")
}

func (g *Gen) emitPop() {
	g.emitRM("LDA", sp, 1, sp, "pop stack")
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
	n0 := node.Children()[0]
	n1 := node.Children()[1]

	for n0 != nil {
		if n0.IsVar() {
			if n0.ExpType() == syntree.INT_EXP_TYPE {
				length := 1
				if n0.IsArray() {
					length = n0.Value()
				}
				g.emitRM("LDC", ac, length, 0, "load "+n0.Name()+" length into ac")
				g.emitRO("SUB", sp, sp, 1, "push stack pointer")
			}
		}
		n0 = n0.Sibling()
	}
	g.gen(n1)
}

func (g *Gen) genFunction(node syntree.Node) {
	n0 := node.Children()[0]
	n1 := node.Children()[1]

	if node.Name() != "main" {
		g.emitPush()
	}
	g.emitRM("ST", fp, 0, fp, "save fp to stack")
	g.emitPush()
	g.emitRM("ST", fp, 0, fp, "save cl to stack")

	g.gen(n0)
	g.gen(n1)
}

func (g *Gen) genIteration(node syntree.Node) {
	n0 := node.Children()[0]
	n1 := node.Children()[1]

	sav0 := g.emitSkip(0)
	g.emitComment("while: jump after body comes back here")
	g.gen(n0)

	sav1 := g.emitSkip(1)
	g.emitComment("while: jump to end belongs here")
	g.gen(n1)
	g.emitRMAbs("LDA", pc, sav0, "while: jump back to body")

	curr := g.emitSkip(0)
	g.emitBackup(sav1)
	g.emitRMAbs("JEQ", ac, curr, "while: jump to end")
	g.emitRestore()
}

func (g *Gen) genReturn(node syntree.Node) {
	//TODO
	n0 := node.Children()[0]
	g.gen(n0)
	//g.emitRM("LD", pc, retFO, mp, "return to caller")
}

func (g *Gen) genSelection(node syntree.Node) {
	n0 := node.Children()[0]
	n1 := node.Children()[1]

	g.gen(n0)
	sav0 := g.emitSkip(1)
	g.emitComment("if: jump to else belongs here")

	g.gen(n1)
	sav1 := g.emitSkip(1)
	g.emitComment("if: jump to end belongs here")

	curr := g.emitSkip(0)
	g.emitBackup(sav0)
	g.emitRMAbs("JEQ", ac, curr, "if: jump to else")
	g.emitRestore()

	if len(node.Children()) == 3 {
		n2 := node.Children()[2]
		g.gen(n2)
		curr = g.emitSkip(0)
		g.emitBackup(sav1)
		g.emitRMAbs("LDA", pc, curr, "if: jump to end")
		g.emitRestore()
	}
}

func (g *Gen) genAssign(node syntree.Node) {
	n0 := node.Children()[0]
	n1 := node.Children()[1]

	memLoc := symtbl.MemLoc(0)
	if symtbl.GlbSymTblMap[n0.SymKey()].HasId(n0.Name()) {
		memLoc = symtbl.GlbSymTblMap[n0.SymKey()].GetMemLoc(n0.Name())
		log.CodeLog.Printf("found %s at %+v offset from fp", n0.Name(), memLoc)
	} else {
		log.ErrorLog.Printf("error could not find id")
	}

	g.gen(n0)
	g.gen(n1)

	g.emitRM("ST", ac, initFO-memLoc.Get(), fp, "store ac into id "+n0.Name())
}

func (g *Gen) genCall(node syntree.Node) {
	n0 := node.Children()[0]
	log.CodeLog.Printf("%+v", n0)
	g.gen(n0)

	if node.Name() == "input" {
		g.emitRO("IN", ac, 0, 0, "read from stdin into ac")
	} else if node.Name() == "output" {
		g.emitRO("OUT", ac, 0, 0, "write to stdout with ac")
	} else {
		// TODO
		log.CodeLog.Printf("TODO %+v", node)
	}
}

func (g *Gen) genConst(node syntree.Node) {
	g.emitRM("LDC", ac, node.Value(), 0, "load constant directly into ac")
}

func (g *Gen) genOp(node syntree.Node) {
	n0 := node.Children()[0]
	n1 := node.Children()[1]

	g.gen(n0)
	g.emitPush()
	g.emitRM("ST", ac, 0, sp, "storing left hand side of operator from ac")

	g.gen(n1)
	g.emitRM("LD", ac1, 0, sp, "loading left hand side of operator into ac1")
	g.emitPop()

	switch node.TokType() {
	case syntree.PLUS:
		g.emitRO("ADD", ac, ac1, ac, "op + [ac = ac1 + ac]")
	case syntree.MINUS:
		g.emitRO("SUB", ac, ac1, ac, "op - [ac = ac1 - ac]")
	case syntree.TIMES:
		g.emitRO("MUL", ac, ac1, ac, "op * [ac = ac1 * ac]")
	case syntree.OVER:
		g.emitRO("DIV", ac, ac1, ac, "op - [ac = ac1 / ac]")
	case syntree.EQ:
		g.emitRO("SUB", ac, ac, ac1, "op substract")
		g.emitRM("JEQ", ac, 2, pc, "branch if true")
		g.emitRM("LDC", ac, 0, ac, "load constant 0 into ac (false)")
		g.emitRM("LDA", pc, 1, pc, "unconditional jump 1")
		g.emitRM("LDC", ac, 1, ac, "load constant 1 into ac (true)")
	case syntree.NEQ:
		g.emitRO("SUB", ac, ac, ac1, "op substract")
		g.emitRM("JNE", ac, 2, pc, "branch if true")
		g.emitRM("LDC", ac, 0, ac, "load constant 0 into ac (false)")
		g.emitRM("LDA", pc, 1, pc, "unconditional jump 1")
		g.emitRM("LDC", ac, 1, ac, "load constant 1 into ac (true)")
	case syntree.LT:
		g.emitRO("SUB", ac, ac, ac1, "op substract")
		g.emitRM("JLT", ac, 2, pc, "branch if true")
		g.emitRM("LDC", ac, 0, ac, "load constant 0 into ac (false)")
		g.emitRM("LDA", pc, 1, pc, "unconditional jump 1")
		g.emitRM("LDC", ac, 1, ac, "load constant 1 into ac (true)")
	case syntree.LTE:
		g.emitRO("SUB", ac, ac, ac1, "op substract")
		g.emitRM("JLE", ac, 2, pc, "branch if true")
		g.emitRM("LDC", ac, 0, ac, "load constant 0 into ac (false)")
		g.emitRM("LDA", pc, 1, pc, "unconditional jump 1")
		g.emitRM("LDC", ac, 1, ac, "load constant 1 into ac (true)")
	case syntree.GT:
		g.emitRO("SUB", ac, ac, ac1, "op substract")
		g.emitRM("JGT", ac, 2, pc, "branch if true")
		g.emitRM("LDC", ac, 0, ac, "load constant 0 into ac (false)")
		g.emitRM("LDA", pc, 1, pc, "unconditional jump 1")
		g.emitRM("LDC", ac, 1, ac, "load constant 1 into ac (true)")
	case syntree.GTE:
		g.emitRO("SUB", ac, ac, ac1, "op substract")
		g.emitRM("JGE", ac, 2, pc, "branch if true")
		g.emitRM("LDC", ac, 0, ac, "load constant 0 into ac (false)")
		g.emitRM("LDA", pc, 1, pc, "unconditional jump 1")
		g.emitRM("LDC", ac, 1, ac, "load constant 1 into ac (true)")
	default:
		log.ErrorLog.Printf("unknown operator type %s", node.TokType())
	}
	g.emitRM("ST", ac, 0, sp, "storing operation result into ac")
}

func (g *Gen) genId(node syntree.Node) {
	memLoc := symtbl.MemLoc(0)
	if symtbl.GlbSymTblMap[node.SymKey()].HasId(node.Name()) {
		memLoc = symtbl.GlbSymTblMap[node.SymKey()].GetMemLoc(node.Name())
		log.CodeLog.Printf("found %s at %+v offset from fp", node.Name(), memLoc)
	} else {
		log.ErrorLog.Printf("error could not find id")
	}

	g.emitRM("LD", ac, initFO-memLoc.Get(), fp, "store ac with id "+node.Name())
}

func (g *Gen) genParam(node syntree.Node) {
	// TODO array
	if node.ExpType() == syntree.INT_EXP_TYPE {
		g.emitComment("pushing " + node.Name() + " param into stack")
		g.emitPush()
	}
}

func (g *Gen) genVar(node syntree.Node) {
	// TODO array
	if node.ExpType() == syntree.INT_EXP_TYPE {
		g.emitComment("pushing " + node.Name() + " var into stack")
		g.emitPush()
	}
}

func (g *Gen) getPrelude() {
	g.emitComment("cminus compilation into tiny machine for " + g.filename)
	g.emitComment("prelude beg")
	g.emitRM("LD", gp, 0, ac, "load global pointer with maxaddress")
	g.emitRM("LDA", fp, 0, gp, "copy global pointer into frame pointer")
	g.emitRM("LDA", sp, 0, gp, "copy global pointer into stack pointer")
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
