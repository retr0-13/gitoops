package github

import (
	"github.com/ovotech/gitoops/pkg/database"
)

// We only support one organization at the moment, so this ingestor is a bit of gimmick

type OrganizationsIngestor struct {
	gqlclient *GraphQLClient
	db        *database.Database
	data      *OrganizationsData
}

type OrganizationsData struct {
	Nodes []struct {
		Login string `json:"login"`
		URL   string `json:"url"`
	} `json:"nodes"`
}

func (ing *OrganizationsIngestor) FetchData() {
	ing.data = &OrganizationsData{
		Nodes: []struct {
			Login string "json:\"login\""
			URL   string "json:\"url\""
		}{
			{
				Login: ing.gqlclient.organization,
				URL:   "https://github.com/" + ing.gqlclient.organization,
			},
		},
	}
}

func (ing *OrganizationsIngestor) Sync() {
	ing.insertOrganizations()
}

func (ing *OrganizationsIngestor) insertOrganizations() {
	organizations := []map[string]string{}

	for _, orgNode := range ing.data.Nodes {
		organizations = append(organizations, map[string]string{
			"url":   orgNode.URL,
			"login": orgNode.Login,
		})
	}

	ing.db.Run(`
	UNWIND $organizations AS organization

	MERGE (o:Organization{id: organization.url})

	SET o.login = organization.login,
	o.url = organization.url
	`, map[string]interface{}{"organizations": organizations})
}