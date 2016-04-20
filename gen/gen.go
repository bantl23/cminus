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
				g.emitRM("LDC", ac, length, 0, "load "+n0.Name()+" length into scratch")
			}
		}
		n0 = n0.Sibling()
	}
	g.gen(n1)
}

func (g *Gen) genFunction(node syntree.Node) {
	g.emitComment("=> func: " + node.Name())
	n0 := node.Children()[0]
	n1 := node.Children()[1]

	g.emitRM("ST", fp, 0, fp, "store frame pointer")

	g.emitRM("ST", fp, 0, fp, "store control link")

	g.emitComment("func: generate declaration here")
	g.gen(n0)
	g.emitComment("func: generate body here")
	g.gen(n1)

	g.emitComment("<= func: " + node.Name())
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

	// left hand of assign
	memLoc := symtbl.MemLoc(0)
	if symtbl.GlbSymTblMap[n0.SymKey()].HasId(n0.Name()) {
		memLoc = symtbl.GlbSymTblMap[n0.SymKey()].GetMemLoc(n0.Name())
		log.CodeLog.Printf("found %s at %+v offset from fp", n0.Name(), memLoc)
	} else {
		log.ErrorLog.Printf("error could not find id")
	}

	n0.SetLeft(true)
	g.gen(n0)

	// right hand of assign
	n1.SetLeft(false)
	g.gen(n1)

	g.emitRM("ST", ac1, initFO-memLoc.Get(), fp, "store ac1 to id "+n0.Name())
}

func (g *Gen) genCall(node syntree.Node) {
	n0 := node.Children()[0]
	log.CodeLog.Printf("%+v", n0)
	g.gen(n0)

	if node.Name() == "input" {
		g.emitRO("IN", ac, 0, 0, "read from stdin into ac")
	} else if node.Name() == "output" {
		g.emitRO("OUT", ac1, 0, 0, "write to stdout with ac1")
	} else {
		// TODO
		log.CodeLog.Printf("TODO %+v", node)
	}
}

func (g *Gen) genConst(node syntree.Node) {
	scratch := ac1
	if node.IsLeft() {
		scratch = ac
	}
	comment := fmt.Sprintf("load constant (%d) directly into scratch", node.Value())
	g.emitRM("LDC", scratch, node.Value(), 0, comment)
}

