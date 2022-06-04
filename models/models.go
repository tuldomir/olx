package models

import "time"

// Appartment .
type Appartment struct {
	ID              int       `json:"id"`
	Notified        bool      `json:"notified"`
	URL             string    `json:"url"`
	Title           string    `json:"title"`
	LastRefreshTime time.Time `json:"last_refresh_time"`
	CreatedTime     time.Time `json:"created_time"`
	ValidToTime     time.Time `json:"valid_to_time"`
	PushupTime      time.Time `json:"pushup_time"`
	// Description     string    `json:"description"`
	Promotion struct {
		Highlighted   bool     `json:"highlighted"`
		Urgent        bool     `json:"urgent"`
		TopAd         bool     `json:"top_ad"`
		Options       []string `json:"options"`
		B2CAdPage     bool     `json:"b2c_ad_page"`
		PremiumAdPage bool     `json:"premium_ad_page"`
	} `json:"promotion"`
	Params []struct {
		Key   string `json:"key"`
		Name  string `json:"name"`
		Type  string `json:"type"`
		Value struct {
			Key   interface{} `json:"key"`
			Label string      `json:"label"`
		} `json:"value"`
	} `json:"params"`
	KeyParams []interface{} `json:"key_params"`
	Business  bool          `json:"business"`
	// User      struct {
	// 	ID                       int         `json:"id"`
	// 	Created                  time.Time   `json:"created"`
	// 	OtherAdsEnabled          bool        `json:"other_ads_enabled"`
	// 	Name                     string      `json:"name"`
	// 	Logo                     interface{} `json:"logo"`
	// 	LogoAdPage               interface{} `json:"logo_ad_page"`
	// 	SocialNetworkAccountType interface{} `json:"social_network_account_type"`
	// 	Photo                    interface{} `json:"photo"`
	// 	BannerMobile             string      `json:"banner_mobile"`
	// 	BannerDesktop            string      `json:"banner_desktop"`
	// 	CompanyName              string      `json:"company_name"`
	// 	About                    string      `json:"about"`
	// 	B2CBusinessPage          bool        `json:"b2c_business_page"`
	// 	IsOnline                 bool        `json:"is_online"`
	// 	LastSeen                 time.Time   `json:"last_seen"`
	// 	SellerType               interface{} `json:"seller_type"`
	// } `json:"user"`
	Status string `json:"status"`
	// Contact struct {
	// 	Name        string `json:"name"`
	// 	Phone       bool   `json:"phone"`
	// 	Chat        bool   `json:"chat"`
	// 	Negotiation bool   `json:"negotiation"`
	// 	Courier     bool   `json:"courier"`
	// } `json:"contact"`
	Map struct {
		Zoom         int     `json:"zoom"`
		Lat          float64 `json:"lat"`
		Lon          float64 `json:"lon"`
		Radius       int     `json:"radius"`
		ShowDetailed bool    `json:"show_detailed"`
	} `json:"map"`
	// Location struct {
	// 	City struct {
	// 		ID             int    `json:"id"`
	// 		Name           string `json:"name"`
	// 		NormalizedName string `json:"normalized_name"`
	// 	} `json:"city"`
	// 	District struct {
	// 		ID   int    `json:"id"`
	// 		Name string `json:"name"`
	// 	} `json:"district"`
	// 	Region struct {
	// 		ID             int    `json:"id"`
	// 		Name           string `json:"name"`
	// 		NormalizedName string `json:"normalized_name"`
	// 	} `json:"region"`
	// } `json:"location"`
	// Photos []struct {
	// 	ID       int    `json:"id"`
	// 	Filename string `json:"filename"`
	// 	Rotation int    `json:"rotation"`
	// 	Width    int    `json:"width"`
	// 	Height   int    `json:"height"`
	// 	Link     string `json:"link"`
	// } `json:"photos"`
	// Partner  interface{} `json:"partner"`
	// Category struct {
	// 	ID   int    `json:"id"`
	// 	Type string `json:"type"`
	// } `json:"category"`
	// Delivery struct {
	// 	Rock struct {
	// 		OfferID interface{} `json:"offer_id"`
	// 		Active  bool        `json:"active"`
	// 		Mode    interface{} `json:"mode"`
	// 	} `json:"rock"`
	// } `json:"delivery"`
	// Safedeal struct {
	// 	Weight          int           `json:"weight"`
	// 	WeightGrams     int           `json:"weight_grams"`
	// 	Status          string        `json:"status"`
	// 	SafedealBlocked bool          `json:"safedeal_blocked"`
	// 	AllowedQuantity []interface{} `json:"allowed_quantity"`
	// } `json:"safedeal"`
	// Shop struct {
	// 	Subdomain interface{} `json:"subdomain"`
	// } `json:"shop"`
	OfferType string `json:"offer_type"`
}

