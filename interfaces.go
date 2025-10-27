package netinterface

import (
	"io"
	"os"
)

const DefaultInterfacesPath = "/etc/network/interfaces"

func Read(r io.Reader) (*File, error) {
	return Parse(r)
}

func ReadDefault() (*File, error) {
	return ReadFile(DefaultInterfacesPath)
}

func ReadFile(path string) (*File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Parse(f)
}

func Write(w io.Writer, f *File) error {
	if err := ValidateAndReport(f); err != nil {
		return err
	}
	return GenerateTo(w, f)
}

func WriteFile(path string, f *File) error {
	w, err := os.Create(path)
	if err != nil {
		return err
	}
	defer w.Close()
	return Write(w, f)
}

func WriteDefault(f *File) error {
	return WriteFile(DefaultInterfacesPath, f)
}

