package storage

import (
	"github.com/google/uuid"
	"github.com/guregu/null"

	"github.com/cmsgov/easi-app/pkg/models"
	"github.com/cmsgov/easi-app/pkg/testhelpers"
)

func (s StoreTestSuite) TestFetchBusinessCaseByID() {
	s.Run("golden path to fetch a business case", func() {
		intake := testhelpers.NewSystemIntake()
		err := s.store.SaveSystemIntake(&intake)
		s.NoError(err)
		businessCase := testhelpers.NewBusinessCase()
		businessCase.SystemIntakeID = intake.ID
		created, err := s.store.CreateBusinessCase(&businessCase)
		s.NoError(err)
		fetched, err := s.store.FetchBusinessCaseByID(created.ID)

		s.NoError(err, "failed to fetch business case")
		s.Equal(created.ID, fetched.ID)
		s.Equal(businessCase.EUAUserID, fetched.EUAUserID)
		s.Len(fetched.LifecycleCostLines, 2)
	})

	s.Run("cannot without an ID that exists in the db", func() {
		badUUID, _ := uuid.Parse("")
		fetched, err := s.store.FetchBusinessCaseByID(badUUID)

		s.Error(err)
		s.Equal("sql: no rows in result set", err.Error())
		s.Equal(&models.BusinessCase{}, fetched)
	})
}

func (s StoreTestSuite) TestFetchBusinessCasesByEuaID() {
	s.Run("golden path to fetch business cases", func() {
		intake := testhelpers.NewSystemIntake()
		intake.Status = models.SystemIntakeStatusSUBMITTED
		err := s.store.SaveSystemIntake(&intake)
		s.NoError(err)

		intake2 := testhelpers.NewSystemIntake()
		intake2.EUAUserID = intake.EUAUserID
		intake2.Status = models.SystemIntakeStatusSUBMITTED
		err = s.store.SaveSystemIntake(&intake2)
		s.NoError(err)

		businessCase := testhelpers.NewBusinessCase()
		businessCase.EUAUserID = intake.EUAUserID
		businessCase.SystemIntakeID = intake.ID

		businessCase2 := testhelpers.NewBusinessCase()
		businessCase2.EUAUserID = intake.EUAUserID
		businessCase2.SystemIntakeID = intake2.ID

		_, err = s.store.CreateBusinessCase(&businessCase)
		s.NoError(err)

		_, err = s.store.CreateBusinessCase(&businessCase2)
		s.NoError(err)

		fetched, err := s.store.FetchBusinessCasesByEuaID(businessCase.EUAUserID)

		s.NoError(err, "failed to fetch business cases")
		s.Len(fetched, 2)
		s.Len(fetched[0].LifecycleCostLines, 2)
		s.Equal(businessCase.EUAUserID, fetched[0].EUAUserID)
	})

	s.Run("fetches no results with other EUA ID", func() {
		fetched, err := s.store.FetchBusinessCasesByEuaID(testhelpers.RandomEUAID())

		s.NoError(err)
		s.Len(fetched, 0)
		s.Equal(models.BusinessCases{}, fetched)
	})
}

func (s StoreTestSuite) TestCreateBusinessCase() {
	s.Run("golden path to create a business case", func() {
		intake := testhelpers.NewSystemIntake()
		err := s.store.SaveSystemIntake(&intake)
		s.NoError(err)
		businessCase := models.BusinessCase{
			SystemIntakeID: intake.ID,
			EUAUserID:      testhelpers.RandomEUAID(),
			LifecycleCostLines: models.EstimatedLifecycleCosts{
				testhelpers.NewEstimatedLifecycleCost(testhelpers.EstimatedLifecycleCostOptions{}),
			},
		}
		created, err := s.store.CreateBusinessCase(&businessCase)

		s.NoError(err, "failed to create a business case")
		s.NotNil(created.ID)
		s.Equal(businessCase.EUAUserID, created.EUAUserID)
		s.Len(created.LifecycleCostLines, 1)
	})

	s.Run("requires a system intake ID", func() {
		businessCase := models.BusinessCase{
			EUAUserID: testhelpers.RandomEUAID(),
		}

		_, err := s.store.CreateBusinessCase(&businessCase)

		s.Error(err)
		s.Equal("pq: Could not complete operation in a failed transaction", err.Error())
	})

	s.Run("requires a system intake ID that exists in the db", func() {
		badintakeID := uuid.New()
		businessCase := models.BusinessCase{
			SystemIntakeID: badintakeID,
			EUAUserID:      testhelpers.RandomEUAID(),
		}

		_, err := s.store.CreateBusinessCase(&businessCase)

		s.Error(err)
		s.Equal("pq: Could not complete operation in a failed transaction", err.Error())
	})

	s.Run("cannot without a eua user id", func() {
		intake := testhelpers.NewSystemIntake()
		err := s.store.SaveSystemIntake(&intake)
		s.NoError(err)
		businessCase := models.BusinessCase{
			SystemIntakeID: intake.ID,
		}
		_, err = s.store.CreateBusinessCase(&businessCase)

		s.Error(err)
		s.Equal("pq: Could not complete operation in a failed transaction", err.Error())
	})
}

