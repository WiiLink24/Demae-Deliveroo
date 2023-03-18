package deliveroo

import (
	"encoding/json"
	"github.com/google/uuid"
)

// GetShopsQuery returns a GraphQL query we need to send to Deliveroo
func GetShopsQuery(longitude float64, latitude float64) ([]byte, error) {
	_uuid, _ := uuid.NewUUID()
	theMap := map[string]any{
		"operationName": "Home",
		"query":         "query Home($location: LocationInput!, $options: SearchOptionsInput!, $capabilities: Capabilities!, $uuid: String!, $fulfillmentMethods: [FulfillmentMethod!]) { fulfillment_times: fulfillment_times(location: $location, capabilities: $capabilities, uuid: $uuid, fulfillment_methods: $fulfillmentMethods) { __typename ...fulfillmentTimesResultFields } results: search(location: $location, uuid: $uuid, options: $options, capabilities: $capabilities) { __typename ui_layout_groups { __typename ...uiLayoutGroups } ui_control_groups { __typename ...uiControlGroups } ui_modals { __typename ...uiAnyModalFields } ui_feed_overlays { __typename ...uiFeedOverlayFields } meta { __typename ...metaFields } } } fragment fulfillmentTimesResultFields on FulfillmentTimesResult { __typename fulfillment_time_methods { __typename fulfillment_method_label fulfillment_method asap { __typename ...optionFields } days { __typename day_label times { __typename ...optionFields } } } } fragment optionFields on FulfillmentTimeOption { __typename option_label selected_label timestamp selected_time { __typename day time } } fragment uiLayoutGroups on UILayoutGroup { __typename id subheader ui_layouts { __typename ...uiLayoutFields } } fragment uiLayoutFields on UILayout { __typename ... on UILayoutCarousel { ui_lines { __typename ...uiLineFields } image_url style key tracking_id target_presentational target { __typename ...targetFields } ui_blocks { __typename ...uiBlockFields } rows carousel_name } ... on UILayoutList { header key tracking_id ui_blocks { __typename ...uiBlockFields } } } fragment uiLineFields on UILine { __typename ... on UITextLine { ui_spans { __typename ...uiSpans } } ... on UITitleLine { text color { __typename ...colorFields } size } ... on UIBulletLine { icon_span { __typename ...uiSpanIcon } bullet_spacer_span { __typename ...uiSpanSpacer } ui_spans { __typename ...uiSpans } } } fragment uiSpans on UISpan { __typename ...uiSpanIcon ...uiSpanSpacer ...uiSpanText ...uiSpanCountdown } fragment uiSpanIcon on UISpanIcon { __typename color { __typename ...colorFields } icon { __typename ...iconFields } iconSize: size } fragment uiSpanSpacer on UISpanSpacer { __typename width } fragment uiSpanText on UISpanText { __typename color { __typename ...colorFields } text size is_bold } fragment uiSpanCountdown on UISpanCountdown { __typename color { __typename ...colorFields } ends_at is_bold size } fragment colorFields on Color { __typename red green blue alpha } fragment iconFields on DeliverooIcon { __typename name image } fragment targetFields on UITarget { __typename ... on UITargetRestaurant { restaurant { __typename ...restaurantTargetFields } ad_id } ... on UITargetParams { ...targetParamFields } ... on UITargetAction { action params { __typename id value } } ... on UITargetMenuItem { menu_item { __typename id } restaurant { __typename ...restaurantTargetFields } } ... on UITargetMenuItemModifier { restaurant_id menu_item_id } ... on UITargetWebPage { ...webPageTargetFields } ... on UITargetDeepLink { uri fallback_target { __typename ...webPageTargetFields } } } fragment restaurantTargetFields on Restaurant { __typename id delivery_status_presentational name is_new_menu_enabled images { __typename default } } fragment targetParamFields on UITargetParams { __typename title subtitle applied_filter_label is_landing_page params { __typename ...paramFields } } fragment paramFields on Param { __typename id value } fragment webPageTargetFields on UITargetWebPage { __typename url } fragment uiBlockFields on UIBlock { __typename ... on UIBanner { banner_id button_caption caption header key tracking_id target { __typename ...targetFields } images { __typename image } ui_theme background_color { __typename ...colorFields } content_description } ... on UICard { key target { __typename ...targetFields } tracking_id restaurant { __typename id } properties { __typename default { __typename bubble { __typename ui_lines { __typename ...uiLineFields } } image ui_lines { __typename ...uiLineFields } overlay { __typename background { __typename ...backgroundColorFields } text { __typename color { __typename ...colorFields } value } promotion_tag { __typename primary_tag_line { __typename ...promotionTagLine } secondary_tag_line { __typename ...promotionTagLine } } } countdown_badge_overlay { __typename background_color { __typename ...colorFields } ui_line { __typename ...uiLineFields } } favourites_overlay { __typename ...uiFavouritesOverlay } } } ui_theme content_description border { __typename border_width top_color { __typename ...colorFields } bottom_color { __typename ...colorFields } left_color { __typename ...colorFields } right_color { __typename ...colorFields } } } ... on UIShortcut { images { __typename default } name name_color { __typename ...colorFields } background_color { __typename ...colorFields } target { __typename ...targetFields } ui_theme tracking_id key } ... on UIButton { text content_description key target { __typename ...targetFields } ui_theme tracking_id } ... on UIMerchandisingCard { tracking_id content_description key target { __typename ...targetFields } button_caption ui_lines { __typename ...uiLineFields } background_color_non_null: background_color { __typename ...colorFields } header_image_url background_image_url } ... on UICategoryPill { content { __typename ...uiLineFields } background_color { __typename ...backgroundColorFields } target { __typename ...targetFields } tracking_id } ... on UITallMenuItemCard { key tracking_id image menu_item_id title price { __typename ... on Currency { code fractional formatted } } target { __typename ...targetFields } } } fragment uiFavouritesOverlay on UIFavouritesOverlay { __typename id entity selected_color { __typename ...colorFields } unselected_color { __typename ...colorFields } background_color { __typename ...backgroundColorFields } is_selected target { __typename ...targetFields } count_data { __typename ...favouriteCountDataFields } } fragment backgroundColorFields on UIBackgroundColor { __typename ... on Color { ...colorFields } ... on ColorGradient { from { __typename ...colorFields } to { __typename ...colorFields } } } fragment favouriteCountDataFields on FavouriteCountData { __typename count is_max_count } fragment promotionTagLine on UICardPromotionTagLine { __typename background_color { __typename ...backgroundColorFields } text { __typename ...uiLineFields } } fragment uiControlGroups on UIControlGroups { __typename applied_filters { __typename ...appliedFilters } filters { __typename ...uiControlFilterFields } fulfillment_methods { __typename ...uiControlFulfillmentMethodFields } layout_groups { __typename ...controlLayoutGroupFields } sort { __typename ...uiControlFilterFields } query_results { __typename ...uiControlQueryResultFields } } fragment appliedFilters on UIControlAppliedFilter { __typename label target_params { __typename ...targetParamFields } } fragment uiControlFilterFields on UIControlFilter { __typename id header images { __typename icon { __typename ...iconFields } } options_type options { __typename id count default disabled header selected target_params { __typename ...targetParamFields } } styling { __typename android { __typename collapse } } } fragment uiControlFulfillmentMethodFields on UIControlFulfillmentMethod { __typename label target_method } fragment controlLayoutGroupFields on UIControlLayoutGroup { __typename label selected_by_default target_layout_group { __typename layout_group_id } } fragment uiControlQueryResultFields on UIControlQueryResult { __typename header options { __typename image { __typename ... on DeliverooIcon { ...iconFields } ... on UIControlQueryResultOptionImageSet { default } } label target { __typename ...targetFields } ui_lines { __typename ...uiLineFields } key tracking_id is_available } tracking_id result_target_presentational result_target { __typename ...targetFields } } fragment uiAnyModalFields on UIAnyModal { __typename ... on UIModal { ...uiModalFields } ... on UIChallengesModal { ...uiChallengesModalFields } ... on UIPlusFullScreenModal { ...uiPlusFullScreenModalFields } } fragment uiModalFields on UIModal { __typename header caption image { __typename ...modalImageSetFields } buttons { __typename ...uiModalButtonFields target { __typename ...targetFields } } tracking_id ui_theme display_id display_only_once } fragment modalImageSetFields on UIModalImageSet { __typename ... on DeliverooIcon { ...iconFields } ... on DeliverooIllustrationBadge { ...illustrationBadgeFields } ... on UIModalImage { ...modalImageFields } } fragment illustrationBadgeFields on DeliverooIllustrationBadge { __typename name image } fragment modalImageFields on UIModalImage { __typename image } fragment uiModalButtonFields on UIModalButton { __typename title ui_theme dismiss_on_action tracking_id } fragment uiChallengesModalFields on UIChallengesModal { __typename challenges_drn_id sparkles full_view { __typename ...uiChallengesFullViewFields } mode display_id tracking_id display_only_once } fragment uiChallengesFullViewFields on UIChallengesFullView { __typename icon { __typename ...uiChallengesIconFields } header header_subtitle body_title body_text confirmation_button { __typename ...uiModalButtonFields target { __typename ...targetFields } } info_button { __typename ...uiModalButtonFields } } fragment uiChallengesIconFields on UIChallengesIcon { __typename ... on UIChallengesIndicator { required completed } ... on UIChallengesBadge { url } ... on UIChallengesSteppedIndicator { steps { __typename ...uiChallengeSteps } steps_completed } } fragment uiChallengeSteps on UIChallengesSteppedStamp { __typename text icon is_highlighted } fragment uiPlusFullScreenModalFields on UIPlusFullScreenModal { __typename image { __typename ...modalImageSetFields } header body footnote primaryButton { __typename ...uiModalButtonFields target { __typename ...targetFields } } secondaryButton { __typename ...uiModalButtonFields target { __typename ...targetFields } } tracking_id display_id display_only_once confetti } fragment uiFeedOverlayFields on UIFeedOverlay { __typename id position overlay_blocks { __typename ...uiFeedOverlayBlockFields } } fragment uiFeedOverlayBlockFields on UIFeedOverlayBlock { __typename ... on UIFeedOverlayBanner { header caption image { __typename ...uiFeedOverlayBlockImageFields } ui_theme display_id tracking_id is_dismissible } } fragment uiFeedOverlayBlockImageFields on UIFeedOverlayBlockImage { __typename ... on DeliverooIllustrationBadge { ...illustrationBadgeFields } } fragment metaFields on SearchResultMeta { __typename title uuid options { __typename query } search_placeholder search_results_title search_results_subtitle validity_ms collection { __typename search_bar_meta { __typename search_bar_placeholder search_bar_params { __typename id value } } } search_pills { __typename id label placeholder params { __typename id value } } }",
		"variables": map[string]any{
			"location": map[string]float64{
				"lat": latitude,
				"lon": longitude,
			},
			"options": map[string]any{
				"deep_link":          "deliveroo://restaurants",
				"fulfillment_method": "DELIVERY",
				"web_column_count":   1,
				"user_preference":    map[string]any{},
			},
			"capabilities": map[string][]string{
				"ui_actions":                {"CHANGE_DELIVERY_TIME", "CLEAR_FILTERS", "NO_DELIVERY_YET", "SHOWCASE_PICKUP", "SHOW_HOME_MAP_VIEW", "SHOW_MEAL_CARD_ISSUERS", "SHOW_PLUS_SIGN_UP", "TOGGLE_FAVOURITE", "COPY_TO_CLIPBOARD"},
				"ui_blocks":                 {"BANNER", "BUTTON", "CARD", "SHORTCUT", "MERCHANDISING_CARD"},
				"ui_controls":               {"APPLIED_FILTER", "FILTER", "LAYOUT_GROUP", "QUERY_RESULT", "SORT"},
				"ui_layouts":                {"CAROUSEL", "LIST"},
				"ui_targets":                {"LAYOUT_GROUP", "MENU_ITEM", "PARAMS", "RESTAURANT", "WEB_PAGE", "DEEP_LINK"},
				"ui_themes":                 {"BANNER_CARD", "BANNER_EMPTY", "BANNER_MARKETING_A", "BANNER_MARKETING_C", "BANNER_PICKUP_SHOWCASE", "BANNER_SERVICE_ADVISORY", "BUTTON_PRIMARY", "BUTTON_SECONDARY", "CARD_INFORMATIVE", "CARD_LARGE", "CARD_MEDIUM", "CARD_MEDIUM_HORIZONTAL", "CARD_SMALL", "CARD_SMALL_DIAGONAL", "CARD_SMALL_HORIZONTAL", "CARD_TALL", "CARD_TALL_GRADIENT", "CARD_WIDE", "MODAL_BUTTON_TERTIARY", "MODAL_DEFAULT", "MODAL_PLUS", "SHORTCUT_DEFAULT", "SHORTCUT_STACKED", "ANY_MODAL"},
				"ui_layout_carousel_styles": {"DEFAULT"},
				"ui_features":               {"HOME_MAP_VIEW", "ILLUSTRATION_BADGES", "SCHEDULED_RANGES", "UI_CARD_BORDER", "UI_CAROUSEL_COLOR", "UI_PROMOTION_TAG", "UI_SPAN_COUNTDOWN", "UNAVAILABLE_RESTAURANTS"},
				"ui_lines":                  {"TITLE", "TEXT", "BULLET"},
				"fulfillment_methods":       {"COLLECTION", "DELIVERY"},
			},
			"uuid":               _uuid.String(),
			"fulfillmentMethods": []string{"COLLECTION", "DELIVERY"},
		},
	}

	return json.Marshal(theMap)
}

