package main

import (
	"strconv"
)

type OperatorTree struct {
	branches     map[rune]*OperatorTree
	childOpCount int
	// a bitmask stating whether any special tokens exist on this branch
	controlOps    uint16
	operatorToken int
}

func NewOperatorTree() *OperatorTree {
	m := make(map[rune]*OperatorTree)
	o := new(OperatorTree)
	o.branches = m
	o.childOpCount = 0
	o.operatorToken = NIL_TOKEN
	return o
}

func (tree *OperatorTree) ToString(formatrune bool) string {
	result := "[" + strconv.FormatInt(int64(tree.operatorToken), 10) + "]{"
	for i, v := range tree.branches {
		if formatrune {
			result += strconv.FormatInt(int64(i), 10)
		} else {
			result += string(i)
		}
		result += ":" + v.ToString(formatrune) + ","
	}
	return result + "}"
}

func (tree *OperatorTree) AddOperator(ra []rune, token int) bool {
	if len(ra) == 0 {
		if tree.operatorToken == NIL_TOKEN {
			tree.operatorToken = token
			if token < NIL_TOKEN {
				tree.controlOps |= 1 << ^token
			}
			return true
		} else {
			return false
		}
	}
	c := ra[0]
	branch, ok := tree.branches[c]
	if !ok {
		branch = NewOperatorTree()
		tree.branches[c] = branch
	}
	success := branch.AddOperator(ra[1:], token)
	if success {
		tree.childOpCount++
		if token < NIL_TOKEN {
			tree.controlOps |= 1 << ^token
		}
	}
	return success
}

func (tree *OperatorTree) AddOperator_rune(r rune, token int) (worked bool) {
	branch, ok := tree.branches[r]
	if ok {
		if branch.operatorToken == NIL_TOKEN {
			branch.operatorToken = token
			tree.childOpCount++
			if token < NIL_TOKEN {
				tree.controlOps |= 1 << ^token
				branch.controlOps |= 1 << ^token
			}
			return true
		} else {
			return false
		}
	} else {
		branch = NewOperatorTree()
		branch.operatorToken = token
		tree.branches[r] = branch
		tree.childOpCount++
		if token < NIL_TOKEN {
			tree.controlOps |= 1 << ^token
			branch.controlOps |= 1 << ^token
		}
		return true
	}
}

func (tree *OperatorTree) PossibleCount(ra []rune) (int, *OperatorTree) {
	count, subtree := tree.PossibleChildCount(ra)
	if subtree.operatorToken != NIL_TOKEN && tree != subtree {
		count++
	}
	return count, subtree
}

func (tree *OperatorTree) PossibleCount_rune(r rune) (int, *OperatorTree) {
	count, subtree := tree.PossibleChildCount_rune(r)
	if subtree.operatorToken != NIL_TOKEN && tree != subtree {
		count++
	}
	return count, subtree
}

func (tree *OperatorTree) PossibleChildCount(ra []rune) (int, *OperatorTree) {
	if len(ra) == 0 {
		return tree.childOpCount, tree
	}
	c := ra[0]
	branch, ok := tree.branches[c]
	if ok {
		return branch.PossibleChildCount(ra[1:])
	} else {
		return 0, tree
	}
}

func (tree *OperatorTree) PossibleChildCount_rune(r rune) (int, *OperatorTree) {
	branch, ok := tree.branches[r]
	if ok {
		return branch.childOpCount, branch
	} else {
		return 0, tree
	}
}

func (tree *OperatorTree) OperatorExists(ra []rune) bool {
	return tree.GetToken(ra) != NIL_TOKEN
}

func (tree *OperatorTree) GetToken(ra []rune) int {
	branch := tree.GetBranch(ra)
	if branch == nil {
		return NIL_TOKEN
	} else {
		return branch.operatorToken
	}
}

func (tree *OperatorTree) GetBranch(ra []rune) *OperatorTree {
	if len(ra) == 0 {
		return tree
	}
	c := ra[0]
	branch, ok := tree.branches[c]
	if ok {
		return branch.GetBranch(ra[1:])
	}
	return nil
}
