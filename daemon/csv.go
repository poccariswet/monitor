package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"
)

const (
	csv_path = "./tmp/error.csv"
)

func WriteServerDown() error {
	file, err := os.OpenFile(csv_path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer file.Close()

	w := csv.NewWriter(file)
	errcode := "server process down"

	y, m, d := time.Now().Date()
	date := fmt.Sprintf("%d/%s/%d %d:%d", y, m, d, time.Now().Hour(), time.Now().Minute())

	if err := w.Write([]string{date, errcode}); err != nil {
		return err
	}

	w.Flush()

	if err := w.Error(); err != nil {
		return err
	}
	return nil
}
