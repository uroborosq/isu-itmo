package ports

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/google/uuid"
	"github.com/uroborosq/isu/internal/isu/adapters"
	"github.com/uroborosq/isu/internal/isu/app"
	"golang.org/x/oauth2"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type UserPublicInfo struct {
	Isu         int
	Email       string
	PhoneNumber string
	FullName    string
}

type NewIDResponse struct {
	ID    int
	Error string
}

type PublicInfoResponse struct {
	Data  UserPublicInfo
	Error string
}

type HttpServer struct {
	service      app.IsuService
	oauth2Config *oauth2.Config
	ctx          context.Context
	provider     *oidc.Provider
	//cache 		 *cache.Cache
	cache map[string]string
}

func NewHttpServer(service app.IsuService, config *oauth2.Config, ctx context.Context, provider *oidc.Provider) *HttpServer {
	return &HttpServer{service: service,
		oauth2Config: config,
		ctx:          ctx,
		provider:     provider}
}

func (h *HttpServer) AddUser(w http.ResponseWriter, r *http.Request) {
	idStr := h.ctx.Value("userId").(string)
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "cant parse id", http.StatusInternalServerError)
	}
	var (
		publicInfo   UserPublicInfo
		bodyBytes    bytes.Buffer
		responseBody NewIDResponse
	)
	_, err = bodyBytes.ReadFrom(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		responseBody.Error = err.Error()
		responseBytes, err := json.Marshal(responseBody)
		if err != nil {
			log.Println(err.Error())
		}
		_, err = w.Write(responseBytes)
		if err != nil {
			log.Println(err.Error())
		}
	}
	err = json.Unmarshal(bodyBytes.Bytes(), &publicInfo)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		responseBytes, err := json.Marshal(responseBody)
		if err != nil {
			log.Println(err.Error())
		}
		_, err = w.Write(responseBytes)
		if err != nil {
			log.Println(err.Error())
		}
	}

	newUser := adapters.UserFullData{
		ID:          id,
		Isu:         publicInfo.Isu,
		Email:       publicInfo.Email,
		PhoneNumber: publicInfo.PhoneNumber,
		FullName:    publicInfo.FullName,
		Role:        adapters.User,
	}

	err = h.service.AddUser(newUser)
	if err != nil {
		responseBody.Error = err.Error()
	}
	responseBytes, err := json.Marshal(responseBody)
	if err != nil {

	}
	_, err = w.Write(responseBytes)
	if err != nil {
		log.Println(err.Error())
	}
}

func (h *HttpServer) GetPublicInfoByPhoneNumber(w http.ResponseWriter, r *http.Request) {
	phoneNumber := r.URL.Query().Get("phoneNumber")
	var responseBody PublicInfoResponse
	user, err := h.service.GetPublicInfo(phoneNumber)
	if err != nil {
		responseBody.Error = err.Error()
		responseBytes, err := json.Marshal(responseBody)
		if err != nil {
			log.Println(err.Error())
		}
		_, err = w.Write(responseBytes)
		if err != nil {
			log.Println(err.Error())
		}
		return
	}
	responseBody.Data = UserPublicInfo{
		Isu:         user.Isu,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		FullName:    user.FullName,
	}
	responseBytes, err := json.Marshal(responseBody)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(responseBytes)
	if err != nil {
		log.Println(err.Error())
	}
}

