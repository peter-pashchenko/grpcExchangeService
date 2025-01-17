-- noinspection SqlDialectInspectionForFile

-- +goose Up
-- +goose StatementBegin
-- noinspection SqlDialectInspection

CREATE TABLE IF NOT EXISTS exchangeRate (
    requestTimeStamp BIGINT PRIMARY KEY,
    askPrice VARCHAR(255),
    bidPrice VARCHAR(255)
    );

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS exchangeRate;

-- +goose StatementEnd