package coinmarketcal

import (
	"context"
	"time"
)

type Event struct {
	ID    int `json:"id"`
	Title struct {
		En string `json:"en"`
		Ko string `json:"ko"`
		Ru string `json:"ru"`
		Tr string `json:"tr"`
		Ja string `json:"ja"`
		Es string `json:"es"`
		Pt string `json:"pt"`
		ID string `json:"id"`
	} `json:"title"`
	Coins []struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Symbol      string `json:"symbol"`
		Fullname    string `json:"fullname"`
		Upcoming    int    `json:"upcoming"`
		Popular     int    `json:"popular"`
		Influential int    `json:"influential"`
		Catalyst    int    `json:"catalyst"`
	} `json:"coins"`
	DateEvent      time.Time `json:"date_event"`
	DisplayedDate  string    `json:"displayed_date"`
	CanOccurBefore bool      `json:"can_occur_before"`
	Categories     []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"categories"`
	Proof       string    `json:"proof"`
	Source      string    `json:"source"`
	CreatedDate time.Time `json:"created_date"`
	Description struct {
		En string `json:"en"`
		Ko string `json:"ko"`
		Ru string `json:"ru"`
		Tr string `json:"tr"`
		Ja string `json:"ja"`
		Es string `json:"es"`
		Pt string `json:"pt"`
		ID string `json:"id"`
	} `json:"description"`
	Percentage           float64 `json:"percentage"`
	VoteCount            int     `json:"vote_count"`
	IsTrending           bool    `json:"is_trending"`
	IsPopular            bool    `json:"is_popular"`
	TrendingIndex        int     `json:"trending_index"`
	PopularIndex         int     `json:"popular_index"`
	InfluentialScore     int     `json:"influential_score"`
	CatalystScore        int     `json:"catalyst_score"`
	ConfirmedByOfficials bool    `json:"confirmed_by_officials"`
	AlertCount           int     `json:"alert_count"`
	OriginalSource       string  `json:"original_source"`
	VoteHistory          []struct {
		Value bool      `json:"value"`
		Date  time.Time `json:"date"`
	} `json:"vote_history"`
	ViewHistory []struct {
		Date time.Time `json:"date"`
	} `json:"view_history"`
}

type EventsRequest struct {
	Page           *int    `url:"page,omitempty"`
	Max            *int    `url:"max,omitempty"`
	DateRangeStart *string `url:"dateRangeStart,omitempty"`
	DateRangeEnd   *string `url:"dateRangeEnd,omitempty"`
	Coins          *string `url:"coins,omitempty"`
	Categories     *string `url:"categories,omitempty"`
	SortBy         *string `url:"sortBy,omitempty"`
	ShowOnly       *string `url:"showOnly,omitempty"`
	ShowViews      *string `url:"showViews,omitempty"`
	ShowVotes      *bool   `url:"showVotes,omitempty"`
	Translations   *string `url:"translations,omitempty"`
}

func (c *Client) GetEvents(ctx context.Context, opt EventsRequest) (DefaultBody[[]Event], error) {
	response := DefaultBody[[]Event]{}
	_, err := c.GetCall(ctx, "/events", opt, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}
