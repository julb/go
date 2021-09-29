package http

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"encoding/base64"

	"github.com/hashicorp/go-retryablehttp"
	cc "github.com/julb/go/pkg/context"
	"github.com/julb/go/pkg/logging"
	"github.com/julb/go/pkg/tracing"
	"github.com/opentracing/opentracing-go"
)

type HttpClient struct{}

type HttpClientRetryOpts struct {
	WaitMin     time.Duration // Minimum time to wait
	WaitMax     time.Duration // Maximum time to wait
	Max         int           // Maximum number of retries
	StatusCodes []int         // status code on which retries should be done
}

type HttpClientTlsOpts struct {
	Insecure   bool
	Truststore *HttpClientTlsTruststoreOpts
	Keystore   *HttpClientTlsKeystoreOpts
}

type HttpClientTlsTruststoreOpts struct {
	CaCertificate string
}

type HttpClientTlsKeystoreOpts struct {
	Certificate    string
	CertificateKey string
}

type HttpClientOpts struct {
	Method      string
	Url         string
	QueryParams map[string]string
	Headers     map[string][]string
	Body        *bytes.Buffer
	Context     context.Context
	Retry       *HttpClientRetryOpts
	Tls         *HttpClientTlsOpts
	logger      *logging.LogWithContext
}

type HttpClientResponse struct {
	RequestOpts *HttpClientOpts
	StatusCode  int
	Headers     map[string][]string
	Body        []byte
}

func (client *HttpClient) Do(opts *HttpClientOpts) (*HttpClientResponse, error) {
	opts.logger.Trace("http client: start configuration")

	// build a retryable client
	httpClient := retryableHttpClient(opts)

	if opts.Tls != nil {
		tlsConfig, err := buildTlsClientConfig(opts)
		if err != nil {
			return nil, err
		}

		httpClient.Transport = &http.Transport{
			TLSClientConfig: tlsConfig,
		}
	}

	// build a request context
	var httpRequestWithContext *http.Request
	var err error
	if opts.Body != nil {
		httpRequestWithContext, err = http.NewRequestWithContext(opts.Context, opts.Method, opts.Url, opts.Body)
	} else {
		httpRequestWithContext, err = http.NewRequestWithContext(opts.Context, opts.Method, opts.Url, nil)
	}
	if err != nil {
		return nil, err
	}

	// add query params if any
	q := httpRequestWithContext.URL.Query()
	for key, value := range opts.QueryParams {
		q.Add(key, value)
	}

	// add headers if any
	for key, values := range opts.Headers {
		for _, value := range values {
			httpRequestWithContext.Header.Add(key, value)
		}
	}

	// propagate context
	propagateContext(httpRequestWithContext)

	// Trace
	opts.logger.Trace("http client: complete configuration")
	opts.logger.Trace("http client: doing call to server")

	// execute request
	httpResponse, err := httpClient.Do(httpRequestWithContext)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	opts.logger.Trace("http client: call to server complete")

	// return response
	return buildCustomHttpResponse(opts, httpResponse)
}

// Default http client opts
func NewHttpClientOpts() *HttpClientOpts {
	return &HttpClientOpts{
		Retry:  NewHttpClientRetryOpts(),
		logger: logging.WithEmptyContext(),
	}
}

// Default http client retry opts
func NewHttpClientRetryOpts() *HttpClientRetryOpts {
	return &HttpClientRetryOpts{
		WaitMin: 1 * time.Second,
		WaitMax: 30 * time.Second,
		Max:     4,
	}
}

// Default http client opts with context
func NewHttpClientOptsWithContext(context context.Context) *HttpClientOpts {
	// USe default constructor
	opts := NewHttpClientOpts()
	opts.Context = context

	// Get logger
	contextualLogger := cc.GetCtxLogger(context)
	if contextualLogger != nil {
		opts.logger = contextualLogger
	}

	// Return opts.
	return opts
}

// Sets the URL
func (opts *HttpClientOpts) WithUrl(url string) *HttpClientOpts {
	opts.Url = url
	return opts
}

// Sets the method
func (opts *HttpClientOpts) WithMethod(method string) *HttpClientOpts {
	opts.Method = method
	return opts
}

