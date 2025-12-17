package items

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/robebs/ts-se-tool-go/internal/save/dataformat"
)

// Economy mirrors Economy in C# CustomClasses/Save/Items/Economy.cs
type Economy struct {
	Bank                           string
	Player                         string
	Companies                      []string
	Garages                        []string
	GarageIgnoreList               []string
	GameProgress                   string
	EventQueue                     string
	MailCtrl                       string
	OversizeOfferCtrl              string
	GameTime                       uint32
	GameTimeSecs                   dataformat.Float
	GameTimeInitial                int
	AchievementsAdded              int
	NewGame                        bool
	TotalDistance                  int
	ExperiencePoints               uint32
	Adr                            byte
	LongDist                       byte
	Heavy                          byte
	Fragile                        byte
	Urgent                         byte
	Mechanical                     byte
	UserColors                     []dataformat.Color
	DeliveryLog                    string
	FerryLog                       string
	StoredCameraMode               int
	StoredActorState               int
	StoredHighBeamStyle            int
	StoredActorWindowsState        dataformat.Vector2f
	StoredActorWiperMode           int
	StoredActorRetarder            int
	StoredDisplayMode              int
	StoredDashboardMapMode         int
	StoredWorldMapZoom             int
	StoredOnlineJobID              int
	StoredOnlineGPSBehind          []string
	StoredOnlineGPSAhead           []string
	StoredOnlineGPSBehindWaypoints []string
	StoredOnlineGPSAheadWaypoints  []string
	StoredOnlineGPSAvoidWaypoints  []string
	StoredSpecialJob               string
	PoliceCtrl                     string
	StoredMapState                 int
	StoredGasPumpMoney             int
	StoredWeatherChangeTimer       dataformat.Float
	StoredCurrentWeather           int
	StoredRainWetness              dataformat.Float
	TimeZone                       int
	TimeZoneName                   string
	LastFerryPosition              dataformat.Vector3i
	StoredShowWeigh                bool
	StoredNeedToWeigh              bool
	StoredNavStartPos              dataformat.Vector3i
	StoredNavEndPos                dataformat.Vector3i
	StoredGPSBehind                []string
	StoredGPSAhead                 []string
	StoredGPSBehindWaypoints       []string
	StoredGPSAheadWaypoints        []string
	StoredGPSAvoidWaypoints        []string
	StoredStartTollgatePos         dataformat.Vector3i
	StoredTutorialState            int
	StoredMapActions               []string
	CleanDistanceCounter           int
	CleanDistanceMax               int
	NoCargoDamageDistanceCounter   int
	NoCargoDamageDistanceMax       int
	NoViolationDistanceCounter     int
	NoViolationDistanceMax         int
	TotalRealTime                  int
	RealTimeSeconds                dataformat.Float
	VisitedCities                  []string
	VisitedCitiesCount             []int
	LastVisitedCity                string
	DiscoveredCutsceneItems        []uint64
	DiscoveredCutsceneItemsStates  []int
	UnlockedDealers                []string
	UnlockedRecruitments           []string
	TotalScreenshotCount           int
	UndamagedCargoRow              int
	ServiceVisitCount              int
	LastServicePos                 dataformat.Vector3f
	GasStationVisitCount           int
	LastGasStationPos              dataformat.Vector3f
	EmergencyCallCount             int
	AICrashCount                   int
	TruckColorChangeCount          int
	RedLightFineCount              int
	CancelledJobCount              int
	TotalFuelLitres                int
	TotalFuelPrice                 int
	TransportedCargoTypes          []string
	AchievedFeats                  int
	DiscoveredRoads                int
	DiscoveredItems                []uint64
	DriversOffer                   []string
	FreelanceTruckOffer            string
	TrucksBoughtOnline             int
	SpecialCargoTimer              int
	ScreenAccessList               []string
	DriverPool                     []string
	Registry                       string
	CompanyJobsInvitationSent      bool
	CompanyCheckHash               uint64
	Relations                      []int
	BusStops                       []string
	BusJobLog                      string
	BusExperiencePoints            int
	BusTotalDistance               int
	BusFinishedJobCount            int
	BusCancelledJobCount           int
	BusTotalPassengers             int
	BusTotalStops                  int
	BusGameTime                    int
	BusPlayingTime                 int
}

