package responses

import (
    "errors"
    "encoding/json"
    "testing"
    "net/http"
    "net/http/httptest"

    gm "github.com/onsi/gomega"
)


func TestResponseDataWrite(t *testing.T) {
    g := gm.NewGomegaWithT(t)
    rr := httptest.NewRecorder()

    expectedBody := `{"foo":5}`

    bodyBytes, err := json.Marshal(map[string]int{
        "foo": 5,
    })
    g.Expect(err).To(gm.BeNil())

    responseData := Data{
        Headers: map[string][]string{
            "Test-Header": []string{"Bar"},
        },
        Body: bodyBytes,
        StatusCode: http.StatusOK,
    }

    err = responseData.Write(rr)
    g.Expect(err).To(gm.BeNil())

    g.Expect(rr.Body.String()).To(gm.Equal(expectedBody))
    g.Expect(rr.Header().Get("Test-Header")).To(gm.Equal("Bar"))
}

func TestErrorResponse(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    expectedError := errors.New("Unexpected error for test")
    responseData := ErrorResponse(expectedError)

    g.Expect(responseData.Error()).To(gm.MatchError(expectedError))
}
