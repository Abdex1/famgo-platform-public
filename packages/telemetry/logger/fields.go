 // STEP 5 — LOGGER FIELDS

//packages/telemetry/go/logger/fields.go

package logger

import "github.com/rs/zerolog"

type Field struct {
    Key   string
    Value any
}

func (f Field) Apply(e *zerolog.Event) {
    e.Interface(f.Key, f.Value)
}

func String(key, value string) Field {
    return Field{Key: key, Value: value}
}

func Int(key string, value int) Field {
    return Field{Key: key, Value: value}
}

func Bool(key string, value bool) Field {
    return Field{Key: key, Value: value}
}
