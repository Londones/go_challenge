package queries

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"go-challenge/internal/database"
	"go-challenge/internal/models"
	"strconv"
	"testing"
	"time"
)

//var testUser = &models.User{
//	Email:    "test@example.com",
//	Name:     "Test User",
//	Password: "password123",
//}
//
//func TestCreateUser(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockService := queries.NewMockUserQueries(ctrl)
//	mockService.EXPECT().CreateUser(testUser).Return(nil)
//
//	err := mockService.CreateUser(testUser)
//	assert.NoError(t, err)
//}

//func TestFindUserByEmail(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockService := queries.NewMockUserQueries(ctrl)
//	mockService.EXPECT().FindUserByEmail("test@example.com").Return(testUser, nil)
//
//	user, err := mockService.FindUserByEmail("test@example.com")
//	assert.NoError(t, err)
//	assert.NotNil(t, user)
//	assert.Equal(t, testUser.Email, user.Email)
//}

var dtb, dbErr = database.TestDatabaseInit()
var db = DatabaseService{s: database.Service{Db: dtb.Db}}

func TestDatabaseService_CreateAnnonce(t *testing.T) {

	if dbErr != nil {
		return
	}

	type test struct {
		name    string
		fields  DatabaseService
		wantId  uint
		wantErr assert.ErrorAssertionFunc
	}

	var annonceTest test
	annonceTest.name = "Test creation annonce"
	annonceTest.fields = db
	annonceTest.wantId = 31
	annonceTest.wantErr = assert.NoError

	s := &DatabaseService{
		s: annonceTest.fields.s,
	}

	users, _ := s.GetAllUsers()
	cats, _ := s.GetAllCats()

	var description string = "annonce de test"

	annonce := models.Annonce{
		Title:       "Annonce de test",
		Description: &description,
		UserID:      users[1].ID,
		CatID:       strconv.Itoa(int(cats[1].ID)),
	}

	t.Run(annonceTest.name, func(t *testing.T) {

		gotId, err := s.CreateAnnonce(&annonce)
		fmt.Println("ici")
		if !annonceTest.wantErr(t, err, fmt.Sprintf("CreateAnnonce(%v)", annonce)) {
			return
		}
		assert.Equalf(t, annonceTest.wantId, gotId, "CreateAnnonce(%v)", annonce)
	})
}

func TestDatabaseService_CreateAssociation(t *testing.T) {
	if dbErr != nil {
		return
	}

	type test struct {
		name    string
		fields  DatabaseService
		wantId  uint
		wantErr assert.ErrorAssertionFunc
	}

	var associationTest test
	associationTest.name = "Test creation association"
	associationTest.fields = db
	associationTest.wantId = 1
	associationTest.wantErr = assert.NoError

	s := &DatabaseService{
		s: associationTest.fields.s,
	}

	var isVerified = true

	users, _ := s.GetAllUsers()

	association := models.Association{
		Name:       "Assoc de TEST",
		AddressRue: "10 rue du TEST",
		Cp:         "12345",
		Ville:      "TEST-CITY",
		Phone:      "0101010101",
		Email:      "emailDeTest@gmail.com",
		OwnerID:    users[1].ID,
		Verified:   &isVerified,
	}

	t.Run(associationTest.name, func(t *testing.T) {

		gotId, err := s.CreateAssociation(&association)
		if !associationTest.wantErr(t, err, fmt.Sprintf("CreateAssociation(%v)", association)) {
			return
		}
		assert.Equalf(t, associationTest.wantId, gotId, "CreateAssociation(%v)", association)
	})
}

