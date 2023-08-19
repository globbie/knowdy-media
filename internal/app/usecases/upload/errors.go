package upload

var (
        ErrUnauthorized  = NewFileErr("UnAuthorized",  "unauthorized access")
        ErrAlreadyExists = NewFileErr("AlreadyExists", "file already exists")
        ErrRepoNotFound  = NewFileErr("NotFound",      "repo not found")
)

type FileError struct {
	code          string
	description   string
}

func NewFileErr(code string, description string) *FileError {
	return &FileError{
	        code:        code,
		description: description,
	}
}

func (e *FileError) Error() string {
	return e.description
}

func (e *FileError) Code() string {
	return e.code
}

