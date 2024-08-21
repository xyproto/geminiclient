package simplegemini

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"cloud.google.com/go/vertexai/genai"
)

var ErrEmptyPrompt = errors.New("empty prompt")

// FunctionCallHandler defines a callback type for handling function responses.
type FunctionCallHandler func(response map[string]any) (map[string]any, error)

// AddFunctionTool registers a custom Go function as a tool that the model can call.
func (gc *GeminiClient) AddFunctionTool(name, description string, fn interface{}) error {
	fnValue := reflect.ValueOf(fn)
	fnType := fnValue.Type()

	if fnType.Kind() != reflect.Func {
		return fmt.Errorf("provided argument is not a function")
	}

	parameters := make(map[string]*genai.Schema)
	var required []string

	for i := 0; i < fnType.NumIn(); i++ {
		paramType := fnType.In(i)
		paramName := fmt.Sprintf("param%d", i+1)

		parameters[paramName] = &genai.Schema{
			Type: mapGoTypeToGenaiType(paramType),
		}
		required = append(required, paramName)
	}

	gc.Functions[name] = fnValue

	functionDecl := &genai.FunctionDeclaration{
		Name:        name,
		Description: description,
		Parameters: &genai.Schema{
			Type:       genai.TypeObject,
			Properties: parameters,
			Required:   required,
		},
	}

	tool := &genai.Tool{
		FunctionDeclarations: []*genai.FunctionDeclaration{functionDecl},
	}
	gc.Tools = append(gc.Tools, tool)

	return nil
}

// MultiQueryWithCallbacks processes a prompt, supports function tools, and uses a callback function to handle function responses.
func (gc *GeminiClient) MultiQueryWithCallbacks(prompt string, base64Data, dataMimeType *string, temperature *float32, callback FunctionCallHandler) (string, error) {
	if strings.TrimSpace(prompt) == "" {
		return "", ErrEmptyPrompt
	}

	gc.ClearParts()
	gc.AddText(prompt)

	if base64Data != nil && dataMimeType != nil {
		data, err := base64.StdEncoding.DecodeString(*base64Data)
		if err != nil {
			return "", fmt.Errorf("failed to decode base64 data: %v", err)
		}
		gc.AddData(*dataMimeType, data)
	}

	ctx, cancel := context.WithTimeout(context.Background(), gc.Timeout)
	defer cancel()

	model := gc.Client.GenerativeModel(gc.ModelName)
	if temperature != nil {
		model.SetTemperature(*temperature)
	}
	model.Tools = gc.Tools
	session := model.StartChat()

	res, err := session.SendMessage(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("failed to send message: %v", err)
	}

	for _, candidate := range res.Candidates {
		for _, part := range candidate.Content.Parts {
			if funcall, ok := part.(genai.FunctionCall); ok {
				responseData, err := gc.invokeFunction(funcall.Name, funcall.Args)
				if err != nil {
					return "", fmt.Errorf("failed to handle function call: %v", err)
				}

				if callback != nil {
					responseData, err = callback(responseData)
					if err != nil {
						return "", fmt.Errorf("callback processing failed: %v", err)
					}
				}

				res, err = session.SendMessage(ctx, genai.FunctionResponse{
					Name:     funcall.Name,
					Response: responseData,
				})
				if err != nil {
					return "", fmt.Errorf("failed to send function response: %v", err)
				}

				var finalResult strings.Builder
				for _, part := range res.Candidates[0].Content.Parts {
					if textPart, ok := part.(genai.Text); ok {
						finalResult.WriteString(string(textPart))
						finalResult.WriteString("\n")
					}
				}
				return strings.TrimSpace(finalResult.String()), nil
			}
		}
	}

	result, err := gc.SubmitToClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to process response: %v", err)
	}

	return strings.TrimSpace(result), nil
}

// MultiQueryWithSequentialCallbacks handles multiple function calls in sequence, using callback functions to manage responses.
func (gc *GeminiClient) MultiQueryWithSequentialCallbacks(prompt string, callbacks map[string]FunctionCallHandler) (string, error) {
	if strings.TrimSpace(prompt) == "" {
		return "", ErrEmptyPrompt
	}

	gc.ClearParts()
	gc.AddText(prompt)

	ctx, cancel := context.WithTimeout(context.Background(), gc.Timeout)
	defer cancel()

	model := gc.Client.GenerativeModel(gc.ModelName)
	model.Tools = gc.Tools
	session := model.StartChat()

	res, err := session.SendMessage(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("failed to send message: %v", err)
	}

	for _, candidate := range res.Candidates {
		for _, part := range candidate.Content.Parts {
			if funcall, ok := part.(genai.FunctionCall); ok {
				handler, exists := callbacks[funcall.Name]
				if !exists {
					return "", fmt.Errorf("no handler found for function: %s", funcall.Name)
				}

				responseData, err := handler(funcall.Args)
				if err != nil {
					return "", fmt.Errorf("handler error for function %s: %v", funcall.Name, err)
				}

				res, err = session.SendMessage(ctx, genai.FunctionResponse{
					Name:     funcall.Name,
					Response: responseData,
				})
				if err != nil {
					return "", fmt.Errorf("failed to send function response: %v", err)
				}
			}
		}
	}

	finalResult, err := gc.SubmitToClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to process final response: %v", err)
	}

	return strings.TrimSpace(finalResult), nil
}

// invokeFunction uses reflection to call the appropriate user-defined function based on the AI's request.
func (gc *GeminiClient) invokeFunction(name string, args map[string]any) (map[string]any, error) {
	fn, exists := gc.Functions[name]
	if !exists {
		return nil, fmt.Errorf("function %s not found", name)
	}

	fnType := fn.Type()

	in := make([]reflect.Value, fnType.NumIn())
	for i := 0; i < fnType.NumIn(); i++ {
		paramName := fmt.Sprintf("param%d", i+1)
		argValue, exists := args[paramName]
		if !exists {
			return nil, fmt.Errorf("missing argument: %s", paramName)
		}
		in[i] = reflect.ValueOf(argValue)
	}

	out := fn.Call(in)

	result := make(map[string]any)
	for i := 0; i < len(out); i++ {
		result[fmt.Sprintf("return%d", i+1)] = out[i].Interface()
	}

	return result, nil
}

// mapGoTypeToGenaiType maps Go types to the corresponding genai.Schema Type values.
func mapGoTypeToGenaiType(goType reflect.Type) genai.Type {
	switch goType.Kind() {
	case reflect.String:
		return genai.TypeString
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return genai.TypeInteger
	case reflect.Float32, reflect.Float64:
		return genai.TypeNumber
	case reflect.Bool:
		return genai.TypeBoolean
	default:
		return genai.TypeString
	}
}

// ClearToolsAndFunctions clears all registered tools and functions.
func (gc *GeminiClient) ClearToolsAndFunctions() {
	gc.Functions = make(map[string]reflect.Value)
	gc.Tools = []*genai.Tool{}
}
