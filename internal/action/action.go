package action

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"
)

// Initialize is used to create a new machine and setup any dependencies
func Initialize(c *cli.Context) error {
	machine := c.String("machine")

	fmt.Println("Creating a new machine:", machine)

	multipass := fmt.Sprintf("%s", c.Context.Value("multipass"))

	// create the machine
	cmd := exec.Command(multipass, "launch", "--name", machine, "--cloud-init", "./cloud-init.yaml")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	// if we are bootstrapping, call that command
	if c.Bool("bootstrap") {
		return c.App.RunContext(c.Context, []string{c.App.Name, "--machine", c.String("machine"), "bootstrap"})
	}

	return nil
}

func Bootstrap(c *cli.Context, e CommandLineExecutor) error {
	machine := c.String("machine")
	php := c.String("php-version")
	database := c.String("database")

	args := []string{"multipass", "exec", machine, "--", "sudo", "bash", "/opt/nitro/bootstrap.sh", php, database}

	return e.Exec(e.Path(), args, os.Environ())
}

// Update will perform system updates on a given machine
func Update(c *cli.Context) error {
	machine := c.String("machine")
	multipass := fmt.Sprintf("%s", c.Context.Value("multipass"))

	fmt.Println("Updating machine:", machine)

	cmd := exec.Command(multipass, "exec", machine, "--", "sudo", "apt", "update", "-y")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func AddHost(c *cli.Context, e CommandLineExecutor) error {
	machine := c.String("machine")
	host := c.Args().First()
	php := c.String("php-version")

	if host == "" {
		return errors.New("missing param host")
	}

	if php == "" {
		fmt.Println("missing php-version")
		php = "7.4"
	}

	fmt.Println("Connecting to machine:", machine)

	args := []string{"multipass", "exec", machine, "--", "sudo", "bash", "/opt/nitro/nginx/add-site.sh", host, php}

	return e.Exec(e.Path(), args, os.Environ())
}

// SSH will login a user to a specific machine
func SSH(m string, e CommandLineExecutor) error {
	fmt.Println("Connecting to machine:", m)

	args := []string{"multipass", "shell", m}
	err := e.Exec(e.Path(), args, os.Environ())
	if err != nil {
		return err
	}

	return nil
}

func Delete(c *cli.Context) error {
	machine := c.String("machine")

	fmt.Println("Deleting machine:", machine)

	multipass := fmt.Sprintf("%s", c.Context.Value("multipass"))
	cmd := exec.Command(multipass, "delete", machine)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func Stop(c *cli.Context) error {
	machine := c.String("machine")

	fmt.Println("Stopping machine:", machine)

	multipass := fmt.Sprintf("%s", c.Context.Value("multipass"))
	cmd := exec.Command(multipass, "stop", machine)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}