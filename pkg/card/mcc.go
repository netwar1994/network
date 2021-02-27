package card

func TranslateMCC(code string) string {
	mcc := map[string]string{
		"5411": "Grocery Stores, Supermarkets",
		"5812": "Eating Places and Restaurants",
	}
	for key, _ := range mcc {
		if key == code {
			return mcc[key]
		}
	}
	return "Category not listed"
}