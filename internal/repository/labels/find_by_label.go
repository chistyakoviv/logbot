package labels

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/chistyakoviv/logbot/internal/model"
)

func (r *repository) FindByLabel(ctx context.Context, label string) ([]*model.Label, error) {
	q := db.Query{
		Name: "repository.labels.find_by_label",
		Sqlizer: r.sq.Select(labelsTableColumns...).
			From(labelsTable).
			Where(sq.Expr("? = ANY("+labelsTableColumnLabels+")", label)),
		// Find rows where labels contains any element of a given array (&& - shares at least one element)
		// Where(sq.Expr(labelsTableColumnLabels+" && ?::varchar[]", []string{"label1", "lable2"})),
		// Find rows where array contains all specified values (@> - contains all given values; <@ - array is subset of given array)
		// Where(sq.Expr(labelsTableColumnLabels+" @> ?::varchar[]", []string{"label1", "lable2"})),
	}

	var row []LabelsRow
	if err := r.db.DB().Selectx(ctx, &row, q); err != nil {
		return nil, fmt.Errorf("%s: %w", q.Name, err)
	}

	out := make([]*model.Label, 0, len(row))
	for _, v := range row {
		out = append(out, ToModel(&v))
	}

	return out, nil
}
