package instance

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Launch spawns the Minecraft process (using a mock Java command for now)
func Launch(ctx context.Context, inst *Instance) error {
	// Construct the massive Java command line (Mocked placeholder)
	// We run "java -version" to prove subprocess execution and logging works.
	cmd := exec.Command("java", "-version")

	// Get pipes for stdout and stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start process: %w", err)
	}

	// Notify frontend that we are running
	runtime.EventsEmit(ctx, "instance:state", map[string]interface{}{
		"id":    inst.ID,
		"state": "Running",
	})

	// Async log readers
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			runtime.EventsEmit(ctx, "instance:log", scanner.Text())
			fmt.Println("[Launcher Log]", scanner.Text())
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			runtime.EventsEmit(ctx, "instance:log", scanner.Text())
			fmt.Println("[Launcher Log]", scanner.Text())
		}
	}()

	// Wait for process to exit in a goroutine
	go func() {
		err := cmd.Wait()
		state := "Stopped"
		if err != nil {
			state = "Crashed"
		}
		runtime.EventsEmit(ctx, "instance:state", map[string]interface{}{
			"id":    inst.ID,
			"state": state,
		})
	}()

	return nil
}
