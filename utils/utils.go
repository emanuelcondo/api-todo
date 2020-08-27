package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-chi/render"
)

type Body struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
	Error  interface{} `json:"error"`
}

type DataType int
type Source int

const (
	DATATYPE_STRING		DataType = 1
	DATATYPE_INTEGER	DataType = 2
	DATATYPE_BOOLEAN	DataType = 3
	DATATYPE_FLOAT		DataType = 4

	SOURCE_PARAMS		Source = 1
	SOURCE_HEADERS		Source = 2
	SOURCE_QUERY		Source = 3
	SOURCE_BODY			Source = 4
)

type InputOptions struct {
	MinLength		uint
	MaxLength		uint
	MinValue		uint
	MaxValue		uint
	AllowedValues	[]interface{}

	ParseToInt		bool
	ParseToFloat	bool
	ParseToBool		bool
	ParseToString	bool
}

type validationError struct {
	message 	string
}

// reply
func Reply(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	body := Body{}

	if status >= 200 && status < 400 {
		body.Status = "success"
		body.Data = data
	} else {
		body.Status = "error"
		body.Error = data
	}

	render.Status(r, status)
	render.JSON(w, r, body)
}

func InputValidator(key string, datatype DataType, source Source, mandatory bool, options ...InputOptions) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		var op *InputOptions = nil
		if len(options) > 0 {
			op = &options[0]
		}

		fn := func(w http.ResponseWriter, r *http.Request) {
			if key == "" || strings.TrimSpace(key) == "" {
				resp := map[string]string {
					"message": "An error occurred when validating input",
				}
				Reply(w, r, http.StatusInternalServerError, resp)
				return
			}

			var err *validationError = nil
			switch source {
				case SOURCE_PARAMS:
					err = validateInputFromParams(key, datatype, mandatory, r, op)
					break
				case SOURCE_HEADERS:
					err = validateInputFromHeaders(key, datatype, mandatory, r, op)
					break
				case SOURCE_QUERY:
					err = validateInputFromQuery(key, datatype, mandatory, r, op)
					break
				case SOURCE_BODY:
					err = validateInputFromBody(key, datatype, mandatory, r, op)
					break
				default:
					break
			}
			if err != nil {
				resp := map[string]string {
					"message": err.message,
				}
				Reply(w, r, http.StatusUnprocessableEntity, resp)
				return
			} else {
				next.ServeHTTP(w, r)
			}
		}
		return http.HandlerFunc(fn)
	}
}

func validateInputFromParams(key string, dataType DataType, mandatory bool, r *http.Request, op *InputOptions) *validationError {
	value := chi.URLParam(r, key)
	err := validateValue(key, value, dataType, mandatory, op)
	return err
}

func validateInputFromHeaders(key string, dataType DataType, mandatory bool, r *http.Request, op *InputOptions) *validationError {
	value := r.Header.Get(key)
	return validateValue(key, value, dataType, mandatory, op)
}

func validateInputFromQuery(key string, dataType DataType, mandatory bool, r *http.Request, op *InputOptions) *validationError {
	qs := r.URL.Query()
	values, _ := qs[key]
	if values != nil {
		var value interface{} = nil
		if len(values) > 0 {
			value = values[0]
		}
		return validateValue(key, value, dataType, mandatory, op)
	} else {
		return nil
	}
}

func validateInputFromBody(key string, dataType DataType, mandatory bool, r *http.Request, op *InputOptions) *validationError {
	var body interface{}

	bodyBytes, _ := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	err := json.Unmarshal(bodyBytes, &body)

	//err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		return &validationError{ message: "Body couldn't be parsed. Check that body be a valid json." }
	} else {
		value := body.(map[string]interface{})[key]
		return validateValue(key, value, dataType, mandatory, op)
	}
}

func validateValue(key string, value interface{}, dataType DataType, mandatory bool, op *InputOptions) *validationError {
	if mandatory && (value == nil || value == "") {
		return &validationError{ message: "Field '" + key + "' is mandatory." }
	}

	var err *validationError = nil
	isInputString := reflect.TypeOf(value).String() == reflect.String.String()

	if op != nil {
		if op.ParseToInt && isInputString {
			param := fmt.Sprintf("%s", value)
			parsed, err := strconv.ParseInt(param, 10, 32)
			if err == nil {
				value = int(parsed)
			} else {
				return &validationError{ message: "Field '"+key+"' could not be parsed to 'integer'" }
			}
		} else if op.ParseToFloat && isInputString {
			param := fmt.Sprintf("%s", value)
			parsed, err := strconv.ParseFloat(param, 32)
			if err == nil {
				value = float32(parsed)
			} else {
				return &validationError{ message: "Field '"+key+"' could not be parsed to 'float'" }
			}
		} else if op.ParseToBool && isInputString {
			param := fmt.Sprintf("%v", value)
			if param == "true" || param == "false" {
				if param == "true" {
					value = true
				} else {
					value = false
				}
			}
		}
	}

	reflectType := reflect.TypeOf(value).String()

	switch dataType {
		case DATATYPE_STRING:
			if reflectType != reflect.String.String() {
				err = &validationError{ message: "Field '" + key + "' contains an invalid value. It should be of type 'string'." }
			} else if op != nil {
				param := fmt.Sprintf("%s", value)
				if op.MinLength > uint(len(param)) {
					err = &validationError{ message: "Field '" + key + "' contains an invalid length. Minimum allowed length: "+strconv.Itoa(int(op.MinLength))+"." }
				}
				if op.MinLength < op.MaxLength && op.MaxLength < uint(len(param)) {
					err = &validationError{ message: "Field '" + key + "' contains an invalid length. Maximum allowed length: "+strconv.Itoa(int(op.MaxLength))+"." }
				}
			}
			break
		case DATATYPE_INTEGER:
			if reflectType != reflect.Int.String() {
				err = &validationError{ message: "Field '" + key + "' contains an invalid value. It should be of type 'integer'." }
			} else if op != nil {
				if int(op.MinValue) > value.(int) {
					err = &validationError{ message: "Field '" + key + "' contains an invalid length. Minimum allowed value: "+strconv.Itoa(int(op.MinValue))+"." }
				}
				if op.MinValue < op.MaxValue && int(op.MaxValue) < value.(int) {
					err = &validationError{ message: "Field '" + key + "' contains an invalid length. Maximum allowed value: "+strconv.Itoa(int(op.MaxValue))+"." }
				}
			}
			break
		case DATATYPE_BOOLEAN:
			if reflectType != reflect.Bool.String() {
				err = &validationError{ message: "Field '" + key + "' contains an invalid value. It should be of type 'boolean'." }
			}
			break
		case DATATYPE_FLOAT:
			if reflectType != reflect.Float32.String() {
				err = &validationError{ message: "Field '" + key + "' contains an invalid value. It should be of type 'float'." }
			} else if op != nil {
				if float32(op.MinValue) > value.(float32) {
					err = &validationError{ message: "Field '" + key + "' contains an invalid length. Minimum allowed value: "+strconv.Itoa(int(op.MinValue))+"." }
				}
				if float32(op.MinValue) < float32(op.MaxValue) && float32(op.MaxValue) < value.(float32) {
					err = &validationError{ message: "Field '" + key + "' contains an invalid length. Maximum allowed value: "+strconv.Itoa(int(op.MaxValue))+"." }
				}
			}
			break
		default:
			break
	}

	return err
}