func TestDatabaseService_CreateCat(t *testing.T) {
	if dbErr != nil {
		return
	}

	type test struct {
		name    string
		fields  DatabaseService
		wantId  uint
		wantErr assert.ErrorAssertionFunc
	}

	var catTest test
	catTest.name = "Test creation cat"
	catTest.fields = db
	catTest.wantId = 31
	catTest.wantErr = assert.NoError

	s := &DatabaseService{
		s: catTest.fields.s,
	}

	var description string = "Il s'appel PABLO et c'est un chat de TEST"
	var date = time.Time{}

	users, _ := s.GetAllUsers()
	races, _ := s.GetAllRace()

	cat := models.Cats{
		Name:            "PABLO",
		BirthDate:       &date,
		Sexe:            "male",
		LastVaccine:     &date,
		LastVaccineName: "ANTI-PABLO",
		Color:           "VIOLET",
		Behavior:        "PABLOCITO",
		Sterilized:      false,
		RaceID:          strconv.Itoa(int(races[1].ID)),
		Description:     &description,
		Reserved:        false,
		UserID:          users[1].ID,
	}

	t.Run(catTest.name, func(t *testing.T) {

		gotId, err := s.CreateCat(&cat)
		if !catTest.wantErr(t, err, fmt.Sprintf("CreateCat(%v)", cat)) {
			return
		}
		assert.Equalf(t, catTest.wantId, gotId, "CreateCat(%v)", cat)
	})
}

