package migrate

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/raylsnetwork/rayls-sovereign-relayer/configinit"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cts/config"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cts/crypto"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cts/domain"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cts/repo"
	"github.com/raylsnetwork/rayls-sovereign-relayer/cts/service"
	"github.com/raylsnetwork/rayls-sovereign-relayer/logger"
	"github.com/raylsnetwork/rayls-sovereign-relayer/otel"
	"github.com/raylsnetwork/rayls-sovereign-relayer/withstack"
)

const (
	encryptorInitTimeout = 5 * time.Second
	migrationTimeout     = 5 * time.Minute
)

func Migrate(path, fromEncryptor, toEncryptor string) error {
	conf, err := configinit.InitConfig[config.Config](path)
	if err != nil {
		slog.Error("Failed to initialize config", slog.Any("error", err))
		return fmt.Errorf("initializing migration config: %w", err)
	}
	slog.Info("Configuration validated successfully")

	otelShutdown, err := otel.SetupOTelSDK(context.Background(), conf.OtelSDKDisabled)
	if err != nil {
		slog.Error("Failed to initlaize OTeL", slog.Any("error", err))
		return fmt.Errorf("initializing OpenTelemetry: %w", err)
	}
	defer func() { _ = otelShutdown(context.Background()) }()

	loggerShutdown, err := logger.InitLogger(conf.LogHandler, conf.LogLevel, conf.OtelSDKDisabled)
	if err != nil {
		slog.Error("Failed to initialize logger", slog.Any("error", err))
		return fmt.Errorf("initializing logger: %w", err)
	}
	defer func() { _ = loggerShutdown(context.Background()) }()

	db, err := repo.NewDatabaseConnection(conf.DatabaseConn)
	if err != nil {
		slog.Error("Failed to connect to the database", slog.Any("error", err))
		return withstack.Wrap(fmt.Errorf("connecting to migration database: %w", err))
	}

	defer db.Close()

	plaintextEncryptorFactory := crypto.NewPlaintextEncryptorFactory()
	awsEncryptorFactory := crypto.NewAWSEncryptorFactory(conf.AWSProfile, conf.AWSAlias)
	gcpEncryptorFactory := crypto.NewGCPEncryptorFactory(
		conf.GCPProject,
		conf.GCPLocation,
		conf.GCPKeyRing,
		conf.GCPCryptoKey,
	)

	ctx, cancel := context.WithTimeout(context.Background(), encryptorInitTimeout)
	defer cancel()

	fromEncryptorClient, cleanup, err := crypto.GetEncryptor(
		ctx,
		fromEncryptor,
		plaintextEncryptorFactory,
		awsEncryptorFactory,
		gcpEncryptorFactory,
	)
	if err != nil {
		slog.Error("Failed to create encryptor client", slog.Any("error", err))
		return fmt.Errorf("creating source encryptor client: %w", err)
	}
	defer func() { _ = cleanup() }()

	toEncryptorClient, cleanup, err := crypto.GetEncryptor(
		ctx,
		toEncryptor,
		plaintextEncryptorFactory,
		awsEncryptorFactory,
		gcpEncryptorFactory,
	)
	if err != nil {
		slog.Error("Failed to create encryptor client", slog.Any("error", err))
		return fmt.Errorf("creating target encryptor client: %w", err)
	}
	defer func() { _ = cleanup() }()

	slog.Info("Beginning to migrate View keys")
	err = migrateRaylsViewKeys(db, fromEncryptor, toEncryptor, fromEncryptorClient, toEncryptorClient)
	if err != nil {
		slog.Error("Failed to migrate View keys:", slog.Any("Migrator Error", err))
	} else {
		slog.Info("Migrated View keys successfully")
	}

	// TODO: Refer to TODO bellow
	// slog.Info("Beginning to migrate rayls sign keys")
	// err = migrateRaylsSignKeys(db, fromEncryptor, toEncryptor, fromEncryptorClient, toEncryptorClient)
	// if err != nil {
	// 	slog.Error("Failed to migrate rayls sign keys:", slog.Any("Migrator Error", err))
	// } else {
	// 	slog.Info("Migrated rayls sign keys successfully")
	// }

	slog.Info("Beginning to migrate payment spend keys")
	err = migratePaymentSpendKeys(db, fromEncryptor, toEncryptor, fromEncryptorClient, toEncryptorClient)
	if err != nil {
		slog.Error("Failed to migrate payment spend keys:", slog.Any("Migrator Error", err))
	} else {
		slog.Info("Migrated payment spend keys successfully")
	}

	return nil
}

func migrateRaylsViewKeys(
	db *repo.PostgresDB,
	fromEncryptor, toEncryptor string,
	fromEncryptorClient, toEncryptorClient service.Encryptor,
) error {
	fromRepository := db.NewRaylsViewKeysRepository(
		prefixedTableName(fromEncryptor, repo.RaylsViewKeysCollectionName),
	)

	toRepository := db.NewRaylsViewKeysRepository(
		prefixedTableName(toEncryptor, repo.RaylsViewKeysCollectionName),
	)

	migrateCtx, cancel := context.WithTimeout(context.Background(), migrationTimeout)
	defer cancel()

	return service.MigrateRepository[domain.EncryptedRaylsViewKeyPair, domain.RaylsViewKeyPair](
		migrateCtx,
		fromRepository,
		toRepository,
		fromEncryptorClient,
		toEncryptorClient,
	)
}

func migratePaymentSpendKeys(
	db *repo.PostgresDB,
	fromEncryptor, toEncryptor string,
	fromEncryptorClient, toEncryptorClient service.Encryptor,
) error {
	fromRepository := db.NewPaymentSpendKeysRepository(
		prefixedTableName(fromEncryptor, repo.PaymentSpendKeysCollectionName),
	)

	toRepository := db.NewPaymentSpendKeysRepository(
		prefixedTableName(toEncryptor, repo.PaymentSpendKeysCollectionName),
	)

	migrateCtx, cancel := context.WithTimeout(context.Background(), migrationTimeout)
	defer cancel()

	return service.MigrateRepository[domain.EncryptedPaymentSpendKeys, domain.PaymentSpendKeys](
		migrateCtx,
		fromRepository,
		toRepository,
		fromEncryptorClient,
		toEncryptorClient,
	)
}

func prefixedTableName(encrName, tableName string) string {
	return encrName + "_" + tableName
}
