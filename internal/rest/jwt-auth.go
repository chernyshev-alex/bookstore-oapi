package rest

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/golang-jwt/jwt/v4"
)

const PublicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA4f5wg5l2hKsTeNem/V41
fGnJm6gOdrj8ym3rFkEU/wT8RDtnSgFEZOQpHEgQ7JL38xUfU0Y3g6aYw9QT0hJ7
mCpz9Er5qLaMXJwZxzHzAahlfA0icqabvJOMvQtzD6uQv6wPEyZtDTWiQi9AXwBp
HssPnpYGIn20ZZuNlX2BrClciHhCPUIIZOQn/MmqTD31jSyjoQoV7MhhMTATKJx2
XrHhR+1DcKJzQBSTAGnpYVaqpsARap+nwRipr3nUTuxyGohBTSmjJ2usSeQXHI3b
ODIRe1AuTyHceAbewn8b462yEWKARdpd9AjQW5SIVPfdsz5B6GlYQ5LdYKtznTuy
7wIDAQAB
-----END PUBLIC KEY-----`

// it must be loaded from an external secure service
const PrivateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA4f5wg5l2hKsTeNem/V41fGnJm6gOdrj8ym3rFkEU/wT8RDtn
SgFEZOQpHEgQ7JL38xUfU0Y3g6aYw9QT0hJ7mCpz9Er5qLaMXJwZxzHzAahlfA0i
cqabvJOMvQtzD6uQv6wPEyZtDTWiQi9AXwBpHssPnpYGIn20ZZuNlX2BrClciHhC
PUIIZOQn/MmqTD31jSyjoQoV7MhhMTATKJx2XrHhR+1DcKJzQBSTAGnpYVaqpsAR
ap+nwRipr3nUTuxyGohBTSmjJ2usSeQXHI3bODIRe1AuTyHceAbewn8b462yEWKA
Rdpd9AjQW5SIVPfdsz5B6GlYQ5LdYKtznTuy7wIDAQABAoIBAQCwia1k7+2oZ2d3
n6agCAbqIE1QXfCmh41ZqJHbOY3oRQG3X1wpcGH4Gk+O+zDVTV2JszdcOt7E5dAy
MaomETAhRxB7hlIOnEN7WKm+dGNrKRvV0wDU5ReFMRHg31/Lnu8c+5BvGjZX+ky9
POIhFFYJqwCRlopGSUIxmVj5rSgtzk3iWOQXr+ah1bjEXvlxDOWkHN6YfpV5ThdE
KdBIPGEVqa63r9n2h+qazKrtiRqJqGnOrHzOECYbRFYhexsNFz7YT02xdfSHn7gM
IvabDDP/Qp0PjE1jdouiMaFHYnLBbgvlnZW9yuVf/rpXTUq/njxIXMmvmEyyvSDn
FcFikB8pAoGBAPF77hK4m3/rdGT7X8a/gwvZ2R121aBcdPwEaUhvj/36dx596zvY
mEOjrWfZhF083/nYWE2kVquj2wjs+otCLfifEEgXcVPTnEOPO9Zg3uNSL0nNQghj
FuD3iGLTUBCtM66oTe0jLSslHe8gLGEQqyMzHOzYxNqibxcOZIe8Qt0NAoGBAO+U
I5+XWjWEgDmvyC3TrOSf/KCGjtu0TSv30ipv27bDLMrpvPmD/5lpptTFwcxvVhCs
2b+chCjlghFSWFbBULBrfci2FtliClOVMYrlNBdUSJhf3aYSG2Doe6Bgt1n2CpNn
/iu37Y3NfemZBJA7hNl4dYe+f+uzM87cdQ214+jrAoGAXA0XxX8ll2+ToOLJsaNT
OvNB9h9Uc5qK5X5w+7G7O998BN2PC/MWp8H+2fVqpXgNENpNXttkRm1hk1dych86
EunfdPuqsX+as44oCyJGFHVBnWpm33eWQw9YqANRI+pCJzP08I5WK3osnPiwshd+
hR54yjgfYhBFNI7B95PmEQkCgYBzFSz7h1+s34Ycr8SvxsOBWxymG5zaCsUbPsL0
4aCgLScCHb9J+E86aVbbVFdglYa5Id7DPTL61ixhl7WZjujspeXZGSbmq0Kcnckb
mDgqkLECiOJW2NHP/j0McAkDLL4tysF8TLDO8gvuvzNC+WQ6drO2ThrypLVZQ+ry
eBIPmwKBgEZxhqa0gVvHQG/7Od69KWj4eJP28kq13RhKay8JOoN0vPmspXJo1HY3
CKuHRG+AP579dncdUnOMvfXOtkdM4vk0+hWASBQzM9xzVcztCa+koAugjVaLS9A+
9uQoqEeVNTckxx0S2bYevRy7hGQmUJTyQm3j1zEUR5jpdbL83Fbq
-----END RSA PRIVATE KEY-----`

