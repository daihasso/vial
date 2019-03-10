package responses

import (
)

import (
    "context"
    "errors"
    "testing"
    "net/http"

    gm "github.com/onsi/gomega"

    "github.com/daihasso/vial/neterr"
)

func TestBuilderBasic(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    expectedString := "foo"
    expectedBodyBytes := []byte(`foo`)

    builder, err := NewBuilder(
        context.Background(),
        JSONEncoding,
        Body(expectedString),
    )
    g.Expect(err).To(gm.BeNil())

    responseData := builder.Finish()

    g.Expect(responseData.Body).To(gm.Equal(expectedBodyBytes))
}

func TestBuilderListBody(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    expectedBody1, expectedBody2 := "foo", 1
    expectedBodyBytes := []byte(`["foo",1]`)

    builder, err := NewBuilder(
        context.Background(),
        JSONEncoding,
        Body(expectedBody1, expectedBody2),
    )
    g.Expect(err).To(gm.BeNil())

    responseData := builder.Finish()

    t.Log(string(responseData.Body))
    g.Expect(responseData.Body).To(gm.Equal(expectedBodyBytes))
}

func TestNewBuilderFailureInAdditional(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    testErr := errors.New("Expected failure!")
    builder, err := NewBuilder(
        context.Background(),
        JSONEncoding,
        func(rb *Builder) error {
            return testErr
        },
    )
    g.Expect(builder).To(gm.BeNil())
    g.Expect(err).ToNot(gm.BeNil())
    t.Log(err.Error())
    g.Expect(err.Error()).To(gm.Equal(
        "Error while running addditionals on Builder: Expected failure!",
    ))
}

func TestBuilderFailureInAdditionalFinish(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    testErr := errors.New("Expected failure!")
    builder, err := NewBuilder(
        context.Background(),
        JSONEncoding,
    )
    g.Expect(err).To(gm.BeNil())

    responseData := builder.Finish(
        func(rb *Builder) error {
            return testErr
        },
    )
    t.Log(responseData.Error().Error())
    g.Expect(responseData.Error().Error()).To(gm.Equal(
        "Error while adding finishing addditionals to Builder: Error while" +
            " running addditionals on Builder: Expected failure!",
    ))
}

func TestBuilderReplaceHeaders(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    builder, err := NewBuilder(
        context.Background(),
        JSONEncoding,
        Headers(map[string][]string{
            "Test-Header": []string{"Foo"},
        }),
    )
    g.Expect(err).To(gm.BeNil())

    builder.ReplaceHeaders(map[string][]string{
        "Better-Header": []string{"Bar"},
    })

    responseData := builder.Finish()
    t.Log(responseData)
    _, ok := responseData.Headers["Test-Header"]
    g.Expect(ok).ToNot(gm.BeTrue())
    betterHeader, ok := responseData.Headers["Better-Header"]
    g.Expect(ok).To(gm.BeTrue())
    g.Expect(betterHeader).To(gm.ConsistOf("Bar"))
}

func TestBuilderAddHeadersNew(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    builder, err := NewBuilder(
        context.Background(),
        JSONEncoding,
        Headers(map[string][]string{
            "Test-Header": []string{"Foo"},
        }),
    )
    g.Expect(err).To(gm.BeNil())

    builder.AddHeaders(map[string][]string{
        "Better-Header": []string{"Bar"},
    })

    responseData := builder.Finish()
    t.Log(responseData)
    testHeader, ok := responseData.Headers["Test-Header"]
    g.Expect(ok).To(gm.BeTrue())
    g.Expect(testHeader).To(gm.ConsistOf("Foo"))
    betterHeader, ok := responseData.Headers["Better-Header"]
    g.Expect(ok).To(gm.BeTrue())
    g.Expect(betterHeader).To(gm.ConsistOf("Bar"))
}

func TestBuilderAddHeadersExisting(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    builder, err := NewBuilder(
        context.Background(),
        JSONEncoding,
        Headers(map[string][]string{
            "Test-Header": []string{"Foo"},
        }),
    )
    g.Expect(err).To(gm.BeNil())

    builder.AddHeaders(map[string][]string{
        "Test-Header": []string{"Bar"},
    })

    responseData := builder.Finish()
    t.Log(responseData)
    testHeader, ok := responseData.Headers["Test-Header"]
    g.Expect(ok).To(gm.BeTrue())
    g.Expect(testHeader).To(gm.ConsistOf("Foo", "Bar"))
}

func TestBuilderSetHeader(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    builder, err := NewBuilder(
        context.Background(),
        JSONEncoding,
        Headers(map[string][]string{
            "Test-Header": []string{"Foo"},
        }),
    )
    g.Expect(err).To(gm.BeNil())

    builder.SetHeader("Test-Header", "Bar")

    responseData := builder.Finish()
    t.Log(responseData)
    testHeader, ok := responseData.Headers["Test-Header"]
    g.Expect(ok).To(gm.BeTrue())
    g.Expect(testHeader).To(gm.ConsistOf("Bar"))
}

