package iban

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	server "github.com/ymakhloufi/pfc/internal/http"
	"go.uber.org/zap"
)

var (
	_ server.Controller = Controller{}

	// faster to compile once, rather than each request.
	// Downside: if it fails, it panics the whole server on startup,
	// instead of doing regexp.Compile() and handling the returned error gracefully.
	validateEndpointRegexp = regexp.MustCompile(`^/v1/iban/([^/?]+)/validate/?$`)
)

// Parser can parse an iban string into an IBAN struct and validate its components.
type Parser interface {
	Parse(iban string) (IBAN, error)
	Validate(IBAN) error
}

// Controller the iban controller that adds routes to the http server.
type Controller struct {
	parser Parser
	logger *zap.Logger
}

func NewController(parser Parser, logger *zap.Logger) *Controller {
	return &Controller{
		parser: parser,
		logger: logger,
	}
}

// SetupRoutes adds the routes to the http server.
func (ctrl Controller) SetupRoutes() {
	// handles all routes prefixed with /iban/ (needed to handle non-query route-params)
	http.HandleFunc("/v1/iban/", func(w http.ResponseWriter, r *http.Request) {
		ctrl.logger.Info("request received", zap.String("method", r.Method), zap.String("path", r.URL.Path))

		// Add sub-routes as new "cases" here.
		switch path := r.URL.Path; {
		case r.Method == http.MethodGet && validateEndpointRegexp.MatchString(path): // /iban/<iban>/validate
			ctrl.validate(w, r)
			return
		default:
			ctrl.writeResponse(w, nil, fmt.Errorf("unsupported route: %s", path), http.StatusNotFound)
		}
	})
}

// swagger:operation GET /v1/iban/{iban}/validate validateIBAN
//
// # Validates a given IBAN string and returns its validity, components and a possible error message.
//
// ---
// parameters:
//   - in: path
//     name: iban
//     required: true
//     type: string
//
// responses:
//
//	'200':
//    description: IBAN was successfully validated, result can be positive or negative
//    schema:
//      $ref: '#/definitions/httpResponse'
//	'422':
//    description: IBAN string could not be parsed, i.e. has a wrong format (Ref https://en.wikipedia.org/wiki/International_Bank_Account_Number#Structure)
//    schema:
//      $ref: '#/definitions/httpResponse'
//	'500':
//	  description: Internal Server Error

// validate parses and validates the iban string.
func (ctrl Controller) validate(w http.ResponseWriter, r *http.Request) {
	ibanStr := validateEndpointRegexp.FindStringSubmatch(r.URL.Path)[1]
	ibanStr = strings.Replace(ibanStr, " ", "", -1) // remove all whitespaces to allow human-friendly spacing

	iban, err := ctrl.parser.Parse(ibanStr)
	if err != nil {
		ctrl.logger.Error("request failed", zap.Error(err))
		ctrl.writeResponse(w, nil, err, http.StatusUnprocessableEntity)
		return
	}

	err = ctrl.parser.Validate(iban)
	if err != nil {
		ctrl.writeResponse(w, &iban, err, http.StatusOK) // failed validation is an expected outcome, thus 200.
		return
	}

	ctrl.writeResponse(w, &iban, nil, http.StatusOK)
}

// writeResponse writes the response to the http response writer.
func (ctrl Controller) writeResponse(w http.ResponseWriter, iban *IBAN, err error, status int) {
	var errStr *string
	if err != nil {
		e := err.Error()
		errStr = &e
	}

	response := httpResponse{Error: errStr, IsValid: err == nil, IBAN: iban}
	l := ctrl.logger.With(
		zap.Any("iban", iban),
		zap.Error(err),
		zap.Int("status", status),
		zap.Any("response", response),
	)

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		l.Error("failed to marshal response", zap.Error(err))
		err = fmt.Errorf("failed to marshal response: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(jsonResponse)
	if err != nil {
		l.Error("failed to write response", zap.Error(err))
		err = fmt.Errorf("failed to write response: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
