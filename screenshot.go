package main

import (
	"log"
	"os"

	"github.com/pkg/browser"
	"github.com/urfave/cli"
)

func anyFuncs(funcs ...func() error) error {
	var err error
	for _, f := range funcs {
		if err = f(); err == nil {
			return nil
		}
	}
	return err
}

func takeScreenshot(serial, output string) error {
	execOut := func() error {
		c := adbCommand(serial, "exec-out", "screencap", "-p")
		imgfile, err := os.Create(output)
		if err != nil {
			return err
		}
		defer func() {
			imgfile.Close()
			if err != nil {
				os.Remove(output)
			}
		}()
		c.Stdout = imgfile
		return c.Run()
	}
	screencap := func() error {
		tmpPath := "/sdcard/fa-screenshot.png"
		c := adbCommand(serial, "shell", "screencap", "-p", tmpPath)
		if err := c.Run(); err != nil {
			return err
		}
		defer adbCommand(serial, "shell", "rm", tmpPath).Run()
		return adbCommand(serial, "pull", tmpPath, output).Run()
	}
	return anyFuncs(execOut, screencap)
}

func actScreenshot(ctx *cli.Context) (err error) {
	serial, err := chooseOne()
	if err != nil {
		return err
	}
	output := ctx.String("output")
	err = takeScreenshot(serial, output)
	if err == nil {
		log.Println("saved to", output)
		if ctx.Bool("open") {
			browser.OpenFile(output)
		}
	}
	return err
}
