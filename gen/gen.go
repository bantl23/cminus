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
	initFO int = -2 // param list
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

var followSibling bool = true

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
	out := fmt.Sprintf("skipping %d amount", amount)
	g.emitComment(out)
	i := g.loc
	g.loc = g.loc + amount
	if g.highLoc < g.loc {
		g.highLoc = g.loc
	}
	return i
}

func (g *Gen) emitBackup(loc int) {
	out := fmt.Sprintf("backing up to %d", loc)
	g.emitComment(out)
	if loc > g.highLoc {
		log.ErrorLog.Printf(">>>>> Error in emitBackup\n")
	}
	g.loc = loc
}

func (g *Gen) emitRestore() {
	out := fmt.Sprintf("restoring to %d", g.highLoc)
	g.emitComment(out)
	g.loc = g.highLoc
}

func (g Gen) getInfo(key string, symbol string) (int, int, symtbl.ReturnType) {
	memLoc := symtbl.MemLoc(0)
	depth := 0
	retType := symtbl.UNK_RET_TYPE
	if symtbl.GlbSymTblMap[key].HasId(symbol) {
		memLoc, depth = symtbl.GlbSymTblMap[key].GetMemLoc(symbol)
		retType = symtbl.GlbSymTblMap[key].GetReturnType(symbol)
		out := fmt.Sprintf("found %s at %+v offset from fp at depth %d [%s]", symbol, memLoc, depth, key)
		g.emitComment(out)
	} else {
		log.ErrorLog.Printf("error could not find id")
	}
	return memLoc.Get(), depth, retType
}

func (g *Gen) getAddress(node syntree.Node) {
	memLoc, depth, _ := g.getInfo(node.SymKey(), node.Name())

	g.emitRM("LDA", ac, 0, fp, "store fp")
	for i := 0; i < depth; i++ {
		g.emitRM("LD", ac, 0, ac, "store scoped fp")
	}
	g.emitRM("LDC", ac1, initFO-memLoc, zero, "store offset from scoped fp")
	g.emitRO("ADD", ac, ac, ac1, "store address of "+node.Name())

	if node.IsArray() {
		n0 := node.Children()[0]
		g.gen(n0)
		g.emitRO("SUB", ac, ac, ac1, "store array address of "+node.Name())
	}
}

func (g *Gen) getValue(node syntree.Node) {
	g.getAddress(node)
	g.emitRM("LD", ac, 0, ac, "store value of "+node.Name())
}

func (g *Gen) genCompound(node syntree.Node) {
	n0 := node.Children()[0]
	n1 := node.Children()[1]

	for n0 != nil {
		if n0.IsVar() && n0.ExpType() == syntree.INT_EXP_TYPE {
			size := 1
			if n0.IsArray() {
				size = n0.Value()
			}
			g.emitPushSize(size, "allocate var "+n0.Name())
		}
		n0 = n0.Sibling()
	}
	g.gen(n1)

	n0 = node.Children()[0]
	for n0 != nil {
		if n0.IsVar() && n0.ExpType() == syntree.INT_EXP_TYPE {
			size := 1
			if n0.IsArray() {
				size = n0.Value()
			}
			g.emitPopSize(size, "deallocate var "+n0.Name())
		}
		n0 = n0.Sibling()
	}
}

func (g *Gen) genFunction(node syntree.Node) {
	memLoc, _, _ := g.getInfo(symtbl.ROOT_KEY, node.Name())

	n0 := node.Children()[0]
	n1 := node.Children()[1]

	sav0 := g.emitSkip(3)

	g.emitRM("LD", ac, 0, sp, "load return pc")
	g.emitPop("deallocate return pc")
	g.emitRO("ADD", ac1, sp, zero, "save sp to to populate params")
	g.emitPush("allocate space for old fp")
	g.emitRM("ST", fp, 0, sp, "save old fp")
	g.emitRO("ADD", fp, sp, zero, "update fp to new frame")
	g.emitPush("allocate space for return pc")
	g.emitRM("ST", ac, 0, sp, "save return pc")

	length := 0
	n0 = node.Children()[0]
	for n0 != nil {
		if n0.ExpType() == syntree.INT_EXP_TYPE {
			g.emitPush("allocate param " + n0.Name())
			length = length + 1
		}
		n0 = n0.Sibling()
	}

	offset := 0
	n0 = node.Children()[0]
	for n0 != nil {
		if n0.ExpType() == syntree.INT_EXP_TYPE {
			if n0.IsArray() {
				g.emitRM("LDA", ac, length-offset-1, ac1, "load param value")
				g.emitRM("ST", ac, initFO-offset, fp, "store param value")
			} else {
				g.emitRM("LD", ac, length-offset-1, ac1, "load param array value")
				g.emitRM("ST", ac, initFO-offset, fp, "store param array value")
			}
			offset = offset + 1
		}
		n0 = n0.Sibling()
	}

	g.gen(n1)

	n0 = node.Children()[0]
	for n0 != nil {
		if n0.ExpType() == syntree.INT_EXP_TYPE {
			g.emitPop("deallocate param " + n0.Name())
		}
		n0 = n0.Sibling()
	}

	g.emitRM("LD", ac, 0, sp, "load return pc")
	g.emitPop("deallocate return pc")
	g.emitRM("LD", fp, 0, sp, "load old fp")
	g.emitPop("dellocate old fp")

	if node.Name() == "main" {
		halt := g.emitSkip(0) + 4
		g.emitRM("LDA", pc, halt, zero, "func: jump to halt")
	} else {
		g.emitRM("LDA", pc, 0, ac, "func: jump back to calling function")
	}

	sav1 := g.emitSkip(0)

	g.emitBackup(sav0)
	g.emitRM("LDC", ac, sav0+3, zero, "save pc loc of "+node.Name())
	g.emitRM("ST", ac, 0-memLoc, gp, "store pc loc of "+node.Name())
	g.emitRM("LDA", pc, sav1, zero, "func: jump to end")
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
	g.emitRM("LD", sp, -1, sp, "set sp to end of old activation")
	g.emitRM("LD", fp, 0, fp, "set fp to old fp")
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

	g.gen(n1)
	g.emitPush("allocate space for right hand assign")
	g.emitRM("ST", ac, 0, sp, "store right hand assign")

	g.getAddress(n0)
	g.emitRM("LD", ac1, 0, sp, "load right hand assign")
	g.emitPop("deallocate space for right hand assign")

	g.emitRM("ST", ac1, 0, ac, "store result into id "+n0.Name())
}

