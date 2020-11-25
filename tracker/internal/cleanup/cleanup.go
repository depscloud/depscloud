package cleanup

// Patch represents a modification to the database.
type Patch = func(s *Servers) error

// Run executes a series of patches to clean up data in the database.
func Run(servers *Servers) error {
	patches := []Patch{
		organizationModule,
	}

	for _, patch := range patches {
		if err := patch(servers); err != nil {
			return err
		}
	}

	return nil
}
