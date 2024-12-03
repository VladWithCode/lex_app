package fetchers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"time"
)

const DGO_URLF = "http://tsjdgo.gob.mx/Recursos/images/flash/ListasAcuerdos/%v/%v.pdf"
const LARGE_DATA_SIZE = 60_000

func DgoFetch(date time.Time, caseType string) (data *[]byte, err error) {
	data, err = dgoFetchResource(date, caseType)

	if err != nil {
		return nil, fmt.Errorf("Fetch file err: %w", err)
	}

	err = PDFToData(data)

	if err != nil {
		return nil, fmt.Errorf("Transform file err: %w", err)
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

	data = new([]byte)
	*data, err = io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func PDFToData(data *[]byte) error {
	outBuf := new(bytes.Buffer)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	transformCmd := exec.CommandContext(ctx, "pdftotext", "-fixed", "5", "-", "-")

	stdIn, err := transformCmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("Error getting stdin pipe:\n\t%w", err)
	}

	stdOut, err := transformCmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("Error getting stdout pipe:\n\t%w", err)
	}

	if err := transformCmd.Start(); err != nil {
		return fmt.Errorf("Start cmd error: %w", err)
	}

	// Concurrenty writes to cmd in pipe
	go func() {
		defer stdIn.Close()
		stdIn.Write(*data)
	}()

	// Councurrently reads from cmd out pipe
	go func() {
		defer stdOut.Close()
		outBuf.ReadFrom(stdOut)
	}()

	// Await read, write and execution of the command
	if err := transformCmd.Wait(); err != nil {
		return fmt.Errorf("Error executing pdftotext:\n  %w", err)
	}

	*data = outBuf.Bytes()

	return nil
}

// Keeping this for reference. It may not be usefull at all in the future
func _pdfToData(data *[]byte) error {
	if len(*data) >= LARGE_DATA_SIZE {
		fmt.Printf("PDF of size %v. Using largePDFToData handler\n", len(*data))
		return PDFToData(data)
	}

	transformCmd := exec.Command("pdftotext", "-fixed", "5", "-", "-")

	stdIn, err := transformCmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("Error getting stdin pipe:\n\t%w", err)
	}
	stdOut, err := transformCmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("Error getting stdout pipe:\n\t%w", err)
	}

	if err := transformCmd.Start(); err != nil {
		return fmt.Errorf("Error starting cmd:\n\t%w", err)
	}

	out := make([]byte, len(*data))
	go func() {
		defer stdIn.Close()
		stdIn.Write(*data)
	}()
	go func() {
		stdOut.Read(out)
	}()

	if err := transformCmd.Wait(); err != nil {
		return fmt.Errorf("Error executing cmd:\n\t%w", err)
	}

	*data = out

	return nil
}
