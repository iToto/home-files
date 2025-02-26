package exchangeclient

type Client interface {
	GetBalance() error
	GetPosition() error
}