// Sets the query param
func (opts *HttpClientOpts) WithQueryParam(key string, value string) *HttpClientOpts {
	if opts.QueryParams == nil {
		opts.QueryParams = make(map[string]string)
	}
	opts.QueryParams[key] = value
	return opts
}

// Sets the header
func (opts *HttpClientOpts) WithHeader(key string, value string) *HttpClientOpts {
	if opts.Headers == nil {
		opts.Headers = make(map[string][]string)
	}
	if opts.Headers[key] == nil {
		opts.Headers[key] = []string{value}
	} else {
		opts.Headers[key] = append(opts.Headers[key], value)
	}
	return opts
}

// Sets the basic auth header
func (opts *HttpClientOpts) WithAuthorizationBasic(username string, password string) *HttpClientOpts {
	value := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password))))
	return opts.WithHeader(HdrAuthorization, value)
}

// Sets the bearer auth header
func (opts *HttpClientOpts) WithAuthorizationBearer(value string) *HttpClientOpts {
	return opts.WithHeader(HdrAuthorization, fmt.Sprintf("Bearer %s", value))
}

// Sets the content-type json
func (opts *HttpClientOpts) WithAcceptJson() *HttpClientOpts {
	return opts.WithHeader(HdrAccept, "application/json")
}

// Sets the content-type json
func (opts *HttpClientOpts) WithContentTypeJson() *HttpClientOpts {
	return opts.WithHeader(HdrContentType, "application/json")
}

// Sets the content-type json and charset utf8
func (opts *HttpClientOpts) WithContentTypeJsonUtf8() *HttpClientOpts {
	return opts.WithHeader(HdrContentType, "application/json;charset=utf8")
}

// Sets the trademark
func (opts *HttpClientOpts) WithTrademark(trademark string) *HttpClientOpts {
	return opts.WithHeader(HdrXJ3Tm, trademark)
}

// Sets the body
func (opts *HttpClientOpts) WithBody(body *bytes.Buffer) *HttpClientOpts {
	opts.Body = body
	return opts
}

// Sets the TLS config
func (opts *HttpClientOpts) WithTLSConfiguration(tlsOpts *HttpClientTlsOpts) *HttpClientOpts {
	opts.Tls = tlsOpts
	return opts
}

// Returns true if the response is 1xx.
func (httpClientResponse *HttpClientResponse) Is1xx() bool {
	return httpClientResponse.StatusCode >= 100 && httpClientResponse.StatusCode < 200
}

// Returns true if the response is 2xx.
func (httpClientResponse *HttpClientResponse) Is2xx() bool {
	return httpClientResponse.StatusCode >= 200 && httpClientResponse.StatusCode < 300
}

// Returns true if the response is 4xx.
func (httpClientResponse *HttpClientResponse) Is3xx() bool {
	return httpClientResponse.StatusCode >= 300 && httpClientResponse.StatusCode < 400
}

// Returns true if the response is 4xx.
func (httpClientResponse *HttpClientResponse) Is4xx() bool {
	return httpClientResponse.StatusCode >= 400 && httpClientResponse.StatusCode < 500
}

// Returns true if the response is 5xx.
func (httpClientResponse *HttpClientResponse) Is5xx() bool {
	return httpClientResponse.StatusCode >= 500 && httpClientResponse.StatusCode < 600
}

// Returns true if the response is 5xx.
func (httpClientResponse *HttpClientResponse) IsError() bool {
	return httpClientResponse.Is4xx() || httpClientResponse.Is5xx()
}

// Returns true if the response is 5xx.
func (httpClientResponse *HttpClientResponse) BodyAsString() string {
	return string(httpClientResponse.Body)
}

// Returns the error message as a string.
func (httpClientResponse *HttpClientResponse) Error() string {
	if httpClientResponse.IsError() {
		return fmt.Sprintf("Got %d: %s", httpClientResponse.StatusCode, httpClientResponse.BodyAsString())
	} else {
		return "unknown error"
	}
}

