package vial

import (
    "testing"

    "github.com/google/uuid"
    gm "github.com/onsi/gomega"
)

func TestParseRouteEmpty(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    route, err := ParseRoute("/test/<name>")
    g.Expect(err).ToNot(gm.HaveOccurred())

    pp, err := route.PathParams("/test/foo")
    g.Expect(err).ToNot(gm.HaveOccurred())

    g.Expect(pp.String("name")).To(gm.BeEquivalentTo("foo"))
}

func TestParseRouteString(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    route, err := ParseRoute("/test/<string:name>")
    g.Expect(err).ToNot(gm.HaveOccurred())

    pp, err := route.PathParams("/test/foo")
    g.Expect(err).ToNot(gm.HaveOccurred())

    g.Expect(pp.String("name")).To(gm.BeEquivalentTo("foo"))
}

func TestParseRouteInt(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    route, err := ParseRoute("/test/<int:id>")
    g.Expect(err).ToNot(gm.HaveOccurred())

    pp, err := route.PathParams("/test/1")
    g.Expect(err).ToNot(gm.HaveOccurred())

    g.Expect(pp.Int("id")).To(gm.BeEquivalentTo(1))
}

func TestParseRouteInteger(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    route, err := ParseRoute("/test/<integer:count>")
    g.Expect(err).ToNot(gm.HaveOccurred())

    pp, err := route.PathParams("/test/1")
    g.Expect(err).ToNot(gm.HaveOccurred())

    g.Expect(pp.Int("count")).To(gm.BeEquivalentTo(1))
}

func TestParseRouteFloat(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    route, err := ParseRoute("/test/<float:totalVelocity>")
    g.Expect(err).ToNot(gm.HaveOccurred())

    pp, err := route.PathParams("/test/1.0")
    g.Expect(err).ToNot(gm.HaveOccurred())

    g.Expect(pp.Float("totalVelocity")).To(gm.BeEquivalentTo(1.0))
}

func TestParseRouteUUID(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    expected, err := uuid.NewRandom()
    g.Expect(err).ToNot(gm.HaveOccurred())

    route, err := ParseRoute("/test/<uuid:post_id>")
    g.Expect(err).ToNot(gm.HaveOccurred())

    pp, err := route.PathParams("/test/" + expected.String())
    g.Expect(err).ToNot(gm.HaveOccurred())

    g.Expect(pp.UUID("post_id")).To(gm.BeEquivalentTo(expected))
}
