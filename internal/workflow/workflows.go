package workflow

type WorkFlow interface {
	Name() string
	Description() string
	Run(*workflowParams)
}

var WorkFlows = make(map[string]WorkFlow)

func RegisterWorkflow(workflow WorkFlow) {
	if _, ok := WorkFlows[workflow.Name()]; ok {
		return
	}
	WorkFlows[workflow.Name()] = workflow
}

func GetFlowNames() []string {
	var names []string
	for name, _ := range WorkFlows {
		names = append(names, name)
	}
	return names
}