// QueryResult .
type QueryResult struct {
	Data     []*Appartment `json:"data"`
	Metadata struct {
		TotalElements     int    `json:"total_elements"`
		VisibleTotalCount int    `json:"visible_total_count"`
		Promoted          []int  `json:"promoted"`
		SearchID          string `json:"search_id"`
		// Adverts           struct {
		// 	Places interface{} `json:"places"`
		// 	Config struct {
		// 		Targeting struct {
		// 			Env               string        `json:"env"`
		// 			Lang              string        `json:"lang"`
		// 			Account           string        `json:"account"`
		// 			DfpUserID         string        `json:"dfp_user_id"`
		// 			CatL0ID           string        `json:"cat_l0_id"`
		// 			CatL1ID           string        `json:"cat_l1_id"`
		// 			CatL0             string        `json:"cat_l0"`
		// 			CatL0Path         string        `json:"cat_l0_path"`
		// 			CatL1             string        `json:"cat_l1"`
		// 			CatL1Path         string        `json:"cat_l1_path"`
		// 			CatL2             string        `json:"cat_l2"`
		// 			CatL0Name         string        `json:"cat_l0_name"`
		// 			CatL1Name         string        `json:"cat_l1_name"`
		// 			CatID             string        `json:"cat_id"`
		// 			PrivateBusiness   string        `json:"private_business"`
		// 			OfferSeek         string        `json:"offer_seek"`
		// 			View              string        `json:"view"`
		// 			RegionID          string        `json:"region_id"`
		// 			Region            string        `json:"region"`
		// 			Subregion         string        `json:"subregion"`
		// 			CityID            string        `json:"city_id"`
		// 			City              string        `json:"city"`
		// 			Currency          string        `json:"currency"`
		// 			SearchEngineInput string        `json:"search_engine_input"`
		// 			Page              string        `json:"page"`
		// 			Segment           []interface{} `json:"segment"`
		// 			AppVersion        string        `json:"app_version"`
		// 		} `json:"targeting"`
		// 	} `json:"config"`
		// } `json:"adverts"`
		Source struct {
			Promoted []int `json:"promoted"`
			Organic  []int `json:"organic"`
		} `json:"source"`
	} `json:"metadata"`
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		Next struct {
			Href string `json:"href"`
		} `json:"next"`
		First struct {
			Href string `json:"href"`
		} `json:"first"`
	} `json:"links"`
}

// TeleMsg .
type TeleMsg struct {
	URL             string    `json:"url"`
	LastRefreshTime time.Time `json:"lastrefreshtime"`
}

// Result .
type Result struct {
	Data *QueryResult
	Err  error
}

// Msgs .
type Msgs struct {
	Msgs []*TeleMsg
	Err  error
}

// BotMsgReq .
type BotMsgReq struct {
	ChatID    int64  `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

// SendMessage .
type SendMessage struct {
	Ok     bool `json:"ok"`
	Result struct {
		MessageID int `json:"message_id"`
		From      struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name"`
			Username  string `json:"username"`
		} `json:"from"`
		Chat struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Username  string `json:"username"`
			Type      string `json:"type"`
		} `json:"chat"`
		Date int    `json:"date"`
		Text string `json:"text"`
	} `json:"result"`
}
