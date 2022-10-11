package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"
	"utilware/dep/jwt"
	"utilware/tig"
)

var (
	defaultOptions = &JwtMiddlewareOptions{
		Secret: []byte("tig"),
		TokenMiddleware: func(c *tig.Context) ([]byte, error) {
			switch {
			case c.Cookie("token") != "":
				return []byte(c.Cookie("token")), nil
			case c.Req.Header.Get("Authorization") != "":
				return []byte(strings.ReplaceAll(c.Req.Header.Get("Authorization"),
					"Bearer ", "")), nil
			case c.PostForm("token") != "":
				return []byte(c.PostForm("token")), nil
			case c.Query("token") != "":
				return []byte(c.Query("token")), nil
			default:
				return nil, errors.New("token not found")
			}
		},
		VerifyErrorMiddleware: func(c *tig.Context, e error) {
			c.String(http.StatusUnauthorized, "401 Unauthorized: "+e.Error())
		},
		HmacLength: 256,
		MaxAge:     time.Duration(1209600) * time.Second,
	}
)

type JwtMiddleware struct {
	secret                []byte
	tokenMiddleware       func(c *tig.Context) ([]byte, error)
	verifyErrorMiddleware func(c *tig.Context, e error)
	hmacAlg               jwt.Alg
	signOption            []jwt.SignOption
}

type JwtMiddlewareOptions struct {
	Secret                []byte                               // 密钥 (默认: "tig")
	TokenMiddleware       func(c *tig.Context) ([]byte, error) // 获取Token字符串 (默认: 从cookie: token, header: Authorization, postForm: token, query: token中获取)
	VerifyErrorMiddleware func(c *tig.Context, e error)        // 验证失败时的处理 (默认: 401 结束请求)
	HmacLength            int                                  // Hmac取模长度 [sha256, sha384, sha512] (默认256)
	MaxAge                time.Duration                        // Token过期时间 (默认: 14天)
	SignOption            []jwt.SignOption                     // 其他签名参数
}

type JwtVerifiedToken struct {
	jwtMiddleware *JwtMiddleware
	jwtContent    []byte
	jwtToken      *jwt.VerifiedToken
}

func NewJwt(options *JwtMiddlewareOptions) *JwtMiddleware {
	if options == nil {
		options = defaultOptions
	}

	if options.Secret == nil {
		options.Secret = defaultOptions.Secret
	}
	if options.TokenMiddleware == nil {
		options.TokenMiddleware = defaultOptions.TokenMiddleware
	}
	if options.VerifyErrorMiddleware == nil {
		options.VerifyErrorMiddleware = defaultOptions.VerifyErrorMiddleware
	}
	if options.HmacLength <= 0 {
		options.HmacLength = defaultOptions.HmacLength
	}
	if options.MaxAge <= 0 {
		options.MaxAge = defaultOptions.MaxAge
	}
	if options.SignOption == nil {
		options.SignOption = []jwt.SignOption{}
	}

	j := &JwtMiddleware{
		secret:                options.Secret,
		tokenMiddleware:       options.TokenMiddleware,
		verifyErrorMiddleware: options.VerifyErrorMiddleware,
		signOption:            options.SignOption,
	}

	switch options.HmacLength {
	case 384:
		j.hmacAlg = jwt.HS384
	case 512:
		j.hmacAlg = jwt.HS512
	default: // And HS256
		j.hmacAlg = jwt.HS256
	}

	j.signOption = append(j.signOption, jwt.MaxAge(options.MaxAge))

	return j
}

func (j *JwtMiddleware) Sign(i interface{}) (*JwtVerifiedToken, error) {
	tokenContent, e := jwt.Sign(j.hmacAlg, j.secret, i, j.signOption...)
	if e != nil {
		return nil, e
	}

	jwtToken, e := jwt.Verify(j.hmacAlg, j.secret, tokenContent)
	if e != nil {
		return nil, e
	}

	return &JwtVerifiedToken{
		jwtMiddleware: j,
		jwtContent:    tokenContent,
		jwtToken:      jwtToken,
	}, nil
}

func (j *JwtMiddleware) JwtVerifyMiddleware() tig.HandlerFunc {
	return func(c *tig.Context) {
		v, e := j.JwtVerify(c)
		if e != nil {
			j.verifyErrorMiddleware(c, e)
			c.End()
			return
		}

		c.SetStore("jwt", v)
		c.Next()
	}
}

func (j *JwtMiddleware) JwtVerifyHandler(handler func(c *tig.Context, v *JwtVerifiedToken, e error) bool) tig.HandlerFunc {
	return func(c *tig.Context) {
		v, e := j.JwtVerify(c)
		if !handler(c, v, e) {
			c.End()
			return
		}

		c.SetStore("jwt", v)
		c.Next()
	}
}

func JwtVerifed(c *tig.Context) *JwtVerifiedToken {
	sv, ok := c.GetStore("jwt")
	if ok && sv != nil {
		if v, ok := sv.(*JwtVerifiedToken); ok {
			return v
		}
	}

	return nil
}

func (j *JwtMiddleware) JwtVerify(c *tig.Context) (*JwtVerifiedToken, error) {
	tokenContent, e := j.tokenMiddleware(c)
	if e != nil {
		return nil, e
	}

	jwtToken, e := jwt.Verify(j.hmacAlg, j.secret, tokenContent)
	if e != nil {
		return nil, e
	}

	if jwtToken.StandardClaims.Expiry < time.Now().Unix() {
		return nil, errors.New("token time expired")
	}

	return &JwtVerifiedToken{
		jwtMiddleware: j,
		jwtContent:    tokenContent,
		jwtToken:      jwtToken,
	}, nil
}

// 解析Token
func (v *JwtVerifiedToken) Claims(dest interface{}) error {
	return v.jwtToken.Claims(dest)
}

// 用Map解析Token
func (v *JwtVerifiedToken) ClaimsMap() (map[string]interface{}, error) {
	var m map[string]interface{}
	return m, v.jwtToken.Claims(&m)
}

// Token过期时间
func (v *JwtVerifiedToken) Expiry() int64 {
	return v.jwtToken.StandardClaims.Expiry
}

// Token Payload
func (v *JwtVerifiedToken) Payload() []byte {
	return v.jwtToken.Payload
}

func (v *JwtVerifiedToken) Token() []byte {
	return v.jwtContent
}

func (v *JwtVerifiedToken) String() string {
	return string(v.jwtContent)
}
