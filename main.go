package main

import (
	"flag"
	"fmt"
	"github.com/kazu22002/dbRelation/config"
	"github.com/kazu22002/dbRelation/repository"
	"io/ioutil"
	"os"

	"github.com/go-yaml/yaml"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	configFlag = flag.String("config",
		"./config/config.toml",
		"-config=\"path-to-your-config-file\" ")

	exclude = []string{
		"admin_authorities",
	}
)

// NewConfig ...
func NewConfig() (*config.Context, error) {
	flags := &config.Flags{
		CfgFileName: *configFlag,
	}
	return config.New(flags)
}

func main() {
	ct, _ := NewConfig()

	db, err := gorm.Open(ct.Config.Database.Db, ct.Config.Database.Dsn)
	defer db.Close()

	if err != nil {
		fmt.Println(err)
	}

	tables := repository.Table(db)

	type Output struct {
		Table         string   `yaml:"table"`
		Columns       []string `yaml:"columns"`
		ParentTable   string   `yaml:"parentTable"`
		ParentColumns []string `yaml:"parentColumns"`
//		Def           string   `yaml:"def"`
	}

	var output []Output
	for _, v := range tables {
		for _, match := range exclude {
			if match == v.Table.TableName {
				continue
			}
		}
		for _, vv := range tables {
			if v.Table.TableName == vv.Table.TableName {
				// 同じテーブル判定
				continue
			}

			for _, column := range vv.Column {
				if v.TableNameId == column.ColumnName {
					t := Output{
						Table:         vv.Table.TableName,
						Columns:       []string{v.TableNameId},
						ParentTable:   v.Table.TableName,
						ParentColumns: []string{"id"},
					}
					output = append(output, t)
				}
			}
		}
	}

	err = WriteOnFile("relation.yml", output)
}

func WriteOnFile(fileName string, data interface{}) error {
	// ここでデータを []byte に変換しています。
	buf, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	// []byte をファイルに上書きしています。
	err = ioutil.WriteFile(fileName, buf, 0644)
	if err != nil {
		return err
	}
	return nil
}