func (s StoreTestSuite) TestUpdateBusinessCase() {
	intake := testhelpers.NewSystemIntake()
	err := s.store.SaveSystemIntake(&intake)
	s.NoError(err)
	euaID := intake.EUAUserID
	businessCaseOriginal := testhelpers.NewBusinessCase()
	businessCaseOriginal.EUAUserID = euaID
	businessCaseOriginal.SystemIntakeID = intake.ID
	createdBusinessCase, err := s.store.CreateBusinessCase(&businessCaseOriginal)
	s.NoError(err)
	id := createdBusinessCase.ID
	year2 := models.LifecycleCostYear2
	year3 := models.LifecycleCostYear3
	solution := models.LifecycleCostSolutionA

	s.Run("golden path to update a business case", func() {
		expectedPhoneNumber := null.StringFrom("3452345678")
		expectedProjectName := null.StringFrom("Fake name")
		businessCaseToUpdate := models.BusinessCase{
			ID:                   id,
			ProjectName:          expectedProjectName,
			RequesterPhoneNumber: expectedPhoneNumber,
			PriorityAlignment:    null.String{},
			LifecycleCostLines: models.EstimatedLifecycleCosts{
				testhelpers.NewEstimatedLifecycleCost(testhelpers.EstimatedLifecycleCostOptions{}),
				testhelpers.NewEstimatedLifecycleCost(testhelpers.EstimatedLifecycleCostOptions{
					Year: &year2,
				}),
				testhelpers.NewEstimatedLifecycleCost(testhelpers.EstimatedLifecycleCostOptions{
					Solution: &solution,
				}),
			},
		}
		_, err := s.store.UpdateBusinessCase(&businessCaseToUpdate)
		s.NoError(err)
		//	fetch the newly updated business case
		updated, err := s.store.FetchBusinessCaseByID(id)
		s.NoError(err)
		s.Equal(expectedPhoneNumber, updated.RequesterPhoneNumber)
		s.Equal(expectedProjectName, updated.ProjectName)
		s.Equal(null.String{}, updated.PriorityAlignment)
		s.Equal(3, len(updated.LifecycleCostLines))
	})

	s.Run("lifecycle costs are recreated", func() {
		businessCaseToUpdate := models.BusinessCase{
			ID: id,
			LifecycleCostLines: models.EstimatedLifecycleCosts{
				testhelpers.NewEstimatedLifecycleCost(testhelpers.EstimatedLifecycleCostOptions{}),
				testhelpers.NewEstimatedLifecycleCost(testhelpers.EstimatedLifecycleCostOptions{
					Year:     &year2,
					Solution: &solution,
				}),
				testhelpers.NewEstimatedLifecycleCost(testhelpers.EstimatedLifecycleCostOptions{
					Year:     &year3,
					Solution: &solution,
				}),
				testhelpers.NewEstimatedLifecycleCost(testhelpers.EstimatedLifecycleCostOptions{}),
				testhelpers.NewEstimatedLifecycleCost(testhelpers.EstimatedLifecycleCostOptions{
					Year: &year2,
				}),
				testhelpers.NewEstimatedLifecycleCost(testhelpers.EstimatedLifecycleCostOptions{
					Year: &year3,
				}),
				testhelpers.NewEstimatedLifecycleCost(testhelpers.EstimatedLifecycleCostOptions{
					Solution: &solution,
				}),
			},
		}
		_, err := s.store.UpdateBusinessCase(&businessCaseToUpdate)
		s.NoError(err)
		//	fetch the newly updated business case
		updated, err := s.store.FetchBusinessCaseByID(id)
		s.NoError(err)
		s.Equal(7, len(updated.LifecycleCostLines))
	})

	s.Run("doesn't update system intake or eua user id", func() {
		unwantedSystemIntakeID := uuid.New()
		unwantedEUAUserID := testhelpers.RandomEUAID()
		businessCaseToUpdate := models.BusinessCase{
			ID:             id,
			SystemIntakeID: unwantedSystemIntakeID,
			EUAUserID:      unwantedEUAUserID,
		}
		_, err := s.store.UpdateBusinessCase(&businessCaseToUpdate)
		s.NoError(err)
		//	fetch the newly updated business case
		updated, err := s.store.FetchBusinessCaseByID(id)
		s.NoError(err)
		s.NotEqual(unwantedSystemIntakeID, updated.SystemIntakeID)
		s.Equal(intake.ID, updated.SystemIntakeID)
		s.NotEqual(unwantedEUAUserID, updated.EUAUserID)
		s.Equal(euaID, updated.EUAUserID)
	})

	s.Run("fails if the business case ID doesn't exist", func() {
		badUUID := uuid.New()
		businessCaseToUpdate := models.BusinessCase{
			ID:                 badUUID,
			LifecycleCostLines: models.EstimatedLifecycleCosts{},
		}
		_, err := s.store.UpdateBusinessCase(&businessCaseToUpdate)
		s.Error(err)
		s.Equal("business case not found", err.Error())
	})
}
