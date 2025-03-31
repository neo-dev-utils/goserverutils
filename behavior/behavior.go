package behavior

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"unicode"
	"unicode/utf8"
)

// Precompute the reflect type for error. Can't use error directly
// because Typeof takes an empty interface value. This is annoying.
var typeOfError = reflect.TypeOf((*error)(nil)).Elem()

// Precompute the reflect type for context.
var typeOfContext = reflect.TypeOf((*context.Context)(nil)).Elem()

// 子服务处理方法
type methodstruct struct {
	name       string
	method     reflect.Method
	classValue reflect.Value
	ctxType    reflect.Type
	argType    reflect.Type
	replyType  reflect.Type
	returnType reflect.Type
}

// Behavior 行为树
type Behavior struct {
	//Node 入口结点
	Node []*Node
	//Param 参数
	Param any
	//方法
	methods map[string]*methodstruct
}

// LoadNodeByJSON 通过json字符串加载数据
func LoadNodeByJSON(action any, judge JudgeInter, buff []byte) (*Behavior, error) {
	body := JSONBody{}
	err := json.Unmarshal(buff, &body)
	if err != nil {
		return nil, err
	}
	pointer := &Behavior{
		methods: make(map[string]*methodstruct),
	}
	typ := reflect.TypeOf(action)
	rcvr := reflect.ValueOf(action)
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		mtype := method.Type
		mname := method.Name
		// Method must be exported.
		if method.PkgPath != "" {
			continue
		}
		if mtype.NumIn() != 2 {
			fmt.Println("in", mtype.NumIn())
			continue
		}
		if mtype.NumOut() > 0 {
			fmt.Println("out", mtype.NumOut())
			continue
		}

		argType := mtype.In(1)
		if !pointer.isExportedOrBuiltinType(argType) {
			continue
		}
		pointer.methods[mname] = &methodstruct{name: mname, method: method, classValue: rcvr, argType: argType}
	}

	//遍历加载
	pointer._RangeJSONNode(action, judge, nil, body.Behavior)
	return pointer, nil
}

// _RangeJSONNode 遍历json node
func (pointer *Behavior) _RangeJSONNode(action ActionInter, judge JudgeInter, parentNode *Node, nodes []JSONNode) {
	//遍历加载
	for _, node := range nodes {
		switch node.Kind {
		case JudgeCondition:
			newNode := CreateJudgeNode(node.IsAction, node.JudeCondition, judge, pointer.methods[node.Action])
			pointer._RangeJSONNode(action, judge, newNode, node.Children)
			if parentNode == nil {
				pointer.Node = append(pointer.Node, newNode)
				continue
			}
			parentNode.AddNode(newNode)
		case Probability:
			//概率结点小于就进入
			newNode := CreateProbabilityNode(node.IsAction, node.ProbabilityStartCondition, node.ProbabilityEndCondition, judge, pointer.methods[node.Action])
			pointer._RangeJSONNode(action, judge, newNode, node.Children)
			if parentNode == nil {
				pointer.Node = append(pointer.Node, newNode)
				continue
			}
			parentNode.AddNode(newNode)
		}
	}
}

// DoAction 执行动作
func (pointer *Behavior) DoAction(param any) {
	//参数
	pointer.Param = param
	//遍历当下所有结点
	pointer._RangeNode(pointer._CallBack)
}

// CallBack  遍历回掉
func (pointer *Behavior) _CallBack(node *Node) bool {
	//允许进入
	if node.AllowAccess(pointer.Param) {
		//判断是否是动作结点
		if node.IsActionNode() {
			//是动作结点执行动作
			node.DoAction(pointer.Param)
			//结束循环
			return false
		}
		//不是动作结点继续遍历
		node.RangeNextNode(pointer._CallBack)
	}
	return true
}

// RangeNode 遍历结点
func (pointer *Behavior) _RangeNode(callback func(node *Node) bool) {
	for _, n := range pointer.Node {
		if callback(n) == false {
			break
		}
	}
}

func (pointer *Behavior) isExported(name string) bool {
	rune, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(rune)
}

func (pointer *Behavior) isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return pointer.isExported(t.Name()) || t.PkgPath() == ""
}
