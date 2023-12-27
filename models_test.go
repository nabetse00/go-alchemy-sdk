package goalchemysdk

import "testing"

func TestAlchemyApiError_Error(t *testing.T) {
	tests := []struct {
		name string
		ae   *AlchemyApiError
		want string
	}{
		{
			name: "test error output",
			ae:   &AlchemyApiError{Code: 22, Message: "some error"},
			want: "Alchemy Api error code=22, message=some error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ae.Error(); got != tt.want {
				t.Errorf("AlchemyApiError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
