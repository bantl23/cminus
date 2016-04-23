package gen

import (
	"fmt"
	"os"

	"github.com/bantl23/cminus/log"
	"github.com/bantl23/cminus/symtbl"
	"github.com/bantl23/cminus/syntree"
)

const (
	ofpFO  int = 0  // old frame pointer
	retFO  int = -1 // return address
	initFO int = -3 // param list
)

const (
	pc   int = 7 // program counter
	gp   int = 6 // global pointer
	fp   int = 5 // frame pointer
	sp   int = 4 // stack pointer
	zero int = 3 // zero
	ac2  int = 2 // accumlator
	ac1  int = 1 // accumlator
	ac   int = 0 // accumlator
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

func (g *Gen) emitPushSize(size int, comment string) {
	c := fmt.Sprintf("push stack (%d): %s", size, comment)
	g.emitRM("LDA", sp, -1*size, sp, c)
}

func (g *Gen) emitPush(comment string) {
	g.emitPushSize(1, comment)
}

func (g *Gen) emitPopSize(size int, comment string) {
	c := fmt.Sprintf("pop stack (%d): %s", size, comment)
	g.emitRM("LDA", sp, 1*size, sp, c)
}

func (g *Gen) emitPop(comment string) {
	g.emitPopSize(1, comment)
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
				size := 1
				if n0.IsArray() {
					size = n0.Value()
				}
				g.emitRM("LDC", ac, size, 0, "load "+n0.Name()+" length into ac")
				g.emitPushSize(size, "allocating vars")
			}
		}
		n0 = n0.Sibling()
	}
	g.gen(n1)

	n0 = node.Children()[0]
	for n0 != nil {
		if n0.IsVar() {
			if n0.ExpType() == syntree.INT_EXP_TYPE {
				size := 1
				if n0.IsArray() {
					size = n0.Value()
				}
				g.emitPopSize(size, "deallocating vars")
			}
		}
		n0 = n0.Sibling()
	}
}

func (g *Gen) genFunction(node syntree.Node) {
	memLoc := symtbl.MemLoc(0)
	if symtbl.GlbSymTblMap[symtbl.ROOT_KEY].HasId(node.Name()) {
		memLoc, _ = symtbl.GlbSymTblMap[symtbl.ROOT_KEY].GetMemLoc(node.Name())
		log.CodeLog.Printf("found %s at %+v offset from gp", node.Name(), memLoc)
	} else {
		log.ErrorLog.Printf("error could not find id")
	}

	n0 := node.Children()[0]
	n1 := node.Children()[1]

	sav0 := g.emitSkip(3)

	g.emitPush("allocating space for fp")
	g.emitRO("ADD", ac1, sp, zero, "save sp to ac1")
	g.emitRM("ST", fp, 0, sp, "save fp to sp")
	g.emitPush("allocating space for sp")
	g.emitRM("ST", sp, 0, sp, "save sp to sp")
	g.emitPush("allocating space for ret pc")
	g.emitRM("ST", ac, 0, sp, "save pc to sp")
	g.emitRO("ADD", fp, ac1, zero, "set fp to sp")

	for n0 != nil {
		if n0.IsParam() {
			if n0.ExpType() == syntree.INT_EXP_TYPE {
				g.emitPush("allocating params")
			}
		}
		n0 = n0.Sibling()
	}

	g.gen(n1)

	n0 = node.Children()[0]
	for n0 != nil {
		if n0.IsParam() {
			if n0.ExpType() == syntree.INT_EXP_TYPE {
				g.emitPop("deallocating params")
			}
		}
		n0 = n0.Sibling()
	}

	g.emitRM("LD", ac, 0, sp, "load pc on sp into ac")
	g.emitPop("deallocating space for fp")
	g.emitRM("LD", sp, 0, sp, "load sp on sp into sp")
	g.emitPop("deallocating space for sp")
	g.emitRM("LD", fp, 0, sp, "load fp on sp into fp")
	g.emitPop("dellocating space for ret pc")

	if node.Name() == "main" {
		halt := g.emitSkip(0) + 2
		g.emitRM("LDA", pc, halt, 0, "func: jump to halt")
	} else {
		g.emitRM("LDA", pc, 0, ac, "func: jump back to calling function")
	}

	sav1 := g.emitSkip(0)

	g.emitBackup(sav0)
	g.emitRM("LDC", ac1, sav0+3, 0, "save pc into ac1 for "+node.Name())
	g.emitRM("ST", ac1, 0-memLoc.Get(), gp, "saving ac1 for "+node.Name())
	g.emitRM("LDA", pc, sav1, 0, "func: jump to end")
	g.emitRestore()
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
	n0 := node.Children()[0]
	g.gen(n0)
	g.emitRM("LD", pc, retFO, fp, "return to caller")
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
	depth := 0
	if symtbl.GlbSymTblMap[n0.SymKey()].HasId(n0.Name()) {
		memLoc, depth = symtbl.GlbSymTblMap[n0.SymKey()].GetMemLoc(n0.Name())
		log.CodeLog.Printf("found %s at %+v offset from fp at depth %d", n0.Name(), memLoc, depth)
	} else {
		log.ErrorLog.Printf("error could not find id")
	}

	g.gen(n0)
	g.gen(n1)

	g.emitRM("LDA", ac2, 0, fp, "store fp into ac2")
	for i := 0; i < depth; i++ {
		g.emitRM("LD", ac2, 0, ac2, "get fp from previous scope")
	}

	if n0.IsArray() {
		g.emitPush("allocating space for assign ret val")
		g.emitRM("ST", ac, 0, sp, "storing value of array")
		nidx := n0.Children()[0]
		g.gen(nidx)
		g.emitRM("LDC", ac1, initFO-memLoc.Get(), 0, "load base address for array")
		g.emitRO("SUB", ac, ac1, ac, "get array offset")
		g.emitRO("ADD", ac, ac2, ac, "get array memory location")
		g.emitRM("LD", ac1, 0, sp, "loading value of array")
		g.emitPop("deallocating space for assign ret val")
		g.emitRM("ST", ac1, 0, ac, "load value into ac")
	} else {
		g.emitRM("ST", ac, initFO-memLoc.Get(), ac2, "store ac into id "+n0.Name())
	}
}

