package judge

import (
	"JudgeCore/internal/global"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types/build"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/moby/go-archive"
)

// BuildImage 构建Sandbox
func BuildImage(currentDir string) bool {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Println("[FeasOJ] Error creating Docker client: ", err)
		return false
	}

	tar, err := archive.TarWithOptions(currentDir, &archive.TarOptions{})
	if err != nil {
		log.Println("[FeasOJ] Error creating tar: ", err)
		return false
	}

	buildOptions := build.ImageBuildOptions{
		Context:    tar,
		Dockerfile: "Sandbox",
		Tags:       []string{"judgecore:latest"},
	}

	log.Println("[FeasOJ] SandBox is being built...")
	buildResponse, err := cli.ImageBuild(ctx, tar, buildOptions)
	if err != nil {
		log.Println("[FeasOJ] Error building Docker image: ", err)
		return false
	}
	defer buildResponse.Body.Close()

	_, err = io.Copy(log.Writer(), buildResponse.Body)
	if err != nil {
		log.Printf("[FeasOJ] Error copying build response: %v", err)
	}

	return true
}

// CompileAndRun 编译并运行代码
func CompileAndRun(filename string, containerID string, problem *global.Problem, testCases []*global.TestCaseRequest) string {
	taskDir := fmt.Sprintf("/workspace/task_%d", time.Now().UnixNano())

	mkdirCmd := exec.Command("docker", "exec", containerID, "mkdir", "-p", taskDir)
	if err := mkdirCmd.Run(); err != nil {
		return global.SystemError
	}

	copyCmd := exec.Command("docker", "exec", containerID, "cp", fmt.Sprintf("/workspace/%s", filename), taskDir)
	if err := copyCmd.Run(); err != nil {
		return global.SystemError
	}

	defer func() {
		if err := resetTaskDirectory(containerID, taskDir); err != nil {
			log.Printf("[FeasOJ] Reset task dir %s error: %v", taskDir, err)
		}
	}()

	ext := filepath.Ext(filename)
	var compileCmd *exec.Cmd

	timeLimitSeconds, memoryLimitKB, err := parseLimits(problem)
	if err != nil {
		return global.SystemError
	}

	switch ext {
	case ".cpp":
		compileCmd = exec.Command("docker", "exec", containerID, "sh", "-c",
			fmt.Sprintf("g++ %s/%s -o %s/%s.out", taskDir, filename, taskDir, filename))
	case ".java":
		renameCmd := exec.Command("docker", "exec", containerID, "sh", "-c",
			fmt.Sprintf("mv %s/%s %s/Main.java", taskDir, filename, taskDir))
		if err := renameCmd.Run(); err != nil {
			return global.CompileError
		}
		compileCmd = exec.Command("docker", "exec", containerID, "sh", "-c",
			fmt.Sprintf("javac %s/Main.java", taskDir))
	case ".rs":
		exeName := strings.TrimSuffix(filename, ".rs")
		compileCmd = exec.Command("docker", "exec", containerID, "sh", "-c",
			fmt.Sprintf("rustc %s/%s -o %s/%s", taskDir, filename, taskDir, exeName))
	case ".php":
		compileCmd = exec.Command("docker", "exec", containerID, "sh", "-c",
			fmt.Sprintf("php -l %s/%s", taskDir, filename))
	case ".pas":
		exeName := strings.TrimSuffix(filename, ".pas")
		compileCmd = exec.Command("docker", "exec", containerID, "sh", "-c",
			fmt.Sprintf("fpc -v0 -O2 %s/%s -o%s/%s", taskDir, filename, taskDir, exeName))
	}

	if compileCmd != nil {
		if err := compileCmd.Run(); err != nil {
			return global.CompileError
		}
	}

	for _, testCase := range testCases {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeLimitSeconds+1)*time.Second)
		defer cancel()

		cmdStr := buildRunCommand(ext, filename, taskDir, timeLimitSeconds, memoryLimitKB)
		if cmdStr == "" {
			return global.SystemError
		}

		runCmd := exec.CommandContext(ctx, "docker", "exec", "-i", containerID, "sh", "-c", cmdStr)
		runCmd.Stdin = strings.NewReader(testCase.InputData)
		output, err := runCmd.CombinedOutput()

		if ctx.Err() == context.DeadlineExceeded {
			return global.TimeLimitExceeded
		}

		if err != nil {
			var exitErr *exec.ExitError
			if errors.As(err, &exitErr) {
				switch exitErr.ExitCode() {
				case 124: // timeout 触发
					return global.TimeLimitExceeded
				case 137: // SIGKILL, 可能是内存超限
					return global.MemoryLimitExceeded
				}
			}
			return global.RuntimeError
		}

		expectedOutput := strings.TrimSpace(testCase.OutputData)
		actualOutput := strings.TrimSpace(string(output))

		if actualOutput != expectedOutput {
			return global.WrongAnswer
		}
	}

	return global.Accepted
}

