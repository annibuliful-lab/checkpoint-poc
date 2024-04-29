package note

import (
	"context"
	"log"
)

type NoteResolver struct{}

func (NoteResolver) CreateStationDashboardReport(ctx context.Context, args struct{ Note string } ) (*StationDashboardReport, error) {
	// authorization := auth.GetAuthorizationContext(ctx)
	log.Println(ctx.Value("accountId").(string))

	createNote := StationDashboardReport{
		Id:        "Mock-ID-02",
		ProjectId: "Mock-ProjectId-02",
		Note:      args.Note,
	}

	return &createNote, nil
}
