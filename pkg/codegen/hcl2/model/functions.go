// Copyright 2016-2020, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/pulumi/pulumi/pkg/codegen/hcl2/syntax"
)

type FunctionSignature interface {
	GetSignature(arguments []Expression) (StaticFunctionSignature, hcl.Diagnostics)
}

// Parameter represents a single function parameter.
type Parameter struct {
	Name string // The name of the parameter.
	Type Type   // The type of the parameter.
}

// StaticFunctionSignature represents the parameters and return type of a function.
type StaticFunctionSignature struct {
	Parameters       []Parameter
	VarargsParameter *Parameter
	ReturnType       Type
}

func (fs StaticFunctionSignature) GetSignature(arguments []Expression) (StaticFunctionSignature, hcl.Diagnostics) {
	return fs, nil
}

type GenericFunctionSignature func(arguments []Expression) (StaticFunctionSignature, hcl.Diagnostics)

func (fs GenericFunctionSignature) GetSignature(arguments []Expression) (StaticFunctionSignature, hcl.Diagnostics) {
	return fs(arguments)
}

// Function represents a function definition.
type Function struct {
	signature FunctionSignature
}

func NewFunction(signature FunctionSignature) *Function {
	return &Function{signature: signature}
}

func (f *Function) SyntaxNode() hclsyntax.Node {
	return syntax.None
}

func (f *Function) Traverse(traverser hcl.Traverser) (Traversable, hcl.Diagnostics) {
	return AnyType, hcl.Diagnostics{cannotTraverseFunction(traverser.SourceRange())}
}

func (f *Function) GetSignature(arguments []Expression) (StaticFunctionSignature, hcl.Diagnostics) {
	return f.signature.GetSignature(arguments)
}

var pulumiBuiltins = map[string]*Function{
	"fileAsset": NewFunction(StaticFunctionSignature{
		Parameters: []Parameter{{
			Name: "path",
			Type: StringType,
		}},
		ReturnType: AssetType,
	}),
	"mimeType": NewFunction(StaticFunctionSignature{
		Parameters: []Parameter{{
			Name: "path",
			Type: StringType,
		}},
		ReturnType: AssetType,
	}),
	"toJSON": NewFunction(StaticFunctionSignature{
		Parameters: []Parameter{{
			Name: "value",
			Type: AnyType,
		}},
		ReturnType: StringType,
	}),
}