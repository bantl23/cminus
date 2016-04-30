package opt

import (
	"github.com/bantl23/cminus/log"
	"github.com/bantl23/cminus/syntree"
)

var CONST_FOLDED bool = false

func ConstantFolding(node syntree.Node) {
	if node != nil {
		for _, n := range node.Children() {
			ConstantFolding(n)
			if n != nil && n.IsOp() {
				if n.Children() != nil {
					parent := n.Parent()
					child := -1
					for i, p := range parent.Children() {
						if n == p {
							child = i
						}
					}
					left := n.Children()[0]
					right := n.Children()[1]
					if left.IsConst() && right.IsConst() {
						if n.TokType() == syntree.PLUS {
							value := left.Value() + right.Value()
							newNode := syntree.NewExpConstNode(n.Pos().Row(), n.Pos().Col(), value)
							parent.Children()[child] = newNode
							CONST_FOLDED = true
							log.OptLog.Printf("new node %+v", newNode)
						} else if n.TokType() == syntree.MINUS {
							value := left.Value() - right.Value()
							newNode := syntree.NewExpConstNode(n.Pos().Row(), n.Pos().Col(), value)
							parent.Children()[child] = newNode
							CONST_FOLDED = true
							log.OptLog.Printf("new node %+v", newNode)
						} else if n.TokType() == syntree.TIMES {
							value := left.Value() * right.Value()
							newNode := syntree.NewExpConstNode(n.Pos().Row(), n.Pos().Col(), value)
							parent.Children()[child] = newNode
							CONST_FOLDED = true
							log.OptLog.Printf("new node %+v", newNode)
						} else if n.TokType() == syntree.OVER {
							value := left.Value() / right.Value()
							newNode := syntree.NewExpConstNode(n.Pos().Row(), n.Pos().Col(), value)
							parent.Children()[child] = newNode
							CONST_FOLDED = true
							log.OptLog.Printf("new node %+v", newNode)
						} else if n.TokType() == syntree.EQ {
							value := 0
							if left.Value() == right.Value() {
								value = 1
							}
							newNode := syntree.NewExpConstNode(n.Pos().Row(), n.Pos().Col(), value)
							parent.Children()[child] = newNode
							CONST_FOLDED = true
							log.OptLog.Printf("new node %+v", newNode)
						} else if n.TokType() == syntree.NEQ {
							value := 0
							if left.Value() != right.Value() {
								value = 1
							}
							newNode := syntree.NewExpConstNode(n.Pos().Row(), n.Pos().Col(), value)
							parent.Children()[child] = newNode
							CONST_FOLDED = true
							log.OptLog.Printf("new node %+v", newNode)
						} else if n.TokType() == syntree.LT {
							value := 0
							if left.Value() < right.Value() {
								value = 1
							}
							newNode := syntree.NewExpConstNode(n.Pos().Row(), n.Pos().Col(), value)
							parent.Children()[child] = newNode
							CONST_FOLDED = true
							log.OptLog.Printf("new node %+v", newNode)
						} else if n.TokType() == syntree.LTE {
							value := 0
							if left.Value() <= right.Value() {
								value = 1
							}
							newNode := syntree.NewExpConstNode(n.Pos().Row(), n.Pos().Col(), value)
							parent.Children()[child] = newNode
							CONST_FOLDED = true
							log.OptLog.Printf("new node %+v", newNode)
						} else if n.TokType() == syntree.GT {
							value := 0
							if left.Value() > right.Value() {
								value = 1
							}
							newNode := syntree.NewExpConstNode(n.Pos().Row(), n.Pos().Col(), value)
							parent.Children()[child] = newNode
							CONST_FOLDED = true
							log.OptLog.Printf("new node %+v", newNode)
						} else if n.TokType() == syntree.GTE {
							value := 0
							if left.Value() < right.Value() {
								value = 1
							}
							newNode := syntree.NewExpConstNode(n.Pos().Row(), n.Pos().Col(), value)
							parent.Children()[child] = newNode
							CONST_FOLDED = true
							log.OptLog.Printf("new node %+v", newNode)
						}
					}
				}
			}
		}
		ConstantFolding(node.Sibling())
	}
}

var vals map[string]int = make(map[string]int)
var containsLeftAssign bool = false
var CONST_PROPAGATED bool = false

