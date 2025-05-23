package advisory

import (
	"context"
	"fmt"
	"sort"

	v2 "github.com/chainguard-dev/advisory-schema/pkg/advisory/v2"
	"github.com/wolfi-dev/wolfictl/pkg/configs"
)

type indexAdapter struct {
	index *configs.Index[v2.Document]
}

// AdaptIndex creates an implementation of advisory.Getter using an existing
// instance of `*configs.Index[v2.Document]`.
func AdaptIndex(index *configs.Index[v2.Document]) Getter {
	return indexAdapter{
		index: index,
	}
}

func (c indexAdapter) PackageNames(_ context.Context) ([]string, error) {
	documents := c.index.Select().Configurations()

	packageNames := make([]string, 0, len(documents))
	for _, d := range documents {
		packageNames = append(packageNames, d.Package.Name)
	}

	// Sort the package names for consistency
	sort.Strings(packageNames)

	return packageNames, nil
}

func (c indexAdapter) Advisories(_ context.Context, packageName string) ([]v2.PackageAdvisory, error) {
	entry, err := c.index.Select().WhereName(packageName).First()
	if err != nil {
		return nil, fmt.Errorf("getting advisory document for %q: %w", packageName, err)
	}

	doc := entry.Configuration()

	name := doc.Package.Name

	pkgAdvs := make([]v2.PackageAdvisory, 0, len(doc.Advisories))
	for _, adv := range doc.Advisories {
		pkgAdv := v2.PackageAdvisory{
			PackageName: name,
			Advisory:    adv,
		}
		pkgAdvs = append(pkgAdvs, pkgAdv)
	}

	return pkgAdvs, nil
}
