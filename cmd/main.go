package main

import (
	"github.com/urfave/cli"
	"os"
)

/**

 查找包含了git的文件
 find ./ -name .git*
 find ./ -name .git*|xargs rm -rf

 find ./ -name "*.DS_Store"


 find ./ -name .git*
 find ./ -name .git*|xargs rm -rf

 find ./ -name *.git
 find ./ -name *.git|xargs rm -rf


 rm -rf $(find ./ -name .git*)

./go-vendors-tool  -p /usr/local/gopath/src/farm -p /usr/local/gopath/src/farm-cr -o /tmp/go-vendors

 */

func main() {
	app := cli.NewApp()
	app.Name = "go vendor aggregator"
	app.Version = "1.0.0"
	app.Usage = "go vendor aggregator"
	app.Action = action
	app.Flags = flags
	app.Run(os.Args)
}

func action(c *cli.Context) error {
	projects := c.StringSlice(OptProject)
	outputDir := c.String(OptOutputDir)
	return aggregate(projects, outputDir)
}