func ConstantPropagation(node syntree.Node) {
	if node != nil {
		if node.IsIteration() || node.IsSelection() {
			vals = make(map[string]int)
		}
		for _, n := range node.Children() {
			ConstantPropagation(n)
		}
		if node.IsAssign() {
			if node.Children()[0].IsArray() == false {
				id := node.Children()[0].Name()
				if node.Children()[1].IsConst() {
					value := node.Children()[1].Value()
					vals[id] = value
					containsLeftAssign = false
				}
				if containsLeftAssign == true {
					delete(vals, id)
				}
			}
		}
		if node.IsId() {
			id := node.Name()
			parent := node.Parent()
			if parent != nil {
				if parent.IsCall() && parent.Name() != "output" {
					sib := node
					var prevSib syntree.Node = nil
					for sib != nil {
						value, ok := vals[sib.Name()]
						if ok == true {
							if sib.Parent() != nil {
								newNode := syntree.NewExpConstNode(node.Pos().Row(), node.Pos().Col(), value)
								newNode.SetSibling(sib.Sibling())
								sib.Parent().Children()[0] = newNode
								sib = newNode
							} else {
								newNode := syntree.NewExpConstNode(node.Pos().Row(), node.Pos().Col(), value)
								newNode.SetSibling(sib.Sibling())
								prevSib.SetSibling(newNode)
								sib = newNode
							}
						}
						prevSib = sib
						sib = sib.Sibling()
					}
				} else {
					children := []int{}
					for i, c := range parent.Children() {
						if c == node {
							children = append(children, i)
						}
					}

					for _, idx := range children {
						if node.Parent().IsAssign() && idx == 0 {
							containsLeftAssign = true
						} else {
							value, ok := vals[id]
							if ok == true {
								newNode := syntree.NewExpConstNode(node.Pos().Row(), node.Pos().Col(), value)
								node.Parent().Children()[idx] = newNode
								log.OptLog.Printf("new node %+v", newNode)
								CONST_PROPAGATED = true
							}
						}
					}
				}
			}
		}
		if node.IsCompound() {
			vals = make(map[string]int)
		}
		ConstantPropagation(node.Sibling())
	}
}

var funcMap map[string]bool = make(map[string]bool)

func FindDeadFuncs(node syntree.Node) {
	if node != nil {
		if node.IsFunc() {
			funcMap[node.Name()] = false
			if node.Name() == "main" {
				funcMap[node.Name()] = true
			}
		}
		for _, n := range node.Children() {
			FindDeadFuncs(n)
		}
		if node.IsCall() {
			funcMap[node.Name()] = true
		}
		FindDeadFuncs(node.Sibling())
	}
}

func RemoveDeadFuncs(node syntree.Node) bool {
	var removeRoot bool = false
	var prevNode syntree.Node = nil
	for node != nil {
		if node.IsFunc() && funcMap[node.Name()] == false {
			log.OptLog.Printf("removing dead function %+v", node)
			if prevNode != nil {
				prevNode.SetSibling(node.Sibling())
			} else {
				removeRoot = true
			}
		}
		prevNode = node
		node = node.Sibling()
	}
	return removeRoot
}

var firstVar syntree.Node = nil
var varMap map[string]bool = make(map[string]bool)

func RemoveDeadVars(node syntree.Node) {
	if node != nil {
		if node.IsCompound() {
			if node.Children() != nil {
				firstVar = node.Children()[0]
				sib := node.Children()[0]
				for sib != nil {
					varMap[sib.Name()] = false
					sib = sib.Sibling()
				}
			}
		}
		if node.IsIteration() || node.IsSelection() {
			varMap = make(map[string]bool)
		}
		for _, n := range node.Children() {
			RemoveDeadVars(n)
		}
		if node.IsId() {
			if node.Parent().IsAssign() && node.Parent().Children()[0] == node {
			} else {
				_, ok := varMap[node.Name()]
				if ok == true {
					varMap[node.Name()] = true
				}
			}
		}
		if node.IsCompound() {
			sib := firstVar
			var prevSib syntree.Node = nil
			for sib != nil {
				if varMap[sib.Name()] == false {
					if sib.Parent() != nil {
						parent := sib.Parent()
						next := sib.Sibling()
						parent.Children()[0] = next
						if next != nil {
							next.SetParent(parent)
						}
					} else {
						prevSib.SetSibling(sib.Sibling())
					}
				}
				prevSib = sib
				sib = sib.Sibling()
			}
			if len(node.Children()) >= 1 {
				var prevSib syntree.Node = nil
				sib := node.Children()[1]
				for sib != nil {
					log.OptLog.Printf("=== %+v %+v %+v %+v", sib, sib.Parent(), prevSib, sib.Sibling())
					if varMap[sib.Children()[0].Name()] == false {
						if sib.IsAssign() {
							if sib.Parent() != nil {
								parent := sib.Parent()
								next := sib.Sibling()
								log.OptLog.Printf("+++ %+v %+v %+v %+v %+v", sib, sib.Parent(), prevSib, sib.Sibling(), parent.Children()[1])
								parent.Children()[1] = next
								if next != nil {
									next.SetParent(parent)
								}
								log.OptLog.Printf("+++ %+v %+v %+v %+v %+v", sib, sib.Parent(), prevSib, sib.Sibling(), parent.Children()[1])
							} else {
								log.OptLog.Printf("||| %+v %+v %+v %+v", sib, sib.Parent(), prevSib, sib.Sibling())
								prevSib.SetSibling(sib.Sibling())
							}
						}
					}
					log.OptLog.Printf("=== %+v %+v %+v %+v", sib, sib.Parent(), prevSib, sib.Sibling())
					prevSib = sib
					sib = sib.Sibling()
				}
			}
			varMap = make(map[string]bool)
		}
		RemoveDeadVars(node.Sibling())
	}
}
