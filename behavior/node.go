package behavior

import (
	"reflect"
)

// NodeType 结点类型
type NodeType uint8

const (
	//JudgeCondition 条件结点
	JudgeCondition NodeType = iota
	//Probability 概率结点【设置进入区间100比】
	Probability
)

// Node 树的结点
type Node struct {
	//Kind 结点类型
	_Kind NodeType
	//IsAction 是否是动作结点
	_IsAction bool
	//JudgeCondition 判断条件
	_JudgeCondition string
	//Judge 判断对象接口实现
	_Judge JudgeInter
	//ProbabilityStartCondition 概率条件--开始区间
	_ProbabilityStartCondition uint32
	//ProbabilityEndCondition 概率条件--结束区间
	_ProbabilityEndCondition uint32
	//ParentNode 父结点
	_ParentNode *Node
	//NextNode 下一个结点
	_NextNode []*Node
	//ActionFunc 执行的动作
	_ActionFunc *methodstruct
}

// CreateJudgeNode 创建条件结点
func CreateJudgeNode(isAction bool, judgeCondition string, judgeInter JudgeInter, actionFunc *methodstruct) *Node {
	node := &Node{
		_Kind:           JudgeCondition,
		_IsAction:       isAction,
		_JudgeCondition: judgeCondition,
		_Judge:          judgeInter,
		_ActionFunc:     actionFunc,
	}
	return node
}

// CreateProbabilityNode 创建概率结点
func CreateProbabilityNode(isAction bool, probabilityStartCondition, probabilityEndCondition uint32, judgeInter JudgeInter, actionFunc *methodstruct) *Node {
	start := probabilityStartCondition
	if probabilityStartCondition > probabilityEndCondition {
		start = probabilityEndCondition
		probabilityEndCondition = probabilityStartCondition
	}
	node := &Node{
		_Kind:                      Probability,
		_IsAction:                  isAction,
		_ProbabilityStartCondition: start,
		_ProbabilityEndCondition:   probabilityEndCondition,
		_Judge:                     judgeInter,
		_ActionFunc:                actionFunc,
	}
	return node
}

// AddNode 添加儿子结点
func (pointer *Node) AddNode(node *Node) {
	pointer._NextNode = append(pointer._NextNode, node)
}

// RangeNextNode 遍历儿子结点
func (pointer *Node) RangeNextNode(callback func(node *Node) bool) {
	for _, n := range pointer._NextNode {
		if callback(n) == false {
			break
		}
	}
}

// GetParentNode 获取父亲结点
func (pointer *Node) GetParentNode() *Node {
	return pointer._ParentNode
}

// AllowAccess 允许进入
func (pointer *Node) AllowAccess(param any) bool {
	switch pointer._Kind {
	case JudgeCondition:
		//条件结点
		if pointer._Judge != nil {
			return pointer._Judge.ConditionJudge(pointer._JudgeCondition, param)
		}
		return false
	case Probability:
		//概率结点小于就进入
		if pointer._Judge != nil {
			return pointer._Judge.ProbabilityJudge(pointer._ProbabilityStartCondition, pointer._ProbabilityEndCondition, param)
		}
		return false
	}
	return false
}

// IsActionNode 是否为动作结点
func (pointer *Node) IsActionNode() bool {
	if pointer._IsAction == true {
		return true
	}
	return false
}

// DoAction 执行动作
func (pointer *Node) DoAction(param any) {
	if pointer._ActionFunc != nil {
		method := pointer._ActionFunc.method.Func
		if method.IsValid() {
			if pointer._ActionFunc.argType.Kind() != reflect.Ptr {
				method.Call([]reflect.Value{pointer._ActionFunc.classValue, reflect.ValueOf(param).Elem()})
			} else {
				method.Call([]reflect.Value{pointer._ActionFunc.classValue, reflect.ValueOf(param)})
			}
		}

	}
}
