package setting

// sections 存放app.yaml配置
var sections = make(map[string]interface{})

// GetSections return sections.
func GetSections() map[string]interface{} {
	return sections
}

// ReadSection read section config.
func (s *Setting) ReadSection(k string, v interface{}) error {
	err := s.vp.UnmarshalKey(k, v)
	if err != nil {
		return err
	}

	if _, ok := sections[k]; !ok {
		sections[k] = v
	}

	return nil
}

// ReloadAllSection if config has changed reload all config.
func (s *Setting) ReloadAllSection() error {
	for k, v := range sections {
		err := s.ReadSection(k, v)
		if err != nil {
			return err
		}
	}

	return nil
}
