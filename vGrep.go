package main

import (
	"fmt"
	"sync"

	ps "github.com/Tobotobo/powershell"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	str := "abcdefg"
	go func() {
		ps.Execute(fmt.Sprintf("Add-Type -AssemblyName System.Speech; $voice = New-Object System.Speech.Synthesis.SpeechSynthesizer; $voice.Speak('%s'); end", str))
		defer wg.Done()
	}()
	wg.Wait()
}
