package fetchers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"
)

const LARGE_DATA_SIZE = 65_000
const DGO_URLF = "http://tsjdgo.gob.mx/Recursos/images/flash/ListasAcuerdos/%v/%v.pdf"

func DgoFetch(date time.Time, caseType string) (data *[]byte, err error) {
	data, err = dgoFetchResource(date, caseType)

	if err != nil {
		return nil, err
	}

	err = dataToPDF(data)

	if err != nil {
		return nil, err
	}

	return
}

func dgoFetchResource(date time.Time, caseType string) (data *[]byte, err error) {
	formattedDate := date.Format("212006")
	resourceUrl := fmt.Sprintf(DGO_URLF, formattedDate, caseType)

	response, err := http.Get(resourceUrl)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 400 {
		return nil, errors.New("No se encontró documento para la fecha solicitada")
	}

	*data, err = io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func dataToPDF(data *[]byte) error {
	if len(*data) >= LARGE_DATA_SIZE {
		return largeDataToPDF(data)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	transformCmd := exec.CommandContext(ctx, "pdftotext", "-fixed 5", "-", "-")
	inPipe, err := transformCmd.StdinPipe()
	defer inPipe.Close()

	if err != nil {
		return err
	}

	_, err = inPipe.Write(*data)

	if err != nil {
		return err
	}

	*data, err = transformCmd.Output()

	return nil
}

func largeDataToPDF(data *[]byte) error {
	inRead, inWrite, _ := os.Pipe()
	outRead, outWrite, _ := os.Pipe()
	outBuf := new(bytes.Buffer)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	transformCmd := exec.CommandContext(ctx, "pdftotext", "-fixed 5", "-", "-")
	transformCmd.Stdin = inRead
	transformCmd.Stdout = outWrite

	if err := transformCmd.Start(); err != nil {
		return err
	}

	// Concurrenty writes to cmd in pipe
	go func() {
		defer inWrite.Close()

		inWrite.Write(*data)
	}()

	// Councurrently reads from cmd out pipe
	go func() {
		defer outRead.Close()

		outBuf.ReadFrom(outRead)
	}()

	// Await read, write and execution of the command
	if err := transformCmd.Wait(); err != nil {
		return nil
	}

	*data = outBuf.Bytes()

	return nil
}
