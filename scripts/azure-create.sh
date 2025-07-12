#!/bin/bash

# ====== User-defined variables ======
APP_NAME="prince-group-backend"
RESOURCE_GROUP="Prince_Group"
ENV_NAME="managedEnvironment-PrinceGroup-813c"
ACR_SERVER="princegroupbackend-dybue9b2fxdkb0au.azurecr.io"
IMAGE_NAME="prince-group-backend:0.5"
PORT=8000

# Secrets
DATABASE_URL="postgres://prince_group_pg:postgres-pg%40123@prince-group-pg.postgres.database.azure.com/postgres"
FIREBASE_PRIVATE_KEY=$(cat <<EOF
-----BEGIN PRIVATE KEY-----\nMIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQC0IQA9VJTNA+Ei\nzJQRzvZWT67lUnwyq50JxaPVhZGYzKHVl4AzctZ/50l3lQK+dbF/+kgo1HTGN7kt\nos2VhTMi361ry9UF0jMwmnSy8E5+EuNmDgL5QqycuPdmAogEHBi7TAtBwPd9Ku1w\ngSI28oByelQsIu46FzCivrReqz4A1r94Ax0SLln21iaTstMUgF9BJm92o/rec846\n4tHgvUi5TUodz3AxrjF9IFK8WqZe99DT0VmGAAz3YLxichyB35PFSxf+RZxkfzt7\nYY07EgGDtuOydDVqSJVRF+13mZUjRZQxTdG+SblCQ1vA1QTQtH/n55La620jiiwV\nYLcpgDo5AgMBAAECgf8mPRci/j/vl3eqLj5WfpTwMtFkffl1cRxsDcuPrTQRJ4Bp\n+W8bs18fTVVlKfo2CkCleV+vViPxUmMBWV6lHyXNlzk4HMqSjl/7JoKc2/yR1J72\n+ZLWVK8W3Ovpfp1ZG98vFKKFmFzmtsqjao/65mz5Xa+0eHR+7hjr2xAgShaXWrk8\ngyNHpS36N0m6Z5dTjGmvzrvjYRtZALE4NFdBfSfnfaEsQ1w3n1turA8hCklHfwi9\nz06+mXI9JZLlj1dmxFyqtAlDQ8uFOUyt15VN1B12k+Z4K9oxC7TzUrkwDKVf1RxS\nFG5TZmJg4lWMQhUVVpf6h6HSrFQKhpXi0VH9LUECgYEA8ikNwNAjqUUwD0I5Znmn\n0noissqMM+k5mw49sp4qsyntKdBhb7ufb6zPo1SdEF6u0eutHIq/pPdFMGRnka61\nnuxOk8eDhsJ7LFeNXv7mb6QyMfp88ZDxV0EFTMBmpHc7yzEPOLrGNnh+RVMU+d9X\naJczYuSRmw1HlNtbCfnQrCECgYEAvmxkIB/8Jefi7ina169CttDzeSaJGueU91ID\nrPcNoCk7YnT5CIvjT0YeAwCkc9kRxgl5kwtf+wQN5YoAgCHn8fO2+zDa6zjXVHoh\nxanOVzHsEW/eGKvkllqZadD/9RC3Mf4E/GFxyKL4oQjdXxoUBRwnzS63uIW/JoNq\nzyV6CxkCgYBcR4S0KxzLzk/IIMZa5JUtQdmjJEhVJ9UJ311niZpf9+QmgQAAYhEZ\nr1LYvM+1gz8/Q34OWFk7dfbpv/kvrNINI6O18NuQKOBjP2HiB2SsundeUEP4kfFF\n/MMWQmNa3QzuG13fkl0iOLx1knl11sQqWSP91XgfC+pxMT36CTaZwQKBgQC2VNMK\nb5XgNcj0gt8o5ofaxPhcaKmfOV9J8R3T4DsLwG88NwS9SjS9E0ZpWZQd2RtLpIbk\nZV/h2l/0Cc+w4MZWxiXPH1h/Ik4MdWUg/xa0JvkDOTpQJUcbMGT1DUoIPZksJS5g\n+m0Yz/OBPhu5lB7XRb5WmQURif8dwXfkIN5bAQKBgQC5js/fgOZ1YQj4+eYT3W5f\nIKHcse0J5CyyiAxTuy3KlKNy5NXuw+DWmYfJrLIxlpk4HC689lRGlvOppGSDz9nD\nFRHW2kbk6XNPqUTLJd1KcvPRyN6IbeZeLXFTNcz6HIKnPCCB9aVaJ56hHk2NA/gl\n9Rqka4ldzBQBBH8JMY2qDA==\n-----END PRIVATE KEY-----\n
EOF
)

# Environment variables (non-secret)
FIREBASE_TYPE="service_account"
FIREBASE_PROJECT_ID="prince-group-50c85"
FIREBASE_PRIVATE_KEY_ID="547a6cfbd9537db359ab16cb8d43bdb87304697b"
FIREBASE_CLIENT_EMAIL="firebase-adminsdk-fbsvc@prince-group-50c85.iam.gserviceaccount.com"
FIREBASE_CLIENT_ID="115799228474237185692"
FIREBASE_AUTH_URI="https://accounts.google.com/o/oauth2/auth"
FIREBASE_TOKEN_URI="https://oauth2.googleapis.com/token"
FIREBASE_AUTH_PROVIDER_CERT_URL="https://www.googleapis.com/oauth2/v1/certs"
FIREBASE_CLIENT_CERT_URL="https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-fbsvc%40prince-group-50c85.iam.gserviceaccount.com"

# ====== Create container app ======

echo "Creating container app: $APP_NAME"

az containerapp create \
  --name "$APP_NAME" \
  --resource-group "$RESOURCE_GROUP" \
  --environment "$ENV_NAME" \
  --image "$ACR_SERVER/$IMAGE_NAME" \
  --target-port $PORT \
  --ingress external \
  --registry-server "$ACR_SERVER" \
  --registry-identity system \
  --secrets \
    database-url="$DATABASE_URL" \
    firebase-private-key="$FIREBASE_PRIVATE_KEY" \
  --env-vars \
    DATABASE_URL=secretref:database-url \
    FIREBASE_PRIVATE_KEY=secretref:firebase-private-key \
    FIREBASE_TYPE="$FIREBASE_TYPE" \
    FIREBASE_PROJECT_ID="$FIREBASE_PROJECT_ID" \
    FIREBASE_PRIVATE_KEY_ID="$FIREBASE_PRIVATE_KEY_ID" \
    FIREBASE_CLIENT_EMAIL="$FIREBASE_CLIENT_EMAIL" \
    FIREBASE_CLIENT_ID="$FIREBASE_CLIENT_ID" \
    FIREBASE_AUTH_URI="$FIREBASE_AUTH_URI" \
    FIREBASE_TOKEN_URI="$FIREBASE_TOKEN_URI" \
    FIREBASE_AUTH_PROVIDER_CERT_URL="$FIREBASE_AUTH_PROVIDER_CERT_URL" \
    FIREBASE_CLIENT_CERT_URL="$FIREBASE_CLIENT_CERT_URL"

echo "✅ Deployment triggered. Use the following to get FQDN:"
echo "az containerapp show --name $APP_NAME --resource-group $RESOURCE_GROUP --query properties.configuration.ingress.fqdn -o tsv"
