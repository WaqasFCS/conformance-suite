const { obtainClientCredentialsAccessToken } = require('../authorise');
const { verifyHeaders, postPayments } = require('./payments');
const { buildPaymentsData } = require('./payment-data-builder');
const { persistPaymentDetails } = require('./persistence');
const debug = require('debug')('debug');

const createRequest = async (resourcePath, headers, paymentData) => {
  verifyHeaders(headers);
  const apiVersion = headers.config.api_version;
  const response = await postPayments(
    resourcePath,
    `/open-banking/v${apiVersion}/payments`,
    headers,
    paymentData,
  );
  let error;
  if (response.Data) {
    const status = response.Data.Status;
    debug(`/payments repsonse Data: ${JSON.stringify(response.Data)}`);
    if (status === 'AcceptedTechnicalValidation' || status === 'AcceptedCustomerProfile') {
      if (response.Data.PaymentId) {
        return response.Data.PaymentId;
      }
    } else {
      error = new Error(`Payment response status: "${status}"`);
      error.status = 500;
      throw error;
    }
  }
  error = new Error('Payment response missing payload');
  error.status = 500;
  throw error;
};

exports.setupPayment = async (authorisationServerId,
  headers, CreditorAccount, InstructedAmount) => {
  const { config } = headers;
  const accessToken = await obtainClientCredentialsAccessToken(config);

  const paymentData = buildPaymentsData(
    {}, // opts
    {}, // risk
    CreditorAccount, InstructedAmount,
  );

  const headersWithToken = Object.assign({ accessToken }, headers);
  const paymentId = await createRequest(
    config.resource_endpoint,
    headersWithToken,
    paymentData,
  );

  const fullPaymentData = {
    Data: {
      PaymentId: paymentId,
      Initiation: paymentData.Data.Initiation,
    },
    Risk: paymentData.Risk,
  };

  persistPaymentDetails(headers.interactionId, fullPaymentData);

  return paymentId;
};