// FromProperties fills the Economy struct from a map of properties.
func (e *Economy) FromProperties(props map[string][]string) error {
	for key, vals := range props {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]

		switch {
		case key == "bank":
			e.Bank = val
		case key == "player":
			e.Player = val
		case key == "companies":
			// capacity hint ignored
		case strings.HasPrefix(key, "companies["):
			e.Companies = append(e.Companies, val)
		case key == "garages":
			// capacity hint ignored
		case strings.HasPrefix(key, "garages["):
			e.Garages = append(e.Garages, val)
		case key == "garage_ignore_list":
			// capacity hint ignored
		case strings.HasPrefix(key, "garage_ignore_list["):
			e.GarageIgnoreList = append(e.GarageIgnoreList, val)
		case key == "game_progress":
			e.GameProgress = val
		case key == "event_queue":
			e.EventQueue = val
		case key == "mail_ctrl":
			e.MailCtrl = val
		case key == "oversize_offer_ctrl":
			e.OversizeOfferCtrl = val
		case key == "game_time":
			v, _ := strconv.ParseUint(val, 10, 32)
			e.GameTime = uint32(v)
		case key == "game_time_secs":
			e.GameTimeSecs = parseEconomyFloat(val)
		case key == "game_time_initial":
			e.GameTimeInitial = parseEconomyInt(val)
		case key == "achievements_added":
			e.AchievementsAdded = parseEconomyInt(val)
		case key == "new_game":
			e.NewGame = parseEconomyBool(val)
		case key == "total_distance":
			e.TotalDistance = parseEconomyInt(val)
		case key == "experience_points":
			v, _ := strconv.ParseUint(val, 10, 32)
			e.ExperiencePoints = uint32(v)
		case key == "adr":
			v, _ := strconv.ParseUint(val, 10, 8)
			e.Adr = byte(v)
		case key == "long_dist":
			v, _ := strconv.ParseUint(val, 10, 8)
			e.LongDist = byte(v)
		case key == "heavy":
			v, _ := strconv.ParseUint(val, 10, 8)
			e.Heavy = byte(v)
		case key == "fragile":
			v, _ := strconv.ParseUint(val, 10, 8)
			e.Fragile = byte(v)
		case key == "urgent":
			v, _ := strconv.ParseUint(val, 10, 8)
			e.Urgent = byte(v)
		case key == "mechanical":
			v, _ := strconv.ParseUint(val, 10, 8)
			e.Mechanical = byte(v)
		case key == "user_colors":
			// capacity hint ignored
		case strings.HasPrefix(key, "user_colors["):
			color, err := parseColor(val)
			if err != nil {
				return fmt.Errorf("parse user_colors: %w", err)
			}
			e.UserColors = append(e.UserColors, color)
		case key == "delivery_log":
			e.DeliveryLog = val
		case key == "ferry_log":
			e.FerryLog = val
		case key == "stored_camera_mode":
			e.StoredCameraMode = parseEconomyInt(val)
		case key == "stored_actor_state":
			e.StoredActorState = parseEconomyInt(val)
		case key == "stored_high_beam_style":
			e.StoredHighBeamStyle = parseEconomyInt(val)
		case key == "stored_actor_windows_state":
			vec, err := parseVector2f(val)
			if err != nil {
				return fmt.Errorf("parse stored_actor_windows_state: %w", err)
			}
			e.StoredActorWindowsState = vec
		case key == "stored_actor_wiper_mode":
			e.StoredActorWiperMode = parseEconomyInt(val)
		case key == "stored_actor_retarder":
			e.StoredActorRetarder = parseEconomyInt(val)
		case key == "stored_display_mode":
			e.StoredDisplayMode = parseEconomyInt(val)
		case key == "stored_dashboard_map_mode":
			e.StoredDashboardMapMode = parseEconomyInt(val)
		case key == "stored_world_map_zoom":
			e.StoredWorldMapZoom = parseEconomyInt(val)
		case key == "stored_online_job_id":
			e.StoredOnlineJobID = parseEconomyInt(val)
		case key == "stored_online_gps_behind":
			// capacity hint ignored
		case strings.HasPrefix(key, "stored_online_gps_behind["):
			e.StoredOnlineGPSBehind = append(e.StoredOnlineGPSBehind, val)
		case key == "stored_online_gps_ahead":
			// capacity hint ignored
		case strings.HasPrefix(key, "stored_online_gps_ahead["):
			e.StoredOnlineGPSAhead = append(e.StoredOnlineGPSAhead, val)
		case key == "stored_online_gps_behind_waypoints":
			// capacity hint ignored
		case strings.HasPrefix(key, "stored_online_gps_behind_waypoints["):
			e.StoredOnlineGPSBehindWaypoints = append(e.StoredOnlineGPSBehindWaypoints, val)
		case key == "stored_online_gps_ahead_waypoints":
			// capacity hint ignored
		case strings.HasPrefix(key, "stored_online_gps_ahead_waypoints["):
			e.StoredOnlineGPSAheadWaypoints = append(e.StoredOnlineGPSAheadWaypoints, val)
		case key == "stored_online_gps_avoid_waypoints":
			// capacity hint ignored
		case strings.HasPrefix(key, "stored_online_gps_avoid_waypoints["):
			e.StoredOnlineGPSAvoidWaypoints = append(e.StoredOnlineGPSAvoidWaypoints, val)
		case key == "stored_special_job":
			e.StoredSpecialJob = val
		case key == "police_ctrl":
			e.PoliceCtrl = val
		case key == "stored_map_state":
			e.StoredMapState = parseEconomyInt(val)
		case key == "stored_gas_pump_money":
			e.StoredGasPumpMoney = parseEconomyInt(val)
		case key == "stored_weather_change_timer":
			e.StoredWeatherChangeTimer = parseEconomyFloat(val)
		case key == "stored_current_weather":
			e.StoredCurrentWeather = parseEconomyInt(val)
		case key == "stored_rain_wetness":
			e.StoredRainWetness = parseEconomyFloat(val)
		case key == "time_zone":
			e.TimeZone = parseEconomyInt(val)
		case key == "time_zone_name":
			e.TimeZoneName = val
		case key == "last_ferry_position":
			vec, err := parseVector3i(val)
			if err != nil {
				return fmt.Errorf("parse last_ferry_position: %w", err)
			}
			e.LastFerryPosition = vec
		case key == "stored_show_weigh":
			e.StoredShowWeigh = parseEconomyBool(val)
		case key == "stored_need_to_weigh":
			e.StoredNeedToWeigh = parseEconomyBool(val)
		case key == "stored_nav_start_pos":
			vec, err := parseVector3i(val)
			if err != nil {
				return fmt.Errorf("parse stored_nav_start_pos: %w", err)
			}
			e.StoredNavStartPos = vec
		case key == "stored_nav_end_pos":
			vec, err := parseVector3i(val)
			if err != nil {
				return fmt.Errorf("parse stored_nav_end_pos: %w", err)
			}
			e.StoredNavEndPos = vec
		case key == "stored_gps_behind":
			// capacity hint ignored
		case strings.HasPrefix(key, "stored_gps_behind["):
			e.StoredGPSBehind = append(e.StoredGPSBehind, val)
		case key == "stored_gps_ahead":
			// capacity hint ignored
		case strings.HasPrefix(key, "stored_gps_ahead["):
			e.StoredGPSAhead = append(e.StoredGPSAhead, val)
		case key == "stored_gps_behind_waypoints":
			// capacity hint ignored
		case strings.HasPrefix(key, "stored_gps_behind_waypoints["):
			e.StoredGPSBehindWaypoints = append(e.StoredGPSBehindWaypoints, val)
		case key == "stored_gps_ahead_waypoints":
			// capacity hint ignored
		case strings.HasPrefix(key, "stored_gps_ahead_waypoints["):
			e.StoredGPSAheadWaypoints = append(e.StoredGPSAheadWaypoints, val)
		case key == "stored_gps_avoid_waypoints":
			// capacity hint ignored
		case strings.HasPrefix(key, "stored_gps_avoid_waypoints["):
			e.StoredGPSAvoidWaypoints = append(e.StoredGPSAvoidWaypoints, val)
		case key == "stored_start_tollgate_pos":
			vec, err := parseVector3i(val)
			if err != nil {
				return fmt.Errorf("parse stored_start_tollgate_pos: %w", err)
			}
			e.StoredStartTollgatePos = vec
		case key == "stored_tutorial_state":
			e.StoredTutorialState = parseEconomyInt(val)
		case key == "stored_map_actions":
			// capacity hint ignored
		case strings.HasPrefix(key, "stored_map_actions["):
			e.StoredMapActions = append(e.StoredMapActions, val)
		case key == "clean_distance_counter":
			e.CleanDistanceCounter = parseEconomyInt(val)
		case key == "clean_distance_max":
			e.CleanDistanceMax = parseEconomyInt(val)
		case key == "no_cargo_damage_distance_counter":
			e.NoCargoDamageDistanceCounter = parseEconomyInt(val)
		case key == "no_cargo_damage_distance_max":
			e.NoCargoDamageDistanceMax = parseEconomyInt(val)
		case key == "no_violation_distance_counter":
			e.NoViolationDistanceCounter = parseEconomyInt(val)
		case key == "no_violation_distance_max":
			e.NoViolationDistanceMax = parseEconomyInt(val)
		case key == "total_real_time":
			e.TotalRealTime = parseEconomyInt(val)
		case key == "real_time_seconds":
			e.RealTimeSeconds = parseEconomyFloat(val)
		case key == "visited_cities":
			// capacity hint ignored
		case strings.HasPrefix(key, "visited_cities["):
			e.VisitedCities = append(e.VisitedCities, val)
		case key == "visited_cities_count":
			// capacity hint ignored
		case strings.HasPrefix(key, "visited_cities_count["):
			e.VisitedCitiesCount = append(e.VisitedCitiesCount, parseEconomyInt(val))
		case key == "last_visited_city":
			e.LastVisitedCity = val
		case key == "discovered_cutscene_items":
			// capacity hint ignored
		case strings.HasPrefix(key, "discovered_cutscene_items["):
			v, _ := strconv.ParseUint(val, 10, 64)
			e.DiscoveredCutsceneItems = append(e.DiscoveredCutsceneItems, v)
		case key == "discovered_cutscene_items_states":
			// capacity hint ignored
		case strings.HasPrefix(key, "discovered_cutscene_items_states["):
			e.DiscoveredCutsceneItemsStates = append(e.DiscoveredCutsceneItemsStates, parseEconomyInt(val))
		case key == "unlocked_dealers":
			// capacity hint ignored
		case strings.HasPrefix(key, "unlocked_dealers["):
			e.UnlockedDealers = append(e.UnlockedDealers, val)
		case key == "unlocked_recruitments":
			// capacity hint ignored
		case strings.HasPrefix(key, "unlocked_recruitments["):
			e.UnlockedRecruitments = append(e.UnlockedRecruitments, val)
		case key == "total_screeshot_count":
			e.TotalScreenshotCount = parseEconomyInt(val)
		case key == "undamaged_cargo_row":
			e.UndamagedCargoRow = parseEconomyInt(val)
		case key == "service_visit_count":
			e.ServiceVisitCount = parseEconomyInt(val)
		case key == "last_service_pos":
			vec, err := parseVector3f(val)
			if err != nil {
				return fmt.Errorf("parse last_service_pos: %w", err)
			}
			e.LastServicePos = vec
		case key == "gas_station_visit_count":
			e.GasStationVisitCount = parseEconomyInt(val)
		case key == "last_gas_station_pos":
			vec, err := parseVector3f(val)
			if err != nil {
				return fmt.Errorf("parse last_gas_station_pos: %w", err)
			}
			e.LastGasStationPos = vec
		case key == "emergency_call_count":
			e.EmergencyCallCount = parseEconomyInt(val)
		case key == "ai_crash_count":
			e.AICrashCount = parseEconomyInt(val)
		case key == "truck_color_change_count":
			e.TruckColorChangeCount = parseEconomyInt(val)
		case key == "red_light_fine_count":
			e.RedLightFineCount = parseEconomyInt(val)
		case key == "cancelled_job_count":
			e.CancelledJobCount = parseEconomyInt(val)
		case key == "total_fuel_litres":
			e.TotalFuelLitres = parseEconomyInt(val)
		case key == "total_fuel_price":
			e.TotalFuelPrice = parseEconomyInt(val)
		case key == "transported_cargo_types":
			// capacity hint ignored
		case strings.HasPrefix(key, "transported_cargo_types["):
			e.TransportedCargoTypes = append(e.TransportedCargoTypes, val)
		case key == "achieved_feats":
			e.AchievedFeats = parseEconomyInt(val)
		case key == "discovered_roads":
			e.DiscoveredRoads = parseEconomyInt(val)
		case key == "discovered_items":
			// capacity hint ignored
		case strings.HasPrefix(key, "discovered_items["):
			v, _ := strconv.ParseUint(val, 10, 64)
			e.DiscoveredItems = append(e.DiscoveredItems, v)
		case key == "drivers_offer":
			// capacity hint ignored
		case strings.HasPrefix(key, "drivers_offer["):
			e.DriversOffer = append(e.DriversOffer, val)
		case key == "freelance_truck_offer":
			e.FreelanceTruckOffer = val
		case key == "trucks_bought_online":
			e.TrucksBoughtOnline = parseEconomyInt(val)
		case key == "special_cargo_timer":
			e.SpecialCargoTimer = parseEconomyInt(val)
		case key == "screen_access_list":
			// capacity hint ignored
		case strings.HasPrefix(key, "screen_access_list["):
			e.ScreenAccessList = append(e.ScreenAccessList, val)
		case key == "driver_pool":
			// capacity hint ignored
		case strings.HasPrefix(key, "driver_pool["):
			e.DriverPool = append(e.DriverPool, val)
		case key == "registry":
			e.Registry = val
		case key == "company_jobs_invitation_sent":
			e.CompanyJobsInvitationSent = parseEconomyBool(val)
		case key == "company_check_hash":
			v, _ := strconv.ParseUint(val, 10, 64)
			e.CompanyCheckHash = v
		case key == "relations":
			// capacity hint ignored
		case strings.HasPrefix(key, "relations["):
			e.Relations = append(e.Relations, parseEconomyInt(val))
		case key == "bus_stops":
			// capacity hint ignored
		case strings.HasPrefix(key, "bus_stops["):
			e.BusStops = append(e.BusStops, val)
		case key == "bus_job_log":
			e.BusJobLog = val
		case key == "bus_experience_points":
			e.BusExperiencePoints = parseEconomyInt(val)
		case key == "bus_total_distance":
			e.BusTotalDistance = parseEconomyInt(val)
		case key == "bus_finished_job_count":
			e.BusFinishedJobCount = parseEconomyInt(val)
		case key == "bus_cancelled_job_count":
			e.BusCancelledJobCount = parseEconomyInt(val)
		case key == "bus_total_passengers":
			e.BusTotalPassengers = parseEconomyInt(val)
		case key == "bus_total_stops":
			e.BusTotalStops = parseEconomyInt(val)
		case key == "bus_game_time":
			e.BusGameTime = parseEconomyInt(val)
		case key == "bus_playing_time":
			e.BusPlayingTime = parseEconomyInt(val)
		}
	}
	return nil
}

