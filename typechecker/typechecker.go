package typechecker

import (
	"fmt"
	"jru-test/ast"
	"jru-test/lexer"
)

type Type interface {
	String() string
	Equals(other Type) bool
}

type PrimitiveType struct {
	Name string
}

func (p PrimitiveType) String() string {
	return p.Name
}

func (p PrimitiveType) Equals(other Type) bool {
	if o, ok := other.(PrimitiveType); ok {
		return p.Name == o.Name
	}
	return false
}

func IsPrimitive(t Type, name string) bool {
	if p, ok := t.(PrimitiveType); ok {
		return p.Name == name
	}
	return false
}

func IsNumeric(t Type) bool {
	if p, ok := t.(PrimitiveType); ok {
		return p.Name == "i8" || p.Name == "i32" || p.Name == "i64" || p.Name == "f32" || p.Name == "f64"
	}
	return false
}

type ArrayType struct {
	ElemType Type
}

func (a ArrayType) String() string {
	return fmt.Sprintf("%s[]", a.ElemType)
}

func (a ArrayType) Equals(other Type) bool {
	if o, ok := other.(ArrayType); ok {
		return a.ElemType.Equals(o.ElemType)
	}
	return false
}

type FuncType struct {
	ReturnType Type
	ParamTypes []Type
}

func (f FuncType) String() string {
	params := ""
	for i, param := range f.ParamTypes {
		if i > 0 {
			params += ","
		}
		params += param.String()
	}
	return fmt.Sprintf("func(%s):%s", params, f.ReturnType)
}

func (f FuncType) Equals(other Type) bool {
	o, ok := other.(FuncType)
	if !ok || len(f.ParamTypes) != len(o.ParamTypes) {
		return false
	}
	if !f.ReturnType.Equals(o.ReturnType) {
		return false
	}
	for i, param := range f.ParamTypes {
		if !param.Equals(o.ParamTypes[i]) {
			return false
		}
	}
	return true
}

type StructType struct {
	Name    string
	Members map[string]Type
}

func (s StructType) String() string {
	return s.Name
}

func (s StructType) Equals(other Type) bool {
	if o, ok := other.(StructType); ok {
		return s.Name == o.Name
	}
	return false
}

type TypeEnv struct {
	parent      *TypeEnv
	vars        map[string]Type
	structTypes map[string]StructType
	funcs       map[string]string
	funcTypes   map[string]FuncType
}

func NewTypeEnv(parent *TypeEnv) *TypeEnv {
	return &TypeEnv{
		parent:      parent,
		vars:        make(map[string]Type),
		structTypes: make(map[string]StructType),
		funcs:       make(map[string]string),
		funcTypes:   make(map[string]FuncType),
	}
}

func (env *TypeEnv) DefineVar(name string, varType Type) {
	env.vars[name] = varType
}

func (env *TypeEnv) LookupVarType(name string) (Type, bool) {
	if varType, ok := env.vars[name]; ok {
		return varType, true
	}
	if env.parent != nil {
		return env.parent.LookupVarType(name)
	}
	return nil, false
}

func (env *TypeEnv) DefineStructType(name string, st StructType) {
	env.structTypes[name] = st
}

func (env *TypeEnv) LookupStructType(name string) (StructType, bool) {
	if st, ok := env.structTypes[name]; ok {
		return st, true
	}
	if env.parent != nil {
		return env.parent.LookupStructType(name)
	}
	return StructType{}, false
}

func (env *TypeEnv) DefineFunc(name string, funcTypeName string) {
	env.funcs[name] = funcTypeName
}

func (env *TypeEnv) LookupFunc(name string) (string, bool) {
	if fn, ok := env.funcs[name]; ok {
		return fn, true
	}
	if env.parent != nil {
		return env.parent.LookupFunc(name)
	}
	return "", false
}

func (env *TypeEnv) DefineFuncType(name string, fn FuncType) {
	env.funcTypes[name] = fn
}

func (env *TypeEnv) LookupFuncType(name string) (FuncType, bool) {
	if fn, ok := env.funcTypes[name]; ok {
		return fn, true
	}
	if env.parent != nil {
		return env.parent.LookupFuncType(name)
	}
	return FuncType{}, false
}

type TypeChecker struct {
	Errors     []string
	env        *TypeEnv
	primitives map[string]Type
}

func NewTypeChecker() *TypeChecker {
	return &TypeChecker{
		Errors: []string{},
		env:    NewTypeEnv(nil),
		primitives: map[string]Type{
			"void":   PrimitiveType{Name: "void"},
			"bool":   PrimitiveType{Name: "bool"},
			"string": PrimitiveType{Name: "string"},
			"i8":     PrimitiveType{Name: "i8"},
			"i32":    PrimitiveType{Name: "i32"},
			"i64":    PrimitiveType{Name: "i64"},
			"f32":    PrimitiveType{Name: "f32"},
			"f64":    PrimitiveType{Name: "f64"},
		},
	}
}

