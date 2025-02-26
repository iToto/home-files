// Package publicrouter provides a HTTP REST router
// that contains all endpoints for performing any publically
// available operation. These endpoints are public and are available to all
// clients. Some require a user be authenticated but other's are completely public.
package publicrouter

import (
	"net/http"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/gorilla/mux"
	"github.com/ulule/limiter/v3"
	legacy_dashboard_restapi "github.com/wingocard/braavos/internal/legacy/handler/restapi/dashboard"
	legacy_health_restapi "github.com/wingocard/braavos/internal/legacy/handler/restapi/health"
	legacy_invitation_restapi "github.com/wingocard/braavos/internal/legacy/handler/restapi/invitation"
	legacy_issuing_restapi "github.com/wingocard/braavos/internal/legacy/handler/restapi/issuing"
	legacy_user_restapi "github.com/wingocard/braavos/internal/legacy/handler/restapi/user"
	legacy_server "github.com/wingocard/braavos/internal/legacy/server"
	legacy_middleware "github.com/wingocard/braavos/internal/legacy/server/middleware"
	legacy_usersvc "github.com/wingocard/braavos/internal/legacy/service/usersvc"
	"github.com/wingocard/braavos/internal/rest/backingservice"
	"github.com/wingocard/braavos/internal/rest/handler/activityhdl"
	"github.com/wingocard/braavos/internal/rest/handler/balancehdl"
	"github.com/wingocard/braavos/internal/rest/handler/disablehdl"
	"github.com/wingocard/braavos/internal/rest/handler/issuingeventshdl"
	"github.com/wingocard/braavos/internal/rest/handler/issuinghdl"
	"github.com/wingocard/braavos/internal/rest/handler/phonehdl"
	"github.com/wingocard/braavos/internal/rest/handler/recurringhdl"
	"github.com/wingocard/braavos/internal/rest/handler/rewardhdl"
	"github.com/wingocard/braavos/internal/rest/handler/userhdl"
	"github.com/wingocard/braavos/internal/rest/middleware"
	"github.com/wingocard/braavos/internal/rest/middleware/eventauth"
	"github.com/wingocard/braavos/internal/wlog"
)

