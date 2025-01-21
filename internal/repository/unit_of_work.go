package repository

import (
	"gorm.io/gorm"
)

// unitOfWork implements UnitOfWork interface
type unitOfWork struct {
	db            *gorm.DB
	userRepo      UserRepository
	communityRepo CommunityRepository
	buildingRepo  BuildingRepository
	unitRepo      UnitRepository
	tx            *gorm.DB
}

// NewUnitOfWork creates a new UnitOfWork instance
func NewUnitOfWork(db *gorm.DB) UnitOfWork {
	return &unitOfWork{
		db: db,
	}
}

// Users returns the UserRepository, initializing it if necessary
func (u *unitOfWork) Users() UserRepository {
	if u.userRepo == nil {
		if u.tx != nil {
			u.userRepo = NewUserRepository(u.tx)
		} else {
			u.userRepo = NewUserRepository(u.db)
		}
	}
	return u.userRepo
}

// Communities returns the CommunityRepository, initializing it if necessary
func (u *unitOfWork) Communities() CommunityRepository {
	if u.communityRepo == nil {
		if u.tx != nil {
			u.communityRepo = NewCommunityRepository(u.tx)
		} else {
			u.communityRepo = NewCommunityRepository(u.db)
		}
	}
	return u.communityRepo
}

// Units returns the UnitRepository, initializing it if necessary
func (u *unitOfWork) Units() UnitRepository {
	if u.unitRepo == nil {
		if u.tx != nil {
			u.unitRepo = NewUnitRepository(u.tx)
		} else {
			u.unitRepo = NewUnitRepository(u.db)
		}
	}
	return u.unitRepo
}

// Buildings returns the BuildingRepository, initializing it if necessary
func (u *unitOfWork) Buildings() BuildingRepository {
	if u.buildingRepo == nil {
		if u.tx != nil {
			u.buildingRepo = NewBuildingRepository(u.tx)
		} else {
			u.buildingRepo = NewBuildingRepository(u.db)
		}
	}
	return u.buildingRepo
}

// Begin starts a new transaction
func (u *unitOfWork) Begin() error {
	tx := u.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	u.tx = tx
	return nil
}

// Commit commits the current transaction
func (u *unitOfWork) Commit() error {
	if u.tx == nil {
		return nil
	}
	err := u.tx.Commit().Error
	u.tx = nil
	u.userRepo = nil
	u.communityRepo = nil
	u.buildingRepo = nil
	u.unitRepo = nil
	return err
}

// Rollback rolls back the current transaction
func (u *unitOfWork) Rollback() error {
	if u.tx == nil {
		return nil
	}
	err := u.tx.Rollback().Error
	u.tx = nil
	u.userRepo = nil
	u.communityRepo = nil
	u.buildingRepo = nil
	u.unitRepo = nil
	return err
}