/*
func TestDatabaseService_CreateFavorite(t *testing.T) {
	type fields struct {
		s database.Service
	}
	type args struct {
		favorite *models.Favorite
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			tt.wantErr(t, s.CreateFavorite(tt.args.favorite), fmt.Sprintf("CreateFavorite(%v)", tt.args.favorite))
		})
	}
}

func TestDatabaseService_CreateRace(t *testing.T) {
	type fields struct {
		s database.Service
	}
	type args struct {
		race *models.Races
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			tt.wantErr(t, s.CreateRace(tt.args.race), fmt.Sprintf("CreateRace(%v)", tt.args.race))
		})
	}
}

func TestDatabaseService_CreateRating(t *testing.T) {
	type fields struct {
		s database.Service
	}
	type args struct {
		rating *models.Rating
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			tt.wantErr(t, s.CreateRating(tt.args.rating), fmt.Sprintf("CreateRating(%v)", tt.args.rating))
		})
	}
}

func TestDatabaseService_CreateRoom(t *testing.T) {
	type fields struct {
		s database.Service
	}
	type args struct {
		room *models.Room
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantId  uint
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			gotId, err := s.CreateRoom(tt.args.room)
			if !tt.wantErr(t, err, fmt.Sprintf("CreateRoom(%v)", tt.args.room)) {
				return
			}
			assert.Equalf(t, tt.wantId, gotId, "CreateRoom(%v)", tt.args.room)
		})
	}
}

func TestDatabaseService_CreateUser(t *testing.T) {
	type fields struct {
		s database.Service
	}
	type args struct {
		user *models.User
		role *models.Roles
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			tt.wantErr(t, s.CreateUser(tt.args.user, tt.args.role), fmt.Sprintf("CreateUser(%v, %v)", tt.args.user, tt.args.role))
		})
	}
}

func TestDatabaseService_DeleteAnnonce(t *testing.T) {
	type fields struct {
		s database.Service
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			tt.wantErr(t, s.DeleteAnnonce(tt.args.id), fmt.Sprintf("DeleteAnnonce(%v)", tt.args.id))
		})
	}
}

func TestDatabaseService_DeleteAssociation(t *testing.T) {
	type fields struct {
		s database.Service
	}
	type args struct {
		id int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			tt.wantErr(t, s.DeleteAssociation(tt.args.id), fmt.Sprintf("DeleteAssociation(%v)", tt.args.id))
		})
	}
}

func TestDatabaseService_DeleteCat(t *testing.T) {
	type fields struct {
		s database.Service
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			tt.wantErr(t, s.DeleteCat(tt.args.id), fmt.Sprintf("DeleteCat(%v)", tt.args.id))
		})
	}
}

func TestDatabaseService_DeleteCatByID(t *testing.T) {
	type fields struct {
		s database.Service
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			tt.wantErr(t, s.DeleteCatByID(tt.args.id), fmt.Sprintf("DeleteCatByID(%v)", tt.args.id))
		})
	}
}

func TestDatabaseService_DeleteFavorite(t *testing.T) {
	type fields struct {
		s database.Service
	}
	type args struct {
		favorite *models.Favorite
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			tt.wantErr(t, s.DeleteFavorite(tt.args.favorite), fmt.Sprintf("DeleteFavorite(%v)", tt.args.favorite))
		})
	}
}

func TestDatabaseService_DeleteRace(t *testing.T) {
	type fields struct {
		s database.Service
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			tt.wantErr(t, s.DeleteRace(tt.args.id), fmt.Sprintf("DeleteRace(%v)", tt.args.id))
		})
	}
}

func TestDatabaseService_DeleteRating(t *testing.T) {
	type fields struct {
		s database.Service
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			tt.wantErr(t, s.DeleteRating(tt.args.id), fmt.Sprintf("DeleteRating(%v)", tt.args.id))
		})
	}
}

func TestDatabaseService_DeleteUser(t *testing.T) {
	type fields struct {
		s database.Service
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			tt.wantErr(t, s.DeleteUser(tt.args.id), fmt.Sprintf("DeleteUser(%v)", tt.args.id))
		})
	}
}

func TestDatabaseService_FindAnnonceByCatID(t *testing.T) {
	type fields struct {
		s database.Service
	}
	type args struct {
		catID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Annonce
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			got, err := s.FindAnnonceByCatID(tt.args.catID)
			if !tt.wantErr(t, err, fmt.Sprintf("FindAnnonceByCatID(%v)", tt.args.catID)) {
				return
			}
			assert.Equalf(t, tt.want, got, "FindAnnonceByCatID(%v)", tt.args.catID)
		})
	}
}

func TestDatabaseService_FindAnnonceByID(t *testing.T) {
	type fields struct {
		s database.Service
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Annonce
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			got, err := s.FindAnnonceByID(tt.args.id)
			if !tt.wantErr(t, err, fmt.Sprintf("FindAnnonceByID(%v)", tt.args.id)) {
				return
			}
			assert.Equalf(t, tt.want, got, "FindAnnonceByID(%v)", tt.args.id)
		})
	}
}

func TestDatabaseService_FindAssociationById(t *testing.T) {
	type fields struct {
		s database.Service
	}
	type args struct {
		id int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Association
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			got, err := s.FindAssociationById(tt.args.id)
			if !tt.wantErr(t, err, fmt.Sprintf("FindAssociationById(%v)", tt.args.id)) {
				return
			}
			assert.Equalf(t, tt.want, got, "FindAssociationById(%v)", tt.args.id)
		})
	}
}

func TestDatabaseService_FindAssociationsByUserId(t *testing.T) {
	type fields struct {
		s database.Service
	}
	type args struct {
		userId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.Association
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			got, err := s.FindAssociationsByUserId(tt.args.userId)
			if !tt.wantErr(t, err, fmt.Sprintf("FindAssociationsByUserId(%v)", tt.args.userId)) {
				return
			}
			assert.Equalf(t, tt.want, got, "FindAssociationsByUserId(%v)", tt.args.userId)
		})
	}
}

func TestDatabaseService_FindCatByID(t *testing.T) {
	type fields struct {
		s database.Service
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Cats
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			got, err := s.FindCatByID(tt.args.id)
			if !tt.wantErr(t, err, fmt.Sprintf("FindCatByID(%v)", tt.args.id)) {
				return
			}
			assert.Equalf(t, tt.want, got, "FindCatByID(%v)", tt.args.id)
		})
	}
}

func TestDatabaseService_FindCatsByUserID(t *testing.T) {
	type fields struct {
		s database.Service
	}
	type args struct {
		userID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.Cats
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			got, err := s.FindCatsByUserID(tt.args.userID)
			if !tt.wantErr(t, err, fmt.Sprintf("FindCatsByUserID(%v)", tt.args.userID)) {
				return
			}
			assert.Equalf(t, tt.want, got, "FindCatsByUserID(%v)", tt.args.userID)
		})
	}
}

func TestDatabaseService_FindFavoriteByID(t *testing.T) {
	type fields struct {
		s database.Service
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Favorite
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			got, err := s.FindFavoriteByID(tt.args.id)
			if !tt.wantErr(t, err, fmt.Sprintf("FindFavoriteByID(%v)", tt.args.id)) {
				return
			}
			assert.Equalf(t, tt.want, got, "FindFavoriteByID(%v)", tt.args.id)
		})
	}
}

func TestDatabaseService_FindFavoritesByAnnonceID(t *testing.T) {
	type fields struct {
		s database.Service
	}
	type args struct {
		annonceID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.Favorite
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			got, err := s.FindFavoritesByAnnonceID(tt.args.annonceID)
			if !tt.wantErr(t, err, fmt.Sprintf("FindFavoritesByAnnonceID(%v)", tt.args.annonceID)) {
				return
			}
			assert.Equalf(t, tt.want, got, "FindFavoritesByAnnonceID(%v)", tt.args.annonceID)
		})
	}
}

func TestDatabaseService_FindFavoritesByUserID(t *testing.T) {
	type fields struct {
		s database.Service
	}
	type args struct {
		userID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.Favorite
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			got, err := s.FindFavoritesByUserID(tt.args.userID)
			if !tt.wantErr(t, err, fmt.Sprintf("FindFavoritesByUserID(%v)", tt.args.userID)) {
				return
			}
			assert.Equalf(t, tt.want, got, "FindFavoritesByUserID(%v)", tt.args.userID)
		})
	}
}

func TestDatabaseService_FindRaceByID(t *testing.T) {
	type fields struct {
		s database.Service
	}
	type args struct {
		id string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantRace models.Races
		wantErr  assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			gotRace, err := s.FindRaceByID(tt.args.id)
			if !tt.wantErr(t, err, fmt.Sprintf("FindRaceByID(%v)", tt.args.id)) {
				return
			}
			assert.Equalf(t, tt.wantRace, gotRace, "FindRaceByID(%v)", tt.args.id)
		})
	}
}

func TestDatabaseService_FindRatingByID(t *testing.T) {
	type fields struct {
		s database.Service
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Rating
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			got, err := s.FindRatingByID(tt.args.id)
			if !tt.wantErr(t, err, fmt.Sprintf("FindRatingByID(%v)", tt.args.id)) {
				return
			}
			assert.Equalf(t, tt.want, got, "FindRatingByID(%v)", tt.args.id)
		})
	}
}

func TestDatabaseService_FindRoomsByUserID(t *testing.T) {
	type fields struct {
		s database.Service
	}
	type args struct {
		userid string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*models.Room
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			got, err := s.FindRoomsByUserID(tt.args.userid)
			if !tt.wantErr(t, err, fmt.Sprintf("FindRoomsByUserID(%v)", tt.args.userid)) {
				return
			}
			assert.Equalf(t, tt.want, got, "FindRoomsByUserID(%v)", tt.args.userid)
		})
	}
}

func TestDatabaseService_FindUserByEmail(t *testing.T) {
	type fields struct {
		s database.Service
	}
	type args struct {
		email string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.User
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			got, err := s.FindUserByEmail(tt.args.email)
			if !tt.wantErr(t, err, fmt.Sprintf("FindUserByEmail(%v)", tt.args.email)) {
				return
			}
			assert.Equalf(t, tt.want, got, "FindUserByEmail(%v)", tt.args.email)
		})
	}
}

func TestDatabaseService_FindUserByGoogleID(t *testing.T) {
	type fields struct {
		s database.Service
	}
	type args struct {
		googleID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.User
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			got, err := s.FindUserByGoogleID(tt.args.googleID)
			if !tt.wantErr(t, err, fmt.Sprintf("FindUserByGoogleID(%v)", tt.args.googleID)) {
				return
			}
			assert.Equalf(t, tt.want, got, "FindUserByGoogleID(%v)", tt.args.googleID)
		})
	}
}

func TestDatabaseService_FindUserByID(t *testing.T) {
	type fields struct {
		s database.Service
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.User
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			got, err := s.FindUserByID(tt.args.id)
			if !tt.wantErr(t, err, fmt.Sprintf("FindUserByID(%v)", tt.args.id)) {
				return
			}
			assert.Equalf(t, tt.want, got, "FindUserByID(%v)", tt.args.id)
		})
	}
}

func TestDatabaseService_GetAddressFromAnnonceID(t *testing.T) {
	type fields struct {
		s database.Service
	}
	type args struct {
		userID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			got, err := s.GetAddressFromAnnonceID(tt.args.userID)
			if !tt.wantErr(t, err, fmt.Sprintf("GetAddressFromAnnonceID(%v)", tt.args.userID)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetAddressFromAnnonceID(%v)", tt.args.userID)
		})
	}
}

func TestDatabaseService_GetAllAnnonces(t *testing.T) {
	type fields struct {
		s database.Service
	}
	tests := []struct {
		name    string
		fields  fields
		want    []models.Annonce
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			got, err := s.GetAllAnnonces()
			if !tt.wantErr(t, err, fmt.Sprintf("GetAllAnnonces()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetAllAnnonces()")
		})
	}
}

func TestDatabaseService_GetAllAssociations(t *testing.T) {
	type fields struct {
		s database.Service
	}
	tests := []struct {
		name    string
		fields  fields
		want    []models.Association
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			got, err := s.GetAllAssociations()
			if !tt.wantErr(t, err, fmt.Sprintf("GetAllAssociations()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetAllAssociations()")
		})
	}
}

func TestDatabaseService_GetAllCats(t *testing.T) {
	type fields struct {
		s database.Service
	}
	tests := []struct {
		name    string
		fields  fields
		want    []models.Cats
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			got, err := s.GetAllCats()
			if !tt.wantErr(t, err, fmt.Sprintf("GetAllCats()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetAllCats()")
		})
	}
}

func TestDatabaseService_GetAllRace(t *testing.T) {
	type fields struct {
		s database.Service
	}
	tests := []struct {
		name    string
		fields  fields
		want    []models.Races
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			got, err := s.GetAllRace()
			if !tt.wantErr(t, err, fmt.Sprintf("GetAllRace()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetAllRace()")
		})
	}
}

func TestDatabaseService_GetAllRatings(t *testing.T) {
	type fields struct {
		s database.Service
	}
	tests := []struct {
		name    string
		fields  fields
		want    []models.Rating
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			got, err := s.GetAllRatings()
			if !tt.wantErr(t, err, fmt.Sprintf("GetAllRatings()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetAllRatings()")
		})
	}
}

func TestDatabaseService_GetAllUsers(t *testing.T) {
	type fields struct {
		s database.Service
	}
	tests := []struct {
		name    string
		fields  fields
		want    []models.User
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			got, err := s.GetAllUsers()
			if !tt.wantErr(t, err, fmt.Sprintf("GetAllUsers()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetAllUsers()")
		})
	}
}

func TestDatabaseService_GetAuthorRatings(t *testing.T) {
	type fields struct {
		s database.Service
	}
	type args struct {
		authorID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.Rating
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			got, err := s.GetAuthorRatings(tt.args.authorID)
			if !tt.wantErr(t, err, fmt.Sprintf("GetAuthorRatings(%v)", tt.args.authorID)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetAuthorRatings(%v)", tt.args.authorID)
		})
	}
}

func TestDatabaseService_GetCatByFilters(t *testing.T) {
	type fields struct {
		s database.Service
	}
	type args struct {
		raceId string
		age    int
		sex    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.Cats
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			got, err := s.GetCatByFilters(tt.args.raceId, tt.args.age, tt.args.sex)
			if !tt.wantErr(t, err, fmt.Sprintf("GetCatByFilters(%v, %v, %v)", tt.args.raceId, tt.args.age, tt.args.sex)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetCatByFilters(%v, %v, %v)", tt.args.raceId, tt.args.age, tt.args.sex)
		})
	}
}

func TestDatabaseService_GetMessagesByRoomID(t *testing.T) {
	type fields struct {
		s database.Service
	}
	type args struct {
		roomID uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*models.Message
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			got, err := s.GetMessagesByRoomID(tt.args.roomID)
			if !tt.wantErr(t, err, fmt.Sprintf("GetMessagesByRoomID(%v)", tt.args.roomID)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetMessagesByRoomID(%v)", tt.args.roomID)
		})
	}
}

func TestDatabaseService_GetOrCreateRoom(t *testing.T) {
	type fields struct {
		s database.Service
	}
	type args struct {
		userID1 string
		userID2 string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Room
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			got, err := s.GetOrCreateRoom(tt.args.userID1, tt.args.userID2)
			if !tt.wantErr(t, err, fmt.Sprintf("GetOrCreateRoom(%v, %v)", tt.args.userID1, tt.args.userID2)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetOrCreateRoom(%v, %v)", tt.args.userID1, tt.args.userID2)
		})
	}
}

func TestDatabaseService_GetRoleByName(t *testing.T) {
	type fields struct {
		s database.Service
	}
	type args struct {
		name models.RoleName
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Roles
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			got, err := s.GetRoleByName(tt.args.name)
			if !tt.wantErr(t, err, fmt.Sprintf("GetRoleByName(%v)", tt.args.name)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetRoleByName(%v)", tt.args.name)
		})
	}
}

func TestDatabaseService_GetRoomByID(t *testing.T) {
	type fields struct {
		s database.Service
	}
	type args struct {
		roomID uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Room
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			got, err := s.GetRoomByID(tt.args.roomID)
			if !tt.wantErr(t, err, fmt.Sprintf("GetRoomByID(%v)", tt.args.roomID)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetRoomByID(%v)", tt.args.roomID)
		})
	}
}

func TestDatabaseService_GetRoomIds(t *testing.T) {
	type fields struct {
		s database.Service
	}
	tests := []struct {
		name    string
		fields  fields
		want    []uint
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			got, err := s.GetRoomIds()
			if !tt.wantErr(t, err, fmt.Sprintf("GetRoomIds()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetRoomIds()")
		})
	}
}

func TestDatabaseService_GetRooms(t *testing.T) {
	type fields struct {
		s database.Service
	}
	tests := []struct {
		name    string
		fields  fields
		want    []*models.Room
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			got, err := s.GetRooms()
			if !tt.wantErr(t, err, fmt.Sprintf("GetRooms()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetRooms()")
		})
	}
}

func TestDatabaseService_GetUserAnnonces(t *testing.T) {
	type fields struct {
		s database.Service
	}
	type args struct {
		userID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.Annonce
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			got, err := s.GetUserAnnonces(tt.args.userID)
			if !tt.wantErr(t, err, fmt.Sprintf("GetUserAnnonces(%v)", tt.args.userID)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetUserAnnonces(%v)", tt.args.userID)
		})
	}
}

func TestDatabaseService_GetUserFavorites(t *testing.T) {
	type fields struct {
		s database.Service
	}
	type args struct {
		UserID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.Favorite
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			got, err := s.GetUserFavorites(tt.args.UserID)
			if !tt.wantErr(t, err, fmt.Sprintf("GetUserFavorites(%v)", tt.args.UserID)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetUserFavorites(%v)", tt.args.UserID)
		})
	}
}

func TestDatabaseService_GetUserIDByAnnonceID(t *testing.T) {
	type fields struct {
		s database.Service
	}
	type args struct {
		annonceID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantId  string
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			gotId, err := s.GetUserIDByAnnonceID(tt.args.annonceID)
			if !tt.wantErr(t, err, fmt.Sprintf("GetUserIDByAnnonceID(%v)", tt.args.annonceID)) {
				return
			}
			assert.Equalf(t, tt.wantId, gotId, "GetUserIDByAnnonceID(%v)", tt.args.annonceID)
		})
	}
}

func TestDatabaseService_GetUserRatings(t *testing.T) {
	type fields struct {
		s database.Service
	}
	type args struct {
		userID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.Rating
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			got, err := s.GetUserRatings(tt.args.userID)
			if !tt.wantErr(t, err, fmt.Sprintf("GetUserRatings(%v)", tt.args.userID)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetUserRatings(%v)", tt.args.userID)
		})
	}
}

func TestDatabaseService_SaveMessage(t *testing.T) {
	type fields struct {
		s database.Service
	}
	type args struct {
		roomID   uint
		senderID string
		content  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Message
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			got, err := s.SaveMessage(tt.args.roomID, tt.args.senderID, tt.args.content)
			if !tt.wantErr(t, err, fmt.Sprintf("SaveMessage(%v, %v, %v)", tt.args.roomID, tt.args.senderID, tt.args.content)) {
				return
			}
			assert.Equalf(t, tt.want, got, "SaveMessage(%v, %v, %v)", tt.args.roomID, tt.args.senderID, tt.args.content)
		})
	}
}

func TestDatabaseService_UpdateAnnonce(t *testing.T) {
	type fields struct {
		s database.Service
	}
	type args struct {
		annonce *models.Annonce
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			tt.wantErr(t, s.UpdateAnnonce(tt.args.annonce), fmt.Sprintf("UpdateAnnonce(%v)", tt.args.annonce))
		})
	}
}

func TestDatabaseService_UpdateAnnonceDescription(t *testing.T) {
	type fields struct {
		s database.Service
	}
	type args struct {
		id          string
		description string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Annonce
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			got, err := s.UpdateAnnonceDescription(tt.args.id, tt.args.description)
			if !tt.wantErr(t, err, fmt.Sprintf("UpdateAnnonceDescription(%v, %v)", tt.args.id, tt.args.description)) {
				return
			}
			assert.Equalf(t, tt.want, got, "UpdateAnnonceDescription(%v, %v)", tt.args.id, tt.args.description)
		})
	}
}

func TestDatabaseService_UpdateAssociation(t *testing.T) {
	type fields struct {
		s database.Service
	}
	type args struct {
		association *models.Association
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			tt.wantErr(t, s.UpdateAssociation(tt.args.association), fmt.Sprintf("UpdateAssociation(%v)", tt.args.association))
		})
	}
}

func TestDatabaseService_UpdateCat(t *testing.T) {
	type fields struct {
		s database.Service
	}
	type args struct {
		cat *models.Cats
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			tt.wantErr(t, s.UpdateCat(tt.args.cat), fmt.Sprintf("UpdateCat(%v)", tt.args.cat))
		})
	}
}

func TestDatabaseService_UpdateFavorite(t *testing.T) {
	type fields struct {
		s database.Service
	}
	type args struct {
		favorite *models.Favorite
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			tt.wantErr(t, s.UpdateFavorite(tt.args.favorite), fmt.Sprintf("UpdateFavorite(%v)", tt.args.favorite))
		})
	}
}

func TestDatabaseService_UpdateRace(t *testing.T) {
	type fields struct {
		s database.Service
	}
	type args struct {
		race models.Races
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			tt.wantErr(t, s.UpdateRace(tt.args.race), fmt.Sprintf("UpdateRace(%v)", tt.args.race))
		})
	}
}

func TestDatabaseService_UpdateRating(t *testing.T) {
	type fields struct {
		s database.Service
	}
	type args struct {
		rating *models.Rating
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			tt.wantErr(t, s.UpdateRating(tt.args.rating), fmt.Sprintf("UpdateRating(%v)", tt.args.rating))
		})
	}
}

func TestDatabaseService_UpdateUser(t *testing.T) {
	type fields struct {
		s database.Service
	}
	type args struct {
		user *models.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &DatabaseService{
				s: tt.fields.s,
			}
			tt.wantErr(t, s.UpdateUser(tt.args.user), fmt.Sprintf("UpdateUser(%v)", tt.args.user))
		})
	}
}

func TestNewQueriesService(t *testing.T) {
	type args struct {
		s *database.Service
	}
	tests := []struct {
		name string
		args args
		want *DatabaseService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewQueriesService(tt.args.s), "NewQueriesService(%v)", tt.args.s)
		})
	}
} */
