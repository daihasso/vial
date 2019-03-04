package vial

import (
    "net/http"
    "testing"

    gm "github.com/onsi/gomega"
)

func TestRequestPathString(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    req, err := http.NewRequest("GET", "/thread/bar", nil)
    g.Expect(err).To(gm.BeNil())

    srvReq := NewServerRequest(req, map[string]string{
        "foo": "bar",
    })

    val, ok := srvReq.PathString("foo")
    g.Expect(ok).To(gm.BeTrue())
    g.Expect(val).To(gm.Equal("bar"))
}

func TestRequestPathInt(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    req, err := http.NewRequest("GET", "/thread/52", nil)
    g.Expect(err).To(gm.BeNil())

    srvReq := NewServerRequest(req, map[string]string{
        "foo": "52",
    })

    val, ok := srvReq.PathInt("foo")
    g.Expect(ok).To(gm.BeTrue())
    g.Expect(val).To(gm.Equal(52))
}

func TestRequestPathUUID(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    req, err := http.NewRequest(
        "GET", "/thread/781d3d17bbbb4b798c4875e326e55275", nil,
    )
    g.Expect(err).To(gm.BeNil())

    srvReq := NewServerRequest(req, map[string]string{
        "foo": "781d3d17bbbb4b798c4875e326e55275",
    })

    val, ok := srvReq.PathUUID("foo")
    g.Expect(ok).To(gm.BeTrue())
    g.Expect(val.String()).To(gm.Equal("781d3d17-bbbb-4b79-8c48-75e326e55275"))
}

func TestRequestQueryParams(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    req, err := http.NewRequest("GET", "/foo?bar=baz", nil)
    g.Expect(err).To(gm.BeNil())

    srvReq := NewServerRequest(req, map[string]string{})

    g.Expect(srvReq.QueryParams()).To(
        gm.HaveKeyWithValue("bar", []string{"baz"}),
    )
}

func TestRequestQueryParam(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    req, err := http.NewRequest("GET", "/foo?bar=baz&bar=baz2", nil)
    g.Expect(err).To(gm.BeNil())

    srvReq := NewServerRequest(req, map[string]string{})

    val, ok := srvReq.QueryParam("bar")
    g.Expect(ok).To(gm.BeTrue())
    g.Expect(val).To(gm.Equal("baz"))
}

func TestRequestQueryParamMultiple(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    req, err := http.NewRequest("GET", "/foo?bar=baz&bar=baz2", nil)
    g.Expect(err).To(gm.BeNil())

    srvReq := NewServerRequest(req, map[string]string{})

    val, ok := srvReq.QueryParamMultiple("bar")
    g.Expect(ok).To(gm.BeTrue())
    g.Expect(val).To(gm.ConsistOf("baz", "baz2"))
}

func TestRequestPathWrongType(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    req, err := http.NewRequest(
        "GET", "/thread/781d3d17bbbb4b798c4875e326e55275", nil,
    )
    g.Expect(err).To(gm.BeNil())

    srvReq := NewServerRequest(req, map[string]string{
        "foo": "781d3d17bbbb4b798c4875e326e55275",
    })

    val, ok := srvReq.PathInt("foo")
    g.Expect(ok).ToNot(gm.BeTrue())
    g.Expect(val).To(gm.BeZero())
}
