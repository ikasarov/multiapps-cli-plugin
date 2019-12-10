// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"
	"os"
	"time"

	"golang.org/x/net/context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"
)

// NewUploadMtaFileParams creates a new UploadMtaFileParams object
// with the default values initialized.
func NewUploadMtaFileParams() *UploadMtaFileParams {
	return &UploadMtaFileParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewUploadMtaFileParamsWithTimeout creates a new UploadMtaFileParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewUploadMtaFileParamsWithTimeout(timeout time.Duration) *UploadMtaFileParams {
	return &UploadMtaFileParams{
		timeout: timeout,
	}
}

// NewUploadMtaFileParamsWithContext creates a new UploadMtaFileParams object
// with the default values initialized, and the ability to set a context for a request
func NewUploadMtaFileParamsWithContext(ctx context.Context) *UploadMtaFileParams {
	return &UploadMtaFileParams{
		Context: ctx,
	}
}

// NewUploadMtaFileParamsWithHTTPClient creates a new UploadMtaFileParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewUploadMtaFileParamsWithHTTPClient(client *http.Client) *UploadMtaFileParams {
	return &UploadMtaFileParams{
		HTTPClient: client,
	}
}

/*UploadMtaFileParams contains all the parameters to send to the API endpoint
for the upload mta file operation typically these are written to a http.Request
*/
type UploadMtaFileParams struct {

	/*Namespace
	  file namespace

	*/
	Namespace *string

	/*File*/
	File os.File

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the upload mta file params
func (o *UploadMtaFileParams) WithTimeout(timeout time.Duration) *UploadMtaFileParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the upload mta file params
func (o *UploadMtaFileParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the upload mta file params
func (o *UploadMtaFileParams) WithContext(ctx context.Context) *UploadMtaFileParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the upload mta file params
func (o *UploadMtaFileParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the upload mta file params
func (o *UploadMtaFileParams) WithHTTPClient(client *http.Client) *UploadMtaFileParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the upload mta file params
func (o *UploadMtaFileParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithNamespace adds the namespace to the upload mta file params
func (o *UploadMtaFileParams) WithNamespace(namespace *string) *UploadMtaFileParams {
	o.SetNamespace(namespace)
	return o
}

// SetNamespace adds the namespace to the upload mta file params
func (o *UploadMtaFileParams) SetNamespace(namespace *string) {
	o.Namespace = namespace
}

// WithFile adds the file to the upload mta file params
func (o *UploadMtaFileParams) WithFile(file os.File) *UploadMtaFileParams {
	o.SetFile(file)
	return o
}

// SetFile adds the file to the upload mta file params
func (o *UploadMtaFileParams) SetFile(file os.File) {
	o.File = file
}

// WriteToRequest writes these params to a swagger request
func (o *UploadMtaFileParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Namespace != nil {

		// query param namespace
		var qrNamespace string
		if o.Namespace != nil {
			qrNamespace = *o.Namespace
		}
		qNamespace := qrNamespace
		if qNamespace != "" {
			if err := r.SetQueryParam("namespace", qNamespace); err != nil {
				return err
			}
		}

	}

	// form file param file
	if err := r.SetFileParam("file", &o.File); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
