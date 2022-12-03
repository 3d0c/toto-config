package helpers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/3d0c/toto-config/pkg/log"
)

// JSONResponder Responder. TODO Make a Responder interface
type JSONResponder struct {
	w http.ResponseWriter
}

func NewJSONResponder(w http.ResponseWriter) JSONResponder {
	return JSONResponder{w: w}
}

func (j JSONResponder) Encode(v interface{}) ([]byte, error) {
	var (
		b   []byte
		err error
	)

	if b, err = json.MarshalIndent(v, "", "    "); err != nil {
		return nil, err
	}

	return b, nil
}

func (j JSONResponder) Write(v interface{}) {
	var (
		b   []byte
		err error
	)

	if b, err = j.Encode(v); err != nil {
		log.TheLogger().Error("error encoding response", zap.Error(err))
		http.Error(j.w, "", http.StatusInternalServerError)
		return
	}

	if _, err = j.w.Write(b); err != nil {
		log.TheLogger().Error("error writing response", zap.Error(err))
	}
}