// ToProperties converts the Economy struct to a map of properties.
func (e *Economy) ToProperties() map[string][]string {
	props := make(map[string][]string)

	props["bank"] = []string{e.Bank}
	props["player"] = []string{e.Player}

	props["companies"] = []string{strconv.Itoa(len(e.Companies))}
	for i, v := range e.Companies {
		props[fmt.Sprintf("companies[%d]", i)] = []string{v}
	}

	props["garages"] = []string{strconv.Itoa(len(e.Garages))}
	for i, v := range e.Garages {
		props[fmt.Sprintf("garages[%d]", i)] = []string{v}
	}

	props["garage_ignore_list"] = []string{strconv.Itoa(len(e.GarageIgnoreList))}
	for i, v := range e.GarageIgnoreList {
		props[fmt.Sprintf("garage_ignore_list[%d]", i)] = []string{v}
	}

	props["game_progress"] = []string{e.GameProgress}
	props["event_queue"] = []string{e.EventQueue}
	props["mail_ctrl"] = []string{e.MailCtrl}
	props["oversize_offer_ctrl"] = []string{e.OversizeOfferCtrl}
	props["game_time"] = []string{strconv.FormatUint(uint64(e.GameTime), 10)}
	props["game_time_secs"] = []string{formatEconomyFloat(e.GameTimeSecs)}
	props["game_time_initial"] = []string{strconv.Itoa(e.GameTimeInitial)}
	props["achievements_added"] = []string{strconv.Itoa(e.AchievementsAdded)}
	props["new_game"] = []string{formatEconomyBool(e.NewGame)}
	props["total_distance"] = []string{strconv.Itoa(e.TotalDistance)}
	props["experience_points"] = []string{strconv.FormatUint(uint64(e.ExperiencePoints), 10)}
	props["adr"] = []string{strconv.FormatUint(uint64(e.Adr), 10)}
	props["long_dist"] = []string{strconv.FormatUint(uint64(e.LongDist), 10)}
	props["heavy"] = []string{strconv.FormatUint(uint64(e.Heavy), 10)}
	props["fragile"] = []string{strconv.FormatUint(uint64(e.Fragile), 10)}
	props["urgent"] = []string{strconv.FormatUint(uint64(e.Urgent), 10)}
	props["mechanical"] = []string{strconv.FormatUint(uint64(e.Mechanical), 10)}

	props["user_colors"] = []string{strconv.Itoa(len(e.UserColors))}
	for i, v := range e.UserColors {
		props[fmt.Sprintf("user_colors[%d]", i)] = []string{formatColor(v)}
	}

	props["delivery_log"] = []string{e.DeliveryLog}
	props["ferry_log"] = []string{e.FerryLog}
	props["stored_camera_mode"] = []string{strconv.Itoa(e.StoredCameraMode)}
	props["stored_actor_state"] = []string{strconv.Itoa(e.StoredActorState)}
	props["stored_high_beam_style"] = []string{strconv.Itoa(e.StoredHighBeamStyle)}
	props["stored_actor_windows_state"] = []string{formatVector2f(e.StoredActorWindowsState)}
	props["stored_actor_wiper_mode"] = []string{strconv.Itoa(e.StoredActorWiperMode)}
	props["stored_actor_retarder"] = []string{strconv.Itoa(e.StoredActorRetarder)}
	props["stored_display_mode"] = []string{strconv.Itoa(e.StoredDisplayMode)}
	props["stored_dashboard_map_mode"] = []string{strconv.Itoa(e.StoredDashboardMapMode)}
	props["stored_world_map_zoom"] = []string{strconv.Itoa(e.StoredWorldMapZoom)}
	props["stored_online_job_id"] = []string{strconv.Itoa(e.StoredOnlineJobID)}

	props["stored_online_gps_behind"] = []string{strconv.Itoa(len(e.StoredOnlineGPSBehind))}
	for i, v := range e.StoredOnlineGPSBehind {
		props[fmt.Sprintf("stored_online_gps_behind[%d]", i)] = []string{v}
	}

	props["stored_online_gps_ahead"] = []string{strconv.Itoa(len(e.StoredOnlineGPSAhead))}
	for i, v := range e.StoredOnlineGPSAhead {
		props[fmt.Sprintf("stored_online_gps_ahead[%d]", i)] = []string{v}
	}

	props["stored_online_gps_behind_waypoints"] = []string{strconv.Itoa(len(e.StoredOnlineGPSBehindWaypoints))}
	for i, v := range e.StoredOnlineGPSBehindWaypoints {
		props[fmt.Sprintf("stored_online_gps_behind_waypoints[%d]", i)] = []string{v}
	}

	props["stored_online_gps_ahead_waypoints"] = []string{strconv.Itoa(len(e.StoredOnlineGPSAheadWaypoints))}
	for i, v := range e.StoredOnlineGPSAheadWaypoints {
		props[fmt.Sprintf("stored_online_gps_ahead_waypoints[%d]", i)] = []string{v}
	}

	props["stored_online_gps_avoid_waypoints"] = []string{strconv.Itoa(len(e.StoredOnlineGPSAvoidWaypoints))}
	for i, v := range e.StoredOnlineGPSAvoidWaypoints {
		props[fmt.Sprintf("stored_online_gps_avoid_waypoints[%d]", i)] = []string{v}
	}

	props["stored_special_job"] = []string{e.StoredSpecialJob}
	props["police_ctrl"] = []string{e.PoliceCtrl}
	props["stored_map_state"] = []string{strconv.Itoa(e.StoredMapState)}
	props["stored_gas_pump_money"] = []string{strconv.Itoa(e.StoredGasPumpMoney)}
	props["stored_weather_change_timer"] = []string{formatEconomyFloat(e.StoredWeatherChangeTimer)}
	props["stored_current_weather"] = []string{strconv.Itoa(e.StoredCurrentWeather)}
	props["stored_rain_wetness"] = []string{formatEconomyFloat(e.StoredRainWetness)}
	props["time_zone"] = []string{strconv.Itoa(e.TimeZone)}
	props["time_zone_name"] = []string{e.TimeZoneName}
	props["last_ferry_position"] = []string{formatVector3i(e.LastFerryPosition)}
	props["stored_show_weigh"] = []string{formatEconomyBool(e.StoredShowWeigh)}
	props["stored_need_to_weigh"] = []string{formatEconomyBool(e.StoredNeedToWeigh)}
	props["stored_nav_start_pos"] = []string{formatVector3i(e.StoredNavStartPos)}
	props["stored_nav_end_pos"] = []string{formatVector3i(e.StoredNavEndPos)}

	props["stored_gps_behind"] = []string{strconv.Itoa(len(e.StoredGPSBehind))}
	for i, v := range e.StoredGPSBehind {
		props[fmt.Sprintf("stored_gps_behind[%d]", i)] = []string{v}
	}

	props["stored_gps_ahead"] = []string{strconv.Itoa(len(e.StoredGPSAhead))}
	for i, v := range e.StoredGPSAhead {
		props[fmt.Sprintf("stored_gps_ahead[%d]", i)] = []string{v}
	}

	props["stored_gps_behind_waypoints"] = []string{strconv.Itoa(len(e.StoredGPSBehindWaypoints))}
	for i, v := range e.StoredGPSBehindWaypoints {
		props[fmt.Sprintf("stored_gps_behind_waypoints[%d]", i)] = []string{v}
	}

	props["stored_gps_ahead_waypoints"] = []string{strconv.Itoa(len(e.StoredGPSAheadWaypoints))}
	for i, v := range e.StoredGPSAheadWaypoints {
		props[fmt.Sprintf("stored_gps_ahead_waypoints[%d]", i)] = []string{v}
	}

	props["stored_gps_avoid_waypoints"] = []string{strconv.Itoa(len(e.StoredGPSAvoidWaypoints))}
	for i, v := range e.StoredGPSAvoidWaypoints {
		props[fmt.Sprintf("stored_gps_avoid_waypoints[%d]", i)] = []string{v}
	}

	props["stored_start_tollgate_pos"] = []string{formatVector3i(e.StoredStartTollgatePos)}
	props["stored_tutorial_state"] = []string{strconv.Itoa(e.StoredTutorialState)}

	props["stored_map_actions"] = []string{strconv.Itoa(len(e.StoredMapActions))}
	for i, v := range e.StoredMapActions {
		props[fmt.Sprintf("stored_map_actions[%d]", i)] = []string{v}
	}

	props["clean_distance_counter"] = []string{strconv.Itoa(e.CleanDistanceCounter)}
	props["clean_distance_max"] = []string{strconv.Itoa(e.CleanDistanceMax)}
	props["no_cargo_damage_distance_counter"] = []string{strconv.Itoa(e.NoCargoDamageDistanceCounter)}
	props["no_cargo_damage_distance_max"] = []string{strconv.Itoa(e.NoCargoDamageDistanceMax)}
	props["no_violation_distance_counter"] = []string{strconv.Itoa(e.NoViolationDistanceCounter)}
	props["no_violation_distance_max"] = []string{strconv.Itoa(e.NoViolationDistanceMax)}
	props["total_real_time"] = []string{strconv.Itoa(e.TotalRealTime)}
	props["real_time_seconds"] = []string{formatEconomyFloat(e.RealTimeSeconds)}

	props["visited_cities"] = []string{strconv.Itoa(len(e.VisitedCities))}
	for i, v := range e.VisitedCities {
		props[fmt.Sprintf("visited_cities[%d]", i)] = []string{v}
	}

	props["visited_cities_count"] = []string{strconv.Itoa(len(e.VisitedCitiesCount))}
	for i, v := range e.VisitedCitiesCount {
		props[fmt.Sprintf("visited_cities_count[%d]", i)] = []string{strconv.Itoa(v)}
	}

	props["last_visited_city"] = []string{e.LastVisitedCity}

	props["discovered_cutscene_items"] = []string{strconv.Itoa(len(e.DiscoveredCutsceneItems))}
	for i, v := range e.DiscoveredCutsceneItems {
		props[fmt.Sprintf("discovered_cutscene_items[%d]", i)] = []string{strconv.FormatUint(v, 10)}
	}

	props["discovered_cutscene_items_states"] = []string{strconv.Itoa(len(e.DiscoveredCutsceneItemsStates))}
	for i, v := range e.DiscoveredCutsceneItemsStates {
		props[fmt.Sprintf("discovered_cutscene_items_states[%d]", i)] = []string{strconv.Itoa(v)}
	}

	props["unlocked_dealers"] = []string{strconv.Itoa(len(e.UnlockedDealers))}
	for i, v := range e.UnlockedDealers {
		props[fmt.Sprintf("unlocked_dealers[%d]", i)] = []string{v}
	}

	props["unlocked_recruitments"] = []string{strconv.Itoa(len(e.UnlockedRecruitments))}
	for i, v := range e.UnlockedRecruitments {
		props[fmt.Sprintf("unlocked_recruitments[%d]", i)] = []string{v}
	}

	props["total_screeshot_count"] = []string{strconv.Itoa(e.TotalScreenshotCount)}
	props["undamaged_cargo_row"] = []string{strconv.Itoa(e.UndamagedCargoRow)}
	props["service_visit_count"] = []string{strconv.Itoa(e.ServiceVisitCount)}
	props["last_service_pos"] = []string{formatVector3f(e.LastServicePos)}
	props["gas_station_visit_count"] = []string{strconv.Itoa(e.GasStationVisitCount)}
	props["last_gas_station_pos"] = []string{formatVector3f(e.LastGasStationPos)}
	props["emergency_call_count"] = []string{strconv.Itoa(e.EmergencyCallCount)}
	props["ai_crash_count"] = []string{strconv.Itoa(e.AICrashCount)}
	props["truck_color_change_count"] = []string{strconv.Itoa(e.TruckColorChangeCount)}
	props["red_light_fine_count"] = []string{strconv.Itoa(e.RedLightFineCount)}
	props["cancelled_job_count"] = []string{strconv.Itoa(e.CancelledJobCount)}
	props["total_fuel_litres"] = []string{strconv.Itoa(e.TotalFuelLitres)}
	props["total_fuel_price"] = []string{strconv.Itoa(e.TotalFuelPrice)}

	props["transported_cargo_types"] = []string{strconv.Itoa(len(e.TransportedCargoTypes))}
	for i, v := range e.TransportedCargoTypes {
		props[fmt.Sprintf("transported_cargo_types[%d]", i)] = []string{v}
	}

	props["achieved_feats"] = []string{strconv.Itoa(e.AchievedFeats)}
	props["discovered_roads"] = []string{strconv.Itoa(e.DiscoveredRoads)}

	props["discovered_items"] = []string{strconv.Itoa(len(e.DiscoveredItems))}
	for i, v := range e.DiscoveredItems {
		props[fmt.Sprintf("discovered_items[%d]", i)] = []string{strconv.FormatUint(v, 10)}
	}

	props["drivers_offer"] = []string{strconv.Itoa(len(e.DriversOffer))}
	for i, v := range e.DriversOffer {
		props[fmt.Sprintf("drivers_offer[%d]", i)] = []string{v}
	}

	props["freelance_truck_offer"] = []string{e.FreelanceTruckOffer}
	props["trucks_bought_online"] = []string{strconv.Itoa(e.TrucksBoughtOnline)}
	props["special_cargo_timer"] = []string{strconv.Itoa(e.SpecialCargoTimer)}

	props["screen_access_list"] = []string{strconv.Itoa(len(e.ScreenAccessList))}
	for i, v := range e.ScreenAccessList {
		props[fmt.Sprintf("screen_access_list[%d]", i)] = []string{v}
	}

	props["driver_pool"] = []string{strconv.Itoa(len(e.DriverPool))}
	for i, v := range e.DriverPool {
		props[fmt.Sprintf("driver_pool[%d]", i)] = []string{v}
	}

	props["registry"] = []string{e.Registry}
	props["company_jobs_invitation_sent"] = []string{formatEconomyBool(e.CompanyJobsInvitationSent)}
	props["company_check_hash"] = []string{strconv.FormatUint(e.CompanyCheckHash, 10)}

	props["relations"] = []string{strconv.Itoa(len(e.Relations))}
	for i, v := range e.Relations {
		props[fmt.Sprintf("relations[%d]", i)] = []string{strconv.Itoa(v)}
	}

	props["bus_stops"] = []string{strconv.Itoa(len(e.BusStops))}
	for i, v := range e.BusStops {
		props[fmt.Sprintf("bus_stops[%d]", i)] = []string{v}
	}

	props["bus_job_log"] = []string{e.BusJobLog}
	props["bus_experience_points"] = []string{strconv.Itoa(e.BusExperiencePoints)}
	props["bus_total_distance"] = []string{strconv.Itoa(e.BusTotalDistance)}
	props["bus_finished_job_count"] = []string{strconv.Itoa(e.BusFinishedJobCount)}
	props["bus_cancelled_job_count"] = []string{strconv.Itoa(e.BusCancelledJobCount)}
	props["bus_total_passengers"] = []string{strconv.Itoa(e.BusTotalPassengers)}
	props["bus_total_stops"] = []string{strconv.Itoa(e.BusTotalStops)}
	props["bus_game_time"] = []string{strconv.Itoa(e.BusGameTime)}
	props["bus_playing_time"] = []string{strconv.Itoa(e.BusPlayingTime)}

	return props
}

