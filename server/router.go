package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/brxie/eluborzyca-backend/config"
	"github.com/brxie/eluborzyca-backend/utils"
	"github.com/brxie/eluborzyca-backend/utils/ilog"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
)

// SwaggerRouter creates Swagger Router
func SwaggerRouter(swaggerFile string) http.Handler {
	loader := *openapi3.NewSwaggerLoader()
	swagger, err := loader.LoadSwaggerFromFile(swaggerFile)
	if err != nil {
		ilog.Panic("Unable to load swagger file.")
		panic(err)
	}

	oa3router := openapi3filter.NewRouter().WithSwagger(swagger)

	openapi3filter.RegisterBodyDecoder("image/jpeg", openapi3filter.FileBodyDecoder)
	openapi3filter.RegisterBodyDecoder("image/png", openapi3filter.FileBodyDecoder)

	httpHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Find route
		route, pathParams, _ := oa3router.FindRoute(r.Method, r.URL)
		ilog.Debug(fmt.Sprintf("Request '%s' parameters: %v", r.URL, pathParams))
		if route == nil {
			utils.WriteMessageResponse(&w, http.StatusNotFound,
				http.StatusText(http.StatusNotFound))
			return
		}
		ilog.Info(fmt.Sprintf("New request %s %s", route.Method, route.Path))
		ctx := context.WithValue(context.TODO(), config.PARAMS, pathParams)

		// Validate
		requestValidationInput := &openapi3filter.RequestValidationInput{
			Request:    r,
			PathParams: pathParams,
			Route:      route,
			Options: &openapi3filter.Options{
				AuthenticationFunc: extractSession,
			},
		}
		if err := validate(ctx, requestValidationInput); err != nil {
			handleValidationError(err, w)
			return
		}

		// Route request
		if err := startHandler(route, w, r.WithContext(ctx)); err != nil {
			utils.WriteMessageResponse(&w, http.StatusInternalServerError, err.Error())
			return
		}
	})

	return httpHandler
}

// Validates request against swagger endpoint schema
func validate(context context.Context, validationInput *openapi3filter.RequestValidationInput) error {
	// Validate request
	if err := openapi3filter.ValidateRequest(context, validationInput); err != nil {
		ilog.Info("Request validation failed: ", err)
		return err
	}
	// Validate body
	if requestBodySchema := validationInput.Route.Operation.RequestBody; requestBodySchema != nil {
		requestBodySchema := requestBodySchema.Value.WithRequired(true)
		if err := openapi3filter.ValidateRequestBody(context, validationInput, requestBodySchema); err != nil {
			ilog.Info("Request body validation failed: ", err)
			return err
		}
	}
	return nil
}

// Gets validation error and sends a proper response to the client
// based on the kind of problem that occurred.
func handleValidationError(validErr error, w http.ResponseWriter) {
	switch validErr.(type) {
	case *openapi3filter.SecurityRequirementsError:
		sessionError := validErr.(*openapi3filter.SecurityRequirementsError).Errors[0]
		unauthResp(w, sessionError)
		return
	default:
		utils.WriteMessageResponse(&w, http.StatusBadRequest,
			"Request didn't match expected schema. "+validErr.Error())
	}
}

func unauthResp(w http.ResponseWriter, authErr error) {
	switch authErr.(type) {
	case *SessionError:
		utils.WriteMessageResponse(&w, authErr.(*SessionError).RespCode, authErr.(*SessionError).Error())
		return
	}
	utils.WriteMessageResponse(&w, http.StatusInternalServerError,
		"Unexpected error: "+authErr.Error())
}

// Starts handler based on provided 'operationId' value in swagger definition.
// Handler is starting with authorization middleware if needed.
func startHandler(route *openapi3filter.Route, w http.ResponseWriter, r *http.Request) error {
	operationID := route.Operation.OperationID
	if _, ok := Handlers[operationID]; !ok {
		return errors.New("Operation '" + operationID + "' not defined in the handler map.")
	}

	Handlers[operationID](w, r)
	return nil
}
