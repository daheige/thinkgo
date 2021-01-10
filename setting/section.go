package setting

// GetSections return sections.
func (s *Setting) GetSections() map[string]interface{} {
	return s.sections
}

// IsSet check key is set.
func (s *Setting) IsSet(key string) bool {
	return s.vp.IsSet(key)
}

// ReadSection read section config.
func (s *Setting) ReadSection(k string, v interface{}) error {
	err := s.vp.UnmarshalKey(k, v)
	if err != nil {
		return err
	}

	if _, ok := s.sections[k]; !ok {
		s.sections[k] = v
	}

	return nil
}

// ReloadAllSection if config has changed reload all config.
func (s *Setting) ReloadAllSection() error {
	for k, v := range s.sections {
		err := s.ReadSection(k, v)
		if err != nil {
			return err
		}
	}

	return nil
}
