-- +goose Up
-- +goose StatementBegin
CREATE TABLE feedbacks (
    id VARCHAR(255) PRIMARY KEY,
    customer_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    rating INT NOT NULL,
    comment TEXT NOT NULL,
    task_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

ALTER TABLE feedbacks ADD CONSTRAINT fk_feedbacks_customers FOREIGN KEY (customer_id) REFERENCES customers (max_id);
ALTER TABLE feedbacks ADD CONSTRAINT fk_feedbacks_users FOREIGN KEY (user_id) REFERENCES users (max_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE feedbacks;
-- +goose StatementEnd