func (g *Gen) genOp(node syntree.Node) {
	n0 := node.Children()[0]
	n1 := node.Children()[1]

	n0.SetLeft(true)
	g.gen(n0)

	n1.SetLeft(false)
	g.gen(n1)

	scratch := ac1
	if node.IsLeft() {
		scratch = ac
	}

	switch node.TokType() {
	case syntree.PLUS:
		g.emitRO("ADD", scratch, ac, ac1, "op + [ac = ac + ac1]")
	case syntree.MINUS:
		g.emitRO("SUB", scratch, ac, ac1, "op - [ac = ac - ac1]")
	case syntree.TIMES:
		g.emitRO("MUL", scratch, ac, ac1, "op * [ac = ac * ac1]")
	case syntree.OVER:
		g.emitRO("DIV", scratch, ac, ac1, "op - [ac = ac / ac1]")
	case syntree.EQ:
		g.emitRO("SUB", scratch, ac, ac1, "op substract")
		g.emitRM("JEQ", scratch, 2, pc, "branch if true")
		g.emitRM("LDC", scratch, 0, scratch, "load constant 0 into scratch (false)")
		g.emitRM("LDA", pc, 1, pc, "unconditional jump 1")
		g.emitRM("LDC", scratch, 1, scratch, "load constant 1 into scratch (true)")
	case syntree.NEQ:
		g.emitRO("SUB", scratch, ac, ac1, "op substract")
		g.emitRM("JNE", scratch, 2, pc, "branch if true")
		g.emitRM("LDC", scratch, 0, scratch, "load constant 0 into scratch (false)")
		g.emitRM("LDA", pc, 1, pc, "unconditional jump 1")
		g.emitRM("LDC", scratch, 1, scratch, "load constant 1 into scratch (true)")
	case syntree.LT:
		g.emitRO("SUB", scratch, ac, ac1, "op substract")
		g.emitRM("JLT", scratch, 2, pc, "branch if true")
		g.emitRM("LDC", scratch, 0, scratch, "load constant 0 into scratch (false)")
		g.emitRM("LDA", pc, 1, pc, "unconditional jump 1")
		g.emitRM("LDC", scratch, 1, scratch, "load constant 1 into scratch (true)")
	case syntree.LTE:
		g.emitRO("SUB", scratch, ac, ac1, "op substract")
		g.emitRM("JLE", scratch, 2, pc, "branch if true")
		g.emitRM("LDC", scratch, 0, scratch, "load constant 0 into scratch (false)")
		g.emitRM("LDA", pc, 1, pc, "unconditional jump 1")
		g.emitRM("LDC", scratch, 1, scratch, "load constant 1 into scratch (true)")
	case syntree.GT:
		g.emitRO("SUB", scratch, ac, ac1, "op substract")
		g.emitRM("JGT", scratch, 2, pc, "branch if true")
		g.emitRM("LDC", scratch, 0, scratch, "load constant 0 into scratch (false)")
		g.emitRM("LDA", pc, 1, pc, "unconditional jump 1")
		g.emitRM("LDC", scratch, 1, scratch, "load constant 1 into scratch (true)")
	case syntree.GTE:
		g.emitRO("SUB", scratch, ac, ac1, "op substract")
		g.emitRM("JGE", scratch, 2, pc, "branch if true")
		g.emitRM("LDC", scratch, 0, scratch, "load constant 0 into scratch (false)")
		g.emitRM("LDA", pc, 1, pc, "unconditional jump 1")
		g.emitRM("LDC", scratch, 1, scratch, "load constant 1 into scratch (true)")
	default:
		log.ErrorLog.Printf("unknown operator type %s", node.TokType())
	}
}

func (g *Gen) genId(node syntree.Node) {
	memLoc := symtbl.MemLoc(0)
	if symtbl.GlbSymTblMap[node.SymKey()].HasId(node.Name()) {
		memLoc = symtbl.GlbSymTblMap[node.SymKey()].GetMemLoc(node.Name())
		log.CodeLog.Printf("found %s at %+v offset from fp", node.Name(), memLoc)
	} else {
		log.ErrorLog.Printf("error could not find id")
	}

	scratch := ac1
	if node.IsLeft() {
		scratch = ac
	}
	comment := fmt.Sprintf("loading %s into scratch", node.Name())
	g.emitRM("LD", scratch, initFO-memLoc.Get(), fp, comment)
}

func (g *Gen) genParam(node syntree.Node) {
	if node.ExpType() == syntree.INT_EXP_TYPE {
		length := 1
		if node.IsArray() {
			length = node.Value()
		}
		if length != 0 {
			g.emitRM("LDC", ac, length, 0, "load "+node.Name()+" length into scratch")
		} else {
			// TODO
			log.CodeLog.Printf("assign a reference")
		}
	}
}

func (g *Gen) genVar(node syntree.Node) {
	if node.ExpType() == syntree.INT_EXP_TYPE {
		length := 1
		if node.IsArray() {
			length = node.Value()
		}
		g.emitRM("LDC", ac, length, 0, "load "+node.Name()+" length into scratch")
		g.emitRO("SUB", fp, fp, ac, "move frame pointer to allocate space for var "+node.Name())
	}
}

func (g *Gen) getPrelude() {
	g.emitComment("cminus compilation into tiny machine for " + g.filename)
	g.emitComment("prelude beg")
	g.emitRM("LD", gp, 0, ac, "load global pointer with maxaddress")
	g.emitRM("LDA", fp, 0, gp, "copy global pointer into frame pointer")
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
