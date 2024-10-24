package constants

const (
	// General errors
	ErrMethodNotAllowed   = "Method not allowed"
	ErrInvalidPayload     = "Invalid request payload"
	ErrFailedToGetUserID  = "Failed to get userID"
	ErrInvalidCredentials = "Invalid credentials"

	// Snippet-related errors
	ErrFailedToGetSnippets     = "Failed to get snippets"
	ErrFailedToCreateSnippet   = "Failed to create snippet"
	ErrFailedToUpdateSnippet   = "Failed to update snippet"
	ErrSnippetNotFound         = "Snippet not found"
	ErrFailedToDeleteSnippet   = "Failed to delete snippet"
	ErrInvalidSnippetID        = "Invalid snippet ID"
	ErrSnippetValidationFailed = "Invalid request payload: "

	// User-related errors
	ErrInvalidEmailFormat     = "Invalid email format"
	ErrFailedToCreateUser     = "Failed to create user"
	ErrFailedToGenerateToken  = "Failed to generate JWT token"
	ErrFailedToDeleteUser     = "Failed to delete user"
	ErrFailedToExtractTokenID = "Failed to extract ID from token"
	ErrFailedToUpdatePassword = "Failed to updated password"

	// Validation messages
	ErrEmptyTitle         = "Title can't be empty!"
	ErrEmptyLanguage      = "Language can't be empty!"
	ErrEmptyContent       = "Content can't be empty!"
	ErrPasswordTooShort   = "Password must be at least 8 characters long"
	ErrPasswordTooLong    = "Password must be at most 20 characters long"
	ErrInvalidPassword    = "Password must contain at least one %s"
	ErrEmptyUsername      = "Username can't be empty"
	ErrUsernameTooShort   = "Username must be at least 3 characters long"
	ErrUsernameTooLong    = "Username must be at most 30 characters long"
	ErrEmptyEmail         = "Email can't be empty"
	ErrEmptyPassword      = "Password can't be empty"
	ErrInvalidSortOptions = "Sort options are invalid"
)