// New creates a new public router.
func New(wl wlog.Logger, bs *backingservice.BackingServices) (*mux.Router, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	router := mux.NewRouter()

	w2wRateLimiter, err := middleware.NewRateLimitMiddleware(bs.RedisClient)
	if err != nil {
		return nil, err
	}
	tm := middleware.NewCloudTraceMiddleware(cfg.GCPProjectID)

	router.Use(legacy_middleware.ExtractIPAddressMiddleware)
	router.Use(middleware.RequestLogVarMiddleware)
	router.Use(tm.TraceHandler)

	router.HandleFunc("/", legacy_server.CreateHandler("Welcome to the Iron Bank!"))
	router.HandleFunc("/health", legacy_health_restapi.HandleHealthCheck(wl, *bs.Legacy.HealthService))

	// public endpoints subrouter
	v1Public := router.PathPrefix("/public/v1").Subrouter()

	v1Public.Handle("/parent", disablehdl.DisabledEndpoint(
		wl,
		"create-parent",
	)).Methods(http.MethodPost, http.MethodOptions)

	v1Public.Handle("/teen", disablehdl.DisabledEndpoint(
		wl,
		"create-teen",
	)).Methods(http.MethodPost, http.MethodOptions)

	phoneRouter := v1Public.PathPrefix("/phone").Subrouter()

	phoneRouter.HandleFunc("/request_otp", phonehdl.RequestOTP(
		wl,
		bs.ValidatePhoneService,
	)).Methods(http.MethodPost, http.MethodOptions)

	phoneRouter.HandleFunc("/check_otp", phonehdl.CheckOTP(
		wl,
		bs.ValidatePhoneService,
	)).Methods(http.MethodPost, http.MethodOptions)

	signupRewardRouter := v1Public.PathPrefix("/signup_reward").Subrouter()

	signupRewardRouter.HandleFunc(
		"/{referral_code}",
		rewardhdl.GetSignupReferralReward(
			wl,
			bs.RewardService,
		)).Methods(http.MethodGet, http.MethodOptions)

	addLegacyV1PublicRoutes(bs, wl, v1Public)

	// private endpoints subrouter
	v1 := router.PathPrefix("/v1").Subrouter()
	v1.Use(legacy_usersvc.AuthMiddleware(bs.Legacy.UserService, wl))

	v1.HandleFunc("/parent",
		userhdl.PatchParent(wl, bs.UserService),
	).Methods(http.MethodPatch)

	v1.HandleFunc("/parent/kyc",
		disablehdl.DisabledEndpoint(wl, "parent-kyc"),
	).Methods(http.MethodPost)

	v1.HandleFunc("/parent/autosend",
		recurringhdl.GetParentAutoSends(wl, bs.RecurringService),
	).Methods(http.MethodGet)

	v1.HandleFunc("/parent/autosend",
		recurringhdl.CreateParentAutoSend(wl, bs.RecurringService),
	).Methods(http.MethodPost)

	v1.HandleFunc("/parent/autosend/{id}",
		recurringhdl.DeleteParentAutoSend(wl, bs.RecurringService),
	).Methods(http.MethodDelete)

	v1.HandleFunc("/card",
		issuinghdl.ViewCard(wl, bs.IssuingService),
	).Methods(http.MethodGet)

	v1.HandleFunc("/referral_info",
		rewardhdl.GetReferralInfo(wl, bs.RewardService),
	).Methods(http.MethodGet)

	v1.HandleFunc("/rewards",
		rewardhdl.GetRewards(wl, bs.RewardService),
	).Methods(http.MethodGet)

	v1.HandleFunc("/qrcode",
		rewardhdl.GetQRCodeImage(wl, bs.RewardService),
	).Methods(http.MethodGet)

	v1.HandleFunc("/balance",
		balancehdl.GetBalance(wl, bs.RewardService, bs.IssuingService),
	).Methods(http.MethodGet)

	v1.Handle("/teen_to_teen_transfer",
		w2wRateLimiter.LimitHandler(limiter.Rate{
			Period: time.Second * time.Duration(cfg.W2WRateInterval),
			Limit:  cfg.W2WRateLimit,
		}, issuinghdl.SendMoneyTeenToTeen(wl, bs.IssuingService)),
	).Methods(http.MethodPost)

	v1.HandleFunc(
		"/parent_to_teen_transfer",
		issuinghdl.SendMoneyParentToTeen(wl, bs.IssuingService),
	).Methods(http.MethodPost)

	v1Legacy := router.PathPrefix("/v1").Subrouter()
	v1Legacy.Use(legacy_usersvc.AuthMiddleware(bs.Legacy.UserService, wl))

	// activity endpoints
	v1.HandleFunc(
		"/activity",
		activityhdl.GetActivities(wl, bs.ActivityService),
	).Methods(http.MethodGet)

	v1.HandleFunc(
		"/activity/merchant/{id}",
		activityhdl.GetMerchantActivity(wl, bs.ActivityService),
	).Methods(http.MethodGet)

	addLegacyV1Routes(bs, wl, v1Legacy)

	// v2 endpoints
	v2 := router.PathPrefix("/v2").Subrouter()
	v2.Use(legacy_usersvc.AuthMiddleware(bs.Legacy.UserService, wl))
	v2.Use(legacy_usersvc.StatusActiveMiddleware(bs.Legacy.UserService, wl))

	v2.HandleFunc("/funding_account",
		disablehdl.DisabledEndpoint(wl, "v2-funding-account"),
	).Methods(http.MethodPost)

	// events
	eventAuthMiddleware, err := eventauth.New()
	if err != nil {
		return nil, err
	}
	v1events := router.PathPrefix("/v1/events").Subrouter()
	v1events.Use(eventAuthMiddleware.AuthRequest(wl))
	v1events.Use(middleware.UnpackEvent(wl))

	v1events.Handle("/issuingsvc", issuinghdl.HandleEvents(wl, bs.IssuingService))
	v1events.Handle("/usersvc", userhdl.HandleEvents(wl, bs.UserService))
	v1events.Handle("/rewardsvc", rewardhdl.HandleEvents(wl, bs.RewardService))
	v1events.Handle("/recurringsvc", recurringhdl.HandleEvents(wl, bs.RecurringService))

	v1events.Handle("/processor-events", issuingeventshdl.HandleEvents(wl, bs.IssuingEventsService))

	return router, nil
}

