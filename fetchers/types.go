package fetchers

import (
	"time"
)

type Commit struct {
	Author struct {
		Raw  string `json:"raw"`
		User struct {
			DisplayName string `json:"display_name"`
			Links       struct {
				Avatar struct {
					Href string `json:"href"`
				} `json:"avatar"`
				HTML struct {
					Href string `json:"href"`
				} `json:"html"`
				Self struct {
					Href string `json:"href"`
				} `json:"self"`
			} `json:"links"`
			Type     string `json:"type"`
			Username string `json:"username"`
			UUID     string `json:"uuid"`
		} `json:"user"`
	} `json:"author"`
	Date  time.Time `json:"date"`
	Hash  string    `json:"hash"`
	Links struct {
		Approve struct {
			Href string `json:"href"`
		} `json:"approve"`
		Comments struct {
			Href string `json:"href"`
		} `json:"comments"`
		Diff struct {
			Href string `json:"href"`
		} `json:"diff"`
		HTML struct {
			Href string `json:"href"`
		} `json:"html"`
		Patch struct {
			Href string `json:"href"`
		} `json:"patch"`
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		Statuses struct {
			Href string `json:"href"`
		} `json:"statuses"`
	} `json:"links"`
	Message string `json:"message"`
	Parents []struct {
		Hash  string `json:"hash"`
		Links struct {
			HTML struct {
				Href string `json:"href"`
			} `json:"html"`
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
		} `json:"links"`
		Type string `json:"type"`
	} `json:"parents"`
	Repository struct {
		FullName string `json:"full_name"`
		Links    struct {
			Avatar struct {
				Href string `json:"href"`
			} `json:"avatar"`
			HTML struct {
				Href string `json:"href"`
			} `json:"html"`
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
		} `json:"links"`
		Name string `json:"name"`
		Type string `json:"type"`
		UUID string `json:"uuid"`
	} `json:"repository"`
	Type string `json:"type"`
}

type Commits struct {
	Next    string   `json:"next"`
	Pagelen int      `json:"pagelen"`
	Values  []Commit `json:"values"`
}
