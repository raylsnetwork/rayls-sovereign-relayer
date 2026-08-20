package testtools

import (
	"context"
	"log"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/testcontainers/testcontainers-go"
	natsmodule "github.com/testcontainers/testcontainers-go/modules/nats"
)

func SetupNATS(t testing.TB) (*nats.Conn, func()) {
	natsContainer, err := natsmodule.Run(context.Background(), "nats:2-alpine")
	if err != nil {
		t.Fatalf("failed to start NATS container: %s", err)
	}

	connStr, err := natsContainer.ConnectionString(context.Background())
	if err != nil {
		t.Fatalf("failed to get NATS connection string: %s", err)
	}

	nc, err := nats.Connect(connStr)
	if err != nil {
		t.Fatalf("failed to connect to NATS: %s", err)
	}

	return nc, func() {
		nc.Close()
		if err := testcontainers.TerminateContainer(natsContainer); err != nil {
			log.Printf("failed to terminate NATS container: %s", err)
		}
	}
}
