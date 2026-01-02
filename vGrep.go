package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"

	ps "github.com/Tobotobo/powershell"
	"github.com/gdamore/tcell/v2"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("使用法: go run main.go <ファイル名>")
		return
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
				drawText(s, 0, i, lines[lineIdx])
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
				if offset < len(lines)-height {
					offset++
				}
			} else if ev.Key() == tcell.KeyUp || ev.Rune() == 'k' {
				if offset > 0 {
					offset--
				}
			}
		case *tcell.EventResize:
			s.Sync()
		}
	}
}

// 指定した座標にテキストを描画するヘルパー関数
func drawText(s tcell.Screen, x, y int, str string) {
	for i, r := range str {
		s.SetContent(x+i, y, r, nil, tcell.StyleDefault)
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

func speak() {
	var wg sync.WaitGroup
	wg.Add(1)
	str := "abcdefg"
	go func() {
		ps.Execute(fmt.Sprintf("Add-Type -AssemblyName System.Speech; $voice = New-Object System.Speech.Synthesis.SpeechSynthesizer; $voice.Speak('%s'); end", str))
		defer wg.Done()
	}()
	wg.Wait()
}
