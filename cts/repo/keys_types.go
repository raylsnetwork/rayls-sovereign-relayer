package repo

import "time"

type PublicRelayerRaylsSignKeysModel struct {
	Kind             string    `db:"kind"` // "public_relayer"
	PublicChainKeys  [][]byte  `db:"public_chain_keys"`
	PrivateChainKeys [][]byte  `db:"private_chain_keys"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}

type PrivateRelayerRaylsSignKeysModel struct {
	Kind                      string    `db:"kind"` // "private_relayer"
	PrivateHubKeys            [][]byte  `db:"private_hub_keys"`
	PrivateHubDvpOperatorKeys [][]byte  `db:"private_hub_dvp_operator_keys"`
	PrivateChainKeys          [][]byte  `db:"private_chain_keys"`
	CreatedAt                 time.Time `db:"created_at"`
	UpdatedAt                 time.Time `db:"updated_at"`
}

type AtomicServiceRaylsSignKeysModel struct {
	Kind             string    `db:"kind"` // "atomic_service"
	PrivateHubKeys   [][]byte  `db:"private_hub_keys"`
	PrivateChainKeys [][]byte  `db:"private_chain_keys"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}

type RaylsViewKeysModel struct {
	InitialBlock       uint64    `db:"initial_block"`
	EncryptedSecretKey []byte    `db:"encrypted_secret_key"`
	RaylsViewPublicKey []byte    `db:"public_key"`
	CreatedAt          time.Time `db:"created_at"`
	UpdatedAt          time.Time `db:"updated_at"`
}

type PaymentSpendKeysModel struct {
	ID        int       `db:"id"`
	SecretKey []byte    `db:"secret_key"`
	PublicKey []byte    `db:"public_key"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type SharedSecretModel struct {
	ChainId         string    `db:"chain_id"`
	EncryptedSecret []byte    `db:"encrypted_secret"`
	InitialBlock    uint64    `db:"initial_block"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}

type EnygmaSelfSecretModel struct {
	EncryptedSecret []byte    `db:"encrypted_secret"`
	InitialBlock    uint64    `db:"initial_block"`
	ResourceID      []byte    `db:"resource_id"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}
