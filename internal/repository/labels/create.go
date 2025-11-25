package labels

import (
	"context"
	"fmt"
	"strings"

	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
)

func (r *repository) Create(ctx context.Context, in *model.Label) (*model.Label, error) {
	row := FromModel(in)

	q := db.Query{
		Name: "repository.labels.create",
		Sqlizer: r.sq.Insert(labelsTable).
			Columns(labelsTableColumns...).
			Values(row.Values()...).
			Suffix("RETURNING " + strings.Join(labelsTableColumns, ",")),
	}

	var out LabelsRow
	if err := r.db.DB().Getx(ctx, &out, q); err != nil {
		return nil, fmt.Errorf("%s: %w", q.Name, err)
	}

	return ToModel(&out), nil
}
