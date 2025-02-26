// Package restapi implements the rest api handlers for Braavos API
package dashboard

import (
	"errors"
	"net/http"
	"sync"

	"github.com/Rhymond/go-money"
	"github.com/wingocard/braavos/internal/legacy/database/dalerr"
	"github.com/wingocard/braavos/internal/legacy/service/invitationsvc"
	"github.com/wingocard/braavos/internal/legacy/service/issuingsvc"
	"github.com/wingocard/braavos/internal/legacy/service/kycsvc"
	"github.com/wingocard/braavos/internal/legacy/service/usersvc"
	"github.com/wingocard/braavos/internal/wlog"
	"github.com/wingocard/braavos/pkg/entities/form"
	"github.com/wingocard/braavos/pkg/entities/invitation"
	"github.com/wingocard/braavos/pkg/entities/issuing"
	"github.com/wingocard/braavos/pkg/entities/user"
	"github.com/wingocard/braavos/pkg/render"
)

type primaryUserResponse struct {
	userResponse
	PassedKYC bool `json:"passed_kyc"`
}

type userResponse struct {
	ID             string       `json:"id,omitempty"`
	Status         user.Status  `json:"status,omitempty"`
	Type           user.Type    `json:"type,omitempty"`
	FirstName      string       `json:"first_name"`
	LastName       string       `json:"last_name"`
	DateOfBirth    form.Date    `json:"date_of_birth"`
	Email          string       `json:"email,omitempty"`
	Phone          string       `json:"phone,omitempty"`
	EmailValidated *bool        `json:"email_validated,omitempty"`
	Username       string       `json:"username"`
	Address        *addressResp `json:"address,omitempty"`
}

type addressResp struct {
	AddressLine1  string  `json:"address_line_1"`
	AddressLine2  *string `json:"address_line_2,omitempty"`
	City          string  `json:"city"`
	ProvState     string  `json:"prov_state"`
	PostalCodeZip string  `json:"postal_code_zip"`
	Country       string  `json:"country"`
}

type cardResponse struct {
	Last4 string `json:"last4"`
}

type bankAccountResponse struct {
	Balance   form.Money    `json:"balance"`
	Card      *cardResponse `json:"card,omitempty"`
	CreatedAt form.DateTime `json:"created_at"`
}

type relationshipResp struct {
	Type          user.RelationshipType `json:"type"`
	IsSponsorship bool                  `json:"is_sponsorship"`
	CreatedAt     form.DateTime         `json:"created_at"`
	User          userResponse          `json:"user"`
	BankAccount   *bankAccountResponse  `json:"bank_account"`
}

type invitationResp struct {
	FirstName        string                `json:"first_name"`
	LastName         string                `json:"last_name"`
	Email            string                `json:"email"`
	Phone            *string               `json:"phone"`
	Token            string                `json:"token"`
	DateOfBirth      *form.DateTime        `json:"date_of_birth"`
	CreatedAt        form.DateTime         `json:"created_at"`
	ExpiresAt        form.DateTime         `json:"expires_at"`
	Type             invitation.Type       `json:"type"`
	RelationshipType user.RelationshipType `json:"relationship_type"`
	IsExpired        bool                  `json:"is_expired"`
}

type dashboardResponse struct {
	User          primaryUserResponse  `json:"user"`
	BankAccount   *bankAccountResponse `json:"bank_account"`
	Relationships []relationshipResp   `json:"relationships"`
	Invites       []invitationResp     `json:"invites"`
	Shutdown      bool                 `json:"shutdown"`
}

type balanceResult struct {
	UserID  string
	Balance *money.Money
	Error   error
}

