package wrapper_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"code.cloudfoundry.org/cli/api/uaa"

	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"

	"code.cloudfoundry.org/cli/api/cloudcontroller/ccerror"

	"code.cloudfoundry.org/cli/api/cloudcontroller"
	"code.cloudfoundry.org/cli/api/cloudcontroller/cloudcontrollerfakes"
	. "code.cloudfoundry.org/cli/api/cloudcontroller/wrapper"
	"code.cloudfoundry.org/cli/api/cloudcontroller/wrapper/wrapperfakes"
	"code.cloudfoundry.org/cli/api/uaa/wrapper/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UAA Authentication", func() {
	var (
		fakeConnection *cloudcontrollerfakes.FakeConnection
		fakeClient     *wrapperfakes.FakeUAAClient
		inMemoryCache  *util.InMemoryCache

		wrapper cloudcontroller.Connection
		request *cloudcontroller.Request
		inner   *UAAAuthentication
	)

	BeforeEach(func() {
		fakeConnection = new(cloudcontrollerfakes.FakeConnection)
		fakeClient = new(wrapperfakes.FakeUAAClient)
		inMemoryCache = util.NewInMemoryTokenCache()
		inner = NewUAAAuthentication(fakeClient, inMemoryCache)
		wrapper = inner.Wrap(fakeConnection)

		request = &cloudcontroller.Request{
			Request: &http.Request{
				Header: http.Header{},
			},
		}
	})

	Describe("Make", func() {
		When("the client is nil", func() {
			BeforeEach(func() {
				inner.SetClient(nil)

				fakeConnection.MakeReturns(ccerror.InvalidAuthTokenError{})
			})

			It("calls the connection without any side effects", func() {
				err := wrapper.Make(request, nil)
				Expect(err).To(MatchError(ccerror.InvalidAuthTokenError{}))

				Expect(fakeClient.RefreshAccessTokenCallCount()).To(Equal(0))
				Expect(fakeConnection.MakeCallCount()).To(Equal(1))
			})
		})

		When("the token is valid", func() {
			var (
				accessToken = buildTokenString(time.Now().AddDate(0, 0, 1))
			)
			BeforeEach(func() {
				inMemoryCache.SetAccessToken(accessToken)
			})

			It("does not refresh the token", func() {
				Expect(fakeClient.RefreshAccessTokenCallCount()).To(Equal(0))
			})

			It("adds authentication headers", func() {
				err := wrapper.Make(request, nil)
				Expect(err).ToNot(HaveOccurred())

				Expect(fakeConnection.MakeCallCount()).To(Equal(1))
				authenticatedRequest, _ := fakeConnection.MakeArgsForCall(0)
				headers := authenticatedRequest.Header
				Expect(headers["Authorization"]).To(ConsistOf([]string{accessToken}))
			})

			When("the request already has headers", func() {
				It("preserves existing headers", func() {
					request.Header.Add("Existing", "header")
					err := wrapper.Make(request, nil)
					Expect(err).ToNot(HaveOccurred())

					Expect(fakeConnection.MakeCallCount()).To(Equal(1))
					authenticatedRequest, _ := fakeConnection.MakeArgsForCall(0)
					headers := authenticatedRequest.Header
					Expect(headers["Existing"]).To(ConsistOf([]string{"header"}))
				})
			})

			When("the wrapped connection returns nil", func() {
				It("returns nil", func() {
					fakeConnection.MakeReturns(nil)

					err := wrapper.Make(request, nil)
					Expect(err).ToNot(HaveOccurred())
				})
			})

			When("the wrapped connection returns an error", func() {
				It("returns the error", func() {
					innerError := errors.New("inner error")
					fakeConnection.MakeReturns(innerError)

					err := wrapper.Make(request, nil)
					Expect(err).To(Equal(innerError))
				})
			})
		})

		When("the token is invalid", func() {
			var (
				expectedBody string
				request      *cloudcontroller.Request
				executeErr   error
			)

			invalidAccessToken := buildTokenString(time.Time{})
			newAccessToken := buildTokenString(time.Now().AddDate(0, 1, 1))
			newRefreshToken := "newRefreshToken"

			BeforeEach(func() {
				expectedBody = "this body content should be preserved"
				body := strings.NewReader(expectedBody)
				request = cloudcontroller.NewRequest(&http.Request{
					Header: http.Header{},
					Body:   ioutil.NopCloser(body),
				}, body)

				inMemoryCache.SetAccessToken(invalidAccessToken)

				fakeClient.RefreshAccessTokenReturns(
					uaa.RefreshedTokens{
						AccessToken:  newAccessToken,
						RefreshToken: newRefreshToken,
						Type:         "bearer",
					},
					nil,
				)
			})

			JustBeforeEach(func() {
				executeErr = wrapper.Make(request, nil)
			})

			It("should refresh the token", func() {
				Expect(executeErr).ToNot(HaveOccurred())
				Expect(fakeClient.RefreshAccessTokenCallCount()).To(Equal(1))
			})

			It("should save the refresh token", func() {
				Expect(inMemoryCache.RefreshToken()).To(Equal(newRefreshToken))
				Expect(inMemoryCache.AccessToken()).To(ContainSubstring(newAccessToken))
			})

			When("token cannot be refreshed", func() {
				JustBeforeEach(func() {
					fakeConnection.MakeReturns(ccerror.InvalidAuthTokenError{})
				})

				It("should not re-try the initial request", func() {
					Expect(fakeConnection.MakeCallCount()).To(Equal(1))
				})
			})

		})
	})
})

func buildTokenString(expiration time.Time) string {
	c := jws.Claims{}
	c.SetExpiration(expiration)
	token := jws.NewJWT(c, crypto.Unsecured)
	tokenBytes, _ := token.Serialize(nil)
	return string(tokenBytes)
}
