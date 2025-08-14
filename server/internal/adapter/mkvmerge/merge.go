package mkvmerge

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/kkiling/goplatform/log"
	"github.com/samber/lo"
)

type Merge struct {
	logger log.Logger
}

func NewMerge(logger log.Logger) *Merge {
	return &Merge{
		logger: logger.Named("mkvmerge"),
	}
}

// вспомогательная функция для разбивки на строки
func splitLines(s string) []string {
	var lines []string
	current := ""
	for _, r := range s {
		if r == '\n' {
			lines = append(lines, current)
			current = ""
		} else {
			current += string(r)
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}

func (s *Merge) Merge(ctx context.Context, params MergeParams, outputChan chan<- OutputMessage) error {
	//info, err := s.GetMediaInfo(params.VideoInputFile)
	//if err != nil {
	//	return fmt.Errorf("get media info: %w", err)
	//}
	//fmt.Printf("%+v\n", info)

	// Проверка существования основного видеофайла
	var err error
	if _, err = os.Stat(params.VideoInputFile); os.IsNotExist(err) {
		return fmt.Errorf("input video file does not exist: %s", params.VideoInputFile)
	}

	// Проверка аудиодорожек
	for _, track := range params.AudioTracks {
		if _, err = os.Stat(track.Path); os.IsNotExist(err) {
			return fmt.Errorf("audio track file does not exist: %s", track.Path)
		}
	}

	// Проверка субтитров
	for _, track := range params.SubtitleTracks {
		if _, err = os.Stat(track.Path); os.IsNotExist(err) {
			return fmt.Errorf("subtitle file does not exist: %s", track.Path)
		}
	}

	args := []string{"-o", filepath.Clean(params.VideoOutputFile)}

	// Снимаем default со всех старых аудио в исходном файле
	//for id, _ := range info.AudioTracks {
	//	// audio.Number — это номер трека в контейнере (1-based), mkvmerge ждёт 0-based
	//	//trackID := audio.Number - 1
	//	args = append(args, "--default-track", fmt.Sprintf("%d:no", id))
	//}

	args = append(args, filepath.Clean(params.VideoInputFile))

	// Добавляем аудиодорожки
	for _, track := range params.AudioTracks {
		{
			if track.Language != "" {
				args = append(args, "--language", "0:"+track.Language)
			}
		}
		args = append(args,
			"--track-name", "0:"+track.Name,
			"--default-track", fmt.Sprintf("0:%s", lo.Ternary(track.Default, "yes", "no")),
			filepath.Clean(track.Path), // Путь к файлу идет ПОСЛЕ флагов!
		)
	}

	// Добавляем субтитры
	for _, track := range params.SubtitleTracks {
		if track.Language != "" {
			args = append(args, "--language", "0:"+track.Language)
		}
		args = append(args,
			"--track-name", "0:"+track.Name,
			"--default-track", fmt.Sprintf("0:%s", lo.Ternary(track.Default, "yes", "no")),
			filepath.Clean(track.Path), // Путь к файлу идет ПОСЛЕ флагов!
		)
	}

	// Для отладки
	debugMsg := "mkvmerge " + strings.Join(args, " ")
	outputChan <- OutputMessage{Type: InfoMessageType, Content: debugMsg}

	cmd := exec.CommandContext(ctx, "mkvmerge", args...)
	//fullArgs := append([]string{"abc", "mkvmerge"}, args...)
	//cmd := exec.CommandContext(ctx, "s6-setuidgid", fullArgs...)

	// Настраиваем пайпы
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("error creating stdout pipe: %v", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("error creating stderr pipe: %v", err)
	}

	// Запускаем команду
	if errStart := cmd.Start(); errStart != nil {
		return fmt.Errorf("error starting command: %v", errStart)
	}

	// Читаем вывод в реальном времени и отправляем в канал
	var wg sync.WaitGroup
	go func() {
		wg.Add(1)
		defer wg.Done()
		s.scanOutput(ctx, stdoutPipe, outputChan, InfoMessageType)
	}()
	go func() {
		wg.Add(1)
		defer wg.Done()
		s.scanOutput(ctx, stderrPipe, outputChan, ErrorMessageType)
	}()

	fmt.Println("*******************")
	fmt.Println("*******************")
	fmt.Println("*******************")

	pid := cmd.Process.Pid
	fmt.Println("PID запущенного процесса:", pid)

	// Выполняем ps, чтобы узнать пользователя
	psCmd := exec.Command("ps", "-o", "user=", "-p", strconv.Itoa(pid))
	out, err := psCmd.Output()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Процесс выполняется от пользователя: %s", out)

	// Получаем UID и GID через /proc/<pid>/status (Linux)
	statusPath := fmt.Sprintf("/proc/%d/status", pid)
	data, err := os.ReadFile(statusPath)
	if err != nil {
		panic(err)
	}

	var uid, gid string
	for _, line := range splitLines(string(data)) {
		if len(line) > 4 && line[:4] == "Uid:" {
			uid = line
		}
		if len(line) > 4 && line[:4] == "Gid:" {
			gid = line
		}
	}

	fmt.Println(uid)
	fmt.Println(gid)

	fmt.Println("*******************")
	fmt.Println("*******************")
	fmt.Println("*******************")

	// Ждем завершения
	if errWait := cmd.Wait(); errWait != nil {
		return fmt.Errorf("mkvmerge failed: %w", errWait)
	}
	wg.Wait()

	return nil
}

// Измененная функция scanOutput теперь является методом Service и принимает канал
func (s *Merge) scanOutput(ctx context.Context, reader io.Reader, outputChan chan<- OutputMessage, msgType MessageType) {
	buf := make([]byte, 1024)
	var leftover []byte

	for {
		n, err := reader.Read(buf)
		if n > 0 {
			data := append(leftover, buf[:n]...)
			lines := bytes.Split(data, []byte{'\r'})

			// Последний элемент может быть неполной строкой
			for i, line := range lines {
				if i == len(lines)-1 {
					leftover = line
					continue
				}

				// Обрабатываем только непустые строки
				if len(line) > 0 {
					outputChan <- OutputMessage{
						Type:    msgType,
						Content: string(line),
					}
				}
			}
		}

		if err != nil {
			if err != io.EOF {
				outputChan <- OutputMessage{
					Type:    ErrorMessageType,
					Content: fmt.Sprintf("read error: %v", err),
				}
			}
			break
		}
	}

	// Выводим оставшиеся данные
	if len(leftover) > 0 {
		outputChan <- OutputMessage{
			Type:    msgType,
			Content: string(leftover),
		}
	}
}
