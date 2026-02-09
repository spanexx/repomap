package parsing

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// ExtractDefinitions parses a Go file and returns a list of simplified definitions.
func ExtractDefinitions(filePath string) ([]string, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var definitions []string

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			definitions = append(definitions, formatFuncDecl(x))
			return false // Don't traverse inside function body
		case *ast.GenDecl:
			if x.Tok == token.TYPE {
				for _, spec := range x.Specs {
					if typeSpec, ok := spec.(*ast.TypeSpec); ok {
						definitions = append(definitions, formatTypeSpec(typeSpec))
					}
				}
			}
			return false
		}
		return true
	})

	return definitions, nil
}

func formatFuncDecl(decl *ast.FuncDecl) string {
	var sb strings.Builder
	sb.WriteString("func ")

	if decl.Recv != nil && len(decl.Recv.List) > 0 {
		sb.WriteString("(")

		field := decl.Recv.List[0]
		typeStr := formatType(field.Type)

		sb.WriteString(typeStr)
		sb.WriteString(") ")
	}

	sb.WriteString(decl.Name.Name)
	sb.WriteString("(")

	if decl.Type.Params != nil {
		sb.WriteString(formatFieldList(decl.Type.Params))
	}
	sb.WriteString(")")

	if decl.Type.Results != nil {
		sb.WriteString(" ")
		res := formatFieldList(decl.Type.Results)

		// If multiple return values (indicated by comma), wrap in parens
		if strings.Contains(res, ",") {
			sb.WriteString("(" + res + ")")
		} else {
			sb.WriteString(res)
		}
	}

	return sb.String()
}

func formatTypeSpec(spec *ast.TypeSpec) string {
	var sb strings.Builder
	sb.WriteString("type ")
	sb.WriteString(spec.Name.Name)
	sb.WriteString(" ")

	switch spec.Type.(type) {
	case *ast.StructType:
		sb.WriteString("struct")
	case *ast.InterfaceType:
		sb.WriteString("interface")
	default:
		sb.WriteString(formatType(spec.Type))
	}

	return sb.String()
}

func formatFieldList(fields *ast.FieldList) string {
	var parts []string
	for _, field := range fields.List {
		typeStr := formatType(field.Type)

		count := len(field.Names)
		if count == 0 {
			count = 1
		}

		for i := 0; i < count; i++ {
			parts = append(parts, typeStr)
		}
	}
	return strings.Join(parts, ", ")
}

func formatType(expr ast.Expr) string {
	if expr == nil {
		return ""
	}
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + formatType(t.X)
	case *ast.SelectorExpr:
		return formatType(t.X) + "." + t.Sel.Name
	case *ast.ArrayType:
		lenStr := ""
		if t.Len != nil {
			if lit, ok := t.Len.(*ast.BasicLit); ok {
				lenStr = lit.Value
			}
		}
		return "[" + lenStr + "]" + formatType(t.Elt)
	case *ast.MapType:
		return "map[" + formatType(t.Key) + "]" + formatType(t.Value)
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.StructType:
		return "struct{}"
	case *ast.FuncType:
		return "func"
	case *ast.Ellipsis:
		return "..." + formatType(t.Elt)
	case *ast.ChanType:
		return "chan " + formatType(t.Value)
	default:
		return fmt.Sprintf("%T", t)
	}
}