func (g *Gen) genCall(node syntree.Node) {
	memLoc, _, _ := g.getInfo(symtbl.ROOT_KEY, node.Name())

	if node.Name() == "input" {
		g.emitRO("IN", ac, 0, 0, "read from stdin into ac")
	} else if node.Name() == "output" {
		n0 := node.Children()[0]
		g.gen(n0)
		g.emitRO("OUT", ac, 0, 0, "write to stdout from ac")
	} else {

		followSibling = false
		n0 := node.Children()[0]
		for n0 != nil {
			g.gen(n0)
			g.emitPush("allocate arg")
			comment := fmt.Sprintf("store arg %s (%d)", n0.Name(), n0.Value())
			g.emitRM("ST", ac, 0, sp, comment)
			n0 = n0.Sibling()
		}
		followSibling = true

		ret := g.emitSkip(0) + 4
		g.emitRM("LDC", ac, ret, 0, "load return pc")
		g.emitPush("allocate space for return pc")
		g.emitRM("ST", ac, 0, sp, "store return pc onto stack")
		g.emitRM("LD", pc, 0-memLoc, gp, "func: jump func "+node.Name())

		followSibling = false
		n0 = node.Children()[0]
		for n0 != nil {
			comment := fmt.Sprintf("deallocate arg %s (%d)", n0.Name(), n0.Value())
			g.emitPop(comment)
			n0 = n0.Sibling()
		}
		followSibling = true
	}
}

func (g *Gen) genConst(node syntree.Node) {
	g.emitRM("LDC", ac, node.Value(), 0, "load const into ac")
}

func (g *Gen) genOp(node syntree.Node) {
	n0 := node.Children()[0]
	n1 := node.Children()[1]

	g.gen(n0)
	g.emitPush("allocate space for left op")
	g.emitRM("ST", ac, 0, sp, "store left of right op")

	g.gen(n1)
	g.emitRM("LD", ac1, 0, sp, "load result of left op")
	g.emitPop("deallocate space for left op")

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
}

func (g *Gen) genId(node syntree.Node) {
	g.getValue(node)
}

func (g *Gen) genParam(node syntree.Node) {
}

func (g *Gen) genVar(node syntree.Node) {
}

func (g *Gen) genPrelude() {
	g.emitComment("cminus compilation into tiny machine for " + g.filename)
	g.emitComment("prelude beg")
	g.emitRM("LD", gp, 0, ac, "load gp with maxaddress")
	g.emitRM("LDA", fp, 0, gp, "copy gp into fp")
	g.emitRM("LDA", sp, 1, gp, "copy gp into sp")
	g.emitRM("ST", ac, 0, ac, "clear location 0")
	g.emitRM("LDC", zero, 0, 0, "set zero")
	g.emitRM("LDC", ac, 0, 0, "clear ac")
	g.emitRM("LDC", ac1, 0, 0, "clear ac1")
	g.emitRM("LDC", ac2, 0, 0, "clear ac2")
	g.emitComment("prelude end")
}

func (g *Gen) genGlobals() {
	g.emitComment("global space beg")
	glbs := symtbl.GlbSymTblMap[symtbl.ROOT_KEY].SymTbl()
	for k, v := range glbs {
		size := v.Size()
		comment := fmt.Sprintf("allocating space for %s (%d)", k, size)
		g.emitPushSize(size, comment)
	}
	g.emitComment("global space end")
}

func (g *Gen) genMain() {
	g.emitComment("main beg")
	memLoc, _, _ := g.getInfo(symtbl.ROOT_KEY, "main")
	g.emitPush("allocate space for fake return pc")
	g.emitRM("ST", zero, 0, sp, "store fake return pc onto stack")
	g.emitRM("LD", pc, 0-memLoc, gp, "jumping to main")
	g.emitComment("main end")
}

func (g *Gen) getHalt() {
	g.emitRO("HALT", 0, 0, 0, "halting program")
}

func (g *Gen) gen(node syntree.Node) {
	if node != nil {
		out := fmt.Sprintf("=> %+v", node)
		g.emitComment(out)
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
		out = fmt.Sprintf("<= %+v", node)
		g.emitComment(out)
		if followSibling == true {
			g.gen(node.Sibling())
		}
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
