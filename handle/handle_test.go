package handle

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestJSONifyingError(t *testing.T) {
	response := jsonifyError(errors.New("Test"))

	if !reflect.DeepEqual(string(response), "{\"error\":\"Test\"}") {
		t.Fail()
	}
}

func TestResponseSetMessage(t *testing.T) {
	var response Response
	response.SetMessage("He sells seashells by the sea shore")

	if !reflect.DeepEqual(response.Message, "He sells seashells by the sea shore") {
		t.Fail()
	}
}

func TestCreatePost(t *testing.T) {
	r, _ := http.NewRequest("POST", "/create", nil)
	w := httptest.NewRecorder()
	Create(w, r)

	if !reflect.DeepEqual(http.StatusOK, w.Code) {
		t.Fail()
	}
}

func TestCreateRedirect(t *testing.T) {
	r, _ := http.NewRequest("GET", "/create", nil)
	w := httptest.NewRecorder()
	Create(w, r)

	if !reflect.DeepEqual(http.StatusMovedPermanently, w.Code) {
		t.Fail()
	}
}