func (h *HttpServer) UpdatePublicInfo(w http.ResponseWriter, r *http.Request) {
	idStr := h.ctx.Value("userId").(string)
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "cant parse id", http.StatusInternalServerError)
	}
	var (
		publicInfo   UserPublicInfo
		bodyBytes    bytes.Buffer
		responseBody NewIDResponse
	)
	_, err = bodyBytes.ReadFrom(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		responseBody.Error = err.Error()
		responseBytes, err := json.Marshal(responseBody)
		if err != nil {
			log.Println(err.Error())
		}
		_, err = w.Write(responseBytes)
		if err != nil {
			log.Println(err.Error())
		}
	}
	err = json.Unmarshal(bodyBytes.Bytes(), &publicInfo)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		responseBytes, err := json.Marshal(responseBody)
		if err != nil {
			log.Println(err.Error())
		}
		_, err = w.Write(responseBytes)
		if err != nil {
			log.Println(err.Error())
		}
	}

	newUser := adapters.UserFullData{
		ID:          id,
		Isu:         publicInfo.Isu,
		Email:       publicInfo.Email,
		PhoneNumber: publicInfo.PhoneNumber,
		FullName:    publicInfo.FullName,
	}

	err = h.service.UpdatePublicInfo(newUser)
	if err != nil {
		responseBody.Error = err.Error()
	}
	responseBytes, err := json.Marshal(responseBody)
	if err != nil {
		http.Error(w, "cant transform response to json", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(responseBytes)
	if err != nil {
		log.Println(err.Error())
	}
}

func (h *HttpServer) UpdateFullInfo(w http.ResponseWriter, r *http.Request) {
	idStr := h.ctx.Value("userId").(string)
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "cant parse id", http.StatusInternalServerError)
		return
	}

	role, err := h.service.GetRole(id)
	if err != nil {
		http.Error(w, "no user with such id", http.StatusBadRequest)
		return
	}

	if role != adapters.Admin {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var (
		fullUserInfo adapters.UserFullData
		bodyBytes    bytes.Buffer
		responseBody NewIDResponse
	)
	_, err = bodyBytes.ReadFrom(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		responseBody.Error = err.Error()
		responseBytes, err := json.Marshal(responseBody)
		if err != nil {
			log.Println(err.Error())
		}
		_, err = w.Write(responseBytes)
		if err != nil {
			log.Println(err.Error())
		}
	}
	err = json.Unmarshal(bodyBytes.Bytes(), &fullUserInfo)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		responseBytes, err := json.Marshal(responseBody)
		if err != nil {
			log.Println(err.Error())
		}
		_, err = w.Write(responseBytes)
		if err != nil {
			log.Println(err.Error())
		}
	}

	err = h.service.UpdateFullInfo(fullUserInfo)
	if err != nil {
		responseBody.Error = err.Error()
	}
	responseBytes, err := json.Marshal(responseBody)
	if err != nil {
		http.Error(w, "cant transform response to json", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(responseBytes)
	if err != nil {
		log.Println(err.Error())
	}
}

func (h *HttpServer) AllUserData(w http.ResponseWriter, r *http.Request) {
	idStr := h.ctx.Value("userId").(string)
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "cant parse id", http.StatusInternalServerError)
		return
	}

	role, err := h.service.GetRole(id)
	if err != nil {
		http.Error(w, "no user with such id", http.StatusBadRequest)
		return
	}

	if role != adapters.Admin {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	users, err := h.service.GetAllUsers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonStr, err := json.Marshal(users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonStr)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *HttpServer) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	state := RandStringBytes(16)
	setCallbackCookie(w, r, "state", state)
	http.Redirect(w, r, h.oauth2Config.AuthCodeURL(state), http.StatusFound)
}

func setCallbackCookie(w http.ResponseWriter, r *http.Request, name, value string) {
	c := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   int(time.Hour.Seconds()),
		Secure:   r.TLS != nil,
		HttpOnly: true,
	}
	http.SetCookie(w, c)
}
func (h *HttpServer) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	state, err := r.Cookie("state")
	if err != nil {
		http.Error(w, "state not found", http.StatusBadRequest)
		return
	}
	if r.URL.Query().Get("state") != state.Value {
		http.Error(w, "state didn't match", http.StatusBadRequest)
		return
	}

	oauth2Token, err := h.oauth2Config.Exchange(h.ctx, r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	accessTokenCookie := &http.Cookie{
		Name:     "access_token",
		Value:    oauth2Token.AccessToken,
		Path:     "/",
		MaxAge:   int(time.Hour.Seconds()),
		Secure:   r.TLS != nil,
		HttpOnly: true,
	}
	refreshTokenCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    oauth2Token.RefreshToken,
		Path:     "/",
		MaxAge:   int(time.Hour.Seconds()),
		Secure:   r.TLS != nil,
		HttpOnly: true,
	}
	http.SetCookie(w, accessTokenCookie)
	http.SetCookie(w, refreshTokenCookie)
	fmt.Fprintf(w, "you are authorized, please, use send request on endpoint again")
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func (h *HttpServer) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawAccessToken, err := r.Cookie("access_token")
		if errors.Is(err, http.ErrNoCookie) {
			http.Redirect(w, r, "/auth", http.StatusFound)
			return
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{
			Timeout:   time.Duration(6000) * time.Second,
			Transport: tr,
		}
		ctx := oidc.ClientContext(context.Background(), client)

		oidcConfig := &oidc.Config{
			ClientID:          h.oauth2Config.ClientID,
			SkipClientIDCheck: true,
		}
		verifier := h.provider.Verifier(oidcConfig)
		idToken, err := verifier.Verify(ctx, rawAccessToken.Value)
		if err != nil {
			rawRefreshToken, err := r.Cookie("refresh_token")
			if errors.Is(err, http.ErrNoCookie) {
				http.Redirect(w, r, "/auth", http.StatusFound)
				return
			} else if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			ts := h.oauth2Config.TokenSource(ctx, &oauth2.Token{RefreshToken: rawRefreshToken.Value})
			tok, err := ts.Token()
			if err != nil {
				http.Error(w, err.Error(), http.StatusForbidden)
				return
			}
			idToken, err = verifier.Verify(ctx, tok.AccessToken)
			if err != nil {
				http.Error(w, err.Error(), http.StatusForbidden)
				return
			}
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		h.ctx = context.WithValue(h.ctx, "userId", idToken.Subject)
		next.ServeHTTP(w, r)
	})
}
