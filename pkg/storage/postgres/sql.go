package postgres

import "text/template"

type Tmpl map[string]*template.Template

func NewTemplate() Tmpl {
	tmpl := Tmpl{}
	tmpl["GetMetric"] = template.Must(template.New("sqlQuery").Funcs(template.FuncMap{
		"GetType": GetType,
	}).Parse(`
      WITH LatestCounter AS (
        SELECT
          MAX(created) AS latest_created
        FROM
          app.{{GetType .}}_part
        WHERE
          app.{{GetType .}}_part.name = $1
      )
      SELECT
        t.name, c.name, c.value
      FROM
        app.{{GetType .}}_part as c
      JOIN
        app.types as t ON c.metric_id = t.uuid
      JOIN
        LatestCounter lc ON c.created = lc.latest_created
      WHERE
        c.name = $1
      ORDER BY c.created;
  `))

	tmpl["UpdateMetric"] = template.Must(template.New("sqlQuery").Funcs(template.FuncMap{
		"GetType": GetType,
	}).Parse(`
      INSERT INTO
    app.{{GetType .}}_part (metric_id, name, value)
	VALUES ((SELECT uuid FROM app.types WHERE name = $1), $2, $3)
	`))
	return tmpl
}