func addLegacyV1PublicRoutes(
	bs *backingservice.BackingServices, wl wlog.Logger, v1Public *mux.Router,
) {
	v1Public.HandleFunc(
		"/login",
		legacy_user_restapi.LoginHandler(wl, bs.Legacy.UserService),
	).Methods(http.MethodPost)
	v1Public.HandleFunc(
		"/password/request_reset",
		legacy_user_restapi.ResetPasswordRequest(
			wl, bs.Legacy.UserService, bs.Legacy.BranchioClient, bs.Legacy.SendgridClient,
		),
	).Methods(http.MethodPost)
	v1Public.HandleFunc(
		"/password/reset",
		legacy_user_restapi.ResetPassword(wl, bs.Legacy.UserService),
	).Methods(http.MethodPost)
	v1Public.HandleFunc(
		"/invitation/{token}",
		legacy_invitation_restapi.GetInvitationHandler(wl, bs.Legacy.InvitationService, bs.Legacy.UserService),
	).Methods(http.MethodGet)
	v1Public.HandleFunc(
		"/email_verify",
		legacy_user_restapi.EmailVerify(wl, bs.Legacy.UserService, bs.Legacy.SegmentClient),
	).Methods(http.MethodPost)
}

func addLegacyV1Routes(
	bs *backingservice.BackingServices, wl wlog.Logger, v1 *mux.Router,
) {
	// Allow users to call the legacy /dashboard and /activity endpoints.
	v1EnabledUsers := v1.NewRoute().Subrouter()
	v1EnabledUsers.Use(legacy_usersvc.StatusUserEnabledMiddleware(bs.Legacy.UserService, wl))

	// dashboard endpoint
	v1EnabledUsers.HandleFunc(
		"/dashboard",
		legacy_dashboard_restapi.Dashboard(
			wl,
			bs.Legacy.UserService,
			bs.Legacy.InvitationService,
			bs.Legacy.IssuingService,
			bs.Legacy.KycService),
	).Methods(http.MethodGet)

	// invitation endpoints
	v1Invitation := v1EnabledUsers.PathPrefix("/invitation").Subrouter()
	v1Invitation.Use(legacy_usersvc.AuthMiddleware(bs.Legacy.UserService, wl))
	v1Invitation.HandleFunc(
		"/teen",
		legacy_invitation_restapi.CreateTeenInvitationHandler(
			wl,
			bs.Legacy.InvitationService,
			bs.Legacy.UserService,
			bs.Legacy.BranchioClient,
			bs.Legacy.SendgridClient,
			bs.Legacy.SegmentClient),
	).Methods(http.MethodPost)
	v1Invitation.HandleFunc(
		"/parent",
		legacy_invitation_restapi.CreateParentInvitationHandler(
			wl,
			bs.Legacy.InvitationService,
			bs.Legacy.UserService,
			bs.Legacy.BranchioClient,
			bs.Legacy.SendgridClient,
			bs.Legacy.SegmentClient),
	).Methods(http.MethodPost)
	// invitation accept endpoint
	v1Invitation.HandleFunc(
		"/{token}/accept",
		legacy_invitation_restapi.AcceptInvitationHandler(
			wl,
			bs.Legacy.InvitationService,
			bs.Legacy.UserService,
			bs.Legacy.SegmentClient,
			bs.Legacy.IssuingService),
	).Methods(http.MethodPost)
	v1Invitation.HandleFunc(
		"/{token}",
		legacy_invitation_restapi.DeleteInvitationHandler(wl, bs.Legacy.InvitationService),
	).Methods(http.MethodDelete)

	v1ActiveUsers := v1.NewRoute().Subrouter()
	v1ActiveUsers.Use(legacy_usersvc.StatusActiveMiddleware(bs.Legacy.UserService, wl))

	// funding account endpoints
	v1ActiveUsers.HandleFunc("/funding_account",
		disablehdl.DisabledEndpoint(wl, "funding-account"),
	).Methods(http.MethodGet)
	v1ActiveUsers.HandleFunc("/funding_account_encryption_key",
		disablehdl.DisabledEndpoint(wl, "funding-account-encryption-key"),
	).Methods(http.MethodGet)
	v1ActiveUsers.HandleFunc("/fund_bank_account",
		disablehdl.DisabledEndpoint(wl, "fund-bank-account"),
	).Methods(http.MethodPost)

	// issuing endpoints
	v1ActiveUsers.HandleFunc(
		"/ach_info",
		issuinghdl.GetACHInfo(wl, bs.IssuingService),
	).Methods(http.MethodGet)
	v1ActiveUsers.HandleFunc(
		"/provision_apple_pay",
		legacy_issuing_restapi.ApplePay(wl, bs.Legacy.IssuingService, bs.Legacy.UserService),
	).Methods(http.MethodPost)

	// user endpoints
	v1User := v1ActiveUsers.PathPrefix("/user").Subrouter()
	v1User.HandleFunc(
		"/me/request_email_verify",
		legacy_user_restapi.RequestEmailVerify(wl, bs.Legacy.UserService, bs.Legacy.BranchioClient, bs.Legacy.SegmentClient),
	).Methods(http.MethodPost)
}
