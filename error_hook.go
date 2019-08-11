package errorhook

import (
	"io"
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

// WriteCallback is a callback that will be trigger when writing an entry
type WriteCallback = func(int, error)

// ErrorHook holds information to hook Logrus logs on error levels
type ErrorHook struct {
	KeepOutput   bool
	Writer       io.Writer
	LogLevels    []logrus.Level
	Formtter     logrus.Formatter
	OnAfterWrite WriteCallback
}

// Init ErrorHook for usage
func Init(
	keepOutput bool, writer io.Writer,
	levels []logrus.Level, onAfterWrite WriteCallback,
	formatter logrus.Formatter,
) *ErrorHook {
	if !keepOutput {
		logrus.SetOutput(ioutil.Discard)
	}
	if levels == nil {
		levels = []logrus.Level{
			logrus.ErrorLevel,
			logrus.FatalLevel,
		}
	}
	if len(levels) == 0 {
		levels = append(levels, logrus.ErrorLevel, logrus.FatalLevel)
	}
	return &ErrorHook{
		KeepOutput:   keepOutput,
		Writer:       writer,
		LogLevels:    levels,
		Formtter:     formatter,
		OnAfterWrite: onAfterWrite,
	}
}

// Level return the supported levels
func (h *ErrorHook) Level() []logrus.Level {
	return h.LogLevels
}

// Fire execute
func (h *ErrorHook) Fire(entry *logrus.Entry) error {
	formatted, err := h.Formtter.Format(entry)
	if err != nil {
		return err
	}
	n, err := h.Writer.Write(formatted)
	if h.OnAfterWrite != nil {
		h.OnAfterWrite(n, err)
	}
	if err != nil {
		return err
	}

	return nil
}
