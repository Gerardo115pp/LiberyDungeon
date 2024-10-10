package dungeon_models

import "fmt"

type DungeonError struct {
	error_context string
	Err           error
}

func (de DungeonError) Error() string {
	return fmt.Sprintf("%s -> %s", de.error_context, de.Err.Error())
}

func (de *DungeonError) AppendContext(context_str string) {
	de.error_context = fmt.Sprintf("%s: %s", de.error_context, context_str)
}

type ErrorLabel string

const (
	ErrGenericError             ErrorLabel = "Generic error"       // Generic error label, try to avoid using this
	ErrProcessError             ErrorLabel = "Process error"       // Any error has no workarounds and completely impedes the process to continue.
	ErrPreconditionFailed       ErrorLabel = "Precondition failed" // An operation required a set of conditions that were not met.
	ErrIOError                  ErrorLabel = "I/O error"           // An error occurred while reading or writing to a file or network.
	ErrDB_CouldNotConnectToDB   ErrorLabel = "Could not connect to database"
	ErrDB_CouldNotCreateTX      ErrorLabel = "Could not create transaction"
	ErrDB_FailedToCommitTX      ErrorLabel = "Failed to commit transaction"
	ErrDB_FailedToRollbackTX    ErrorLabel = "Failed to rollback transaction"
	ErrLErrorVariableNotFound   ErrorLabel = "a variable in LabelError was not found"
	ErrLErrVariableTypeMismatch ErrorLabel = "a variable in LabelError was not of the expected type"
	ErrFS_NoSuchFileOrDirectory ErrorLabel = "No such file or directory"
	ErrFS_NoSuchDirectory       ErrorLabel = "No such directory"
	ErrFS_FileExists            ErrorLabel = "File already exists"
	ErrFS_PermissionDenied      ErrorLabel = "Permission denied"
	ErrPlatform_NoSuchCategory  ErrorLabel = "No such category"
	ErrPlatform_NoSuchCluster   ErrorLabel = "No such cluster"
	ErrPlatform_NoSuchMedia     ErrorLabel = "No media found"
)

type LabeledError struct {
	DungeonError
	Label             ErrorLabel
	relevantVariables map[string]interface{}
}

func (le LabeledError) Error() string {
	return fmt.Sprintf("+++%s+++\n%s", le.Label, le.DungeonError.Error())
}

func (le *LabeledError) StoreVariable(variable_name string, variable_value interface{}) {
	if le.relevantVariables == nil {
		le.relevantVariables = make(map[string]interface{})
	}

	le.relevantVariables[variable_name] = variable_value
}

func (le *LabeledError) GetVariable(variable_name string) (interface{}, *LabeledError) {
	if le.relevantVariables == nil {
		return nil, NewLabeledError(fmt.Errorf("Variable '%s' not found", variable_name), "In LabelError.GetVariable", ErrGenericError)
	}

	value, exists := le.relevantVariables[variable_name]
	if !exists {
		return nil, NewLabeledError(fmt.Errorf("Variable '%s' not found", variable_name), "In LabelError.GetVariable", ErrLErrorVariableNotFound)
	}

	return value, nil
}

/**
 * Returns the value of the variable as a string. If the variable is not found, or if it is not a string, an error is returned.
 */
func (le *LabeledError) GetStringVariable(variable_name string) (string, *LabeledError) {
	value, labeled_err := le.GetVariable(variable_name)
	if labeled_err != nil {
		return "", labeled_err
	}

	string_value, is_string := value.(string)
	if !is_string {
		return "", NewLabeledError(fmt.Errorf("Variable '%s' is not a string", variable_name), "In LabelError.GetStringVariable", ErrLErrVariableTypeMismatch)
	}

	return string_value, nil
}

func NewLabeledError(err error, context_string string, label ErrorLabel) *LabeledError {
	var new_error *LabeledError = new(LabeledError)

	new_error.error_context = context_string
	new_error.Err = err
	new_error.Label = label

	return new_error
}
