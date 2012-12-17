package post

import (
    "time"
)

type Time struct {
    time.Time
}

func (t *Time) SetYAML(tag string, value interface{}) bool {
    s, ok := value.(string)
    if !ok {
        return false
    }
    parsed, err := time.Parse(time.RFC822, s)
    if err != nil {
        return false
    }
    t.Time = parsed
    return true
}