func buildCustomHttpResponse(opts *HttpClientOpts, response *http.Response) (*HttpClientResponse, error) {
	// Get body as bytes
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// Get headers as a map of string,string
	var headersMap = make(map[string][]string)
	for name, values := range response.Header {
		headersMap[strings.ToLower(name)] = values
	}

	// build a http client response
	var httpClientResponse = &HttpClientResponse{
		RequestOpts: opts,
		StatusCode:  response.StatusCode,
		Body:        bodyBytes,
		Headers:     headersMap,
	}

	// if the request has failed, returns an error.
	if httpClientResponse.IsError() {
		return httpClientResponse, httpClientResponse
	}

	// return the result.
	return httpClientResponse, nil
}

func buildTlsClientConfig(opts *HttpClientOpts) (*tls.Config, error) {
	opts.logger.Trace("http client: tls - start configuration")

	tlsConfig := &tls.Config{}

	if opts.Tls.Truststore != nil {
		if opts.Tls.Truststore.CaCertificate != "" {
			// Read CA certificates
			caCertificatesBytes, err := ioutil.ReadFile(opts.Tls.Truststore.CaCertificate)
			if err != nil {
				opts.logger.Errorf("http client: tls - unable to read ca certificate: %s", err.Error())
				return nil, err
			}

			// Create a cert pool
			rootCAs := x509.NewCertPool()
			ok := rootCAs.AppendCertsFromPEM(caCertificatesBytes)
			if ok {
				opts.logger.Debugf("http client: tls - ca certificates successfully retrieved from %s", opts.Tls.Truststore.CaCertificate)
			} else {
				opts.logger.Warnf("http client: tls - no valid certificates to import from %s", opts.Tls.Truststore.CaCertificate)
			}
			tlsConfig.RootCAs = rootCAs
		}
	}

	if opts.Tls.Insecure {
		opts.logger.Warn("http client: tls - connection configured as insecure")
		tlsConfig.InsecureSkipVerify = true
	}

	opts.logger.Trace("http client: tls - configuration completed")

	return tlsConfig, nil
}

// Apply the context of the request as headers.
func propagateContext(r *http.Request) {
	// propagate x-request-id
	r.Header.Set(HdrXRequestId, GetCtxRequestId(r))

	// propagate opentracing span
	if tracing.IsTracerConfigured() {
		if span := opentracing.SpanFromContext(r.Context()); span != nil {
			err := opentracing.GlobalTracer().Inject(
				span.Context(),
				opentracing.HTTPHeaders,
				opentracing.HTTPHeadersCarrier(r.Header))
			if err != nil {
				GetCtxLogger(r).Errorf("fail to inject opentracing context: %s", err.Error())
			}
		}
	}
}

// Builds a retryable client
func retryableHttpClient(opts *HttpClientOpts) *http.Client {
	// build a client
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = opts.Retry.Max
	retryClient.RetryWaitMin = opts.Retry.WaitMin
	retryClient.RetryWaitMax = opts.Retry.WaitMax
	retryClient.Logger = &retryableHttpClientLoggerProxy{
		logger: opts.logger,
	}
	retryClient.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		if len(opts.Retry.StatusCodes) == 0 {
			return true, nil
		}
		for _, statusCode := range opts.Retry.StatusCodes {
			if statusCode == resp.StatusCode {
				return true, nil
			}
		}
		return false, nil
	}
	retryClient.ErrorHandler = retryablehttp.PassthroughErrorHandler

	// Get standard HTTP client back
	return retryClient.StandardClient()
}

// Wrapper for retryable logger
type retryableHttpClientLoggerProxy struct {
	logger *logging.LogWithContext
}

func (proxy *retryableHttpClientLoggerProxy) Error(msg string, keysAndValues ...interface{}) {
	proxy.logger.Errorf(msg, keysAndValues...)
}

func (proxy *retryableHttpClientLoggerProxy) Warn(msg string, keysAndValues ...interface{}) {
	proxy.logger.Warnf(msg, keysAndValues...)
}

func (proxy *retryableHttpClientLoggerProxy) Info(msg string, keysAndValues ...interface{}) {
	proxy.logger.Infof(msg, keysAndValues...)
}

func (proxy *retryableHttpClientLoggerProxy) Debug(msg string, keysAndValues ...interface{}) {
	proxy.logger.Debugf(msg, keysAndValues...)
}
