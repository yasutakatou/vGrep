package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	ps "github.com/Tobotobo/powershell"
	"github.com/gdamore/tcell/v2"
)

type configData struct {
	LABEL string
	COUNT int
	MAX   int
	VOICE string
}

var (
	configs []configData
)

func main() {
	configs = nil

	if len(os.Args) < 2 {
		fmt.Println("使用法: go run main.go <ファイル名>")
		os.Exit(1)
	}

	if loadConfig("vGrep.ini") == false {
		log.Fatalf("Fail to read config file")
		os.Exit(1)
	}

	// ファイルの読み込み
	lines, err := readLines(os.Args[1])
	if err != nil {
		log.Fatalf("ファイルの読み込みに失敗しました: %v", err)
	}

	// 画面の初期化
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%v", err)
	}
	defer s.Fini()

	offset := 0 // 現在の表示開始行

	for {
		s.Clear()
		_, height := s.Size()

		// テキストの描画
		for i := 0; i < height; i++ {
			lineIdx := i + offset
			if lineIdx < len(lines) {
				for m := 0; m < len(configs); m++ {
					if strings.Index(lines[lineIdx], configs[m].LABEL) != -1 {
						drawText(s, 0, i, lines[lineIdx], true)
						configs[m].COUNT++
						if configs[m].COUNT >= configs[m].MAX {
							go speak(configs[m].VOICE)
							configs[m].COUNT = 0
						}
						break
					} else {
						drawText(s, 0, i, lines[lineIdx], false)
					}
				}
			}
		}

		s.Show()

		// イベントハンドリング
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Rune() == 'q' {
				return
			} else if ev.Key() == tcell.KeyDown || ev.Rune() == 'j' {
				if offset+height > len(lines) {
					offset = len(lines) - height
				} else {
					offset += height
				}
			} else if ev.Key() == tcell.KeyUp || ev.Rune() == 'k' {
				if offset-height < 0 {
					offset = 0
				} else {
					offset -= height
				}
			}
		case *tcell.EventResize:
			s.Sync()
		}
	}
}

// 指定した座標にテキストを描画するヘルパー関数
func drawText(s tcell.Screen, x, y int, str string, cFlag bool) {
	for i, r := range str {
		if cFlag == true {
			s.SetContent(x+i, y, r, nil, tcell.StyleDefault.Foreground(tcell.ColorRed))
		} else {
			s.SetContent(x+i, y, r, nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
		}
	}
}

// ファイルを全行読み込んでスライスに格納
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func speak(str string) {
	//var wg sync.WaitGroup
	//wg.Add(1)
	//go func() {
	ps.Execute(fmt.Sprintf("Add-Type -AssemblyName System.Speech; $voice = New-Object System.Speech.Synthesis.SpeechSynthesizer; $voice.Speak('%s'); end", str))
	//defer wg.Done()
	//}()
	//wg.Wait()
}

func loadConfig(configFile string) bool {
	var fp *os.File
	var err error
	fp, err = os.Open(configFile)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	reader := csv.NewReader(fp)
	reader.Comma = '\t'
	reader.LazyQuotes = true
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		if len(record) == 3 {
			i, err := strconv.Atoi(record[1])
			if err == nil {
				configs = append(configs, configData{LABEL: record[0], COUNT: 0, MAX: i, VOICE: record[2]})
				fmt.Println(record)
			}
		}
	}
	if configs == nil {
		return false
	}
	return true
}
