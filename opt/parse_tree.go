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

var vals map[string]int = nil
var containsLeftAssign bool = false
var CONST_PROPAGATED bool = false

func ConstantPropagation(node syntree.Node) {
	if node != nil {
		log.OptLog.Printf("PRENODE %+v", node)
		if node.IsCompound() {
			vals = make(map[string]int)
		}
		for _, n := range node.Children() {
			ConstantPropagation(n)
		}
		if node.IsAssign() {
			id := node.Children()[0].Name()
			if node.Children()[1].IsConst() {
				value := node.Children()[1].Value()
				vals[id] = value
				containsLeftAssign = false
			}
			if containsLeftAssign == true {
				delete(vals, id)
				log.OptLog.Printf("+++++++++ DELETE %+v", vals)
			}
			log.OptLog.Printf("+++++++++++ ISASSIGN vals %+v", vals)
		}
		if node.IsId() {
			log.OptLog.Printf("++++++++++ ISID %+v", node)
			id := node.Name()
			parent := node.Parent()
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
						log.OptLog.Printf("++++++++ UPDATED %+v", newNode)
						CONST_PROPAGATED = true
					}
				}
			}
		}
		if node.IsCompound() {
			vals = nil
		}
		log.OptLog.Printf("POSTNODE %+v", node)
		ConstantPropagation(node.Sibling())
	}

}

var funcNode syntree.Node = nil

func TailRecursion(node syntree.Node) {
	if node != nil {
		if node.IsFunc() {
			funcNode = node
		}
		for _, n := range node.Children() {
			TailRecursion(n)
		}
		if node.IsReturn() {
			if funcNode != nil && node.Children() != nil && len(node.Children()) == 1 && node.Children()[0].IsCall() {
				callName := node.Children()[0].Name()
				if callName == funcNode.Name() {
					funcNode.SetTail(true)
					log.OptLog.Printf("Tail recursion found %+v", funcNode)
				}
			}
		}
		if node.IsFunc() {
			funcNode = nil
		}
		node = node.Sibling()
	}
}