// TODO move to vault or service
type JWSValidator interface {
	ValidateJWS(jws string) (*jwt.Token, error)
}

type Validator struct {
	PrivateKey *rsa.PrivateKey
	PubKey     *rsa.PublicKey
}

type UserInfo struct {
}

type UserClaim struct {
	*jwt.StandardClaims
	TokenType string
	UserInfo
}

var ErrNoAuthHeader = errors.New("authorization header is missing")
var ErrInvalidAuthHeader = errors.New("authorization header is malformed")
var ErrClaimsInvalid = errors.New("provided claims do not match expected scopes")

var _ JWSValidator = (*Validator)(nil)

func NewValidator() (*Validator, error) {
	pub := LoadPublicRSAFKey()
	priv := LoadPrivateRSAKey()
	return &Validator{PrivateKey: priv, PubKey: pub}, nil
}

func NewAuthenticator(v JWSValidator) openapi3filter.AuthenticationFunc {
	// create the token to pretend we have logged in already
	tokenStr := CreateToken(jwt.MapClaims{"foo": "bar", "c1": "rw"})
	fmt.Printf("debug:token:%s\n", tokenStr)

	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		return Authenticate(v, ctx, input)
	}
}

func (v *Validator) ValidateJWS(tokenStr string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return v.PubKey, nil
	})
	return token, err
}

func Authenticate(v JWSValidator, ctx context.Context, input *openapi3filter.AuthenticationInput) error {
	if input.SecuritySchemeName != "BearerAuth" {
		return fmt.Errorf("security scheme %s != 'BearerAuth'", input.SecuritySchemeName)
	}

	jws, err := GetJWSFromRequest(input.RequestValidationInput.Request)
	if err != nil {
		return fmt.Errorf("getting jws: %w", err)
	}

	token, err := v.ValidateJWS(jws)
	if err != nil {
		return fmt.Errorf("validating JWS: %w", err)
	}

	err = CheckTokenClaims(input.Scopes, *token)
	if err != nil {
		return fmt.Errorf("token claims don't match: %w", err)
	}
	return nil
}

func GetJWSFromRequest(req *http.Request) (string, error) {
	authHdr := req.Header.Get("Authorization")
	if authHdr == "" {
		return "", ErrNoAuthHeader
	}
	const prefix = "Bearer "
	if !strings.HasPrefix(authHdr, prefix) {
		return "", ErrInvalidAuthHeader
	}
	return strings.TrimPrefix(authHdr, prefix), nil
}

func GetClaimsFromToken(t jwt.Token) map[string]string {
	result := make(map[string]string)
	if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		for tag, v := range claims {
			result[tag] = v.(string)
		}
	}
	return result
}

func CheckTokenClaims(expectedClaims []string, t jwt.Token) error {
	claims := GetClaimsFromToken(t)
	for _, e := range expectedClaims {
		if _, ok := claims[e]; !ok {
			return ErrClaimsInvalid
		}
	}
	return nil
}

// returns a encoded JWT token that has been signed with the specified cryptographic key.
func CreateToken(claims jwt.Claims) string {
	// actually load from secure place
	privKey := LoadPrivateRSAKey()

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenStr, e := token.SignedString(privKey)
	if e != nil {
		panic(e.Error())
	}
	return tokenStr
}

func LoadPrivateRSAKey() *rsa.PrivateKey {
	privKey, e := jwt.ParseRSAPrivateKeyFromPEM([]byte(PrivateKey))
	if e != nil {
		panic(e.Error())
	}
	return privKey
}

func LoadPublicRSAFKey() *rsa.PublicKey {
	pubKey, e := jwt.ParseRSAPublicKeyFromPEM([]byte(PublicKey))
	if e != nil {
		panic(e.Error())
	}
	return pubKey
}
