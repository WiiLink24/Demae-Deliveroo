package deliveroo

type CategoryCode string

const (
	Pizza            CategoryCode = "01"
	BentoBox         CategoryCode = "02"
	Sushi            CategoryCode = "03"
	Japanese         CategoryCode = "04"
	Chinese          CategoryCode = "05"
	Western          CategoryCode = "06"
	FastFood         CategoryCode = "07"
	Curry            CategoryCode = "08"
	PartyFood        CategoryCode = "09"
	DrinksAndDessert CategoryCode = "10"
)

// HasFoodType stores the availability of food types
type HasFoodType struct {
	Pizza            bool
	Bento            bool
	Sushi            bool
	Japanese         bool
	Chinese          bool
	Western          bool
	FastFood         bool
	Curry            bool
	PartyFood        bool
	DessertAndDrinks bool
}

type Store struct {
	Name         string
	StoreID      string
	Address      string
	WaitTime     float64
	MinPrice     string
	IsOpen       bool
	DetailedWait string
	Phone        string
	ServiceHours ServiceHours
	Information  string
	Amenity      string
}

type ServiceHours struct {
	OpenTime  string
	CloseTime string
}

type Category struct {
	Name string
	Code int
}

type Item struct {
	Name           string
	Description    string
	ImgID          string
	Price          string
	SoldOut        bool
	ModifierGroups []ModifierGroup
}

type ModifierGroup struct {
	ID           string
	Name         string
	MinSelection float64
	MaxSelection float64
	Modifiers    []Modifier
}

type Modifier struct {
	ID          string
	Name        string
	Description string
	Price       string
}
