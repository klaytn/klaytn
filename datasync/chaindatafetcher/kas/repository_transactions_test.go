package kas

func (s *SuiteRepository) TestInsertTransaction() {
	s.T().Log(s.repo.db.Dialect().GetName())
}
