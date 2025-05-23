CREATE TABLE invoices
(
    id            int4 GENERATED BY DEFAULT AS IDENTITY NOT NULL,
    user_id       int4                                  NOT NULL,
    ticket_id     int4                                  NOT NULL,
    status_id     int4                                  NOT NULL,
    date_invoice  timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
    price         float4                                NOT NULL,
    status_change timestamptz DEFAULT CURRENT_TIMESTAMP NULL,
    CONSTRAINT invoices_pkey PRIMARY KEY (id)
);

CREATE INDEX idx_invoices_user_id ON invoices (user_id);
CREATE INDEX idx_invoices_ticket_id ON invoices (ticket_id);