func (tc *TypeChecker) Err(msg string) {
	tc.Errors = append(tc.Errors, msg)
}

func (tc *TypeChecker) ResolveType(astType ast.Type) Type {
	switch t := astType.(type) {
	case ast.NamedType:
		if prim, ok := tc.primitives[t.TypeName]; ok {
			return prim
		}
		if structType, ok := tc.env.LookupStructType(t.TypeName); ok {
			return structType
		}
		tc.Err(fmt.Sprintf("undefined type: %s", t.TypeName))
		return nil
	case ast.ArrayType:
		elemType := tc.ResolveType(t.UnderlyingType)
		if elemType == nil {
			return nil
		}
		return ArrayType{ElemType: elemType}
	case ast.FuncType:
		paramTypes := []Type{}
		for _, astParamType := range t.ParamTypes {
			paramType := tc.ResolveType(astParamType)
			if paramType == nil {
				continue
			}
			paramTypes = append(paramTypes, paramType)
		}
		returnType := tc.ResolveType(t.ReturnType)
		if returnType == nil {
			return nil
		}
		return FuncType{
			ReturnType: returnType,
			ParamTypes: paramTypes,
		}
	default:
		tc.Err(fmt.Sprintf("unknown type: %T", astType))
		return nil
	}
}

func Check(program ast.BlockStmt) []string {
	tc := NewTypeChecker()
	tc.CheckBlockStmt(program)
	return tc.Errors
}

func (tc *TypeChecker) CheckBlockStmt(block ast.BlockStmt) {
	oldEnv := tc.env
	tc.env = NewTypeEnv(oldEnv)
	for _, stmt := range block.Body {
		tc.CheckStmt(stmt)
	}
	tc.env = oldEnv
}

func (tc *TypeChecker) CheckStmt(stmt ast.Stmt) {
	switch s := stmt.(type) {
	case ast.BlockStmt:
		tc.CheckBlockStmt(s)
	case ast.VarDeclStmt:
		tc.CheckVarDeclStmt(s)
	case ast.StructDeclStmt:
		tc.CheckStructDeclStmt(s)
	case ast.FuncDeclStmt:
		tc.CheckFuncDeclStmt(s)
	case ast.IfStmt:
		tc.CheckIfStmt(s)
	case ast.ForStmt:
		tc.CheckForStmt(s)
	case ast.ExpressionStmt:
		tc.InferType(s.Expr)
	default:
		tc.Err(fmt.Sprintf("unknown statement type: %T", stmt))
	}
}

func (tc *TypeChecker) CheckVarDeclStmt(stmt ast.VarDeclStmt) {
	declaredType := tc.ResolveType(stmt.Var.Type)
	if declaredType == nil {
		return
	}
	if stmt.InitVal != nil {
		initType := tc.InferType(stmt.InitVal)
		if initType == nil {
			return
		}
		if !declaredType.Equals(initType) {
			tc.Err(fmt.Sprintf("type mismatch: variable %s declared as %s but initialized with %s", stmt.Var.Name, declaredType, initType))
		}
	}
	tc.env.DefineVar(stmt.Var.Name, declaredType)
}

func (tc *TypeChecker) CheckStructDeclStmt(stmt ast.StructDeclStmt) {
	if _, ok := tc.env.LookupStructType(stmt.Name); ok {
		tc.Err(fmt.Sprintf("redeclared struct %s in the same scope", stmt.Name))
		return
	}
	members := make(map[string]Type)
	for _, member := range stmt.Members {
		if _, ok := members[member.Name]; ok {
			tc.Err(fmt.Sprintf("duplicate member %s in struct %s", member.Name, stmt.Name))
			continue
		}
		members[member.Name] = tc.ResolveType(member.Type)
	}
	tc.env.DefineStructType(stmt.Name, StructType{
		Name:    stmt.Name,
		Members: members,
	})
}

func (tc *TypeChecker) CheckFuncDeclStmt(stmt ast.FuncDeclStmt) {
	if _, ok := tc.env.LookupFuncType(stmt.Name); ok {
		tc.Err(fmt.Sprintf("redeclared function %s in the same scope", stmt.Name))
		return
	}
	returnType := tc.primitives["void"]
	if stmt.ReturnType != nil {
		returnType = tc.ResolveType(stmt.ReturnType)
		if returnType == nil {
			return
		}
	}
	paramTypes := make([]Type, 0, len(stmt.Parameters))
	paramEnv := NewTypeEnv(tc.env)
	for _, param := range stmt.Parameters {
		paramType := tc.ResolveType(param.Type)
		if paramType == nil {
			return
		}
		paramTypes = append(paramTypes, paramType)
		paramEnv.DefineVar(param.Name, paramType)
	}
	funcType := FuncType{
		ReturnType: returnType,
		ParamTypes: paramTypes,
	}
	funcTypeName := fmt.Sprintf("%s", funcType)
	tc.env.DefineFunc(stmt.Name, funcTypeName)
	tc.env.DefineFuncType(funcTypeName, funcType)
	oldEnv := tc.env
	tc.env = paramEnv
	tc.CheckBlockStmt(stmt.Body)
	tc.env = oldEnv
}