func (g *Gen) genCall(node syntree.Node) {
	memLoc := symtbl.MemLoc(0)
	if symtbl.GlbSymTblMap[symtbl.ROOT_KEY].HasId(node.Name()) {
		memLoc, _ = symtbl.GlbSymTblMap[symtbl.ROOT_KEY].GetMemLoc(node.Name())
		log.CodeLog.Printf("found %s at %+v offset from gp", node.Name(), memLoc)
	} else {
		log.ErrorLog.Printf("error could not find id")
	}
	n0 := node.Children()[0]
	g.gen(n0)

	if node.Name() == "input" {
		g.emitRO("IN", ac, 0, 0, "read from stdin into ac")
	} else if node.Name() == "output" {
		g.emitRO("OUT", ac, 0, 0, "write to stdout with ac")
	} else {
		ret := g.emitSkip(0) + 2
		g.emitRM("LDC", ac, ret, 0, "loading ret pc into ac")
		g.emitRM("LD", pc, 0-memLoc.Get(), gp, "func: jump func "+node.Name())
	}
}

func (g *Gen) genConst(node syntree.Node) {
	g.emitRM("LDC", ac, node.Value(), 0, "load constant directly into ac")
}

func (g *Gen) genOp(node syntree.Node) {
	n0 := node.Children()[0]
	n1 := node.Children()[1]

	g.gen(n0)
	g.emitPush("allocating space for op ret val")
	g.emitRM("ST", ac, 0, sp, "storing left hand side of operator from ac")

	g.gen(n1)
	g.emitRM("LD", ac1, 0, sp, "loading left hand side of operator into ac1")
	g.emitPop("deallocating space for op ret val")

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
		g.emitRO("SUB", ac, ac1, ac, "op substract")
		g.emitRM("JEQ", ac, 2, pc, "branch if true")
		g.emitRM("LDC", ac, 0, ac, "load constant 0 into ac (false)")
		g.emitRM("LDA", pc, 1, pc, "unconditional jump 1")
		g.emitRM("LDC", ac, 1, ac, "load constant 1 into ac (true)")
	case syntree.NEQ:
		g.emitRO("SUB", ac, ac1, ac, "op substract")
		g.emitRM("JNE", ac, 2, pc, "branch if true")
		g.emitRM("LDC", ac, 0, ac, "load constant 0 into ac (false)")
		g.emitRM("LDA", pc, 1, pc, "unconditional jump 1")
		g.emitRM("LDC", ac, 1, ac, "load constant 1 into ac (true)")
	case syntree.LT:
		g.emitRO("SUB", ac, ac1, ac, "op substract")
		g.emitRM("JLT", ac, 2, pc, "branch if true")
		g.emitRM("LDC", ac, 0, ac, "load constant 0 into ac (false)")
		g.emitRM("LDA", pc, 1, pc, "unconditional jump 1")
		g.emitRM("LDC", ac, 1, ac, "load constant 1 into ac (true)")
	case syntree.LTE:
		g.emitRO("SUB", ac, ac1, ac, "op substract")
		g.emitRM("JLE", ac, 2, pc, "branch if true")
		g.emitRM("LDC", ac, 0, ac, "load constant 0 into ac (false)")
		g.emitRM("LDA", pc, 1, pc, "unconditional jump 1")
		g.emitRM("LDC", ac, 1, ac, "load constant 1 into ac (true)")
	case syntree.GT:
		g.emitRO("SUB", ac, ac1, ac, "op substract")
		g.emitRM("JGT", ac, 2, pc, "branch if true")
		g.emitRM("LDC", ac, 0, ac, "load constant 0 into ac (false)")
		g.emitRM("LDA", pc, 1, pc, "unconditional jump 1")
		g.emitRM("LDC", ac, 1, ac, "load constant 1 into ac (true)")
	case syntree.GTE:
		g.emitRO("SUB", ac, ac1, ac, "op substract")
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
	depth := 0
	if symtbl.GlbSymTblMap[node.SymKey()].HasId(node.Name()) {
		memLoc, depth = symtbl.GlbSymTblMap[node.SymKey()].GetMemLoc(node.Name())
		log.CodeLog.Printf("found %s at %+v offset from fp at depth %d", node.Name(), memLoc, depth)
	} else {
		log.ErrorLog.Printf("error could not find id")
	}

	g.emitRM("LDA", ac2, 0, fp, "store fp into ac2")
	for i := 0; i < depth; i++ {
		g.emitRM("LD", ac2, 0, ac2, "get fp from previous scope")
	}

	if node.IsArray() {
		n0 := node.Children()[0]
		g.gen(n0)
		g.emitRM("LDC", ac1, initFO-memLoc.Get(), 0, "load base address for array")
		g.emitRO("SUB", ac, ac1, ac, "get array offset")
		g.emitRO("ADD", ac, ac2, ac, "get array memory location")
		g.emitRM("LD", ac, 0, ac, "load value into ac")
	} else {
		g.emitRM("LD", ac, initFO-memLoc.Get(), ac2, "store ac with id "+node.Name())
	}
}

