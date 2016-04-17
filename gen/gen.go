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
	g.emitRM("LD", g.gp, 0, g.ac, "load maxaddress from location 0")
	g.emitRM("LDA", g.mp, 0, g.gp, "copy gp to mp")
	g.emitRM("ST", g.ac, 0, g.ac, "clear location 0")
	g.emitComment("prelude end")
}

func (g *Gen) halt() {
	g.emitRO("HALT", 0, 0, 0, "halting program")
}

func (g *Gen) genCompound(node syntree.Node) {
	log.CodeLog.Printf("=> compound %+v", node)
	for _, n := range node.Children() {
		g.gen(n)
	}
	log.CodeLog.Printf("<= compound %+v", node)
}

func (g *Gen) genFunction(node syntree.Node) {
	log.CodeLog.Printf("=> function %+v", node)
	for _, n := range node.Children() {
		g.gen(n)
	}
	log.CodeLog.Printf("<= function %+v", node)
}

func (g *Gen) genIteration(node syntree.Node) {
	log.CodeLog.Printf("=> iteration %+v %d", node, len(node.Children()))
	for _, n := range node.Children() {
		g.gen(n)
	}
	log.CodeLog.Printf("<= iteration %+v", node)
}

func (g *Gen) genReturn(node syntree.Node) {
	log.CodeLog.Printf("=> return %+v", node)
	g.gen(node)
	log.CodeLog.Printf("<= return %+v", node)
}

func (g *Gen) genSelection(node syntree.Node) {
	log.CodeLog.Printf("=> selection %+v (%d)", node, len(node.Children()))
	for _, n := range node.Children() {
		g.gen(n)
	}
	log.CodeLog.Printf("<= selection %+v (%d)", node, len(node.Children()))
}

func (g *Gen) genAssign(node syntree.Node) {
	log.CodeLog.Printf("=> assign %+v", node)
	for _, n := range node.Children() {
		log.CodeLog.Printf("%+v", n)
		g.gen(n)
	}
	log.CodeLog.Printf("<= assign %+v", node)
}

func (g *Gen) genCall(node syntree.Node) {
	log.CodeLog.Printf("=> call %+v", node)
	for _, n := range node.Children() {
		g.gen(n)
	}
	log.CodeLog.Printf("<= call %+v", node)
}

func (g *Gen) genConst(node syntree.Node) {
	log.CodeLog.Printf("=> const %+v", node)
	comment := fmt.Sprintf("load const with %d", node.Value())
	g.emitRM("LDC", g.ac, node.Value(), 0, comment)
	log.CodeLog.Printf("<= const %+v", node)
}

func (g *Gen) genOp(node syntree.Node) {
	log.CodeLog.Printf("=> op %+v", node)
	for _, n := range node.Children() {
		g.gen(n)
	}
	log.CodeLog.Printf("<= op %+v", node)
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
		log.CodeLog.Printf("=> id_arr %+v", node)
		g.gen(node.Children()[0])
		log.CodeLog.Printf("<= id_arr %+v", node)
	} else {
		log.CodeLog.Printf("=> id %+v", node)
		log.CodeLog.Printf("<= id %+v", node)
	}
}

func (g *Gen) genParam(node syntree.Node) {
	if node.IsArray() {
		log.CodeLog.Printf("=> param_arr %+v", node)
		log.CodeLog.Printf("<= param_arr %+v", node)
	} else {
		log.CodeLog.Printf("=> param %+v", node)
		log.CodeLog.Printf("<= param %+v", node)
	}
}

func (g *Gen) genVar(node syntree.Node) {
	if node.IsArray() {
		log.CodeLog.Printf("=> var_arr %+v", node)
		log.CodeLog.Printf("<= var_arr %+v", node)
	} else {
		log.CodeLog.Printf("=> var %+v", node)
		log.CodeLog.Printf("<= var %+v", node)
	}
}

func (g *Gen) gen(node syntree.Node) {
	if node != nil {
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
		g.gen(node.Sibling())
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
