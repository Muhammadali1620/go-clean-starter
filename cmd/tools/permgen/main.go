package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"new_project/internal/core/utils"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	log.Println("Starting Permission Generation (permgen)...")

	fset := token.NewFileSet()
	entries, err := os.ReadDir("internal/models")
	if err != nil {
		log.Fatalf("Failed to read models directory: %v", err)
	}

	var permissions []string

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") {
			continue
		}

		filePath := filepath.Join("internal/models", entry.Name())
		file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
		if err != nil {
			log.Printf("Warning: failed to parse %s: %v", filePath, err)
			continue
		}

		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.TYPE {
				continue
			}

			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}

				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					continue
				}

				modelName := strings.ToUpper(camelToSnake(typeSpec.Name.Name))

				// 1. Generate Default CRUD Permissions
				permissions = append(permissions, fmt.Sprintf("%s_CREATE", modelName))
				permissions = append(permissions, fmt.Sprintf("%s_READ", modelName))
				permissions = append(permissions, fmt.Sprintf("%s_UPDATE", modelName))
				permissions = append(permissions, fmt.Sprintf("%s_DELETE", modelName))

				// 2. Check for @manager-access tag
				isManagerAccess := false
				if genDecl.Doc != nil {
					for _, comment := range genDecl.Doc.List {
						if strings.Contains(comment.Text, "@manager-access") {
							isManagerAccess = true
							break
						}
					}
				}

				// 3. If @manager-access tag is present, generate permissions for club manager
				if isManagerAccess {
					permissions = append(permissions, fmt.Sprintf("CMANAGER_%s_CREATE", modelName))
					permissions = append(permissions, fmt.Sprintf("CMANAGER_%s_READ", modelName))
					permissions = append(permissions, fmt.Sprintf("CMANAGER_%s_UPDATE", modelName))
					permissions = append(permissions, fmt.Sprintf("CMANAGER_%s_DELETE", modelName))
				}

				// 4. Look for custom permissions in struct docs
				if genDecl.Doc != nil {
					for _, comment := range genDecl.Doc.List {
						if strings.HasPrefix(comment.Text, "// @permission:") {
							customPerm := strings.TrimSpace(strings.TrimPrefix(comment.Text, "// @permission:"))
							permissions = append(permissions, strings.ToUpper(customPerm))
						}
					}
				}

				// 5. Look for custom permissions in struct fields
				if structType.Fields != nil {
					for _, field := range structType.Fields.List {
						if field.Doc != nil {
							for _, comment := range field.Doc.List {
								if strings.HasPrefix(comment.Text, "// @permission:") {
									customPerm := strings.TrimSpace(strings.TrimPrefix(comment.Text, "// @permission:"))
									permissions = append(permissions, strings.ToUpper(customPerm))
								}
							}
						}
					}
				}
			}
		}
	}

	generateSQL(permissions)
}

func camelToSnake(s string) string {
	var result []rune
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return string(result)
}

func generateSQL(permissions []string) {
	permMap := make(map[string]bool)
	var uniquePerms []string
	for _, p := range permissions {
		if !permMap[p] {
			permMap[p] = true
			uniquePerms = append(uniquePerms, p)
		}
	}

	var sb strings.Builder
	sb.WriteString("-- Auto-generated permissions seed file\n")
	sb.WriteString("INSERT INTO permissions (code, description) VALUES\n")

	for i, p := range uniquePerms {
		desc := strings.ReplaceAll(p, "_", " ")
		sb.WriteString(fmt.Sprintf("\t('%s', '{\"en\": \"%s\"}')", p, desc))

		if i < len(uniquePerms)-1 {
			sb.WriteString(",\n")
		} else {
			sb.WriteString("\nON CONFLICT (code) DO NOTHING;\n")
		}
	}

	filename := fmt.Sprintf("migrations/%s_seed_permissions.up.sql", utils.NowUTC().Format("20060102150405"))
	err := os.WriteFile(filename, []byte(sb.String()), 0644)
	if err != nil {
		log.Fatalf("Failed to write SQL file: %v", err)
	}

	log.Printf("Successfully generated %d permissions in %s\n", len(uniquePerms), filename)
}