// TerminateContainer 终止并删除Docker容器
func TerminateContainer(containerID string) {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Printf("[FeasOJ] Error creating Docker client for termination: %v", err)
		return
	}

	if err := cli.ContainerStop(ctx, containerID, container.StopOptions{}); err != nil {
		log.Printf("[FeasOJ] Error stopping container %s: %v", containerID, err)
	}
}

func resetTaskDirectory(containerID, taskDir string) error {
	resetCmd := exec.Command("docker", "exec", containerID, "sh", "-c", fmt.Sprintf("rm -rf %s", taskDir))
	if err := resetCmd.Run(); err != nil {
		log.Printf("[FeasOJ] Error cleaning task directory %s in container %s: %v", taskDir, containerID, err)
		return err
	}
	return nil
}

func parseLimits(problem *global.Problem) (timeLimit int, memoryLimit int, err error) {
	re := regexp.MustCompile(`\d+`)

	timeMatches := re.FindAllString(problem.Timelimit, -1)
	if len(timeMatches) == 0 {
		return 0, 0, fmt.Errorf("no time limit found")
	}
	timeLimit, err = strconv.Atoi(timeMatches[0])
	if err != nil {
		return 0, 0, err
	}

	memMatches := re.FindAllString(problem.Memorylimit, -1)
	if len(memMatches) == 0 {
		return 0, 0, fmt.Errorf("no memory limit found")
	}
	memoryLimitMB, err := strconv.Atoi(memMatches[0])
	if err != nil {
		return 0, 0, err
	}
	memoryLimit = memoryLimitMB * 1024 // 转换为KB

	return timeLimit, memoryLimit, nil
}

func buildRunCommand(ext, filename, taskDir string, timeLimit, memoryLimit int) string {
	switch ext {
	case ".cpp":
		return fmt.Sprintf("ulimit -v %d && timeout -s SIGKILL %ds %s/%s.out", memoryLimit, timeLimit, taskDir, filename)
	case ".java":
		heapSizeMB := max(memoryLimit/1024, 32)
		return fmt.Sprintf(
			"ulimit -v %d && timeout -s SIGKILL %ds java -cp %s -Xms%dm -Xmx%dm -XX:MaxRAMPercentage=80.0 Main",
			memoryLimit, timeLimit, taskDir, heapSizeMB, heapSizeMB,
		)
	case ".py":
		return fmt.Sprintf("ulimit -v %d && timeout -s SIGKILL %ds python %s/%s", memoryLimit, timeLimit, taskDir, filename)
	case ".rs":
		exeName := strings.TrimSuffix(filename, ".rs")
		return fmt.Sprintf("ulimit -v %d && timeout -s SIGKILL %ds %s/%s",
			memoryLimit, timeLimit, taskDir, exeName)
	case ".php":
		return fmt.Sprintf("ulimit -v %d && timeout -s SIGKILL %ds php %s/%s",
			memoryLimit, timeLimit, taskDir, filename)
	case ".pas":
		exeName := strings.TrimSuffix(filename, ".pas")
		return fmt.Sprintf("ulimit -v %d && timeout -s SIGKILL %ds %s/%s",
			memoryLimit, timeLimit, taskDir, exeName)
	default:
		return ""
	}
}