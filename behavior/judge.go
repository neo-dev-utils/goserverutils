package behavior

// JudgeInter 条件判断实行接口
type JudgeInter interface {
	//条件判断
	ConditionJudge(condition string, param any) bool
	//概率判断
	ProbabilityJudge(min, max uint32, param any) bool
}
