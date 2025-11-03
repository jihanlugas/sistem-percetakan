package constant

type RefTable string
type TransactionType int64
type PaymentType string

const (
	GENDER_MALE   = "MALE"
	GENDER_FEMALE = "FEMALE"

	// Form
	FORM_VARIANT_NO_ACTION string = ""
	FORM_VARIANT_NEW       string = "new"
	FORM_VARIANT_UPDATE    string = "update"
	FORM_VARIANT_DELETE    string = "delete"

	// Ref Table
	REF_TABLE_ITEM        RefTable = "item"
	REF_TABLE_ITEMVARIANT RefTable = "itemvariant"
	REF_TABLE_ADDON       RefTable = "addon"
	REF_TABLE_USER        RefTable = "user"

	TRANSACTION_TYPE_DEBIT  TransactionType = 1
	TRANSACTION_TYPE_KREDIT TransactionType = -1

	PAYMENT_TYPE_CASH     PaymentType = "CASH"
	PAYMENT_TYPE_TRANSFER PaymentType = "TRANSFER"
)