func GetCreatePaymentQuery() ([]byte, error) {
	theMap := map[string]any{
		"operationName": "CreatePaymentPlan",
		"query":         "query CreatePaymentPlan($delivery_address_id: String, $payment_limitations: [InputPaymentOptionState!]!, $capabilities: Capabilities, $clientEvent: ClientEvent) { payment_plan: create_payment_plan(delivery_address_id: $delivery_address_id, payment_limitations: $payment_limitations, capabilities: $capabilities, client_event: $clientEvent) { __typename id hide_old_payment_section blocks { __typename ...sectionTarget } fulfillment_details { __typename fulfillment_type restaurant eta { __typename title description } in_person_fulfillment_address { __typename address1 city post_code phone distance_presentational coordinates { __typename latitude longitude } } } delivery_addresses { __typename available { __typename ...address } selected { __typename ...address } add_new_address_cta } payment_options { __typename completing { __typename ... payment_option ... completing_payment_option } selected_completing { __typename ... payment_option ... completing_payment_option } promoted { __typename action_url description title icon_url type payment_token_id } fund_balances { __typename ... payment_option } new_card_config { __typename add_card_cta tokenizer api_key description store_option_opt_in_out { __typename selected label } } google_pay_config { __typename tokenizer api_key } paypal_config { __typename api_key add_card_cta tokenizer } auth_redirect_payment_option_configs { __typename action_url icon_url title description } } line_item_groups(filter_by_type: [TOTAL]) { __typename group_type line_items { __typename title cost } } execution_state { __typename execution_cta is_executable banner { __typename ...uiBanner } show_add_payment_method_cta } footer_note marketing_preferences { __typename id content selected type } loyalty_card { __typename ...ui_loyalty_card } payment_breakdown_block { __typename ...paymentBreakdownTarget } meal_card_toggle_block { __typename ...mealCardPresentTarget ...mealCardAbsentTarget } tracking_data { __typename ... on TrackingData { currency_code } } } } fragment sectionTarget on UISection { __typename key title id blocks { __typename ...checkboxTarget ...bannerPickPaymentMethodTarget ...bannerFeaturedPaymentMethodTarget ...selectedPaymentMethodTarget ...selectedPaymentMethodCreditTarget ...paymentBreakdownTarget ...mealCardAbsentTarget ...mealCardPresentTarget ...ui_loyalty_card } } fragment checkboxTarget on UICheckbox { __typename ... on UICheckbox { key id default required submission { __typename ...submission } event_tracking { __typename ...eventTarget } label { __typename ...labelTarget } invalid_state { __typename event_tracking { __typename ...eventTarget } description { __typename ...labelTarget } } } } fragment bannerPickPaymentMethodTarget on UIBannerPickPaymentMethod { __typename ... on UIBannerPickPaymentMethod { key ui_lines { __typename ...uiLineFields } background_image { __typename url title } cta methodsScreen: methods_screen { __typename ...pickPaymentMethodsScreenFields } } } fragment bannerFeaturedPaymentMethodTarget on UIBannerFeaturedPaymentMethod { __typename ... on UIBannerFeaturedPaymentMethod { key title ui_lines { __typename ...uiLineFields } background_image { __typename url title } featured_payment_methods_icons { __typename title url } additional_payment_methods_counter_label methodsScreen: methods_screen { __typename ...pickPaymentMethodsScreenFields } } } fragment selectedPaymentMethodTarget on UISelectedPaymentMethod { __typename key ui_lines { __typename ...uiLineFields } icon { __typename url } cta selected_option { __typename ...paymentTokenWithoutVariants } methods_screen { __typename ...changeMethodScreen } } fragment selectedPaymentMethodCreditTarget on UISelectedPaymentMethodCredit { __typename key ui_lines { __typename ...uiLineFields } icon { __typename url } } fragment paymentBreakdownTarget on UIPaymentBreakdown { __typename key rows { __typename label { __typename ...uiLineFields } cost { __typename ...uiLineFields } } } fragment mealCardAbsentTarget on UIMealCardAbsentToggle { __typename key meal_card_toggle { __typename ...toggleTarget } methods_screen { __typename ...pickPaymentMethodsScreenFields } banner { __typename ...uiBanner } } fragment mealCardPresentTarget on UIMealCardPresentToggle { __typename key meal_card_toggle { __typename ...toggleTarget } meal_card { __typename ... paymentTokenWithoutVariants } banner { __typename ...uiBanner } } fragment ui_loyalty_card on UILoyaltyCard { __typename key ui_lines { __typename ...uiLineFields } text_box { __typename ...ui_text_box } cta header_image { __typename url } submit_path } fragment submission on UISubmission { __typename field_name submitted_with } fragment eventTarget on EventTracking { __typename event_name default_properties { __typename name value } } fragment labelTarget on UIAttributedString { __typename color { __typename ...colorRGBAFields } content prepend_icon } fragment colorRGBAFields on UIColorRGBA { __typename red green blue alpha } fragment uiLineFields on UILine { __typename ... on UITitleLine { key text color { __typename ... colorFields } } ... on UITextLine { key ui_spans { __typename ...uiSpans } } } fragment colorFields on Color { __typename hex red green blue alpha } fragment uiSpans on UISpan { __typename ...uiSpanIcon ...uiSpanSpacer ...uiSpanText } fragment uiSpanIcon on UISpanIcon { __typename color { __typename ...colorFields } icon { __typename ...iconFields } iconSize: size } fragment uiSpanSpacer on UISpanSpacer { __typename width } fragment uiSpanText on UISpanText { __typename color { __typename ...colorFields } text size is_bold } fragment iconFields on DeliverooIcon { __typename name image } fragment pickPaymentMethodsScreenFields on UIBannerPickPaymentMethodsScreen { __typename banner { __typename ui_lines { __typename ...uiLineFields } cta } title cta cancel_cta blocks { __typename ...paypalPaymentMethodFields ...newCardFields ...authRedirectPaymentMethodFields ...selectablePaymentOptionFields ...mobileWalletPaymentMethodFields ...selectablePaymentOptionWithVariantsFields } } fragment paypalPaymentMethodFields on UIPaypalPaymentMethod { __typename key icon_url api_key ui_lines { __typename ...uiLineFields } } fragment newCardFields on UINewCard { __typename key cta cancel_cta icon_url ui_lines { __typename ...uiLineFields } tokenizer api_key storeOptionOptInOut: store_option_opt_in_out { __typename selected label } cardIcons: card_icons { __typename url title } } fragment authRedirectPaymentMethodFields on UIAuthRedirectPaymentMethod { __typename key title action_url icon_url ui_lines { __typename ...uiLineFields } } fragment selectablePaymentOptionFields on UISelectablePaymentOption { __typename key id is_selectable icon_url ui_lines { __typename ...uiLineFields } } fragment mobileWalletPaymentMethodFields on UIMobileWalletPaymentMethod { __typename id key api_key icon_url tokenizer ui_lines { __typename ...uiLineFields } } fragment selectablePaymentOptionWithVariantsFields on UISelectablePaymentOptionWithVariants { __typename key id is_selectable icon_url variantsScreen: variants_screen { __typename ...selectablePaymentOptionWithVariantsScreenFields } ui_lines { __typename ...uiLineFields } } fragment selectablePaymentOptionWithVariantsScreenFields on UISelectablePaymentOptionWithVariantsScreen { __typename title cancel_cta blocks { __typename ...paymentVariantFields } } fragment paymentVariantFields on UIPaymentVariant { __typename id name } fragment paymentTokenWithoutVariants on UIPaymentTokenOptionWithoutVariants { __typename id ui_lines { __typename ...uiLineFields } icon { __typename url } is_selectable is_selected is_deletable } fragment changeMethodScreen on UIChangePaymentMethodScreen { __typename ui_lines { __typename ...uiLineFields } try_another_way_to_pay_cta methods_screen { __typename ...pickPaymentMethodsScreenFields } blocks { __typename ...paymentTokenWithoutVariants } } fragment toggleTarget on UIToggle { __typename title is_selected is_selectable ui_lines { __typename ...uiLineFields } event_tracking { __typename ...eventTarget } } fragment uiBanner on UIBanner { __typename title description cta } fragment ui_text_box on UITextbox { __typename default placeholder submission { __typename ...submission } } fragment address on DeliveryAddress { __typename id title short_description long_description location { __typename latitude longitude } is_selectable delivery_note phone_number } fragment payment_option on PaymentOption { __typename id title description icon_url is_selectable proposed_amount { __typename numerical currency_code } } fragment completing_payment_option on CompletingPaymentOption { __typename payment_type variants { __typename title selected { __typename name id } variants { __typename name id } } }",
		"variables": map[string]any{
			"delivery_address_id": nil,
			"payment_limitations": []any{},
			"capabilities": map[string]any{
				"wallets":                []any{},
				"payment_capabilities":   []string{"RETURN_PAYPAL_PAYMENT_OPTIONS", "PAYPAL_UPSELL", "RETURN_IDEAL", "RETURN_PAYMENT_TOKEN_TYPE", "PAYMENT_TOKEN_UPSELL"},
				"ui_blocks_capabilities": []string{"PAYMENT_SECTION", "TERMS_AND_CONDITIONS_SECTION"},
			},
			"clientEvent": "INITIAL_PLAN_LOAD",
		},
	}

	return json.Marshal(theMap)
}