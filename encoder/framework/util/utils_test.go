package utils_test

import (
	utils "encoder/framework/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsJson(t *testing.T) {
	json := `{
		"id": "ba915b00-76ea-11eb-9439-0242ac130002",
		"file_path": "convite.mp4",
		"status": "pending"
	}`

	err := utils.IsJSON(json)
	require.Nil(t, err)

	json = "invalid"
	err = utils.IsJSON(json)
	require.Error(t, err)

}