// Dashboard handles HTTP requests for the app home screen.
func Dashboard(
	wl wlog.Logger,
	userService usersvc.UserSVC,
	invitationService invitationsvc.InvitationSVC,
	issuingService issuingsvc.IssuingSVC,
	kycService kycsvc.Service,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		wl := wlog.WithServiceRequest(ctx, wl, "dashboard")

		userID, err := usersvc.GetUserIDContext(ctx)
		if err != nil {
			render.InternalError(ctx, wl, w, err)
			return
		}

		u, err := userService.GetUserByID(ctx, userID)
		if err != nil {
			render.InternalError(ctx, wl, w, err)
			return
		}

		relationships, err := userService.GetRelationshipDetails(ctx, u)
		if err != nil {
			render.InternalError(ctx, wl, w, err)
			return
		}

		invites, err := invitationService.GetInvitationsByUserID(ctx, userID)
		if err != nil {
			render.InternalError(ctx, wl, w, err)
			return
		}

		passedKYC, err := kycService.DidUserPassKYC(ctx, userID)
		if err != nil {
			render.InternalError(ctx, wl, w, err)
			return
		}

		userIDs := []string{u.ID}
		if u.Type == user.Parent {
			for _, rel := range relationships {
				userIDs = append(userIDs, rel.User.ID)
			}
		}

		// Slice for parent and teens balances
		balances := make(map[string]money.Money, len(userIDs))
		bankAccounts := make(map[string]*issuing.BankAccount, len(userIDs))
		if u.Type == user.Parent || len(relationships) > 0 {
			bankAccounts, err = issuingService.GetBankAccountsByUserIDs(ctx, userIDs)
			if err != nil && !errors.Is(err, issuingsvc.ErrAccountNotFound) {
				render.InternalError(ctx, wl, w, err)
				return
			}

			if len(bankAccounts) > 0 {
				balanceCh := make(chan balanceResult, len(bankAccounts))
				getBalanceAsync := func(ba *issuing.BankAccount, wg *sync.WaitGroup, balanceCh chan balanceResult) {
					defer wg.Done()
					balance, err := issuingService.GetBalance(ctx, ba)
					balanceCh <- balanceResult{UserID: ba.UserID, Balance: balance, Error: err}
				}

				// Get parent and teens balance asynchronously
				var wg sync.WaitGroup
				wg.Add(len(bankAccounts))
				for _, bankAccount := range bankAccounts {
					go getBalanceAsync(bankAccount, &wg, balanceCh)
				}
				wg.Wait()
				close(balanceCh)

				// Create a map[userID]balance for easy lookup
				for br := range balanceCh {
					if br.Error != nil { // we return an error if any call fails
						render.InternalError(ctx, wl, w, br.Error)
						return
					}
					balances[br.UserID] = *br.Balance
				}
			}
		}

		var userBankAccount *bankAccountResponse
		if b, ok := balances[u.ID]; ok {
			userBankAccount = &bankAccountResponse{
				Balance:   form.Money(b),
				CreatedAt: form.DateTime(bankAccounts[u.ID].CreatedAt),
			}
			if u.Type == user.Teen {
				userBankAccount.Card = &cardResponse{Last4: bankAccounts[u.ID].CardLast4}
			}
		}
		resp := dashboardResponse{
			User: primaryUserResponse{
				userResponse: userResponse{
					ID:             u.ID,
					Status:         u.Status,
					Type:           u.Type,
					FirstName:      u.FirstName,
					LastName:       u.LastName,
					Email:          u.Email,
					Phone:          u.Phone,
					EmailValidated: &u.EmailValidated,
					Username:       u.Username,
					DateOfBirth:    form.Date(u.DateOfBirth),
				},
				PassedKYC: passedKYC,
			},
			BankAccount:   userBankAccount,
			Invites:       make([]invitationResp, len(invites)),
			Relationships: make([]relationshipResp, len(relationships)),
		}
		if u.Type == user.Parent {
			a, err := userService.GetUserAddressByUserID(ctx, u.ID)
			if err != nil && !errors.Is(&dalerr.EntityNotFoundError{}, err) {
				render.InternalError(ctx, wl, w, err)
				return
			}
			if a != nil {
				resp.User.Address = &addressResp{
					AddressLine1:  a.AddressLine1,
					AddressLine2:  a.AddressLine2,
					City:          a.City,
					ProvState:     a.ProvState,
					PostalCodeZip: a.PostalCodeZip,
					Country:       a.Country,
				}
			}
		}
		for i, invite := range invites {
			invResp := invitationResp{
				FirstName:        invite.InviteeFirstName,
				LastName:         invite.InviteeLastName,
				Email:            invite.InviteeEmail,
				Phone:            invite.InviteePhone,
				Token:            invite.Token,
				CreatedAt:        form.DateTime(invite.CreatedAt),
				ExpiresAt:        form.DateTime(invite.ExpiresAt),
				Type:             invite.Type,
				RelationshipType: invite.InviteeRelationshipType,
				IsExpired:        invite.IsExpired(),
			}
			if invite.InviteeDateOfBirth != nil {
				dob := form.DateTime(*invite.InviteeDateOfBirth)
				invResp.DateOfBirth = &dob
			}
			resp.Invites[i] = invResp
		}
		for i, rs := range relationships {
			resp.Relationships[i] = relationshipResp{
				Type:          rs.Relationship.Type,
				IsSponsorship: rs.Relationship.IsSponsorship,
				CreatedAt:     form.DateTime(rs.Relationship.CreatedAt),
				User: userResponse{
					ID:          rs.User.ID,
					Type:        rs.User.Type,
					FirstName:   rs.User.FirstName,
					LastName:    rs.User.LastName,
					Username:    rs.User.Username,
					DateOfBirth: form.Date(rs.User.DateOfBirth),
				},
			}
			if u.Type == user.Parent {
				resp.Relationships[i].BankAccount = &bankAccountResponse{
					Balance:   form.Money(balances[rs.User.ID]),
					CreatedAt: form.DateTime(bankAccounts[rs.User.ID].CreatedAt),
				}
			}
		}

		// program shutdown
		resp.Shutdown = true

		render.JSON(ctx, wl, w, resp, http.StatusOK)
	}
}
