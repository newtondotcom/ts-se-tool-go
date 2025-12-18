package itemsextra

// This package mirrors classes from CustomClasses/Save/ItemsExtra in the
// original C# TS SE Tool code base. The types defined here are *domain*
// types, built on top of the low-level save-game items found in
// internal/save/items (Economy, Company, Job_offer_Data, Garage, Player, etc.).
//
// The goal of these structs is not to be a 1‑to‑1 port of every C# field but
// rather to provide the minimal, convenient view required by higher-level
// logic such as "prepare cities/companies", "prepare cargos/trailers/trucks",
// and the future CLI features.

// Cargo roughly mirrors CustomClasses/Save/ItemsExtra/Cargo.cs.
// It represents a single logical cargo (e.g. "furniture") with some metadata.
type Cargo struct {
	// ID is the internal cargo identifier extracted from strings like
	// "cargo.furniture" -> "furniture".
	ID string

	// CargoType encodes the type of cargo:
	//   0 = normal
	//   1 = heavy
	//   2 = double
	// (mirrors the C# logic around company_truck containing "heavy"/"double").
	CargoType int

	// TrailerDefName is the SII name of the trailer_definition block.
	TrailerDefName string

	// UnitsCount is the amount/units requested in the job_offer, when known.
	UnitsCount int
}

// City roughly mirrors CustomClasses/Save/ItemsExtra/City.cs.
type City struct {
	// Name is the internal city name, e.g. "bruxelles".
	Name string

	// Country is the internal country identifier (filled via CountryDictionary).
	Country string

	// Companies holds the list of company identifiers located in this city
	// (e.g. "ikea", "sellplan").
	Companies []string

	// JobOffersCountByCompany stores, when available, the number of job offers
	// per company in this city.
	JobOffersCountByCompany map[string]int
}

// Company roughly mirrors CustomClasses/Save/ItemsExtra/Company.cs.
type Company struct {
	// Name is the company identifier (e.g. "ikea").
	Name string

	// CityName is the name of the city where the company is located.
	CityName string

	// JobOffers is the list of SII names for job_offer_data blocks related
	// to this company.
	JobOffers []string
}

// CompanyTruck roughly mirrors CustomClasses/Save/ItemsExtra/CompanyTruck.cs.
type CompanyTruck struct {
	// TruckID is the string stored in job_offer_data.company_truck, such as
	// "truck.volvo.fh16.2013".
	TruckID string

	// CargoType follows the same convention as Cargo.CargoType.
	CargoType int
}

// Country roughly mirrors CustomClasses/Save/ItemsExtra/Country.cs.
type Country struct {
	// Name is the internal country identifier (e.g. "belgium").
	Name string

	// Cities lists the names of cities belonging to this country.
	Cities []string
}

// Garages roughly mirrors CustomClasses/Save/ItemsExtra/Garages.cs.
// It is a light wrapper over the low-level items.Garage data.
type Garages struct {
	// Name is the internal garage name (part after "garage."), often tied to
	// a city identifier.
	Name string

	// Status is copied from the Garage item (e.g. owned / for sale / etc.).
	Status int

	// Vehicles, Drivers and Trailers contain the raw SII names as stored in the
	// garage item.
	Vehicles []string
	Drivers  []string
	Trailers []string
}

// PlayerJob is a small helper for jobs currently associated to the player.
// It can be extended later if we need more data (reward, cargo, company, ...).
type PlayerJob struct {
	// JobOfferName is the SII name of the job_offer_data block.
	JobOfferName string

	// CompanyName and CityName describe where the job originates.
	CompanyName string
	CityName    string
}

// TrailerDefinition roughly mirrors CustomClasses/Save/ItemsExtra/TrailerDefinition.cs.
type TrailerDefinition struct {
	// Name is the SII name of the trailer_definition block.
	Name string

	// Variants collects all trailer_variant identifiers seen for this definition.
	Variants []string
}

// UserCompanyDriverData is a placeholder for future driver-related logic.
// For the world-loading step we mostly need trucks/trailers, so keep this
// intentionally small for now.
type UserCompanyDriverData struct {
	// DriverName is the SII name of the driver (ai_driver.*).
	DriverName string
}

// UserCompanyTrailerData represents a trailer owned by the player company.
type UserCompanyTrailerData struct {
	// TrailerName is the SII name of the trailer block.
	TrailerName string

	// DefName is the associated trailer_def name, when known.
	DefName string
}

// UserCompanyTruckData represents a truck owned by the player company.
type UserCompanyTruckData struct {
	// TruckName is the SII name of the truck block.
	TruckName string

	// CurrentGarage is the name of the garage where the truck is parked,
	// if this information is available.
	CurrentGarage string
}

// UserCompanyTruckDataPart is kept as a placeholder to mirror the original
// C# design. It can later host split data such as accessories, paint jobs, etc.
type UserCompanyTruckDataPart struct {
	TruckName string
}

// VisitedCity mirrors CustomClasses/Save/ItemsExtra/VisitedCity.cs in a
// condensed form.
type VisitedCity struct {
	CityName string

	// Count is the number of times the city has been visited, when known.
	Count int
}
