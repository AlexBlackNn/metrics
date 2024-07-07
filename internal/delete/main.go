package main

import (
	"fmt"
	"github.com/AlexBlackNn/metrics/internal/domain/models"
	"os"
	"text/template"
)

// Helper function to get the type
func GetType(m models.MetricGetter) string {
	return m.GetType()
}

// Helper function to get the name
func GetName(m models.MetricGetter) string {
	return m.GetName()
}

func main() {
	metric := &models.Metric[int64]{Type: "counter", Name: "test_counter"}

	// Create the template with a function to call GetType
	tpl := template.Must(template.New("sqlQuery").Funcs(template.FuncMap{
		"GetType": GetType,
		"GetName": GetName,
	}).Parse(`
      WITH LatestCounter AS (
        SELECT
          MAX(created) AS latest_created
        FROM
          {{GetType .}}
        WHERE
          {{GetType .}}.name = '{{GetName .}}'
      )
      SELECT
        t.name, c.name, c.value
      FROM
        {{GetType .}} as c
      JOIN
        app.types as t ON c.metric_id = t.uuid
      JOIN
        LatestCounter lc ON c.created = lc.latest_created
      WHERE
        c.name = '{{GetName .}}'
      ORDER BY c.created;
  `))

	// Execute the template and print the result
	err := tpl.Execute(os.Stdout, metric)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return
	}
}
