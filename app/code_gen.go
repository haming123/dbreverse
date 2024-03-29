package app

import (
	"bytes"
	"fmt"
	"github.com/haming123/wego/worm"
	"io/ioutil"
	"strings"
)

func FirstToUpper(s string) string {
	if len(s) < 1 {
		return s
	}

	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if i == 0 && 'a' <= c && c <= 'z' {
			c -= 'a' - 'A'
		}
		b.WriteByte(c)
	}

	return b.String()
}

func gen_model_header(dialect worm.Dialect, flds []worm.ColumnInfo) string {
	has_time := false
	for _, field := range flds {
		if dialect.DbType2GoType(field.DbType) == "time.Time" {
			has_time = true
			break
		}
	}

	var buff bytes.Buffer
	buff.WriteString("package ")
	buff.WriteString(AppCfg.PkgName)
	buff.WriteString("\n")
	if has_time {
		buff.WriteString("\nimport \"time\"\n")
	}
	return buff.String()
}

func gen_model_struct(dialect worm.Dialect, flds []worm.ColumnInfo, table_name string) string {
	strs := strings.Split(table_name, ".")
	if len(strs) == 2 {
		table_name = strs[1]
	}

	var buff bytes.Buffer
	buff.WriteString("\ntype ")
	buff.WriteString(FirstToUpper(table_name))
	buff.WriteString(" struct {\n")
	for _, field := range flds {
		go_type := dialect.DbType2GoType(field.DbType)
		if go_type == "int32" {
			go_type = "int"
		}
		if go_type == "int" && field.Length > 10 {
			go_type = "int64"
		}

		buff.WriteString("\t//")
		buff.WriteString(fmt.Sprintf("%s:", field.SQLType))
		if len(field.Comment) > 0 {
			commit := field.Comment
			commit = strings.ReplaceAll(commit, "\n", "")
			buff.WriteString(commit)
		}
		buff.WriteString("\n")
		if AppCfg.UseTag {
			name_str := fmt.Sprintf("\t%-20s", FirstToUpper(field.Name))
			buff.WriteString(name_str)
			go_type_str := fmt.Sprintf("\t%-6s", go_type)
			buff.WriteString(go_type_str)
			buff.WriteString("\t`db:\"")
			buff.WriteString(field.Name)
			if field.IsAutoIncrement {
				buff.WriteString(";autoincr")
			} else if AppCfg.CreateTime != "" && field.Name == AppCfg.CreateTime {
				buff.WriteString(";n_update")
			}
			buff.WriteString("\"`\n")
		} else {
			name_str := fmt.Sprintf("\t%-20s", "DB_"+field.Name)
			buff.WriteString(name_str)
			go_type_str := fmt.Sprintf("\t%-6s", go_type)
			buff.WriteString(go_type_str)
			buff.WriteString("\n")
		}
	}
	buff.WriteString("}\n")

	return buff.String()
}

/*
func (ent *DB_User)TableName() string {
	return "user"
}
*/
func gen_func_table_name(table_name string) string {
	strs := strings.Split(table_name, ".")
	if len(strs) == 2 {
		table_name = strs[1]
	}
	struct_name := FirstToUpper(table_name)

	var buff bytes.Buffer

	buff.WriteString("\nfunc (ent *")
	buff.WriteString(struct_name)
	buff.WriteString(") TableName() string {\n")
	buff.WriteString("\treturn \"")
	buff.WriteString(table_name)
	buff.WriteString("\"\n}\n\n")

	return buff.String()
}

/*
var G_user User
var P_user = worm.NewModelPool(G_user)
*/
func gen_func_model_pool(table_name string) string {
	strs := strings.Split(table_name, ".")
	if len(strs) == 2 {
		table_name = strs[1]
	}
	var_name := strings.ToLower(table_name)
	struct_name := FirstToUpper(table_name)

	var buff bytes.Buffer
	buff.WriteString(fmt.Sprintf("var G_%s %s\n", var_name, struct_name))
	buff.WriteString(fmt.Sprintf("var P_%s = worm.NewModelPool(G_%s)\n", var_name, var_name))
	return buff.String()
}

func gen_model_code(db *worm.DbEngine, table_name string) (string, error) {
	dialect := db.GetDialect()
	flds, err := dialect.GetColumns(db.DB(), table_name)
	if err != nil {
		return "", err
	}
	//fmt.Println(flds)

	code_str := gen_model_header(dialect, flds)
	code_str += gen_model_struct(dialect, flds, table_name)
	code_str += gen_func_table_name(table_name)
	if AppCfg.UsePool {
		code_str += gen_func_model_pool(table_name)
	}
	return code_str, nil
}

func CodeGen4Table(db *worm.DbEngine, table_name string, file_name string) error {
	code_str, err := gen_model_code(db, table_name)
	if err != nil {
		return err
	}

	if file_name == "" {
		fmt.Println(code_str)
		return nil
	} else {
		return ioutil.WriteFile(file_name, []byte(code_str), 0666)
	}
}
