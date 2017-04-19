package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"

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
		{
			Name:   "init",
			Usage:  "Generate initialized config json",
			Action: initAction,
			Flags: []cli.Flag{
				uFlag, pFlag, dFlag,
				userFlag, passwordFlag, databaseFlag,
			},
		},
		{
			Name:      "pull",
			Usage:     "Generate tables JSON from database",
			ArgsUsage: "{config file path}",
			Action:    pullAction,
			Flags:     []cli.Flag{dFlag, databaseFlag},
		},
		{
			Name:      "addtype",
			Usage:     "Generate tables JSON from database",
			ArgsUsage: "{config file path}",
			Action:    addTypeAction,
		},
		{
			Name:      "gen",
			Usage:     "Generate source code from JSON",
			ArgsUsage: "{config file path}",
			Action:    genAction,
			Flags:     []cli.Flag{dFlag, tFlag, databaseFlag, tableFlag},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
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
	b, err := commands.NewConfig(user, password, dbname).ExportJSON()
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", b)
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

var keyReg = regexp.MustCompile(`^\\w+\\.{1}\\w+$`)

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
		if _, err := fmt.Scan(&key); err != nil {
			return err
		}
		if keyReg.MatchString(key) {
			break
		}
		fmt.Println(fmt.Sprintf("Invalid key: \"%s\"", key))
	}

	// check duplicate
	if cmd.Config.CustomColumnType[key] != nil {
		var answer string
		fmt.Println("The selected key already exists.\nDo you want to overwrite? (yes) ")
		if _, err := fmt.Scan(&answer); err != nil {
			return err
		}
		if answer != "yes" {
			return nil
		}
	}

	fmt.Print("type: ")
	if _, err := fmt.Scan(&typ); err != nil {
		return err
	}
	fmt.Print("sample value: ")
	if _, err := fmt.Scan(&sampleValue); err != nil {
		return err
	}

	fmt.Print("package: (github.com/path/to) ")
	if _, err := fmt.Scan(&pkg); err != nil {
		return err
	}
	if pkg != "" {
		fmt.Print("package alias: ")
		if _, err := fmt.Scan(&pkgAlias); err != nil {
			return err
		}
	}

	data := commands.CustomColumnType{
		Type:         typ,
		SampleValue:  sampleValue,
		Package:      pkg,
		PackageAlias: pkgAlias,
	}
	preview := map[string]commands.CustomColumnType{key: data}
	b, err := json.MarshalIndent(preview, "", "  ")
	if err != nil {
		return err
	}
	fmt.Printf("About to write to \"%s\"\n", path)
	fmt.Println(string(b))

	var answer string
	fmt.Println("Is this ok? (yes) ")
	if _, err := fmt.Scan(&answer); err != nil {
		return err
	}
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
	if err := exec.Command("go", "fmt", cmd.Config.OutputSourcePath+"/...").Run(); err != nil {
		return err
	}
	fmt.Println("ok.")
	return nil
}

func getConfig(path, dbName string) (*commands.Command, error) {
	if path == "" {
		fmt.Println("Please set the config.json created with the \"init\" command")
		return nil, errors.New("config path is empty")
	}
	cmd, err := commands.NewCommandFromJSON(path, dbName)
	if err != nil {
		return nil, err
	}
	if dbName != "" {
		cmd.Config.MysqlConfig.DbName = dbName
	}
	return cmd, nil
}
