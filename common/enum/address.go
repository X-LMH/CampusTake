package enum

type AddressDefaultType int8

const (
	AddressNotDefault AddressDefaultType = 0
	AddressIsDefault  AddressDefaultType = 1
)

type AddressType int8

const (
	AddressTypeDelivery AddressType = 1
	AddressTypePickup   AddressType = 2
)
