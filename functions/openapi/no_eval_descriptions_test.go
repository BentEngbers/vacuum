package openapi

import (
	"github.com/daveshanley/vacuum/model"
	"github.com/daveshanley/vacuum/utils"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"regexp"
	"testing"
)

func TestNoEvalInDescriptions_GetSchema(t *testing.T) {
	def := NoEvalInDescriptions{}
	assert.Equal(t, "no_eval_descriptions", def.GetSchema().Name)
}

func TestNoEvalInDescriptions_RunRule(t *testing.T) {
	def := NoEvalInDescriptions{}
	res := def.RunRule(nil, model.RuleFunctionContext{})
	assert.Len(t, res, 0)
}

func TestNoEvalInDescriptions_RunRule_SuccessCheck(t *testing.T) {

	yml := `paths:
  /pizza/:
    description: do you do crisps?"
  /cake/:
    description: nah mate, only onions.
components:
  schemas:
    CrispsOnion:
      description: a lovely bunch of coconuts`

	path := "$"

	var rootNode yaml.Node
	yaml.Unmarshal([]byte(yml), &rootNode)

	nodes, _ := utils.FindNodes([]byte(yml), path)

	rule := buildOpenApiTestRuleAction(path, "no_eval_description", "", nil)
	fo := make(map[string]string)
	fo["pattern"] = "eval\\("
	comp, _ := regexp.Compile(fo["pattern"])
	rule.PrecomiledPattern = comp
	ctx := buildOpenApiTestContext(model.CastToRuleAction(rule.Then), fo)
	ctx.Rule = &rule
	ctx.Index = model.NewSpecIndex(&rootNode)

	def := NoEvalInDescriptions{}
	res := def.RunRule(nodes, ctx)

	assert.Len(t, res, 0)
}

func TestNoEvalInDescriptions_RunRule_EvalFail(t *testing.T) {

	yml := `paths:
  /pizza/:
    description: eval("alert('hax0r')")"
  /cake/:
    description: nah mate, only onions.
components:
  schemas:
    CrispsOnion:
      description: eval("/*scripkiddy.js*/")`

	path := "$"

	var rootNode yaml.Node
	yaml.Unmarshal([]byte(yml), &rootNode)

	nodes, _ := utils.FindNodes([]byte(yml), path)

	rule := buildOpenApiTestRuleAction(path, "no_eval_description", "", nil)
	fo := make(map[string]string)
	fo["pattern"] = "eval\\("
	comp, _ := regexp.Compile(fo["pattern"])
	rule.PrecomiledPattern = comp

	ctx := buildOpenApiTestContext(model.CastToRuleAction(rule.Then), fo)
	ctx.Index = model.NewSpecIndex(&rootNode)
	ctx.Rule = &rule

	def := NoEvalInDescriptions{}
	res := def.RunRule(nodes, ctx)

	assert.Len(t, res, 2)
}

func TestNoScriptInDescriptions_RunRule_EvalFail(t *testing.T) {

	yml := `paths:
  /pizza/:
    description: <script>console.log('hax0r')</script>"
  /cake/:
    description: nah mate, only onions.
components:
  schemas:
    CrispsOnion:
      description: no hack`

	path := "$"

	var rootNode yaml.Node
	yaml.Unmarshal([]byte(yml), &rootNode)

	nodes, _ := utils.FindNodes([]byte(yml), path)

	rule := buildOpenApiTestRuleAction(path, "no_script_description", "", nil)

	fo := make(map[string]string)
	fo["pattern"] = "<script"
	comp, _ := regexp.Compile(fo["pattern"])
	rule.PrecomiledPattern = comp
	ctx := buildOpenApiTestContext(model.CastToRuleAction(rule.Then), fo)
	ctx.Index = model.NewSpecIndex(&rootNode)
	ctx.Rule = &rule

	def := NoEvalInDescriptions{}
	res := def.RunRule(nodes, ctx)

	assert.Len(t, res, 1)
}
