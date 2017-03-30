package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/suzujun/gendao/commands"

	"errors"

	"gopkg.in/urfave/cli.v1"
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