func (tc *TypeChecker) CheckIfStmt(stmt ast.IfStmt) {
	condType := tc.InferType(stmt.Cond)
	if !IsPrimitive(condType, "bool") {
		tc.Err("if- statement condition does not evaluate to a boolean type")
	}
	tc.CheckStmt(stmt.Then)
	if stmt.Else != nil {
		tc.CheckStmt(stmt.Else)
	}
}

func (tc *TypeChecker) CheckForStmt(stmt ast.ForStmt) {
	tc.CheckStmt(stmt.Init)
	condType := tc.InferType(stmt.Cond)
	if !IsPrimitive(condType, "bool") {
		tc.Err("for- statement condition does not evaluate to a boolean type")
	}
	tc.CheckStmt(stmt.Iter)
	tc.CheckBlockStmt(stmt.Body)
}

func (tc *TypeChecker) InferType(expr ast.Expr) Type {
	switch e := expr.(type) {
	case ast.NumberLiteralExpr:
		return tc.primitives["i32"] // todo; evaluate the number literal to determine exact type
	case ast.StringLiteralExpr:
		return tc.primitives["string"]
	case ast.BoolLiteralExpr:
		return tc.primitives["bool"]
	case ast.IdentExpr:
		if varType, ok := tc.env.LookupVarType(e.Value); ok {
			return varType
		}
		if structType, ok := tc.env.LookupStructType(e.Value); ok {
			return structType
		}
		if funcTypeName, ok := tc.env.LookupFunc(e.Value); ok {
			if funcType, ok := tc.env.LookupFuncType(funcTypeName); ok {
				return funcType
			}
		}
		tc.Err(fmt.Sprintf("undefined variable: %s", e.Value))
		return nil
	case ast.BinaryExpr:
		return tc.CheckBinaryExpr(e)
	case ast.UnaryExpr:
		return tc.CheckUnaryExpr(e)
	case ast.GroupExpr:
		return tc.InferType(e.Expr)
	case ast.FuncCallExpr:
		return tc.CheckFuncCallExpr(e)
	case ast.StructLiteralExpr:
		return tc.CheckStructLiteralExpr(e)
	case ast.StructMemberExpr:
		return tc.CheckStructMemberExpr(e)
	case ast.ArrayIndexExpr:
		return tc.CheckArrayIndexExpr(e)
	case ast.AssignExpr:
		return tc.CheckAssignExpr(e)
	default:
		tc.Err(fmt.Sprintf("unknown expression type: %T", expr))
		return nil
	}
}

func (tc *TypeChecker) CheckBinaryExpr(expr ast.BinaryExpr) Type {
	leftType := tc.InferType(expr.Lhs)
	rightType := tc.InferType(expr.Rhs)
	if leftType == nil || rightType == nil {
		return nil
	}
	switch expr.Operator.Type {
	case lexer.PLUS, lexer.DASH, lexer.STAR, lexer.SLASH, lexer.PERCENT:
		if IsNumeric(leftType) && IsNumeric(rightType) {
			return leftType // no specific reason, just pick one arbitrarily until we have e.g. type promotion (i32 -> f32 etc.)
		}
		if expr.Operator.Type == lexer.PLUS && IsPrimitive(leftType, "string") && IsPrimitive(rightType, "string") {
			return tc.primitives["string"]
		}
		tc.Err(fmt.Sprintf("invalid operands for %s: %s and %s", expr.Operator.Value, leftType, rightType))
		return nil
	case lexer.EQUALS, lexer.NOT_EQUALS:
		if !leftType.Equals(rightType) {
			tc.Err(fmt.Sprintf("cannot compare %s and %s", leftType, rightType))
			return nil
		}
		return tc.primitives["bool"]
	case lexer.LESS, lexer.LESS_EQUALS, lexer.GREATER, lexer.GREATER_EQUALS:
		if IsNumeric(leftType) && IsNumeric(rightType) {
			return tc.primitives["bool"]
		}
		tc.Err(fmt.Sprintf("invalid operands for %s: %s and %s", expr.Operator.Value, leftType, rightType))
		return nil

		// TODO: Add checks for all operators
	default:
		tc.Err(fmt.Sprintf("unsupported binary operator: %s", expr.Operator.Value))
		return nil
	}
}

