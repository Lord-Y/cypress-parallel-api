// Package hooks will manage all hooks requirements
package hooks

import (
	"database/sql"

	"github.com/Lord-Y/cypress-parallel-api/commons"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/syyongx/php2go"
)

// getProjectInfos collect requirements to start the unit testing
func (p *plain) getProjectInfos() (z map[string]string, err error) {
	db, err := sql.Open(
		"postgres",
		commons.BuildDSN(),
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to DB")
		return z, err
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT * FROM projects WHERE project_name = $1 LIMIT 1")
	if err != nil && err != sql.ErrNoRows {
		return z, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(
		php2go.Addslashes(p.ProjectName),
	)
	if err != nil && err != sql.ErrNoRows {
		return z, err
	}

	columns, err := rows.Columns()
	if err != nil {
		return z, err
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	m := make(map[string]string)
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return
		}
		var value string
		for i, col := range values {
			if col == nil {
				value = ""
			} else {
				value = php2go.Stripslashes(string(col))
			}
			m[columns[i]] = value
		}
	}
	if err = rows.Err(); err != nil {
		return z, err
	}
	return m, nil
}

// create will insert executions in DB
func (p *execution) create() (z int64, err error) {
	db, err := sql.Open(
		"postgres",
		commons.BuildDSN(),
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to DB")
		return z, err
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO executions(project_id, branch, execution_status, uniq_id, spec, result) VALUES($1, $2, $3, $4, $5, $6) RETURNING execution_id")
	if err != nil && err != sql.ErrNoRows {
		return z, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(
		p.projectID,
		php2go.Addslashes(p.branch),
		php2go.Addslashes(p.executionStatus),
		php2go.Addslashes(p.uniqID),
		php2go.Addslashes(p.spec),
		php2go.Addslashes(p.result),
	).Scan(&z)
	if err != nil && err != sql.ErrNoRows {
		return z, err
	}
	return z, nil
}

// getProjectAnnotations collect requirements to start the unit testing
func (p *projects) getProjectAnnotations() (z []map[string]interface{}, err error) {
	db, err := sql.Open(
		"postgres",
		commons.BuildDSN(),
	)
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT * FROM annotations WHERE project_id = $1 LIMIT 1")
	if err != nil && err != sql.ErrNoRows {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(
		p.Project_id,
	)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	columns, err := rows.Columns()
	if err != nil {
		return
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	m := make([]map[string]interface{}, 0)
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return
		}
		var value string
		sub := make(map[string]interface{})
		for i, col := range values {
			if col == nil {
				value = ""
			} else {
				value = php2go.Stripslashes(string(col))
			}
			sub[columns[i]] = value
		}
		m = append(m, sub)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return m, nil
}

// getProjectEnvironments collect requirements to start the unit testing
func (p *projects) getProjectEnvironments() (z []map[string]interface{}, err error) {
	db, err := sql.Open(
		"postgres",
		commons.BuildDSN(),
	)
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT * FROM environments WHERE project_id = $1 LIMIT 1")
	if err != nil && err != sql.ErrNoRows {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(
		p.Project_id,
	)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	columns, err := rows.Columns()
	if err != nil {
		return
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	m := make([]map[string]interface{}, 0)
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return
		}
		var value string
		sub := make(map[string]interface{})
		for i, col := range values {
			if col == nil {
				value = ""
			} else {
				value = php2go.Stripslashes(string(col))
			}
			sub[columns[i]] = value
		}
		m = append(m, sub)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return m, nil
}

// update will update pod_name field in DB
func (p *updatePodName) update() (err error) {
	db, err := sql.Open(
		"postgres",
		commons.BuildDSN(),
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to DB")
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("UPDATE executions SET pod_name = $1, execution_status = 'RUNNING' WHERE uniq_id = $2 AND spec = $3")
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(
		php2go.Addslashes(p.podName),
		php2go.Addslashes(p.uniqID),
		php2go.Addslashes(p.spec),
	).Scan()
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	return nil
}

// executionStatus get executions by status
func executionStatus(execution_status string) (z []map[string]interface{}, err error) {
	db, err := sql.Open(
		"postgres",
		commons.BuildDSN(),
	)
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT DISTINCT uniq_id FROM executions WHERE execution_status = $1")
	if err != nil && err != sql.ErrNoRows {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(
		execution_status,
	)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	columns, err := rows.Columns()
	if err != nil {
		return
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	m := make([]map[string]interface{}, 0)
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return
		}
		var value string
		sub := make(map[string]interface{})
		for i, col := range values {
			if col == nil {
				value = ""
			} else {
				value = php2go.Stripslashes(string(col))
			}
			sub[columns[i]] = value
		}
		m = append(m, sub)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return m, nil
}

// pgqueued get queued executions
func pgqueued(uniqId string) (z []map[string]interface{}, err error) {
	db, err := sql.Open(
		"postgres",
		commons.BuildDSN(),
	)
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT e.*, p.project_name FROM executions e LEFT JOIN projects p ON e.project_id = p.project_id WHERE e.execution_status = 'QUEUED' AND e.uniq_id = $1")
	if err != nil && err != sql.ErrNoRows {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(
		php2go.Addslashes(uniqId),
	)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	columns, err := rows.Columns()
	if err != nil {
		return
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	m := make([]map[string]interface{}, 0)
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return
		}
		var value string
		sub := make(map[string]interface{})
		for i, col := range values {
			if col == nil {
				value = ""
			} else {
				value = php2go.Stripslashes(string(col))
			}
			sub[columns[i]] = value
		}
		m = append(m, sub)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return m, nil
}

// countExecutions will count number of executions not in specified values
func countExecutions(uniq_id string) (z map[string]string, err error) {
	db, err := sql.Open(
		"postgres",
		commons.BuildDSN(),
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to DB")
		return z, err
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT COUNT(execution_id) FROM executions WHERE uniq_id = $1 AND execution_status = 'RUNNING'")
	if err != nil && err != sql.ErrNoRows {
		return z, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(
		php2go.Addslashes(uniq_id),
	)
	if err != nil && err != sql.ErrNoRows {
		return z, err
	}

	columns, err := rows.Columns()
	if err != nil {
		return z, err
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	m := make(map[string]string)
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return
		}
		var value string
		for i, col := range values {
			if col == nil {
				value = ""
			} else {
				value = php2go.Stripslashes(string(col))
			}
			m[columns[i]] = value
		}
	}
	if err = rows.Err(); err != nil {
		return z, err
	}
	return m, nil
}
