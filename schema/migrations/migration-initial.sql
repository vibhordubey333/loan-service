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