func TestBuilderAddHeaderNew(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    builder, err := NewBuilder(
        context.Background(),
        JSONEncoding,
        Headers(map[string][]string{
            "Test-Header": []string{"Foo"},
        }),
    )
    g.Expect(err).To(gm.BeNil())

    builder.AddHeader("Better-Header", "Bar")

    responseData := builder.Finish()
    t.Log(responseData)
    testHeader, ok := responseData.Headers["Test-Header"]
    g.Expect(ok).To(gm.BeTrue())
    g.Expect(testHeader).To(gm.ConsistOf("Foo"))
    betterHeader, ok := responseData.Headers["Better-Header"]
    g.Expect(ok).To(gm.BeTrue())
    g.Expect(betterHeader).To(gm.ConsistOf("Bar"))
}

func TestBuilderAddHeaderExisting(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    builder, err := NewBuilder(
        context.Background(),
        JSONEncoding,
        Headers(map[string][]string{
            "Test-Header": []string{"Foo"},
        }),
    )
    g.Expect(err).To(gm.BeNil())

    builder.AddHeader("Test-Header", "Bar")

    responseData := builder.Finish()
    t.Log(responseData)
    testHeader, ok := responseData.Headers["Test-Header"]
    g.Expect(ok).To(gm.BeTrue())
    g.Expect(testHeader).To(gm.ConsistOf("Foo", "Bar"))
}

func TestBuilderAdditionals(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    expectedBodyBytes := []byte(`{"baz":1,"foo":"bar"}`)

    builder, err := NewBuilder(
        context.Background(),
        JSONEncoding,
        Headers(map[string][]string{
            "Test-Header": []string{"Foo"},
        }),
        Body(map[string]interface{}{
            "foo": "bar",
            "baz": 1,
        }),
        Status(http.StatusOK),
    )
    g.Expect(err).To(gm.BeNil())

    g.Expect(builder.encodingType).To(gm.Equal(JSONEncoding))

    responseData := builder.Finish(
        AddHeader("Second-header", "FooBar"),
    )
    t.Log(responseData)

    testHeader, ok := responseData.Headers["Test-Header"]
    g.Expect(ok).To(gm.BeTrue())
    g.Expect(testHeader).To(gm.ConsistOf("Foo"))

    secondHeader, ok := responseData.Headers["Second-Header"]
    g.Expect(ok).To(gm.BeTrue())
    g.Expect(secondHeader).To(gm.ConsistOf("FooBar"))

    g.Expect(responseData.StatusCode).To(gm.Equal(http.StatusOK))

    g.Expect(responseData.Body).To(gm.Equal(expectedBodyBytes))
}

func TestBuilderCustomContentType(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    builder, err := NewBuilder(
        context.Background(),
        JSONEncoding,
    )
    g.Expect(err).To(gm.BeNil())

    customContentType := "application/vnd.my.company+json"
    builder.SetContentType(customContentType)

    responseData := builder.Finish()
    t.Log(responseData)

    contentTypeHeader, ok := responseData.Headers["Content-Type"]
    g.Expect(ok).To(gm.BeTrue())
    g.Expect(contentTypeHeader).To(gm.ConsistOf(customContentType))
}

func TestBuilderAbort(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    expectedBodyBytes := []byte(
        `{"errors":[{"code":0,"message":"Expected failure while testing"}]}`,
    )

    builder, err := NewBuilder(
        context.Background(),
        JSONEncoding,
        Headers(map[string][]string{
            "Test-Header": []string{"Foo"},
        }),
    )
    g.Expect(err).To(gm.BeNil())

    responseData := builder.Abort(
        http.StatusConflict,
        neterr.NewCodedError(
            0, "Expected failure while testing",
        ),
    )
    t.Log(responseData)

    testHeader, ok := responseData.Headers["Test-Header"]
    g.Expect(ok).To(gm.BeTrue())
    g.Expect(testHeader).To(gm.ConsistOf("Foo"))

    g.Expect(responseData.Body).To(gm.Equal(expectedBodyBytes))
}

func TestAbort(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    expectedBodyBytes := []byte(
        `{"errors":[{"code":0,"message":"Expected failure while testing"}]}`,
    )

    responseData := Abort(
        context.Background(),
        JSONEncoding,
        http.StatusConflict,
        neterr.NewCodedError(
            0, "Expected failure while testing",
        ),
    )
    t.Log(responseData)

    g.Expect(responseData.Body).To(gm.Equal(expectedBodyBytes))
}

func TestRespond(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    expectedBodyBytes := []byte(`{"foo":"bar"}`)

    responseData := Respond(
        context.Background(),
        JSONEncoding,
        http.StatusOK,
        Body(map[string]string{
            "foo": "bar",
        }),
    )
    t.Log(responseData)

    g.Expect(responseData.Body).To(gm.Equal(expectedBodyBytes))
}
