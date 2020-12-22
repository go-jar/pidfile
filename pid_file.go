package pidfile

import (
	"errors"
	"os"
	"strconv"
)

const (
	MinPid        = 0
	TmpPathSuffix = ".tmp"
)

type PidFile struct {
	Pid     *Pid
	File    *File
	TmpFile *File
}

func NewPidFile(path string) *PidFile {
	return &PidFile{
		Pid:     NewPid(os.Getpid()),
		File:    NewFile(path),
		TmpFile: NewFile(path + TmpPathSuffix),
	}
}

func (pf *PidFile) Create() error {
	file := pf.File
	pid, err := ReadFromFile(file)
	if err == nil && pid.ProcessExits() {
		file = pf.TmpFile
	}

	if err = WritePidToFile(file, pf.Pid); err != nil {
		return err
	}

	return nil
}

func (pf *PidFile) Clear() error {
	pid, err := ReadFromFile(pf.File)
	tmpPid, tmpErr := ReadFromFile(pf.TmpFile)

	if err != nil && tmpErr != nil {
		return errors.New("clear pid error: " + err.Error() + ", clear tmp pid error: " + tmpErr.Error())
	}

	if err == nil && pf.Pid.Id == pid.Id {
		if err = pf.File.Remove(); err != nil {
			return err
		}

		if tmpErr == nil {
			if err = pf.TmpFile.Rename(pf.File.Path); err != nil {
				return err
			}
		}
	} else {
		if tmpErr == nil && pf.Pid.Id == tmpPid.Id {
			if err = pf.TmpFile.Remove(); err != nil {
				return err
			}
		}
	}

	return nil
}

func ReadFromFile(file *File) (*Pid, error) {
	fb, err := file.Read()
	if err != nil {
		return nil, err
	}

	id, err := strconv.Atoi(string(fb))
	if err != nil || id <= MinPid {
		return nil, errors.New("pid file data error")
	}

	return NewPid(id), nil
}

func WritePidToFile(file *File, pid *Pid) error {
	fb := []byte(strconv.Itoa(pid.Id))
	return file.Write(fb)
}

func CreatePidFile(path string) (*PidFile, error) {
	pf := NewPidFile(path)
	err := pf.Create()
	return pf, err
}

func ClearPidFile(pf *PidFile) error {
	return pf.Clear()
}
