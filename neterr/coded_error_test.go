package neterr

import (
    "testing"
    "errors"
    "encoding/json"

    gm "github.com/onsi/gomega"
)

func TestNewCodedError(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    codedError := NewCodedError(1, "My custom error")
    g.Expect(codedError.code).To(gm.Equal(1))
    g.Expect(codedError.isVialError).To(gm.BeFalse())

    g.Expect(codedError.message).To(gm.Equal("My custom error"))
}

func TestNewCodedErrorFromExistingError(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    codedError := CodedErrorFromError(1, errors.New("Heyooo"))
    g.Expect(codedError.code).To(gm.Equal(1))
    g.Expect(codedError.isVialError).To(gm.BeFalse())

    g.Expect(codedError.message).To(gm.Equal("Heyooo"))
}

func TestCodedErrorMarshal(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    codedError := NewCodedError(1, "My custom error")
    marshaled, err := json.Marshal(codedError)
    g.Expect(err).To(gm.BeNil())

    g.Expect(marshaled).To(gm.BeEquivalentTo(
        `{"code":1,"message":"My custom error"}`,
    ))
}

func TestFrameworkCodedErrorMarshal(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    codedError := CodedError{
        code: 1,
        message: "My custom error",
        isVialError: true,
    }
    marshaled, err := json.Marshal(codedError)
    g.Expect(err).To(gm.BeNil())

    g.Expect(marshaled).To(gm.BeEquivalentTo(
        `{"code":1,"message":"My custom error","vial_error":true}`,
    ))
}

func TestCodedErrorUnmarshal(t *testing.T) {
    g := gm.NewGomegaWithT(t)


    marshaled := []byte(`{"code":1,"message":"Framework error"}`)
    codedError := CodedError{}
    err := json.Unmarshal(marshaled, &codedError)
    g.Expect(err).To(gm.BeNil())

    g.Expect(codedError.code).To(gm.Equal(1))
    g.Expect(codedError.isVialError).To(gm.BeFalse())
}

func TestFrameworkCodedErrorUnmarshal(t *testing.T) {
    g := gm.NewGomegaWithT(t)


    marshaled := []byte(
        `{"code":1,"message":"Framework error","vial_error":true}`,
    )
    codedError := CodedError{}
    err := json.Unmarshal(marshaled, &codedError)
    g.Expect(err).To(gm.BeNil())

    g.Expect(codedError.code).To(gm.Equal(1))
    g.Expect(codedError.isVialError).To(gm.BeTrue())
}

func TestCodedErrorAccessors(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    codedError := NewCodedError(1, "My custom error")

    g.Expect(codedError.Code()).To(gm.Equal(1))
    g.Expect(codedError.Message()).To(gm.Equal("My custom error"))
    g.Expect(codedError.IsVialError()).To(gm.BeFalse())
}
