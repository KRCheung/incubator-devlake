/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package e2e

import (
	"fmt"
	"github.com/apache/incubator-devlake/models/domainlayer/code"
	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"testing"

	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/github/impl"
	"github.com/apache/incubator-devlake/plugins/github/tasks"
)

func TestRepoDataFlow(t *testing.T) {
	var plugin impl.Github
	dataflowTester := e2ehelper.NewDataFlowTester(t, "gitlab", plugin)

	githubRepository := &models.GithubRepo{
		GithubId: 134018330,
	}
	taskData := &tasks.GithubTaskData{
		Options: &tasks.GithubOptions{
			ConnectionId: 1,
			Owner:        "panjf2000",
			Repo:         "ants",
			Config: models.Config{
				PrType:      "type/(.*)$",
				PrComponent: "component/(.*)$",
			},
		},
		Repo: githubRepository,
	}

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_github_api_repositories.csv", "_raw_github_api_repositories")

	// verify extraction
	dataflowTester.FlushTabler(&models.GithubRepo{})
	dataflowTester.Subtask(tasks.ExtractApiRepoMeta, taskData)
	dataflowTester.VerifyTable(
		models.GithubRepo{},
		fmt.Sprintf("./snapshot_tables/%s.csv", models.GithubRepo{}.TableName()),
		[]string{"connection_id", "github_id"},
		[]string{
			"name",
			"html_url",
			"description",
			"owner_id",
			"owner_login",
			"language",
			"created_date",
			"updated_date",
			"_raw_data_params",
			"_raw_data_table",
			"_raw_data_id",
			"_raw_data_remark",
		},
	)

	// verify extraction
	dataflowTester.FlushTabler(&code.Repo{})
	dataflowTester.FlushTabler(&ticket.Board{})
	dataflowTester.Subtask(tasks.ConvertRepoMeta, taskData)
	dataflowTester.VerifyTable(
		code.Repo{},
		fmt.Sprintf("./snapshot_tables/%s.csv", code.Repo{}.TableName()),
		[]string{"id"},
		[]string{
			"name",
			"url",
			"description",
			"owner_id",
			"language",
			"forked_from",
			"created_date",
			"updated_date",
			"deleted",
		},
	)
	dataflowTester.VerifyTable(
		ticket.Board{},
		fmt.Sprintf("./snapshot_tables/%s.csv", ticket.Board{}.TableName()),
		[]string{"id"},
		[]string{
			"name",
			"description",
			"url",
			"created_date",
		},
	)
}