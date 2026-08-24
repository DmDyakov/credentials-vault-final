//go:build integration

package integration

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"credentials-vault/client-cli/internal/client"
	"credentials-vault/client-cli/internal/config"
	authpb "credentials-vault/gen/go/auth/v1"
	vaultpb "credentials-vault/gen/go/vault/v1"
)

// mockAuthServer — моковый AuthService.
type mockAuthServer struct {
	authpb.UnimplementedAuthServiceServer
}

// TestClientFullFlow  — полный цикл клиента с моковым сервером.
func TestClientFullFlow(t *testing.T) {
	certFile, keyFile := generateTestCert(t)

	addr, cleanup := startMockServer(t, certFile, keyFile)
	defer cleanup()

	cfg := &config.Config{
		ServerAddress: addr,
		CAFile:        certFile,
	}

	cl, err := client.New(cfg)
	require.NoError(t, err)
	defer cl.Close()

	ctx := context.Background()
	username := "testuser"
	password := "TestPass123"

	t.Run("Register", func(t *testing.T) {
		err = cl.Register(ctx, username, password)
		require.NoError(t, err)
	})

	t.Run("AddCredentials", func(t *testing.T) {
		err = cl.AddCredentials(ctx, "example.com", "bob", "secret")
		require.NoError(t, err)
	})

	var itemID string
	t.Run("ListItems", func(t *testing.T) {
		items, listErr := cl.ListItems(ctx)
		require.NoError(t, listErr)
		require.NotEmpty(t, items)

		itemID = items[0].ID
		assert.Equal(t, "ITEM_TYPE_LOGIN", items[0].Type)
		assert.Equal(t, "example.com", items[0].Metadata["site"])
	})

	t.Run("GetItem_Decrypt", func(t *testing.T) {
		item, getErr := cl.GetItem(ctx, itemID)
		require.NoError(t, getErr)

		assert.Equal(t, "bob", item.Secret["username"])
		assert.Equal(t, "secret", item.Secret["password"])
		assert.Equal(t, "example.com", item.Metadata["site"])
	})
}

func (m *mockAuthServer) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	return authpb.RegisterResponse_builder{
		User: authpb.User_builder{
			Id:       proto.String("test-user-id"),
			Username: proto.String(req.GetUsername()),
		}.Build(),
		Message: proto.String("registered"),
	}.Build(), nil
}

func (m *mockAuthServer) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	header := metadata.Pairs("authorization", "Bearer test-token")
	grpc.SendHeader(ctx, header)

	return authpb.LoginResponse_builder{
		User: authpb.User_builder{
			Id:       proto.String("test-user-id"),
			Username: proto.String(req.GetUsername()),
		}.Build(),
		Salt: []byte("test-salt-1234567890123456"),
	}.Build(), nil
}

// mockVaultServer — моковый VaultService.
type mockVaultServer struct {
	vaultpb.UnimplementedVaultServiceServer
	items map[string]*vaultpb.VaultItem
}

func (m *mockVaultServer) CreateItem(ctx context.Context, req *vaultpb.CreateItemRequest) (*vaultpb.CreateItemResponse, error) {
	item := vaultpb.VaultItem_builder{
		Id:            proto.String("test-item-id"),
		Type:          req.GetType().Enum(), // ← Было req.Type
		EncryptedData: req.GetEncryptedData(),
		Metadata:      req.GetMetadata(),
		CreatedAt:     timestamppb.Now(),
		UpdatedAt:     timestamppb.Now(),
	}.Build()

	m.items[item.GetId()] = item

	return vaultpb.CreateItemResponse_builder{
		Item: item,
	}.Build(), nil
}

func (m *mockVaultServer) GetItem(ctx context.Context, req *vaultpb.GetItemRequest) (*vaultpb.GetItemResponse, error) {
	item, ok := m.items[req.GetId()]
	if !ok {
		return nil, status.Error(codes.NotFound, "item not found")
	}

	return vaultpb.GetItemResponse_builder{
		Item: item,
	}.Build(), nil
}

func (m *mockVaultServer) ListItems(ctx context.Context, req *vaultpb.ListItemsRequest) (*vaultpb.ListItemsResponse, error) {
	items := make([]*vaultpb.VaultItem, 0, len(m.items))
	for _, item := range m.items {
		items = append(items, item)
	}

	return vaultpb.ListItemsResponse_builder{
		Items: items,
	}.Build(), nil
}

// generateTestCert генерирует самоподписанный сертификат и ключ.
func generateTestCert(t *testing.T) (certFile, keyFile string) {
	t.Helper()

	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "localhost",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
		DNSNames:              []string{"localhost"},
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	require.NoError(t, err)

	tmpDir := t.TempDir()

	certFile = filepath.Join(tmpDir, "server.crt")
	certOut, err := os.Create(certFile)
	require.NoError(t, err)
	defer certOut.Close()
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	keyFile = filepath.Join(tmpDir, "server.key")
	keyOut, err := os.Create(keyFile)
	require.NoError(t, err)
	defer keyOut.Close()

	keyBytes, err := x509.MarshalECPrivateKey(priv)
	require.NoError(t, err)
	pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: keyBytes})

	return certFile, keyFile
}

// startMockServer запускает моковый TLS gRPC сервер.
func startMockServer(t *testing.T, certFile, keyFile string) (addr string, cleanup func()) {
	t.Helper()

	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	require.NoError(t, err)

	server := grpc.NewServer(grpc.Creds(creds))

	authpb.RegisterAuthServiceServer(server, &mockAuthServer{})
	vaultpb.RegisterVaultServiceServer(server, &mockVaultServer{
		items: make(map[string]*vaultpb.VaultItem),
	})

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	go server.Serve(listener)

	return listener.Addr().String(), func() {
		server.Stop()
	}
}
