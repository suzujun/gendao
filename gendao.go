package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"gopkg.in/urfave/cli.v1"

	"github.com/suzujun/gendao/commands"
)

func main() {

	userFlag := cli.StringFlag{
		Name:  "user",
		Usage: "User name to connect to mysql",
	}
	uFlag := userFlag
	uFlag.Name = "u"

	passwordFlag := cli.StringFlag{
		Name:  "password",
		Usage: "Password to connect to mysql",
	}
	pFlag := passwordFlag
	pFlag.Name = "p"

	databaseFlag := cli.StringFlag{
		Name:  "database",
		Usage: "target database name",
	}
	dFlag := databaseFlag
	dFlag.Name = "d"

	tableFlag := cli.StringFlag{
		Name:  "table",
		Usage: "target table names",
	}
	tFlag := databaseFlag
	tFlag.Name = "t"

	app := cli.NewApp()
	app.Name = "gendao"
	app.Usage = "make an dao and model source code for golang"
	app.Version = "0.1"

	app.Commands = []cli.Command{
		cli.Command{
			Name:   "init",
			Usage:  "Generate initialized config json",
			Action: initAction,
			Flags: []cli.Flag{
				uFlag, pFlag, dFlag,
				userFlag, passwordFlag, databaseFlag,
			},
		},
		cli.Command{
			Name:      "pull",
			Usage:     "Generate tables JSON from database",
			ArgsUsage: "{config file path}",
			Action:    pullAction,
			Flags:     []cli.Flag{dFlag, databaseFlag},
		},
		cli.Command{
			Name:      "addtype",
			Usage:     "Generate tables JSON from database",
			ArgsUsage: "{config file path}",
			Action:    addTypeAction,
		},
		cli.Command{
			Name:      "gen",
			Usage:     "Generate source code from JSON",
			ArgsUsage: "{config file path}",
			Action:    genAction,
			Flags:     []cli.Flag{dFlag, tFlag, databaseFlag, tableFlag},
		},
	}
	app.Run(os.Args)
}

func getFlag(c *cli.Context, names ...string) string {
	for _, name := range names {
		if v := c.String(name); v != "" {
			return v
		}
	}
	return ""
}

func initAction(c *cli.Context) error {
	user := getFlag(c, "user", "u")
	password := getFlag(c, "password", "p")
	dbname := getFlag(c, "database", "d")
	b, err := commands.GenerateConfigJSON(user, password, dbname)
	if err != nil {
		return err
	}
	fmt.Println(fmt.Sprintf("%s", b))
	return nil
}

func pullAction(c *cli.Context) error {
	path := c.Args().First()
	dbname := getFlag(c, "database", "d")
	cmd, err := getConfig(path, dbname)
	if err != nil {
		return err
	}
	if err := cmd.GenerateJSON(); err != nil {
		return err
	}
	fmt.Println("ok.")
	return nil
}

func addTypeAction(c *cli.Context) error {
	path := c.Args().First()
	dbname := getFlag(c, "database", "d")
	cmd, err := getConfig(path, dbname)
	if err != nil {
		return err
	}

	fmt.Println(`
This utility will walk you through setting a custom type for table column.
When inputting to the end, it saves it in the specified json file.
Press ^C at any time to quit.`)

	var key, typ, sampleValue, pkg, pkgAlias string

	for {
		fmt.Print("key: (table.column) ")
		fmt.Scanf("%s", &key)
		matched, err := regexp.MatchString("^\\w+\\.{1}\\w+$", key)
		if err != nil {
			return err
		}
		if matched {
			break
		}
		fmt.Println(fmt.Sprintf("Invalid key: \"%s\"", key))
	}

	// check duplicate
	if cmd.Config.CustomColumnType[key] != nil {
		var answer string
		fmt.Println("The selected key already exists.\nDo you want to overwrite? (yes) ")
		fmt.Scanf("%s", &answer)
		if answer != "yes" {
			return nil
		}
	}

	fmt.Print("type: ")
	fmt.Scanf("%s", &typ)
	fmt.Print("sample value: ")
	fmt.Scanf("%s", &sampleValue)

	fmt.Print("package: (github.com/path/to) ")
	fmt.Scanf("%s", &pkg)
	if pkg != "" {
		fmt.Print("package alias: ")
		fmt.Scanf("%s", &pkgAlias)
	}

	data := commands.CustomColumnType{
		Type:         typ,
		SampleValue:  sampleValue,
		Package:      pkg,
		PackageAlias: pkgAlias,
	}
	preview := map[string]interface{}{key: data}
	b, err := json.MarshalIndent(preview, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(fmt.Sprintf("About to write to \"%s\"", path))
	fmt.Println(string(b))

	var answer string
	fmt.Println("Is this ok? (yes) ")
	fmt.Scanf("%s", &answer)
	if answer != "yes" {
		return nil
	}

	// update config
	cmd.Config.CustomColumnType[key] = &data
	if err := cmd.Config.Write(path); err != nil {
		return err
	}

	fmt.Println("ok.")
	return nil
}

func genAction(c *cli.Context) error {
	path := c.Args().First()
	dbname := getFlag(c, "database", "d")
	cmd, err := getConfig(path, dbname)
	if err != nil {
		return err
	}
	table := getFlag(c, "table", "t")
	if err := cmd.GenerateSourceFromJSON(table); err != nil {
		return err
	}
	fmt.Println("run \"go fmt " + cmd.Config.OutputSourcePath + "/...\"")
	exec.Command("go", "fmt", cmd.Config.OutputSourcePath+"/...").Run()
	fmt.Println("ok.")
	return nil
}

func getGoPath() string {
	paths := strings.Split(os.Getenv("GOPATH"), ":")
	return paths[len(paths)-1]
}

func getConfig(path, dbName string) (*commands.Command, error) {
	if path == "" {
		fmt.Println("Please set the config.json created with the \"init\" command")
		return nil, errors.New("")
	}
	cmd, err := commands.GetCommand(path, dbName)
	if err != nil {
		return nil, err
	}
	if dbName != "" {
		cmd.Config.MysqlConfig.DbName = dbName
	}
	return cmd, nil
}