func (tc *TypeChecker) CheckUnaryExpr(expr ast.UnaryExpr) Type {
	operandType := tc.InferType(expr.Rhs)
	if operandType == nil {
		return nil
	}
	switch expr.Operator.Type {
	case lexer.PLUS, lexer.DASH:
		if IsNumeric(operandType) {
			return operandType
		}
		tc.Err(fmt.Sprintf("invalid operand for %s: %s", expr.Operator.Value, operandType))
		return nil
	case lexer.NOT:
		if IsPrimitive(operandType, "bool") {
			return tc.primitives["bool"]
		}
		tc.Err(fmt.Sprintf("invalid operand for %s: %s", expr.Operator.Value, operandType))
		return nil
	default:
		tc.Err(fmt.Sprintf("unsupported unary operator: %s", expr.Operator.Value))
		return nil
	}
}

func (tc *TypeChecker) CheckFuncCallExpr(expr ast.FuncCallExpr) Type {
	funcType := tc.InferType(expr.Func)
	if funcType == nil {
		return nil
	}
	ft, ok := funcType.(FuncType)
	if !ok {
		tc.Err(fmt.Sprintf("cannot call non-function value of type %s", funcType))
		return nil
	}
	if len(expr.Args) != len(ft.ParamTypes) {
		tc.Err(fmt.Sprintf("wrong number of arguments, expected %d, found %d", len(ft.ParamTypes), len(expr.Args)))
		return nil
	}
	for i, arg := range expr.Args {
		argType := tc.InferType(arg)
		if argType == nil {
			return nil
		}
		if !ft.ParamTypes[i].Equals(argType) {
			tc.Err(fmt.Sprintf("argument %d type mismatch: expected %s, found %s", i+1, ft.ParamTypes[i], argType))
			return nil
		}
	}
	return ft.ReturnType
}

func (tc *TypeChecker) CheckStructLiteralExpr(expr ast.StructLiteralExpr) Type {
	var structType StructType
	structTypeValue := tc.InferType(expr.Struct)
	structType, ok := structTypeValue.(StructType)
	if !ok {
		tc.Err(fmt.Sprintf("expression of type %s cannot be used as a struct", structTypeValue))
		return nil
	}
	assignedMembers := make(map[string]bool, len(structType.Members))
	for memberName := range structType.Members {
		assignedMembers[memberName] = false
	}
	for _, member := range expr.Members {
		memberName, ok := member.Assigne.(ast.IdentExpr)
		if !ok {
			tc.Err(fmt.Sprintf("not a valid struct member identifier: %s", member.Assigne))
			continue
		}
		assigneType, ok := structType.Members[memberName.Value]
		if !ok {
			tc.Err(fmt.Sprintf("%s is not a member of struct %s", memberName.Value, structType.Name))
			continue
		}
		if assignedMembers[memberName.Value] == true {
			tc.Err(fmt.Sprintf("struct member %s assigned multiple times", memberName.Value))
			continue
		}
		assignedValueType := tc.InferType(member.AssignedValue)
		if assignedValueType == nil {
			continue
		}
		if !assigneType.Equals(assignedValueType) {
			tc.Err(fmt.Sprintf("cannot assign %s to %s of struct member %s", assignedValueType, assigneType, memberName.Value))
			continue
		}
		assignedMembers[memberName.Value] = true
	}
	for memberName, assigned := range assignedMembers {
		if !assigned {
			tc.Err(fmt.Sprintf("struct member %s is not assigned a value", memberName))
		}
	}
	return structType
}

func (tc *TypeChecker) CheckStructMemberExpr(expr ast.StructMemberExpr) Type {
	structTypeValue := tc.InferType(expr.Struct)
	structType, ok := structTypeValue.(StructType)
	if !ok {
		tc.Err(fmt.Sprintf("expression of type %s cannot be used as a struct", structTypeValue))
		return nil
	}
	memberName, ok := expr.Member.(ast.IdentExpr)
	if !ok {
		tc.Err(fmt.Sprintf("not a valid struct member identifier: %s", expr.Member))
	}
	memberType, ok := structType.Members[memberName.Value]
	if !ok {
		tc.Err(fmt.Sprintf("%s is not a member of struct %s", memberName.Value, structType.Name))
	}
	return memberType
}

func (tc *TypeChecker) CheckArrayIndexExpr(expr ast.ArrayIndexExpr) Type {
	// TODO
	// expr.Array must be a valid array
	// expr.Index must be a numeric type
	return nil
}

func (tc *TypeChecker) CheckAssignExpr(expr ast.AssignExpr) Type {
	assigneType := tc.InferType(expr.Assigne)
	assignedValueType := tc.InferType(expr.AssignedValue)
	if !assigneType.Equals(assignedValueType) {
		tc.Err(fmt.Sprintf("cannot assign %s to %s", assignedValueType, assigneType))
	}
	return assigneType
}