// Helper functions for parsing and formatting

func parseEconomyInt(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}

func parseEconomyFloat(s string) dataformat.Float {
	f, _ := strconv.ParseFloat(s, 32)
	return dataformat.Float(f)
}

func formatEconomyFloat(f dataformat.Float) string {
	return strconv.FormatFloat(float64(f), 'f', -1, 32)
}

func parseEconomyBool(s string) bool {
	return strings.ToLower(s) == "true"
}

func formatEconomyBool(b bool) string {
	return strconv.FormatBool(b)
}

// parseVector2f parses a Vector2f from a string like "(&bd85bf17, &bd5ecfd4)"
func parseVector2f(s string) (dataformat.Vector2f, error) {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "(")
	s = strings.TrimSuffix(s, ")")
	parts := strings.Split(s, ",")
	if len(parts) != 2 {
		return dataformat.Vector2f{}, fmt.Errorf("invalid vector2f format: expected 2 components, got %d", len(parts))
	}

	x, err := parseHexFloat(strings.TrimSpace(parts[0]))
	if err != nil {
		return dataformat.Vector2f{}, fmt.Errorf("parse X: %w", err)
	}
	y, err := parseHexFloat(strings.TrimSpace(parts[1]))
	if err != nil {
		return dataformat.Vector2f{}, fmt.Errorf("parse Y: %w", err)
	}

	return dataformat.Vector2f{X: x, Y: y}, nil
}

