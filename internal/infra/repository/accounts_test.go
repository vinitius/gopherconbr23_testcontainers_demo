package repository

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type AccountsSQLRepositorySuite struct {
	SQLRepositorySuite
}

func TestAccountsSQLRepositorySuite(t *testing.T) {
	suite.Run(t, new(AccountsSQLRepositorySuite))
}

func (s *AccountsSQLRepositorySuite) assertAccountInDatabase(expected AccountModel) {
	s.T().Helper()
	var actual AccountModel
	err := s.DB.NewSelect().Model(&expected).WherePK().Scan(s.Ctx, &actual)
	s.Require().NoError(err)
	s.Require().Empty(cmp.Diff(expected, actual))
}

func (s *AccountsSQLRepositorySuite) SetupTest() {
	s.Ctx = context.TODO()

	// Truncate before every test
	_, err := s.DB.NewTruncateTable().Model(new(AccountModel)).Exec(s.Ctx)
	s.Require().NoError(err)

	// Add fixtures
	_, err = s.DB.NewInsert().Model(&[]AccountModel{
		{
			ID:       "1",
			Type:     "regular",
			Label:    "misc",
			Currency: "ETH",
		},
		{
			ID:       "2",
			Type:     "regular",
			Label:    "misc",
			Currency: "DOGE",
		},
	}).Exec(s.Ctx)
	s.Require().NoError(err)
}

func (s *AccountsSQLRepositorySuite) TestFindAll() {
	// Arrange
	underTest := NewAccountsSQLRepository(s.DB)
	expected := []AccountModel{
		{
			ID:       "1",
			Type:     "regular",
			Label:    "misc",
			Currency: "ETH",
		},
		{
			ID:       "2",
			Type:     "regular",
			Label:    "misc",
			Currency: "DOGE",
		},
	}

	// Act
	actual, err := underTest.FindAll(s.Ctx)

	// Assert
	s.Require().NoError(err)
	s.Require().Empty(cmp.Diff(actual, expected))
}

func (s *AccountsSQLRepositorySuite) TestCreate() {
	// Arrange
	underTest := NewAccountsSQLRepository(s.DB)
	id := uuid.NewString()
	account := AccountModel{
		ID:       id,
		Type:     "regular",
		Label:    "miscellaneous",
		Currency: "BTC",
	}

	// Act
	err := underTest.Create(s.Ctx, account)

	// Assert
	s.Require().NoError(err)
	s.assertAccountInDatabase(account)
}
