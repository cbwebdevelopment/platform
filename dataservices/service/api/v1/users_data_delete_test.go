package v1_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"net/http"

	"github.com/tidepool-org/platform/app"
	"github.com/tidepool-org/platform/dataservices/service/api/v1"
	"github.com/tidepool-org/platform/service"
)

var _ = Describe("UsersDataDelete", func() {
	Context("Unit Tests", func() {
		var targetUserID string
		var context *TestContext

		BeforeEach(func() {
			targetUserID = app.NewID()
			context = NewTestContext()
			context.RequestImpl.PathParams["userid"] = targetUserID
			context.AuthenticationDetailsImpl.IsServerOutputs = []bool{true}
			context.DataStoreSessionImpl.DestroyDataForUserByIDOutputs = []error{nil}
			context.TaskStoreSessionImpl.DestroyTasksForUserByIDOutputs = []error{nil}
			context.MetricServicesClientImpl.RecordMetricOutputs = []error{nil}
		})

		It("succeeds if authenticated as server", func() {
			v1.UsersDataDelete(context)
			Expect(context.DataStoreSessionImpl.DestroyDataForUserByIDInputs).To(Equal([]string{targetUserID}))
			Expect(context.TaskStoreSessionImpl.DestroyTasksForUserByIDInputs).To(Equal([]string{targetUserID}))
			Expect(context.MetricServicesClientImpl.RecordMetricInputs).To(Equal([]RecordMetricInput{{context, "users_data_delete", nil}}))
			Expect(context.RespondWithStatusAndDataInputs).To(Equal([]RespondWithStatusAndDataInput{{http.StatusOK, []struct{}{}}}))
			Expect(context.ValidateTest()).To(BeTrue())
		})

		It("panics if context is missing", func() {
			context.AuthenticationDetailsImpl.IsServerOutputs = []bool{}
			context.DataStoreSessionImpl.DestroyDataForUserByIDOutputs = []error{}
			context.TaskStoreSessionImpl.DestroyTasksForUserByIDOutputs = []error{}
			context.MetricServicesClientImpl.RecordMetricOutputs = []error{}
			Expect(func() { v1.UsersDataDelete(nil) }).To(Panic())
			Expect(context.ValidateTest()).To(BeTrue())
		})

		It("panics if request is missing", func() {
			context.RequestImpl = nil
			context.AuthenticationDetailsImpl.IsServerOutputs = []bool{}
			context.DataStoreSessionImpl.DestroyDataForUserByIDOutputs = []error{}
			context.TaskStoreSessionImpl.DestroyTasksForUserByIDOutputs = []error{}
			context.MetricServicesClientImpl.RecordMetricOutputs = []error{}
			Expect(func() { v1.UsersDataDelete(nil) }).To(Panic())
			Expect(context.ValidateTest()).To(BeTrue())
		})

		It("responds with error if user id not provided as a parameter", func() {
			delete(context.RequestImpl.PathParams, "userid")
			context.AuthenticationDetailsImpl.IsServerOutputs = []bool{}
			context.DataStoreSessionImpl.DestroyDataForUserByIDOutputs = []error{}
			context.TaskStoreSessionImpl.DestroyTasksForUserByIDOutputs = []error{}
			context.MetricServicesClientImpl.RecordMetricOutputs = []error{}
			v1.UsersDataDelete(context)
			Expect(context.RespondWithErrorInputs).To(Equal([]*service.Error{v1.ErrorUserIDMissing()}))
			Expect(context.ValidateTest()).To(BeTrue())
		})

		It("panics if authentication details is missing", func() {
			context.AuthenticationDetailsImpl = nil
			context.DataStoreSessionImpl.DestroyDataForUserByIDOutputs = []error{}
			context.TaskStoreSessionImpl.DestroyTasksForUserByIDOutputs = []error{}
			context.MetricServicesClientImpl.RecordMetricOutputs = []error{}
			Expect(func() { v1.UsersDataDelete(context) }).To(Panic())
			Expect(context.ValidateTest()).To(BeTrue())
		})

		It("responds with error if not server", func() {
			context.AuthenticationDetailsImpl.IsServerOutputs = []bool{false}
			context.DataStoreSessionImpl.DestroyDataForUserByIDOutputs = []error{}
			context.TaskStoreSessionImpl.DestroyTasksForUserByIDOutputs = []error{}
			context.MetricServicesClientImpl.RecordMetricOutputs = []error{}
			v1.UsersDataDelete(context)
			Expect(context.RespondWithErrorInputs).To(Equal([]*service.Error{service.ErrorUnauthorized()}))
			Expect(context.ValidateTest()).To(BeTrue())
		})

		It("panics if data store session is missing", func() {
			context.DataStoreSessionImpl = nil
			context.TaskStoreSessionImpl.DestroyTasksForUserByIDOutputs = []error{}
			context.MetricServicesClientImpl.RecordMetricOutputs = []error{}
			Expect(func() { v1.UsersDataDelete(context) }).To(Panic())
			Expect(context.ValidateTest()).To(BeTrue())
		})

		It("responds with error if data store session delete data for user by id returns error", func() {
			err := errors.New("other")
			context.DataStoreSessionImpl.DestroyDataForUserByIDOutputs = []error{err}
			context.TaskStoreSessionImpl.DestroyTasksForUserByIDOutputs = []error{}
			context.MetricServicesClientImpl.RecordMetricOutputs = []error{}
			v1.UsersDataDelete(context)
			Expect(context.RespondWithInternalServerFailureInputs).To(Equal([]RespondWithInternalServerFailureInput{{"Unable to delete data for user by id", []interface{}{err}}}))
			Expect(context.ValidateTest()).To(BeTrue())
		})

		It("panics if tasks store session is missing", func() {
			context.TaskStoreSessionImpl = nil
			context.MetricServicesClientImpl.RecordMetricOutputs = []error{}
			Expect(func() { v1.UsersDataDelete(context) }).To(Panic())
			Expect(context.ValidateTest()).To(BeTrue())
		})

		It("responds with error if task store session delete tasks for user by id returns error", func() {
			err := errors.New("other")
			context.TaskStoreSessionImpl.DestroyTasksForUserByIDOutputs = []error{err}
			context.MetricServicesClientImpl.RecordMetricOutputs = []error{}
			v1.UsersDataDelete(context)
			Expect(context.RespondWithInternalServerFailureInputs).To(Equal([]RespondWithInternalServerFailureInput{{"Unable to delete tasks for user by id", []interface{}{err}}}))
			Expect(context.ValidateTest()).To(BeTrue())
		})

		It("panics if metric services client is missing", func() {
			context.MetricServicesClientImpl = nil
			Expect(func() { v1.UsersDataDelete(context) }).To(Panic())
			Expect(context.DataStoreSessionImpl.DestroyDataForUserByIDInputs).To(Equal([]string{targetUserID}))
			Expect(context.ValidateTest()).To(BeTrue())
		})

		It("logs and ignores if metric services record metric returns an error", func() {
			context.MetricServicesClientImpl.RecordMetricOutputs = []error{errors.New("other")}
			v1.UsersDataDelete(context)
			Expect(context.DataStoreSessionImpl.DestroyDataForUserByIDInputs).To(Equal([]string{targetUserID}))
			Expect(context.MetricServicesClientImpl.RecordMetricInputs).To(Equal([]RecordMetricInput{{context, "users_data_delete", nil}}))
			Expect(context.RespondWithStatusAndDataInputs).To(Equal([]RespondWithStatusAndDataInput{{http.StatusOK, []struct{}{}}}))
			Expect(context.ValidateTest()).To(BeTrue())
		})
	})
})