// formatVector2f formats a Vector2f to string like "(&bd85bf17, &bd5ecfd4)"
func formatVector2f(v dataformat.Vector2f) string {
	return fmt.Sprintf("(%s, %s)", formatHexFloat(v.X), formatHexFloat(v.Y))
}

// parseVector3i parses a Vector3i from a string like "(2147483647, 2147483647, 2147483647)"
func parseVector3i(s string) (dataformat.Vector3i, error) {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "(")
	s = strings.TrimSuffix(s, ")")
	parts := strings.Split(s, ",")
	if len(parts) != 3 {
		return dataformat.Vector3i{}, fmt.Errorf("invalid vector3i format: expected 3 components, got %d", len(parts))
	}

	x, err := strconv.ParseInt(strings.TrimSpace(parts[0]), 10, 32)
	if err != nil {
		return dataformat.Vector3i{}, fmt.Errorf("parse X: %w", err)
	}
	y, err := strconv.ParseInt(strings.TrimSpace(parts[1]), 10, 32)
	if err != nil {
		return dataformat.Vector3i{}, fmt.Errorf("parse Y: %w", err)
	}
	z, err := strconv.ParseInt(strings.TrimSpace(parts[2]), 10, 32)
	if err != nil {
		return dataformat.Vector3i{}, fmt.Errorf("parse Z: %w", err)
	}

	return dataformat.Vector3i{X: int32(x), Y: int32(y), Z: int32(z)}, nil
}

// formatVector3i formats a Vector3i to string like "(2147483647, 2147483647, 2147483647)"
func formatVector3i(v dataformat.Vector3i) string {
	return fmt.Sprintf("(%d, %d, %d)", v.X, v.Y, v.Z)
}

// parseColor parses a Color from a string (uses dataformat.Color parsing)
func parseColor(s string) (dataformat.Color, error) {
	return dataformat.NewColorFromString(s), nil
}

// formatColor formats a Color to string
func formatColor(c dataformat.Color) string {
	return c.ToString()
}

// parseVector3f, formatVector3f, parseHexFloat, formatHexFloat are defined in company.go
// They are reused here since they're in the same package
