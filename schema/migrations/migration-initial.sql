/* Storing loan information and state*/
CREATE TABLE IF NOT EXISTS loans (
    id UUID PRIMARY KEY,
    borrower_id_number TEXT NOT NULL,
    principal_amount DECIMAL(15,2) NOT NULL,
    rate DECIMAL(5,2) NOT NULL,
    roi DECIMAL(5,2) NOT NULL,
    state TEXT NOT NULL,
    approval_details JSONB,
    disbursement_details JSONB,
    agreement_letter_url TEXT,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
    );

CREATE INDEX IF NOT EXISTS idx_loans_state ON loans(state);
CREATE INDEX IF NOT EXISTS idx_loans_borrower ON loans(borrower_id_number);


/* Storing loan investments */
CREATE TABLE IF NOT EXISTS investments (
    id UUID PRIMARY KEY,
    loan_id UUID NOT NULL,
    investor_id UUID NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_investments_loan ON investments(loan_id);
CREATE INDEX IF NOT EXISTS idx_investments_investor ON investments(investor_id);

-- Inserting mock data

INSERT INTO loans (
    id, borrower_id_number, principal_amount, rate, roi, state,
    approval_details, disbursement_details, agreement_letter_url,
    created_at, updated_at
) VALUES
(
'550e8400-e29b-41d4-a716-446655440000', '123-45-6789', 10000.00, 5.00, 6.00, 'approved',
'{"approved_by": "Alice Johnson", "approval_date": "2023-01-15"}'::jsonb,
'{"disbursed_by": "Bob Smith", "disbursement_date": "2023-01-20"}'::jsonb,
'https://example.com/agreements/loan_550e8400-e29b-41d4-a716-446655440000.pdf',
'2023-01-15T10:00:00Z', '2023-01-20T10:00:00Z'
)

INSERT INTO investments (
id, loan_id, investor_id, amount, created_at
) VALUES
(
'860e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-446655440000', '550e8400-e29b-41d4-a716-446655440000', 5000.00, '2023-01-15T10:30:00Z'
)