package items

// This file defines lightweight Go structs mirroring the C# classes in
// CustomClasses/Save/Items. The fields and behavior will be filled in
// progressively as the SII parser/serializer is completed.

// SiiNBlockCore is the base type for many save-game blocks in the C# code.
// Here it is just a marker that can be embedded or referenced later.
type SiiNBlockCore struct{}

type BusJobLog struct{}

type BusStop struct{}

type Company struct{}

type DeliveryLogEntry struct{}

type DeliveryLog struct{}

type DriverAI struct{}

type DriverPlayer struct{}

type EconomyEventQueue struct{}

type EconomyEvent struct{}

type Economy struct{}

type FerryLogEntry struct{}

type FerryLog struct{}

type GameProgress struct{}

type Garage struct{}

type GPSWaypointStorage struct{}

type JobInfo struct{}

type JobOfferData struct{}

type MailCtrl struct{}

type MailDef struct{}

type MapAction struct{}

type OversizeBlockRuleSave struct{}

type OversizeJobSave struct{}

type OversizeOfferCtrl struct{}

type OversizeOffer struct{}

type OversizeRouteOffers struct{}

type PlayerJob struct{}

type PoliceCtrl struct{}

type ProfitLogEntry struct{}

type ProfitLog struct{}

type Registry struct{}

type SiiNunit struct{}

type TrailerDef struct{}

type TrailerUtilizationLogEntry struct{}

type TrailerUtilizationLog struct{}

type Trailer struct{}

type TrajectoryOrdersSave struct{}

type TransportData struct{}

type Unidentified struct{}

type VehicleAccessory struct{}

type VehicleAddonAccessory struct{}

type VehicleCargoAccessory struct{}

type VehicleDrvPlateAccessory struct{}

type VehiclePaintJobAccessory struct{}

type VehicleSoundAccessory struct{}

type VehicleWheelAccessory struct{}

type Vehicle struct{}
