package vial

import (
    "testing"

    gm "github.com/onsi/gomega"
    "github.com/google/uuid"
)

var (
    pathParamTestStrings = map[string]string {
        "": "foo",
        "string": "bar",
        "int": "1",
        "integer": "2",
        "float": "4.0",
        "uuid": "e02d6750-75c7-4a7e-9baa-6ffd70d6af9f",
    }
    pathParamTestValues = map[string]interface{} {
        "": "foo",
        "string": "bar",
        "int": 1,
        "integer": 2,
        "float": 4.0,
        "uuid": (func(u uuid.UUID, err error) uuid.UUID {
            if err != nil {
                panic(err)
            }
            return u
        })(uuid.Parse(
            "e02d6750-75c7-4a7e-9baa-6ffd70d6af9f",
        )),
    }
)

func TestPathParamMatchers(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    for key, matcher := range pathParamMatchers {
        s := pathParamTestStrings[key]
        v := pathParamTestValues[key]

        result, err := matcher.Coercer(s)
        g.Expect(err).ToNot(gm.HaveOccurred())
        g.Expect(result).To(gm.BeEquivalentTo(v))
    }
}