func (g *Gen) genParam(node syntree.Node) {
	if node.ExpType() == syntree.INT_EXP_TYPE {
		size := 1
		if node.IsArray() {
			if node.IsArray() {
				size = node.Value()
			}
		}
		g.emitComment("pushing " + node.Name() + " param into stack")
		g.emitPushSize(size, "allocating space for params")
	}
}

func (g *Gen) genVar(node syntree.Node) {
}

func (g *Gen) genPrelude() {
	g.emitComment("cminus compilation into tiny machine for " + g.filename)
	g.emitComment("prelude beg")
	g.emitRM("LD", gp, 0, ac, "load global pointer with maxaddress")
	g.emitRM("LDA", fp, 0, gp, "copy global pointer into frame pointer")
	g.emitRM("LDA", sp, 1, gp, "copy global pointer into stack pointer")
	g.emitRM("ST", ac, 0, ac, "clear location 0")
	g.emitRM("LDC", zero, 0, 0, "set zero")
	g.emitComment("prelude end")
}

func (g *Gen) genGlobals() {
	g.emitComment("global space beg")
	g.emitPushSize(3, "pushing blank spaces so scoping will work consistantly")
	glbs := symtbl.GlbSymTblMap[symtbl.ROOT_KEY].SymTbl()
	for k, v := range glbs {
		size := v.Size()
		comment := fmt.Sprintf("allocating space for %s (%d)", k, size)
		g.emitPushSize(size, comment)
	}
	g.emitComment("global space end")
}

func (g *Gen) genMain() {
	memLoc := symtbl.MemLoc(0)
	if symtbl.GlbSymTblMap[symtbl.ROOT_KEY].HasId("main") {
		memLoc, _ = symtbl.GlbSymTblMap[symtbl.ROOT_KEY].GetMemLoc("main")
		log.CodeLog.Printf("found main at %+v offset from gp", memLoc)
	} else {
		log.ErrorLog.Printf("error could not find id")
	}

	g.emitRM("LD", pc, 0-memLoc.Get(), gp, "jumping to main")
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
		g.genPrelude()
		g.genGlobals()
		g.gen(node)
		g.genMain()
		g.getHalt()
	} else {
		log.ErrorLog.Printf(">>>>> Error opening %s", filename)
	}
}
