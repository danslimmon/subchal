package subchal

import (
    "log"
    "testing"
)

type TestLogger struct {
    T *testing.T
}
func (logger TestLogger) Write(p []byte) (n int, err error) {
    logger.T.Log(string(p))
    return len(p), nil
}

func SetTestLogger(t *testing.T) {
    tl := TestLogger{
        T: t,
    }
    log.SetOutput(tl)
}
