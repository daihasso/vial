package vial

import (
    "testing"

    "github.com/google/uuid"
    gm "github.com/onsi/gomega"
)

func TestPathParamString(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    pathParams := PathParams{
        "foo": "bar",
    }

    val, err := pathParams.String("foo")
    g.Expect(err).ToNot(gm.HaveOccurred())
    g.Expect(val).To(gm.Equal("bar"))
}

func TestPathParamFloat(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    pathParams := PathParams{
        "foo": 55.1,
    }

    val, err := pathParams.Float("foo")
    g.Expect(err).ToNot(gm.HaveOccurred())
    g.Expect(val).To(gm.Equal(55.1))
}

func TestPathParamInt(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    pathParams := PathParams{
        "foo": 55,
    }

    val, err := pathParams.Int("foo")
    g.Expect(err).ToNot(gm.HaveOccurred())
    g.Expect(val).To(gm.Equal(55))
}

func TestPathParamUUID(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    expectedUuid := "781d3d17bbbb4b798c4875e326e55275"
    actualUuid, err := uuid.Parse(expectedUuid)
    g.Expect(err).To(gm.BeNil())
    pathParams := PathParams{
        "foo": actualUuid,
    }

    val, err := pathParams.UUID("foo")
    g.Expect(err).ToNot(gm.HaveOccurred())
    g.Expect(val).To(gm.Equal(actualUuid))
}

func TestPathParamWrongType(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    pathParams := PathParams{
        "foo": 1,
    }

    _, err := pathParams.String("foo")
    g.Expect(err).To(gm.HaveOccurred())
    g.Expect(IsWrongPathParamType(err)).To(gm.BeTrue())
}

func TestPathParamNotExist(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    pathParams := PathParams{
        "foo": 1,
    }

    _, err := pathParams.String("bar")
    g.Expect(err).To(gm.HaveOccurred())
    g.Expect(IsPathParamDoesNotExist(err)).To(gm.BeTrue())
}
