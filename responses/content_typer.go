package responses

import (
)

type ContentTyper interface {
    ContentType(string) string
}
