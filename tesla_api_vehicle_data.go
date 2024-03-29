package tesla

// autogenerated from https://mholt.github.io/json-to-go/
// copied json file there

type TeslaVehicleData struct {
	Response struct {
		ID                     int64    `json:"id"`
		UserID                 int      `json:"user_id"`
		VehicleID              int      `json:"vehicle_id"`
		Vin                    string   `json:"vin"`
		DisplayName            string   `json:"display_name"`
		Color                  any      `json:"color"`
		AccessType             string   `json:"access_type"`
		Tokens                 []string `json:"tokens"`
		State                  string   `json:"state"`
		InService              bool     `json:"in_service"`
		IDS                    string   `json:"id_s"`
		CalendarEnabled        bool     `json:"calendar_enabled"`
		APIVersion             int      `json:"api_version"`
		BackseatToken          any      `json:"backseat_token"`
		BackseatTokenUpdatedAt any      `json:"backseat_token_updated_at"`
		DriveState             struct {
			GpsAsOf                 int     `json:"gps_as_of"`
			Heading                 int     `json:"heading"`
			Latitude                float64 `json:"latitude"`
			Longitude               float64 `json:"longitude"`
			NativeLatitude          float64 `json:"native_latitude"`
			NativeLocationSupported int     `json:"native_location_supported"`
			NativeLongitude         float64 `json:"native_longitude"`
			NativeType              string  `json:"native_type"`
			Power                   int     `json:"power"`
			ShiftState              any     `json:"shift_state"`
			Speed                   any     `json:"speed"`
			Timestamp               int64   `json:"timestamp"`
		} `json:"drive_state"`
		ClimateState struct {
			BatteryHeater              bool    `json:"battery_heater"`
			BatteryHeaterNoPower       bool    `json:"battery_heater_no_power"`
			ClimateKeeperMode          string  `json:"climate_keeper_mode"`
			DefrostMode                int     `json:"defrost_mode"`
			DriverTempSetting          float64 `json:"driver_temp_setting"`
			FanStatus                  int     `json:"fan_status"`
			InsideTemp                 float64 `json:"inside_temp"`
			IsAutoConditioningOn       bool    `json:"is_auto_conditioning_on"`
			IsClimateOn                bool    `json:"is_climate_on"`
			IsFrontDefrosterOn         bool    `json:"is_front_defroster_on"`
			IsPreconditioning          bool    `json:"is_preconditioning"`
			IsRearDefrosterOn          bool    `json:"is_rear_defroster_on"`
			LeftTempDirection          int     `json:"left_temp_direction"`
			MaxAvailTemp               float64 `json:"max_avail_temp"`
			MinAvailTemp               float64 `json:"min_avail_temp"`
			OutsideTemp                float64 `json:"outside_temp"`
			PassengerTempSetting       float64 `json:"passenger_temp_setting"`
			RemoteHeaterControlEnabled bool    `json:"remote_heater_control_enabled"`
			RightTempDirection         int     `json:"right_temp_direction"`
			SeatHeaterLeft             int     `json:"seat_heater_left"`
			SeatHeaterRight            int     `json:"seat_heater_right"`
			SideMirrorHeaters          bool    `json:"side_mirror_heaters"`
			Timestamp                  int64   `json:"timestamp"`
			WiperBladeHeater           bool    `json:"wiper_blade_heater"`
		} `json:"climate_state"`
		ChargeState struct {
			BatteryHeaterOn             bool    `json:"battery_heater_on"`
			BatteryLevel                int     `json:"battery_level"`
			BatteryRange                float64 `json:"battery_range"`
			ChargeCurrentRequest        int     `json:"charge_current_request"`
			ChargeCurrentRequestMax     int     `json:"charge_current_request_max"`
			ChargeEnableRequest         bool    `json:"charge_enable_request"`
			ChargeEnergyAdded           float64 `json:"charge_energy_added"`
			ChargeLimitSoc              int     `json:"charge_limit_soc"`
			ChargeLimitSocMax           int     `json:"charge_limit_soc_max"`
			ChargeLimitSocMin           int     `json:"charge_limit_soc_min"`
			ChargeLimitSocStd           int     `json:"charge_limit_soc_std"`
			ChargeMilesAddedIdeal       float64 `json:"charge_miles_added_ideal"`
			ChargeMilesAddedRated       float64 `json:"charge_miles_added_rated"`
			ChargePortColdWeatherMode   any     `json:"charge_port_cold_weather_mode"`
			ChargePortDoorOpen          bool    `json:"charge_port_door_open"`
			ChargePortLatch             string  `json:"charge_port_latch"`
			ChargeRate                  float64 `json:"charge_rate"`
			ChargeToMaxRange            bool    `json:"charge_to_max_range"`
			ChargerActualCurrent        int     `json:"charger_actual_current"`
			ChargerPhases               int     `json:"charger_phases"`
			ChargerPilotCurrent         int     `json:"charger_pilot_current"`
			ChargerPower                int     `json:"charger_power"`
			ChargerVoltage              int     `json:"charger_voltage"`
			ChargingState               string  `json:"charging_state"`
			ConnChargeCable             string  `json:"conn_charge_cable"`
			EstBatteryRange             float64 `json:"est_battery_range"`
			FastChargerBrand            string  `json:"fast_charger_brand"`
			FastChargerPresent          bool    `json:"fast_charger_present"`
			FastChargerType             string  `json:"fast_charger_type"`
			IdealBatteryRange           float64 `json:"ideal_battery_range"`
			ManagedChargingActive       bool    `json:"managed_charging_active"`
			ManagedChargingStartTime    any     `json:"managed_charging_start_time"`
			ManagedChargingUserCanceled bool    `json:"managed_charging_user_canceled"`
			MaxRangeChargeCounter       int     `json:"max_range_charge_counter"`
			MinutesToFullCharge         int     `json:"minutes_to_full_charge"`
			NotEnoughPowerToHeat        bool    `json:"not_enough_power_to_heat"`
			ScheduledChargingPending    bool    `json:"scheduled_charging_pending"`
			ScheduledChargingStartTime  any     `json:"scheduled_charging_start_time"`
			TimeToFullCharge            float64 `json:"time_to_full_charge"`
			Timestamp                   int64   `json:"timestamp"`
			TripCharging                bool    `json:"trip_charging"`
			UsableBatteryLevel          int     `json:"usable_battery_level"`
			UserChargeEnableRequest     any     `json:"user_charge_enable_request"`
		} `json:"charge_state"`
		GuiSettings struct {
			Gui24HourTime       bool   `json:"gui_24_hour_time"`
			GuiChargeRateUnits  string `json:"gui_charge_rate_units"`
			GuiDistanceUnits    string `json:"gui_distance_units"`
			GuiRangeDisplay     string `json:"gui_range_display"`
			GuiTemperatureUnits string `json:"gui_temperature_units"`
			ShowRangeUnits      bool   `json:"show_range_units"`
			Timestamp           int64  `json:"timestamp"`
		} `json:"gui_settings"`
		VehicleState struct {
			APIVersion          int    `json:"api_version"`
			AutoparkStateV2     string `json:"autopark_state_v2"`
			AutoparkStyle       string `json:"autopark_style"`
			CalendarSupported   bool   `json:"calendar_supported"`
			CarVersion          string `json:"car_version"`
			CenterDisplayState  int    `json:"center_display_state"`
			Df                  int    `json:"df"`
			Dr                  int    `json:"dr"`
			FdWindow            int    `json:"fd_window"`
			FpWindow            int    `json:"fp_window"`
			Ft                  int    `json:"ft"`
			HomelinkDeviceCount int    `json:"homelink_device_count"`
			HomelinkNearby      bool   `json:"homelink_nearby"`
			IsUserPresent       bool   `json:"is_user_present"`
			LastAutoparkError   string `json:"last_autopark_error"`
			Locked              bool   `json:"locked"`
			MediaState          struct {
				RemoteControlEnabled bool `json:"remote_control_enabled"`
			} `json:"media_state"`
			NotificationsSupported  bool    `json:"notifications_supported"`
			Odometer                float64 `json:"odometer"`
			ParsedCalendarSupported bool    `json:"parsed_calendar_supported"`
			Pf                      int     `json:"pf"`
			Pr                      int     `json:"pr"`
			RdWindow                int     `json:"rd_window"`
			RemoteStart             bool    `json:"remote_start"`
			RemoteStartEnabled      bool    `json:"remote_start_enabled"`
			RemoteStartSupported    bool    `json:"remote_start_supported"`
			RpWindow                int     `json:"rp_window"`
			Rt                      int     `json:"rt"`
			SentryMode              bool    `json:"sentry_mode"`
			SentryModeAvailable     bool    `json:"sentry_mode_available"`
			SmartSummonAvailable    bool    `json:"smart_summon_available"`
			SoftwareUpdate          struct {
				DownloadPerc        int    `json:"download_perc"`
				ExpectedDurationSec int    `json:"expected_duration_sec"`
				InstallPerc         int    `json:"install_perc"`
				Status              string `json:"status"`
				Version             string `json:"version"`
			} `json:"software_update"`
			SpeedLimitMode struct {
				Active          bool    `json:"active"`
				CurrentLimitMph float64 `json:"current_limit_mph"`
				MaxLimitMph     int     `json:"max_limit_mph"`
				MinLimitMph     int     `json:"min_limit_mph"`
				PinCodeSet      bool    `json:"pin_code_set"`
			} `json:"speed_limit_mode"`
			SummonStandbyModeEnabled bool   `json:"summon_standby_mode_enabled"`
			SunRoofPercentOpen       int    `json:"sun_roof_percent_open"`
			SunRoofState             string `json:"sun_roof_state"`
			Timestamp                int64  `json:"timestamp"`
			ValetMode                bool   `json:"valet_mode"`
			ValetPinNeeded           bool   `json:"valet_pin_needed"`
			VehicleName              any    `json:"vehicle_name"`
		} `json:"vehicle_state"`
		VehicleConfig struct {
			CanAcceptNavigationRequests bool   `json:"can_accept_navigation_requests"`
			CanActuateTrunks            bool   `json:"can_actuate_trunks"`
			CarSpecialType              string `json:"car_special_type"`
			CarType                     string `json:"car_type"`
			ChargePortType              string `json:"charge_port_type"`
			DefaultChargeToMax          bool   `json:"default_charge_to_max"`
			EceRestrictions             bool   `json:"ece_restrictions"`
			EuVehicle                   bool   `json:"eu_vehicle"`
			ExteriorColor               string `json:"exterior_color"`
			HasAirSuspension            bool   `json:"has_air_suspension"`
			HasLudicrousMode            bool   `json:"has_ludicrous_mode"`
			MotorizedChargePort         bool   `json:"motorized_charge_port"`
			Plg                         bool   `json:"plg"`
			RearSeatHeaters             int    `json:"rear_seat_heaters"`
			RearSeatType                int    `json:"rear_seat_type"`
			Rhd                         bool   `json:"rhd"`
			RoofColor                   string `json:"roof_color"`
			SeatType                    int    `json:"seat_type"`
			SpoilerType                 string `json:"spoiler_type"`
			SunRoofInstalled            int    `json:"sun_roof_installed"`
			ThirdRowSeats               string `json:"third_row_seats"`
			Timestamp                   int64  `json:"timestamp"`
			TrimBadging                 string `json:"trim_badging"`
			UseRangeBadging             bool   `json:"use_range_badging"`
			WheelType                   string `json:"wheel_type"`
		} `json:"vehicle_config"`
	} `json:"response"`
}
