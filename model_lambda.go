package main

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/service/lambda"
)

type ModelLambda struct{}

func NewModelLambda() ModelLambda {
	return ModelLambda{}
}

type ModelLambdaDeployOutput lambda.InvokeOutput

func (self ModelLambdaDeployOutput) Status() DeployStatus {
	if *self.StatusCode != 200 {
		return DeployStatusFail
	}
	return DeployStatusSuccess
}

func (self ModelLambdaDeployOutput) Message() string {
	resByte, err := json.Marshal(self)
	if err != nil {
		return ""
	}
	return string(resByte)
}

func (self ModelLambda) Deploy(pj DeployProject, phase string, option DeployOption) (o DeployOutput, err error) {
	lambda, err := CreateLambdaInstance()
	if err != nil {
		return
	}
	ecr, err := CreateECRInstance()
	if err != nil {
		return
	}
	tag, err := ecr.FindImageTagByRegexp(pj.ECRRepository(), pj.FilterRegexp(), pj.TargetRegexp(), ImageTagVars{Branch: option.Branch, Phase: phase})
	if err != nil {
		return
	}
	ph := pj.FindPhase(phase)
	payload, err := PayloadVars{Tag: tag}.Parse(ph.Payload)
	if err != nil {
		return
	}

	res, err := lambda.Invoke(pj.FuncName(), payload)
	return ModelLambdaDeployOutput(*res), err
}
