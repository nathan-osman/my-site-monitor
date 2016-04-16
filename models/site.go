package models

const (
	StatusUnknown = iota
	StatusOK
	StatusWarning
	StatusError
)

// Site represents an individual website being monitored.
type Site struct {
	ID         int    // Unique identifier
	Name       string // Descriptive name
	URL        string // Full URI of site
	Status     int    // Current status of the site
	IntervalMS int    // Time between pings
	WarningMS  int    // Time before raising a warning
	ErrorMS    int    // Time before raising an error
}

// InitSites executes the SQL necessary to create the sites table.
func InitSites() error {
	_, err := db.Exec(
		`
		CREATE TABLE Sites (
			ID         serial PRIMARY KEY,
			Name       text NOT NULL,
			URL        text NOT NULL,
			Status     integer NOT NULL,
			IntervalMS integer NOT NULL,
			WarningMS  integer NOT NULL,
			ErrorMS    integer NOT NULL
		)
		`,
	)
	return err
}

// ListSites retrieves all sites in the database.
func ListSites() ([]*Site, error) {
	sites := []*Site{}
	rows, err := db.Query(
		`
		SELECT
			ID, Name, URL, Status, IntervalMS, WarningMS, ErrorMS
		FROM
			Sites
		ORDER BY
		    URL ASC
		`,
	)
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		site := &Site{}
		err := row.Scan(
			&site.ID,
			&site.Name,
			&site.URL,
			&site.Status,
			&site.IntervalMS,
			&site.WarningMS,
			&site.ErrorMS,
		)
		if err != nil {
			return nil, err
		}
		sites = append(sites, site)
	}
	return sites
}

// Create commits the site to the database.
func (s *Site) Create() error {
	r, err := db.Exec(
		`
		INSERT INTO
			Sites (Name, URL, Status, IntervalMS, WarningMS, ErrorMS)
		VALUES
			(?, ?, ?, ?, ?, ?)
		`,
		s.Name,
		s.URL,
		s.Status,
		s.IntervalMS,
		s.WarningMS,
		s.ErrorMS,
	)
	if err != nil {
		return err
	}
	i, err := r.LastInsertId()
	if err != nil {
		return err
	}
	s.ID = i
	return nil
}

// Update commits changes to the site to the database.
func (s *Site) Update() error {
	_, err := db.Exec(
		`
		UPDATE
			Sites
		SET
		    Name       = ?,
		    URL        = ?,
		    Status     = ?,
		    IntervalMS = ?,
		    WarningMS  = ?,
		    ErrorMs    = ?
		WHERE
			ID = ?
		`,
		s.Name,
		s.URL,
		s.Status,
		s.IntervalMS,
		s.WarningMS,
		s.ErrorMS,
		s.ID,
	)
	return err
}

// Delete removes the site from the database.
func (s *Site) Delete() error {
	_, err := db.Exec(
		`
		DELETE FROM
		    Sites
		WHERE
		    ID = ?
		`,
		s.ID,
	)
	return err
}
